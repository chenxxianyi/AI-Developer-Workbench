package dto

// DBSchemaRequest is the input for DB Schema Review.
type DBSchemaRequest struct {
	Title            string `json:"title" binding:"required"`
	SchemaType       string `json:"schema_type" binding:"required"`
	DatabaseType     string `json:"database_type"`
	BusinessContext  string `json:"business_context"`
	SchemaContent    string `json:"schema_content" binding:"required"`
	TargetGoal       string `json:"target_goal"`
}

// DBSchemaResult is the output for DB Schema Review.
type DBSchemaResult struct {
	Scores             []ScoreItem  `json:"scores"`
	Issues             []IssueItem  `json:"issues"`
	OptimizedSchema    *string      `json:"optimized_schema,omitempty"`
	MigrationSuggestions []string   `json:"migration_suggestions,omitempty"`
	Recommendations    []string     `json:"recommendations"`
	CodexPrompt        string       `json:"codex_prompt"`
}