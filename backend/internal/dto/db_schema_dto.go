package dto

// DBSchemaRequest is the input for DB Schema Review.
type DBSchemaRequest struct {
	Title           string `json:"title" binding:"required"`
	ProjectID       string `json:"project_id,omitempty"`
	SchemaType      string `json:"schema_type" binding:"required"`
	DatabaseType    string `json:"database_type"`
	BusinessContext string `json:"business_context"`
	SchemaContent   string `json:"schema_content" binding:"required"`
	TargetGoal      string `json:"target_goal"`
	ParentReportID  string `json:"parent_report_id,omitempty"`
}

// DBSchemaResult is the output for DB Schema Review.
type DBSchemaResult struct {
	Scores               []ScoreItem      `json:"scores"`
	Issues               []IssueItem      `json:"issues"`
	OptimizedSchema      *string          `json:"optimized_schema,omitempty"`
	MigrationSuggestions []string         `json:"migration_suggestions,omitempty"`
	ERDiagramMermaid     *string          `json:"er_diagram_mermaid,omitempty"`
	IndexRecommendations []IndexRec       `json:"index_recommendations,omitempty"`
	MigrationRisks       []MigrationRisk  `json:"migration_risks,omitempty"`
	Recommendations      []string         `json:"recommendations"`
	ActionItems          []ActionItem     `json:"action_items,omitempty"`
	CodexPrompt          string           `json:"codex_prompt"`
}

// IndexRec represents an index recommendation.
type IndexRec struct {
	Table   string `json:"table"`
	Columns string `json:"columns"`
	Reason  string `json:"reason"`
	Impact  string `json:"impact"` // high|medium|low
}

// MigrationRisk represents a migration risk assessment.
type MigrationRisk struct {
	Operation    string `json:"operation"`
	Risk         string `json:"risk"` // high|medium|low
	Description  string `json:"description"`
	RollbackPlan string `json:"rollback_plan,omitempty"`
}
