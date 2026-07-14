package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"ai-developer-workbench/internal/dto"
	"ai-developer-workbench/internal/model"
	"ai-developer-workbench/internal/prompts"
	"ai-developer-workbench/internal/service"
	"ai-developer-workbench/internal/util"
)

// AgentConfigService handles Agent Config Studio operations.
type AgentConfigService struct {
	aiService     service.AIService
	reportService service.ReportService
}

// NewAgentConfigService creates a new Agent Config service.
func NewAgentConfigService(aiService service.AIService, reportService service.ReportService) *AgentConfigService {
	return &AgentConfigService{
		aiService:     aiService,
		reportService: reportService,
	}
}

// Run executes the Agent Config Studio tool.
func (s *AgentConfigService) Run(ctx context.Context, req dto.AgentConfigRequest) (*dto.ReportDTO, error) {
	if req.ParentReportID != "" {
		if _, err := s.reportService.ValidateParentReport(ctx, model.ToolTypeAgentConfig, req.ParentReportID); err != nil {
			return nil, fmt.Errorf("invalid parent report: %w", err)
		}
	}
	project, err := s.reportService.ResolveProject(ctx, req.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project: %w", err)
	}
	// 1. Create processing report.
	inputData, _ := json.Marshal(req)
	report, err := s.reportService.CreateProcessingReport(ctx, model.ToolTypeAgentConfig, req.Title, "json", inputData, req.ParentReportID, req.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("failed to create report: %w", err)
	}

	// 2. Build prompt. Redact secrets in free-text fields before reaching the AI.
	systemPrompt, userPrompt := prompts.BuildAgentConfigPrompt(
		req.ProjectName, req.ProjectType, req.FrontendStack,
		req.BackendStack, req.Database, req.UIStyle,
		util.RedactText(req.CodingPreferences), util.RedactText(req.StrictRules),
	)
	userPrompt = prompts.AppendTrustedProjectContext(userPrompt, project)

	// 3. Call AI service.
	aiResult, err := s.aiService.GenerateJSON(ctx, service.AIRequest{
		ToolType:     model.ToolTypeAgentConfig,
		SystemPrompt: systemPrompt,
		UserPrompt:   userPrompt,
	})
	if err != nil {
		slog.Error("AI call failed for Agent Config", "error", err)
		_ = s.reportService.FailReport(ctx, report.ID, err.Error())
		return nil, err
	}

	// 4. Parse AI response.
	var result dto.AgentConfigResult
	if err := util.ParseAIResponseInto(aiResult.JSONText, &result); err != nil {
		parseErr := fmt.Errorf("parse AI response: %w", err)
		slog.Error("Failed to parse AI response", "error", err)
		_ = s.reportService.FailReport(ctx, report.ID, parseErr.Error())
		return nil, parseErr
	}

	// 5. Normalize result.
	s.normalizeResult(&result)

	// 6. Build generated files.
	generatedFiles := s.buildGeneratedFiles(result)

	// 7. Save succeeded report.
	summary := fmt.Sprintf("Generated %d configuration files for %s", len(result.GeneratedFilesContent), req.ProjectName)
	reportJSON, _ := json.Marshal(result)

	return s.reportService.SucceedReport(ctx, report.ID, reportJSON, summary, nil, nil, generatedFiles)
}

// normalizeResult normalizes the result data.
func (s *AgentConfigService) normalizeResult(result *dto.AgentConfigResult) {
	if result.Recommendations == nil {
		result.Recommendations = []string{}
	}
	if result.GeneratedFilesContent == nil {
		result.GeneratedFilesContent = map[string]string{}
	}
	result.ActionItems = dto.NormalizeActionItems(result.ActionItems)
	// Validate filenames.
	for filename := range result.GeneratedFilesContent {
		if !util.IsAllowedGeneratedFilename(filename) {
			slog.Warn("Unexpected generated filename", "filename", filename)
		}
	}
}

// buildGeneratedFiles creates model.GeneratedFile entries from the result.
func (s *AgentConfigService) buildGeneratedFiles(result dto.AgentConfigResult) []model.GeneratedFile {
	files := make([]model.GeneratedFile, 0, len(result.GeneratedFilesContent))
	for filename, content := range result.GeneratedFilesContent {
		mimeType := "text/markdown"
		language := "markdown"
		if filename == "openapi.json" {
			mimeType = "application/json"
			language = "json"
		} else if filename == "migration.sql" {
			mimeType = "text/x-sql"
			language = "sql"
		}
		files = append(files, model.GeneratedFile{
			Filename: filename,
			Content:  content,
			MimeType: mimeType,
			Language: language,
		})
	}
	return files
}
