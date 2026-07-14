package config

import (
	"os"
	"strings"
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

func TestLoadConfig_RequiresAPIKey(t *testing.T) {
	restoreDriver := setTestEnv("DATABASE_DRIVER", "sqlite")
	defer restoreDriver()
	restoreKey := setTestEnv("AI_API_KEY", "")
	defer restoreKey()

	_, err := LoadConfig("")
	if err == nil || !strings.Contains(err.Error(), "AI_API_KEY is required") {
		t.Fatalf("LoadConfig() error = %v, want missing AI_API_KEY error", err)
	}
}

func TestLoadConfig_WithAPIKey(t *testing.T) {
	restoreDriver := setTestEnv("DATABASE_DRIVER", "sqlite")
	defer restoreDriver()
	restoreKey := setTestEnv("AI_API_KEY", "sk-test-real-key")
	defer restoreKey()

	cfg, err := LoadConfig("")
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}
	if cfg.AI.APIKey != "sk-test-real-key" {
		t.Fatalf("AI.APIKey = %q, want configured key", cfg.AI.APIKey)
	}
}

func TestDatabaseConfig_DSNIncludesMultiStatements(t *testing.T) {
	db := DatabaseConfig{
		Driver:   "mysql",
		Host:     "127.0.0.1",
		Port:     "3306",
		Name:     "ai_workbench",
		User:     "workbench",
		Password: "secret",
		Charset:  "utf8mb4",
		Loc:      "UTC",
	}

	got := db.DSN()
	if !strings.Contains(got, "multiStatements=true") {
		t.Fatalf("DSN() missing multiStatements=true: %s", got)
	}
	if !strings.Contains(got, "parseTime=True") {
		t.Fatalf("DSN() missing parseTime=True: %s", got)
	}
}
