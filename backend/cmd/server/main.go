package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ai-developer-workbench/internal/config"
	"ai-developer-workbench/internal/database"
	"ai-developer-workbench/internal/handler"
	"ai-developer-workbench/internal/middleware"
	"ai-developer-workbench/internal/repository"
	"ai-developer-workbench/internal/service"
	toolservice "ai-developer-workbench/internal/service/tools"
	"ai-developer-workbench/pkg/sse"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const shutdownTimeout = 15 * time.Second

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx); err != nil {
		slog.Error("Backend stopped", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}

	configureLogging(cfg)
	configureGin(cfg)

	if err := createStorageDirectories(cfg); err != nil {
		return err
	}

	db, err := database.Connect(&cfg.Database, cfg.IsDevelopment())
	if err != nil {
		return err
	}
	defer func() {
		if err := database.Close(db); err != nil {
			slog.Error("Failed to close database connection", "error", err)
		}
	}()

	if err := database.RunMigrations(db, cfg.Database.AutoMigrate); err != nil {
		return err
	}

	server := newHTTPServer(cfg, buildRouter(cfg, db))
	serverErrors := make(chan error, 1)
	go func() {
		slog.Info("Backend server started",
			"address", server.Addr,
			"environment", cfg.App.Env,
			"version", cfg.App.Version,
		)
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("serve HTTP: %w", err)
	case <-ctx.Done():
		slog.Info("Shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("graceful HTTP shutdown: %w", err)
	}

	slog.Info("Backend server stopped gracefully")
	return nil
}

func configureLogging(cfg *config.Config) {
	level := slog.LevelInfo
	if cfg.IsDevelopment() {
		level = slog.LevelDebug
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	slog.SetDefault(slog.New(handler))
}

func configureGin(cfg *config.Config) {
	if cfg.IsDevelopment() {
		gin.SetMode(gin.DebugMode)
		return
	}
	gin.SetMode(gin.ReleaseMode)
}

func createStorageDirectories(cfg *config.Config) error {
	for _, dir := range []string{cfg.Upload.Dir, cfg.Upload.TempDir} {
		if err := os.MkdirAll(dir, 0o750); err != nil {
			return fmt.Errorf("create storage directory %q: %w", dir, err)
		}
	}
	return nil
}

func buildRouter(cfg *config.Config, db *gorm.DB) *gin.Engine {
	reportRepo := repository.NewReportRepository(db)
	generatedFileRepo := repository.NewGeneratedFileRepository(db)
	reportAssetRepo := repository.NewReportAssetRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	jobRepo := repository.NewJobRepository(db)
	aiRunRepo := repository.NewAIRunRepository(db)

	reportService := service.NewReportService(
		cfg,
		reportRepo,
		generatedFileRepo,
		reportAssetRepo,
		projectRepo,
		db,
	)
	projectService := service.NewProjectService(projectRepo)
	jobService := service.NewJobService(db, jobRepo)
	fileService := service.NewFileService(cfg, reportAssetRepo)
	zipService := service.NewZipService(cfg.Upload.TempDir)

	var aiService service.AIService = service.NewOpenAICompatibleService(&cfg.AI)
	aiService = service.NewInstrumentedAIService(aiService, aiRunRepo)
	exportService := service.NewExportService(reportRepo, generatedFileRepo)

	agentConfigService := toolservice.NewAgentConfigService(aiService, reportService)
	dbSchemaService := toolservice.NewDBSchemaService(aiService, reportService)
	uiReviewService := toolservice.NewUIReviewService(
		aiService,
		reportService,
		fileService,
		zipService,
		cfg.Upload.Dir,
	)
	projectDoctorService := toolservice.NewProjectDoctorService(
		aiService,
		reportService,
		fileService,
		zipService,
		cfg.Upload.Dir,
		cfg.Upload.TempDir,
	)
	apiDocService := toolservice.NewAPIDocService(
		aiService,
		reportService,
		fileService,
		zipService,
		cfg.Upload.Dir,
	)
	toolRunHandler := handler.NewToolRunHandler(
		agentConfigService,
		dbSchemaService,
		uiReviewService,
		projectDoctorService,
		apiDocService,
	)
	authHandler := handler.NewAuthHandler(db, cfg.JWT.Secret, cfg.JWT.Expire)
	adminHandler := handler.NewAdminHandler(db)
	aiGenerationService := service.NewAIGenerationService(db, aiService)
	blueprintHandler := handler.NewBlueprintHandler(db, aiGenerationService)
	requirementHandler := handler.NewRequirementHandler(db)
	broker := sse.NewBroker()
	workspaceService := service.NewWorkspaceService(cfg.Workspace.RootDir)
	taskService := service.NewTaskService(db, broker)
	taskHandler := handler.NewTaskHandler(taskService, broker, workspaceService, aiGenerationService)
	fileHandler := handler.NewFileHandler(workspaceService)
	buildHandler := handler.NewBuildHandler(workspaceService)

	router := gin.New()
	// Cap uploaded bodies at the configured limit and return a clear 413 for
	// oversize JSON or multipart requests. Gin's default MultipartMemory is
	// 32MB; align it with the configured upload cap so large multipart streams
	// are rejected at the boundary rather than streamed to disk unbounded.
	uploadCap := int64(cfg.Upload.MaxUploadSizeMB) << 20
	router.MaxMultipartMemory = uploadCap
	router.Use(
		middleware.RequestID(),
		middleware.Recovery(),
		middleware.Logger(),
		middleware.BodyLimit(uploadCap),
		middleware.CORS(cfg.CORS.AllowOrigins, cfg.App.Env == "production"),
	)

	api := router.Group("/api")
	handler.RegisterHealthRoutes(api, db)
	handler.RegisterSystemRoutes(api, cfg)
	handler.RegisterDashboardRoutes(api, reportService)
	handler.RegisterToolRoutes(api, reportRepo)
	handler.RegisterReportRoutes(api, reportService)
	handler.RegisterProjectRoutes(api, projectService)
	handler.RegisterRequirementRoutes(api, requirementHandler)
	handler.RegisterBlueprintRoutes(api, blueprintHandler)
	handler.RegisterTaskRoutes(api, taskHandler)
	handler.RegisterFileRoutes(api, fileHandler)
	handler.RegisterBuildRoutes(api, buildHandler)
	handler.RegisterExportRoutes(api, exportService)
	handler.RegisterToolRunRoutes(api, toolRunHandler)
	handler.RegisterJobRoutes(api, jobService)
	handler.RegisterObservabilityRoutes(api, aiRunRepo)
	handler.RegisterAuthRoutes(api, authHandler, cfg.JWT.Secret)

	adminAPI := api.Group("")
	adminAPI.Use(middleware.JWTAuth(cfg.JWT.Secret), middleware.RequireAdmin())
	handler.RegisterAdminRoutes(adminAPI, adminHandler)

	return router
}

func newHTTPServer(cfg *config.Config, router http.Handler) *http.Server {
	writeTimeout := time.Duration(cfg.AI.TimeoutSeconds+30) * time.Second
	if writeTimeout < 2*time.Minute {
		writeTimeout = 2 * time.Minute
	}

	return &http.Server{
		Addr:              ":" + cfg.App.Port,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       2 * time.Minute,
		MaxHeaderBytes:    1 << 20,
	}
}
