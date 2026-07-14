package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ai-developer-workbench/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestPasswordHashing(t *testing.T) {
	hash, err := hashPassword("AdminTest@2026!")
	require.NoError(t, err)
	require.NotEqual(t, "AdminTest@2026!", hash)
	require.True(t, checkPassword(hash, "AdminTest@2026!"))
	require.False(t, checkPassword(hash, "wrong-password"))
}

func TestAdminLoginAndProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.User{}))

	hash, err := hashPassword("AdminTest@2026!")
	require.NoError(t, err)
	user := model.User{
		ID: uuid.New().String(), Username: "admin", Email: "admin@example.test",
		PasswordHash: hash, Role: "admin", Status: "active",
	}
	require.NoError(t, db.Create(&user).Error)

	router := gin.New()
	api := router.Group("/api")
	RegisterAuthRoutes(api, NewAuthHandler(db, "test-secret", 1), "test-secret")

	body, _ := json.Marshal(map[string]string{"username": "admin", "password": "AdminTest@2026!"})
	loginReq := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(body))
	loginReq.Header.Set("Content-Type", "application/json")
	loginResp := httptest.NewRecorder()
	router.ServeHTTP(loginResp, loginReq)
	require.Equal(t, http.StatusOK, loginResp.Code, loginResp.Body.String())

	var payload struct {
		Data struct {
			Token string     `json:"token"`
			User  model.User `json:"user"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(loginResp.Body.Bytes(), &payload))
	require.NotEmpty(t, payload.Data.Token)
	require.Equal(t, "admin", payload.Data.User.Role)

	profileReq := httptest.NewRequest(http.MethodGet, "/api/auth/profile", nil)
	profileReq.Header.Set("Authorization", "Bearer "+payload.Data.Token)
	profileResp := httptest.NewRecorder()
	router.ServeHTTP(profileResp, profileReq)
	require.Equal(t, http.StatusOK, profileResp.Code, profileResp.Body.String())
}
