package githubcopilotapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type Token struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}

type Copilot struct {
	githubToken string
	token       Token
	sessionID   string
	machineID   string
}

func NewCopilot() (*Copilot, error) {
	c := new(Copilot)
	c.githubToken = getCacheToken()
	err := c.withAuth()
	return c, err
}

func (c *Copilot) withAuth() error {
	if c.githubToken == "" {
		return errors.New("no GitHub token found")
	}

	if c.token.Token == "" || (c.token.ExpiresAt <= time.Now().Unix()) {
		sessionID := uuid.NewString() + strconv.Itoa(int(time.Now().UnixNano()/1000))

		url := "https://api.github.com/copilot_internal/v2/token"
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}

		req.Header.Set("Authorization", "token "+c.githubToken)
		req.Header.Set("Accept", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return errors.New("Failed to authenticate: " + resp.Status)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		var token Token
		err = json.Unmarshal(body, &token)
		if err != nil {
			return err
		}

		c.token = token
		c.sessionID = sessionID

	}

	return nil
}

type UserData struct {
	OAuthToken string `json:"oauth_token"`
}

func getCacheToken() string {
	// loading token from the environment only in GitHub Codespaces
	token := os.Getenv("GITHUB_TOKEN")
	codespaces := os.Getenv("CODESPACES")
	if token != "" && codespaces != "" {
		return token
	}

	// loading token from the file
	configPath := findConfigPath()

	// token can be sometimes in apps.json sometimes in hosts.json
	filePaths := []string{
		filepath.Join(configPath, "github-copilot", "hosts.json"),
		filepath.Join(configPath, "github-copilot", "apps.json"),
	}

	for _, filePath := range filePaths {
		if _, err := os.Stat(filePath); err == nil {
			file, err := os.ReadFile(filePath)
			if err != nil {
				continue
			}

			var userData map[string]UserData
			err = json.Unmarshal(file, &userData)
			if err != nil {
				continue
			}

			for _, value := range userData {
				if value.OAuthToken != "" {
					return value.OAuthToken
				}
			}
		}
	}

	return ""
}

type EmbedInput struct {
	Filename string `json:"filename"`
	Filetype string `json:"filetype"`
	Prompt   string `json:"prompt,omitempty"`
	Content  string `json:"content,omitempty"`
}

type EmbedOpts struct {
	Model     string
	ChunkSize int
	OnDone    func([]EmbedInput)
	OnError   func(string)
}

func (c *Copilot) Embed(inputs []EmbedInput, opts EmbedOpts) {
	if len(inputs) == 0 {
		if opts.OnDone != nil {
			opts.OnDone([]EmbedInput{})
		}
		return
	}

	url := "https://api.githubcopilot.com/embeddings"

	chunks := make([][]EmbedInput, 0)
	for i := 0; i < len(inputs); i += opts.ChunkSize {
		end := i + opts.ChunkSize
		if end > len(inputs) {
			end = len(inputs)
		}
		chunks = append(chunks, inputs[i:end])
	}

	for _, chunk := range chunks {
		body, err := json.Marshal(chunk)
		if err != nil {
			if opts.OnError != nil {
				opts.OnError("Failed to encode request body: " + err.Error())
			}
			return
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
		if err != nil {
			if opts.OnError != nil {
				opts.OnError("Failed to create request: " + err.Error())
			}
			return
		}

		req.Header = generateHeaders(c.token.Token, c.sessionID, c.machineID)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			if opts.OnError != nil {
				opts.OnError("Failed to make request: " + err.Error())
			}
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			if opts.OnError != nil {
				opts.OnError("Failed to get response: " + resp.Status)
			}
			return
		}

		var result []EmbedInput
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			if opts.OnError != nil {
				opts.OnError("Failed to decode response: " + err.Error())
			}
			return
		}

		if opts.OnDone != nil {
			opts.OnDone(result)
		}
	}
}

func generateHeaders(token, sessionID, machineID string) http.Header {
	headers := http.Header{
		"Authorization":          []string{"Bearer " + token},
		"x-request-id":           []string{uuid.NewString()},
		"vscode-sessionid":       []string{sessionID},
		"vscode-machineid":       []string{machineID},
		"copilot-integration-id": []string{"vscode-chat"},
		"openai-organization":    []string{"github-copilot"},
		"openai-intent":          []string{"conversation-panel"},
		"Content-Type":           []string{"application/json"},
	}

	versionHeaders := map[string]string{
		"client-version": "1.0.0",
		"user-agent":     "github.com/stong1994/github-copilot-api/1.0.0",
	}

	for key, value := range versionHeaders {
		headers.Add(key, value)
	}

	return headers
}

type CompletionOpts struct {
	Model       string
	Temperature float64
	Stream      bool
	TopP        int
	N           int
}

type completionRequest struct {
	Intent      bool      `json:"intent"`
	Model       string    `json:"model"`
	N           int       `json:"n"`
	Stream      bool      `json:"stream"`
	Temperature float64   `json:"temperature"`
	TopP        int       `json:"top_p"`
	Messages    []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type CompletionResponse struct {
	Choices []struct {
		ContentFilterResults struct {
			Hate struct {
				Filtered bool   `json:"filtered,omitempty"`
				Severity string `json:"severity,omitempty"`
			} `json:"hate,omitempty"`
			SelfHarm struct {
				Filtered bool   `json:"filtered,omitempty"`
				Severity string `json:"severity,omitempty"`
			} `json:"self_harm,omitempty"`
			Sexual struct {
				Filtered bool   `json:"filtered,omitempty"`
				Severity string `json:"severity,omitempty"`
			} `json:"sexual,omitempty"`
			Violence struct {
				Filtered bool   `json:"filtered,omitempty"`
				Severity string `json:"severity,omitempty"`
			} `json:"violence,omitempty"`
		} `json:"content_filter_results,omitempty"`
		FinishReason string  `json:"finish_reason,omitempty"`
		Index        float64 `json:"index,omitempty"`
		Message      struct {
			Content string `json:"content,omitempty"`
			Role    string `json:"role,omitempty"`
		} `json:"message,omitempty"`
	} `json:"choices,omitempty"`
	Created             float64 `json:"created,omitempty"`
	ID                  string  `json:"id,omitempty"`
	Model               string  `json:"model,omitempty"`
	PromptFilterResults []struct {
		ContentFilterResults struct {
			Hate struct {
				Filtered bool   `json:"filtered,omitempty"`
				Severity string `json:"severity,omitempty"`
			} `json:"hate,omitempty"`
			SelfHarm struct {
				Filtered bool   `json:"filtered,omitempty"`
				Severity string `json:"severity,omitempty"`
			} `json:"self_harm,omitempty"`
			Sexual struct {
				Filtered bool   `json:"filtered,omitempty"`
				Severity string `json:"severity,omitempty"`
			} `json:"sexual,omitempty"`
			Violence struct {
				Filtered bool   `json:"filtered,omitempty"`
				Severity string `json:"severity,omitempty"`
			} `json:"violence,omitempty"`
		} `json:"content_filter_results,omitempty"`
		PromptIndex float64 `json:"prompt_index,omitempty"`
	} `json:"prompt_filter_results,omitempty"`
	Usage struct {
		CompletionTokens float64 `json:"completion_tokens,omitempty"`
		PromptTokens     float64 `json:"prompt_tokens,omitempty"`
		TotalTokens      float64 `json:"total_tokens,omitempty"`
	} `json:"usage,omitempty"`
}

func (c *Copilot) Complete(systemPrompt, userAsk string, opts CompletionOpts) (CompletionResponse, error) {
	url := "https://api.githubcopilot.com/chat/completions"
	cr := completionRequest{
		Intent:      true,
		Model:       opts.Model,
		N:           opts.N,
		Stream:      opts.Stream,
		TopP:        opts.TopP,
		Temperature: opts.Temperature,
	}

	if systemPrompt != "" {
		cr.Messages = []Message{{Content: systemPrompt, Role: "system"}}
	}
	cr.Messages = append(cr.Messages, Message{Content: userAsk, Role: "user"})
	body, err := json.Marshal(cr)
	if err != nil {
		return CompletionResponse{}, err
	}

	if err = c.withAuth(); err != nil {
		return CompletionResponse{}, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return CompletionResponse{}, err
	}

	req.Header = generateHeaders(c.token.Token, c.sessionID, c.machineID)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return CompletionResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return CompletionResponse{}, errors.New("Failed to complete: " + resp.Status)
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return CompletionResponse{}, err
	}

	var response CompletionResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return CompletionResponse{}, err
	}
	return response, nil
}

func findConfigPath() string {
	// Check XDG_CONFIG_HOME first
	config := os.Getenv("XDG_CONFIG_HOME")
	if config != "" && isDirectory(config) {
		return config
	}

	// On Windows, check AppData/Local
	if os.Getenv("OS") == "Windows_NT" {
		config = filepath.Join(os.Getenv("APPDATA"), "Local")
		if isDirectory(config) {
			return config
		}
	}

	// Fallback to ~/.config
	config = filepath.Join(os.Getenv("HOME"), ".config")
	if isDirectory(config) {
		return config
	}

	// If all else fails, return an empty string
	return ""
}

func isDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
