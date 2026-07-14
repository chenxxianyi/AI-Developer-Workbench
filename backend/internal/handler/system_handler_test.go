package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ai-developer-workbench/internal/config"
	"ai-developer-workbench/internal/dto"
	"ai-developer-workbench/internal/util"

	"github.com/gin-gonic/gin"
)

func TestGetStatus_ReturnsRealAIConfiguration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		AI:     config.AIConfig{Provider: "openai", Model: "gpt-4.1", VisionModel: "gpt-4.1"},
		Upload: config.UploadConfig{MaxUploadSizeMB: 20, MaxProjectFiles: 120, MaxZipUncompressedMB: 100},
	}

	router := gin.New()
	RegisterSystemRoutes(&router.RouterGroup, cfg)
	req := httptest.NewRequest(http.MethodGet, "/system/status", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	if json.Valid(w.Body.Bytes()) == false {
		t.Fatalf("invalid JSON response: %s", w.Body.String())
	}
	var resp util.Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	dataBytes, _ := json.Marshal(resp.Data)
	var status dto.SystemStatusDTO
	if err := json.Unmarshal(dataBytes, &status); err != nil {
		t.Fatal(err)
	}
	if status.Provider != "openai" || status.TextModel != "gpt-4.1" {
		t.Fatalf("unexpected AI status: %+v", status)
	}
}

func TestGetStatus_ReturnsUploadLimits(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		AI:     config.AIConfig{Provider: "test", Model: "test-model", VisionModel: "test-vision"},
		Upload: config.UploadConfig{MaxUploadSizeMB: 15, MaxProjectFiles: 80, MaxZipUncompressedMB: 50},
	}
	router := gin.New()
	RegisterSystemRoutes(&router.RouterGroup, cfg)
	req := httptest.NewRequest(http.MethodGet, "/system/status", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var resp util.Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	dataBytes, _ := json.Marshal(resp.Data)
	var status dto.SystemStatusDTO
	if err := json.Unmarshal(dataBytes, &status); err != nil {
		t.Fatal(err)
	}
	if status.UploadLimits.ImageMaxBytes != 15*1024*1024 {
		t.Fatalf("unexpected image limit: %d", status.UploadLimits.ImageMaxBytes)
	}
	if status.UploadLimits.ZipMaxFiles != 80 {
		t.Fatalf("unexpected file limit: %d", status.UploadLimits.ZipMaxFiles)
	}
}
