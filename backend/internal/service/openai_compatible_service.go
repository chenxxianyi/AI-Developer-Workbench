package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"ai-developer-workbench/internal/config"
)

// OpenAICompatibleService calls OpenAI-compatible Chat Completions API.
type OpenAICompatibleService struct {
	baseURL     string
	apiKey      string
	model       string
	visionModel string
	httpClient  *http.Client
	maxRetries  int
}

// NewOpenAICompatibleService creates a new OpenAI-compatible service.
func NewOpenAICompatibleService(cfg *config.AIConfig) *OpenAICompatibleService {
	timeout := time.Duration(cfg.TimeoutSeconds) * time.Second
	return &OpenAICompatibleService{
		baseURL:     cfg.BaseURL,
		apiKey:      cfg.APIKey,
		model:       cfg.Model,
		visionModel: cfg.VisionModel,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		maxRetries: cfg.MaxRetries,
	}
}

// chatRequest represents an OpenAI Chat Completions request.
type chatRequest struct {
	Model          string          `json:"model"`
	Messages       []chatMessage   `json:"messages"`
	Temperature    float64         `json:"temperature,omitempty"`
	ResponseFormat *responseFormat `json:"response_format,omitempty"`
}

type chatMessage struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"`
}

type responseFormat struct {
	Type string `json:"type"`
}

// chatResponse represents an OpenAI Chat Completions response.
type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error,omitempty"`
}

// GenerateJSON calls the OpenAI-compatible API and returns the response.
func (s *OpenAICompatibleService) GenerateJSON(ctx context.Context, input AIRequest) (*AIResult, error) {
	model := s.model
	messages := []chatMessage{
		{Role: "system", Content: input.SystemPrompt},
	}

	// Build user message content.
	imagePaths := input.ImagePaths
	if len(imagePaths) == 0 && input.ImagePath != "" {
		imagePaths = []string{input.ImagePath}
	}
	if input.NeedVision && len(imagePaths) > 0 {
		model = s.visionModel
		imageContent, err := s.buildVisionContent(input.UserPrompt, imagePaths)
		if err != nil {
			return nil, fmt.Errorf("failed to build vision content: %w", err)
		}
		messages = append(messages, chatMessage{Role: "user", Content: imageContent})
	} else {
		messages = append(messages, chatMessage{Role: "user", Content: input.UserPrompt})
	}

	reqBody := chatRequest{
		Model:          model,
		Messages:       messages,
		Temperature:    0.3,
		ResponseFormat: &responseFormat{Type: "json_object"},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	var lastErr error
	attempts := s.maxRetries + 1
	for i := 0; i < attempts; i++ {
		result, err := s.doRequest(ctx, bodyBytes, model)
		if err == nil {
			return result, nil
		}
		lastErr = err
		if !isRetryableError(err) {
			break
		}
		slog.Warn("Retrying AI request", "attempt", i+1, "error", err)
		select {
		case <-time.After(time.Second):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	return nil, fmt.Errorf("AI request failed after %d attempts: %w", attempts, lastErr)
}

// doRequest sends a single request to the API.
func (s *OpenAICompatibleService) doRequest(ctx context.Context, body []byte, model string) (*AIResult, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, &retryableError{fmt.Errorf("rate limited (429)")}
	}

	if resp.StatusCode >= 500 {
		return nil, &retryableError{fmt.Errorf("server error (%d)", resp.StatusCode)}
	}

	if resp.StatusCode != http.StatusOK {
		var errResp chatResponse
		if json.Unmarshal(respBody, &errResp) == nil && errResp.Error != nil {
			return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, errResp.Error.Message)
		}
		return nil, fmt.Errorf("API error: status %d", resp.StatusCode)
	}

	var chatResp chatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	content := chatResp.Choices[0].Message.Content

	return &AIResult{
		RawText:  content,
		JSONText: content,
		Provider: "openai_compatible",
		Model:    model,
	}, nil
}

// buildVisionContent builds the message content for vision requests.
func (s *OpenAICompatibleService) buildVisionContent(prompt string, imagePaths []string) ([]interface{}, error) {
	content := []interface{}{map[string]interface{}{"type": "text", "text": prompt}}
	for _, imagePath := range imagePaths {
		imgData, err := os.ReadFile(imagePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read image: %w", err)
		}
		if len(imgData) > 20*1024*1024 {
			return nil, fmt.Errorf("image too large for vision API")
		}
		mimeType := "image/png"
		switch strings.ToLower(filepath.Ext(imagePath)) {
		case ".jpg", ".jpeg":
			mimeType = "image/jpeg"
		case ".webp":
			mimeType = "image/webp"
		}
		dataURL := fmt.Sprintf("data:%s;base64,%s", mimeType, base64.StdEncoding.EncodeToString(imgData))
		content = append(content, map[string]interface{}{
			"type":      "image_url",
			"image_url": map[string]interface{}{"url": dataURL},
		})
	}
	return content, nil
}

// retryableError indicates an error that can be retried.
type retryableError struct {
	error
}

// isRetryableError checks if an error is retryable.
func isRetryableError(err error) bool {
	_, ok := err.(*retryableError)
	return ok
}
