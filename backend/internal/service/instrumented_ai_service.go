package service

import (
	"context"
	"time"

	"ai-developer-workbench/internal/model"
	"ai-developer-workbench/internal/repository"
)

// instrumentedAIService wraps AIService with observability recording.
type instrumentedAIService struct {
	inner    AIService
	aiRunRepo repository.AIRunRepository
}

// NewInstrumentedAIService creates an AI service that records every call.
func NewInstrumentedAIService(inner AIService, aiRunRepo repository.AIRunRepository) AIService {
	return &instrumentedAIService{inner: inner, aiRunRepo: aiRunRepo}
}

func (s *instrumentedAIService) GenerateJSON(ctx context.Context, input AIRequest) (*AIResult, error) {
	start := time.Now()
	result, err := s.inner.GenerateJSON(ctx, input)
	durationMs := time.Since(start).Milliseconds()

	run := &model.AIRun{
		ReportID:   "", // caller must set
		ToolType:   input.ToolType,
		Provider:   "unknown",
		Model:      "unknown",
		IsMock:     false,
		DurationMs: durationMs,
	}

	if result != nil {
		run.Provider = result.Provider
		run.Model = result.Model
	}

	if err != nil {
		run.ParseSuccess = false
		run.ErrorType = classifyAIError(err)
	}

	// Best-effort write; don't fail the main request.
	_ = s.aiRunRepo.Create(ctx, run)

	return result, err
}

func classifyAIError(err error) string {
	if err == nil {
		return ""
	}
	msg := err.Error()
	// Simple classification by error message patterns.
	if contains(msg, "timeout") || contains(msg, "deadline exceeded") {
		return model.AIErrorTimeout
	}
	if contains(msg, "429") || contains(msg, "rate limit") {
		return model.AIErrorRateLimit
	}
	if contains(msg, "cancel") {
		return model.AIErrorCanceled
	}
	if contains(msg, "json") || contains(msg, "parse") || contains(msg, "unmarshal") {
		return model.AIErrorInvalidJSON
	}
	if contains(msg, "status 4") {
		return model.AIErrorProvider4xx
	}
	if contains(msg, "status 5") {
		return model.AIErrorProvider5xx
	}
	return model.AIErrorInternal
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchSubstring(s, substr)
}

func searchSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Ensure instrumentedAIService satisfies AIService.
var _ AIService = (*instrumentedAIService)(nil)
