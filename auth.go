package githubcopilotapi

import (
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

func (c *Copilot) withAuth() error {
	if c.githubOAuthToken == "" {
		return errors.New("no GitHub token found")
	}

	if c.copilotToken.Token == "" || (c.copilotToken.ExpiresAt <= time.Now().Unix()) {
		sessionID := uuid.NewString() + strconv.Itoa(int(time.Now().UnixNano()/1000))

		url := "https://api.github.com/copilot_internal/v2/token"
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}

		req.Header.Set("Authorization", "token "+c.githubOAuthToken)
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

		c.copilotToken = token
		c.sessionID = sessionID

	}

	return nil
}

type UserData struct {
	OAuthToken string `json:"oauth_token"`
}

func getOAuthTokenInLocal() string {
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
