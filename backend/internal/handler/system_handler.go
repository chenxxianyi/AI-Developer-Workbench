package handler

import (
	"ai-developer-workbench/internal/config"
	"ai-developer-workbench/internal/dto"
	"ai-developer-workbench/internal/model"
	"ai-developer-workbench/internal/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SystemHandler handles system status requests.
type SystemHandler struct {
	cfg *config.Config
	db  *gorm.DB
}

// NewSystemHandler creates a new system handler.
func NewSystemHandler(cfg *config.Config, db ...*gorm.DB) *SystemHandler {
	h := &SystemHandler{cfg: cfg}
	if len(db) > 0 {
		h.db = db[0]
	}
	return h
}

// GetStatus handles GET /api/system/status.
func (h *SystemHandler) GetStatus(c *gin.Context) {
	provider := h.cfg.AI.Provider
	textModel := h.cfg.AI.Model
	visionModel := h.cfg.AI.VisionModel
	if h.db != nil {
		var preset model.ModelPreset
		if err := h.db.Where("is_default = ? AND status = ?", true, "active").Order("updated_at desc").First(&preset).Error; err == nil {
			provider = preset.Provider
			textModel = preset.Model
			visionModel = preset.VisionModel
		}
	}

	status := dto.SystemStatusDTO{
		Healthy:     true,
		Provider:    provider,
		TextModel:   textModel,
		VisionModel: visionModel,
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
func RegisterSystemRoutes(r *gin.RouterGroup, cfg *config.Config, db ...*gorm.DB) {
	handler := NewSystemHandler(cfg, db...)
	r.GET("/system/status", handler.GetStatus)
}
