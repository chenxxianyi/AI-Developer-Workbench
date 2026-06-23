package tools

import (
	"context"
	"encoding/json"
	"errors"
	"mime/multipart"
	"path/filepath"
	"strings"
	"testing"

	"ai-developer-workbench/internal/dto"
	"ai-developer-workbench/internal/model"
	"ai-developer-workbench/internal/service"
	"gorm.io/datatypes"
)

func TestResolveUploadPathIncludesConfiguredUploadDirectory(t *testing.T) {
	uploadDir := filepath.Join("var", "ai-workbench", "uploads")
	relativePath := filepath.ToSlash(filepath.Join(
		"report-id",
		"source",
		"screenshot.png",
	))

	got := resolveUploadPath(uploadDir, relativePath)
	want := filepath.Join(uploadDir, "report-id", "source", "screenshot.png")

	if got != want {
		t.Fatalf("resolveUploadPath() = %q, want %q", got, want)
	}
}

func TestUIReviewAcceptsProjectZipAsCodeSource(t *testing.T) {
	ai := &fakeUIReviewAIService{
		result: &service.AIResult{
			JSONText: `{"scores":[{"name":"视觉层级","score":80,"max_score":100,"comment":"结构清晰"}],"issues":[],"recommendations":["补充按钮状态"],"codex_prompt":"优化组件状态"}`,
		},
	}
	reports := newFakeUIReviewReportService()
	files := &fakeUIReviewFileService{}
	zips := &fakeUIReviewZipService{
		summary: &service.ProjectSummary{
			Tree:          []string{"package.json", "src/App.vue"},
			DetectedStack: []string{"vue"},
			ImportantFiles: []service.FileSummary{
				{Path: "src/App.vue", Content: "<template><button>提交</button></template>", Size: 42},
			},
		},
	}

	svc := NewUIReviewService(ai, reports, files, zips, filepath.Join("tmp", "uploads"))

	_, err := svc.Run(context.Background(), UIReviewFormInput{
		Title:      "完整项目审查",
		ReviewMode: "code",
		CodeSource: "project_zip",
		ProjectZip: &multipart.FileHeader{
			Filename: "frontend.zip",
			Header:   map[string][]string{"Content-Type": {"application/zip"}},
			Size:     123,
		},
	})
	if err != nil {
		t.Fatalf("Run() unexpected error: %v", err)
	}

	if files.lastAssetType != model.AssetTypeProjectZip {
		t.Fatalf("SaveUpload asset type = %q, want %q", files.lastAssetType, model.AssetTypeProjectZip)
	}
	if !zips.called {
		t.Fatal("expected ZipService.ExtractAndAnalyze to be called")
	}
	if ai.lastRequest.NeedVision {
		t.Fatal("code mode with project_zip should not request vision")
	}
	if !strings.Contains(ai.lastRequest.UserPrompt, "前端项目 ZIP 源码摘要") {
		t.Fatalf("AI prompt missing project ZIP summary marker: %s", ai.lastRequest.UserPrompt)
	}
	if !strings.Contains(ai.lastRequest.UserPrompt, "src/App.vue") {
		t.Fatalf("AI prompt missing scanned file path: %s", ai.lastRequest.UserPrompt)
	}
}

func TestUIReviewCodeModeRequiresCodeOrProjectZip(t *testing.T) {
	svc := NewUIReviewService(
		&fakeUIReviewAIService{},
		newFakeUIReviewReportService(),
		&fakeUIReviewFileService{},
		&fakeUIReviewZipService{},
		filepath.Join("tmp", "uploads"),
	)

	_, err := svc.Run(context.Background(), UIReviewFormInput{
		Title:      "缺少源码",
		ReviewMode: "code",
	})
	if err == nil {
		t.Fatal("Run() error = nil, want validation error")
	}
	if !strings.Contains(err.Error(), "code or project_zip") {
		t.Fatalf("Run() error = %q, want code or project_zip validation", err.Error())
	}
}

type fakeUIReviewAIService struct {
	result      *service.AIResult
	lastRequest service.AIRequest
}

func (f *fakeUIReviewAIService) GenerateJSON(ctx context.Context, input service.AIRequest) (*service.AIResult, error) {
	f.lastRequest = input
	if f.result != nil {
		return f.result, nil
	}
	return &service.AIResult{JSONText: `{"scores":[],"issues":[],"recommendations":[],"codex_prompt":""}`}, nil
}

type fakeUIReviewReportService struct {
	report *model.Report
}

func newFakeUIReviewReportService() *fakeUIReviewReportService {
	return &fakeUIReviewReportService{
		report: &model.Report{ID: "report-1", Title: "report"},
	}
}

func (f *fakeUIReviewReportService) CreateProcessingReport(ctx context.Context, toolType, title, inputMode string, inputData json.RawMessage) (*model.Report, error) {
	f.report.Title = title
	f.report.ToolType = toolType
	f.report.InputMode = inputMode
	f.report.InputJSON = datatypes.JSON(inputData)
	return f.report, nil
}

func (f *fakeUIReviewReportService) SucceedReport(ctx context.Context, id string, reportJSON json.RawMessage, summary string, totalScore *int, grade *string, generatedFiles []model.GeneratedFile) (*dto.ReportDTO, error) {
	return &dto.ReportDTO{ID: id, Title: f.report.Title, Summary: summary, ReportData: reportJSON}, nil
}

func (f *fakeUIReviewReportService) FailReport(ctx context.Context, id string, errorMessage string) error {
	return nil
}

func (f *fakeUIReviewReportService) FallbackReport(ctx context.Context, id string, reportJSON json.RawMessage, summary string) error {
	return nil
}

func (f *fakeUIReviewReportService) GetReport(ctx context.Context, id string) (*dto.ReportDTO, error) {
	return &dto.ReportDTO{ID: id, Title: f.report.Title}, nil
}

func (f *fakeUIReviewReportService) ListReports(ctx context.Context, query dto.ListReportsQuery) (*dto.PaginatedResponse[dto.ReportDTO], error) {
	return nil, errors.New("not implemented")
}

func (f *fakeUIReviewReportService) DeleteReport(ctx context.Context, id string) error {
	return errors.New("not implemented")
}

func (f *fakeUIReviewReportService) GetDashboardStats(ctx context.Context) (*dto.DashboardStatsDTO, error) {
	return nil, errors.New("not implemented")
}

type fakeUIReviewFileService struct {
	lastAssetType string
}

func (f *fakeUIReviewFileService) SaveUpload(ctx context.Context, reportID string, fileHeader *multipart.FileHeader, assetType string, allowedTypes map[string]string) (*model.ReportAsset, error) {
	f.lastAssetType = assetType
	return &model.ReportAsset{
		ReportID:     reportID,
		AssetType:    assetType,
		OriginalName: fileHeader.Filename,
		RelativePath: filepath.ToSlash(filepath.Join(reportID, "source", fileHeader.Filename)),
	}, nil
}

func (f *fakeUIReviewFileService) DeleteReportDir(uploadDir, tempDir, reportID string) error {
	return nil
}

func (f *fakeUIReviewFileService) ValidateFile(fileHeader *multipart.FileHeader, allowedTypes map[string]string) error {
	return nil
}

type fakeUIReviewZipService struct {
	called  bool
	summary *service.ProjectSummary
}

func (f *fakeUIReviewZipService) ExtractAndAnalyze(zipPath string, limits service.ZipLimits) (*service.ProjectSummary, error) {
	f.called = true
	if f.summary != nil {
		return f.summary, nil
	}
	return &service.ProjectSummary{}, nil
}

func (f *fakeUIReviewZipService) Cleanup(dir string) error {
	return nil
}
