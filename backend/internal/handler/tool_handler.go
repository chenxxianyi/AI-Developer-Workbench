package handler

import (
	"ai-developer-workbench/internal/dto"
	"ai-developer-workbench/internal/repository"
	"ai-developer-workbench/internal/util"

	"github.com/gin-gonic/gin"
)

// ToolHandler handles tool metadata requests.
type ToolHandler struct {
	reportRepo repository.ReportRepository
}

// NewToolHandler creates a new tool handler.
func NewToolHandler(reportRepo repository.ReportRepository) *ToolHandler {
	return &ToolHandler{reportRepo: reportRepo}
}

// GetTools handles GET /api/tools.
func (h *ToolHandler) GetTools(c *gin.Context) {
	// Get usage counts from database.
	stats, err := h.reportRepo.GetDashboardStats(c.Request.Context())
	if err != nil {
		util.InternalError(c, "Failed to get tool usage stats")
		return
	}

	// Build tool metadata with usage counts.
	tools := make([]dto.ToolMetaDTO, 0, len(dto.ToolMetaList))
	for _, meta := range dto.ToolMetaList {
		tool := meta
		if stats.ToolUsage != nil {
			tool.UsageCount = stats.ToolUsage[tool.ToolType]
		}
		tools = append(tools, tool)
	}

	util.SuccessResponse(c, tools)
}

// RegisterToolRoutes registers tool routes.
func RegisterToolRoutes(r *gin.RouterGroup, reportRepo repository.ReportRepository) {
	handler := NewToolHandler(reportRepo)
	r.GET("/tools", handler.GetTools)
}