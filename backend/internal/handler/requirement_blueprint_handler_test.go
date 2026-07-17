package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"ai-developer-workbench/internal/model"
	"ai-developer-workbench/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type failingBlueprintGenerator struct{ err error }

func (f failingBlueprintGenerator) GenerateBlueprint(context.Context, string) (*service.BlueprintAIResult, error) {
	return nil, f.err
}

func newRequirementBlueprintTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.Project{}, &model.Requirement{}, &model.Blueprint{}))
	return db
}

func TestBlueprintConfirmRejectsSupersededVersionButGetStillReturnsIt(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := newRequirementBlueprintTestDB(t)
	projectID := uuid.NewString()
	require.NoError(t, db.Create(&model.Project{ID: projectID, Name: "测试", ProjectType: "utility_app"}).Error)
	require.NoError(t, db.Create(&model.Requirement{ID: uuid.NewString(), ProjectID: projectID, Version: 2, Content: `{"schema_version":2,"app_type":"utility_app","goal":"测试","target_users":["用户"],"must_have_features":["计算"],"acceptance_criteria":["显示结果"]}`}).Error)
	require.NoError(t, db.Create(&model.Blueprint{ID: uuid.NewString(), ProjectID: projectID, Version: 1, Status: "superseded", Content: `{"product_positioning":"旧蓝图","pages":[{"name":"首页","route":"/"}]}`}).Error)

	router := gin.New()
	api := router.Group("/api")
	RegisterBlueprintRoutes(api, NewBlueprintHandler(db, nil))

	getRecorder := httptest.NewRecorder()
	router.ServeHTTP(getRecorder, httptest.NewRequest(http.MethodGet, "/api/projects/"+projectID+"/blueprint", nil))
	require.Equal(t, http.StatusOK, getRecorder.Code)
	require.Contains(t, getRecorder.Body.String(), "superseded")

	confirmRecorder := httptest.NewRecorder()
	router.ServeHTTP(confirmRecorder, httptest.NewRequest(http.MethodPost, "/api/projects/"+projectID+"/blueprint/confirm", nil))
	require.Equal(t, http.StatusBadRequest, confirmRecorder.Code)
	require.Contains(t, confirmRecorder.Body.String(), "当前蓝图已失效")
}

func TestRequirementSaveWithEquivalentJSONKeepsConfirmedBlueprint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := newRequirementBlueprintTestDB(t)
	projectID := uuid.NewString()
	require.NoError(t, db.Create(&model.Project{ID: projectID, Name: "测试", ProjectType: "utility_app"}).Error)
	existingContent := `{"schema_version":2,"app_type":"utility_app","goal":"计算工具","target_users":["用户"],"must_have_features":["计算"],"acceptance_criteria":["显示结果"]}`
	require.NoError(t, db.Create(&model.Requirement{ID: uuid.NewString(), ProjectID: projectID, Version: 1, Content: existingContent}).Error)
	blueprintID := uuid.NewString()
	require.NoError(t, db.Create(&model.Blueprint{ID: blueprintID, ProjectID: projectID, Version: 1, Status: "confirmed", Content: `{}`}).Error)

	equivalentContent := `{"goal":"计算工具","schema_version":2,"target_users":["用户"],"app_type":"utility_app","acceptance_criteria":["显示结果"],"must_have_features":["计算"]}`
	body, err := json.Marshal(map[string]string{"content": equivalentContent})
	require.NoError(t, err)
	router := gin.New()
	api := router.Group("/api")
	RegisterRequirementRoutes(api, NewRequirementHandler(db))
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPut, "/api/projects/"+projectID+"/requirements", bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusOK, recorder.Code)
	var count int64
	require.NoError(t, db.Model(&model.Requirement{}).Where("project_id = ?", projectID).Count(&count).Error)
	require.EqualValues(t, 1, count)
	var blueprint model.Blueprint
	require.NoError(t, db.First(&blueprint, "id = ?", blueprintID).Error)
	require.Equal(t, "confirmed", blueprint.Status)
}

func TestEquivalentJSONDetectsMeaningfulChanges(t *testing.T) {
	require.True(t, equivalentJSON(`{"a":1,"b":[2]}`, `{"b":[2],"a":1}`))
	require.False(t, equivalentJSON(`{"a":1}`, `{"a":2}`))
}

func TestBlueprintGenerateReturnsGatewayTimeoutForAIProviderTimeout(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api := router.Group("/api")
	RegisterBlueprintRoutes(api, NewBlueprintHandler(nil, failingBlueprintGenerator{err: errors.New("context deadline exceeded while waiting for model")}))
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/api/projects/project-timeout/blueprint/generate", nil))

	require.Equal(t, http.StatusGatewayTimeout, recorder.Code)
	require.Contains(t, recorder.Body.String(), "AI_TIMEOUT")
	require.Contains(t, recorder.Body.String(), "AI 模型响应超时")
}
