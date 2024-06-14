package githubcopilotapi

import (
	"net/http"

	"github.com/google/uuid"
)

const (
	defaultCompletionModel = "gpt-4"
	defaultEmbeddingModel  = "copilot-text-embedding-ada-002"
	defaultBaseURL         = "https://api.githubcopilot.com"
)

// Doer performs a HTTP request.
type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

type Copilot struct {
	githubOAuthToken string
	copilotToken     Token
	sessionID        string
	machineID        string
	userAgent        string

	model      string
	baseURL    string
	httpClient Doer
}

func NewCopilot(opts ...Option) (*Copilot, error) {
	c := new(Copilot)
	for _, opt := range opts {
		opt(c)
	}
	if c.githubOAuthToken == "" {
		c.githubOAuthToken = getOAuthTokenInLocal()
	}
	err := c.withAuth()
	return c, err
}

func (c *Copilot) setHeaders(req *http.Request) {
	req.Header = generateHeaders(c.copilotToken.Token, c.sessionID, c.machineID)
}

// TODO: set with option
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
