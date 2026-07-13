package service

import (
	"context"
	"errors"
	"testing"

	"ai-developer-workbench/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/datatypes"
)

func TestExportGitHubIssues_WithActionItems(t *testing.T) {
	reportRepo := &stubReportRepo{
		report: &model.Report{
			ID:       "r-issues",
			Title:    "UI <script> Report",
			ToolType: model.ToolTypeUIReview,
			Status:   model.StatusSucceeded,
			ReportJSON: datatypes.JSON([]byte(`{
				"action_items": [
					{
						"id": "fix-upload",
						"title": "Fix upload",
						"priority": "high",
						"effort": "small",
						"category": "accessibility",
						"reason": "Upload area lacks keyboard support",
						"suggested_prompt": "Add Enter and Space support",
						"issue_title": "fix(ui): upload keyboard support",
						"issue_body": "## Acceptance\n- [ ] Keyboard works"
					}
				]
			}`)),
		},
	}
	exportSvc := NewExportService(reportRepo, &stubFileRepo{files: map[string]*model.GeneratedFile{}})

	content, filename, err := exportSvc.ExportGitHubIssues(context.Background(), "r-issues")
	require.NoError(t, err)

	md := string(content)
	assert.Equal(t, "ui_review_github_issues.md", filename)
	assert.Contains(t, md, "fix(ui): upload keyboard support")
	assert.Contains(t, md, "Suggested labels: `ui`, `accessibility`, `quality`")
	assert.Contains(t, md, "Report ID: `r-issues`")
	assert.Contains(t, md, "Action item ID: `fix-upload`")
	assert.Contains(t, md, "### Acceptance Criteria")
	assert.Contains(t, md, "- [ ] Relevant automated tests pass")
	assert.Contains(t, md, "- [ ] Keyboard works")
	assert.Contains(t, md, "UI &lt;script&gt; Report")
	assert.NotContains(t, md, "UI <script> Report")
}

func TestExportGitHubIssues_NoActionItems(t *testing.T) {
	reportRepo := &stubReportRepo{
		report: &model.Report{
			ID:         "r-empty",
			Title:      "Old Report",
			ToolType:   model.ToolTypeProjectDoctor,
			Status:     model.StatusSucceeded,
			ReportJSON: datatypes.JSON([]byte(`{"recommendations":["legacy"]}`)),
		},
	}
	exportSvc := NewExportService(reportRepo, &stubFileRepo{files: map[string]*model.GeneratedFile{}})

	_, _, err := exportSvc.ExportGitHubIssues(context.Background(), "r-empty")
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrNoActionItems))
}
