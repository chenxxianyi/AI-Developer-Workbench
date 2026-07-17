package handler

import (
	"context"
	"net/http"
	"strings"

	"ai-developer-workbench/internal/model"
	"ai-developer-workbench/internal/service"
	"ai-developer-workbench/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BlueprintHandler handles blueprint CRUD endpoints.
type BlueprintHandler struct {
	db       *gorm.DB
	aiGenSvc blueprintGenerator
}

type blueprintGenerator interface {
	GenerateBlueprint(context.Context, string) (*service.BlueprintAIResult, error)
}

func NewBlueprintHandler(db *gorm.DB, aiGenSvc blueprintGenerator) *BlueprintHandler {
	return &BlueprintHandler{db: db, aiGenSvc: aiGenSvc}
}

func (h *BlueprintHandler) Get(c *gin.Context) {
	projectID := c.Param("id")
	var bp model.Blueprint
	if err := h.db.Where("project_id = ?", projectID).Order("version desc").First(&bp).Error; err != nil {
		response.NotFound(c, "蓝图不存在，请先生成")
		return
	}
	response.Success(c, bp)
}

func (h *BlueprintHandler) Save(c *gin.Context) {
	projectID := c.Param("id")
	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "请提供蓝图内容")
		return
	}

	var existing model.Blueprint
	h.db.Where("project_id = ?", projectID).Order("version desc").First(&existing)

	bp := model.Blueprint{
		ID:        uuid.New().String(),
		ProjectID: projectID,
		Content:   req.Content,
		Status:    "draft",
		Version:   existing.Version + 1,
	}
	if err := h.db.Create(&bp).Error; err != nil {
		response.InternalError(c, "保存蓝图失败")
		return
	}
	response.Created(c, bp)
}

func (h *BlueprintHandler) Confirm(c *gin.Context) {
	projectID := c.Param("id")
	var bp model.Blueprint
	if err := h.db.Where("project_id = ?", projectID).Order("version desc").First(&bp).Error; err != nil {
		response.NotFound(c, "蓝图不存在")
		return
	}
	if bp.Status == "superseded" {
		response.ValidationError(c, "需求已经更新，当前蓝图已失效，请根据最新需求重新生成蓝图")
		return
	}
	var requirement model.Requirement
	if err := h.db.Where("project_id = ?", projectID).Order("version desc").First(&requirement).Error; err != nil {
		response.ValidationError(c, "请先保存有效需求再确认蓝图")
		return
	}
	if err := service.ValidateBlueprintAgainstRequirements(bp.Content, requirement.Content); err != nil {
		response.ValidationError(c, "蓝图尚未满足确认条件: "+err.Error())
		return
	}
	h.db.Model(&model.Blueprint{}).
		Where("project_id = ? AND id <> ? AND status = ?", projectID, bp.ID, "confirmed").
		Update("status", "superseded")
	h.db.Model(&bp).Update("status", "confirmed")
	bp.Status = "confirmed"
	response.Success(c, bp)
}

func (h *BlueprintHandler) Generate(c *gin.Context) {
	projectID := c.Param("id")
	if h.aiGenSvc == nil {
		response.InternalError(c, "AI 蓝图生成服务未初始化")
		return
	}

	result, err := h.aiGenSvc.GenerateBlueprint(c.Request.Context(), projectID)
	if err != nil {
		if isAITimeoutError(err) {
			response.Error(c, http.StatusGatewayTimeout, "AI_TIMEOUT", "AI 模型响应超时。请重试；如果持续出现，请在 AI 模型设置中检查服务状态或适当提高超时时间")
			return
		}
		response.InternalError(c, "AI 蓝图生成失败: "+err.Error())
		return
	}
	content, err := service.MarshalBlueprintContent(result)
	if err != nil {
		response.InternalError(c, "蓝图序列化失败")
		return
	}

	var existing model.Blueprint
	h.db.Where("project_id = ?", projectID).Order("version desc").First(&existing)

	bp := model.Blueprint{
		ID:        uuid.New().String(),
		ProjectID: projectID,
		Content:   content,
		Status:    "generated",
		Version:   existing.Version + 1,
	}
	if err := h.db.Create(&bp).Error; err != nil {
		response.InternalError(c, "保存 AI 蓝图失败")
		return
	}
	response.Created(c, bp)
}

func isAITimeoutError(err error) bool {
	if err == nil {
		return false
	}
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "timeout") || strings.Contains(message, "deadline exceeded")
}

func RegisterBlueprintRoutes(r *gin.RouterGroup, h *BlueprintHandler) {
	bp := r.Group("/projects/:id/blueprint")
	bp.GET("", h.Get)
	bp.POST("/generate", h.Generate)
	bp.PUT("", h.Save)
	bp.POST("/confirm", h.Confirm)
}
