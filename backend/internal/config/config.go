package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all application configuration.
type Config struct {
	App     AppConfig
	Database DatabaseConfig
	Upload  UploadConfig
	AI      AIConfig
	CORS    CORSConfig
}

// AppConfig holds application-level settings.
type AppConfig struct {
	Env     string
	Port    string
	Version string
}

// DatabaseConfig holds MySQL connection settings.
type DatabaseConfig struct {
	Driver          string
	Host            string
	Port            string
	Name            string
	User            string
	Password        string
	Charset         string
	Loc             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int // minutes
	AutoMigrate     bool
}

// DSN returns the data source name for the configured driver.
// This value must never be logged.
func (d DatabaseConfig) DSN() string {
	if d.Driver == "sqlite" {
		return d.Name // For SQLite, Name is the file path
	}
	// multiStatements is required because the migration runner executes each
	// versioned SQL file as a single Exec call, and the migration files contain
	// multiple statements separated by semicolons.
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=%s&multiStatements=true",
		d.User, d.Password, d.Host, d.Port, d.Name, d.Charset, d.Loc)
}

// IsSQLite returns true if the database driver is SQLite.
func (d DatabaseConfig) IsSQLite() bool {
	return d.Driver == "sqlite"
}

// UploadConfig holds file upload and ZIP processing limits.
type UploadConfig struct {
	Dir                   string
	TempDir               string
	MaxUploadSizeMB       int
	MaxProjectFiles       int
	MaxFileReadBytes      int
	MaxProjectTotalReadKB int
	MaxZipUncompressedMB  int
}

// AIConfig holds AI provider settings.
type AIConfig struct {
	Provider       string
	BaseURL        string
	APIKey         string
	Model          string
	VisionModel    string
	MockMode       bool
	TimeoutSeconds int
	MaxRetries     int
}

// CORSConfig holds CORS allowlist settings.
type CORSConfig struct {
	AllowOrigins []string
}

// LoadConfig loads configuration from .env file and environment variables.
// If envFile is empty, no .env file is loaded.
func LoadConfig(envFile string) (*Config, error) {
	if envFile != "" {
		if _, err := os.Stat(envFile); err == nil {
			if err := godotenv.Load(envFile); err != nil {
				return nil, fmt.Errorf("failed to load .env file: %w", err)
			}
		}
	}

	cfg := &Config{
		App: AppConfig{
			Env:     getEnv("APP_ENV", "development"),
			Port:    getEnv("APP_PORT", "8080"),
			Version: getEnv("APP_VERSION", "0.1.0"),
		},
		Database: DatabaseConfig{
			Driver:          getEnv("DATABASE_DRIVER", "mysql"),
			Host:            getEnv("DATABASE_HOST", "127.0.0.1"),
			Port:            getEnv("DATABASE_PORT", "3306"),
			Name:            getEnv("DATABASE_NAME", "ai_workbench"),
			User:            getEnv("DATABASE_USER", "workbench"),
			Password:        getEnv("DATABASE_PASSWORD", ""),
			Charset:         getEnv("DATABASE_CHARSET", "utf8mb4"),
			Loc:             getEnv("DATABASE_LOC", "UTC"),
			MaxOpenConns:    getEnvInt("DATABASE_MAX_OPEN_CONNS", 20),
			MaxIdleConns:    getEnvInt("DATABASE_MAX_IDLE_CONNS", 10),
			ConnMaxLifetime: getEnvInt("DATABASE_CONN_MAX_LIFETIME_MINUTES", 30),
			AutoMigrate:     getEnvBool("DB_AUTO_MIGRATE", true),
		},
		Upload: UploadConfig{
			Dir:                   getEnv("UPLOAD_DIR", "./uploads"),
			TempDir:               getEnv("TEMP_DIR", "./temp"),
			MaxUploadSizeMB:       getEnvInt("MAX_UPLOAD_SIZE_MB", 20),
			MaxProjectFiles:       getEnvInt("MAX_PROJECT_FILES", 120),
			MaxFileReadBytes:      getEnvInt("MAX_FILE_READ_BYTES", 12000),
			MaxProjectTotalReadKB: getEnvInt("MAX_PROJECT_TOTAL_READ_BYTES", 300000),
			MaxZipUncompressedMB:  getEnvInt("MAX_ZIP_UNCOMPRESSED_MB", 100),
		},
		AI: AIConfig{
			Provider:       getEnv("AI_PROVIDER", "openai"),
			BaseURL:        getEnv("AI_BASE_URL", "https://api.openai.com/v1"),
			APIKey:         getEnv("AI_API_KEY", ""),
			Model:          getEnv("AI_MODEL", "gpt-4.1"),
			VisionModel:    getEnv("AI_VISION_MODEL", "gpt-4.1"),
			MockMode:       getEnvBool("AI_MOCK_MODE", false),
			TimeoutSeconds: getEnvInt("AI_TIMEOUT_SECONDS", 90),
			MaxRetries:     getEnvInt("AI_MAX_RETRIES", 1),
		},
		CORS: CORSConfig{
			AllowOrigins: getEnvSlice("CORS_ALLOW_ORIGINS", []string{"http://localhost:5173"}),
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// validate checks required configuration values.
func (c *Config) validate() error {
	// SQLite doesn't require host/password
	if !c.Database.IsSQLite() {
		if c.Database.Host == "" {
			return fmt.Errorf("DATABASE_HOST is required")
		}
		if c.Database.Password == "" {
			return fmt.Errorf("DATABASE_PASSWORD is required")
		}
	}

	// Normalize mock mode: explicit MockMode=true forces mock; empty API key auto-mocks.
	if c.AI.MockMode || c.AI.APIKey == "" {
		c.AI.MockMode = true
		return nil
	}

	// Only validate API key when NOT in mock mode.
	if c.AI.APIKey == "" {
		return fmt.Errorf("AI_API_KEY is required")
	}
	return nil
}

// IsMockMode returns true when the AI is running in mock mode.
func (c *Config) IsMockMode() bool {
	return c.AI.MockMode
}

// IsDevelopment returns true if running in development mode.
func (c *Config) IsDevelopment() bool {
	return c.App.Env == "development"
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return n
}

func getEnvBool(key string, defaultVal bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	b, err := strconv.ParseBool(val)
	if err != nil {
		return defaultVal
	}
	return b
}

func getEnvSlice(key string, defaultVal []string) []string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	parts := strings.Split(val, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	if len(result) == 0 {
		return defaultVal
	}
	return result
}
