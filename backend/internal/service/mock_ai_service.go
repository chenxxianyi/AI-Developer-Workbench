package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"ai-developer-workbench/internal/model"
)

// MockAIService returns deterministic JSON for each tool type without any external calls.
type MockAIService struct{}

// NewMockAIService creates a new MockAIService.
func NewMockAIService() *MockAIService {
	return &MockAIService{}
}

// GenerateJSON returns a deterministic mock result based on ToolType.
func (s *MockAIService) GenerateJSON(_ context.Context, input AIRequest) (*AIResult, error) {
	var data interface{}

	switch input.ToolType {
	case model.ToolTypeUIReview:
		data = s.mockUIReviewData()
	case model.ToolTypeProjectDoctor:
		data = s.mockProjectDoctorData()
	case model.ToolTypeAgentConfig:
		data = s.mockAgentConfigData()
	case model.ToolTypeAPIDoc:
		data = s.mockAPIDocData()
	case model.ToolTypeDBSchema:
		data = s.mockDBSchemaData()
	default:
		return nil, fmt.Errorf("mock: unknown tool type %q", input.ToolType)
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("mock: failed to marshal JSON: %w", err)
	}
	jsonText := string(jsonBytes)

	return &AIResult{
		RawText:  jsonText,
		JSONText: jsonText,
		Provider: "mock",
		Model:    "mock-mode",
	}, nil
}

// --- Mock data types ---

type mockScoreItem struct {
	Name     string `json:"name"`
	Score    int    `json:"score"`
	MaxScore int    `json:"max_score"`
	Comment  string `json:"comment"`
}

type mockIssueRegion struct{ X, Y, Width, Height float64 }
type mockIssueItem struct {
	Title              string           `json:"title"`
	Severity           string           `json:"severity"`
	Category           string           `json:"category"`
	Problem            string           `json:"problem"`
	Suggestion         string           `json:"suggestion"`
	Action             string           `json:"action"`
	Viewport           string           `json:"viewport,omitempty"`
	Region             *mockIssueRegion `json:"region,omitempty"`
	ContrastSuggestion string           `json:"contrast_suggestion,omitempty"`
	ComponentPrompt    string           `json:"component_prompt,omitempty"`
}

type mockActionItem struct {
	ID              string `json:"id"`
	Title           string `json:"title"`
	Priority        string `json:"priority"`
	Effort          string `json:"effort"`
	Category        string `json:"category"`
	Reason          string `json:"reason"`
	SuggestedPrompt string `json:"suggested_prompt"`
	IssueTitle      string `json:"issue_title"`
	IssueBody       string `json:"issue_body"`
}

type mockEndpointItem struct {
	Method      string `json:"method"`
	Path        string `json:"path"`
	Description string `json:"description"`
}

type mockModuleItem struct {
	Name      string             `json:"name"`
	Endpoints []mockEndpointItem `json:"endpoints"`
}

// --- Tool-specific mock data ---

func (s *MockAIService) mockUIReviewData() interface{} {
	return map[string]interface{}{
		"screenshot_contexts": []map[string]string{{"kind": "desktop", "viewport": "1440x900"}, {"kind": "mobile", "viewport": "390x844"}},
		"scores": []mockScoreItem{
			{Name: "视觉层级", Score: 78, MaxScore: 100, Comment: "整体层级清晰，但部分区域间距不一致。"},
			{Name: "一致性", Score: 82, MaxScore: 100, Comment: "按钮和表单风格统一，但图标风格有轻微差异。"},
			{Name: "可访问性", Score: 65, MaxScore: 100, Comment: "部分元素缺少 aria-label，键盘导航不完整。"},
			{Name: "响应式", Score: 72, MaxScore: 100, Comment: "桌面端表现良好，移动端导航栏溢出。"},
		},
		"issues": []mockIssueItem{
			{
				Title: "移动端导航栏横向溢出", Severity: "high", Category: "responsive",
				Problem:    "在 320px 宽度下导航栏出现横向滚动条。",
				Suggestion: "将导航链接改为汉堡菜单或使用 flex-wrap。",
				Action:     "检查并修复导航栏在移动端的布局。",
				Viewport:   "mobile", Region: &mockIssueRegion{X: 5, Y: 2, Width: 90, Height: 12},
				ComponentPrompt: "Update the navigation component to collapse below 640px and preserve keyboard focus.",
			},
			{
				Title: "上传按钮缺少键盘支持", Severity: "high", Category: "accessibility",
				Problem:    "文件上传区域仅响应点击事件，键盘用户无法使用。",
				Suggestion: "添加 Enter/Space 事件处理和 aria-label。",
				Action:     "为上传区域添加键盘事件和 ARIA 属性。",
			},
			{
				Title: "评分文字颜色对比度不足", Severity: "medium", Category: "contrast",
				Problem:    "浅灰色文字在白色背景上对比度仅 2.3:1。",
				Suggestion: "将文字颜色改为 #555 或更深的灰色。",
				Action:     "调整评分文字颜色以符合 WCAG AA 标准。",
			},
			{
				Title: "表单字段缺少 label 关联", Severity: "medium", Category: "form",
				Problem:    "部分输入框使用 placeholder 代替 label，屏幕阅读器无法识别。",
				Suggestion: "为每个输入框添加显式 label 元素。",
				Action:     "为所有表单字段添加 label 和 htmlFor 关联。",
			},
		},
		"recommendations": []string{
			"统一图标库，避免混用 Material Icons 和自定义 SVG。",
			"为移动端导航实现汉堡菜单。",
			"补充所有交互元素的 aria-label。",
		},
		"action_items": mockActionItems("ui-review", "accessibility", "UI 审查"),
		"codex_prompt": "请修复 UIReviewPage 中的可访问性问题：1) 上传区域增加键盘 Enter/Space 触发；2) 所有表单字段添加显式 label；3) 移动端导航栏改为汉堡菜单。",
	}
}

func (s *MockAIService) mockProjectDoctorData() interface{} {
	return map[string]interface{}{
		"scores": []mockScoreItem{
			{Name: "结构清晰度", Score: 75, MaxScore: 100, Comment: "项目目录结构合理，但缺少明确的模块划分文档。"},
			{Name: "可维护性", Score: 68, MaxScore: 100, Comment: "存在硬编码配置和重复代码。证据：多处直接使用环境变量字符串。"},
			{Name: "可测试性", Score: 55, MaxScore: 100, Comment: "测试覆盖率较低，缺少集成测试。证据：仅有 1 个测试文件。"},
			{Name: "可部署性", Score: 80, MaxScore: 100, Comment: "Dockerfile 存在，但缺少健康检查和 docker-compose。"},
			{Name: "文档完整度", Score: 45, MaxScore: 100, Comment: "缺少 README，API 文档不完整。"},
			{Name: "Agent 可接手程度", Score: 60, MaxScore: 100, Comment: "项目缺少 AGENTS.md，AI Coding 工具难以快速上手。"},
		},
		"evidence_files": []map[string]interface{}{
			{"path": "README.md", "type": "readme", "present": false, "notes": "根目录缺少 README，建议补充安装和开发指南"},
			{"path": "AGENTS.md", "type": "agents_md", "present": false, "notes": "需创建以帮助 AI Coding 工具理解项目"},
			{"path": "backend/go.sum", "type": "lockfile", "present": true, "notes": "Go 依赖锁定文件已存在"},
			{"path": "frontend/package-lock.json", "type": "lockfile", "present": true, "notes": "前端依赖已锁定"},
			{"path": "backend/Dockerfile", "type": "dockerfile", "present": true, "notes": "后端已有 Dockerfile，前端缺失"},
			{"path": ".github/workflows", "type": "ci", "present": false, "notes": "无 CI 配置，建议添加 GitHub Actions"},
			{"path": ".env.example", "type": "docs", "present": true, "notes": "环境变量模板已存在"},
		},
		"issues": []mockIssueItem{
			{
				Title: "缺少 README 文档", Severity: "high", Category: "documentation",
				Problem:    "项目中没有 README.md，新人无法快速了解项目。",
				Suggestion: "创建 README.md，包含项目简介、安装步骤和开发指南。",
				Action:     "编写完整的 README.md，包含启动命令和架构说明。",
			},
			{
				Title: "测试覆盖率不足", Severity: "high", Category: "testing",
				Problem:    "仅有一个单元测试文件，核心业务逻辑未覆盖。",
				Suggestion: "为关键模块添加单元测试和集成测试。",
				Action:     "补充测试用例，目标行覆盖率 > 60%。",
			},
			{
				Title: "硬编码配置", Severity: "medium", Category: "config",
				Problem:    "数据库连接和 API Key 等配置硬编码在代码中。",
				Suggestion: "使用环境变量或配置文件管理所有配置。",
				Action:     "将所有硬编码配置迁移到 .env 文件。",
			},
			{
				Title: "缺少 CI/CD 流程", Severity: "medium", Category: "deploy",
				Problem:    "项目没有 .github/workflows 目录，缺少自动化构建和测试。",
				Suggestion: "添加 GitHub Actions 或等效 CI 配置。",
				Action:     "创建包含 build + test 的 CI 工作流。",
			},
		},
		"tech_debt": []map[string]interface{}{
			{"title": "补全前后端单元测试", "impact": "high", "cost": "high", "category": "testing",
				"description": "核心业务逻辑未覆盖，回归风险高。", "suggested_fix": "优先为 service 层添加单元测试，使用 mock 隔离外部依赖。"},
			{"title": "消除硬编码配置", "impact": "high", "cost": "medium", "category": "config",
				"description": "多处数据库和 API 配置硬编码，不同环境切换困难。", "suggested_fix": "统一使用 config.LoadConfig() 读取环境变量，删除所有硬编码。"},
			{"title": "补充项目文档", "impact": "medium", "cost": "low", "category": "documentation",
				"description": "缺少 README 和 AGENTS.md，协作效率低。", "suggested_fix": "优先创建 README.md 和 AGENTS.md，后续补充 API 文档。"},
			{"title": "前端组件单元测试", "impact": "medium", "cost": "medium", "category": "testing",
				"description": "Vue 组件缺少 vitest 测试，UI 回归靠手动。", "suggested_fix": "使用 @vue/test-utils 为关键页面组件添加渲染和行为测试。"},
		},
		"recommendations": []string{
			"优先补充 README.md 和 AGENTS.md 以提升 Agent Readiness。",
			"提取所有硬编码配置到环境变量，避免多环境部署风险。",
			"为 service 层添加单元测试，覆盖核心业务流程。",
		},
		"action_items": mockActionItems("project-doctor", "project-health", "项目诊断"),
		"codex_prompt": "请修复项目的以下问题：1) 创建 README.md 包含项目介绍、安装和开发指南；2) 创建 AGENTS.md 提供 AI Coding 工具使用说明；3) 将所有硬编码配置迁移到环境变量。",
	}
}

func (s *MockAIService) mockAgentConfigData() interface{} {
	return map[string]interface{}{
		"generated_files_content": map[string]string{
			"AGENTS.md": `# AGENTS.md

## Project Overview
This is a full-stack web application built with Vue 3 and Go.

## Tech Stack
- Frontend: Vue 3, TypeScript, Pinia, Vite
- Backend: Go, Gin, GORM
- Database: MySQL 8

## Commands
- Start backend: cd backend && go run ./cmd/server
- Start frontend: cd frontend && npm run dev
- Run tests: cd backend && go test ./...

## Project Structure
backend/   - Go API server
  cmd/server/    - Entry point
  internal/      - Application code
frontend/  - Vue 3 SPA
  src/pages/     - Page components
  src/components/- Shared components

## Coding Conventions
- Use TypeScript strict mode
- All Go functions must have comments
- Follow Vue 3 Composition API patterns
`,
			"TASK_PLAN.md": `# Task Plan

## Current Sprint
- [ ] Implement report list page
- [ ] Implement report detail page
- [ ] Add mock mode for demos

## Backlog
- [ ] Add project profiles
- [ ] Implement async job processing
- [ ] Add GitHub integration

## Bug Fixes
- [ ] Fix mobile navigation overflow
- [ ] Fix keyboard accessibility on upload area
`,
			"CODING_RULES.md": `# Coding Rules

## General
- No hardcoded secrets or API keys
- Use environment variables for all configuration
- Keep functions small (< 50 lines)

## Frontend
- Use composables for reusable logic
- All user-visible strings get i18n support
- ARIA labels on all interactive elements

## Backend
- All handlers validate input before processing
- Repository layer handles all DB access
- Use context-aware operations for cancellation

## Testing
- Unit tests for all business logic
- Integration tests for API endpoints
- E2E tests for critical user flows
`,
		},
		"target_format": "codex",
		"missing_confirmations": []string{
			"[CONFIRM] FRONTEND_STYLE_GUIDE.md has not been generated — provide frontend framework details for accurate style guide.",
		},
		"recommendations": []string{
			"AGENTS.md 已生成，建议根据实际项目调整命令和路径。",
			"CODING_RULES.md 包含前后端通用规则，可根据团队偏好修改。",
			"TASK_PLAN.md 列出了当前待办事项，建议纳入项目管理工具。",
		},
		"action_items": mockActionItems("agent-config", "agent-config", "Agent 配置"),
		"codex_prompt": "请根据 CODING_RULES.md 中的规则审查当前代码库，并修复所有违规项。",
	}
}

func (s *MockAIService) mockAPIDocData() interface{} {
	mdContent := `# API Documentation

## Health

### GET /api/health
返回服务健康状态。

**Response 200:**
{"status": "ok"}

## Reports

### GET /api/reports
获取报告列表。

**Query Parameters:**
- tool_type (string, optional): 按工具类型筛选
- sort (string, optional): 排序方式 (newest/oldest/score_desc/score_asc)
- page (int, optional): 页码 (default: 1)
- page_size (int, optional): 每页数量 (default: 10, max: 100)
`

	return map[string]interface{}{
		"modules": []mockModuleItem{
			{
				Name: "Health",
				Endpoints: []mockEndpointItem{
					{Method: "GET", Path: "/api/health", Description: "返回服务健康状态"},
				},
			},
			{
				Name: "System",
				Endpoints: []mockEndpointItem{
					{Method: "GET", Path: "/api/system/status", Description: "返回系统状态和 AI 配置信息"},
				},
			},
			{
				Name: "Reports",
				Endpoints: []mockEndpointItem{
					{Method: "GET", Path: "/api/reports", Description: "获取报告列表，支持筛选和分页"},
					{Method: "GET", Path: "/api/reports/:id", Description: "获取报告详情"},
					{Method: "DELETE", Path: "/api/reports/:id", Description: "删除报告及其关联文件"},
				},
			},
			{
				Name: "Tools",
				Endpoints: []mockEndpointItem{
					{Method: "POST", Path: "/api/tools/ui-review/run", Description: "运行 UI 质量审查"},
					{Method: "POST", Path: "/api/tools/project-doctor/run", Description: "运行项目诊断"},
					{Method: "POST", Path: "/api/tools/agent-config/run", Description: "运行 Agent 配置生成"},
					{Method: "POST", Path: "/api/tools/api-doc/run", Description: "运行 API 文档生成"},
					{Method: "POST", Path: "/api/tools/db-schema/run", Description: "运行数据库结构审查"},
				},
			},
		},
		"markdown_content": &mdContent,
		"recommendations": []string{
			"建议为每个端点补充请求体和响应体的完整示例。",
			"补充认证方式的说明文档。",
			"为前端开发者提供调用示例。",
		},
		"action_items": mockActionItems("api-doc", "api-docs", "API 文档"),
		"codex_prompt": "请根据 API 文档补全 reports 模块的前端 API 客户端调用函数。",
	}
}

func (s *MockAIService) mockDBSchemaData() interface{} {
	return map[string]interface{}{
		"scores": []mockScoreItem{
			{Name: "结构评分", Score: 75, MaxScore: 100, Comment: "表设计合理，但缺少外键约束。"},
			{Name: "索引评分", Score: 60, MaxScore: 100, Comment: "核心表缺少常用查询的索引。"},
			{Name: "扩展性评分", Score: 70, MaxScore: 100, Comment: "单表数据量增长后可能影响查询性能。"},
			{Name: "数据完整性评分", Score: 65, MaxScore: 100, Comment: "缺少部分 NOT NULL 约束和数据校验。"},
		},
		"issues": []mockIssueItem{
			{
				Title: "reports 表缺少 report_data 索引", Severity: "high", Category: "index",
				Problem:    "按 tool_type 和 status 查询时需要全表扫描。",
				Suggestion: "在 (tool_type, status, created_at) 上创建复合索引。",
				Action:     "添加 idx_reports_tool_status_created 索引。",
			},
			{
				Title: "字符串主键性能问题", Severity: "medium", Category: "structure",
				Problem:    "使用 UUID 字符串作为主键在 InnoDB 中可能导致页分裂。",
				Suggestion: "考虑使用自增 ID 作为主键，UUID 作为业务标识。",
				Action:     "评估迁移成本后决定是否调整主键策略。",
			},
			{
				Title: "缺少 updated_at 自动更新", Severity: "medium", Category: "integrity",
				Problem:    "部分表没有自动更新 updated_at 的触发器。",
				Suggestion: "使用 GORM 的 BeforeUpdate hook 或数据库触发器。",
				Action:     "确保所有表在更新记录时自动更新 updated_at。",
			},
		},
		"optimized_schema": strPtr("-- 建议优化\nALTER TABLE reports ADD INDEX idx_tool_status_created (tool_type, status, created_at);\nALTER TABLE reports MODIFY COLUMN title VARCHAR(255) NOT NULL;\n-- ⚠ 以下 DDL 有锁表风险，请在低峰期执行\n-- ALTER TABLE reports ADD FOREIGN KEY ..."),
		"migration_suggestions": []string{
			"阶段 1 (安全): 添加索引 ALTER TABLE reports ADD INDEX idx_tool_status_created (tool_type, status, created_at);",
			"阶段 2 (低峰): 修改列约束 ALTER TABLE reports MODIFY COLUMN title VARCHAR(255) NOT NULL;",
			"阶段 3 (评估): 评估 UUID 主键改为自增 ID 的迁移方案",
		},
		"er_diagram_mermaid": strPtr("erDiagram\n  reports ||--o{ generated_files : contains\n  reports {\n    char id PK\n    varchar tool_type\n    varchar title\n    varchar status\n    int total_score\n  }\n  generated_files {\n    char id PK\n    char report_id FK\n    varchar filename\n    text content\n  }"),
		"index_recommendations": []map[string]interface{}{
			{"table": "reports", "columns": "(tool_type, status, created_at)", "reason": "报告列表按工具和状态筛选时避免全表扫描", "impact": "high"},
			{"table": "generated_files", "columns": "(report_id, filename)", "reason": "按报告 ID 和文件名检索时加速查询", "impact": "medium"},
		},
		"migration_risks": []map[string]interface{}{
			{"operation": "ALTER TABLE reports MODIFY COLUMN title VARCHAR(255) NOT NULL", "risk": "medium", "description": "如果存在 NULL 值会导致语句失败，执行前需检查并填充默认值", "rollback_plan": "ALTER TABLE reports MODIFY COLUMN title VARCHAR(255);"},
			{"operation": "ALTER TABLE reports ADD FOREIGN KEY (project_id) REFERENCES projects(id)", "risk": "high", "description": "外键约束可能阻塞现有报表写入，需在业务低峰期添加", "rollback_plan": "ALTER TABLE reports DROP FOREIGN KEY fk_reports_project;"},
		},
		"recommendations": []string{
			"优先添加高频查询列的索引。",
			"审查所有 VARCHAR 列的长度是否合理。",
			"考虑定期清理历史数据或分区策略。",
		},
		"action_items": mockActionItems("db-schema", "database", "数据库结构"),
		"codex_prompt": "请为 reports 表添加复合索引 idx_tool_status_created (tool_type, status, created_at)，并确保所有 VARCHAR 列有合理的长度约束。",
	}
}

func mockActionItems(prefix, category, target string) []mockActionItem {
	return []mockActionItem{
		{
			ID:              prefix + "-fix-high-priority",
			Title:           "修复" + target + "中的高优问题",
			Priority:        "high",
			Effort:          "small",
			Category:        category,
			Reason:          "报告已识别出会影响可用性或交付质量的高优问题。",
			SuggestedPrompt: "请根据报告中的 high severity 问题逐项修复，并在完成后列出修改文件和验证命令。",
			IssueTitle:      "fix(" + category + "): resolve high priority findings",
			IssueBody:       "## 背景\n报告发现高优问题需要优先处理。\n\n## 修复要求\n- 定位报告中的 high severity 项\n- 修改对应文件或配置\n- 保留兼容行为\n\n## 验收\n- [ ] 相关测试通过\n- [ ] 报告中的高优问题已复查",
		},
		{
			ID:              prefix + "-add-regression-tests",
			Title:           "补充" + target + "回归测试",
			Priority:        "medium",
			Effort:          "medium",
			Category:        "testing",
			Reason:          "当前问题需要自动化测试防止后续回归。",
			SuggestedPrompt: "请为本报告中的关键修复补充单元测试或组件测试，并说明覆盖的失败路径。",
			IssueTitle:      "test(" + category + "): add regression coverage",
			IssueBody:       "## 背景\n报告建议补充自动化测试。\n\n## 修复要求\n- 为关键行为添加测试\n- 覆盖成功和失败状态\n\n## 验收\n- [ ] 新增测试先失败后通过\n- [ ] 现有测试全部通过",
		},
		{
			ID:              prefix + "-document-follow-up",
			Title:           "记录" + target + "修复说明",
			Priority:        "low",
			Effort:          "small",
			Category:        "documentation",
			Reason:          "修复后的行为和验证命令需要沉淀，便于后续复查。",
			SuggestedPrompt: "请把本次修复的影响范围、验证命令和剩余风险整理到对应文档或 PR 描述中。",
			IssueTitle:      "docs(" + category + "): record validation notes",
			IssueBody:       "## 背景\n需要记录修复结果和验证方式。\n\n## 修复要求\n- 写明修改范围\n- 写明验证命令\n- 写明剩余风险\n\n## 验收\n- [ ] 文档包含可复查的命令和结果",
		},
	}
}

func strPtr(s string) *string {
	return &s
}

// Ensure MockAIService satisfies AIService interface.
var _ AIService = (*MockAIService)(nil)

// Avoid unused import error for time in dbschema.
var _ = time.Now
