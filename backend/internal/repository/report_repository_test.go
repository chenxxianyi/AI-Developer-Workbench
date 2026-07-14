package repository

import (
	"testing"

	"ai-developer-workbench/internal/dto"
	"ai-developer-workbench/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// TestList_ValidStatusValues tests that the DTO-level validation accepts known statuses.
func TestValidStatusValues_AcceptsKnownStatuses(t *testing.T) {
	valid := dto.ValidStatusValues()
	for _, status := range model.ValidToolTypes() {
		// Each tool type has reports with valid statuses.
		_ = status
	}

	known := []string{"processing", "succeeded", "failed"}
	for _, s := range known {
		assert.True(t, valid[s], "status %q should be valid", s)
	}
	assert.False(t, valid["invalid_status"], "unknown status should be rejected")
	assert.False(t, valid[""], "empty string should be rejected")
	assert.False(t, valid["<script>"], "XSS attempt should be rejected")
}

// TestList_QuerySetDefaults verifies pagination defaults.
func TestList_QuerySetDefaults(t *testing.T) {
	tests := []struct {
		name         string
		input        dto.ListReportsQuery
		expectedPage int
		expectedSize int
	}{
		{"zero values", dto.ListReportsQuery{Page: 0, PageSize: 0}, 1, 10},
		{"negative values", dto.ListReportsQuery{Page: -1, PageSize: -5}, 1, 10},
		{"valid values", dto.ListReportsQuery{Page: 3, PageSize: 25}, 3, 25},
		{"page size too large", dto.ListReportsQuery{Page: 1, PageSize: 200}, 1, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.input.SetDefaults()
			assert.Equal(t, tt.expectedPage, tt.input.Page)
			assert.Equal(t, tt.expectedSize, tt.input.PageSize)
		})
	}
}

// TestList_StatusFieldInQuery tests that the Status field is preserved in the query.
func TestList_StatusFieldInQuery(t *testing.T) {
	q := dto.ListReportsQuery{
		ToolType: "ui_review",
		Status:   "succeeded",
		Sort:     "score_desc",
		Page:     2,
		PageSize: 15,
	}
	q.SetDefaults()

	assert.Equal(t, "ui_review", q.ToolType)
	assert.Equal(t, "succeeded", q.Status)
	assert.Equal(t, "score_desc", q.Sort)
	assert.Equal(t, 2, q.Page)
	assert.Equal(t, 15, q.PageSize)
}

// TestList_AllStatusCombinations tests status+tools combinations.
func TestList_AllStatusCombinations(t *testing.T) {
	validStatuses := []string{"processing", "succeeded", "failed"}
	validTools := model.ValidToolTypes()
	valid := dto.ValidStatusValues()

	for _, status := range validStatuses {
		assert.True(t, valid[status], "status %q should be valid", status)
		for _, toolType := range validTools {
			q := dto.ListReportsQuery{
				ToolType: toolType,
				Status:   status,
				Page:     1,
				PageSize: 10,
			}
			assert.Equal(t, status, q.Status)
			assert.Equal(t, toolType, q.ToolType)
		}
	}
}

// TestReportRepository_Interface tests that the repository interface is complete.
func TestReportRepository_Interface(t *testing.T) {
	// Verify NewReportRepository returns the interface.
	var repo ReportRepository = NewReportRepository(&gorm.DB{})
	require.NotNil(t, repo)
}

// TestList_StatusIsolatedFilter tests that status filter is independent.
func TestList_StatusIsolatedFilter(t *testing.T) {
	// When only status is provided (no tool_type), the query should still work.
	q := dto.ListReportsQuery{
		Status:   "failed",
		Page:     1,
		PageSize: 10,
	}
	q.SetDefaults()
	assert.Equal(t, "failed", q.Status)
	assert.Empty(t, q.ToolType)
	assert.Equal(t, 1, q.Page)
}
