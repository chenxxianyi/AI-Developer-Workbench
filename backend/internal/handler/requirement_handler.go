package handler

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"ai-developer-workbench/internal/dto"
	"ai-developer-workbench/internal/model"
	"ai-developer-workbench/internal/service"
	"ai-developer-workbench/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RequirementHandler handles project requirement CRUD.
type RequirementHandler struct {
	db        *gorm.DB
	assistSvc *service.RequirementAssistService
}

func NewRequirementHandler(db *gorm.DB, assistSvc ...*service.RequirementAssistService) *RequirementHandler {
	var assistant *service.RequirementAssistService
	if len(assistSvc) > 0 {
		assistant = assistSvc[0]
	}
	return &RequirementHandler{db: db, assistSvc: assistant}
}

func (h *RequirementHandler) Assist(c *gin.Context) {
	if h.assistSvc == nil {
		response.InternalError(c, "需求整理服务未初始化")
		return
	}
	var input service.RequirementAssistInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.ValidationError(c, "请提供产品描述和当前需求")
		return
	}
	if len(input.Description) > 12000 {
		response.ValidationError(c, "产品描述不能超过 12000 个字符")
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 90*time.Second)
	defer cancel()
	result, err := h.assistSvc.Assist(ctx, c.Param("id"), input)
	if err != nil {
		response.BusinessError(c, err.Error())
		return
	}
	response.Success(c, result)
}

func (h *RequirementHandler) Get(c *gin.Context) {
	projectID := c.Param("id")
	var req model.Requirement
	if err := h.db.Where("project_id = ?", projectID).Order("version desc").First(&req).Error; err != nil {
		response.NotFound(c, "需求不存在")
		return
	}
	response.Success(c, req)
}

func (h *RequirementHandler) Save(c *gin.Context) {
	projectID := c.Param("id")
	var body struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.ValidationError(c, "请提供需求内容")
		return
	}
	if len(body.Content) > 64*1024 || !json.Valid([]byte(body.Content)) {
		response.ValidationError(c, "需求内容必须是有效且不超过 64KB 的 JSON")
		return
	}
	var spec struct {
		SchemaVersion      int      `json:"schema_version"`
		AppType            string   `json:"app_type"`
		Goal               string   `json:"goal"`
		TargetUsers        []string `json:"target_users"`
		MustHaveFeatures   []string `json:"must_have_features"`
		AcceptanceCriteria []string `json:"acceptance_criteria"`
	}
	if err := json.Unmarshal([]byte(body.Content), &spec); err != nil {
		response.ValidationError(c, "无法解析需求内容")
		return
	}
	if spec.SchemaVersion >= 2 && (strings.TrimSpace(spec.Goal) == "" || len(spec.TargetUsers) == 0 || len(spec.MustHaveFeatures) == 0 || len(spec.AcceptanceCriteria) == 0) {
		response.ValidationError(c, "结构化需求缺少目标、目标用户、必须功能或验收标准")
		return
	}
	if spec.SchemaVersion >= 2 && !dto.ValidProjectTypes[spec.AppType] {
		response.ValidationError(c, "结构化需求包含不支持的项目类型")
		return
	}

	var existing model.Requirement
	h.db.Where("project_id = ?", projectID).Order("version desc").First(&existing)
	if existing.ID != "" && equivalentJSON(existing.Content, body.Content) {
		response.Success(c, existing)
		return
	}

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
	if spec.SchemaVersion >= 2 {
		h.db.Model(&model.Project{}).Where("id = ?", projectID).Update("project_type", spec.AppType)
	}
	// A requirement change invalidates previously confirmed generation input.
	h.db.Model(&model.Blueprint{}).
		Where("project_id = ? AND status = ?", projectID, "confirmed").
		Update("status", "superseded")
	response.Created(c, req)
}

func equivalentJSON(left, right string) bool {
	var leftValue, rightValue any
	if json.Unmarshal([]byte(left), &leftValue) != nil || json.Unmarshal([]byte(right), &rightValue) != nil {
		return left == right
	}
	leftBytes, leftErr := json.Marshal(leftValue)
	rightBytes, rightErr := json.Marshal(rightValue)
	return leftErr == nil && rightErr == nil && string(leftBytes) == string(rightBytes)
}

func RegisterRequirementRoutes(r *gin.RouterGroup, h *RequirementHandler) {
	req := r.Group("/projects/:id/requirements")
	req.GET("", h.Get)
	req.PUT("", h.Save)
	req.POST("/assist", h.Assist)
}
