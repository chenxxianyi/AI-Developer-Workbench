package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"ai-developer-workbench/internal/dto"
	"ai-developer-workbench/internal/model"
	"ai-developer-workbench/internal/prompts"
	"ai-developer-workbench/internal/service"
	"ai-developer-workbench/internal/util"
)

// DBSchemaService handles DB Schema Review operations.
type DBSchemaService struct {
	aiService     service.AIService
	reportService service.ReportService
}

// NewDBSchemaService creates a new DB Schema service.
func NewDBSchemaService(aiService service.AIService, reportService service.ReportService) *DBSchemaService {
	return &DBSchemaService{aiService: aiService, reportService: reportService}
}

// Run executes the DB Schema Review tool.
func (s *DBSchemaService) Run(ctx context.Context, req dto.DBSchemaRequest) (*dto.ReportDTO, error) {
	if req.ParentReportID != "" {
		if _, err := s.reportService.ValidateParentReport(ctx, model.ToolTypeDBSchema, req.ParentReportID); err != nil {
			return nil, fmt.Errorf("invalid parent report: %w", err)
		}
	}
	project, err := s.reportService.ResolveProject(ctx, req.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project: %w", err)
	}
	inputData, _ := json.Marshal(req)
	report, err := s.reportService.CreateProcessingReport(ctx, model.ToolTypeDBSchema, req.Title, "json", inputData, req.ParentReportID, req.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("failed to create report: %w", err)
	}

	// Redact secrets before the content reaches the AI.
	schemaContent := util.RedactText(req.SchemaContent)
	businessContext := util.RedactText(req.BusinessContext)

	systemPrompt, userPrompt := prompts.BuildDBSchemaPrompt(
		req.SchemaType, req.DatabaseType, businessContext,
		schemaContent, req.TargetGoal,
	)
	userPrompt = prompts.AppendTrustedProjectContext(userPrompt, project)

	aiResult, err := s.aiService.GenerateJSON(ctx, service.AIRequest{
		ToolType:     model.ToolTypeDBSchema,
		SystemPrompt: systemPrompt,
		UserPrompt:   userPrompt,
	})
	if err != nil {
		slog.Error("AI call failed for DB Schema", "error", err)
		_ = s.reportService.FailReport(ctx, report.ID, err.Error())
		return nil, err
	}

	var result dto.DBSchemaResult
	if err := util.ParseAIResponseInto(aiResult.JSONText, &result); err != nil {
		parseErr := fmt.Errorf("parse AI response: %w", err)
		slog.Error("Failed to parse AI response", "error", err)
		_ = s.reportService.FailReport(ctx, report.ID, parseErr.Error())
		return nil, parseErr
	}

	// Normalize.
	s.normalizeResult(&result)

	// Compute total score.
	var totalScore *int
	var grade *string
	if len(result.Scores) > 0 {
		avg := computeAverageScore(result.Scores)
		totalScore = &avg
		g := util.ComputeGrade(avg)
		grade = &g
	}

	// Build generated files.
	generatedFiles := []model.GeneratedFile{
		{
			Filename: "DB_SCHEMA_REVIEW.md",
			Content:  s.buildMarkdownReport(result, req),
			MimeType: "text/markdown",
			Language: "markdown",
		},
	}
	if result.OptimizedSchema != nil && *result.OptimizedSchema != "" {
		generatedFiles = append(generatedFiles, model.GeneratedFile{
			Filename: "migration.sql",
			Content:  "-- ⚠ Suggestion only — review before execution\n\n" + *result.OptimizedSchema,
			MimeType: "text/x-sql",
			Language: "sql",
		})
	}

	summary := fmt.Sprintf("Schema review complete. Found %d issues.", len(result.Issues))
	reportJSON, _ := json.Marshal(result)

	return s.reportService.SucceedReport(ctx, report.ID, reportJSON, summary, totalScore, grade, generatedFiles)
}

func (s *DBSchemaService) normalizeResult(result *dto.DBSchemaResult) {
	if result.Issues == nil {
		result.Issues = []dto.IssueItem{}
	}
	if result.Recommendations == nil {
		result.Recommendations = []string{}
	}
	if result.MigrationSuggestions == nil {
		result.MigrationSuggestions = []string{}
	}
	result.ActionItems = dto.NormalizeActionItems(result.ActionItems)
	for i := range result.Scores {
		result.Scores[i].Score = util.NormalizeScore(result.Scores[i].Score)
		result.Scores[i].MaxScore = 100
	}
	for i := range result.Issues {
		result.Issues[i].Severity = util.NormalizeSeverity(result.Issues[i].Severity)
	}
}

func (s *DBSchemaService) buildMarkdownReport(result dto.DBSchemaResult, req dto.DBSchemaRequest) string {
	md := "# DB Schema Review\n\n"
	md += "## Scores\n\n"
	for _, s := range result.Scores {
		md += fmt.Sprintf("- **%s**: %d/100 - %s\n", s.Name, s.Score, s.Comment)
	}
	md += "\n## Issues\n\n"
	for _, i := range result.Issues {
		md += fmt.Sprintf("- **[%s] %s**: %s\n", i.Severity, i.Title, i.Problem)
	}
	md += "\n## Recommendations\n\n"
	for _, r := range result.Recommendations {
		md += fmt.Sprintf("- %s\n", r)
	}
	return md
}

func computeAverageScore(scores []dto.ScoreItem) int {
	if len(scores) == 0 {
		return 0
	}
	total := 0
	for _, s := range scores {
		total += s.Score
	}
	return total / len(scores)
}
