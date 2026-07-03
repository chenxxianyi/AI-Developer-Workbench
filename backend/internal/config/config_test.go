package config

import (
	"os"
	"testing"
)

func setTestEnv(key, value string) func() {
	old, existed := os.LookupEnv(key)
	os.Setenv(key, value)
	return func() {
		if existed {
			os.Setenv(key, old)
		} else {
			os.Unsetenv(key)
		}
	}
}

func TestLoadConfig_MockModeExplicitTrue(t *testing.T) {
	restore := setTestEnv("DATABASE_DRIVER", "sqlite")
	defer restore()

	os.Setenv("AI_MOCK_MODE", "true")
	os.Setenv("AI_API_KEY", "")
	defer func() {
		os.Unsetenv("AI_MOCK_MODE")
		os.Unsetenv("AI_API_KEY")
	}()

	cfg, err := LoadConfig("")
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}
	if !cfg.AI.MockMode {
		t.Error("expected MockMode=true when AI_MOCK_MODE=true")
	}
	if !cfg.IsMockMode() {
		t.Error("IsMockMode()=false, want true")
	}
}

func TestLoadConfig_AutoMockWhenNoAPIKey(t *testing.T) {
	restore := setTestEnv("DATABASE_DRIVER", "sqlite")
	defer restore()

	os.Setenv("AI_API_KEY", "")
	os.Unsetenv("AI_MOCK_MODE")
	defer os.Unsetenv("AI_API_KEY")

	cfg, err := LoadConfig("")
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}
	if !cfg.AI.MockMode {
		t.Error("expected MockMode=true when AI_API_KEY is empty (auto-mock)")
	}
}

func TestLoadConfig_RealModeWithAPIKey(t *testing.T) {
	restore := setTestEnv("DATABASE_DRIVER", "sqlite")
	defer restore()

	os.Setenv("AI_API_KEY", "sk-test-real-key")
	os.Setenv("AI_MOCK_MODE", "false")
	defer func() {
		os.Unsetenv("AI_API_KEY")
		os.Unsetenv("AI_MOCK_MODE")
	}()

	cfg, err := LoadConfig("")
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}
	if cfg.AI.MockMode {
		t.Error("expected MockMode=false when AI_MOCK_MODE=false and API key is set")
	}
	if cfg.IsMockMode() {
		t.Error("IsMockMode()=true, want false")
	}
}

func TestLoadConfig_MockModeOverridesAPIKey(t *testing.T) {
	restore := setTestEnv("DATABASE_DRIVER", "sqlite")
	defer restore()

	os.Setenv("AI_MOCK_MODE", "true")
	os.Setenv("AI_API_KEY", "sk-some-key")
	defer func() {
		os.Unsetenv("AI_MOCK_MODE")
		os.Unsetenv("AI_API_KEY")
	}()

	cfg, err := LoadConfig("")
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}
	if !cfg.AI.MockMode {
		t.Error("expected MockMode=true when AI_MOCK_MODE=true regardless of API key")
	}
}

func TestIsMockMode(t *testing.T) {
	tests := []struct {
		name     string
		cfg      Config
		expected bool
	}{
		{"mock mode", Config{AI: AIConfig{MockMode: true}}, true},
		{"real mode", Config{AI: AIConfig{MockMode: false}}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cfg.IsMockMode(); got != tt.expected {
				t.Errorf("IsMockMode()=%v, want %v", got, tt.expected)
			}
		})
	}
}
