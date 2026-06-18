package handler

import (
	"ai-developer-workbench/internal/service"
	"ai-developer-workbench/internal/util"

	"github.com/gin-gonic/gin"
)

// DashboardHandler handles dashboard statistics requests.
type DashboardHandler struct {
	reportService service.ReportService
}

// NewDashboardHandler creates a new dashboard handler.
func NewDashboardHandler(reportService service.ReportService) *DashboardHandler {
	return &DashboardHandler{reportService: reportService}
}

// GetStats handles GET /api/dashboard/stats.
func (h *DashboardHandler) GetStats(c *gin.Context) {
	stats, err := h.reportService.GetDashboardStats(c.Request.Context())
	if err != nil {
		util.InternalError(c, "Failed to get dashboard stats")
		return
	}

	util.SuccessResponse(c, stats)
}

// RegisterDashboardRoutes registers dashboard routes.
func RegisterDashboardRoutes(r *gin.RouterGroup, reportService service.ReportService) {
	handler := NewDashboardHandler(reportService)
	r.GET("/dashboard/stats", handler.GetStats)
}