package handler

import (
	"ai-developer-workbench/internal/model"
	"ai-developer-workbench/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RequirementHandler handles project requirement CRUD.
type RequirementHandler struct {
	db *gorm.DB
}

func NewRequirementHandler(db *gorm.DB) *RequirementHandler {
	return &RequirementHandler{db: db}
}

func (h *RequirementHandler) Get(c *gin.Context) {
	projectID := c.Param("projectId")
	var req model.Requirement
	if err := h.db.Where("project_id = ?", projectID).Order("version desc").First(&req).Error; err != nil {
		response.NotFound(c, "需求不存在")
		return
	}
	response.Success(c, req)
}

func (h *RequirementHandler) Save(c *gin.Context) {
	projectID := c.Param("projectId")
	var body struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.ValidationError(c, "请提供需求内容")
		return
	}

	var existing model.Requirement
	h.db.Where("project_id = ?", projectID).Order("version desc").First(&existing)

	req := model.Requirement{
		ID:        uuid.New().String(),
		ProjectID: projectID,
		Content:   body.Content,
		Version:   existing.Version + 1,
	}
	if err := h.db.Create(&req).Error; err != nil {
		response.InternalError(c, "保存需求失败")
		return
	}
	response.Created(c, req)
}

func RegisterRequirementRoutes(r *gin.RouterGroup, h *RequirementHandler) {
	req := r.Group("/projects/:id/requirements")
	req.GET("", h.Get)
	req.PUT("", h.Save)
}
