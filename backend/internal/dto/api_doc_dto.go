package dto

// APIDocRequest is the input for API Doc Builder.
type APIDocRequest struct {
	Title          string `json:"title" binding:"required"`
	SourceType     string `json:"source_type" binding:"required"`
	BackendStack   string `json:"backend_stack"`
	Code           string `json:"code"`
	APIDescription string `json:"api_description"`
	OutputFormat   string `json:"output_format" binding:"required"`
}

// APIDocResult is the output for API Doc Builder.
type APIDocResult struct {
	Modules         []ModuleItem  `json:"modules"`
	MarkdownContent *string       `json:"markdown_content,omitempty"`
	OpenAPIContent  *string       `json:"openapi_content,omitempty"`
	Recommendations []string      `json:"recommendations"`
	CodexPrompt     string        `json:"codex_prompt"`
}

// ModuleItem represents an API module.
type ModuleItem struct {
	Name      string       `json:"name"`
	Endpoints []EndpointItem `json:"endpoints"`
}

// EndpointItem represents an API endpoint.
type EndpointItem struct {
	Method      string `json:"method"`
	Path        string `json:"path"`
	Description string `json:"description"`
}