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

func WithGithubToken(githubtoken string) Option {
	return func(c *Copilot) {
		c.githubToken = githubtoken
	}
}

func WithHTTPCopilot(httpclient Doer) Option {
	return func(c *Copilot) {
		c.httpClient = httpclient
	}
}
