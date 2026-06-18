package prompts

// DBSchemaPromptSchema describes the expected output format.
const DBSchemaPromptSchema = `
Output a JSON object with this exact structure:
{
  "scores": [
    {"name": "string - dimension name", "score": number (0-100), "max_score": 100, "comment": "string"}
  ],
  "issues": [
    {"title": "string", "severity": "high|medium|low", "category": "string", "problem": "string", "suggestion": "string", "action": "string"}
  ],
  "optimized_schema": "string - SQL DDL for optimized schema (optional)",
  "migration_suggestions": ["string - specific migration steps"],
  "recommendations": ["string - general recommendations"],
  "codex_prompt": "string - prompt for generating migration scripts"
}
`

// BuildDBSchemaPrompt creates the prompt for DB Schema Review.
func BuildDBSchemaPrompt(schemaType, databaseType, businessContext, schemaContent, targetGoal string) (string, string) {
	systemPrompt := `You are an expert database architect reviewing schema definitions.
Your task is to analyze the provided schema for quality, performance, and best practices.

Scoring dimensions:
1. Naming Conventions (0-100): Consistent naming, clear purpose
2. Normalization (0-100): Appropriate level, no redundancy
3. Indexing Strategy (0-100): Efficient indexes, covers query patterns
4. Data Integrity (0-100): Constraints, types, validation
5. Extensibility (0-100): Audit columns, soft delete, versioning

IMPORTANT: You are analyzing schema text only. Do NOT connect to any database or execute any SQL.
` + DBSchemaPromptSchema

	userPrompt := `Schema Information:
- Type: ` + schemaType + ` (SQL, GORM, Prisma, or natural language)
- Database: ` + databaseType + `
- Business Context: ` + businessContext + `
- Target Goal: ` + targetGoal + `

Schema Content:
` + schemaContent + `

Analyze this schema and provide scores, issues, and optimization suggestions.`

	return BuildPrompt(systemPrompt, userPrompt)
}