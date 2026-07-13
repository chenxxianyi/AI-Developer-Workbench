package prompts

// ProjectDoctorPromptSchema describes the expected output format.
const ProjectDoctorPromptSchema = `
Output a JSON object with this exact structure:
{
  "scores": [
    {"name": "string - dimension name", "score": number (0-100), "max_score": 100, "comment": "string - evidence-based explanation"}
  ],
  "evidence_files": [
    {"path": "README.md", "type": "readme", "present": true, "notes": "missing project setup steps"}
  ],
  "issues": [
    {"title": "string", "severity": "high|medium|low", "category": "string", "problem": "string", "suggestion": "string", "action": "string"}
  ],
  "tech_debt": [
    {"title": "string", "impact": "high|medium|low", "cost": "high|medium|low", "category": "string", "description": "string", "suggested_fix": "string"}
  ],
  "recommendations": ["string - actionable recommendations"],
  "action_items": [
    {"id": "stable-kebab-case-id", "title": "string", "priority": "high|medium|low", "effort": "small|medium|large", "category": "string", "reason": "string", "suggested_prompt": "string", "issue_title": "string", "issue_body": "string - Markdown"}
  ],
  "codex_prompt": "string - prompt for generating fixes"
}
` + ActionItemsPromptSchema

// BuildProjectDoctorPrompt creates the prompt for Project Doctor.
func BuildProjectDoctorPrompt(projectName, techStack, projectDescription, analysisDepth string, projectSummary string) (string, string) {
	systemPrompt := `You are an expert software architect analyzing project health.
Your task is to review a project's structure, code quality, and engineering practices.

Scoring dimensions (0-100 each, scoring MUST cite evidence from the project summary):
1. Structure Clarity (结构清晰度): Directory layout, separation of concerns, module boundaries
2. Maintainability (可维护性): Code duplication, hardcoded config, code complexity
3. Testability (可测试性): Test presence, test commands found, test framework usage
4. Deployability (可部署性): Dockerfile, compose, CI config, health checks
5. Documentation (文档完整度): README, AGENTS.md, API docs, inline comments
6. Agent Readiness (Agent 可接手程度): Are AGENTS.md, TASK_PLAN.md, clear entry points present? Can an AI coding agent immediately start contributing?

Evidence check — for each of the following, report whether found in the project:
- README.md: root-level readme with setup instructions
- AGENTS.md: AI coding agent instructions
- Lock file(s): package-lock.json, yarn.lock, go.sum, Cargo.lock, etc.
- Dockerfile: any containerization config
- CI config: GitHub Actions, Jenkins, GitLab CI, etc.
- Test commands: recognizable test scripts or test directories
- .env.example: template for required environment variables

For dependency risks: only report findings visible in lock files or manifest files from the project summary.
Do NOT guess online versions. Mark uncertain findings with "[Needs verification]".

Output issues sorted by severity (high first), and tech debt items sorted by (impact high, cost low) first.

IMPORTANT: You are analyzing static code only. Do NOT execute any build commands, tests, or scripts.
` + ProjectDoctorPromptSchema

	userPrompt := `Project Information:
- Name: ` + projectName + `
- Tech Stack: ` + techStack + `
- Description: ` + projectDescription + `
- Analysis Depth: ` + analysisDepth + `

Project Summary (from static scan):
` + projectSummary + `

Analyze this project's health. For each scoring dimension, cite specific evidence files seen in the project summary.
List missing artifacts in the evidence_files array. Prioritize tech debt by (impact × 1/cost).`

	return BuildPrompt(systemPrompt, userPrompt)
}
