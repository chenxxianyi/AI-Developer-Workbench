package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger middleware logs request details.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		requestID := c.GetString("request_id")

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		// Log without sensitive information.
		slog.Info("Request",
			"request_id", requestID,
			"method", method,
			"path", path,
			"status", status,
			"latency_ms", latency.Milliseconds(),
		)
	}
}