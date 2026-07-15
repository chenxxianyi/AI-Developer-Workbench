package service

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"ai-developer-workbench/internal/model"
	"ai-developer-workbench/internal/util"

	"gorm.io/gorm"
)

const (
	ToolTypeBlueprintGeneration = "blueprint_generation"
	ToolTypeCodeGeneration      = "code_generation"
)

type AIGenerationService struct {
	db *gorm.DB
	ai AIService
}

func NewAIGenerationService(db *gorm.DB, ai AIService) *AIGenerationService {
	return &AIGenerationService{db: db, ai: ai}
}

type BlueprintPageSpec struct {
	Name        string   `json:"name"`
	Route       string   `json:"route"`
	Purpose     string   `json:"purpose,omitempty"`
	KeySections []string `json:"key_sections,omitempty"`
}

type BlueprintAPIEndpointSpec struct {
	Method      string `json:"method"`
	Path        string `json:"path"`
	Description string `json:"description"`
}

type BlueprintDataModelSpec struct {
	Name   string   `json:"name"`
	Fields []string `json:"fields"`
}

type BlueprintAIResult struct {
	ProductPositioning  string                     `json:"product_positioning"`
	TechStack           string                     `json:"tech_stack"`
	UIStyle             string                     `json:"ui_style,omitempty"`
	Pages               []BlueprintPageSpec        `json:"pages"`
	Components          []string                   `json:"components,omitempty"`
	APIEndpoints        []BlueprintAPIEndpointSpec `json:"api_endpoints,omitempty"`
	DataModels          []BlueprintDataModelSpec   `json:"data_models,omitempty"`
	ImplementationNotes []string                   `json:"implementation_notes,omitempty"`
}

type GeneratedProjectFile struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

type CodeGenerationAIResult struct {
	Files []GeneratedProjectFile `json:"files"`
	Notes []string               `json:"notes,omitempty"`
}

func (s *AIGenerationService) GenerateBlueprint(ctx context.Context, projectID string) (*BlueprintAIResult, error) {
	project, err := s.loadProject(projectID)
	if err != nil {
		return nil, err
	}
	requirements := s.loadLatestRequirementContent(projectID)

	systemPrompt := `你是资深产品架构师和全栈技术负责人。你必须只返回合法 JSON 对象，不要使用 Markdown，不要解释。`
	userPrompt := fmt.Sprintf(`请基于以下项目信息生成一个可执行的网站/应用蓝图。

项目信息：
- 名称：%s
- 描述：%s
- 前端技术栈偏好：%s
- 后端技术栈偏好：%s
- 数据库：%s
- UI 风格：%s
- 编码规则：%s

用户需求 JSON：
%s

请返回 JSON，结构必须为：
{
  "product_positioning": "一句话产品定位",
  "tech_stack": "推荐技术栈",
  "ui_style": "视觉风格说明",
  "pages": [{"name":"页面名","route":"/route","purpose":"页面目标","key_sections":["区块"]}],
  "components": ["关键组件"],
  "api_endpoints": [{"method":"GET","path":"/api/example","description":"用途"}],
  "data_models": [{"name":"模型名","fields":["字段"]}],
  "implementation_notes": ["实现注意事项"]
}

约束：
- pages 至少包含首页。
- route 必须以 / 开头。
- 内容必须贴合用户需求，不要返回示例占位。`,
		project.Name,
		util.RedactText(project.Description),
		project.FrontendStack,
		project.BackendStack,
		project.Database,
		project.UIStyle,
		util.RedactText(project.CodingRules),
		emptyAsJSON(requirements),
	)

	aiResult, err := s.ai.GenerateJSON(ctx, AIRequest{
		ToolType:     ToolTypeBlueprintGeneration,
		SystemPrompt: systemPrompt,
		UserPrompt:   userPrompt,
	})
	if err != nil {
		return nil, err
	}

	var result BlueprintAIResult
	if err := util.ParseAIResponseInto(aiResult.JSONText, &result); err != nil {
		return nil, fmt.Errorf("parse blueprint AI response: %w", err)
	}
	normalizeBlueprint(&result, project)
	if err := validateBlueprint(result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *AIGenerationService) GenerateProjectFiles(ctx context.Context, projectID string) (*CodeGenerationAIResult, error) {
	project, err := s.loadProject(projectID)
	if err != nil {
		return nil, err
	}
	requirements := s.loadLatestRequirementContent(projectID)
	blueprint := s.loadLatestBlueprintContent(projectID)
	if strings.TrimSpace(blueprint) == "" {
		return nil, fmt.Errorf("蓝图不存在，请先生成并确认蓝图")
	}

	systemPrompt := `你是资深前端工程师。你必须只返回合法 JSON 对象，不要使用 Markdown，不要解释。`
	userPrompt := fmt.Sprintf(`请基于项目需求和蓝图生成一个可运行的前端项目文件集合。

项目信息：
- 名称：%s
- 描述：%s
- 前端技术栈：%s
- 后端技术栈：%s
- UI 风格：%s
- 编码规则：%s

需求 JSON：
%s

蓝图 JSON：
%s

请返回 JSON，结构必须为：
{
  "files": [
    {"path":"package.json","content":"文件完整内容"},
    {"path":"index.html","content":"文件完整内容"},
    {"path":"src/main.ts","content":"文件完整内容"},
    {"path":"src/App.vue","content":"文件完整内容"}
  ],
  "notes": ["实现说明"]
}

硬性要求：
- files 必须包含 package.json、index.html、src/main.ts、src/App.vue。
- 只能生成相对路径，禁止绝对路径和 ..。
- 文件数量控制在 16 个以内。
- 代码必须完整，不要用省略号，不要写 TODO 占位。
- package.json 必须包含 build 脚本。
- 以 Vue 3 + TypeScript + CSS 为默认实现；如果蓝图指定其他前端栈，也要保证文件自洽。`,
		project.Name,
		util.RedactText(project.Description),
		project.FrontendStack,
		project.BackendStack,
		project.UIStyle,
		util.RedactText(project.CodingRules),
		emptyAsJSON(requirements),
		blueprint,
	)

	aiResult, err := s.ai.GenerateJSON(ctx, AIRequest{
		ToolType:     ToolTypeCodeGeneration,
		SystemPrompt: systemPrompt,
		UserPrompt:   userPrompt,
	})
	if err != nil {
		return nil, err
	}

	var result CodeGenerationAIResult
	if err := util.ParseAIResponseInto(aiResult.JSONText, &result); err != nil {
		return nil, fmt.Errorf("parse code generation AI response: %w", err)
	}
	if err := validateGeneratedFiles(result.Files); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *AIGenerationService) loadProject(projectID string) (*model.Project, error) {
	var project model.Project
	if err := s.db.First(&project, "id = ?", projectID).Error; err != nil {
		return nil, fmt.Errorf("项目不存在: %w", err)
	}
	return &project, nil
}

func (s *AIGenerationService) loadLatestRequirementContent(projectID string) string {
	var req model.Requirement
	if err := s.db.Where("project_id = ?", projectID).Order("version desc").First(&req).Error; err != nil {
		return ""
	}
	return req.Content
}

func (s *AIGenerationService) loadLatestBlueprintContent(projectID string) string {
	var bp model.Blueprint
	if err := s.db.Where("project_id = ?", projectID).Order("version desc").First(&bp).Error; err != nil {
		return ""
	}
	return bp.Content
}

func emptyAsJSON(value string) string {
	if strings.TrimSpace(value) == "" {
		return "{}"
	}
	return value
}

func normalizeBlueprint(result *BlueprintAIResult, project *model.Project) {
	if strings.TrimSpace(result.TechStack) == "" {
		result.TechStack = strings.TrimSpace(strings.Join([]string{project.FrontendStack, project.BackendStack}, " + "))
	}
	if strings.TrimSpace(result.ProductPositioning) == "" {
		result.ProductPositioning = project.Description
	}
	for i := range result.Pages {
		if result.Pages[i].Route == "" {
			result.Pages[i].Route = "/"
		}
		if !strings.HasPrefix(result.Pages[i].Route, "/") {
			result.Pages[i].Route = "/" + result.Pages[i].Route
		}
	}
}

func validateBlueprint(result BlueprintAIResult) error {
	if strings.TrimSpace(result.ProductPositioning) == "" {
		return fmt.Errorf("AI 返回的蓝图缺少 product_positioning")
	}
	if len(result.Pages) == 0 {
		return fmt.Errorf("AI 返回的蓝图缺少 pages")
	}
	return nil
}

func validateGeneratedFiles(files []GeneratedProjectFile) error {
	if len(files) == 0 {
		return fmt.Errorf("AI 未返回任何项目文件")
	}
	if len(files) > 32 {
		return fmt.Errorf("AI 返回文件过多: %d", len(files))
	}

	required := map[string]bool{
		"package.json": false,
		"index.html":   false,
		"src/main.ts":  false,
		"src/App.vue":  false,
	}
	seen := map[string]bool{}
	for _, file := range files {
		path := filepath.ToSlash(strings.TrimSpace(file.Path))
		if path == "" {
			return fmt.Errorf("AI 返回了空文件路径")
		}
		if strings.HasPrefix(path, "/") || strings.Contains(path, "..") || filepath.IsAbs(path) {
			return fmt.Errorf("AI 返回了非法文件路径: %s", file.Path)
		}
		if seen[path] {
			return fmt.Errorf("AI 返回了重复文件路径: %s", path)
		}
		seen[path] = true
		if strings.TrimSpace(file.Content) == "" {
			return fmt.Errorf("AI 返回的文件内容为空: %s", path)
		}
		if _, ok := required[path]; ok {
			required[path] = true
		}
	}
	for path, ok := range required {
		if !ok {
			return fmt.Errorf("AI 返回的文件缺少必需文件: %s", path)
		}
	}
	return nil
}

func MarshalBlueprintContent(result *BlueprintAIResult) (string, error) {
	bytes, err := json.Marshal(result)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
