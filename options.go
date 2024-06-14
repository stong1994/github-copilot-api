package githubcopilotapi

// Option is an option for the Lingyi client.
type Option func(*Copilot)

func WithCompletionModel(model string) Option {
	return func(c *Copilot) {
		c.completionModel = model
	}
}

func WithEmbeddingModel(model string) Option {
	return func(c *Copilot) {
		c.embeddingModel = model
	}
}

func WithBaseURL(baseURL string) Option {
	return func(c *Copilot) {
		c.baseURL = baseURL
	}
}

func WithGithubOAuthToken(githubOAuthToken string) Option {
	return func(c *Copilot) {
		c.githubOAuthToken = githubOAuthToken
	}
}

func WithHTTPCopilot(httpclient Doer) Option {
	return func(c *Copilot) {
		c.httpClient = httpclient
	}
}

func WithCopilotIntegrationID(copilotintegrationID string) Option {
	return func(c *Copilot) {
		c.copilotintegrationID = copilotintegrationID
	}
}

func WithOpenAIOrganization(openaiOrganization string) Option {
	return func(c *Copilot) {
		c.openaiOrganization = openaiOrganization
	}
}

func WithOpenAIIntent(openaiIntent string) Option {
	return func(c *Copilot) {
		c.openaiIntent = openaiIntent
	}
}

func WithUserAgent(userAgent string) Option {
	return func(c *Copilot) {
		c.userAgent = userAgent
	}
}

func WithClientVersion(version string) Option {
	return func(c *Copilot) {
		c.clientVersion = version
	}
}
