package handler

import (
	"strings"

	"ai-developer-workbench/internal/config"
	"ai-developer-workbench/internal/middleware"
	"ai-developer-workbench/internal/model"
	"ai-developer-workbench/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AdminHandler handles admin endpoints (models, prompts, users, projects).
type AdminHandler struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewAdminHandler(db *gorm.DB, cfg *config.Config) *AdminHandler {
	return &AdminHandler{db: db, cfg: cfg}
}

// ── Users ──

func (h *AdminHandler) ListUsers(c *gin.Context) {
	var users []model.User
	h.db.Order("created_at desc").Find(&users)
	response.Success(c, users)
}

func (h *AdminHandler) GetUser(c *gin.Context) {
	var user model.User
	if err := h.db.First(&user, "id = ?", c.Param("userId")).Error; err != nil {
		response.NotFound(c, "用户不存在")
		return
	}
	response.Success(c, user)
}

func (h *AdminHandler) UpdateUserStatus(c *gin.Context) {
	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "请提供账号状态")
		return
	}
	if req.Status != "active" && req.Status != "disabled" {
		response.ValidationError(c, "账号状态只能是 active 或 disabled")
		return
	}
	if c.Param("userId") == middleware.GetUserID(c) && req.Status == "disabled" {
		response.ValidationError(c, "不能停用当前登录账号")
		return
	}
	result := h.db.Model(&model.User{}).Where("id = ?", c.Param("userId")).Update("status", req.Status)
	if result.Error != nil {
		response.InternalError(c, "更新用户状态失败")
		return
	}
	if result.RowsAffected == 0 {
		response.NotFound(c, "用户不存在")
		return
	}
	response.Success(c, nil)
}

// ── Projects ──

func (h *AdminHandler) ListProjects(c *gin.Context) {
	var projects []model.Project
	h.db.Order("created_at desc").Find(&projects)
	response.Success(c, projects)
}

// AI Models

type modelPresetRequest struct {
	Name           string `json:"name"`
	Provider       string `json:"provider"`
	BaseURL        string `json:"base_url"`
	Model          string `json:"model"`
	VisionModel    string `json:"vision_model"`
	TimeoutSeconds int    `json:"timeout_seconds"`
	MaxRetries     int    `json:"max_retries"`
	Status         string `json:"status"`
	IsDefault      bool   `json:"is_default"`
}

func (h *AdminHandler) ListModels(c *gin.Context) {
	if err := h.ensureConfiguredModelPreset(); err != nil {
		response.InternalError(c, "failed to initialize model presets")
		return
	}

	var models []model.ModelPreset
	if err := h.db.Order("is_default desc, updated_at desc").Find(&models).Error; err != nil {
		response.InternalError(c, "failed to load model presets")
		return
	}
	response.Success(c, models)
}

func (h *AdminHandler) CreateModel(c *gin.Context) {
	var req modelPresetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "invalid model preset payload")
		return
	}
	preset, ok := h.modelPresetFromRequest(req, nil)
	if !ok {
		response.ValidationError(c, "name, provider, base URL, and model ID are required")
		return
	}
	if err := h.saveModelPreset(preset); err != nil {
		response.InternalError(c, "failed to save model preset")
		return
	}
	response.Created(c, preset)
}

func (h *AdminHandler) UpdateModel(c *gin.Context) {
	var existing model.ModelPreset
	if err := h.db.First(&existing, "id = ?", c.Param("id")).Error; err != nil {
		response.NotFound(c, "model preset not found")
		return
	}

	var req modelPresetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "invalid model preset payload")
		return
	}
	preset, ok := h.modelPresetFromRequest(req, &existing)
	if !ok {
		response.ValidationError(c, "name, provider, base URL, and model ID are required")
		return
	}
	if err := h.saveModelPreset(preset); err != nil {
		response.InternalError(c, "failed to save model preset")
		return
	}
	response.Success(c, preset)
}

func (h *AdminHandler) DeleteModel(c *gin.Context) {
	var existing model.ModelPreset
	if err := h.db.First(&existing, "id = ?", c.Param("id")).Error; err != nil {
		response.NotFound(c, "model preset not found")
		return
	}
	if existing.IsDefault {
		response.ValidationError(c, "default model cannot be deleted; switch default model first")
		return
	}
	if err := h.db.Delete(&existing).Error; err != nil {
		response.InternalError(c, "failed to save model preset")
		return
	}
	response.Success(c, nil)
}

func (h *AdminHandler) ensureConfiguredModelPreset() error {
	if h.cfg == nil {
		return nil
	}
	name := strings.TrimSpace(h.cfg.AI.Model)
	if name == "" {
		name = "Current AI Model"
	}
	var existing model.ModelPreset
	err := h.db.Where("provider = ? AND base_url = ? AND model = ?", h.cfg.AI.Provider, h.cfg.AI.BaseURL, h.cfg.AI.Model).First(&existing).Error
	if err == nil {
		updates := map[string]interface{}{
			"vision_model":    h.cfg.AI.VisionModel,
			"timeout_seconds": h.cfg.AI.TimeoutSeconds,
			"max_retries":     h.cfg.AI.MaxRetries,
			"status":          "active",
		}
		return h.db.Model(&existing).Updates(updates).Error
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}

	var count int64
	if err := h.db.Model(&model.ModelPreset{}).Count(&count).Error; err != nil {
		return err
	}
	preset := &model.ModelPreset{
		Name:           name,
		Provider:       h.cfg.AI.Provider,
		BaseURL:        h.cfg.AI.BaseURL,
		Model:          h.cfg.AI.Model,
		VisionModel:    h.cfg.AI.VisionModel,
		TimeoutSeconds: h.cfg.AI.TimeoutSeconds,
		MaxRetries:     h.cfg.AI.MaxRetries,
		Status:         "active",
		IsDefault:      count == 0,
	}
	return h.db.Create(preset).Error
}

func (h *AdminHandler) modelPresetFromRequest(req modelPresetRequest, existing *model.ModelPreset) (*model.ModelPreset, bool) {
	preset := &model.ModelPreset{}
	if existing != nil {
		*preset = *existing
	}

	name := strings.TrimSpace(req.Name)
	provider := strings.TrimSpace(req.Provider)
	baseURL := strings.TrimRight(strings.TrimSpace(req.BaseURL), "/")
	modelID := strings.TrimSpace(req.Model)
	visionModel := strings.TrimSpace(req.VisionModel)
	status := strings.TrimSpace(req.Status)
	if status == "" {
		status = "active"
	}
	if status != "active" && status != "disabled" {
		return nil, false
	}
	if name == "" || provider == "" || baseURL == "" || modelID == "" {
		return nil, false
	}
	if visionModel == "" {
		visionModel = modelID
	}
	preset.Name = name
	preset.Provider = provider
	preset.BaseURL = baseURL
	preset.Model = modelID
	preset.VisionModel = visionModel
	preset.Status = status
	preset.IsDefault = req.IsDefault
	preset.TimeoutSeconds = req.TimeoutSeconds
	if preset.TimeoutSeconds <= 0 {
		preset.TimeoutSeconds = 180
	}
	preset.MaxRetries = req.MaxRetries
	if preset.MaxRetries < 0 {
		preset.MaxRetries = 0
	}
	return preset, true
}

func (h *AdminHandler) saveModelPreset(preset *model.ModelPreset) error {
	return h.db.Transaction(func(tx *gorm.DB) error {
		if preset.IsDefault {
			if err := tx.Model(&model.ModelPreset{}).Where("id <> ?", preset.ID).Update("is_default", false).Error; err != nil {
				return err
			}
			preset.Status = "active"
		}
		return tx.Save(preset).Error
	})
}

// ── Prompts ──

func (h *AdminHandler) ListPrompts(c *gin.Context) {
	response.Success(c, []gin.H{
		{"id": "1", "name": "UI Review Default", "type": "ui_review", "status": "active", "version": 1},
	})
}

func (h *AdminHandler) UpdatePrompt(c *gin.Context) {
	response.Success(c, nil)
}

// RegisterAdminRoutes registers admin routes (requires admin role).
func RegisterAdminRoutes(r *gin.RouterGroup, h *AdminHandler) {
	admin := r.Group("/admin")

	admin.GET("/users", h.ListUsers)
	admin.GET("/users/:userId", h.GetUser)
	admin.PUT("/users/:userId/status", h.UpdateUserStatus)

	admin.GET("/projects", h.ListProjects)

	admin.GET("/models", h.ListModels)
	admin.POST("/models", h.CreateModel)
	admin.PUT("/models/:id", h.UpdateModel)
	admin.DELETE("/models/:id", h.DeleteModel)

	admin.GET("/prompts", h.ListPrompts)
	admin.PUT("/prompts/:id", h.UpdatePrompt)
}
