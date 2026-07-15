package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all application configuration (unified Workbench + Builder).
type Config struct {
	App       AppConfig
	Database  DatabaseConfig
	Redis     RedisConfig
	JWT       JWTConfig
	Security  SecurityConfig
	Upload    UploadConfig
	Workspace WorkspaceConfig
	AI        AIConfig
	CORS      CORSConfig
}

// ── App ──

type AppConfig struct {
	Env     string
	Port    string
	Host    string
	Version string
}

// ── Database ──

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

func (d DatabaseConfig) DSN() string {
	if d.Driver == "sqlite" {
		return d.Name
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=%s&multiStatements=true",
		d.User, d.Password, d.Host, d.Port, d.Name, d.Charset, d.Loc)
}

func (d DatabaseConfig) IsSQLite() bool { return d.Driver == "sqlite" }

// ── Redis ──

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

// ── JWT ──

type JWTConfig struct {
	Secret string
	Expire int // hours
}

// ── Security ──

type SecurityConfig struct {
	EncryptionKey string
}

// ── Upload ──

type UploadConfig struct {
	Dir                   string
	TempDir               string
	MaxUploadSizeMB       int
	MaxProjectFiles       int
	MaxFileReadBytes      int
	MaxProjectTotalReadKB int
	MaxZipUncompressedMB  int
}

// ── Workspace (for generation / build / preview) ──

type WorkspaceConfig struct {
	RootDir             string
	BuildTimeoutSec     int
	PreviewSessionTTL   int // minutes
	MaxConcurrentBuilds int
}

// ── AI ──

type AIConfig struct {
	Provider       string
	BaseURL        string
	APIKey         string
	Model          string
	VisionModel    string
	TimeoutSeconds int
	MaxRetries     int
}

// ── CORS ──

type CORSConfig struct {
	AllowOrigins []string
}

// ── Load ──

func LoadConfig(envFile string) (*Config, error) {
	if envFile != "" {
		if _, err := os.Stat(envFile); err == nil {
			if err := godotenv.Load(envFile); err != nil {
				return nil, fmt.Errorf("load .env: %w", err)
			}
		}
	}

	cfg := &Config{
		App: AppConfig{
			Env:     getEnv("APP_ENV", "development"),
			Port:    getEnv("APP_PORT", "8080"),
			Host:    getEnv("APP_HOST", "0.0.0.0"),
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
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "127.0.0.1:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "change-me-in-production"),
			Expire: getEnvInt("JWT_EXPIRE_HOURS", 72),
		},
		Security: SecurityConfig{
			EncryptionKey: getEnv("APP_ENCRYPTION_KEY", ""),
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
		Workspace: WorkspaceConfig{
			RootDir:             getEnv("WORKSPACE_DIR", "./workspace"),
			BuildTimeoutSec:     getEnvInt("BUILD_TIMEOUT_SEC", 600),
			PreviewSessionTTL:   getEnvInt("PREVIEW_SESSION_TTL_MINUTES", 120),
			MaxConcurrentBuilds: getEnvInt("MAX_CONCURRENT_BUILDS", 3),
		},
		AI: AIConfig{
			Provider:       getEnv("AI_PROVIDER", "openai"),
			BaseURL:        getEnv("AI_BASE_URL", "https://api.openai.com/v1"),
			APIKey:         getEnv("AI_API_KEY", ""),
			Model:          getEnv("AI_MODEL", "gpt-4.1"),
			VisionModel:    getEnv("AI_VISION_MODEL", "gpt-4.1"),
			TimeoutSeconds: getEnvInt("AI_TIMEOUT_SECONDS", 180),
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

func (c *Config) validate() error {
	if !c.Database.IsSQLite() {
		if c.Database.Host == "" {
			return fmt.Errorf("DATABASE_HOST is required")
		}
	}
	if c.AI.APIKey == "" {
		return fmt.Errorf("AI_API_KEY is required")
	}
	// Warn on weak JWT secret in production
	if c.App.Env == "production" && c.JWT.Secret == "change-me-in-production" {
		return fmt.Errorf("JWT_SECRET must be changed in production")
	}
	return nil
}

func (c *Config) IsDevelopment() bool { return c.App.Env == "development" }

// ── Helpers ──

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
