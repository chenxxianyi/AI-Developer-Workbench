package handler

import (
	"errors"

	"ai-developer-workbench/internal/dto"
	"ai-developer-workbench/internal/service"
	"ai-developer-workbench/internal/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ProjectHandler handles project CRUD requests.
type ProjectHandler struct {
	projectService service.ProjectService
}

// NewProjectHandler creates a new project handler.
func NewProjectHandler(projectService service.ProjectService) *ProjectHandler {
	return &ProjectHandler{projectService: projectService}
}

// Create handles POST /api/projects.
func (h *ProjectHandler) Create(c *gin.Context) {
	var input dto.ProjectCreateDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		util.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	project, err := h.projectService.Create(c.Request.Context(), input)
	if err != nil {
		util.BadRequest(c, err.Error())
		return
	}
	util.SuccessResponse(c, project)
}

// List handles GET /api/projects.
func (h *ProjectHandler) List(c *gin.Context) {
	var query dto.ListProjectsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		util.BadRequest(c, "Invalid query parameters")
		return
	}
	result, err := h.projectService.List(c.Request.Context(), query)
	if err != nil {
		util.InternalError(c, "Failed to list projects")
		return
	}
	util.SuccessResponse(c, result)
}

// Get handles GET /api/projects/:id.
func (h *ProjectHandler) Get(c *gin.Context) {
	id := c.Param("id")
	project, err := h.projectService.Get(c.Request.Context(), id)
	if err != nil {
		util.NotFound(c, "Project not found")
		return
	}
	util.SuccessResponse(c, project)
}

// Update handles PATCH /api/projects/:id.
func (h *ProjectHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var input dto.ProjectUpdateDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		util.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	project, err := h.projectService.Update(c.Request.Context(), id, input)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			util.NotFound(c, "Project not found")
			return
		}
		if containsValidation(err.Error()) {
			util.BadRequest(c, err.Error())
			return
		}
		util.InternalError(c, "Failed to update project")
		return
	}
	util.SuccessResponse(c, project)
}

// Delete handles DELETE /api/projects/:id.
func (h *ProjectHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	result, err := h.projectService.Delete(c.Request.Context(), id)
	if err != nil {
		util.NotFound(c, "Project not found")
		return
	}
	util.SuccessResponse(c, result)
}

// GetStats handles GET /api/projects/:id/stats.
func (h *ProjectHandler) GetStats(c *gin.Context) {
	id := c.Param("id")
	stats, err := h.projectService.GetStats(c.Request.Context(), id)
	if err != nil {
		util.NotFound(c, "Project not found")
		return
	}
	util.SuccessResponse(c, stats)
}

// ListReports handles GET /api/projects/:id/reports.
func (h *ProjectHandler) ListReports(c *gin.Context) {
	var query dto.ListReportsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		util.BadRequest(c, "Invalid query parameters")
		return
	}
	result, err := h.projectService.ListReports(c.Request.Context(), c.Param("id"), query)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			util.NotFound(c, "Project not found")
			return
		}
		util.InternalError(c, "Failed to list project reports")
		return
	}
	util.SuccessResponse(c, result)
}

// RegisterProjectRoutes registers project routes.
func RegisterProjectRoutes(r *gin.RouterGroup, projectService service.ProjectService) {
	h := NewProjectHandler(projectService)
	r.POST("/projects", h.Create)
	r.GET("/projects", h.List)
	r.GET("/projects/:id", h.Get)
	r.PATCH("/projects/:id", h.Update)
	r.DELETE("/projects/:id", h.Delete)
	r.GET("/projects/:id/stats", h.GetStats)
	r.GET("/projects/:id/reports", h.ListReports)
}

// containsValidation is a tiny helper to surface field-level errors as 400.
func containsValidation(msg string) bool {
	for _, k := range []string{"must be", "cannot be empty", "is required"} {
		if containsStr(msg, k) {
			return true
		}
	}
	return false
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || indexOf(s, sub) >= 0)
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
