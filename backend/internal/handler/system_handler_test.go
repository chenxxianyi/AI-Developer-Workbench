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

func TestGetStatus_ReturnsMockMode(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		mockMode     bool
		expectMocked bool
	}{
		{"mock mode enabled", true, true},
		{"real mode (mock disabled)", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				AI: config.AIConfig{
					Provider:    "openai",
					Model:       "gpt-4.1",
					VisionModel: "gpt-4.1",
					MockMode:    tt.mockMode,
				},
				Upload: config.UploadConfig{
					MaxUploadSizeMB:      20,
					MaxProjectFiles:      120,
					MaxZipUncompressedMB: 100,
				},
				CORS: config.CORSConfig{
					AllowOrigins: []string{"*"},
				},
			}

			router := gin.New()
			RegisterSystemRoutes(&router.RouterGroup, cfg)

			req := httptest.NewRequest(http.MethodGet, "/system/status", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
			}

			var resp util.Response
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("failed to parse response: %v", err)
			}

			// Extract data from the response.
			dataBytes, _ := json.Marshal(resp.Data)
			var status dto.SystemStatusDTO
			if err := json.Unmarshal(dataBytes, &status); err != nil {
				t.Fatalf("failed to parse status data: %v", err)
			}

			if status.MockMode != tt.expectMocked {
				t.Errorf("mock_mode=%v, want %v", status.MockMode, tt.expectMocked)
			}
			if status.Provider != "openai" {
				t.Errorf("provider=%q, want %q", status.Provider, "openai")
			}
		})
	}
}

func TestGetStatus_ReturnsUploadLimits(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := &config.Config{
		AI: config.AIConfig{
			Provider:    "test",
			Model:       "test-model",
			VisionModel: "test-vision",
			MockMode:    true,
		},
		Upload: config.UploadConfig{
			MaxUploadSizeMB:      15,
			MaxProjectFiles:      80,
			MaxZipUncompressedMB: 50,
		},
		CORS: config.CORSConfig{
			AllowOrigins: []string{"*"},
		},
	}

	router := gin.New()
	RegisterSystemRoutes(&router.RouterGroup, cfg)

	req := httptest.NewRequest(http.MethodGet, "/system/status", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp util.Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	dataBytes, _ := json.Marshal(resp.Data)
	var status dto.SystemStatusDTO
	if err := json.Unmarshal(dataBytes, &status); err != nil {
		t.Fatalf("failed to parse status data: %v", err)
	}

	if status.UploadLimits.ImageMaxBytes != 15*1024*1024 {
		t.Errorf("ImageMaxBytes=%d, want %d", status.UploadLimits.ImageMaxBytes, 15*1024*1024)
	}
	if status.UploadLimits.ZipMaxFiles != 80 {
		t.Errorf("ZipMaxFiles=%d, want 80", status.UploadLimits.ZipMaxFiles)
	}
}
