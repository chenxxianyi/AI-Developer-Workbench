package handler

import (
	"ai-developer-workbench/internal/service"
	"ai-developer-workbench/pkg/response"
	"ai-developer-workbench/pkg/sse"

	"github.com/gin-gonic/gin"
)

// TaskHandler handles task endpoints, including SSE streaming.
type TaskHandler struct {
	taskSvc *service.TaskService
	broker  *sse.Broker
}

func NewTaskHandler(taskSvc *service.TaskService, broker *sse.Broker) *TaskHandler {
	return &TaskHandler{taskSvc: taskSvc, broker: broker}
}

func (h *TaskHandler) Get(c *gin.Context) {
	task, err := h.taskSvc.Get(c.Param("id"))
	if err != nil {
		response.NotFound(c, "任务不存在")
		return
	}
	response.Success(c, task)
}

func (h *TaskHandler) Stream(c *gin.Context) {
	sse.StreamHandler(h.broker)(c)
}

func (h *TaskHandler) Cancel(c *gin.Context) {
	if err := h.taskSvc.Cancel(c.Param("id")); err != nil {
		response.BusinessError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *TaskHandler) Retry(c *gin.Context) {
	task, err := h.taskSvc.Retry(c.Param("id"))
	if err != nil {
		response.BusinessError(c, err.Error())
		return
	}
	response.Success(c, task)
}

func RegisterTaskRoutes(r *gin.RouterGroup, h *TaskHandler) {
	t := r.Group("/tasks")
	t.GET("/:id", h.Get)
	t.GET("/:id/stream", h.Stream)
	t.POST("/:id/retry", h.Retry)
	t.POST("/:id/cancel", h.Cancel)
}
