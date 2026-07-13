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
    "BACKEND_ARCHITECTURE.md": "string - backend architecture doc (optional)",
    "REVIEW_CHECKLIST.md": "string - code review checklist (optional)"
  },
  "target_format": "codex|copilot|cursor|windsurf",
  "missing_confirmations": ["string - items that need user verification before use"],
  "recommendations": ["string - actionable recommendations"],
  "action_items": [
    {"id": "stable-kebab-case-id", "title": "string", "priority": "high|medium|low", "effort": "small|medium|large", "category": "string", "reason": "string", "suggested_prompt": "string", "issue_title": "string", "issue_body": "string - Markdown"}
  ],
  "codex_prompt": "string - prompt for generating additional files"
}

CRITICAL RULES:
- All generated content MUST reference real project commands and paths from the user's input.
- When you cannot confirm a command or path, mark the claim with "[CONFIRM: <what to check>]".
- Do NOT fabricate paths, commands, or tooling that are not described in the project information.
- If coding_preferences and strict_rules conflict, list the conflicts in missing_confirmations.
` + ActionItemsPromptSchema

// BuildAgentConfigPrompt creates the prompt for Agent Config Studio.
func BuildAgentConfigPrompt(projectName, projectType, frontendStack, backendStack, database, uiStyle, codingPrefs, strictRules string) (string, string) {
	systemPrompt := `You are an expert at creating AI assistant configuration files for software projects.
Your task is to generate comprehensive AGENTS.md and related configuration files that help AI coding assistants understand and work effectively with a project.

Output format: Codex (AGENTS.md + CODING_RULES.md). Also generate REVIEW_CHECKLIST.md.

Focus on:
1. Tech stack and architecture context — reference real project names and paths
2. Coding conventions and patterns — derive from user input, do not invent
3. Common pitfalls and gotchas — based on the described tech stack
4. Testing and quality requirements
5. Development workflow — use commands from the user's input where available

For target_format: choose the most appropriate format based on the tech stack.
- codex: general purpose, works with any AI coding tool
- copilot: GitHub Copilot instructions format
- cursor: Cursor/Windsurf rules format

Mark anything you're uncertain about with "[CONFIRM: ...]" in the generated content.
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

Generate comprehensive AI assistant configuration files for this project.
Use real commands and paths from the project information where available.`

	return BuildPrompt(systemPrompt, userPrompt)
}
