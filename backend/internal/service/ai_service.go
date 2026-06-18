package service

import "context"

// AIRequest is the input to an AI service call.
type AIRequest struct {
	ToolType     string
	SystemPrompt string
	UserPrompt   string
	ImagePath    string // Optional: path to image file for vision
	NeedVision   bool   // Whether to use vision model
}

// AIResult is the output from an AI service call.
type AIResult struct {
	RawText  string // Full raw response text
	JSONText string // Extracted JSON text
	Provider string // Provider name
	Model    string // Model used
}

// AIService defines the interface for AI providers.
type AIService interface {
	GenerateJSON(ctx context.Context, input AIRequest) (*AIResult, error)
}