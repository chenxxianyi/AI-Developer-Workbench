package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"ai-developer-workbench/internal/util"

	"github.com/gin-gonic/gin"
)

// Recovery middleware catches panics and returns a unified error response.
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the stack trace internally.
				stack := string(debug.Stack())
				requestID := c.GetString("request_id")
				slog.Error("Panic recovered",
					"error", err,
					"request_id", requestID,
					"stack", stack,
				)

				// Return unified error without exposing stack trace.
				util.ErrorResponse(c, http.StatusInternalServerError, util.CodeInternalError, "An unexpected error occurred")
				c.Abort()
			}
		}()
		c.Next()
	}
}