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
	ToolTypeGenerationPlan      = "generation_plan"
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

type BlueprintUserFlowSpec struct {
	Name  string   `json:"name"`
	Steps []string `json:"steps"`
}

type BlueprintFeatureSpec struct {
	ID                 string   `json:"id"`
	Name               string   `json:"name"`
	Priority           string   `json:"priority"`
	Description        string   `json:"description"`
	AcceptanceCriteria []string `json:"acceptance_criteria"`
}

type BlueprintStateSpec struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type BlueprintVisualSystemSpec struct {
	Style         string `json:"style"`
	Colors        string `json:"colors"`
	Layout        string `json:"layout"`
	Accessibility string `json:"accessibility"`
}

type BlueprintAIResult struct {
	SchemaVersion       int                        `json:"schema_version"`
	AppType             string                     `json:"app_type"`
	ProductPositioning  string                     `json:"product_positioning"`
	TechStack           string                     `json:"tech_stack"`
	UIStyle             string                     `json:"ui_style,omitempty"`
	Pages               []BlueprintPageSpec        `json:"pages"`
	Components          []string                   `json:"components,omitempty"`
	UserFlows           []BlueprintUserFlowSpec    `json:"user_flows,omitempty"`
	Features            []BlueprintFeatureSpec     `json:"features,omitempty"`
	InteractionRules    []string                   `json:"interaction_rules,omitempty"`
	StateModel          []BlueprintStateSpec       `json:"state_model,omitempty"`
	DomainRules         []string                   `json:"domain_rules,omitempty"`
	APIEndpoints        []BlueprintAPIEndpointSpec `json:"api_endpoints,omitempty"`
	DataModels          []BlueprintDataModelSpec   `json:"data_models,omitempty"`
	VisualSystem        BlueprintVisualSystemSpec  `json:"visual_system,omitempty"`
	AcceptanceCriteria  []string                   `json:"acceptance_criteria,omitempty"`
	TestPlan            []string                   `json:"test_plan,omitempty"`
	ImplementationNotes []string                   `json:"implementation_notes,omitempty"`
	OpenQuestions       []string                   `json:"open_questions,omitempty"`
}

type GeneratedProjectFile struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

type CodeGenerationAIResult struct {
	Files []GeneratedProjectFile `json:"files"`
	Notes []string               `json:"notes,omitempty"`
}

type GenerationPlanFile struct {
	Path         string   `json:"path"`
	Purpose      string   `json:"purpose"`
	FeatureIDs   []string `json:"feature_ids"`
	Dependencies []string `json:"dependencies,omitempty"`
}

type GenerationPlan struct {
	ApplicationType   string               `json:"application_type"`
	Architecture      string               `json:"architecture"`
	Files             []GenerationPlanFile `json:"files"`
	VerificationSteps []string             `json:"verification_steps"`
}

type CodeGenerationContentSpec struct {
	ProductName     string                  `json:"product_name"`
	Tagline         string                  `json:"tagline"`
	HeroTitle       string                  `json:"hero_title"`
	HeroSubtitle    string                  `json:"hero_subtitle"`
	PrimaryCTA      string                  `json:"primary_cta"`
	SecondaryCTA    string                  `json:"secondary_cta"`
	Metrics         []CodeGenerationMetric  `json:"metrics"`
	Features        []CodeGenerationFeature `json:"features"`
	Cases           []CodeGenerationCase    `json:"cases"`
	Pricing         []CodeGenerationPlan    `json:"pricing"`
	ContactTitle    string                  `json:"contact_title"`
	ContactSubtitle string                  `json:"contact_subtitle"`
	Notes           []string                `json:"notes,omitempty"`
}

type CodeGenerationMetric struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type CodeGenerationFeature struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CodeGenerationCase struct {
	Name   string `json:"name"`
	Result string `json:"result"`
}

type CodeGenerationPlan struct {
	Name        string   `json:"name"`
	Price       string   `json:"price"`
	Description string   `json:"description"`
	Items       []string `json:"items"`
}

func (s *AIGenerationService) GenerateBlueprint(ctx context.Context, projectID string) (*BlueprintAIResult, error) {
	project, err := s.loadProject(projectID)
	if err != nil {
		return nil, err
	}
	requirements := s.loadLatestRequirementContent(projectID)

	systemPrompt := `你是资深产品架构师和全栈技术负责人。你必须只返回合法 JSON 对象，不要使用 Markdown，不要解释。`
	userPrompt := fmt.Sprintf(`请基于以下项目信息生成一个可执行、可评审、可测试的应用蓝图。

项目信息：
- 名称：%s
- 项目类型：%s
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
  "schema_version": 2,
  "app_type": "项目类型",
  "product_positioning": "一句话产品定位",
  "tech_stack": "推荐技术栈",
  "ui_style": "视觉风格说明",
  "pages": [{"name":"页面名","route":"/route","purpose":"页面目标","key_sections":["区块"]}],
  "user_flows": [{"name":"流程名","steps":["步骤"]}],
  "features": [{"id":"F-001","name":"功能名","priority":"must|should","description":"具体行为","acceptance_criteria":["可验证条件"]}],
  "components": ["关键组件"],
  "interaction_rules": ["交互规则"],
  "state_model": [{"name":"状态名","description":"状态含义和变化"}],
  "domain_rules": ["业务或领域规则"],
  "api_endpoints": [{"method":"GET","path":"/api/example","description":"用途"}],
  "data_models": [{"name":"模型名","fields":["字段"]}],
  "visual_system": {"style":"视觉风格","colors":"颜色策略","layout":"布局策略","accessibility":"无障碍要求"},
  "acceptance_criteria": ["全局验收标准"],
  "test_plan": ["核心测试"],
  "implementation_notes": ["实现注意事项"],
  "open_questions": []
}

约束：
- pages 至少包含首页。
- route 必须以 / 开头。
- 每项 must 功能必须有明确、可执行、可测试的行为和验收标准。
- 互动应用必须描述核心状态、领域规则和交互流程；管理后台必须描述导航、列表、表单和数据状态。
- 只有 landing_page 或用户明确提出时才允许默认生成价格、客户案例和联系表单。
- 不得虚构用户数量、价格、客户或业务数据。
- 内容必须贴合用户需求，不要返回示例占位。`,
		project.Name,
		project.ProjectType,
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
	requirements := safePromptJSON(s.loadLatestRequirementContent(projectID), 3500)
	blueprint := s.loadConfirmedBlueprintContent(projectID)
	if strings.TrimSpace(blueprint) == "" {
		return nil, fmt.Errorf("confirmed blueprint is missing; generate, review and confirm a blueprint first")
	}
	blueprint = compactBlueprintForCodegen(blueprint)
	plan, err := s.generateProjectPlan(ctx, project, requirements, blueprint)
	if err != nil {
		return nil, err
	}
	planJSON, _ := json.Marshal(plan)

	systemPrompt := `You are a senior Vue 3 product engineer. Generate complete, runnable application source files. Return one valid JSON object only, without Markdown or explanations.`
	userPrompt := fmt.Sprintf(`Generate a functional Vue 3 + TypeScript application based strictly on the confirmed requirement and blueprint.

Project info:
- Name: %s
- Application type: %s
- Description: %s
- UI style: %s
- Coding rules: %s

Requirement summary JSON:
%s

Blueprint summary JSON:
%s

Approved generation plan JSON:
%s

Return JSON exactly in this shape:
{
  "files": [
    {"path":"package.json","content":"complete file content"},
    {"path":"tsconfig.json","content":"complete file content"},
    {"path":"index.html","content":"complete file content"},
    {"path":"src/main.ts","content":"complete file content"},
    {"path":"src/App.vue","content":"complete file content"}
  ],
  "notes": ["implemented scope or known limitation"]
}

Hard requirements:
- Generate an actual application, not a page that merely describes the requested product.
- Implement every must-have feature and acceptance criterion from the confirmed inputs.
- Use multiple focused files for domain logic, state and reusable components when the application is interactive.
- Include package.json, tsconfig.json, index.html, src/main.ts and src/App.vue. Keep the total at 24 files or fewer.
- package.json must use pinned compatible versions and only scripts named dev, build, preview and test. Never add lifecycle scripts.
- Include vue-tsc so the build service can run an explicit typecheck. Allowed packages are vue, vue-router, pinia, vite, typescript, vue-tsc, @vitejs/plugin-vue, vitest, @vue/test-utils and jsdom. Prefer no extra dependency.
- Use simple non-shell scripts: vite, vite build, vite preview and vitest run. The test script must terminate without watch mode.
- Do not use external APIs, remote assets, secrets, fake metrics, invented pricing or customer claims unless explicitly required.
- For interactive applications, implement real state transitions and domain behavior. For games, include legal actions, turn state, restart and a deterministic playable opponent when required.
- For dashboards, implement navigation, meaningful local sample data, filters, loading/empty states and forms described by the blueprint.
- Only landing_page applications may default to marketing sections; other application types must prioritize the core task UI.
- Use accessible semantic HTML, responsive CSS and visible keyboard focus.
- Do not wrap file content in Markdown fences.
- Use the same language as the requirement and blueprint.`,
		project.Name,
		project.ProjectType,
		util.RedactText(project.Description),
		project.UIStyle,
		util.RedactText(project.CodingRules),
		requirements,
		blueprint,
		string(planJSON),
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

func (s *AIGenerationService) generateProjectPlan(ctx context.Context, project *model.Project, requirements, blueprint string) (*GenerationPlan, error) {
	systemPrompt := `You are a senior software architect. Return one valid JSON object only. Plan a compact Vue application before code generation.`
	userPrompt := fmt.Sprintf(`Create a file and verification plan for this application.

Project: %s
Application type: %s
Requirements: %s
Confirmed blueprint: %s

Return:
{"application_type":"type","architecture":"short architecture summary","files":[{"path":"src/example.ts","purpose":"why this file exists","feature_ids":["F-001"],"dependencies":[]}],"verification_steps":["test behavior"]}

Constraints:
- Include package.json, tsconfig.json, index.html, src/main.ts and src/App.vue.
- Plan 24 files or fewer and map every must feature to at least one file.
- Interactive applications must separate domain rules/state from visual components.
- Do not plan pricing, contact forms or marketing proof unless explicitly required.`,
		project.Name, project.ProjectType, requirements, blueprint)
	result, err := s.ai.GenerateJSON(ctx, AIRequest{ToolType: ToolTypeGenerationPlan, SystemPrompt: systemPrompt, UserPrompt: userPrompt})
	if err != nil {
		return nil, fmt.Errorf("generate project plan: %w", err)
	}
	var plan GenerationPlan
	if err := util.ParseAIResponseInto(result.JSONText, &plan); err != nil {
		return nil, fmt.Errorf("parse generation plan: %w", err)
	}
	if len(plan.Files) == 0 || len(plan.Files) > 24 {
		return nil, fmt.Errorf("generation plan must contain 1 to 24 files")
	}
	required := map[string]bool{"package.json": false, "tsconfig.json": false, "index.html": false, "src/main.ts": false, "src/App.vue": false}
	for _, file := range plan.Files {
		path := filepath.ToSlash(strings.TrimSpace(file.Path))
		if _, ok := required[path]; ok {
			required[path] = true
		}
	}
	for path, found := range required {
		if !found {
			return nil, fmt.Errorf("generation plan is missing %s", path)
		}
	}
	return &plan, nil
}

// RepairProjectFiles performs a scoped repair after deterministic verification
// fails. The returned result is the complete merged project, not only patches.
func (s *AIGenerationService) RepairProjectFiles(ctx context.Context, projectID string, current *CodeGenerationAIResult, verificationError string) (*CodeGenerationAIResult, error) {
	project, err := s.loadProject(projectID)
	if err != nil {
		return nil, err
	}
	type repairFile struct {
		Path    string `json:"path"`
		Content string `json:"content"`
	}
	filesForPrompt := make([]repairFile, 0, len(current.Files))
	remaining := 70000
	for _, file := range current.Files {
		if remaining <= 0 {
			break
		}
		content := util.TruncateText(file.Content, min(remaining, 14000))
		remaining -= len(content)
		filesForPrompt = append(filesForPrompt, repairFile{Path: file.Path, Content: content})
	}
	currentJSON, _ := json.Marshal(filesForPrompt)
	systemPrompt := `You repair Vue 3 + TypeScript projects from deterministic test or build errors. Return one valid JSON object only, without Markdown.`
	userPrompt := fmt.Sprintf(`Repair only files related to the verification error. Preserve working functionality and the confirmed product scope.

Project: %s
Application type: %s

Verification error:
%s

Current project files (some large files may be truncated):
%s

Return only changed complete files:
{"files":[{"path":"relative/path","content":"complete corrected content"}],"notes":["repair summary"]}

Constraints:
- Do not add dependencies unless absolutely required and allowed by the existing package policy.
- Do not remove requested functionality to make the build pass.
- Never use Markdown fences, absolute paths or path traversal.
- Keep changes scoped to the error.`,
		project.Name,
		project.ProjectType,
		util.TruncateText(verificationError, 10000),
		string(currentJSON),
	)
	aiResult, err := s.ai.GenerateJSON(ctx, AIRequest{ToolType: ToolTypeCodeGeneration, SystemPrompt: systemPrompt, UserPrompt: userPrompt})
	if err != nil {
		return nil, err
	}
	var patch CodeGenerationAIResult
	if err := util.ParseAIResponseInto(aiResult.JSONText, &patch); err != nil {
		return nil, fmt.Errorf("parse repair AI response: %w", err)
	}
	if len(patch.Files) == 0 {
		return nil, fmt.Errorf("AI repair returned no changed files")
	}
	merged := make(map[string]GeneratedProjectFile, len(current.Files)+len(patch.Files))
	for _, file := range current.Files {
		merged[filepath.ToSlash(file.Path)] = file
	}
	for _, file := range patch.Files {
		path := filepath.ToSlash(strings.TrimSpace(file.Path))
		if path == "" || strings.HasPrefix(path, "/") || strings.Contains(path, "..") || filepath.IsAbs(path) || strings.TrimSpace(file.Content) == "" {
			return nil, fmt.Errorf("AI repair returned an invalid file: %s", file.Path)
		}
		if _, exists := merged[path]; !exists {
			return nil, fmt.Errorf("AI repair attempted to add an unplanned file: %s", path)
		}
		file.Path = path
		merged[path] = file
	}
	result := &CodeGenerationAIResult{Notes: append(current.Notes, patch.Notes...)}
	for _, original := range current.Files {
		result.Files = append(result.Files, merged[filepath.ToSlash(original.Path)])
	}
	if err := validateGeneratedFiles(result.Files); err != nil {
		return nil, fmt.Errorf("validate repaired files: %w", err)
	}
	return result, nil
}

func normalizeContentSpec(spec *CodeGenerationContentSpec, project *model.Project) {
	// Legacy landing-page content normalization is intentionally disabled.
	// New generation parses CodeGenerationAIResult directly.
	_ = spec
	_ = project
	return

	// Deprecated implementation below is unreachable.
	if strings.TrimSpace(spec.ProductName) == "" {
		spec.ProductName = project.Name
	}
	if strings.TrimSpace(spec.Tagline) == "" {
		spec.Tagline = "AI generated product experience"
	}
	if strings.TrimSpace(spec.HeroTitle) == "" {
		spec.HeroTitle = project.Name
	}
	if strings.TrimSpace(spec.HeroSubtitle) == "" {
		spec.HeroSubtitle = project.Description
	}
	if strings.TrimSpace(spec.PrimaryCTA) == "" {
		spec.PrimaryCTA = "Get started"
	}
	if strings.TrimSpace(spec.SecondaryCTA) == "" {
		spec.SecondaryCTA = "View details"
	}
	if len(spec.Metrics) == 0 {
		spec.Metrics = []CodeGenerationMetric{{Value: "24/7", Label: "Always on"}, {Value: "3x", Label: "Faster delivery"}, {Value: "99%", Label: "Reliable experience"}}
	}
	if len(spec.Features) == 0 {
		spec.Features = []CodeGenerationFeature{{Title: "Smart workflow", Description: "Turn project requirements into a clear and usable digital experience."}, {Title: "Responsive design", Description: "A polished layout that works across desktop and mobile devices."}, {Title: "Conversion focused", Description: "Clear sections guide visitors from discovery to contact."}}
	}
	if len(spec.Cases) == 0 {
		spec.Cases = []CodeGenerationCase{{Name: "Growth team", Result: "Launched a focused web experience faster."}, {Name: "Operations team", Result: "Improved visitor understanding and lead quality."}}
	}
	if len(spec.Pricing) == 0 {
		spec.Pricing = []CodeGenerationPlan{{Name: "Starter", Price: "Contact us", Description: "For teams validating the idea.", Items: []string{"Core page", "Responsive layout", "Lead form"}}, {Name: "Growth", Price: "Custom", Description: "For teams ready to scale.", Items: []string{"Advanced sections", "Case studies", "Priority support"}}, {Name: "Enterprise", Price: "Custom", Description: "For complex organizations.", Items: []string{"Custom workflow", "Integration planning", "Dedicated success"}}}
	}
	if strings.TrimSpace(spec.ContactTitle) == "" {
		spec.ContactTitle = "Talk to us"
	}
	if strings.TrimSpace(spec.ContactSubtitle) == "" {
		spec.ContactSubtitle = "Share your goals and we will follow up with a tailored plan."
	}
	if len(spec.Metrics) > 3 {
		spec.Metrics = spec.Metrics[:3]
	}
	if len(spec.Features) > 6 {
		spec.Features = spec.Features[:6]
	}
	if len(spec.Cases) > 3 {
		spec.Cases = spec.Cases[:3]
	}
	if len(spec.Pricing) > 3 {
		spec.Pricing = spec.Pricing[:3]
	}
}

func renderVueProjectFiles(project *model.Project, spec CodeGenerationContentSpec) ([]GeneratedProjectFile, error) {
	// Retained temporarily for source compatibility with old branches. The
	// generic generation pipeline must never fall back to this marketing page.
	_ = project
	_ = spec
	return nil, fmt.Errorf("legacy landing-page renderer is disabled")

	// Deprecated implementation below is unreachable and will be removed after
	// legacy branch compatibility is no longer required.
	features, err := json.MarshalIndent(spec.Features, "", "  ")
	if err != nil {
		return nil, err
	}
	metrics, err := json.MarshalIndent(spec.Metrics, "", "  ")
	if err != nil {
		return nil, err
	}
	cases, err := json.MarshalIndent(spec.Cases, "", "  ")
	if err != nil {
		return nil, err
	}
	pricing, err := json.MarshalIndent(spec.Pricing, "", "  ")
	if err != nil {
		return nil, err
	}

	packageJSON := `{"scripts":{"dev":"vite --host 0.0.0.0","build":"vite build","preview":"vite preview --host 0.0.0.0"},"dependencies":{"@vitejs/plugin-vue":"latest","typescript":"latest","vite":"latest","vue":"latest"},"devDependencies":{}}`
	indexHTML := fmt.Sprintf(`<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>%s</title>
  </head>
  <body>
    <div id="app"></div>
    <script type="module" src="/src/main.ts"></script>
  </body>
</html>
`, htmlText(spec.ProductName))
	mainTS := `import { createApp } from 'vue'
import App from './App.vue'

createApp(App).mount('#app')
`
	viteConfig := `import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
})
`
	appVue := fmt.Sprintf(`<script setup lang="ts">
const productName = %s
const tagline = %s
const heroTitle = %s
const heroSubtitle = %s
const primaryCTA = %s
const secondaryCTA = %s
const contactTitle = %s
const contactSubtitle = %s
const features = %s
const metrics = %s
const cases = %s
const pricing = %s
</script>

<template>
  <main class="page-shell">
    <section class="hero-section">
      <nav class="nav-bar">
        <div class="brand"><span class="brand-mark">AI</span>{{ productName }}</div>
        <div class="nav-actions"><a href="#features">Features</a><a href="#pricing">Pricing</a><a href="#contact">Contact</a></div>
      </nav>
      <div class="hero-grid">
        <div class="hero-copy">
          <p class="eyebrow">{{ tagline }}</p>
          <h1>{{ heroTitle }}</h1>
          <p class="hero-subtitle">{{ heroSubtitle }}</p>
          <div class="cta-row"><a class="btn primary" href="#contact">{{ primaryCTA }}</a><a class="btn ghost" href="#features">{{ secondaryCTA }}</a></div>
        </div>
        <div class="hero-card" aria-label="Product highlights">
          <div v-for="metric in metrics" :key="metric.label" class="metric"><strong>{{ metric.value }}</strong><span>{{ metric.label }}</span></div>
        </div>
      </div>
    </section>

    <section id="features" class="section">
      <div class="section-heading"><p class="eyebrow">Core capabilities</p><h2>Built for measurable outcomes</h2></div>
      <div class="feature-grid"><article v-for="feature in features" :key="feature.title" class="feature-card"><h3>{{ feature.title }}</h3><p>{{ feature.description }}</p></article></div>
    </section>

    <section class="section proof-section">
      <div class="section-heading"><p class="eyebrow">Customer proof</p><h2>Trusted by teams that move fast</h2></div>
      <div class="case-grid"><article v-for="item in cases" :key="item.name" class="case-card"><h3>{{ item.name }}</h3><p>{{ item.result }}</p></article></div>
    </section>

    <section id="pricing" class="section">
      <div class="section-heading"><p class="eyebrow">Plans</p><h2>Choose the right starting point</h2></div>
      <div class="pricing-grid"><article v-for="plan in pricing" :key="plan.name" class="price-card"><h3>{{ plan.name }}</h3><strong>{{ plan.price }}</strong><p>{{ plan.description }}</p><ul><li v-for="item in plan.items" :key="item">{{ item }}</li></ul></article></div>
    </section>

    <section id="contact" class="contact-section">
      <div><p class="eyebrow">Next step</p><h2>{{ contactTitle }}</h2><p>{{ contactSubtitle }}</p></div>
      <form class="lead-form"><label>Name<input placeholder="Your name" /></label><label>Email<input placeholder="you@example.com" /></label><label>Goal<textarea placeholder="Tell us what you want to build"></textarea></label><button type="button">{{ primaryCTA }}</button></form>
    </section>
  </main>
</template>

<style scoped>
:global(body){margin:0;font-family:Inter,ui-sans-serif,system-ui,-apple-system,BlinkMacSystemFont,"Segoe UI",sans-serif;background:#07111f;color:#ecf5ff}*{box-sizing:border-box}.page-shell{min-height:100vh;background:radial-gradient(circle at top left,rgba(56,189,248,.25),transparent 32rem),linear-gradient(135deg,#08111f,#10213d 52%%,#07111f)}.hero-section,.section,.contact-section{width:min(1120px,calc(100%% - 32px));margin:0 auto}.hero-section{padding:28px 0 72px}.nav-bar{display:flex;align-items:center;justify-content:space-between;padding:14px 18px;border:1px solid rgba(148,163,184,.24);border-radius:22px;background:rgba(15,23,42,.72);backdrop-filter:blur(18px)}.brand{display:flex;align-items:center;gap:10px;font-weight:800}.brand-mark{display:grid;place-items:center;width:34px;height:34px;border-radius:12px;background:#38bdf8;color:#06111f}.nav-actions{display:flex;gap:18px}.nav-actions a{color:#cbd5e1;text-decoration:none}.hero-grid{display:grid;grid-template-columns:1.25fr .75fr;gap:36px;align-items:center;padding-top:82px}.eyebrow{margin:0 0 12px;color:#67e8f9;text-transform:uppercase;letter-spacing:.14em;font-size:12px;font-weight:800}.hero-copy h1{margin:0;font-size:clamp(44px,7vw,78px);line-height:.95;letter-spacing:-.06em}.hero-subtitle{max-width:680px;color:#cbd5e1;font-size:20px;line-height:1.7}.cta-row{display:flex;gap:14px;flex-wrap:wrap;margin-top:28px}.btn,.lead-form button{border:0;border-radius:999px;padding:14px 22px;font-weight:800;text-decoration:none;cursor:pointer}.primary,.lead-form button{background:linear-gradient(135deg,#22d3ee,#60a5fa);color:#06111f}.ghost{border:1px solid rgba(203,213,225,.3);color:#e2e8f0}.hero-card,.feature-card,.case-card,.price-card,.contact-section{border:1px solid rgba(148,163,184,.22);box-shadow:0 24px 80px rgba(0,0,0,.28);background:rgba(15,23,42,.72);backdrop-filter:blur(18px)}.hero-card{border-radius:30px;padding:24px;display:grid;gap:16px}.metric{padding:18px;border-radius:22px;background:rgba(15,23,42,.78)}.metric strong{display:block;font-size:34px;color:#7dd3fc}.metric span{color:#cbd5e1}.section{padding:72px 0}.section-heading{max-width:680px;margin-bottom:28px}.section-heading h2,.contact-section h2{margin:0;font-size:clamp(30px,4vw,48px);letter-spacing:-.04em}.feature-grid,.case-grid,.pricing-grid{display:grid;grid-template-columns:repeat(3,1fr);gap:18px}.feature-card,.case-card,.price-card{border-radius:24px;padding:24px}.feature-card h3,.case-card h3,.price-card h3{margin:0 0 10px}.feature-card p,.case-card p,.price-card p,.contact-section p,.price-card li{color:#cbd5e1;line-height:1.7}.price-card strong{display:block;font-size:28px;color:#7dd3fc;margin:12px 0}.price-card ul{padding-left:18px}.contact-section{display:grid;grid-template-columns:.9fr 1.1fr;gap:28px;align-items:start;border-radius:32px;padding:34px;margin-bottom:56px}.lead-form{display:grid;gap:14px}.lead-form label{display:grid;gap:7px;color:#cbd5e1;font-weight:700}.lead-form input,.lead-form textarea{width:100%%;border:1px solid rgba(148,163,184,.26);border-radius:16px;background:rgba(2,6,23,.5);color:#fff;padding:13px 14px;font:inherit}.lead-form textarea{min-height:104px;resize:vertical}@media (max-width:900px){.hero-grid,.contact-section{grid-template-columns:1fr}.feature-grid,.case-grid,.pricing-grid{grid-template-columns:1fr}.nav-actions{display:none}.hero-section{padding-bottom:40px}.section{padding:44px 0}}
</style>
`,
		jsonString(spec.ProductName),
		jsonString(spec.Tagline),
		jsonString(spec.HeroTitle),
		jsonString(spec.HeroSubtitle),
		jsonString(spec.PrimaryCTA),
		jsonString(spec.SecondaryCTA),
		jsonString(spec.ContactTitle),
		jsonString(spec.ContactSubtitle),
		string(features),
		string(metrics),
		string(cases),
		string(pricing),
	)

	_ = project
	return []GeneratedProjectFile{
		{Path: "package.json", Content: packageJSON},
		{Path: "index.html", Content: indexHTML},
		{Path: "vite.config.ts", Content: viteConfig},
		{Path: "src/main.ts", Content: mainTS},
		{Path: "src/App.vue", Content: appVue},
	}, nil
}

func jsonString(value string) string {
	bytes, err := json.Marshal(value)
	if err != nil {
		return `""`
	}
	return string(bytes)
}

func htmlText(value string) string {
	return strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;", "'", "&#39;").Replace(value)
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

func (s *AIGenerationService) loadConfirmedBlueprintContent(projectID string) string {
	var bp model.Blueprint
	if err := s.db.Where("project_id = ? AND status = ?", projectID, "confirmed").Order("version desc").First(&bp).Error; err != nil {
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

func safePromptJSON(value string, maxBytes int) string {
	value = strings.TrimSpace(util.RedactText(value))
	if value == "" {
		return "{}"
	}
	return util.TruncateText(value, maxBytes)
}

func compactBlueprintForCodegen(value string) string {
	value = strings.TrimSpace(util.RedactText(value))
	if value == "" {
		return "{}"
	}
	var bp BlueprintAIResult
	if err := json.Unmarshal([]byte(value), &bp); err == nil {
		compact := BlueprintAIResult{
			SchemaVersion:       bp.SchemaVersion,
			AppType:             bp.AppType,
			ProductPositioning:  util.TruncateText(bp.ProductPositioning, 800),
			TechStack:           util.TruncateText(bp.TechStack, 300),
			UIStyle:             util.TruncateText(bp.UIStyle, 500),
			Pages:               bp.Pages,
			Components:          bp.Components,
			UserFlows:           bp.UserFlows,
			Features:            bp.Features,
			InteractionRules:    bp.InteractionRules,
			StateModel:          bp.StateModel,
			DomainRules:         bp.DomainRules,
			AcceptanceCriteria:  bp.AcceptanceCriteria,
			TestPlan:            bp.TestPlan,
			ImplementationNotes: bp.ImplementationNotes,
		}
		if len(compact.Pages) > 8 {
			compact.Pages = compact.Pages[:8]
		}
		if len(compact.Components) > 12 {
			compact.Components = compact.Components[:12]
		}
		if len(compact.ImplementationNotes) > 8 {
			compact.ImplementationNotes = compact.ImplementationNotes[:8]
		}
		if bytes, err := json.Marshal(compact); err == nil {
			return util.TruncateText(string(bytes), 6000)
		}
	}
	return util.TruncateText(value, 6000)
}

func normalizeBlueprint(result *BlueprintAIResult, project *model.Project) {
	result.SchemaVersion = 2
	if strings.TrimSpace(result.AppType) == "" {
		result.AppType = project.ProjectType
	}
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
	for i := range result.Features {
		if strings.TrimSpace(result.Features[i].ID) == "" {
			result.Features[i].ID = fmt.Sprintf("F-%03d", i+1)
		}
		if strings.TrimSpace(result.Features[i].Priority) == "" {
			result.Features[i].Priority = "must"
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
	if result.SchemaVersion >= 2 {
		if strings.TrimSpace(result.AppType) == "" {
			return fmt.Errorf("AI 返回的蓝图缺少 app_type")
		}
		if len(result.Features) == 0 {
			return fmt.Errorf("AI 返回的蓝图缺少 features")
		}
		for _, feature := range result.Features {
			if strings.EqualFold(feature.Priority, "must") && (strings.TrimSpace(feature.Name) == "" || len(feature.AcceptanceCriteria) == 0) {
				return fmt.Errorf("必须功能 %s 缺少名称或验收标准", feature.ID)
			}
		}
		if len(result.OpenQuestions) > 0 {
			return fmt.Errorf("蓝图仍有 %d 个未决问题", len(result.OpenQuestions))
		}
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
		"package.json":  false,
		"tsconfig.json": false,
		"index.html":    false,
		"src/main.ts":   false,
		"src/App.vue":   false,
	}
	seen := map[string]bool{}
	totalBytes := 0
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
		if len(file.Content) > 256*1024 {
			return fmt.Errorf("AI 返回的单个文件过大: %s", path)
		}
		totalBytes += len(file.Content)
		if totalBytes > 2*1024*1024 {
			return fmt.Errorf("AI 返回的项目总大小超过 2MB")
		}
		if _, ok := required[path]; ok {
			required[path] = true
		}
		if path == "package.json" {
			var packageFile struct {
				Scripts         map[string]string `json:"scripts"`
				Dependencies    map[string]string `json:"dependencies"`
				DevDependencies map[string]string `json:"devDependencies"`
			}
			if err := json.Unmarshal([]byte(file.Content), &packageFile); err != nil {
				return fmt.Errorf("AI 返回的 package.json 无效: %w", err)
			}
			if strings.TrimSpace(packageFile.Scripts["build"]) == "" {
				return fmt.Errorf("package.json 缺少 build 脚本")
			}
			for script := range packageFile.Scripts {
				switch script {
				case "dev", "build", "preview", "test":
				default:
					return fmt.Errorf("package.json 包含不允许的脚本: %s", script)
				}
			}
			for name, script := range packageFile.Scripts {
				if strings.ContainsAny(script, ";&|><`$") {
					return fmt.Errorf("package.json 脚本包含不安全的 shell 操作: %s", name)
				}
				fields := strings.Fields(script)
				if len(fields) == 0 || (fields[0] != "vite" && fields[0] != "vitest") {
					return fmt.Errorf("package.json 脚本命令不在允许列表: %s", name)
				}
			}
			allowedPackages := map[string]bool{
				"vue": true, "vue-router": true, "pinia": true,
				"vite": true, "typescript": true, "vue-tsc": true, "@vitejs/plugin-vue": true,
				"vitest": true, "@vue/test-utils": true, "jsdom": true,
			}
			for name := range packageFile.Dependencies {
				if !allowedPackages[name] {
					return fmt.Errorf("package.json 包含不允许的依赖: %s", name)
				}
			}
			for name := range packageFile.DevDependencies {
				if !allowedPackages[name] {
					return fmt.Errorf("package.json 包含不允许的开发依赖: %s", name)
				}
			}
			if packageFile.Dependencies["vue-tsc"] == "" && packageFile.DevDependencies["vue-tsc"] == "" {
				return fmt.Errorf("package.json 缺少用于类型检查的 vue-tsc")
			}
		}
	}
	for path, ok := range required {
		if !ok {
			return fmt.Errorf("AI 返回的文件缺少必需文件: %s", path)
		}
	}
	return nil
}

// ValidateBlueprintContent validates a stored blueprint before confirmation.
func ValidateBlueprintContent(content string) error {
	var result BlueprintAIResult
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return fmt.Errorf("蓝图 JSON 无效: %w", err)
	}
	return validateBlueprint(result)
}

// ValidateBlueprintAgainstRequirements applies deterministic confirmation
// gates using both versioned inputs.
func ValidateBlueprintAgainstRequirements(blueprintContent, requirementContent string) error {
	var blueprint BlueprintAIResult
	if err := json.Unmarshal([]byte(blueprintContent), &blueprint); err != nil {
		return fmt.Errorf("蓝图 JSON 无效: %w", err)
	}
	if err := validateBlueprint(blueprint); err != nil {
		return err
	}
	var requirement struct {
		SchemaVersion    int      `json:"schema_version"`
		MustHaveFeatures []string `json:"must_have_features"`
	}
	if err := json.Unmarshal([]byte(requirementContent), &requirement); err != nil {
		return fmt.Errorf("需求 JSON 无效: %w", err)
	}
	if requirement.SchemaVersion < 2 {
		return nil
	}
	mustCount := 0
	for _, feature := range blueprint.Features {
		if strings.EqualFold(feature.Priority, "must") {
			mustCount++
		}
	}
	if mustCount < len(requirement.MustHaveFeatures) {
		return fmt.Errorf("必须功能覆盖不足：需求 %d 项，蓝图仅 %d 项", len(requirement.MustHaveFeatures), mustCount)
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
