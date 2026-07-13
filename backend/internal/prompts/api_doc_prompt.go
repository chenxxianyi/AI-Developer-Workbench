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

Output format requirements:
- If output_format is "markdown": only generate markdown_content
- If output_format is "openapi": only generate openapi_content
- If output_format is "both": generate both fields

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

Generate API documentation for these endpoints.`

	return BuildPrompt(systemPrompt, userPrompt)
}
