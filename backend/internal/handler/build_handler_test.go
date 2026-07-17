package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"ai-developer-workbench/internal/service"
	"ai-developer-workbench/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type fakeProjectBuilder struct {
	build func(context.Context, string) (string, error)
}

func (f fakeProjectBuilder) Build(ctx context.Context, projectID string) (string, error) {
	return f.build(ctx, projectID)
}

func TestBuildHandlerBuildReturnsServedPreviewURL(t *testing.T) {
	gin.SetMode(gin.TestMode)
	projectID := uuid.NewString()
	ws := service.NewWorkspaceService(t.TempDir())
	builder := fakeProjectBuilder{build: func(_ context.Context, gotProjectID string) (string, error) {
		require.Equal(t, projectID, gotProjectID)
		require.NoError(t, ws.WriteFile(projectID, "dist/index.html", []byte("<!doctype html><h1>built</h1>")))
		return "built", nil
	}}
	router := buildTestRouter(NewBuildHandler(ws, builder))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/projects/"+projectID+"/build", nil)
	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusOK, recorder.Code)
	var body response.APIResponse
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &body))
	data, ok := body.Data.(map[string]interface{})
	require.True(t, ok)
	require.Equal(t, "/api/projects/"+projectID+"/preview/", data["preview_url"])
}

func TestBuildHandlerBuildReportsBuildFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	projectID := uuid.NewString()
	ws := service.NewWorkspaceService(t.TempDir())
	builder := fakeProjectBuilder{build: func(context.Context, string) (string, error) {
		return "", errors.New("vite failed")
	}}
	router := buildTestRouter(NewBuildHandler(ws, builder))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/projects/"+projectID+"/build", nil)
	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
	require.Contains(t, recorder.Body.String(), "vite failed")
}

func TestBuildHandlerPreviewServesAssetsAndSPAFallback(t *testing.T) {
	gin.SetMode(gin.TestMode)
	projectID := uuid.NewString()
	ws := service.NewWorkspaceService(t.TempDir())
	require.NoError(t, ws.WriteFile(projectID, "dist/index.html", []byte("<!doctype html><h1>preview</h1>")))
	require.NoError(t, ws.WriteFile(projectID, "dist/assets/app.js", []byte("console.log('preview')")))
	handler := NewBuildHandler(ws, fakeProjectBuilder{})
	router := buildTestRouter(handler)

	tests := []struct {
		path        string
		wantStatus  int
		wantBody    string
		wantContent string
	}{
		{path: "/api/projects/" + projectID + "/preview/", wantStatus: http.StatusOK, wantBody: "<h1>preview</h1>", wantContent: "text/html"},
		{path: "/api/projects/" + projectID + "/preview/assets/app.js", wantStatus: http.StatusOK, wantBody: "console.log", wantContent: "text/javascript"},
		{path: "/api/projects/" + projectID + "/preview/game/room", wantStatus: http.StatusOK, wantBody: "<h1>preview</h1>", wantContent: "text/html"},
		{path: "/api/projects/" + projectID + "/preview/assets/missing.js", wantStatus: http.StatusNotFound, wantBody: "预览文件不存在", wantContent: "application/json"},
	}

	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, test.path, nil)
			router.ServeHTTP(recorder, request)

			require.Equal(t, test.wantStatus, recorder.Code)
			require.Contains(t, recorder.Body.String(), test.wantBody)
			require.Contains(t, recorder.Header().Get("Content-Type"), test.wantContent)
			if test.wantStatus == http.StatusOK {
				require.Contains(t, recorder.Header().Get("Content-Security-Policy"), "sandbox")
			}
		})
	}
}

func TestBuildHandlerRejectsInvalidProjectID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ws := service.NewWorkspaceService(t.TempDir())
	router := buildTestRouter(NewBuildHandler(ws, fakeProjectBuilder{}))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/projects/not-a-uuid/preview/", nil)
	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestBuildHandlerStatusReportsExistingPreview(t *testing.T) {
	gin.SetMode(gin.TestMode)
	projectID := uuid.NewString()
	ws := service.NewWorkspaceService(t.TempDir())
	require.NoError(t, ws.WriteFile(projectID, "dist/index.html", []byte("<!doctype html>")))
	router := buildTestRouter(NewBuildHandler(ws, fakeProjectBuilder{}))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/projects/"+projectID+"/build", nil)
	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Contains(t, recorder.Body.String(), `"ready":true`)
	require.Contains(t, recorder.Body.String(), "/preview/")
}

func buildTestRouter(handler *BuildHandler) *gin.Engine {
	router := gin.New()
	api := router.Group("/api")
	RegisterBuildRoutes(api, handler)
	return router
}

func TestIsWithinDir(t *testing.T) {
	root := filepath.Join(string(os.PathSeparator), "workspace", "project", "dist")
	require.True(t, isWithinDir(root, filepath.Join(root, "assets", "app.js")))
	require.False(t, isWithinDir(root, filepath.Join(root, "..", "secret.txt")))
}
