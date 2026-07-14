package middleware

import (
	"net"
	"strings"

	"ai-developer-workbench/pkg/auth"
	"ai-developer-workbench/pkg/response"

	"github.com/gin-gonic/gin"
)

const (
	// ContextKeyUserID stores the authenticated user's ID.
	ContextKeyUserID = "user_id"
	// ContextKeyUserRole stores the authenticated user's role.
	ContextKeyUserRole = "user_role"
	// ContextKeyProjectID stores the project ID from URL params (set by ProjectAccess).
	ContextKeyProjectID = "project_id"
)

// JWTAuth validates the JWT token and injects user info into the context.
func JWTAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "未提供认证令牌")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			response.Unauthorized(c, "认证格式错误")
			return
		}

		claims, err := auth.ParseToken(jwtSecret, parts[1])
		if err != nil {
			response.Unauthorized(c, "认证令牌无效或已过期")
			return
		}

		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyUserRole, claims.Role)
		c.Next()
	}
}

// RequireAuth is a simplified version that only checks if user is authenticated.
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, exists := c.Get(ContextKeyUserID); !exists {
			response.Unauthorized(c, "请先登录")
			return
		}
		c.Next()
	}
}

// RequireAdmin checks if the authenticated user has admin role.
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get(ContextKeyUserRole)
		if role != "admin" {
			response.Forbidden(c, "需要管理员权限")
			return
		}
		c.Next()
	}
}

// GetUserID extracts the user ID from the Gin context.
func GetUserID(c *gin.Context) string {
	id, _ := c.Get(ContextKeyUserID)
	if s, ok := id.(string); ok {
		return s
	}
	return ""
}

// GetUserRole extracts the user role from the Gin context.
func GetUserRole(c *gin.Context) string {
	role, _ := c.Get(ContextKeyUserRole)
	if s, ok := role.(string); ok {
		return s
	}
	return ""
}

// ProjectAccess validates that the current user owns the project identified by the :id or :projectId URL param.
// Admins bypass ownership check.
func ProjectAccess(projectRepo interface {
	GetUserIDByProjectID(projectID string) (string, error)
}) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID := c.Param("projectId")
		if projectID == "" {
			projectID = c.Param("id")
		}
		if projectID == "" {
			c.Next()
			return
		}

		userID := GetUserID(c)
		role := GetUserRole(c)

		// Admin can access any project
		if role == "admin" {
			c.Set(ContextKeyProjectID, projectID)
			c.Next()
			return
		}

		ownerID, err := projectRepo.GetUserIDByProjectID(projectID)
		if err != nil {
			response.NotFound(c, "项目不存在")
			return
		}
		if ownerID != userID {
			response.Forbidden(c, "无权访问该项目")
			return
		}

		c.Set(ContextKeyProjectID, projectID)
		c.Next()
	}
}

// GetProjectID extracts the project ID from the Gin context.
func GetProjectID(c *gin.Context) string {
	id, _ := c.Get(ContextKeyProjectID)
	if s, ok := id.(string); ok {
		return s
	}
	return ""
}

// GetClientIP safely extracts the client IP, respecting X-Forwarded-For and X-Real-Ip headers.
// Falls back to c.ClientIP().
func GetClientIP(c *gin.Context) string {
	if ip := c.GetHeader("X-Real-Ip"); ip != "" {
		return ip
	}
	if fwd := c.GetHeader("X-Forwarded-For"); fwd != "" {
		parts := strings.Split(fwd, ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}
	if ip, _, err := net.SplitHostPort(c.Request.RemoteAddr); err == nil && ip != "" {
		return ip
	}
	return c.ClientIP()
}
