package handler

import (
	"net/http"
	"strconv"

	"ai-developer-workbench/internal/repository"
	"ai-developer-workbench/internal/util"

	"github.com/gin-gonic/gin"
)

type ObservabilityHandler struct {
	aiRunRepo repository.AIRunRepository
}

func NewObservabilityHandler(aiRunRepo repository.AIRunRepository) *ObservabilityHandler {
	return &ObservabilityHandler{aiRunRepo: aiRunRepo}
}

func RegisterObservabilityRoutes(rg *gin.RouterGroup, aiRunRepo repository.AIRunRepository) {
	h := NewObservabilityHandler(aiRunRepo)
	rg.GET("/ai/stats", h.GetStats)
}

func (h *ObservabilityHandler) GetStats(c *gin.Context) {
	toolType := c.Query("tool_type")
	days := 30
	if d := c.Query("days"); d != "" {
		if v, err := strconv.Atoi(d); err == nil && v >= 1 && v <= 90 {
			days = v
		}
	}

	stats, err := h.aiRunRepo.GetStats(c.Request.Context(), toolType, days)
	if err != nil {
		util.ErrorResponse(c, http.StatusInternalServerError, util.CodeInternalError, err.Error())
		return
	}

	util.SuccessResponse(c, stats)
}
