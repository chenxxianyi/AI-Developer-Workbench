package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// CORS middleware handles cross-origin requests with an allowlist.
//
// In production (APP_ENV=production) a wildcard ("*") or empty allowlist is
// rejected: no ACAO header is emitted for unmatched origins. This keeps the
// no-auth workbench from being exposed publicly with permissive CORS.
func CORS(allowOrigins []string, isProduction bool) gin.HandlerFunc {
	// Build a normalized allowlist. Drop wildcards in production.
	normalized := make([]string, 0, len(allowOrigins))
	for _, o := range allowOrigins {
		o = strings.TrimSpace(o)
		if o == "" {
			continue
		}
		if o == "*" {
			if isProduction {
				continue // never emit * in production
			}
			normalized = append(normalized, o)
			continue
		}
		normalized = append(normalized, o)
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// Check if origin is in allowlist.
		allowed := false
		for _, o := range normalized {
			if o == "*" || o == origin {
				allowed = true
				break
			}
		}

		if allowed && origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers", "Content-Type, X-Request-ID")
			c.Header("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
			c.Header("Access-Control-Expose-Headers", "Content-Disposition, X-Request-ID")
			c.Header("Vary", "Origin")
		}

		// Handle preflight: always return 204, but only echo ACAO for allowed origins.
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}