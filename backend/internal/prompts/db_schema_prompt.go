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
  "migration_suggestions": ["string - phased migration steps with risk labels"],
  "er_diagram_mermaid": "string - Mermaid ER diagram in text format (optional)",
  "index_recommendations": [
    {"table": "string", "columns": "string", "reason": "string", "impact": "high|medium|low"}
  ],
  "migration_risks": [
    {"operation": "string - DDL statement", "risk": "high|medium|low", "description": "string", "rollback_plan": "string"}
  ],
  "recommendations": ["string - general recommendations"],
  "action_items": [
    {"id": "stable-kebab-case-id", "title": "string", "priority": "high|medium|low", "effort": "small|medium|large", "category": "string", "reason": "string", "suggested_prompt": "string", "issue_title": "string", "issue_body": "string - Markdown"}
  ],
  "codex_prompt": "string - prompt for generating migration scripts"
}

CRITICAL: For migration_risks, every DDL operation that could lock tables or cause downtime MUST be marked with its risk level and a rollback plan.
For er_diagram_mermaid, output valid Mermaid erDiagram syntax.
Distinguish between MySQL, PostgreSQL, and SQLite capabilities — don't suggest features that don't exist in the target database_type.
` + ActionItemsPromptSchema

// BuildDBSchemaPrompt creates the prompt for DB Schema Review.
func BuildDBSchemaPrompt(schemaType, databaseType, businessContext, schemaContent, targetGoal string) (string, string) {
	systemPrompt := `You are an expert database architect reviewing schema definitions.
Your task is to analyze the provided schema for quality, performance, and best practices.

Scoring dimensions:
1. Structure (结构评分) (0-100): Table design, normalization level, data types
2. Indexing (索引评分) (0-100): Index coverage for common query patterns, index types
3. Extensibility (扩展性评分) (0-100): Audit columns, soft delete, partitioning readiness
4. Data Integrity (数据完整性评分) (0-100): Constraints, NOT NULL, foreign keys, CHECK

For each index recommendation, explain WHY it's needed (cite a query pattern from the business context).
For migration risks, explicitly label operations that require table locks (ALTER TABLE) vs online-safe operations (CREATE INDEX CONCURRENTLY for PostgreSQL).

Generate a Mermaid erDiagram that visualizes all tables and relationships found in the schema.

IMPORTANT: You are analyzing schema text only. Do NOT connect to any database or execute any SQL.
All recommendations must be specific to the database_type: ` + databaseType + `.
` + DBSchemaPromptSchema

	userPrompt := `Schema Information:
- Type: ` + schemaType + ` (SQL, GORM, Prisma, or natural language)
- Database: ` + databaseType + `
- Business Context: ` + businessContext + `
- Target Goal: ` + targetGoal + `

Schema Content:
` + schemaContent + `

Analyze this schema. For each issue found, provide a rollback plan in the migration_risks section.
Generate an ER diagram in Mermaid format.`

	return BuildPrompt(systemPrompt, userPrompt)
}
