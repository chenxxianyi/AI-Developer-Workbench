package handler

import (
	"errors"

	"ai-developer-workbench/internal/service"
	"ai-developer-workbench/internal/util"

	"github.com/gin-gonic/gin"
)

// ExportHandler handles report export requests.
type ExportHandler struct {
	exportService service.ExportService
}

// NewExportHandler creates a new export handler.
func NewExportHandler(exportService service.ExportService) *ExportHandler {
	return &ExportHandler{exportService: exportService}
}

// ExportMarkdown handles GET /api/reports/:id/export?format=markdown.
func (h *ExportHandler) ExportMarkdown(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		util.BadRequest(c, "Report ID is required")
		return
	}

	format := c.DefaultQuery("format", "markdown")
	var (
		content  []byte
		filename string
		err      error
	)
	switch format {
	case "markdown":
		content, filename, err = h.exportService.ExportMarkdown(c.Request.Context(), id)
	case "github-issues":
		content, filename, err = h.exportService.ExportGitHubIssues(c.Request.Context(), id)
	default:
		util.BadRequest(c, "Unsupported export format")
		return
	}
	if err != nil {
		if errors.Is(err, service.ErrNoActionItems) {
			util.BadRequest(c, "Report has no action items")
			return
		}
		util.NotFound(c, "Report not found")
		return
	}

	util.WriteDownloadResponse(c, filename, content, "text/markdown")
}

// DownloadFile handles GET /api/reports/:id/files/:filename.
func (h *ExportHandler) DownloadFile(c *gin.Context) {
	id := c.Param("id")
	filename := c.Param("filename")
	if id == "" || filename == "" {
		util.BadRequest(c, "Report ID and filename are required")
		return
	}

	content, name, mimeType, err := h.exportService.GetFileContent(c.Request.Context(), id, filename)
	if err != nil {
		util.NotFound(c, "File not found")
		return
	}

	util.WriteDownloadResponse(c, name, content, mimeType)
}

// RegisterExportRoutes registers export routes.
func RegisterExportRoutes(r *gin.RouterGroup, exportService service.ExportService) {
	handler := NewExportHandler(exportService)
	r.GET("/reports/:id/export", handler.ExportMarkdown)
	r.GET("/reports/:id/files/:filename", handler.DownloadFile)
}
