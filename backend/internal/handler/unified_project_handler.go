package handler

import (
	"ai-developer-workbench/internal/model"
	"ai-developer-workbench/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UnifiedProjectHandler handles unified project CRUD (Workbench + Builder).
type UnifiedProjectHandler struct {
	db *gorm.DB
}

func NewUnifiedProjectHandler(db *gorm.DB) *UnifiedProjectHandler {
	return &UnifiedProjectHandler{db: db}
}

type createProjectReq struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func (h *UnifiedProjectHandler) Create(c *gin.Context) {
	var req createProjectReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "请填写项目名称")
		return
	}
	_ = c.GetString("user_id") // reserved for future user binding
	project := model.Project{
		ID: uuid.New().String(), Name: req.Name, Description: req.Description,
	}
	if err := h.db.Create(&project).Error; err != nil {
		response.InternalError(c, "创建项目失败")
		return
	}
	response.Created(c, project)
}

func (h *UnifiedProjectHandler) List(c *gin.Context) {
	userID := c.GetString("user_id")
	var projects []model.Project
	h.db.Where("user_id = ? AND deleted_at IS NULL", userID).Order("created_at desc").Find(&projects)
	response.Success(c, projects)
}

func (h *UnifiedProjectHandler) Get(c *gin.Context) {
	id := c.Param("projectId")
	var project model.Project
	if err := h.db.First(&project, "id = ?", id).Error; err != nil {
		response.NotFound(c, "项目不存在")
		return
	}
	response.Success(c, project)
}

func (h *UnifiedProjectHandler) Update(c *gin.Context) {
	id := c.Param("projectId")
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "无效的请求参数")
		return
	}
	delete(req, "id")
	delete(req, "user_id")
	if err := h.db.Model(&model.Project{}).Where("id = ?", id).Updates(req).Error; err != nil {
		response.InternalError(c, "更新失败")
		return
	}
	var project model.Project
	h.db.First(&project, "id = ?", id)
	response.Success(c, project)
}

func (h *UnifiedProjectHandler) Delete(c *gin.Context) {
	id := c.Param("projectId")
	if err := h.db.Where("id = ?", id).Delete(&model.Project{}).Error; err != nil {
		response.InternalError(c, "删除失败")
		return
	}
	response.Success(c, nil)
}

func RegisterUnifiedProjectRoutes(r *gin.RouterGroup, h *UnifiedProjectHandler) {
	p := r.Group("/projects")
	p.POST("", h.Create)
	p.GET("", h.List)
	p.GET("/:projectId", h.Get)
	p.PUT("/:projectId", h.Update)
	p.DELETE("/:projectId", h.Delete)
}
