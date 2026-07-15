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
	requirements := safePromptJSON(s.loadLatestRequirementContent(projectID), 3500)
	blueprint := s.loadLatestBlueprintContent(projectID)
	if strings.TrimSpace(blueprint) == "" {
		return nil, fmt.Errorf("blueprint is missing; generate and confirm blueprint first")
	}
	blueprint = compactBlueprintForCodegen(blueprint)

	systemPrompt := `You are a senior product copywriter and frontend content architect. Return a valid JSON object only. Do not use Markdown or explanations.`
	userPrompt := fmt.Sprintf(`Create compact page content for a runnable Vue 3 landing page based on the requirement and blueprint.

Project info:
- Name: %s
- Description: %s
- UI style: %s
- Coding rules: %s

Requirement summary JSON:
%s

Blueprint summary JSON:
%s

Return JSON exactly in this shape:
{
  "product_name": "product or project name",
  "tagline": "short tagline",
  "hero_title": "conversion focused hero title",
  "hero_subtitle": "one paragraph hero subtitle",
  "primary_cta": "primary button text",
  "secondary_cta": "secondary button text",
  "metrics": [{"value":"98%%","label":"metric label"}],
  "features": [{"title":"feature title","description":"feature description"}],
  "cases": [{"name":"customer/case name","result":"case result"}],
  "pricing": [{"name":"plan name","price":"price text","description":"plan description","items":["item"]}],
  "contact_title": "lead form title",
  "contact_subtitle": "lead form subtitle",
  "notes": ["implementation note"]
}

Hard requirements:
- Keep the response compact: 3 metrics, 4 to 6 features, 2 to 3 cases, and 3 pricing plans.
- Use the same language as the requirement/blueprint when possible.
- The content must be specific to this project; do not return generic placeholders.`,
		project.Name,
		util.RedactText(project.Description),
		project.UIStyle,
		util.RedactText(project.CodingRules),
		requirements,
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

	var spec CodeGenerationContentSpec
	if err := util.ParseAIResponseInto(aiResult.JSONText, &spec); err != nil {
		return nil, fmt.Errorf("parse code generation AI response: %w", err)
	}
	normalizeContentSpec(&spec, project)
	files, err := renderVueProjectFiles(project, spec)
	if err != nil {
		return nil, err
	}
	if err := validateGeneratedFiles(files); err != nil {
		return nil, err
	}
	return &CodeGenerationAIResult{Files: files, Notes: spec.Notes}, nil
}

func normalizeContentSpec(spec *CodeGenerationContentSpec, project *model.Project) {
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
			ProductPositioning:  util.TruncateText(bp.ProductPositioning, 800),
			TechStack:           util.TruncateText(bp.TechStack, 300),
			UIStyle:             util.TruncateText(bp.UIStyle, 500),
			Pages:               bp.Pages,
			Components:          bp.Components,
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
