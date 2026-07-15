package handler

import (
	"context"
	"time"

	"ai-developer-workbench/internal/service"
	"ai-developer-workbench/pkg/response"
	"ai-developer-workbench/pkg/sse"

	"github.com/gin-gonic/gin"
)

// TaskHandler handles task endpoints, including SSE streaming.
type TaskHandler struct {
	taskSvc  *service.TaskService
	broker   *sse.Broker
	ws       *service.WorkspaceService
	aiGenSvc *service.AIGenerationService
}

func NewTaskHandler(taskSvc *service.TaskService, broker *sse.Broker, ws *service.WorkspaceService, aiGenSvc *service.AIGenerationService) *TaskHandler {
	return &TaskHandler{taskSvc: taskSvc, broker: broker, ws: ws, aiGenSvc: aiGenSvc}
}

func (h *TaskHandler) Create(c *gin.Context) {
	var req struct {
		ProjectID string `json:"project_id" binding:"required"`
		Type      string `json:"type" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "请提供项目 ID 和任务类型")
		return
	}

	userID := ""
	if value, ok := c.Get("user_id"); ok {
		if id, ok := value.(string); ok {
			userID = id
		}
	}

	task, err := h.taskSvc.Create(req.ProjectID, userID, req.Type)
	if err != nil {
		response.InternalError(c, "创建任务失败")
		return
	}

	if req.Type == "generation" {
		go h.runGenerationTask(task.ID, req.ProjectID)
	}

	response.Created(c, task)
}

func (h *TaskHandler) runGenerationTask(taskID, projectID string) {
	time.Sleep(150 * time.Millisecond)
	if err := h.taskSvc.Start(taskID); err != nil {
		h.taskSvc.PublishFailed(taskID, err.Error())
		return
	}

	h.taskSvc.PublishProgress(taskID, "prepare", 10, "正在准备项目工作区")
	if h.aiGenSvc == nil {
		errMsg := "AI 代码生成服务未初始化"
		_ = h.taskSvc.Fail(taskID, "AI_SERVICE_NOT_READY", errMsg)
		h.taskSvc.PublishFailed(taskID, errMsg)
		return
	}
	if err := h.ws.EnsureProjectDir(projectID); err != nil {
		_ = h.taskSvc.Fail(taskID, "WORKSPACE_ERROR", err.Error())
		h.taskSvc.PublishFailed(taskID, err.Error())
		return
	}

	h.taskSvc.PublishProgress(taskID, "ai_generation", 25, "正在调用真实 AI 生成项目代码")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	result, err := h.aiGenSvc.GenerateProjectFiles(ctx, projectID)
	if err != nil {
		_ = h.taskSvc.Fail(taskID, "AI_GENERATION_ERROR", err.Error())
		h.taskSvc.PublishFailed(taskID, err.Error())
		return
	}

	h.taskSvc.PublishProgress(taskID, "write_files", 70, "正在写入 AI 生成的项目文件")
	for _, file := range result.Files {
		if err := h.ws.WriteFile(projectID, file.Path, []byte(file.Content)); err != nil {
			_ = h.taskSvc.Fail(taskID, "WRITE_FILE_ERROR", err.Error())
			h.taskSvc.PublishFailed(taskID, err.Error())
			return
		}
	}

	h.taskSvc.PublishProgress(taskID, "finalize", 90, "正在整理资源、配置和依赖文件")
	if err := h.taskSvc.Complete(taskID, "AI code generation completed"); err != nil {
		h.taskSvc.PublishFailed(taskID, err.Error())
		return
	}
	h.taskSvc.PublishCompleted(taskID)
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
	t.POST("", h.Create)
	t.GET("/:id", h.Get)
	t.GET("/:id/stream", h.Stream)
	t.POST("/:id/retry", h.Retry)
	t.POST("/:id/cancel", h.Cancel)
}
