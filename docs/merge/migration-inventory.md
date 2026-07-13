# M0 功能盘点和迁移映射表

> 基于 2026-07-13 代码基线分析

## 1. 前端页面

### Workbench（主仓库）— `frontend/src/`

| 页面 | 文件 | 处置 |
|------|------|------|
| 着陆页 | `pages/LandingPage.vue` | 保留（改为登录后重定向到 Dashboard） |
| Dashboard | `pages/DashboardPage.vue` | **合并**（扩展项目/任务/生成统计） |
| 项目列表 | `pages/ProjectsPage.vue` | **合并**（统一 Builder 项目卡片） |
| 项目创建/编辑 | `pages/ProjectFormPage.vue` | **合并**（增加生成项目类型） |
| 项目详情 | `pages/ProjectDetailPage.vue` | **合并**（增加子页面导航） |
| 设置 | `pages/SettingsPage.vue` | 保留 |
| UI Review | `pages/tools/UIReviewPage.vue` | 保留（增加项目上下文绑定） |
| Project Doctor | `pages/tools/ProjectDoctorPage.vue` | 保留（增加项目上下文绑定） |
| Agent Config | `pages/tools/AgentConfigPage.vue` | 保留 |
| API Doc | `pages/tools/APIDocPage.vue` | 保留 |
| DB Schema | `pages/tools/DBSchemaPage.vue` | 保留 |
| 报告列表 | `pages/ReportsPage.vue` | 保留（关联项目筛选） |
| 报告详情 | `pages/ReportDetailPage.vue` | 保留 |
| 报告对比 | `pages/ReportComparePage.vue` | 保留 |

### Builder — `apps/web/src/`

| 页面 | 文件 | 处置 |
|------|------|------|
| 登录 | `views/studio/LoginView.vue` | **迁移** |
| 首页/Dashboard | `views/studio/HomeView.vue` | 废弃（合并到 Workbench Dashboard） |
| 项目列表 | `views/studio/ProjectsView.vue` | 废弃（合并到 Workbench 项目列表） |
| 创建项目 | `views/studio/CreateProjectView.vue` | **迁移**（增加项目类型选择） |
| 蓝图评审 | `views/studio/BlueprintReviewView.vue` | **迁移** |
| 生成进度 | `views/studio/GenerationView.vue` | **迁移**（含 SSE Client） |
| 预览 | `views/studio/PreviewView.vue` | **迁移** |
| 个人资料 | `views/studio/ProfileView.vue` | **迁移** |
| 管理后台-模型 | `views/admin/ModelsView.vue` | **迁移** |
| 管理后台-Prompt | `views/admin/PromptsView.vue` | **迁移** |
| 管理后台-用户 | `views/admin/UsersView.vue` | **迁移** |
| 管理后台-项目 | `views/admin/ProjectsView.vue` | **迁移** |

---

## 2. 前端 Store

### Workbench

| Store | 文件 | 处置 |
|-------|------|------|
| projectStore | `stores/projectStore.ts` | **合并**（增加 Builder 项目字段） |
| reportStore | `stores/reportStore.ts` | 保留 |
| languageStore | `stores/languageStore.ts` | 保留 |

### Builder

| Store | 文件 | 处置 |
|-------|------|------|
| auth | `stores/auth.ts` | **迁移**（JWT + 用户角色） |
| project | `stores/project.ts` | 废弃（合并到 Workbench projectStore） |
| task | `stores/task.ts` | **迁移**（SSE + 进度） |

---

## 3. 前端 API Client

### Workbench

| API | 文件 | 处置 |
|-----|------|------|
| client | `api/client.ts` (内联) | **合并**（统一 JWT/Base URL/错误处理） |

### Builder

| API | 文件 | 处置 |
|-----|------|------|
| client | `api/client.ts` | 废弃（合并到统一 client） |
| auth | `api/auth.ts` | **迁移** |
| projects | `api/projects.ts` | **迁移** |
| requirements | `api/requirements.ts` | **迁移** |
| blueprints | `api/blueprints.ts` | **迁移** |
| tasks | `api/tasks.ts` | **迁移** |
| files | `api/files.ts` | **迁移** |
| preview | `api/preview.ts` | **迁移** |
| sse | `api/sse.ts` | **迁移** |
| admin | `api/admin.ts` | **迁移** |

---

## 4. 前端路由

| 来源 | 路由 | 处置 |
|------|------|------|
| Workbench | `/` → LandingPage | 保留 |
| Workbench | `/dashboard` | **合并** |
| Workbench | `/tools/ui-review` | 保留（加 project_id 参数） |
| Workbench | `/tools/project-doctor` | 保留（加 project_id 参数） |
| Workbench | `/tools/agent-config` | 保留 |
| Workbench | `/tools/api-doc` | 保留 |
| Workbench | `/tools/db-schema` | 保留 |
| Workbench | `/reports` | 保留 |
| Workbench | `/reports/:id` | 保留 |
| Workbench | `/projects` | **合并** |
| Workbench | `/projects/new` | **合并** |
| Workbench | `/projects/:id` | **合并** |
| Workbench | `/settings` | 保留 |
| Builder | `/login` | **迁移** |
| Builder | `/projects/create` | **迁移** |
| Builder | `/projects/:id/blueprint` | **迁移** |
| Builder | `/projects/:id/progress` | **迁移** |
| Builder | `/projects/:id/preview` | **迁移** |
| Builder | `/profile` | **迁移** |
| Builder | `/admin/models` | **迁移** |
| Builder | `/admin/prompts` | **迁移** |
| Builder | `/admin/users` | **迁移** |
| Builder | `/admin/projects` | **迁移** |

---

## 5. 后端 API

### Workbench — `/api` 前缀

| 路由 | Handler | 处置 |
|------|---------|------|
| GET `/api/health` | HealthRoutes | 保留（改 `/api/v1/health`） |
| GET `/api/config` | SystemRoutes | 保留 |
| GET `/api/observability/ai-runs` | ObservabilityRoutes | 保留 |
| GET `/api/dashboard` | DashboardRoutes | **合并** |
| POST `/api/tools/:tool/run` | ToolRunRoutes | 保留（改 `/api/v1`） |
| GET `/api/reports` | ReportRoutes | 保留 |
| GET/PUT `/api/reports/:id` | ReportRoutes | 保留 |
| GET `/api/projects` | ProjectRoutes | **合并** |
| POST `/api/projects` | ProjectRoutes | **合并** |
| GET/PUT/DELETE `/api/projects/:id` | ProjectRoutes | **合并** |
| POST `/api/export/reports/:id` | ExportRoutes | 保留 |
| GET `/api/jobs/:id` | JobRoutes | 保留 |

### Builder — `/api/v1` 前缀（已统一）

| 路由 | Handler | 处置 |
|------|---------|------|
| POST `/api/v1/auth/register` | AuthHandler | **迁移** |
| POST `/api/v1/auth/login` | AuthHandler | **迁移** |
| GET `/api/v1/auth/profile` | AuthHandler | **迁移** |
| PUT `/api/v1/auth/profile` | AuthHandler | **迁移** |
| GET/POST `/api/v1/projects` | ProjectHandler | **迁移**（合并到 Workbench 模块） |
| GET/PUT/DELETE `/api/v1/projects/:id` | ProjectHandler | **迁移** |
| GET/PUT `/api/v1/projects/:id/requirements` | RequirementHandler | **迁移** |
| POST `/api/v1/projects/:id/requirements/parse` | RequirementHandler | **迁移** |
| GET `/api/v1/projects/:id/blueprint` | BlueprintHandler | **迁移** |
| POST `/api/v1/projects/:id/blueprint/generate` | BlueprintHandler | **迁移** |
| PUT `/api/v1/projects/:id/blueprint` | BlueprintHandler | **迁移** |
| POST `/api/v1/projects/:id/blueprint/confirm` | BlueprintHandler | **迁移** |
| GET `/api/v1/tasks/:id` | TaskHandler | **迁移** |
| GET `/api/v1/tasks/:id/stream` | TaskHandler (SSE) | **迁移** |
| POST `/api/v1/tasks/:id/retry` | TaskHandler | **迁移** |
| POST `/api/v1/tasks/:id/cancel` | TaskHandler | **迁移** |
| GET `/api/v1/projects/:id/files` | FileHandler | **迁移** |
| GET `/api/v1/projects/:id/files/content` | FileHandler | **迁移** |
| POST `/api/v1/projects/:id/build` | PreviewHandler | **迁移** |
| GET `/api/v1/preview/:sessionId/*` | PreviewHandler | **迁移** |
| GET `/api/v1/projects/:id/export` | — | **迁移** |
| GET/POST/PUT/DELETE `/api/v1/admin/models` | AIModelHandler | **迁移** |
| GET/POST/PUT/DELETE `/api/v1/admin/prompts` | PromptTemplateHandler | **迁移** |
| GET `/api/v1/admin/users` | AdminHandler | **迁移** |
| GET `/api/v1/admin/projects` | AdminHandler | **迁移** |

---

## 6. 后端模型（GORM）

### 冲突模型映射

| 模型 | Workbench | Builder | 处置 |
|------|----------|---------|------|
| **Project** | `internal/model/` (字段少：name, description, status) | `internal/model/project.go` (字段多：type, source, quality_score, blueprint_id 等) | **合并**：以 Builder 为基础扩展 |
| **Task** | `internal/model/` (Job 模型) | `internal/model/task.go` | **合并**：统一 Job/Task |
| **File** | 无独立模型 | `internal/model/file.go` | **迁移** |
| **Report** | `internal/model/` | 无 | 保留 |
| **User** | 无 | `internal/model/user.go` | **迁移** |

### Builder 独有模型 → 迁移

| 模型 | 文件 | 处置 |
|------|------|------|
| User | `model/user.go` | **迁移** |
| Requirement | `model/requirement.go` | **迁移** |
| Blueprint | `model/blueprint.go` | **迁移** |
| Page | `model/page.go` | **迁移** |
| Component | `model/component.go` | **迁移** |
| File | `model/file.go` | **迁移** |
| AIModel | `model/ai_model.go` | **迁移** |
| PromptTemplate | `model/prompt_template.go` | **迁移** |

### Workbench 独有模型 → 保留

| 模型 | 处置 |
|------|------|
| Report | 保留 |
| GeneratedFile | 保留 |
| ReportAsset | 保留 |
| AIRun | 保留 |
| Job | 合并到 Task |

---

## 7. 数据库表

| 表 | 来源 | 处置 |
|----|------|------|
| users | Builder | **迁移** |
| projects | 两项目均有 | **合并**（以 Builder schema 为基础） |
| project_requirements | Builder | **迁移** |
| project_blueprints | Builder | **迁移** |
| tasks | Builder | **迁移**（合并 Workbench jobs） |
| task_progress | Builder | **迁移** |
| files | Builder | **迁移** |
| ai_models | Builder | **迁移** |
| prompt_templates | Builder | **迁移** |
| reports | Workbench | 保留 |
| report_assets | Workbench | 保留 |
| generated_files | Workbench | 保留 |
| jobs | Workbench | 废弃（合并到 tasks） |
| ai_runs | Workbench | 保留 |
| report_lineage | Workbench | 保留 |

---

## 8. 依赖差异

| 依赖 | Workbench | Builder | 目标 |
|------|----------|---------|------|
| Tailwind CSS | v4.3.1 | v3.4.19 | **v4** |
| 图标库 | `@lucide/vue` | `lucide-vue-next` | **`@lucide/vue`** |
| 路由模式 | `createWebHistory` | `createWebHashHistory` | **`createWebHistory`** |
| API 前缀 | `/api` | `/api/v1` | **`/api/v1`** |
| i18n | `vue-i18n` | 无 | 保留 |
| 测试 | Vitest + Playwright | 无 | 保留 |
| 数据库 | SQLite/MySQL | MySQL | **MySQL**（开发可用 SQLite） |
| 缓存 | 无 | Redis | **Redis** |
| 认证 | 无 | JWT | **JWT** |
