package handler

import (
	"ai-developer-workbench/internal/model"
	"ai-developer-workbench/pkg/auth"
	"ai-developer-workbench/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuthHandler handles authentication endpoints.
type AuthHandler struct {
	db        *gorm.DB
	jwtSecret string
	jwtExpire int
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(db *gorm.DB, jwtSecret string, jwtExpire int) *AuthHandler {
	return &AuthHandler{db: db, jwtSecret: jwtSecret, jwtExpire: jwtExpire}
}

type registerReq struct {
	Username string `json:"username" binding:"required,min=2,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=100"`
}

type loginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type updateProfileReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "请填写完整的注册信息")
		return
	}

	var existing model.User
	if err := h.db.Where("username = ? OR email = ?", req.Username, req.Email).First(&existing).Error; err == nil {
		response.Conflict(c, "用户名或邮箱已被注册")
		return
	}

	user := model.User{
		ID:           uuid.New().String(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashPassword(req.Password),
		Role:         "user",
		Status:       "active",
	}
	if err := h.db.Create(&user).Error; err != nil {
		response.InternalError(c, "注册失败")
		return
	}

	token, _ := auth.GenerateToken(h.jwtSecret, user.ID, user.Username, user.Role, h.jwtExpire)
	response.Created(c, gin.H{"token": token, "user": user})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "请输入用户名和密码")
		return
	}

	var user model.User
	if err := h.db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		response.Unauthorized(c, "用户名或密码错误")
		return
	}
	if user.Status != "active" {
		response.Forbidden(c, "账户已被禁用")
		return
	}
	if !checkPassword(user.PasswordHash, req.Password) {
		response.Unauthorized(c, "用户名或密码错误")
		return
	}

	token, _ := auth.GenerateToken(h.jwtSecret, user.ID, user.Username, user.Role, h.jwtExpire)
	response.Success(c, gin.H{"token": token, "user": user})
}

func (h *AuthHandler) Profile(c *gin.Context) {
	userID := c.GetString("user_id")
	var user model.User
	if err := h.db.First(&user, "id = ?", userID).Error; err != nil {
		response.NotFound(c, "用户不存在")
		return
	}
	response.Success(c, user)
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	var req updateProfileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "无效的请求参数")
		return
	}

	userID := c.GetString("user_id")
	updates := map[string]interface{}{}
	if req.Username != "" {
		updates["username"] = req.Username
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if len(updates) == 0 {
		response.ValidationError(c, "没有要更新的字段")
		return
	}

	if err := h.db.Model(&model.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		response.InternalError(c, "更新失败")
		return
	}

	var user model.User
	h.db.First(&user, "id = ?", userID)
	response.Success(c, user)
}

// RegisterAuthRoutes registers auth routes.
func RegisterAuthRoutes(r *gin.RouterGroup, h *AuthHandler) {
	auth := r.Group("/auth")
	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)
	auth.GET("/profile", h.Profile)
	auth.PUT("/profile", h.UpdateProfile)
}

// Simple password hashing (replace with bcrypt in production).
func hashPassword(pwd string) string            { return pwd } // TODO: bcrypt
func checkPassword(hash, pwd string) bool        { return hash == pwd }
