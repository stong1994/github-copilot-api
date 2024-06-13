package githubcopilotapi

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ResponseMessage struct {
	Role    *string `json:"role,omitempty"`
	Content string  `json:"content"`
}

// CompletionRequest is a request to complete a completion.
type CompletionRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
	N           int       `json:"n,omitempty"`
	TopP        float64   `json:"top_p,omitempty"`
	Stream      bool      `json:"stream"`

	// StreamingFunc is a function to be called for each chunk of a streaming response.
	// Return an error to stop streaming early.
	StreamingFunc func(ctx context.Context, chunk []byte) error `json:"-"`
}

type CompletionResponseChoice struct {
	ContentFilterResults ContentFilterResults `json:"content_filter_results,omitempty"`
	FinishReason         string               `json:"finish_reason,omitempty"`
	Index                float64              `json:"index,omitempty"`
	Message              ResponseMessage      `json:"message,omitempty"`
}
type ContentFilterResults struct {
	Error struct {
		Code    string `json:"code,omitempty"`
		Message string `json:"message,omitempty"`
	} `json:"error,omitempty"`
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
}
type CompletionResponse struct {
	Choices             []CompletionResponseChoice `json:"choices,omitempty"`
	Created             float64                    `json:"created,omitempty"`
	ID                  string                     `json:"id,omitempty"`
	Model               string                     `json:"model,omitempty"`
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

// StreamedChatResponsePayload is a chunk from the stream.
type StreamedChatResponsePayload struct {
	Choices []struct {
		ContentFilterOffsets struct {
			CheckOffset float64 `json:"check_offset,omitempty"`
			EndOffset   float64 `json:"end_offset,omitempty"`
			StartOffset float64 `json:"start_offset,omitempty"`
		} `json:"content_filter_offsets,omitempty"`
		ContentFilterResults ContentFilterResults `json:"content_filter_results,omitempty"`
		Delta                ResponseMessage      `json:"delta,omitempty"`
		Index                float64              `json:"index,omitempty"`
		FinishReason         string               `json:"finish_reason,omitempty"`
	} `json:"choices,omitempty"`
	Created float64 `json:"created,omitempty"`
	ID      string  `json:"id,omitempty"`
	Error   error   `json:"-"`
}

type errorMessage struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Param   any    `json:"param,omitempty"`
		Type    string `json:"type,omitempty"`
	} `json:"error"`
}

func (c *Copilot) setCompletionDefaults(payload *CompletionRequest) {
	// Set defaults
	if payload.N == 0 {
		payload.N = 1
	}
	if payload.TopP == 0 {
		payload.TopP = 1
	}
	if payload.Temperature == 0 {
		payload.Temperature = 0.1
	}

	switch {
	// Prefer the model specified in the payload.
	case payload.Model != "":

	// If no model is set in the payload, take the one specified in the client.
	case c.model != "":
		payload.Model = c.model
	// Fallback: use the default model
	default:
		payload.Model = defaultCompletionModel
	}
	if c.httpClient == nil {
		c.httpClient = http.DefaultClient
	}
}

func (c *Copilot) CreateCompletion(ctx context.Context, payload *CompletionRequest) (*CompletionResponse, error) {
	if err := c.withAuth(); err != nil {
		return nil, err
	}
	c.setCompletionDefaults(payload)
	if payload.StreamingFunc != nil {
		payload.Stream = true
	}
	// Build request payload

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// Build request
	body := bytes.NewReader(payloadBytes)
	if c.baseURL == "" {
		c.baseURL = defaultBaseURL
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/chat/completions", c.baseURL), body)
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
	if payload.StreamingFunc != nil {
		return parseStreamingChatResponse(ctx, r, payload)
	}
	// Parse response
	var response CompletionResponse
	return &response, json.NewDecoder(r.Body).Decode(&response)
}

func parseStreamingChatResponse(
	ctx context.Context,
	r *http.Response,
	payload *CompletionRequest,
) (*CompletionResponse, error) {
	scanner := bufio.NewScanner(r.Body)
	responseChan := make(chan StreamedChatResponsePayload)
	go func() {
		defer close(responseChan)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				continue
			}

			data := strings.TrimPrefix(line, "data:") // here use `data:` instead of `data: ` for compatibility
			data = strings.TrimSpace(data)
			if data == "[DONE]" {
				return
			}
			var streamPayload StreamedChatResponsePayload
			err := json.NewDecoder(bytes.NewReader([]byte(data))).Decode(&streamPayload)
			if err != nil {
				streamPayload.Error = fmt.Errorf("error decoding streaming response: %w", err)
				responseChan <- streamPayload
				return
			}
			responseChan <- streamPayload
		}
		if err := scanner.Err(); err != nil {
			responseChan <- StreamedChatResponsePayload{Error: fmt.Errorf("error reading streaming response: %w", err)}
			return
		}
	}()
	// Combine response
	return combineStreamingChatResponse(ctx, payload, responseChan)
}

func combineStreamingChatResponse(
	ctx context.Context,
	payload *CompletionRequest,
	responseChan chan StreamedChatResponsePayload,
) (*CompletionResponse, error) {
	response := CompletionResponse{
		Choices: []CompletionResponseChoice{
			{},
		},
	}

	for streamResponse := range responseChan {
		if streamResponse.Error != nil {
			return nil, streamResponse.Error
		}

		if len(streamResponse.Choices) == 0 {
			continue
		}
		choice := streamResponse.Choices[0]
		chunk := []byte(choice.Delta.Content)
		response.Choices[0].Message.Content += choice.Delta.Content
		response.Choices[0].FinishReason = choice.FinishReason
		response.Choices[0].ContentFilterResults = choice.ContentFilterResults

		if payload.StreamingFunc != nil {
			err := payload.StreamingFunc(ctx, chunk)
			if err != nil {
				return nil, fmt.Errorf("streaming func returned an error: %w", err)
			}
		}
	}
	return &response, nil
}
