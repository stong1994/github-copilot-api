package githubcopilotapi

import (
	"net/http"

	"github.com/google/uuid"
)

const (
	defaultCompletionModel      = "gpt-4"
	defaultEmbeddingModel       = "copilot-text-embedding-ada-002"
	defaultBaseURL              = "https://api.githubcopilot.com"
	defaultCopilotIntegrationID = "vscode-chat"
	defaultOpenAIOrganization   = "github-copilot"
	defaultOpenAIIntent         = "conversation-panel"
	defaultUserAgent            = "github.com/stong1994/github-copilot-api/1.0.0"
	defaultClientVersion        = "1.0.0"
)

// Doer performs a HTTP request.
type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

type Copilot struct {
	githubOAuthToken string
	copilotToken     Token

	sessionID string
	machineID string

	completionModel      string
	embeddingModel       string
	userAgent            string
	copilotintegrationID string
	openaiOrganization   string
	openaiIntent         string
	clientVersion        string

	baseURL    string
	httpClient Doer
}

func NewCopilot(opts ...Option) (*Copilot, error) {
	c := new(Copilot)
	for _, opt := range opts {
		opt(c)
	}
	err := c.withAuth()
	return c, err
}

func (c *Copilot) setDefault() {
	if c.completionModel == "" {
		c.completionModel = defaultCompletionModel
	}
	if c.embeddingModel == "" {
		c.embeddingModel = defaultEmbeddingModel
	}
	if c.baseURL == "" {
		c.baseURL = defaultBaseURL
	}
	if c.copilotintegrationID == "" {
		c.copilotintegrationID = defaultCopilotIntegrationID
	}
	if c.openaiOrganization == "" {
		c.openaiOrganization = defaultOpenAIOrganization
	}
	if c.openaiIntent == "" {
		c.openaiIntent = defaultOpenAIIntent
	}
	if c.userAgent == "" {
		c.userAgent = defaultUserAgent
	}
	if c.githubOAuthToken == "" {
		c.githubOAuthToken = getOAuthTokenInLocal()
	}
	if c.httpClient == nil {
		c.httpClient = http.DefaultClient
	}
	if c.clientVersion == "" {
		c.clientVersion = defaultClientVersion
	}
}

func (c *Copilot) setHeaders(req *http.Request) {
	req.Header = c.generateHeaders(c.copilotToken.Token, c.sessionID, c.machineID)
}

// TODO: set with option
func (c *Copilot) generateHeaders(token, sessionID, machineID string) http.Header {
	headers := http.Header{
		"Authorization":          []string{"Bearer " + token},
		"x-request-id":           []string{uuid.NewString()},
		"vscode-sessionid":       []string{sessionID},
		"vscode-machineid":       []string{machineID},
		"copilot-integration-id": []string{c.copilotintegrationID},
		"openai-organization":    []string{c.openaiOrganization},
		"openai-intent":          []string{c.openaiIntent},
		"Content-Type":           []string{"application/json"},
	}

	versionHeaders := map[string]string{
		"client-version": c.clientVersion,
		"user-agent":     c.userAgent,
	}

	for key, value := range versionHeaders {
		headers.Add(key, value)
	}

	return headers
}
