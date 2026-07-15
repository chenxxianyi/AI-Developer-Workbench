package handler

import (
	"ai-developer-workbench/internal/service"
	"ai-developer-workbench/pkg/response"

	"github.com/gin-gonic/gin"
)

// FileHandler handles project file operations.
type FileHandler struct {
	ws *service.WorkspaceService
}

func NewFileHandler(ws *service.WorkspaceService) *FileHandler {
	return &FileHandler{ws: ws}
}

func (h *FileHandler) ListFiles(c *gin.Context) {
	projectID := c.Param("id")
	subPath := c.DefaultQuery("path", ".")

	entries, err := h.ws.ListDir(projectID, subPath)
	if err != nil {
		response.BusinessError(c, err.Error())
		return
	}

	type FileEntry struct {
		Name  string `json:"name"`
		IsDir bool   `json:"is_dir"`
	}
	result := make([]FileEntry, 0)
	for _, e := range entries {
		result = append(result, FileEntry{Name: e.Name(), IsDir: e.IsDir()})
	}
	response.Success(c, result)
}

func (h *FileHandler) GetFileContent(c *gin.Context) {
	projectID := c.Param("id")
	filePath := c.Query("path")
	if filePath == "" {
		response.ValidationError(c, "请指定文件路径")
		return
	}

	data, err := h.ws.ReadFile(projectID, filePath)
	if err != nil {
		response.BusinessError(c, err.Error())
		return
	}

	isText := service.IsTextFile(filePath)
	response.Success(c, gin.H{
		"path":      filePath,
		"content":   string(data),
		"size":      len(data),
		"is_binary": !isText,
	})
}

func (h *FileHandler) ExportProject(c *gin.Context) {
	projectID := c.Param("id")
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", "attachment; filename=project-"+projectID+".zip")

	exportSvc := service.NewProjectExportService(h.ws)
	if err := exportSvc.ExportProject(projectID, c.Writer); err != nil {
		c.AbortWithStatus(500)
	}
}

func RegisterFileRoutes(r *gin.RouterGroup, h *FileHandler) {
	f := r.Group("/projects/:id")
	f.GET("/files", h.ListFiles)
	f.GET("/files/content", h.GetFileContent)
	f.GET("/export", h.ExportProject)
}
