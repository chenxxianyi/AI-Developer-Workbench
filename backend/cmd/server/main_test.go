package main

import (
	"sort"
	"testing"

	"ai-developer-workbench/internal/config"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TestBuildRouterRegistersAllAPIRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := &config.Config{
		Upload: config.UploadConfig{
			Dir:     t.TempDir(),
			TempDir: t.TempDir(),
		},
		AI: config.AIConfig{
			Provider:    "openai",
			BaseURL:     "https://example.com/v1",
			APIKey:      "test-key",
			Model:       "test-model",
			VisionModel: "test-vision-model",
		},
		CORS: config.CORSConfig{
			AllowOrigins: []string{"http://localhost:5173"},
		},
	}

	router := buildRouter(cfg, &gorm.DB{})

	got := make([]string, 0, len(router.Routes()))
	for _, route := range router.Routes() {
		got = append(got, route.Method+" "+route.Path)
	}
	sort.Strings(got)

	want := []string{
		"DELETE /api/projects/:id",
		"DELETE /api/reports/:id",
		"GET /api/admin/models",
		"GET /api/admin/projects",
		"GET /api/admin/prompts",
		"GET /api/admin/users",
		"GET /api/admin/users/:userId",
		"GET /api/ai/stats",
		"GET /api/auth/profile",
		"GET /api/dashboard/stats",
		"GET /api/health",
		"GET /api/jobs/:id",
		"GET /api/projects",
		"GET /api/projects/:id",
		"GET /api/projects/:id/reports",
		"GET /api/projects/:id/stats",
		"GET /api/projects/:id/blueprint",
		"GET /api/projects/:id/requirements",
		"GET /api/reports",
		"GET /api/reports/:id",
		"GET /api/reports/:id/compare/:targetId",
		"GET /api/reports/:id/export",
		"GET /api/reports/:id/files/:filename",
		"GET /api/system/status",
		"GET /api/tools",
		"PATCH /api/projects/:id",
		"POST /api/auth/login",
		"POST /api/auth/register",
		"POST /api/jobs/:id/cancel",
		"POST /api/jobs/:id/retry",
		"POST /api/projects",
		"POST /api/projects/:id/blueprint/confirm",
		"POST /api/projects/:id/blueprint/generate",
		"POST /api/tools/agent-config/run",
		"POST /api/tools/api-doc/run",
		"POST /api/tools/db-schema/run",
		"POST /api/tools/project-doctor/run",
		"POST /api/tools/ui-review/run",
		"PUT /api/admin/models/:id",
		"PUT /api/admin/prompts/:id",
		"PUT /api/admin/users/:userId/status",
		"PUT /api/auth/profile",
		"PUT /api/projects/:id/blueprint",
		"PUT /api/projects/:id/requirements",
	}
	sort.Strings(want)

	if len(got) != len(want) {
		t.Fatalf("registered route count = %d, want %d\nroutes: %v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("registered routes = %v, want %v", got, want)
		}
	}
}
