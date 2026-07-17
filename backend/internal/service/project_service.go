package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"ai-developer-workbench/internal/dto"
	"ai-developer-workbench/internal/model"
	"ai-developer-workbench/internal/repository"
)

// ProjectService handles project lifecycle operations.
type ProjectService interface {
	Create(ctx context.Context, input dto.ProjectCreateDTO) (*dto.ProjectDTO, error)
	Get(ctx context.Context, id string) (*dto.ProjectDTO, error)
	List(ctx context.Context, query dto.ListProjectsQuery) (*dto.PaginatedResponse[dto.ProjectSummaryDTO], error)
	Update(ctx context.Context, id string, input dto.ProjectUpdateDTO) (*dto.ProjectDTO, error)
	Delete(ctx context.Context, id string) (*dto.ProjectDeleteDTO, error)
	GetStats(ctx context.Context, id string) (*dto.ProjectStatsDTO, error)
	ListReports(ctx context.Context, id string, query dto.ListReportsQuery) (*dto.PaginatedResponse[dto.ReportDTO], error)
}

type projectService struct {
	repo repository.ProjectRepository
}

// NewProjectService creates a new project service.
func NewProjectService(repo repository.ProjectRepository) ProjectService {
	return &projectService{repo: repo}
}

func (s *projectService) Create(ctx context.Context, input dto.ProjectCreateDTO) (*dto.ProjectDTO, error) {
	input = normalizeProjectInput(input)
	if err := validateProjectInput(input); err != nil {
		return nil, err
	}

	project := &model.Project{
		Name:          input.Name,
		ProjectType:   input.ProjectType,
		Description:   input.Description,
		RepoURL:       input.RepoURL,
		FrontendStack: input.FrontendStack,
		BackendStack:  input.BackendStack,
		Database:      input.Database,
		UIStyle:       input.UIStyle,
		CodingRules:   input.CodingRules,
	}

	if err := s.repo.Create(ctx, project); err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}
	return s.toDTO(project), nil
}

func (s *projectService) Get(ctx context.Context, id string) (*dto.ProjectDTO, error) {
	project, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.toDTO(project), nil
}

func (s *projectService) List(ctx context.Context, query dto.ListProjectsQuery) (*dto.PaginatedResponse[dto.ProjectSummaryDTO], error) {
	query.SetDefaults()
	projects, total, err := s.repo.List(ctx, query)
	if err != nil {
		return nil, err
	}

	items := make([]dto.ProjectSummaryDTO, 0, len(projects))
	for _, item := range projects {
		p := item.Project
		items = append(items, dto.ProjectSummaryDTO{
			ID:           p.ID,
			Name:         p.Name,
			ProjectType:  p.ProjectType,
			Description:  p.Description,
			RepoURL:      p.RepoURL,
			ReportCount:  item.ReportCount,
			AverageScore: item.AverageScore,
			CreatedAt:    p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:    p.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return &dto.PaginatedResponse[dto.ProjectSummaryDTO]{
		Items:    items,
		Total:    total,
		Page:     query.Page,
		PageSize: query.PageSize,
	}, nil
}

func (s *projectService) Update(ctx context.Context, id string, input dto.ProjectUpdateDTO) (*dto.ProjectDTO, error) {
	project, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	next := dto.ProjectCreateDTO{
		Name:          project.Name,
		ProjectType:   project.ProjectType,
		Description:   project.Description,
		RepoURL:       project.RepoURL,
		FrontendStack: project.FrontendStack,
		BackendStack:  project.BackendStack,
		Database:      project.Database,
		UIStyle:       project.UIStyle,
		CodingRules:   project.CodingRules,
	}
	if input.Name != nil {
		next.Name = *input.Name
	}
	if input.ProjectType != nil {
		next.ProjectType = *input.ProjectType
	}
	if input.Description != nil {
		next.Description = *input.Description
	}
	if input.RepoURL != nil {
		next.RepoURL = *input.RepoURL
	}
	if input.FrontendStack != nil {
		next.FrontendStack = *input.FrontendStack
	}
	if input.BackendStack != nil {
		next.BackendStack = *input.BackendStack
	}
	if input.Database != nil {
		next.Database = *input.Database
	}
	if input.UIStyle != nil {
		next.UIStyle = *input.UIStyle
	}
	if input.CodingRules != nil {
		next.CodingRules = *input.CodingRules
	}
	next = normalizeProjectInput(next)
	if err := validateProjectInput(next); err != nil {
		return nil, err
	}

	project.Name = next.Name
	project.ProjectType = next.ProjectType
	project.Description = next.Description
	project.RepoURL = next.RepoURL
	project.FrontendStack = next.FrontendStack
	project.BackendStack = next.BackendStack
	project.Database = next.Database
	project.UIStyle = next.UIStyle
	project.CodingRules = next.CodingRules

	if err := s.repo.Update(ctx, project); err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}
	return s.toDTO(project), nil
}

func (s *projectService) Delete(ctx context.Context, id string) (*dto.ProjectDeleteDTO, error) {
	detached, err := s.repo.Delete(ctx, id)
	if err != nil {
		return nil, err
	}
	return &dto.ProjectDeleteDTO{DetachedReportCount: detached}, nil
}

func (s *projectService) GetStats(ctx context.Context, id string) (*dto.ProjectStatsDTO, error) {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return nil, err
	}
	return s.repo.GetStats(ctx, id)
}

func (s *projectService) ListReports(ctx context.Context, id string, query dto.ListReportsQuery) (*dto.PaginatedResponse[dto.ReportDTO], error) {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return nil, err
	}
	query.SetDefaults()
	reports, total, err := s.repo.ListReports(ctx, id, query)
	if err != nil {
		return nil, err
	}
	items := make([]dto.ReportDTO, 0, len(reports))
	for _, report := range reports {
		inputData := json.RawMessage(report.InputJSON)
		if len(inputData) == 0 {
			inputData = json.RawMessage(`{}`)
		}
		items = append(items, dto.ReportDTO{
			ID:             report.ID,
			ToolType:       report.ToolType,
			Title:          report.Title,
			InputMode:      report.InputMode,
			Status:         report.Status,
			Summary:        report.Summary,
			TotalScore:     report.TotalScore,
			Grade:          report.Grade,
			InputData:      inputData,
			ReportData:     []byte("{}"),
			GeneratedFiles: []dto.GeneratedFileDTO{},
			ParentReportID: report.ParentReportID,
			ProjectID:      report.ProjectID,
			CreatedAt:      report.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:      report.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}
	return &dto.PaginatedResponse[dto.ReportDTO]{
		Items:    items,
		Total:    total,
		Page:     query.Page,
		PageSize: query.PageSize,
	}, nil
}

func (s *projectService) toDTO(p *model.Project) *dto.ProjectDTO {
	return &dto.ProjectDTO{
		ID:            p.ID,
		Name:          p.Name,
		ProjectType:   p.ProjectType,
		Description:   p.Description,
		RepoURL:       p.RepoURL,
		FrontendStack: p.FrontendStack,
		BackendStack:  p.BackendStack,
		Database:      p.Database,
		UIStyle:       p.UIStyle,
		CodingRules:   p.CodingRules,
		CreatedAt:     p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     p.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func normalizeProjectInput(input dto.ProjectCreateDTO) dto.ProjectCreateDTO {
	input.Name = strings.TrimSpace(input.Name)
	input.ProjectType = strings.TrimSpace(input.ProjectType)
	if input.ProjectType == "" {
		input.ProjectType = "utility_app"
	}
	input.Description = strings.TrimSpace(input.Description)
	input.RepoURL = strings.TrimSpace(input.RepoURL)
	input.FrontendStack = strings.TrimSpace(input.FrontendStack)
	input.BackendStack = strings.TrimSpace(input.BackendStack)
	input.Database = strings.TrimSpace(input.Database)
	input.UIStyle = strings.TrimSpace(input.UIStyle)
	input.CodingRules = strings.TrimSpace(input.CodingRules)
	return input
}

func validateProjectInput(input dto.ProjectCreateDTO) error {
	if input.Name == "" {
		return fmt.Errorf("project name is required")
	}
	if !dto.ValidProjectTypes[input.ProjectType] {
		return fmt.Errorf("project_type is not supported")
	}
	if len(input.Name) > dto.ProjectNameMaxLength {
		return fmt.Errorf("project name must be %d characters or fewer", dto.ProjectNameMaxLength)
	}
	if len(input.Description) > dto.ProjectDescriptionMaxLength {
		return fmt.Errorf("description must be %d characters or fewer", dto.ProjectDescriptionMaxLength)
	}
	if len(input.RepoURL) > dto.ProjectRepoURLMaxLength {
		return fmt.Errorf("repo_url must be %d characters or fewer", dto.ProjectRepoURLMaxLength)
	}
	if input.RepoURL != "" {
		parsed, err := url.ParseRequestURI(input.RepoURL)
		if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
			return fmt.Errorf("repo_url must be a valid http or https URL")
		}
	}
	if len(input.FrontendStack) > dto.ProjectStackMaxLength || len(input.BackendStack) > dto.ProjectStackMaxLength || len(input.UIStyle) > dto.ProjectStackMaxLength {
		return fmt.Errorf("project stack and UI style fields must be %d characters or fewer", dto.ProjectStackMaxLength)
	}
	if len(input.Database) > dto.ProjectDatabaseMaxLength {
		return fmt.Errorf("database must be %d characters or fewer", dto.ProjectDatabaseMaxLength)
	}
	if len(input.CodingRules) > dto.ProjectCodingRulesMaxLength {
		return fmt.Errorf("coding_rules must be %d characters or fewer", dto.ProjectCodingRulesMaxLength)
	}
	return nil
}
