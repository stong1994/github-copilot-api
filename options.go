package githubcopilotapi

// Option is an option for the Lingyi client.
type Option func(*Copilot)

func WithModel(model string) Option {
	return func(c *Copilot) {
		c.model = model
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
