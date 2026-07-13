package prompts

// APIDocPromptSchema describes the expected output format.
const APIDocPromptSchema = `
Output a JSON object with this exact structure:
{
  "modules": [
    {
      "name": "string - module name",
      "endpoints": [
        {"method": "GET|POST|PUT|DELETE|PATCH", "path": "string", "description": "string"}
      ]
    }
  ],
  "markdown_content": "string - full Markdown documentation",
  "openapi_content": "string - OpenAPI 3.0 JSON (if requested)",
  "curl_examples": ["string - curl commands for each endpoint"],
  "frontend_guide": "string - frontend API client guide with code examples (optional)",
  "documentation_gaps": ["string - missing docs, unclear auth, inconsistent error codes"],
  "recommendations": ["string - recommendations"],
  "action_items": [
    {"id": "stable-kebab-case-id", "title": "string", "priority": "high|medium|low", "effort": "small|medium|large", "category": "string", "reason": "string", "suggested_prompt": "string", "issue_title": "string", "issue_body": "string - Markdown"}
  ],
  "codex_prompt": "string - prompt for generating implementation stubs"
}
` + ActionItemsPromptSchema

// BuildAPIDocPrompt creates the prompt for API Doc Builder.
func BuildAPIDocPrompt(sourceType, backendStack, code, apiDescription, outputFormat string, projectSummary string) (string, string) {
	systemPrompt := `You are an expert API documentation writer.
Your task is to analyze code or descriptions and generate comprehensive API documentation.

Required outputs:
- curl_examples: at least one curl command per module showing request format
- frontend_guide: JavaScript/TypeScript code snippets for common API calls
- documentation_gaps: list any missing information (auth method unclear, error codes undocumented, DTO fields missing)

Output format requirements:
- If output_format is "markdown": only generate markdown_content
- If output_format is "openapi": only generate openapi_content (must be valid OpenAPI 3.0 JSON)
- If output_format is "both": generate both fields

IMPORTANT: openapi_content must pass basic OpenAPI 3.0 structural validation.
Include "openapi": "3.0.0" at the root, with info, paths, and components sections.
` + APIDocPromptSchema

	userPrompt := `Documentation Request:
- Source Type: ` + sourceType + ` (code, project_zip, or manual)
- Backend Stack: ` + backendStack + `
- Output Format: ` + outputFormat + `

Code (if source_type is code):
` + code + `

API Description (if source_type is manual):
` + apiDescription + `

Project Summary (if source_type is project_zip):
` + projectSummary + `

Generate API documentation with curl examples and frontend API client guide.`

	return BuildPrompt(systemPrompt, userPrompt)
}
