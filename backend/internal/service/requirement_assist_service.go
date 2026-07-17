package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"ai-developer-workbench/internal/dto"
	"ai-developer-workbench/internal/model"
	"ai-developer-workbench/internal/util"

	"gorm.io/gorm"
)

const ToolTypeRequirementAssist = "requirement_assist"

type RequirementAssistSpec struct {
	SchemaVersion             int                      `json:"schema_version"`
	AppType                   string                   `json:"app_type"`
	Goal                      string                   `json:"goal"`
	TargetUsers               []string                 `json:"target_users"`
	PrimaryScenarios          []string                 `json:"primary_scenarios"`
	MustHaveFeatures          []string                 `json:"must_have_features"`
	ShouldHaveFeatures        []string                 `json:"should_have_features"`
	Screens                   []string                 `json:"screens"`
	InteractionRules          []string                 `json:"interaction_rules"`
	DataAndStorage            RequirementAssistStorage `json:"data_and_storage"`
	VisualPreferences         RequirementAssistVisual  `json:"visual_preferences"`
	ResponsiveTargets         []string                 `json:"responsive_targets"`
	NonFunctionalRequirements []string                 `json:"non_functional_requirements"`
	AcceptanceCriteria        []string                 `json:"acceptance_criteria"`
	OutOfScope                []string                 `json:"out_of_scope"`
}

type RequirementAssistStorage struct {
	Persistence     string `json:"persistence"`
	BackendRequired bool   `json:"backend_required"`
}

type RequirementAssistVisual struct {
	Style        string `json:"style"`
	PrimaryColor string `json:"primary_color"`
}

type RequirementAssistQuestion struct {
	ID       string   `json:"id"`
	Question string   `json:"question"`
	Options  []string `json:"options"`
	Blocking bool     `json:"blocking"`
}

type RequirementAssistSummary struct {
	Product  string   `json:"product"`
	Users    string   `json:"users"`
	Features []string `json:"features"`
	Success  []string `json:"success"`
}

type RequirementAssistInput struct {
	Description     string                `json:"description"`
	CurrentSpec     RequirementAssistSpec `json:"current_spec"`
	CurrentStep     string                `json:"current_step"`
	ConfirmedFields []string              `json:"confirmed_fields"`
}

type RequirementAssistResult struct {
	Spec           RequirementAssistSpec       `json:"spec"`
	InferredFields []string                    `json:"inferred_fields"`
	Questions      []RequirementAssistQuestion `json:"questions"`
	Ready          bool                        `json:"ready"`
	Summary        RequirementAssistSummary    `json:"summary"`
}

type RequirementAssistService struct {
	db *gorm.DB
	ai AIService
}

func NewRequirementAssistService(db *gorm.DB, ai AIService) *RequirementAssistService {
	return &RequirementAssistService{db: db, ai: ai}
}

func (s *RequirementAssistService) Assist(ctx context.Context, projectID string, input RequirementAssistInput) (*RequirementAssistResult, error) {
	var project model.Project
	if err := s.db.First(&project, "id = ?", projectID).Error; err != nil {
		return nil, fmt.Errorf("load project: %w", err)
	}
	description := strings.TrimSpace(input.Description)
	if description == "" {
		description = strings.TrimSpace(input.CurrentSpec.Goal)
	}
	if description == "" {
		return nil, fmt.Errorf("请先简单描述想做的产品")
	}
	currentJSON, _ := json.Marshal(input.CurrentSpec)
	systemPrompt := `你是面向普通用户的产品需求助手。你必须只返回合法 JSON，不使用 Markdown，不解释。使用用户所用语言，问题必须口语化，不能出现持久化、非功能、领域模型等专业术语。`
	userPrompt := fmt.Sprintf(`请把用户的自然语言想法整理成可生成应用的结构化需求。

项目名称：%s
当前项目类型：%s
用户描述：%s
当前步骤：%s
当前需求（用户已填写的内容优先）：%s

返回格式：
{
  "spec": {
    "schema_version": 2,
    "app_type": "interactive_app|dashboard|data_product|content_site|ecommerce|utility_app|landing_page|analysis_existing",
    "goal": "一句清晰目标",
    "target_users": ["普通语言用户类型"],
    "primary_scenarios": ["用户行为"],
    "must_have_features": ["必须功能"],
    "should_have_features": ["可选功能"],
    "screens": ["页面或屏幕"],
    "interaction_rules": ["可验证规则"],
    "data_and_storage": {"persistence":"不用保存|浏览器本地存储|登录后保存|跨设备同步", "backend_required":false},
    "visual_preferences": {"style":"视觉风格", "primary_color":"颜色偏好"},
    "responsive_targets": ["desktop","mobile"],
    "non_functional_requirements": ["普通语言质量要求"],
    "acceptance_criteria": ["完成后可以验证的行为"],
    "out_of_scope": ["本次暂不实现"]
  },
  "inferred_fields": ["AI 补充的字段名"],
  "questions": [{"id":"stable-id","question":"普通用户能理解的问题","options":["选项一","选项二","交给 AI 推荐"],"blocking":false}],
  "ready": true,
  "summary": {"product":"产品摘要","users":"用户摘要","features":["核心功能"],"success":["完成效果"]}
}

规则：
- 最多返回 3 个问题，只有会改变核心方向的问题才 blocking=true。
- 不得擅自增加收费、客户案例、销售联系表单、登录或远程 API。
- 非 landing_page 不得默认生成营销区块。
- 必须功能控制在 3-8 项，每项用用户能理解的行为名称。
- 每个必须功能都要有对应的验收条件。
- 对游戏应包含真实规则、回合、重新开始以及需求明确时的电脑对手。
- 当前需求中已有的非空内容不得被改写或删除。`,
		project.Name, project.ProjectType, util.RedactText(description), input.CurrentStep, string(currentJSON))
	result, err := s.ai.GenerateJSON(ctx, AIRequest{ToolType: ToolTypeRequirementAssist, SystemPrompt: systemPrompt, UserPrompt: userPrompt})
	if err != nil {
		return nil, fmt.Errorf("AI 整理需求失败: %w", err)
	}
	var assisted RequirementAssistResult
	if err := util.ParseAIResponseInto(result.JSONText, &assisted); err != nil {
		return nil, fmt.Errorf("解析 AI 需求建议失败: %w", err)
	}
	preserveConfirmedRequirementFields(&assisted.Spec, input.CurrentSpec, input.ConfirmedFields)
	normalizeRequirementAssist(&assisted, project.ProjectType)
	return &assisted, nil
}

func preserveConfirmedRequirementFields(next *RequirementAssistSpec, current RequirementAssistSpec, fields []string) {
	for _, field := range fields {
		switch field {
		case "app_type":
			if current.AppType != "" {
				next.AppType = current.AppType
			}
		case "goal":
			if current.Goal != "" {
				next.Goal = current.Goal
			}
		case "target_users":
			if len(current.TargetUsers) > 0 {
				next.TargetUsers = current.TargetUsers
			}
		case "primary_scenarios":
			if len(current.PrimaryScenarios) > 0 {
				next.PrimaryScenarios = current.PrimaryScenarios
			}
		case "must_have_features":
			if len(current.MustHaveFeatures) > 0 {
				next.MustHaveFeatures = current.MustHaveFeatures
			}
		case "should_have_features":
			if len(current.ShouldHaveFeatures) > 0 {
				next.ShouldHaveFeatures = current.ShouldHaveFeatures
			}
		case "visual_preferences":
			if current.VisualPreferences.Style != "" || current.VisualPreferences.PrimaryColor != "" {
				next.VisualPreferences = current.VisualPreferences
			}
		case "responsive_targets":
			if len(current.ResponsiveTargets) > 0 {
				next.ResponsiveTargets = current.ResponsiveTargets
			}
		}
	}
}

func normalizeRequirementAssist(result *RequirementAssistResult, fallbackType string) {
	result.Spec.SchemaVersion = 2
	if !dto.ValidProjectTypes[result.Spec.AppType] {
		result.Spec.AppType = fallbackType
	}
	result.Spec.TargetUsers = compactRequirementItems(result.Spec.TargetUsers, 8)
	result.Spec.PrimaryScenarios = compactRequirementItems(result.Spec.PrimaryScenarios, 10)
	result.Spec.MustHaveFeatures = compactRequirementItems(result.Spec.MustHaveFeatures, 8)
	result.Spec.ShouldHaveFeatures = compactRequirementItems(result.Spec.ShouldHaveFeatures, 8)
	result.Spec.Screens = compactRequirementItems(result.Spec.Screens, 12)
	result.Spec.InteractionRules = compactRequirementItems(result.Spec.InteractionRules, 16)
	result.Spec.AcceptanceCriteria = compactRequirementItems(result.Spec.AcceptanceCriteria, 16)
	result.Spec.OutOfScope = compactRequirementItems(result.Spec.OutOfScope, 10)
	result.Spec.NonFunctionalRequirements = compactRequirementItems(result.Spec.NonFunctionalRequirements, 10)
	if len(result.Spec.ResponsiveTargets) == 0 {
		result.Spec.ResponsiveTargets = []string{"desktop", "mobile"}
	}
	if len(result.Questions) > 3 {
		result.Questions = result.Questions[:3]
	}
	for i := range result.Questions {
		result.Questions[i].Options = compactRequirementItems(result.Questions[i].Options, 4)
	}
	if strings.TrimSpace(result.Summary.Product) == "" {
		result.Summary.Product = result.Spec.Goal
	}
	if strings.TrimSpace(result.Summary.Users) == "" {
		result.Summary.Users = strings.Join(result.Spec.TargetUsers, "、")
	}
	if len(result.Summary.Features) == 0 {
		result.Summary.Features = append([]string(nil), result.Spec.MustHaveFeatures...)
	}
	if len(result.Summary.Success) == 0 {
		result.Summary.Success = append([]string(nil), result.Spec.AcceptanceCriteria...)
	}
	blocking := false
	for _, question := range result.Questions {
		blocking = blocking || question.Blocking
	}
	result.Ready = strings.TrimSpace(result.Spec.Goal) != "" && len(result.Spec.TargetUsers) > 0 && len(result.Spec.MustHaveFeatures) > 0 && len(result.Spec.AcceptanceCriteria) > 0 && !blocking
}

func compactRequirementItems(items []string, limit int) []string {
	seen := map[string]bool{}
	result := make([]string, 0, len(items))
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" || seen[item] {
			continue
		}
		seen[item] = true
		result = append(result, item)
		if len(result) == limit {
			break
		}
	}
	return result
}
