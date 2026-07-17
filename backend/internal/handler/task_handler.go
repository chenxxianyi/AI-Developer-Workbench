package handler

import (
	"context"
	"errors"
	"fmt"
	"sync"
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
	builder  projectBuilder
	cancelMu sync.Mutex
	cancels  map[string]context.CancelFunc
}

func NewTaskHandler(taskSvc *service.TaskService, broker *sse.Broker, ws *service.WorkspaceService, aiGenSvc *service.AIGenerationService, builder projectBuilder) *TaskHandler {
	return &TaskHandler{taskSvc: taskSvc, broker: broker, ws: ws, aiGenSvc: aiGenSvc, builder: builder, cancels: make(map[string]context.CancelFunc)}
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
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		h.cancelMu.Lock()
		h.cancels[task.ID] = cancel
		h.cancelMu.Unlock()
		go h.runGenerationTask(ctx, task.ID, req.ProjectID)
	}

	response.Created(c, task)
}

func (h *TaskHandler) runGenerationTask(ctx context.Context, taskID, projectID string) {
	defer func() {
		h.cancelMu.Lock()
		if cancel := h.cancels[taskID]; cancel != nil {
			cancel()
			delete(h.cancels, taskID)
		}
		h.cancelMu.Unlock()
	}()
	select {
	case <-time.After(150 * time.Millisecond):
	case <-ctx.Done():
		return
	}
	if err := h.taskSvc.Start(taskID); err != nil {
		if errors.Is(ctx.Err(), context.Canceled) {
			return
		}
		h.taskSvc.PublishFailed(taskID, err.Error())
		return
	}

	h.taskSvc.PublishProgress(taskID, "prepare", 5, "正在加载结构化需求和已确认蓝图")
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

	h.taskSvc.PublishProgress(taskID, "ai_generation", 20, "正在按功能和模块生成真实应用代码")
	result, err := h.aiGenSvc.GenerateProjectFiles(ctx, projectID)
	if err != nil {
		if errors.Is(ctx.Err(), context.Canceled) {
			return
		}
		_ = h.taskSvc.Fail(taskID, "AI_GENERATION_ERROR", err.Error())
		h.taskSvc.PublishFailed(taskID, err.Error())
		return
	}

	h.taskSvc.PublishProgress(taskID, "write_files", 55, fmt.Sprintf("正在写入 %d 个生成文件", len(result.Files)))
	for _, file := range result.Files {
		if ctx.Err() != nil {
			return
		}
		if err := h.ws.WriteFile(projectID, file.Path, []byte(file.Content)); err != nil {
			_ = h.taskSvc.Fail(taskID, "WRITE_FILE_ERROR", err.Error())
			h.taskSvc.PublishFailed(taskID, err.Error())
			return
		}
	}

	if h.builder == nil {
		errMsg := "构建验证服务未初始化"
		_ = h.taskSvc.Fail(taskID, "BUILD_SERVICE_NOT_READY", errMsg)
		h.taskSvc.PublishFailed(taskID, errMsg)
		return
	}
	h.taskSvc.PublishProgress(taskID, "build", 70, "正在运行测试并执行生产构建")
	var buildOutput string
	for attempt := 0; attempt < 3; attempt++ {
		buildOutput, err = h.builder.Build(ctx, projectID)
		if err == nil {
			break
		}
		if errors.Is(ctx.Err(), context.Canceled) {
			return
		}
		if attempt == 2 {
			_ = h.taskSvc.Fail(taskID, "BUILD_VALIDATION_ERROR", err.Error())
			h.taskSvc.PublishFailed(taskID, "生成代码经过两轮自动修复后仍未通过验证: "+err.Error())
			return
		}
		h.taskSvc.PublishProgress(taskID, "repair", 76+attempt*7, fmt.Sprintf("构建验证失败，正在进行第 %d 轮定向修复", attempt+1))
		repaired, repairErr := h.aiGenSvc.RepairProjectFiles(ctx, projectID, result, err.Error())
		if repairErr != nil {
			_ = h.taskSvc.Fail(taskID, "AUTO_REPAIR_ERROR", repairErr.Error())
			h.taskSvc.PublishFailed(taskID, "自动修复失败: "+repairErr.Error())
			return
		}
		for _, file := range repaired.Files {
			if err := h.ws.WriteFile(projectID, file.Path, []byte(file.Content)); err != nil {
				_ = h.taskSvc.Fail(taskID, "REPAIR_WRITE_ERROR", err.Error())
				h.taskSvc.PublishFailed(taskID, err.Error())
				return
			}
		}
		result = repaired
	}

	h.taskSvc.PublishProgress(taskID, "verify", 92, "构建已通过，正在确认预览产物")
	if _, err := h.ws.ReadFile(projectID, "dist/index.html"); err != nil {
		_ = h.taskSvc.Fail(taskID, "PREVIEW_ARTIFACT_MISSING", err.Error())
		h.taskSvc.PublishFailed(taskID, "构建未生成可用的预览入口")
		return
	}
	resultSummary := fmt.Sprintf("Generated %d files and passed the production build", len(result.Files))
	if len(buildOutput) > 0 {
		resultSummary += "."
	}
	if err := h.taskSvc.Complete(taskID, resultSummary); err != nil {
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
	h.cancelMu.Lock()
	cancel := h.cancels[c.Param("id")]
	h.cancelMu.Unlock()
	if cancel != nil {
		cancel()
	}
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
