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
		slog.Warn("Failed to parse AI response, using fallback", "error", err)
		result = s.buildFallbackResult(input)
		fallbackJSON, _ := json.Marshal(result)
		_ = s.reportService.FallbackReport(ctx, report.ID, fallbackJSON, "AI response parsing failed")
		return s.reportService.GetReport(ctx, report.ID)
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

	summaryStr := fmt.Sprintf("Project health check complete. Found %d issues.", len(result.Issues))
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

func (s *ProjectDoctorService) buildFallbackResult(input ProjectDoctorFormInput) dto.ProjectDoctorResult {
	return dto.ProjectDoctorResult{
		Scores:          []dto.ScoreItem{{Name: "Overall", Score: 0, MaxScore: 100, Comment: "AI parsing failed"}},
		Issues:          []dto.IssueItem{},
		Recommendations: []string{"AI response parsing failed. Please try again."},
		CodexPrompt:     "Retry project health check for " + input.ProjectName,
	}
}

func (s *ProjectDoctorService) buildMarkdownReport(result dto.ProjectDoctorResult, input ProjectDoctorFormInput, summary *service.ProjectSummary) string {
	md := "# Project Doctor Report\n\n"
	md += "## Project: " + input.ProjectName + "\n\n"
	if summary != nil && len(summary.DetectedStack) > 0 {
		md += "## Detected Stack\n\n"
		for _, s := range summary.DetectedStack {
			md += fmt.Sprintf("- %s\n", s)
		}
	}
	md += "\n## Scores\n\n"
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
