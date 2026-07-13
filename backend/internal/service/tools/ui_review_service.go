package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"mime/multipart"
	"path/filepath"

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
	zipService    service.ZipService
	uploadDir     string
}

// NewUIReviewService creates a new UI Review service.
func NewUIReviewService(
	aiService service.AIService,
	reportService service.ReportService,
	fileService service.FileService,
	zipService service.ZipService,
	uploadDir string,
) *UIReviewService {
	return &UIReviewService{
		aiService:     aiService,
		reportService: reportService,
		fileService:   fileService,
		zipService:    zipService,
		uploadDir:     uploadDir,
	}
}

// UIReviewFormInput holds multipart form data for UI Review.
type UIReviewFormInput struct {
	Title             string
	ReviewMode        string
	CodeSource        string
	PageType          string
	TargetStyle       string
	Description       string
	Code              string
	Screenshot        *multipart.FileHeader // legacy desktop screenshot field
	DesktopScreenshot *multipart.FileHeader
	MobileScreenshot  *multipart.FileHeader
	DesktopViewport   string
	MobileViewport    string
	ProjectZip        *multipart.FileHeader
	ParentReportID    string
	ProjectID         string
}

// Run executes the UI Review tool.
func (s *UIReviewService) Run(ctx context.Context, input UIReviewFormInput) (*dto.ReportDTO, error) {
	// Validate based on review mode.
	if err := s.validateInput(input); err != nil {
		return nil, err
	}

	// Build request for report creation.
	req := dto.UIReviewRequest{
		Title:           input.Title,
		ProjectID:       input.ProjectID,
		ReviewMode:      input.ReviewMode,
		CodeSource:      normalizeUIReviewCodeSource(input.CodeSource),
		PageType:        input.PageType,
		TargetStyle:     input.TargetStyle,
		Description:     input.Description,
		DesktopViewport: normalizeViewport(input.DesktopViewport, "1440x900"),
		MobileViewport:  normalizeViewport(input.MobileViewport, "390x844"),
		Code:            input.Code,
	}
	inputData, _ := json.Marshal(req)

	// Validate parent report lineage (same tool type) before creating the child.
	if input.ParentReportID != "" {
		if _, err := s.reportService.ValidateParentReport(ctx, model.ToolTypeUIReview, input.ParentReportID); err != nil {
			return nil, fmt.Errorf("invalid parent report: %w", err)
		}
	}
	project, err := s.reportService.ResolveProject(ctx, input.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project: %w", err)
	}

	report, err := s.reportService.CreateProcessingReport(ctx, model.ToolTypeUIReview, input.Title, input.ReviewMode, inputData, input.ParentReportID, input.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("failed to create report: %w", err)
	}

	// Save screenshots in viewport order. The legacy screenshot field is treated
	// as desktop so existing API clients continue to work.
	desktopScreenshot := input.DesktopScreenshot
	if desktopScreenshot == nil {
		desktopScreenshot = input.Screenshot
	}
	var imagePaths []string
	saveScreenshot := func(file *multipart.FileHeader) error {
		if file == nil {
			return nil
		}
		asset, err := s.fileService.SaveUpload(ctx, report.ID, file, model.AssetTypeScreenshot, service.AllowedImageTypes())
		if err != nil {
			return err
		}
		imagePaths = append(imagePaths, resolveUploadPath(s.uploadDir, asset.RelativePath))
		return nil
	}
	if err := saveScreenshot(desktopScreenshot); err != nil {
		_ = s.reportService.FailReport(ctx, report.ID, fmt.Sprintf("desktop screenshot upload failed: %v", err))
		return nil, fmt.Errorf("desktop screenshot upload failed: %w", err)
	}
	if err := saveScreenshot(input.MobileScreenshot); err != nil {
		_ = s.reportService.FailReport(ctx, report.ID, fmt.Sprintf("mobile screenshot upload failed: %v", err))
		return nil, fmt.Errorf("mobile screenshot upload failed: %w", err)
	}

	var projectSummaryText string
	if normalizeUIReviewCodeSource(input.CodeSource) == "project_zip" && input.ProjectZip != nil {
		asset, err := s.fileService.SaveUpload(ctx, report.ID, input.ProjectZip, model.AssetTypeProjectZip, service.AllowedArchiveTypes())
		if err != nil {
			_ = s.reportService.FailReport(ctx, report.ID, fmt.Sprintf("project zip upload failed: %v", err))
			return nil, fmt.Errorf("project zip upload failed: %w", err)
		}

		zipPath := resolveUploadPath(s.uploadDir, asset.RelativePath)
		limits := service.DefaultZipLimits(120, 100, 12000, 300000)
		summary, err := s.zipService.ExtractAndAnalyze(zipPath, limits)
		if err != nil {
			_ = s.reportService.FailReport(ctx, report.ID, fmt.Sprintf("project zip analysis failed: %v", err))
			return nil, fmt.Errorf("project zip analysis failed: %w", err)
		}
		projectSummaryJSON, _ := json.Marshal(summary)
		projectSummaryText = util.TruncateText(string(projectSummaryJSON), 300000)
	}

	// Truncate code if needed. Redact secrets before the content reaches the AI.
	code := util.RedactText(util.TruncateText(input.Code, 12000))

	// Build prompt with explicit viewport-to-image mapping.
	viewportDescription := fmt.Sprintf("%s\nScreenshot viewports: desktop=%s; mobile=%s. Images are ordered desktop then mobile.", input.Description, normalizeViewport(input.DesktopViewport, "1440x900"), normalizeViewport(input.MobileViewport, "390x844"))
	systemPrompt, userPrompt := prompts.BuildUIReviewPrompt(
		input.ReviewMode, normalizeUIReviewCodeSource(input.CodeSource),
		input.PageType, input.TargetStyle,
		viewportDescription, code, projectSummaryText,
	)
	userPrompt = prompts.AppendTrustedProjectContext(userPrompt, project)

	// Call AI.
	needVision := input.ReviewMode == "screenshot" || input.ReviewMode == "screenshot_code"
	aiResult, err := s.aiService.GenerateJSON(ctx, service.AIRequest{
		ToolType:     model.ToolTypeUIReview,
		SystemPrompt: systemPrompt,
		UserPrompt:   userPrompt,
		ImagePaths:   imagePaths,
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

	summary := fmt.Sprintf("UI 审查完成，发现 %d 个问题。", len(result.Issues))
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
		if input.Screenshot == nil && input.DesktopScreenshot == nil && input.MobileScreenshot == nil {
			return fmt.Errorf("at least one screenshot is required for screenshot mode")
		}
	case "code":
		if input.Code == "" && input.ProjectZip == nil {
			return fmt.Errorf("code or project_zip is required for code mode")
		}
	case "screenshot_code":
		if input.Screenshot == nil && input.DesktopScreenshot == nil && input.MobileScreenshot == nil {
			return fmt.Errorf("at least one screenshot is required for screenshot_code mode")
		}
		if input.Code == "" && input.ProjectZip == nil {
			return fmt.Errorf("code or project_zip is required for screenshot_code mode")
		}
	default:
		return fmt.Errorf("invalid review_mode: %s", input.ReviewMode)
	}
	return nil
}

func normalizeViewport(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

func normalizeUIReviewCodeSource(codeSource string) string {
	if codeSource == "project_zip" {
		return "project_zip"
	}
	return "paste"
}

func (s *UIReviewService) normalizeResult(result *dto.UIReviewResult) {
	if result.Issues == nil {
		result.Issues = []dto.IssueItem{}
	}
	if result.Recommendations == nil {
		result.Recommendations = []string{}
	}
	result.ActionItems = dto.NormalizeActionItems(result.ActionItems)
	if len(result.ScreenshotContexts) == 0 {
		result.ScreenshotContexts = []dto.ScreenshotContext{{Kind: "desktop", Viewport: "1440x900"}, {Kind: "mobile", Viewport: "390x844"}}
	}
	for i := range result.Issues {
		issue := &result.Issues[i]
		if issue.Severity == "high" && issue.ComponentPrompt == "" {
			issue.ComponentPrompt = fmt.Sprintf("Fix the %s component issue: %s. Apply the suggested change and add a regression test.", issue.Category, issue.Title)
		}
		if result.Issues[i].Region != nil {
			r := result.Issues[i].Region
			r.X, r.Y = clampPercent(r.X), clampPercent(r.Y)
			r.Width, r.Height = clampPercent(r.Width), clampPercent(r.Height)
		}
	}
	for i := range result.Scores {
		result.Scores[i].Score = util.NormalizeScore(result.Scores[i].Score)
		result.Scores[i].MaxScore = 100
	}
	for i := range result.Issues {
		result.Issues[i].Severity = util.NormalizeSeverity(result.Issues[i].Severity)
	}
}

func clampPercent(value float64) float64 {
	if value < 0 {
		return 0
	}
	if value > 100 {
		return 100
	}
	return value
}

func (s *UIReviewService) buildFallbackResult(input UIReviewFormInput) dto.UIReviewResult {
	return dto.UIReviewResult{
		Scores:          []dto.ScoreItem{{Name: "总体评分", Score: 0, MaxScore: 100, Comment: "AI 结果解析失败"}},
		Issues:          []dto.IssueItem{},
		Recommendations: []string{"AI 结果解析失败，请重试。"},
		CodexPrompt:     "请重新执行 " + input.ReviewMode + " 模式的 UI 审查。",
	}
}

func (s *UIReviewService) buildMarkdownReport(result dto.UIReviewResult, input UIReviewFormInput) string {
	md := "# UI 审查报告\n\n"
	md += "## 评分\n\n"
	for _, s := range result.Scores {
		md += fmt.Sprintf("- **%s**: %d/100 - %s\n", s.Name, s.Score, s.Comment)
	}
	md += "\n## 发现的问题\n\n"
	for _, i := range result.Issues {
		md += fmt.Sprintf("- **[%s] %s** (%s): %s\n", i.Severity, i.Title, i.Category, i.Problem)
	}
	md += "\n## 改进建议\n\n"
	for _, r := range result.Recommendations {
		md += fmt.Sprintf("- %s\n", r)
	}
	return md
}

func resolveUploadPath(uploadDir, relativePath string) string {
	return filepath.Join(uploadDir, filepath.FromSlash(relativePath))
}
