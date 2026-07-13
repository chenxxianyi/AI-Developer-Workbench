package handler

import (
	"ai-developer-workbench/internal/model"
	"ai-developer-workbench/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BlueprintHandler handles blueprint CRUD endpoints.
type BlueprintHandler struct {
	db *gorm.DB
}

func NewBlueprintHandler(db *gorm.DB) *BlueprintHandler {
	return &BlueprintHandler{db: db}
}

func (h *BlueprintHandler) Get(c *gin.Context) {
	projectID := c.Param("projectId")
	var bp model.Blueprint
	if err := h.db.Where("project_id = ?", projectID).Order("version desc").First(&bp).Error; err != nil {
		response.NotFound(c, "蓝图不存在，请先生成")
		return
	}
	response.Success(c, bp)
}

func (h *BlueprintHandler) Save(c *gin.Context) {
	projectID := c.Param("projectId")
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
	projectID := c.Param("projectId")
	var bp model.Blueprint
	if err := h.db.Where("project_id = ?", projectID).Order("version desc").First(&bp).Error; err != nil {
		response.NotFound(c, "蓝图不存在")
		return
	}
	h.db.Model(&bp).Update("status", "confirmed")
	response.Success(c, bp)
}

func (h *BlueprintHandler) Generate(c *gin.Context) {
	projectID := c.Param("projectId")
	// TODO: integrate AI blueprint generator
	bp := model.Blueprint{
		ID: uuid.New().String(), ProjectID: projectID,
		Content: `{"product_positioning":"AI生成的企业网站","pages":[{"name":"首页","route":"/"}],"tech_stack":"Vue 3 + Tailwind CSS"}`,
		Status: "generated", Version: 1,
	}
	h.db.Create(&bp)
	response.Created(c, bp)
}

func RegisterBlueprintRoutes(r *gin.RouterGroup, h *BlueprintHandler) {
	bp := r.Group("/projects/:projectId/blueprint")
	bp.GET("", h.Get)
	bp.POST("/generate", h.Generate)
	bp.PUT("", h.Save)
	bp.POST("/confirm", h.Confirm)
}
