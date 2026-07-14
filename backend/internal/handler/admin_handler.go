package handler

import (
	"ai-developer-workbench/internal/middleware"
	"ai-developer-workbench/internal/model"
	"ai-developer-workbench/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AdminHandler handles admin endpoints (models, prompts, users, projects).
type AdminHandler struct {
	db *gorm.DB
}

func NewAdminHandler(db *gorm.DB) *AdminHandler {
	return &AdminHandler{db: db}
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

// ── AI Models ──

func (h *AdminHandler) ListModels(c *gin.Context) {
	// Placeholder — AI models table not yet migrated
	response.Success(c, []gin.H{
		{"id": "1", "name": "GPT-4.1", "provider": "openai", "status": "active"},
		{"id": "2", "name": "Claude 3.5", "provider": "anthropic", "status": "active"},
	})
}

func (h *AdminHandler) UpdateModel(c *gin.Context) {
	response.Success(c, nil)
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
	admin.PUT("/models/:id", h.UpdateModel)

	admin.GET("/prompts", h.ListPrompts)
	admin.PUT("/prompts/:id", h.UpdatePrompt)
}
