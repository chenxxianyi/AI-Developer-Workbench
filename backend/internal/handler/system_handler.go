package handler

import (
	"ai-developer-workbench/internal/config"
	"ai-developer-workbench/internal/dto"
	"ai-developer-workbench/internal/util"

	"github.com/gin-gonic/gin"
)

// SystemHandler handles system status requests.
type SystemHandler struct {
	cfg *config.Config
}

// NewSystemHandler creates a new system handler.
func NewSystemHandler(cfg *config.Config) *SystemHandler {
	return &SystemHandler{cfg: cfg}
}

// GetStatus handles GET /api/system/status.
func (h *SystemHandler) GetStatus(c *gin.Context) {
	status := dto.SystemStatusDTO{
		Healthy:     true,
		Provider:    h.cfg.AI.Provider,
		TextModel:   h.cfg.AI.Model,
		VisionModel: h.cfg.AI.VisionModel,
		UploadLimits: dto.UploadLimitsDTO{
			ImageMaxBytes:    int64(h.cfg.Upload.MaxUploadSizeMB) * 1024 * 1024,
			ZipMaxBytes:      int64(h.cfg.Upload.MaxUploadSizeMB) * 1024 * 1024,
			ZipMaxFiles:      h.cfg.Upload.MaxProjectFiles,
			ZipMaxTotalBytes: int64(h.cfg.Upload.MaxZipUncompressedMB) * 1024 * 1024,
		},
	}

	util.SuccessResponse(c, status)
}

// RegisterSystemRoutes registers system routes.
func RegisterSystemRoutes(r *gin.RouterGroup, cfg *config.Config) {
	handler := NewSystemHandler(cfg)
	r.GET("/system/status", handler.GetStatus)
}