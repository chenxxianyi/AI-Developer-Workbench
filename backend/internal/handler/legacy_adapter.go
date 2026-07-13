package handler

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// LegacyAdapter maps old Workbench /api routes to new /api/v1 endpoints.
// It responds with 301 redirects and logs deprecation warnings.
type LegacyAdapter struct{}

func NewLegacyAdapter() *LegacyAdapter {
	return &LegacyAdapter{}
}

// legacyMappings maps old paths to new paths.
var legacyMappings = map[string]string{
	"/api/health":   "/api/v1/health",
	"/api/dashboard": "/api/v1/dashboard",
	"/api/projects":  "/api/v1/projects",
	"/api/reports":   "/api/v1/reports",
	"/api/tools/":    "/api/v1/tools/",
	"/api/jobs/":     "/api/v1/tasks/",
}

// Handle redirects legacy /api/* requests to /api/v1/* with 308 Permanent Redirect.
func (a *LegacyAdapter) Handle(c *gin.Context) {
	oldPath := c.Request.URL.Path

	for oldPrefix, newPrefix := range legacyMappings {
		if strings.HasPrefix(oldPath, oldPrefix) {
			newPath := strings.Replace(oldPath, oldPrefix, newPrefix, 1)
			c.Redirect(308, newPath)
			return
		}
	}

	// If no mapping found, just proxy to /api/v1 with path replaced
	newPath := "/api/v1" + strings.TrimPrefix(oldPath, "/api")
	c.Redirect(308, newPath)
}

// RegisterLegacyRoutes registers the legacy adapter on /api group.
func RegisterLegacyRoutes(r *gin.RouterGroup, adapter *LegacyAdapter) {
	r.Any("/*path", adapter.Handle)
}
