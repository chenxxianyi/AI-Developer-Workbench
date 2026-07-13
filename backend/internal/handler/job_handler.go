package handler

import (
	"net/http"

	"ai-developer-workbench/internal/service"
	"ai-developer-workbench/internal/util"

	"github.com/gin-gonic/gin"
)

type JobHandler struct {
	jobService service.JobService
}

func NewJobHandler(jobService service.JobService) *JobHandler {
	return &JobHandler{jobService: jobService}
}

func RegisterJobRoutes(rg *gin.RouterGroup, jobService service.JobService) {
	h := NewJobHandler(jobService)
	rg.GET("/jobs/:id", h.GetJob)
	rg.POST("/jobs/:id/cancel", h.CancelJob)
	rg.POST("/jobs/:id/retry", h.RetryJob)
}

func (h *JobHandler) GetJob(c *gin.Context) {
	id := c.Param("id")
	job, err := h.jobService.GetJob(c.Request.Context(), id)
	if err != nil {
		util.ErrorResponse(c, http.StatusNotFound, util.CodeReportNotFound, "job not found")
		return
	}
	util.SuccessResponse(c, job)
}

func (h *JobHandler) CancelJob(c *gin.Context) {
	id := c.Param("id")
	if err := h.jobService.RequestCancel(c.Request.Context(), id); err != nil {
		util.ErrorResponse(c, http.StatusConflict, util.CodeStateConflict, err.Error())
		return
	}
	job, _ := h.jobService.GetJob(c.Request.Context(), id)
	util.SuccessResponse(c, job)
}

func (h *JobHandler) RetryJob(c *gin.Context) {
	id := c.Param("id")
	newJob, err := h.jobService.RetryJob(c.Request.Context(), id)
	if err != nil {
		util.ErrorResponse(c, http.StatusConflict, util.CodeStateConflict, err.Error())
		return
	}
	util.SuccessResponse(c, newJob)
}
