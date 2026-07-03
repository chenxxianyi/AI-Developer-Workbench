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

type mockIssueItem struct {
	Title      string `json:"title"`
	Severity   string `json:"severity"`
	Category   string `json:"category"`
	Problem    string `json:"problem"`
	Suggestion string `json:"suggestion"`
	Action     string `json:"action"`
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
		"scores": []mockScoreItem{
			{Name: "视觉层级", Score: 78, MaxScore: 100, Comment: "整体层级清晰，但部分区域间距不一致。"},
			{Name: "一致性", Score: 82, MaxScore: 100, Comment: "按钮和表单风格统一，但图标风格有轻微差异。"},
			{Name: "可访问性", Score: 65, MaxScore: 100, Comment: "部分元素缺少 aria-label，键盘导航不完整。"},
			{Name: "响应式", Score: 72, MaxScore: 100, Comment: "桌面端表现良好，移动端导航栏溢出。"},
		},
		"issues": []mockIssueItem{
			{
				Title: "移动端导航栏横向溢出", Severity: "high", Category: "responsive",
				Problem: "在 320px 宽度下导航栏出现横向滚动条。",
				Suggestion: "将导航链接改为汉堡菜单或使用 flex-wrap。",
				Action: "检查并修复导航栏在移动端的布局。",
			},
			{
				Title: "上传按钮缺少键盘支持", Severity: "high", Category: "accessibility",
				Problem: "文件上传区域仅响应点击事件，键盘用户无法使用。",
				Suggestion: "添加 Enter/Space 事件处理和 aria-label。",
				Action: "为上传区域添加键盘事件和 ARIA 属性。",
			},
			{
				Title: "评分文字颜色对比度不足", Severity: "medium", Category: "contrast",
				Problem: "浅灰色文字在白色背景上对比度仅 2.3:1。",
				Suggestion: "将文字颜色改为 #555 或更深的灰色。",
				Action: "调整评分文字颜色以符合 WCAG AA 标准。",
			},
			{
				Title: "表单字段缺少 label 关联", Severity: "medium", Category: "form",
				Problem: "部分输入框使用 placeholder 代替 label，屏幕阅读器无法识别。",
				Suggestion: "为每个输入框添加显式 label 元素。",
				Action: "为所有表单字段添加 label 和 htmlFor 关联。",
			},
		},
		"recommendations": []string{
			"统一图标库，避免混用 Material Icons 和自定义 SVG。",
			"为移动端导航实现汉堡菜单。",
			"补充所有交互元素的 aria-label。",
		},
		"codex_prompt": "请修复 UIReviewPage 中的可访问性问题：1) 上传区域增加键盘 Enter/Space 触发；2) 所有表单字段添加显式 label；3) 移动端导航栏改为汉堡菜单。",
	}
}

func (s *MockAIService) mockProjectDoctorData() interface{} {
	return map[string]interface{}{
		"scores": []mockScoreItem{
			{Name: "结构清晰度", Score: 75, MaxScore: 100, Comment: "项目目录结构合理，但缺少明确的模块划分文档。"},
			{Name: "可维护性", Score: 68, MaxScore: 100, Comment: "存在硬编码配置和重复代码。"},
			{Name: "可测试性", Score: 55, MaxScore: 100, Comment: "测试覆盖率较低，缺少集成测试。"},
			{Name: "可部署性", Score: 80, MaxScore: 100, Comment: "Dockerfile 存在，但缺少健康检查和 docker-compose。"},
			{Name: "文档完整度", Score: 45, MaxScore: 100, Comment: "缺少 README 和 API 文档。"},
			{Name: "Agent 可接手程度", Score: 60, MaxScore: 100, Comment: "项目缺少 AGENTS.md 和 AI Coding 配置。"},
		},
		"issues": []mockIssueItem{
			{
				Title: "缺少 README 文档", Severity: "high", Category: "documentation",
				Problem: "项目中没有 README.md，新人无法快速了解项目。",
				Suggestion: "创建 README.md，包含项目简介、安装步骤和开发指南。",
				Action: "编写完整的 README.md。",
			},
			{
				Title: "测试覆盖率不足", Severity: "high", Category: "testing",
				Problem: "仅有一个单元测试文件，核心业务逻辑未覆盖。",
				Suggestion: "为关键模块添加单元测试和集成测试。",
				Action: "补充测试用例，目标行覆盖率 > 60%。",
			},
			{
				Title: "硬编码配置", Severity: "medium", Category: "config",
				Problem: "数据库连接和 API Key 等配置硬编码在代码中。",
				Suggestion: "使用环境变量或配置文件管理所有配置。",
				Action: "将所有硬编码配置迁移到 .env 文件。",
			},
			{
				Title: "缺少 Docker Compose", Severity: "medium", Category: "deploy",
				Problem: "仅有 Dockerfile，没有 docker-compose.yml。",
				Suggestion: "创建 docker-compose.yml 编排前后端和数据库。",
				Action: "编写 docker-compose.yml 并验证一键启动。",
			},
		},
		"recommendations": []string{
			"优先补充 README.md 和 AGENTS.md。",
			"提取所有硬编码配置到环境变量。",
			"为关键路径添加测试。",
		},
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
		"recommendations": []string{
			"AGENTS.md 已生成，建议根据实际项目调整命令和路径。",
			"CODING_RULES.md 包含前后端通用规则，可根据团队偏好修改。",
			"TASK_PLAN.md 列出了当前待办事项，建议纳入项目管理工具。",
		},
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
				Problem: "按 tool_type 和 status 查询时需要全表扫描。",
				Suggestion: "在 (tool_type, status, created_at) 上创建复合索引。",
				Action: "添加 idx_reports_tool_status_created 索引。",
			},
			{
				Title: "字符串主键性能问题", Severity: "medium", Category: "structure",
				Problem: "使用 UUID 字符串作为主键在 InnoDB 中可能导致页分裂。",
				Suggestion: "考虑使用自增 ID 作为主键，UUID 作为业务标识。",
				Action: "评估迁移成本后决定是否调整主键策略。",
			},
			{
				Title: "缺少 updated_at 自动更新", Severity: "medium", Category: "integrity",
				Problem: "部分表没有自动更新 updated_at 的触发器。",
				Suggestion: "使用 GORM 的 BeforeUpdate hook 或数据库触发器。",
				Action: "确保所有表在更新记录时自动更新 updated_at。",
			},
		},
		"optimized_schema": strPtr("-- 建议优化\nALTER TABLE reports ADD INDEX idx_tool_status_created (tool_type, status, created_at);\nALTER TABLE reports MODIFY COLUMN title VARCHAR(255) NOT NULL;\n-- ⚠ 以下 DDL 有锁表风险，请在低峰期执行\n-- ALTER TABLE reports ADD FOREIGN KEY ..."),
		"migration_suggestions": []string{
			"阶段 1 (安全): 添加索引 ALTER TABLE reports ADD INDEX idx_tool_status_created (tool_type, status, created_at);",
			"阶段 2 (低峰): 修改列约束 ALTER TABLE reports MODIFY COLUMN title VARCHAR(255) NOT NULL;",
			"阶段 3 (评估): 评估 UUID 主键改为自增 ID 的迁移方案",
		},
		"recommendations": []string{
			"优先添加高频查询列的索引。",
			"审查所有 VARCHAR 列的长度是否合理。",
			"考虑定期清理历史数据或分区策略。",
		},
		"codex_prompt": "请为 reports 表添加复合索引 idx_tool_status_created (tool_type, status, created_at)，并确保所有 VARCHAR 列有合理的长度约束。",
	}
}

func strPtr(s string) *string {
	return &s
}

// Ensure MockAIService satisfies AIService interface.
var _ AIService = (*MockAIService)(nil)

// Avoid unused import error for time in dbschema.
var _ = time.Now
