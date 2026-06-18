package handler

import (
	"context"

	"ai-developer-workbench/internal/database"
	"ai-developer-workbench/internal/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HealthHandler handles health check requests.
type HealthHandler struct {
	db *gorm.DB
}

// NewHealthHandler creates a new health handler.
func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// HealthStatus represents the health check response data.
type HealthStatus struct {
	Status    string `json:"status"`
	Database  string `json:"database,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}

// Check handles GET /api/health.
func (h *HealthHandler) Check(c *gin.Context) {
	ctx := context.Background()

	status := HealthStatus{
		Status: "ok",
	}

	if h.db != nil {
		if err := database.Ping(ctx, h.db); err != nil {
			status.Status = "degraded"
			status.Database = "error"
		} else {
			status.Database = "ok"
		}
	}

	util.SuccessResponse(c, status)
}

// RegisterHealthRoutes registers health check routes.
func RegisterHealthRoutes(r *gin.RouterGroup, db *gorm.DB) {
	handler := NewHealthHandler(db)
	r.GET("/health", handler.Check)
}