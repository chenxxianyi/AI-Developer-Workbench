package dto

// ObservabilityDTO is the API response for AI call analytics.
type ObservabilityDTO struct {
	TotalCalls    int64             `json:"total_calls"`
	SuccessRate   float64           `json:"success_rate"`
	FallbackRate  float64           `json:"fallback_rate"`
	ParseFailRate float64           `json:"parse_fail_rate"`
	AvgDurationMs float64           `json:"avg_duration_ms"`
	P50DurationMs float64           `json:"p50_duration_ms"`
	P95DurationMs float64           `json:"p95_duration_ms"`
	RetryRate     float64           `json:"retry_rate"`
	ByTool        []ToolAIStatDTO   `json:"by_tool,omitempty"`
	ByModel       []ModelAIStatDTO  `json:"by_model,omitempty"`
}

type ToolAIStatDTO struct {
	ToolType    string  `json:"tool_type"`
	TotalCalls  int64   `json:"total_calls"`
	SuccessRate float64 `json:"success_rate"`
}

type ModelAIStatDTO struct {
	Model         string  `json:"model"`
	TotalCalls    int64   `json:"total_calls"`
	AvgDurationMs float64 `json:"avg_duration_ms"`
}
