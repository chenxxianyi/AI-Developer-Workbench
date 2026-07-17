package handler

import (
	"context"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"ai-developer-workbench/internal/service"
	"ai-developer-workbench/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type projectBuilder interface {
	Build(ctx context.Context, projectID string) (string, error)
}

// BuildHandler handles project builds and serves their generated preview files.
type BuildHandler struct {
	ws      *service.WorkspaceService
	builder projectBuilder
}

func NewBuildHandler(ws *service.WorkspaceService, builder projectBuilder) *BuildHandler {
	return &BuildHandler{ws: ws, builder: builder}
}

func (h *BuildHandler) Build(c *gin.Context) {
	projectID := c.Param("id")
	if !validProjectID(projectID) {
		response.ValidationError(c, "无效的项目 ID")
		return
	}

	if _, err := h.builder.Build(c.Request.Context(), projectID); err != nil {
		response.BusinessError(c, "项目构建失败: "+err.Error())
		return
	}

	// A successful command without an entry document is not a usable preview.
	if _, err := h.ws.ReadFile(projectID, "dist/index.html"); err != nil {
		response.InternalError(c, "构建完成，但未生成 dist/index.html")
		return
	}

	response.Success(c, gin.H{
		"preview_url": "/api/projects/" + url.PathEscape(projectID) + "/preview/",
	})
}

func (h *BuildHandler) Status(c *gin.Context) {
	projectID := c.Param("id")
	if !validProjectID(projectID) {
		response.ValidationError(c, "无效的项目 ID")
		return
	}
	previewURL := "/api/projects/" + url.PathEscape(projectID) + "/preview/"
	if _, err := h.ws.ReadFile(projectID, "dist/index.html"); err != nil {
		response.Success(c, gin.H{"ready": false, "preview_url": ""})
		return
	}
	response.Success(c, gin.H{"ready": true, "preview_url": previewURL})
}

func (h *BuildHandler) Preview(c *gin.Context) {
	projectID := c.Param("id")
	if !validProjectID(projectID) {
		response.ValidationError(c, "无效的项目 ID")
		return
	}

	target, isIndex, err := h.resolvePreviewFile(projectID, c.Param("filepath"))
	if err != nil {
		response.NotFound(c, "预览文件不存在")
		return
	}

	// Generated applications are untrusted. The response-level sandbox also
	// protects previews opened in a new tab, outside the workbench iframe.
	c.Header("Content-Security-Policy", "sandbox allow-scripts allow-forms allow-modals allow-popups allow-downloads")
	c.Header("Cross-Origin-Resource-Policy", "cross-origin")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("X-Content-Type-Options", "nosniff")
	if isIndex {
		c.Header("Cache-Control", "no-store")
	} else {
		c.Header("Cache-Control", "public, max-age=3600")
	}
	c.File(target)
}

func (h *BuildHandler) resolvePreviewFile(projectID, requestPath string) (string, bool, error) {
	// Backslashes are path separators on Windows and must not bypass URL path
	// normalization. path.Clean handles URL paths; filepath.Rel enforces the
	// final filesystem boundary.
	if strings.Contains(requestPath, "\\") {
		return "", false, os.ErrNotExist
	}

	cleanPath := strings.TrimPrefix(path.Clean("/"+requestPath), "/")
	if cleanPath == "." || cleanPath == "" {
		cleanPath = "index.html"
	}

	distDir, err := h.ws.SafeResolve(projectID, "dist")
	if err != nil {
		return "", false, err
	}
	target, err := h.ws.SafeResolve(projectID, filepath.Join("dist", filepath.FromSlash(cleanPath)))
	if err != nil {
		return "", false, err
	}
	if !isWithinDir(distDir, target) {
		return "", false, os.ErrNotExist
	}

	info, statErr := os.Stat(target)
	if statErr == nil && !info.IsDir() {
		return target, cleanPath == "index.html", nil
	}

	// History-mode SPA routes have no extension. Missing assets retain a 404
	// instead of receiving HTML with the wrong content type.
	if filepath.Ext(cleanPath) != "" {
		return "", false, os.ErrNotExist
	}
	indexPath := filepath.Join(distDir, "index.html")
	if info, err := os.Stat(indexPath); err != nil || info.IsDir() {
		return "", false, os.ErrNotExist
	}
	return indexPath, true, nil
}

func isWithinDir(root, target string) bool {
	rel, err := filepath.Rel(root, target)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

func validProjectID(projectID string) bool {
	_, err := uuid.Parse(projectID)
	return err == nil
}

func RegisterBuildRoutes(r *gin.RouterGroup, h *BuildHandler) {
	r.POST("/projects/:id/build", h.Build)
	r.GET("/projects/:id/build", h.Status)
	r.GET("/projects/:id/preview/*filepath", h.Preview)
}
