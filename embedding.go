package githubcopilotapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type EmbeddingRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

type EmbeddingResponse struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
		Index     int       `json:"index"`
	} `json:"data"`
	Usage struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
}

func (c *Copilot) setEmbeddingDefaults(payload *EmbeddingRequest) {
	if payload.Model == "" {
		payload.Model = c.embeddingModel
	}
}

func (c *Copilot) CreateEmbedding(ctx context.Context, payload *EmbeddingRequest) (*EmbeddingResponse, error) {
	if err := c.withAuth(); err != nil {
		return nil, err
	}
	c.setEmbeddingDefaults(payload)

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// Build request
	body := bytes.NewReader(payloadBytes)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/embeddings", c.baseURL), body)
	if err != nil {
		return nil, err
	}

	c.setHeaders(req)

	// Send request
	r, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned unexpected status code: %d", r.StatusCode)
	}
	// Parse response
	var response EmbeddingResponse
	return &response, json.NewDecoder(r.Body).Decode(&response)
}
