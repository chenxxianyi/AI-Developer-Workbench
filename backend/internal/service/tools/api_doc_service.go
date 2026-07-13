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

// APIDocService handles API Doc Builder operations.
type APIDocService struct {
	aiService     service.AIService
	reportService service.ReportService
	fileService   service.FileService
	zipService    service.ZipService
	uploadDir     string
}

// NewAPIDocService creates a new API Doc service.
func NewAPIDocService(
	aiService service.AIService,
	reportService service.ReportService,
	fileService service.FileService,
	zipService service.ZipService,
	uploadDir string,
) *APIDocService {
	return &APIDocService{
		aiService:     aiService,
		reportService: reportService,
		fileService:   fileService,
		zipService:    zipService,
		uploadDir:     uploadDir,
	}
}

// APIDocFormInput holds multipart form data for API Doc Builder.
type APIDocFormInput struct {
	Title          string
	SourceType     string
	BackendStack   string
	Code           string
	APIDescription string
	OutputFormat   string
	ProjectZip     *multipart.FileHeader
	ParentReportID string
	ProjectID      string
}

// Run executes the API Doc Builder tool.
func (s *APIDocService) Run(ctx context.Context, input APIDocFormInput) (*dto.ReportDTO, error) {
	if err := s.validateInput(input); err != nil {
		return nil, err
	}

	if input.ParentReportID != "" {
		if _, err := s.reportService.ValidateParentReport(ctx, model.ToolTypeAPIDoc, input.ParentReportID); err != nil {
			return nil, fmt.Errorf("invalid parent report: %w", err)
		}
	}
	project, err := s.reportService.ResolveProject(ctx, input.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project: %w", err)
	}

	req := dto.APIDocRequest{
		Title:          input.Title,
		SourceType:     input.SourceType,
		BackendStack:   input.BackendStack,
		Code:           input.Code,
		APIDescription: input.APIDescription,
		OutputFormat:   input.OutputFormat,
		ParentReportID: input.ParentReportID,
		ProjectID:      input.ProjectID,
	}
	inputData, _ := json.Marshal(req)

	report, err := s.reportService.CreateProcessingReport(ctx, model.ToolTypeAPIDoc, input.Title, input.SourceType, inputData, input.ParentReportID, input.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("failed to create report: %w", err)
	}

	// Handle project_zip mode.
	var projectSummaryText string
	if input.SourceType == "project_zip" && input.ProjectZip != nil {
		asset, err := s.fileService.SaveUpload(ctx, report.ID, input.ProjectZip, model.AssetTypeProjectZip, service.AllowedArchiveTypes())
		if err != nil {
			_ = s.reportService.FailReport(ctx, report.ID, fmt.Sprintf("zip upload failed: %v", err))
			return nil, fmt.Errorf("zip upload failed: %w", err)
		}

		zipPath := s.uploadDir + "/" + asset.RelativePath
		limits := service.DefaultZipLimits(120, 100, 12000, 300000)
		summary, err := s.zipService.ExtractAndAnalyze(zipPath, limits)
		if err != nil {
			_ = s.reportService.FailReport(ctx, report.ID, fmt.Sprintf("zip analysis failed: %v", err))
			return nil, fmt.Errorf("zip analysis failed: %w", err)
		}
		summaryJSON, _ := json.Marshal(summary)
		projectSummaryText = util.TruncateText(string(summaryJSON), 300000)
	}

	// Redact secrets before the content reaches the AI.
	code := util.RedactText(util.TruncateText(input.Code, 12000))
	apiDescription := util.RedactText(input.APIDescription)

	systemPrompt, userPrompt := prompts.BuildAPIDocPrompt(
		input.SourceType, input.BackendStack, code,
		apiDescription, input.OutputFormat, projectSummaryText,
	)
	userPrompt = prompts.AppendTrustedProjectContext(userPrompt, project)

	aiResult, err := s.aiService.GenerateJSON(ctx, service.AIRequest{
		ToolType:     model.ToolTypeAPIDoc,
		SystemPrompt: systemPrompt,
		UserPrompt:   userPrompt,
	})
	if err != nil {
		slog.Error("AI call failed for API Doc", "error", err)
		_ = s.reportService.FailReport(ctx, report.ID, err.Error())
		return nil, err
	}

	var result dto.APIDocResult
	if err := util.ParseAIResponseInto(aiResult.JSONText, &result); err != nil {
		slog.Warn("Failed to parse AI response, using fallback", "error", err)
		result = s.buildFallbackResult(input)
		fallbackJSON, _ := json.Marshal(result)
		_ = s.reportService.FallbackReport(ctx, report.ID, fallbackJSON, "AI response parsing failed")
		return s.reportService.GetReport(ctx, report.ID)
	}

	s.normalizeResult(&result, input)

	// Build generated files.
	generatedFiles := s.buildGeneratedFiles(result, input)

	summary := fmt.Sprintf("Generated API documentation with %d modules", len(result.Modules))
	reportJSON, _ := json.Marshal(result)

	// API Doc has no scores.
	return s.reportService.SucceedReport(ctx, report.ID, reportJSON, summary, nil, nil, generatedFiles)
}

func (s *APIDocService) validateInput(input APIDocFormInput) error {
	if input.Title == "" {
		return fmt.Errorf("title is required")
	}
	if input.SourceType == "" {
		return fmt.Errorf("source_type is required")
	}
	switch input.SourceType {
	case "code":
		if input.Code == "" {
			return fmt.Errorf("code is required for code mode")
		}
	case "project_zip":
		if input.ProjectZip == nil {
			return fmt.Errorf("project_zip is required for project_zip mode")
		}
	case "manual":
		if input.APIDescription == "" {
			return fmt.Errorf("api_description is required for manual mode")
		}
	default:
		return fmt.Errorf("invalid source_type: %s", input.SourceType)
	}
	if input.OutputFormat == "" {
		return fmt.Errorf("output_format is required")
	}
	return nil
}

func (s *APIDocService) normalizeResult(result *dto.APIDocResult, input APIDocFormInput) {
	if result.Modules == nil {
		result.Modules = []dto.ModuleItem{}
	}
	if result.Recommendations == nil {
		result.Recommendations = []string{}
	}
	result.ActionItems = dto.NormalizeActionItems(result.ActionItems)
	// Ensure content matches output_format.
	if input.OutputFormat == "markdown" || input.OutputFormat == "both" {
		if result.MarkdownContent == nil {
			empty := ""
			result.MarkdownContent = &empty
		}
	}
	if input.OutputFormat == "openapi" || input.OutputFormat == "both" {
		if result.OpenAPIContent == nil {
			empty := ""
			result.OpenAPIContent = &empty
		}
	}
}

func (s *APIDocService) buildGeneratedFiles(result dto.APIDocResult, input APIDocFormInput) []model.GeneratedFile {
	files := []model.GeneratedFile{}

	if result.MarkdownContent != nil && *result.MarkdownContent != "" {
		files = append(files, model.GeneratedFile{
			Filename: "API_DOCUMENTATION.md",
			Content:  *result.MarkdownContent,
			MimeType: "text/markdown",
			Language: "markdown",
		})
	}

	if result.OpenAPIContent != nil && *result.OpenAPIContent != "" {
		files = append(files, model.GeneratedFile{
			Filename: "openapi.json",
			Content:  *result.OpenAPIContent,
			MimeType: "application/json",
			Language: "json",
		})
	}

	return files
}

func (s *APIDocService) buildFallbackResult(input APIDocFormInput) dto.APIDocResult {
	return dto.APIDocResult{
		Modules:         []dto.ModuleItem{},
		Recommendations: []string{"AI response parsing failed. Please try again."},
		CodexPrompt:     "Retry API doc generation for " + input.SourceType + " mode",
	}
}
