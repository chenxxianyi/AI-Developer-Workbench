package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORS middleware handles cross-origin requests with an allowlist.
func CORS(allowOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// Check if origin is in allowlist.
		allowed := false
		for _, o := range allowOrigins {
			if o == origin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers", "Content-Type, X-Request-ID")
			c.Header("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
			c.Header("Access-Control-Expose-Headers", "Content-Disposition, X-Request-ID")
		}

		// Handle preflight.
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}