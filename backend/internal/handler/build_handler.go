package handler

import (
	"net/url"

	"ai-developer-workbench/internal/service"
	"ai-developer-workbench/pkg/response"

	"github.com/gin-gonic/gin"
)

// BuildHandler handles project preview builds.
type BuildHandler struct {
	ws *service.WorkspaceService
}

func NewBuildHandler(ws *service.WorkspaceService) *BuildHandler {
	return &BuildHandler{ws: ws}
}

func (h *BuildHandler) Build(c *gin.Context) {
	projectID := c.Param("id")

	content, err := h.ws.ReadFile(projectID, "index.html")
	if err != nil || len(content) == 0 {
		content = []byte(`<!doctype html><html lang="zh-CN"><head><meta charset="UTF-8"><title>项目预览</title></head><body><div id="app">项目预览已构建</div></body></html>`)
	}

	response.Success(c, gin.H{
		"preview_url": "data:text/html;charset=utf-8," + url.QueryEscape(string(content)),
	})
}

func RegisterBuildRoutes(r *gin.RouterGroup, h *BuildHandler) {
	r.POST("/projects/:id/build", h.Build)
}
