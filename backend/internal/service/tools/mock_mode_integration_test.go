package tools

import (
	"context"
	"encoding/json"
	"mime/multipart"
	"testing"
	"time"

	"ai-developer-workbench/internal/dto"
	"ai-developer-workbench/internal/model"
	"ai-developer-workbench/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/datatypes"
)

// fakeReportService records calls to verify the full report lifecycle in mock mode.
type fakeReportService struct {
	lastModelReport    *model.Report
	lastReportDTO      *dto.ReportDTO
	createCalled       bool
	succeedCalled      bool
	failCalled         bool
	fallbackCalled     bool
	lastSummary        string
	lastScore          *int
	lastGrade          *string
	lastGeneratedFiles []model.GeneratedFile
	lastReportJSON     json.RawMessage
}

func newFakeReportService() *fakeReportService {
	return &fakeReportService{}
}

func (f *fakeReportService) CreateProcessingReport(_ context.Context, toolType, title, inputMode string, inputData json.RawMessage, parentReportID, projectID string) (*model.Report, error) {
	f.createCalled = true
	r := &model.Report{
		ID:        "mock-report-" + toolType,
		ToolType:  toolType,
		Title:     title,
		InputMode: inputMode,
		Status:    model.StatusProcessing,
		InputJSON: datatypes.JSON(inputData),
	}
	if parentReportID != "" {
		r.ParentReportID = &parentReportID
	}
	if projectID != "" {
		r.ProjectID = &projectID
	}
	f.lastModelReport = r
	return r, nil
}

func (f *fakeReportService) SucceedReport(_ context.Context, id string, reportJSON json.RawMessage, summary string, totalScore *int, grade *string, generatedFiles []model.GeneratedFile) (*dto.ReportDTO, error) {
	f.succeedCalled = true
	f.lastSummary = summary
	f.lastScore = totalScore
	f.lastGrade = grade
	f.lastGeneratedFiles = generatedFiles
	f.lastReportJSON = reportJSON

	gfDTOs := make([]dto.GeneratedFileDTO, len(generatedFiles))
	for i, gf := range generatedFiles {
		gfDTOs[i] = dto.GeneratedFileDTO{
			ID:       gf.ID,
			Filename: gf.Filename,
			Language: gf.Language,
			MimeType: gf.MimeType,
		}
	}

	r := &dto.ReportDTO{
		ID:             id,
		ToolType:       f.lastModelReport.ToolType,
		Title:          f.lastModelReport.Title,
		Status:         model.StatusSucceeded,
		Summary:        summary,
		TotalScore:     totalScore,
		Grade:          grade,
		ReportData:     reportJSON,
		GeneratedFiles: gfDTOs,
	}
	f.lastReportDTO = r
	return r, nil
}

func (f *fakeReportService) FailReport(_ context.Context, id string, errorMessage string) error {
	f.failCalled = true
	return nil
}

func (f *fakeReportService) FallbackReport(_ context.Context, id string, reportJSON json.RawMessage, summary string) error {
	f.fallbackCalled = true
	f.lastReportJSON = reportJSON
	return nil
}

func (f *fakeReportService) GetReport(_ context.Context, id string) (*dto.ReportDTO, error) {
	if f.lastReportDTO != nil {
		return f.lastReportDTO, nil
	}
	return nil, assert.AnError
}

func (f *fakeReportService) ListReports(_ context.Context, _ dto.ListReportsQuery) (*dto.PaginatedResponse[dto.ReportDTO], error) {
	return &dto.PaginatedResponse[dto.ReportDTO]{Items: nil, Total: 0}, nil
}

func (f *fakeReportService) DeleteReport(_ context.Context, _ string) error {
	return nil
}

func (f *fakeReportService) GetDashboardStats(_ context.Context) (*dto.DashboardStatsDTO, error) {
	return &dto.DashboardStatsDTO{}, nil
}

func (f *fakeReportService) ValidateParentReport(_ context.Context, toolType, parentReportID string) (*model.Report, error) {
	if parentReportID == "" {
		return nil, nil
	}
	// For testing, always return a valid parent report with matching tool type
	return &model.Report{
		ID:       parentReportID,
		ToolType: toolType,
		Title:    "Parent Report",
		Status:   model.StatusSucceeded,
	}, nil
}

func (f *fakeReportService) ResolveProject(_ context.Context, projectID string) (*model.Project, error) {
	if projectID == "" {
		return nil, nil
	}
	return &model.Project{ID: projectID, Name: "Test Project", FrontendStack: "Vue", BackendStack: "Go"}, nil
}

func (f *fakeReportService) CompareReports(_ context.Context, _, _ string) (*dto.ReportCompareDTO, error) {
	return &dto.ReportCompareDTO{}, nil
}

var _ service.ReportService = (*fakeReportService)(nil)

func TestMockMode_AgentConfig_EndToEnd(t *testing.T) {
	reportSvc := newFakeReportService()
	aiSvc := service.NewMockAIService()
	svc := NewAgentConfigService(aiSvc, reportSvc)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := dto.AgentConfigRequest{
		Title:         "Mock Agent Config E2E",
		ProjectName:   "TestProject",
		ProjectType:   "fullstack",
		FrontendStack: "Vue 3",
		BackendStack:  "Go/Gin",
		Database:      "MySQL",
	}

	report, err := svc.Run(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, report)

	assert.True(t, reportSvc.createCalled, "CreateProcessingReport should be called")
	assert.True(t, reportSvc.succeedCalled, "SucceedReport should be called")
	assert.False(t, reportSvc.failCalled, "FailReport should not be called")
	assert.False(t, reportSvc.fallbackCalled, "FallbackReport should not be called")

	assert.Equal(t, model.StatusSucceeded, report.Status)
	assert.Nil(t, report.TotalScore)
	assert.Nil(t, report.Grade)
	assert.NotEmpty(t, report.GeneratedFiles)
	assert.NotEmpty(t, report.Summary)

	var result dto.AgentConfigResult
	err = json.Unmarshal(report.ReportData, &result)
	require.NoError(t, err)
	assert.NotEmpty(t, result.GeneratedFilesContent)
	assert.NotEmpty(t, result.CodexPrompt)
}

func TestMockMode_DBSchema_EndToEnd(t *testing.T) {
	reportSvc := newFakeReportService()
	aiSvc := service.NewMockAIService()
	svc := NewDBSchemaService(aiSvc, reportSvc)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := dto.DBSchemaRequest{
		Title:           "Mock DB Schema E2E",
		SchemaType:      "review",
		DatabaseType:    "MySQL",
		SchemaContent:   "CREATE TABLE users (id INT PRIMARY KEY);",
		BusinessContext: "E-commerce",
	}

	report, err := svc.Run(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, report)

	assert.True(t, reportSvc.createCalled)
	assert.True(t, reportSvc.succeedCalled)
	assert.False(t, reportSvc.failCalled)

	assert.Equal(t, model.StatusSucceeded, report.Status)
	require.NotNil(t, report.TotalScore)
	assert.GreaterOrEqual(t, *report.TotalScore, 0)
	assert.LessOrEqual(t, *report.TotalScore, 100)
	require.NotNil(t, report.Grade)
	assert.NotEmpty(t, report.GeneratedFiles)

	var result dto.DBSchemaResult
	err = json.Unmarshal(report.ReportData, &result)
	require.NoError(t, err)
	assert.NotEmpty(t, result.Scores)
	assert.NotEmpty(t, result.Issues)
	assert.NotEmpty(t, result.CodexPrompt)
}

func TestMockMode_UIReview_CodeMode_EndToEnd(t *testing.T) {
	reportSvc := newFakeReportService()
	aiSvc := service.NewMockAIService()
	svc := NewUIReviewService(aiSvc, reportSvc, &fakeUIReviewFileService{}, &fakeUIReviewZipService{}, "/tmp/test")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	input := UIReviewFormInput{
		Title:       "Mock UI Review E2E",
		ReviewMode:  "code",
		CodeSource:  "paste",
		Code:        "<div class='container'><button>Click me</button></div>",
		PageType:    "dashboard",
		Description: "Test dashboard",
	}

	report, err := svc.Run(ctx, input)
	require.NoError(t, err)
	require.NotNil(t, report)

	assert.True(t, reportSvc.createCalled)
	assert.True(t, reportSvc.succeedCalled)

	require.NotNil(t, report.TotalScore)
	assert.GreaterOrEqual(t, *report.TotalScore, 0)
	assert.LessOrEqual(t, *report.TotalScore, 100)
	require.NotNil(t, report.Grade)

	var result dto.UIReviewResult
	err = json.Unmarshal(report.ReportData, &result)
	require.NoError(t, err)
	assert.NotEmpty(t, result.Scores)
	assert.NotEmpty(t, result.Issues)
	assert.NotEmpty(t, result.Recommendations)
	assert.NotEmpty(t, result.CodexPrompt)
}

func TestMockMode_ProjectDoctor_EndToEnd(t *testing.T) {
	reportSvc := newFakeReportService()
	aiSvc := service.NewMockAIService()
	svc := NewProjectDoctorService(aiSvc, reportSvc, &fakeUIReviewFileService{}, &fakeUIReviewZipService{}, "/tmp/test", "/tmp/temp")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	input := ProjectDoctorFormInput{
		Title:              "Mock Project Doctor E2E",
		ProjectName:        "TestProject",
		TechStack:          "Vue 3, Go, MySQL",
		ProjectDescription: "A full-stack web application",
		AnalysisDepth:      "standard",
		ProjectZip: &multipart.FileHeader{
			Filename: "test.zip",
			Header:   map[string][]string{"Content-Type": {"application/zip"}},
			Size:     1024,
		},
	}

	report, err := svc.Run(ctx, input)
	require.NoError(t, err)
	require.NotNil(t, report)

	assert.True(t, reportSvc.createCalled)
	assert.True(t, reportSvc.succeedCalled)

	require.NotNil(t, report.TotalScore)
	assert.GreaterOrEqual(t, *report.TotalScore, 0)
	assert.LessOrEqual(t, *report.TotalScore, 100)
	require.NotNil(t, report.Grade)

	var result dto.ProjectDoctorResult
	err = json.Unmarshal(report.ReportData, &result)
	require.NoError(t, err)
	assert.NotEmpty(t, result.Scores)
	assert.NotEmpty(t, result.Issues)
	assert.NotEmpty(t, result.CodexPrompt)
}

func TestMockMode_APIDoc_PasteMode_EndToEnd(t *testing.T) {
	reportSvc := newFakeReportService()
	aiSvc := service.NewMockAIService()
	svc := NewAPIDocService(aiSvc, reportSvc, &fakeUIReviewFileService{}, &fakeUIReviewZipService{}, "/tmp/test")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	input := APIDocFormInput{
		Title:          "Mock API Doc E2E",
		SourceType:     "code",
		BackendStack:   "Go/Gin",
		Code:           "// GET /api/users\nfunc ListUsers() {}",
		APIDescription: "User management API",
		OutputFormat:   "markdown",
	}

	report, err := svc.Run(ctx, input)
	require.NoError(t, err)
	require.NotNil(t, report)

	assert.True(t, reportSvc.createCalled)
	assert.True(t, reportSvc.succeedCalled)

	assert.Nil(t, report.TotalScore)
	assert.Nil(t, report.Grade)
	assert.NotEmpty(t, report.GeneratedFiles)

	var result dto.APIDocResult
	err = json.Unmarshal(report.ReportData, &result)
	require.NoError(t, err)
	assert.NotEmpty(t, result.Modules)
	assert.NotEmpty(t, result.CodexPrompt)
}

func TestMockMode_DeterministicOutput(t *testing.T) {
	reportSvc1 := newFakeReportService()
	svc1 := NewDBSchemaService(service.NewMockAIService(), reportSvc1)

	reportSvc2 := newFakeReportService()
	svc2 := NewDBSchemaService(service.NewMockAIService(), reportSvc2)

	req := dto.DBSchemaRequest{
		Title:         "Deterministic Test",
		SchemaType:    "review",
		DatabaseType:  "PostgreSQL",
		SchemaContent: "CREATE TABLE items (id SERIAL PRIMARY KEY, name TEXT);",
	}

	ctx := context.Background()

	report1, err := svc1.Run(ctx, req)
	require.NoError(t, err)

	report2, err := svc2.Run(ctx, req)
	require.NoError(t, err)

	assert.Equal(t, *report1.TotalScore, *report2.TotalScore, "scores should be identical")
	assert.Equal(t, string(report1.ReportData), string(report2.ReportData), "report data should be identical")
}

func TestMockMode_NonScoringToolsDoNotFakeScores(t *testing.T) {
	reportSvc := newFakeReportService()
	aiSvc := service.NewMockAIService()
	svc := NewAgentConfigService(aiSvc, reportSvc)

	ctx := context.Background()

	req := dto.AgentConfigRequest{
		Title:       "Non-scoring Test",
		ProjectName: "TestProject",
	}

	report, err := svc.Run(ctx, req)
	require.NoError(t, err)

	assert.Nil(t, report.TotalScore, "Agent Config should not have a score")
	assert.Nil(t, report.Grade, "Agent Config should not have a grade")
}

func TestMockMode_ScoringToolsHaveValidScores(t *testing.T) {
	reportSvc := newFakeReportService()
	svc := NewDBSchemaService(service.NewMockAIService(), reportSvc)

	ctx := context.Background()

	req := dto.DBSchemaRequest{
		Title:         "Scoring Test",
		SchemaType:    "review",
		DatabaseType:  "MySQL",
		SchemaContent: "CREATE TABLE t (id INT);",
	}

	report, err := svc.Run(ctx, req)
	require.NoError(t, err)

	require.NotNil(t, report.TotalScore, "DB Schema should have a score")
	assert.GreaterOrEqual(t, *report.TotalScore, 0, "score must be >= 0")
	assert.LessOrEqual(t, *report.TotalScore, 100, "score must be <= 100")

	validGrades := map[string]bool{"A": true, "B": true, "C": true, "D": true, "F": true}
	require.NotNil(t, report.Grade)
	assert.True(t, validGrades[*report.Grade], "grade must be A-F, got %q", *report.Grade)
}
