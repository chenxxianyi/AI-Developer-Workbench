package service

import (
	"encoding/json"
	"testing"

	"ai-developer-workbench/internal/dto"

	"github.com/stretchr/testify/assert"
)

func mkReportJSON(t *testing.T, issues []dto.IssueItem, actions []dto.ActionItem) json.RawMessage {
	t.Helper()
	type result struct {
		Issues      []dto.IssueItem  `json:"issues"`
		ActionItems []dto.ActionItem `json:"action_items"`
	}
	b, err := json.Marshal(result{Issues: issues, ActionItems: actions})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return b
}

func TestExtractIssues(t *testing.T) {
	raw := mkReportJSON(t, []dto.IssueItem{
		{Title: "A", Severity: "high", Category: "bug"},
		{Title: "B", Severity: "low", Category: "ui"},
	}, nil)
	got := extractIssues(raw)
	assert.Len(t, got, 2)
	assert.Equal(t, "A", got[0].Title)
}

func TestExtractActionItems(t *testing.T) {
	raw := mkReportJSON(t, nil, []dto.ActionItem{
		{ID: "fix-a", Title: "Fix A", Priority: "high", Effort: "small"},
	})
	got := extractActionItems(raw)
	assert.Len(t, got, 1)
	assert.Equal(t, "fix-a", got[0].ID)
}

func TestIssueCountDeltaFrom(t *testing.T) {
	baseline := []dto.IssueItem{
		{Severity: "high"}, {Severity: "high"}, {Severity: "medium"}, {Severity: "low"},
	}
	target := []dto.IssueItem{
		{Severity: "high"}, {Severity: "medium"}, {Severity: "medium"}, {Severity: "low"}, {Severity: "low"},
	}
	d := IssueCountDeltaFrom(baseline, target)
	assert.Equal(t, -1, d.High, "high: 1 - 2")
	assert.Equal(t, 1, d.Medium, "medium: 2 - 1")
	assert.Equal(t, 1, d.Low, "low: 2 - 1")
	assert.Equal(t, 1, d.Total, "total: 5 - 4")
}

func TestCompareIssues(t *testing.T) {
	baseline := []dto.IssueItem{
		{Title: "Fixed", Category: "bug", Severity: "high"},
		{Title: "Persists", Category: "ui", Severity: "medium"},
	}
	target := []dto.IssueItem{
		{Title: "Persists", Category: "ui", Severity: "medium"},
		{Title: "NewOne", Category: "perf", Severity: "low"},
	}
	c := compareIssues(baseline, target)
	assert.Len(t, c.Resolved, 1)
	assert.Equal(t, "Fixed", c.Resolved[0].Title)
	assert.Len(t, c.New, 1)
	assert.Equal(t, "NewOne", c.New[0].Title)
	assert.Len(t, c.Persist, 1)
	assert.Equal(t, "Persists", c.Persist[0].Title)
}

func TestCompareActionItems(t *testing.T) {
	baseline := []dto.ActionItem{
		{ID: "keep-1", Title: "Keep", Category: "a"},
		{ID: "gone-1", Title: "Gone", Category: "b"},
	}
	target := []dto.ActionItem{
		{ID: "keep-1", Title: "Keep", Category: "a"},
		{ID: "new-1", Title: "Fresh", Category: "c"},
	}
	d := compareActionItems(baseline, target)
	assert.Len(t, d.Resolved, 1)
	assert.Equal(t, "gone-1", d.Resolved[0].ID)
	assert.Len(t, d.New, 1)
	assert.Equal(t, "new-1", d.New[0].ID)
	assert.Len(t, d.Persist, 1)
	assert.Equal(t, "keep-1", d.Persist[0].ID)
}

func TestCompareActionItems_FallbackByCategoryAndTitle(t *testing.T) {
	// IDs differ but category+title match → should persist, not double-count.
	baseline := []dto.ActionItem{{ID: "old-id", Title: "Same", Category: "x"}}
	target := []dto.ActionItem{{ID: "new-id", Title: "Same", Category: "x"}}
	d := compareActionItems(baseline, target)
	assert.Empty(t, d.Resolved, "fallback match should not be resolved")
	assert.Empty(t, d.New, "fallback match should not be new")
	assert.Len(t, d.Persist, 1)
}

func TestGradeDelta(t *testing.T) {
	assert.Equal(t, "B → A", gradeDelta("B", "A"))
	assert.Equal(t, "A → A (unchanged)", gradeDelta("A", "A"))
}

func TestExtractIssues_EmptyOrInvalid(t *testing.T) {
	assert.Nil(t, extractIssues(nil))
	assert.Nil(t, extractIssues(json.RawMessage(`{}`)))
	assert.Nil(t, extractIssues(json.RawMessage(`not-json`)))
}
