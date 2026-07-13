# P1-12 to P1-16 Project Profile Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Complete the project-profile vertical slice: validated project CRUD, report ownership, project-aware tool runs, and a project detail view with history, trends, and generated artifacts.

**Architecture:** Keep project persistence behind `ProjectRepository`, expose project-specific read models through `ProjectService`, and let `ReportService` resolve a selected project before a tool creates its processing report. The Vue application receives a shared project picker for all five tool forms and uses the project-detail API for history and aggregate displays.

**Tech Stack:** Go 1.25, Gin, GORM, SQLite test driver, Vue 3, Pinia, TypeScript, Vitest.

---

### Task 1: Close P1-12/P1-13 Backend Contracts

**Files:**
- Modify: `backend/internal/dto/project_dto.go`
- Modify: `backend/internal/repository/project_repository.go`
- Modify: `backend/internal/service/project_service.go`
- Modify: `backend/internal/handler/project_handler.go`
- Modify: `backend/cmd/server/main_test.go`
- Create: `backend/internal/service/project_service_test.go`

- [ ] **Step 1: Define bounded project inputs and explicit delete output**

```go
const (
    ProjectNameMaxLength = 128
    ProjectURLMaxLength = 512
    ProjectDescriptionMaxLength = 4000
    ProjectCodingRulesMaxLength = 12000
)

type ProjectDeleteDTO struct {
    DetachedReportCount int64 `json:"detached_report_count"`
}
```

- [ ] **Step 2: Make project summaries include aggregate report data**

```go
SELECT p.*, COALESCE(s.report_count, 0) AS report_count, s.average_score
FROM projects p
LEFT JOIN (
    SELECT project_id, COUNT(*) AS report_count, AVG(total_score) AS average_score
    FROM reports
    WHERE project_id IS NOT NULL
    GROUP BY project_id
) s ON s.project_id = p.id
ORDER BY p.updated_at DESC
```

- [ ] **Step 3: Validate, create, update, list, delete, and fetch project reports**

```go
func validateProjectInput(input dto.ProjectCreateDTO) error {
    if strings.TrimSpace(input.Name) == "" {
        return fmt.Errorf("project name is required")
    }
    // Validate all documented text limits and an optional http/https repository URL.
    return nil
}
```

- [ ] **Step 4: Test service behavior with SQLite-backed repositories**

```go
func TestProjectServiceCreateRejectsInvalidRepositoryURL(t *testing.T) {
    _, err := svc.Create(ctx, dto.ProjectCreateDTO{Name: "Workbench", RepoURL: "git@host:repo"})
    require.ErrorContains(t, err, "repo_url")
}
```

- [ ] **Step 5: Verify the backend slice**

Run: `cd backend; go test ./internal/service ./internal/handler ./cmd/server`

Expected: project validation, pagination, statistics, deletion semantics, and the route inventory pass.

### Task 2: Implement P1-15 Project Ownership for Tool Runs

**Files:**
- Modify: `backend/internal/service/report_service.go`
- Modify: `backend/internal/dto/{ui_review,project_doctor,agent_config,api_doc,db_schema}_dto.go`
- Modify: `backend/internal/handler/tool_run_handler.go`
- Modify: `backend/internal/service/tools/{ui_review,project_doctor,agent_config,api_doc,db_schema}_service.go`
- Modify: `backend/internal/prompts/prompt_builder.go`
- Modify: affected `*_test.go` fake `ReportService` implementations

- [ ] **Step 1: Resolve a selected project before report creation**

```go
ResolveProjectContext(ctx context.Context, projectID string) (*model.Project, error)
CreateProcessingReport(
    ctx context.Context,
    toolType, title, inputMode string,
    inputData json.RawMessage,
    parentReportID, projectID string,
) (*model.Report, error)
```

- [ ] **Step 2: Carry `project_id` through all five JSON and multipart request paths**

```go
type UIReviewFormInput struct {
    // Existing fields...
    ProjectID string
}

projectID := c.PostForm("project_id")
```

- [ ] **Step 3: Append trusted project context without treating uploaded content as trusted**

```go
func AppendTrustedProjectContext(userPrompt string, project *model.Project) string {
    if project == nil {
        return userPrompt
    }
    return userPrompt + "\n\nTrusted project profile:\n" +
        "- Frontend: " + project.FrontendStack + "\n" +
        "- Backend: " + project.BackendStack + "\n" +
        "- UI style: " + project.UIStyle + "\n" +
        "- Coding rules: " + project.CodingRules
}
```

- [ ] **Step 4: Test valid ownership and invalid-project rejection**

```go
func TestCreateProcessingReportRejectsMissingProject(t *testing.T) {
    _, err := svc.CreateProcessingReport(ctx, model.ToolTypeUIReview, "Review", "code", []byte(`{}`), "", "missing")
    require.ErrorContains(t, err, "project not found")
}
```

- [ ] **Step 5: Verify tool and report tests**

Run: `cd backend; go test ./internal/service/... ./internal/handler/...`

Expected: all fake services match the interface and no invalid project can create an orphan processing report.

### Task 3: Complete P1-16 Aggregates and Project Detail APIs

**Files:**
- Modify: `backend/internal/dto/project_dto.go`
- Modify: `backend/internal/repository/project_repository.go`
- Modify: `backend/internal/service/project_service.go`
- Modify: `backend/internal/handler/project_handler.go`
- Modify: `frontend/src/api/projects.ts`
- Modify: `frontend/src/types/project.ts`
- Modify: `frontend/src/stores/projectStore.ts`
- Modify: `frontend/src/pages/ProjectDetailPage.vue`

- [ ] **Step 1: Return report history, score trend, high-risk count, and latest generated artifacts**

```go
type ProjectTrendPointDTO struct {
    Date         string   `json:"date"`
    ReportCount  int64    `json:"report_count"`
    AverageScore *float64 `json:"average_score"`
    HighRiskCount int64   `json:"high_risk_count"`
}

type ProjectArtifactDTO struct {
    ToolType string `json:"tool_type"`
    ReportID string `json:"report_id"`
    Filename string `json:"filename"`
}
```

- [ ] **Step 2: Register project report history**

```go
r.GET("/projects/:id/reports", h.ListReports)
```

- [ ] **Step 3: Render trends, paginated history, and latest artifact links without per-project fan-out**

```ts
await Promise.all([
  store.fetchProject(projectId.value),
  store.fetchStats(projectId.value),
  store.fetchReports(projectId.value),
])
```

- [ ] **Step 4: Verify output**

Run: `cd backend; go test ./...`

Expected: project stats handle empty, unscored, and scored reports; artifacts and history remain scoped to the requested project.

### Task 4: Complete P1-14/P1-15 Frontend Workflows

**Files:**
- Create: `frontend/src/components/tool/ProjectPicker.vue`
- Create: `frontend/src/components/tool/ProjectPicker.test.ts`
- Modify: `frontend/src/pages/ProjectFormPage.vue`
- Modify: `frontend/src/pages/ProjectDetailPage.vue`
- Modify: `frontend/src/pages/tools/{UIReviewPage,ProjectDoctorPage,AgentConfigPage,APIDocPage,DBSchemaPage}.vue`
- Modify: `frontend/src/components/layout/Sidebar.vue`

- [ ] **Step 1: Add a searchable, optional project picker**

```vue
<ProjectPicker
  v-model="projectId"
  input-id="ui-review-project"
  help-id="ui-review-project-help"
/>
```

- [ ] **Step 2: Preselect project IDs from project-detail tool links**

```ts
const initialProjectID = typeof route.query.project_id === 'string'
  ? route.query.project_id
  : ''
```

- [ ] **Step 3: Add safe unsaved-form navigation confirmation**

```ts
onBeforeRouteLeave(() => {
  if (!saving.value && isDirty.value && !window.confirm('离开后未保存的项目修改将丢失，确定离开吗？')) {
    return false
  }
})
```

- [ ] **Step 4: Verify frontend behavior**

Run: `cd frontend; npm.cmd run test:unit; npm.cmd run build`

Expected: the picker filters projects, project links carry `project_id`, forms protect unsaved edits, and TypeScript succeeds.

### Task 5: Regression Gate and Documentation

**Files:**
- Modify: `项目优化开发任务拆分.md`

- [ ] **Step 1: Run all project checks**

Run:

```powershell
cd backend; go test ./...
cd ../frontend; npm.cmd run test:unit
npm.cmd run build
npm.cmd run test:e2e
```

Expected: all commands pass, or any external-environment exception is recorded precisely.

- [ ] **Step 2: Update P1-12 through P1-16 only for verified behavior**

```markdown
## [x] P1-12 新增 Project 数据模型
## [x] P1-13 实现项目 CRUD 与统计 API
## [x] P1-14 实现项目列表与详情页面
## [x] P1-15 让五工具运行关联项目
## [x] P1-16 聚合项目趋势与最新产物
```

- [ ] **Step 3: Commit**

```bash
git add backend frontend 项目优化开发任务拆分.md docs/superpowers/plans/2026-07-10-p1-12-project-profile.md
git commit -m "feat(P1-12): complete project profile workflow"
```
