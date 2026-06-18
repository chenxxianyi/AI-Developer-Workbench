package handler

import (
	"net/http"

	"ai-developer-workbench/internal/dto"
	"ai-developer-workbench/internal/service/tools"
	"ai-developer-workbench/internal/util"

	"github.com/gin-gonic/gin"
)

// ToolRunHandler handles tool execution requests.
type ToolRunHandler struct {
	agentConfigService    *tools.AgentConfigService
	dbSchemaService       *tools.DBSchemaService
	uiReviewService       *tools.UIReviewService
	projectDoctorService  *tools.ProjectDoctorService
	apiDocService         *tools.APIDocService
}

// NewToolRunHandler creates a new tool run handler.
func NewToolRunHandler(
	agentConfigService *tools.AgentConfigService,
	dbSchemaService *tools.DBSchemaService,
	uiReviewService *tools.UIReviewService,
	projectDoctorService *tools.ProjectDoctorService,
	apiDocService *tools.APIDocService,
) *ToolRunHandler {
	return &ToolRunHandler{
		agentConfigService:    agentConfigService,
		dbSchemaService:       dbSchemaService,
		uiReviewService:       uiReviewService,
		projectDoctorService:  projectDoctorService,
		apiDocService:         apiDocService,
	}
}

// RunAgentConfig handles POST /api/tools/agent-config/run.
func (h *ToolRunHandler) RunAgentConfig(c *gin.Context) {
	var req dto.AgentConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	report, err := h.agentConfigService.Run(c.Request.Context(), req)
	if err != nil {
		util.ErrorResponse(c, http.StatusInternalServerError, util.CodeInternalError, err.Error())
		return
	}

	util.SuccessResponse(c, report)
}

// RunDBSchema handles POST /api/tools/db-schema/run.
func (h *ToolRunHandler) RunDBSchema(c *gin.Context) {
	var req dto.DBSchemaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	report, err := h.dbSchemaService.Run(c.Request.Context(), req)
	if err != nil {
		util.ErrorResponse(c, http.StatusInternalServerError, util.CodeInternalError, err.Error())
		return
	}

	util.SuccessResponse(c, report)
}

// RunUIReview handles POST /api/tools/ui-review/run.
func (h *ToolRunHandler) RunUIReview(c *gin.Context) {
	input := tools.UIReviewFormInput{
		Title:       c.PostForm("title"),
		ReviewMode:  c.PostForm("review_mode"),
		PageType:    c.PostForm("page_type"),
		TargetStyle: c.PostForm("target_style"),
		Description: c.PostForm("description"),
		Code:        c.PostForm("code"),
	}

	// Handle screenshot file.
	if file, err := c.FormFile("screenshot"); err == nil {
		input.Screenshot = file
	}

	report, err := h.uiReviewService.Run(c.Request.Context(), input)
	if err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, util.CodeInvalidRequest, err.Error())
		return
	}

	util.SuccessResponse(c, report)
}

// RunProjectDoctor handles POST /api/tools/project-doctor/run.
func (h *ToolRunHandler) RunProjectDoctor(c *gin.Context) {
	input := tools.ProjectDoctorFormInput{
		Title:              c.PostForm("title"),
		ProjectName:        c.PostForm("project_name"),
		TechStack:          c.PostForm("tech_stack"),
		ProjectDescription: c.PostForm("project_description"),
		AnalysisDepth:      c.PostForm("analysis_depth"),
	}

	if file, err := c.FormFile("project_zip"); err == nil {
		input.ProjectZip = file
	}

	report, err := h.projectDoctorService.Run(c.Request.Context(), input)
	if err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, util.CodeInvalidRequest, err.Error())
		return
	}

	util.SuccessResponse(c, report)
}

// RunAPIDoc handles POST /api/tools/api-doc/run.
func (h *ToolRunHandler) RunAPIDoc(c *gin.Context) {
	input := tools.APIDocFormInput{
		Title:          c.PostForm("title"),
		SourceType:     c.PostForm("source_type"),
		BackendStack:   c.PostForm("backend_stack"),
		Code:           c.PostForm("code"),
		APIDescription: c.PostForm("api_description"),
		OutputFormat:   c.PostForm("output_format"),
	}

	if file, err := c.FormFile("project_zip"); err == nil {
		input.ProjectZip = file
	}

	report, err := h.apiDocService.Run(c.Request.Context(), input)
	if err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, util.CodeInvalidRequest, err.Error())
		return
	}

	util.SuccessResponse(c, report)
}

// RegisterToolRunRoutes registers tool execution routes.
func RegisterToolRunRoutes(r *gin.RouterGroup, h *ToolRunHandler) {
	r.POST("/tools/agent-config/run", h.RunAgentConfig)
	r.POST("/tools/db-schema/run", h.RunDBSchema)
	r.POST("/tools/ui-review/run", h.RunUIReview)
	r.POST("/tools/project-doctor/run", h.RunProjectDoctor)
	r.POST("/tools/api-doc/run", h.RunAPIDoc)
}