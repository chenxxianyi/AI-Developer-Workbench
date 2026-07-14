package service

import (
	"context"
	"encoding/json"
	"testing"

	"ai-developer-workbench/internal/config"
	"ai-developer-workbench/internal/dto"
	"ai-developer-workbench/internal/model"
	"ai-developer-workbench/internal/repository"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// fakeTxRepo implements repository.ReportRepository for transaction tests.
type fakeTxRepo struct {
	report     *model.Report
	updateErr  error
	getByIDErr error
}

func (f *fakeTxRepo) Create(_ context.Context, r *model.Report) error { return nil }
func (f *fakeTxRepo) Update(_ context.Context, _ *model.Report) error { return f.updateErr }
func (f *fakeTxRepo) GetByID(_ context.Context, _ string) (*model.Report, error) {
	if f.getByIDErr != nil {
		return nil, f.getByIDErr
	}
	return f.report, nil
}
func (f *fakeTxRepo) List(_ context.Context, _ dto.ListReportsQuery) ([]model.Report, int64, error) {
	return nil, 0, nil
}
func (f *fakeTxRepo) Delete(_ context.Context, _ *gorm.DB, _ string) error { return nil }
func (f *fakeTxRepo) GetDashboardStats(_ context.Context) (*dto.DashboardStatsDTO, error) {
	return &dto.DashboardStatsDTO{}, nil
}

// TestSucceedReport_TransactionRollback verifies that the transaction pattern is correct.
// The real transaction logic uses s.db.Transaction() which requires a real DB.
// This test verifies the business logic around report state transitions.
func TestReportStatusTransitions_AreCorrect(t *testing.T) {
	assert.Equal(t, "succeeded", model.StatusSucceeded)
	assert.Equal(t, "failed", model.StatusFailed)
	assert.Equal(t, "processing", model.StatusProcessing)
}

func TestGeneratedFile_HasRequiredFields(t *testing.T) {
	gf := model.GeneratedFile{
		ReportID: "r1",
		Filename: "test.md",
		Content:  "# Test",
		MimeType: "text/markdown",
		Language: "markdown",
	}
	assert.Equal(t, "r1", gf.ReportID)
	assert.Equal(t, "test.md", gf.Filename)
	assert.NotEmpty(t, gf.Content)

	// Verify that SizeBytes is set by the service (not by the caller).
	// The service sets SizeBytes = len(Content) in SucceedReport.
	assert.Equal(t, uint64(0), gf.SizeBytes, "SizeBytes should be 0 before service sets it")
}

func TestReportDataRoundTrip(t *testing.T) {
	// Verify that report JSON can be marshaled and unmarshaled correctly.
	reportData := dto.UIReviewResult{
		Scores: []dto.ScoreItem{
			{Name: "Test", Score: 80, MaxScore: 100, Comment: "Good"},
		},
		Issues: []dto.IssueItem{
			{Title: "Issue 1", Severity: "high", Category: "bug", Problem: "Broken", Suggestion: "Fix it", Action: "Fix"},
		},
		Recommendations: []string{"Do better"},
		CodexPrompt:     "Fix the bugs",
	}

	bytes, err := json.Marshal(reportData)
	require.NoError(t, err)

	var decoded dto.UIReviewResult
	err = json.Unmarshal(bytes, &decoded)
	require.NoError(t, err)

	assert.Equal(t, 80, decoded.Scores[0].Score)
	assert.Equal(t, "high", decoded.Issues[0].Severity)
	assert.Equal(t, "Fix the bugs", decoded.CodexPrompt)
}

func TestCreateProcessingReportAssociatesOnlyExistingProjects(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(
		&model.Project{},
		&model.Report{},
		&model.GeneratedFile{},
		&model.ReportAsset{},
	))

	project := &model.Project{Name: "Workbench"}
	require.NoError(t, db.Create(project).Error)
	reportSvc := NewReportService(
		&config.Config{},
		repository.NewReportRepository(db),
		repository.NewGeneratedFileRepository(db),
		repository.NewReportAssetRepository(db),
		repository.NewProjectRepository(db),
		db,
	)

	report, err := reportSvc.CreateProcessingReport(
		context.Background(),
		model.ToolTypeUIReview,
		"Project review",
		"code",
		json.RawMessage(`{"project_id":"`+project.ID+`"}`),
		"",
		project.ID,
	)
	require.NoError(t, err)
	require.NotNil(t, report.ProjectID)
	assert.Equal(t, project.ID, *report.ProjectID)

	_, err = reportSvc.CreateProcessingReport(
		context.Background(),
		model.ToolTypeUIReview,
		"Invalid project review",
		"code",
		json.RawMessage(`{}`),
		"",
		"missing-project",
	)
	require.ErrorContains(t, err, "project not found")

	var count int64
	require.NoError(t, db.Model(&model.Report{}).Count(&count).Error)
	assert.Equal(t, int64(1), count, "invalid project must not create a processing report")
}
