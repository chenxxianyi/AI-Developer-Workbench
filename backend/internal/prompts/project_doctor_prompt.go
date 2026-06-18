package prompts

// ProjectDoctorPromptSchema describes the expected output format.
const ProjectDoctorPromptSchema = `
Output a JSON object with this exact structure:
{
  "scores": [
    {"name": "string - dimension name", "score": number (0-100), "max_score": 100, "comment": "string"}
  ],
  "issues": [
    {"title": "string", "severity": "high|medium|low", "category": "string", "problem": "string", "suggestion": "string", "action": "string"}
  ],
  "recommendations": ["string - actionable recommendations"],
  "codex_prompt": "string - prompt for generating fixes"
}
`

// BuildProjectDoctorPrompt creates the prompt for Project Doctor.
func BuildProjectDoctorPrompt(projectName, techStack, projectDescription, analysisDepth string, projectSummary string) (string, string) {
	systemPrompt := `You are an expert software architect analyzing project health.
Your task is to review a project's structure, code quality, and engineering practices.

Scoring dimensions:
1. Project Structure (0-100): Organization, separation of concerns
2. Documentation (0-100): README, API docs, inline comments
3. Dependency Management (0-100): Versioning, security, updates
4. Code Organization (0-100): Naming, patterns, modularity
5. Testing Coverage (0-100): Unit tests, integration tests

IMPORTANT: You are analyzing static code only. Do NOT execute any build commands, tests, or scripts.
` + ProjectDoctorPromptSchema

	userPrompt := `Project Information:
- Name: ` + projectName + `
- Tech Stack: ` + techStack + `
- Description: ` + projectDescription + `
- Analysis Depth: ` + analysisDepth + `

Project Summary (from static scan):
` + projectSummary + `

Analyze this project's health and provide scores, issues, and recommendations.`

	return BuildPrompt(systemPrompt, userPrompt)
}