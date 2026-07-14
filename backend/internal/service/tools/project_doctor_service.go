package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"mime/multipart"

	"ai-developer-workbench/internal/dto"
	"ai-developer-workbench/internal/model"
	"ai-developer-workbench/internal/prompts"
	"ai-developer-workbench/internal/service"
	"ai-developer-workbench/internal/util"
)

// ProjectDoctorService handles Project Doctor operations.
type ProjectDoctorService struct {
	aiService     service.AIService
	reportService service.ReportService
	fileService   service.FileService
	zipService    service.ZipService
	uploadDir     string
	tempDir       string
}

// NewProjectDoctorService creates a new Project Doctor service.
func NewProjectDoctorService(
	aiService service.AIService,
	reportService service.ReportService,
	fileService service.FileService,
	zipService service.ZipService,
	uploadDir, tempDir string,
) *ProjectDoctorService {
	return &ProjectDoctorService{
		aiService:     aiService,
		reportService: reportService,
		fileService:   fileService,
		zipService:    zipService,
		uploadDir:     uploadDir,
		tempDir:       tempDir,
	}
}

// ProjectDoctorFormInput holds multipart form data for Project Doctor.
type ProjectDoctorFormInput struct {
	Title              string
	ProjectName        string
	TechStack          string
	ProjectDescription string
	AnalysisDepth      string
	ProjectZip         *multipart.FileHeader
	ParentReportID     string
	ProjectID          string
}

// Run executes the Project Doctor tool.
func (s *ProjectDoctorService) Run(ctx context.Context, input ProjectDoctorFormInput) (*dto.ReportDTO, error) {
	if input.Title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if input.ProjectZip == nil {
		return nil, fmt.Errorf("project_zip is required")
	}

	if input.ParentReportID != "" {
		if _, err := s.reportService.ValidateParentReport(ctx, model.ToolTypeProjectDoctor, input.ParentReportID); err != nil {
			return nil, fmt.Errorf("invalid parent report: %w", err)
		}
	}
	project, err := s.reportService.ResolveProject(ctx, input.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project: %w", err)
	}

	req := dto.ProjectDoctorRequest{
		Title:              input.Title,
		ProjectName:        input.ProjectName,
		TechStack:          input.TechStack,
		ProjectDescription: input.ProjectDescription,
		AnalysisDepth:      input.AnalysisDepth,
		ParentReportID:     input.ParentReportID,
		ProjectID:          input.ProjectID,
	}
	inputData, _ := json.Marshal(req)

	report, err := s.reportService.CreateProcessingReport(ctx, model.ToolTypeProjectDoctor, input.Title, "project_zip", inputData, input.ParentReportID, input.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("failed to create report: %w", err)
	}

	// Save ZIP file.
	asset, err := s.fileService.SaveUpload(ctx, report.ID, input.ProjectZip, model.AssetTypeProjectZip, service.AllowedArchiveTypes())
	if err != nil {
		_ = s.reportService.FailReport(ctx, report.ID, fmt.Sprintf("zip upload failed: %v", err))
		return nil, fmt.Errorf("zip upload failed: %w", err)
	}

	// Extract and analyze ZIP.
	zipPath := s.uploadDir + "/" + asset.RelativePath
	limits := service.DefaultZipLimits(120, 100, 12000, 300000)
	summary, err := s.zipService.ExtractAndAnalyze(zipPath, limits)
	if err != nil {
		_ = s.reportService.FailReport(ctx, report.ID, fmt.Sprintf("zip analysis failed: %v", err))
		return nil, fmt.Errorf("zip analysis failed: %w", err)
	}

	// Build project summary text.
	summaryJSON, _ := json.Marshal(summary)
	summaryText := util.TruncateText(string(summaryJSON), 300000)

	// Build prompt. Redact the free-text project description before it reaches the AI.
	systemPrompt, userPrompt := prompts.BuildProjectDoctorPrompt(
		input.ProjectName, input.TechStack,
		util.RedactText(input.ProjectDescription), input.AnalysisDepth, summaryText,
	)
	userPrompt = prompts.AppendTrustedProjectContext(userPrompt, project)

	// Call AI.
	aiResult, err := s.aiService.GenerateJSON(ctx, service.AIRequest{
		ToolType:     model.ToolTypeProjectDoctor,
		SystemPrompt: systemPrompt,
		UserPrompt:   userPrompt,
	})
	if err != nil {
		slog.Error("AI call failed for Project Doctor", "error", err)
		_ = s.reportService.FailReport(ctx, report.ID, err.Error())
		return nil, err
	}

	var result dto.ProjectDoctorResult
	if err := util.ParseAIResponseInto(aiResult.JSONText, &result); err != nil {
		parseErr := fmt.Errorf("parse AI response: %w", err)
		slog.Error("Failed to parse AI response", "error", err)
		_ = s.reportService.FailReport(ctx, report.ID, parseErr.Error())
		return nil, parseErr
	}

	s.normalizeResult(&result)

	var totalScore *int
	var grade *string
	if len(result.Scores) > 0 {
		avg := computeAverageScore(result.Scores)
		totalScore = &avg
		g := util.ComputeGrade(avg)
		grade = &g
	}

	generatedFiles := []model.GeneratedFile{
		{
			Filename: "PROJECT_DOCTOR_REPORT.md",
			Content:  s.buildMarkdownReport(result, input, summary),
			MimeType: "text/markdown",
			Language: "markdown",
		},
	}

	summaryStr := fmt.Sprintf("Project health check complete. %d issues found, %d tech debt items identified.", len(result.Issues), len(result.TechDebtItems))
	reportJSON, _ := json.Marshal(result)

	return s.reportService.SucceedReport(ctx, report.ID, reportJSON, summaryStr, totalScore, grade, generatedFiles)
}

func (s *ProjectDoctorService) normalizeResult(result *dto.ProjectDoctorResult) {
	if result.Issues == nil {
		result.Issues = []dto.IssueItem{}
	}
	if result.Recommendations == nil {
		result.Recommendations = []string{}
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

func (s *ProjectDoctorService) buildMarkdownReport(result dto.ProjectDoctorResult, input ProjectDoctorFormInput, summary *service.ProjectSummary) string {
	md := "# Project Doctor Report\n\n"
	md += "## Project: " + input.ProjectName + "\n\n"
	if summary != nil && len(summary.DetectedStack) > 0 {
		md += "## Detected Stack\n\n"
		for _, st := range summary.DetectedStack {
			md += fmt.Sprintf("- %s\n", st)
		}
	}

	// Evidence files
	if len(result.EvidenceFiles) > 0 {
		md += "\n## Evidence Files\n\n"
		md += "| Artifact | Present | Notes |\n"
		md += "|----------|---------|-------|\n"
		for _, ev := range result.EvidenceFiles {
			status := "❌"
			if ev.Present {
				status = "✅"
			}
			md += fmt.Sprintf("| %s (%s) | %s | %s |\n", ev.Path, ev.Type, status, ev.Notes)
		}
	}

	md += "\n## Scores\n\n"
	for _, sc := range result.Scores {
		md += fmt.Sprintf("- **%s**: %d/100 - %s\n", sc.Name, sc.Score, sc.Comment)
	}

	// Tech debt
	if len(result.TechDebtItems) > 0 {
		md += "\n## Technical Debt (Prioritized)\n\n"
		md += "| Item | Impact | Cost | Category |\n"
		md += "|------|--------|------|----------|\n"
		for _, td := range result.TechDebtItems {
			md += fmt.Sprintf("| %s | %s | %s | %s |\n", td.Title, td.Impact, td.Cost, td.Category)
			md += fmt.Sprintf("  _Description:_ %s\n\n", td.Description)
			if td.SuggestedFix != "" {
				md += fmt.Sprintf("  _Fix:_ %s\n\n", td.SuggestedFix)
			}
		}
	}

	md += "\n## Issues\n\n"
	for _, i := range result.Issues {
		md += fmt.Sprintf("- **[%s] %s**: %s\n", i.Severity, i.Title, i.Problem)
	}
	md += "\n## Recommendations\n\n"
	for _, r := range result.Recommendations {
		md += fmt.Sprintf("- %s\n", r)
	}

	// Action items
	if len(result.ActionItems) > 0 {
		md += "\n## Action Items\n\n"
		for _, a := range result.ActionItems {
			md += fmt.Sprintf("- [%s] **%s** (_%s effort_): %s\n", a.Priority, a.Title, a.Effort, a.Reason)
		}
	}

	return md
}
