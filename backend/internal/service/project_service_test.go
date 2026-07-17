package service

import (
	"context"
	"testing"
	"time"

	"ai-developer-workbench/internal/dto"
	"ai-developer-workbench/internal/model"
	"ai-developer-workbench/internal/repository"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func newProjectServiceForTest(t *testing.T) (ProjectService, *gorm.DB) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(
		&model.Project{},
		&model.Report{},
		&model.GeneratedFile{},
		&model.ReportAsset{},
	))
	return NewProjectService(repository.NewProjectRepository(db)), db
}

func TestProjectServiceCreateValidatesAndNormalizesInput(t *testing.T) {
	svc, _ := newProjectServiceForTest(t)

	project, err := svc.Create(context.Background(), dto.ProjectCreateDTO{
		Name:        "  Workbench  ",
		Description: "  Local quality workspace  ",
		RepoURL:     "https://example.com/workbench",
	})
	require.NoError(t, err)
	assert.Equal(t, "Workbench", project.Name)
	assert.Equal(t, "Local quality workspace", project.Description)
	assert.Equal(t, "utility_app", project.ProjectType)

	_, err = svc.Create(context.Background(), dto.ProjectCreateDTO{
		Name:    "Invalid URL",
		RepoURL: "git@example.com:org/repo.git",
	})
	require.ErrorContains(t, err, "repo_url")

	_, err = svc.Create(context.Background(), dto.ProjectCreateDTO{
		Name:        "Too long",
		Description: string(make([]byte, dto.ProjectDescriptionMaxLength+1)),
	})
	require.ErrorContains(t, err, "description")

	_, err = svc.Create(context.Background(), dto.ProjectCreateDTO{Name: "Invalid type", ProjectType: "unknown"})
	require.ErrorContains(t, err, "project_type")
}

func TestProjectServiceListIncludesReportAggregates(t *testing.T) {
	svc, db := newProjectServiceForTest(t)
	ctx := context.Background()
	first, err := svc.Create(ctx, dto.ProjectCreateDTO{Name: "First"})
	require.NoError(t, err)
	second, err := svc.Create(ctx, dto.ProjectCreateDTO{Name: "Second"})
	require.NoError(t, err)

	scoreA, scoreB := 80, 100
	require.NoError(t, db.Create(&model.Report{
		ToolType:   model.ToolTypeUIReview,
		Title:      "A",
		InputMode:  "code",
		Status:     model.StatusSucceeded,
		ProjectID:  &first.ID,
		TotalScore: &scoreA,
		ReportJSON: datatypes.JSON([]byte(`{}`)),
	}).Error)
	require.NoError(t, db.Create(&model.Report{
		ToolType:   model.ToolTypeDBSchema,
		Title:      "B",
		InputMode:  "json",
		Status:     model.StatusSucceeded,
		ProjectID:  &first.ID,
		TotalScore: &scoreB,
		ReportJSON: datatypes.JSON([]byte(`{}`)),
	}).Error)

	list, err := svc.List(ctx, dto.ListProjectsQuery{Page: 1, PageSize: 10})
	require.NoError(t, err)
	require.Len(t, list.Items, 2)
	var firstSummary dto.ProjectSummaryDTO
	for _, item := range list.Items {
		if item.ID == first.ID {
			firstSummary = item
		}
	}
	assert.Equal(t, int64(2), firstSummary.ReportCount)
	require.NotNil(t, firstSummary.AverageScore)
	assert.Equal(t, 90.0, *firstSummary.AverageScore)

	search, err := svc.List(ctx, dto.ListProjectsQuery{Search: "Second", Page: 1, PageSize: 10})
	require.NoError(t, err)
	require.Len(t, search.Items, 1)
	assert.Equal(t, second.ID, search.Items[0].ID)
}

func TestProjectServiceDeleteDetachesReports(t *testing.T) {
	svc, db := newProjectServiceForTest(t)
	ctx := context.Background()
	project, err := svc.Create(ctx, dto.ProjectCreateDTO{Name: "Delete me"})
	require.NoError(t, err)
	report := &model.Report{
		ToolType:  model.ToolTypeProjectDoctor,
		Title:     "Historical report",
		InputMode: "project_zip",
		Status:    model.StatusSucceeded,
		ProjectID: &project.ID,
		ReportJSON: datatypes.JSON([]byte(`{
			"issues": [{"severity": "high"}]
		}`)),
	}
	require.NoError(t, db.Create(report).Error)

	result, err := svc.Delete(ctx, project.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(1), result.DetachedReportCount)

	var stored model.Report
	require.NoError(t, db.First(&stored, "id = ?", report.ID).Error)
	assert.Nil(t, stored.ProjectID)
	_, err = svc.Get(ctx, project.ID)
	require.Error(t, err)
}

func TestProjectServiceStatsIncludesTrendArtifactsAndHistory(t *testing.T) {
	svc, db := newProjectServiceForTest(t)
	ctx := context.Background()
	project, err := svc.Create(ctx, dto.ProjectCreateDTO{Name: "Analytics"})
	require.NoError(t, err)
	score := 88
	report := &model.Report{
		ToolType:   model.ToolTypeDBSchema,
		Title:      "Schema review",
		InputMode:  "json",
		Status:     model.StatusSucceeded,
		ProjectID:  &project.ID,
		TotalScore: &score,
		InputJSON: datatypes.JSON([]byte(`{
			"schema_type": "SQL",
			"schema_content": "CREATE TABLE reports (id CHAR(36) PRIMARY KEY);"
		}`)),
		ReportJSON: datatypes.JSON([]byte(`{
			"issues": [{"severity": "high"}, {"severity": "low"}]
		}`)),
		CreatedAt: time.Now().UTC(),
	}
	require.NoError(t, db.Create(report).Error)
	require.NoError(t, db.Create(&model.GeneratedFile{
		ReportID: report.ID,
		Filename: "migration.sql",
		MimeType: "text/x-sql",
		Language: "sql",
		Content:  "ALTER TABLE reports ADD INDEX idx_status(status);",
	}).Error)

	stats, err := svc.GetStats(ctx, project.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(1), stats.TotalReports)
	assert.Equal(t, int64(1), stats.HighSeverityCount)
	require.NotNil(t, stats.AverageScore)
	assert.Equal(t, 88.0, *stats.AverageScore)
	require.Len(t, stats.QualityTrend, 1)
	require.Len(t, stats.LatestArtifacts, 1)
	assert.Equal(t, "migration.sql", stats.LatestArtifacts[0].Filename)

	history, err := svc.ListReports(ctx, project.ID, dto.ListReportsQuery{Page: 1, PageSize: 10})
	require.NoError(t, err)
	require.Len(t, history.Items, 1)
	assert.Equal(t, report.ID, history.Items[0].ID)
	assert.JSONEq(t, `{
		"schema_type": "SQL",
		"schema_content": "CREATE TABLE reports (id CHAR(36) PRIMARY KEY);"
	}`, string(history.Items[0].InputData))
	assert.JSONEq(t, `{}`, string(history.Items[0].ReportData))
}
