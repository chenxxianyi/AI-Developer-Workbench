package prompts

// AgentConfigPromptSchema describes the expected output format.
const AgentConfigPromptSchema = `
Output a JSON object with this exact structure:
{
  "generated_files_content": {
    "AGENTS.md": "string - project context for AI assistants",
    "TASK_PLAN.md": "string - task planning template",
    "CODING_RULES.md": "string - coding standards",
    "FRONTEND_STYLE_GUIDE.md": "string - frontend styling guide (optional)",
    "BACKEND_ARCHITECTURE.md": "string - backend architecture doc (optional)"
  },
  "recommendations": ["string - actionable recommendations"],
  "codex_prompt": "string - prompt for generating additional files"
}
`

// BuildAgentConfigPrompt creates the prompt for Agent Config Studio.
func BuildAgentConfigPrompt(projectName, projectType, frontendStack, backendStack, database, uiStyle, codingPrefs, strictRules string) (string, string) {
	systemPrompt := `You are an expert at creating AI assistant configuration files for software projects.
Your task is to generate comprehensive AGENTS.md and related configuration files that help AI coding assistants understand and work effectively with a project.

Focus on:
1. Tech stack and architecture context
2. Coding conventions and patterns used in the project
3. Common pitfalls and gotchas
4. Testing and quality requirements
5. Development workflow

` + AgentConfigPromptSchema

	userPrompt := `Project Information:
- Name: ` + projectName + `
- Type: ` + projectType + `
- Frontend Stack: ` + frontendStack + `
- Backend Stack: ` + backendStack + `
- Database: ` + database + `
- UI Style: ` + uiStyle + `
- Coding Preferences: ` + codingPrefs + `
- Strict Rules: ` + strictRules + `

Generate comprehensive AI assistant configuration files for this project.`

	return BuildPrompt(systemPrompt, userPrompt)
}