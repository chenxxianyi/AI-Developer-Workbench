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

// UIReviewService handles UI Review operations.
type UIReviewService struct {
	aiService     service.AIService
	reportService service.ReportService
	fileService   service.FileService
}

// NewUIReviewService creates a new UI Review service.
func NewUIReviewService(aiService service.AIService, reportService service.ReportService, fileService service.FileService) *UIReviewService {
	return &UIReviewService{aiService: aiService, reportService: reportService, fileService: fileService}
}

// UIReviewFormInput holds multipart form data for UI Review.
type UIReviewFormInput struct {
	Title       string
	ReviewMode  string
	PageType    string
	TargetStyle string
	Description string
	Code        string
	Screenshot  *multipart.FileHeader
}

// Run executes the UI Review tool.
func (s *UIReviewService) Run(ctx context.Context, input UIReviewFormInput) (*dto.ReportDTO, error) {
	// Validate based on review mode.
	if err := s.validateInput(input); err != nil {
		return nil, err
	}

	// Build request for report creation.
	req := dto.UIReviewRequest{
		Title:       input.Title,
		ReviewMode:  input.ReviewMode,
		PageType:    input.PageType,
		TargetStyle: input.TargetStyle,
		Description: input.Description,
		Code:        input.Code,
	}
	inputData, _ := json.Marshal(req)

	report, err := s.reportService.CreateProcessingReport(ctx, model.ToolTypeUIReview, input.Title, input.ReviewMode, inputData)
	if err != nil {
		return nil, fmt.Errorf("failed to create report: %w", err)
	}

	// Save screenshot if provided.
	var imagePath string
	if input.Screenshot != nil {
		asset, err := s.fileService.SaveUpload(ctx, report.ID, input.Screenshot, model.AssetTypeScreenshot, service.AllowedImageTypes())
		if err != nil {
			_ = s.reportService.FailReport(ctx, report.ID, fmt.Sprintf("screenshot upload failed: %v", err))
			return nil, fmt.Errorf("screenshot upload failed: %w", err)
		}
		imagePath = asset.RelativePath
	}

	// Truncate code if needed.
	code := util.TruncateText(input.Code, 12000)

	// Build prompt.
	systemPrompt, userPrompt := prompts.BuildUIReviewPrompt(
		input.ReviewMode, input.PageType, input.TargetStyle,
		input.Description, code,
	)

	// Call AI.
	needVision := input.ReviewMode == "screenshot" || input.ReviewMode == "screenshot_code"
	aiResult, err := s.aiService.GenerateJSON(ctx, service.AIRequest{
		ToolType:     model.ToolTypeUIReview,
		SystemPrompt: systemPrompt,
		UserPrompt:   userPrompt,
		ImagePath:    imagePath,
		NeedVision:   needVision,
	})
	if err != nil {
		slog.Error("AI call failed for UI Review", "error", err)
		_ = s.reportService.FailReport(ctx, report.ID, err.Error())
		return nil, err
	}

	var result dto.UIReviewResult
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
			Filename: "UI_REVIEW_REPORT.md",
			Content:  s.buildMarkdownReport(result, input),
			MimeType: "text/markdown",
			Language: "markdown",
		},
	}

	summary := fmt.Sprintf("UI review complete. Found %d issues.", len(result.Issues))
	reportJSON, _ := json.Marshal(result)

	return s.reportService.SucceedReport(ctx, report.ID, reportJSON, summary, totalScore, grade, generatedFiles)
}

func (s *UIReviewService) validateInput(input UIReviewFormInput) error {
	if input.Title == "" {
		return fmt.Errorf("title is required")
	}
	if input.ReviewMode == "" {
		return fmt.Errorf("review_mode is required")
	}
	switch input.ReviewMode {
	case "screenshot":
		if input.Screenshot == nil {
			return fmt.Errorf("screenshot is required for screenshot mode")
		}
	case "code":
		if input.Code == "" {
			return fmt.Errorf("code is required for code mode")
		}
	case "screenshot_code":
		if input.Screenshot == nil {
			return fmt.Errorf("screenshot is required for screenshot_code mode")
		}
		if input.Code == "" {
			return fmt.Errorf("code is required for screenshot_code mode")
		}
	default:
		return fmt.Errorf("invalid review_mode: %s", input.ReviewMode)
	}
	return nil
}

func (s *UIReviewService) normalizeResult(result *dto.UIReviewResult) {
	if result.Issues == nil {
		result.Issues = []dto.IssueItem{}
	}
	if result.Recommendations == nil {
		result.Recommendations = []string{}
	}
	for i := range result.Scores {
		result.Scores[i].Score = util.NormalizeScore(result.Scores[i].Score)
		result.Scores[i].MaxScore = 100
	}
	for i := range result.Issues {
		result.Issues[i].Severity = util.NormalizeSeverity(result.Issues[i].Severity)
	}
}

func (s *UIReviewService) buildFallbackResult(input UIReviewFormInput) dto.UIReviewResult {
	return dto.UIReviewResult{
		Scores:         []dto.ScoreItem{{Name: "Overall", Score: 0, MaxScore: 100, Comment: "AI parsing failed"}},
		Issues:         []dto.IssueItem{},
		Recommendations: []string{"AI response parsing failed. Please try again."},
		CodexPrompt:    "Retry UI review for " + input.ReviewMode + " mode",
	}
}

func (s *UIReviewService) buildMarkdownReport(result dto.UIReviewResult, input UIReviewFormInput) string {
	md := "# UI Review Report\n\n"
	md += "## Scores\n\n"
	for _, s := range result.Scores {
		md += fmt.Sprintf("- **%s**: %d/100 - %s\n", s.Name, s.Score, s.Comment)
	}
	md += "\n## Issues\n\n"
	for _, i := range result.Issues {
		md += fmt.Sprintf("- **[%s] %s** (%s): %s\n", i.Severity, i.Title, i.Category, i.Problem)
	}
	md += "\n## Recommendations\n\n"
	for _, r := range result.Recommendations {
		md += fmt.Sprintf("- %s\n", r)
	}
	return md
}