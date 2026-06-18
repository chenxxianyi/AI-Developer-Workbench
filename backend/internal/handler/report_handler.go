package handler

import (
	"ai-developer-workbench/internal/dto"
	"ai-developer-workbench/internal/service"
	"ai-developer-workbench/internal/util"

	"github.com/gin-gonic/gin"
)

// ReportHandler handles report CRUD requests.
type ReportHandler struct {
	reportService service.ReportService
}

// NewReportHandler creates a new report handler.
func NewReportHandler(reportService service.ReportService) *ReportHandler {
	return &ReportHandler{reportService: reportService}
}

// ListReports handles GET /api/reports.
func (h *ReportHandler) ListReports(c *gin.Context) {
	var query dto.ListReportsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		util.BadRequest(c, "Invalid query parameters")
		return
	}

	result, err := h.reportService.ListReports(c.Request.Context(), query)
	if err != nil {
		util.InternalError(c, "Failed to list reports")
		return
	}

	util.SuccessResponse(c, result)
}

// GetReport handles GET /api/reports/:id.
func (h *ReportHandler) GetReport(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		util.BadRequest(c, "Report ID is required")
		return
	}

	report, err := h.reportService.GetReport(c.Request.Context(), id)
	if err != nil {
		util.NotFound(c, "Report not found")
		return
	}

	util.SuccessResponse(c, report)
}

// DeleteReport handles DELETE /api/reports/:id.
func (h *ReportHandler) DeleteReport(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		util.BadRequest(c, "Report ID is required")
		return
	}

	if err := h.reportService.DeleteReport(c.Request.Context(), id); err != nil {
		util.NotFound(c, "Report not found or deletion failed")
		return
	}

	util.SuccessResponse(c, nil)
}

// RegisterReportRoutes registers report routes.
func RegisterReportRoutes(r *gin.RouterGroup, reportService service.ReportService) {
	handler := NewReportHandler(reportService)
	r.GET("/reports", handler.ListReports)
	r.GET("/reports/:id", handler.GetReport)
	r.DELETE("/reports/:id", handler.DeleteReport)
}