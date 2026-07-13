# A-02 统一前端路由表

> 统一使用 `createWebHistory`，废除 Builder 的 `createWebHashHistory`

## 路由清单

### 公共路由

| 路径 | 名称 | 页面 | 来源 | 权限 |
|------|------|------|------|------|
| `/` | landing | LandingPage | Workbench | 公开 |
| `/login` | login | LoginView | Builder | 公开 |
| `/dashboard` | dashboard | DashboardPage | Workbench（扩展） | 需登录 |
| `/settings` | settings | SettingsPage | Workbench | 需登录 |
| `/profile` | profile | ProfileView | Builder | 需登录 |

### 项目工作区路由 `/projects/:projectId/`

| 路径 | 名称 | 页面 | 来源 |
|------|------|------|------|
| `/projects` | projects | ProjectsPage | 合并 |
| `/projects/new` | project-create | CreateProjectView | 合并（Builder 表单 + 项目类型选择） |
| `/projects/:projectId` | project-overview | ProjectDetailPage | Workbench（扩展） |
| `/projects/:projectId/edit` | project-edit | ProjectFormPage | Workbench |
| `/projects/:projectId/requirements` | project-requirements | 需求编辑页 | Builder（迁移） |
| `/projects/:projectId/blueprint` | project-blueprint | BlueprintReviewView | Builder（迁移） |
| `/projects/:projectId/generation` | project-generation | GenerationView | Builder（迁移 + SSE） |
| `/projects/:projectId/preview` | project-preview | PreviewView | Builder（迁移） |
| `/projects/:projectId/files` | project-files | FilesPage | 新建（含 Builder 文件 API） |
| `/projects/:projectId/reports` | project-reports | ReportsPage | Workbench（扩展项目筛选） |

### AI 工具路由

| 路径 | 名称 | 页面 | 来源 |
|------|------|------|------|
| `/tools/ui-review` | ui-review | UIReviewPage | Workbench（加 project_id 参数） |
| `/tools/project-doctor` | project-doctor | ProjectDoctorPage | Workbench（加 project_id 参数） |
| `/tools/agent-config` | agent-config | AgentConfigPage | Workbench |
| `/tools/api-doc` | api-doc | APIDocPage | Workbench |
| `/tools/db-schema` | db-schema | DBSchemaPage | Workbench |

### 报告路由

| 路径 | 名称 | 页面 | 来源 |
|------|------|------|------|
| `/reports` | reports | ReportsPage | Workbench |
| `/reports/:id` | report-detail | ReportDetailPage | Workbench |
| `/reports/:id/compare/:targetId` | report-compare | ReportComparePage | Workbench |

### 管理后台路由 `/admin/`

| 路径 | 名称 | 页面 | 来源 | 权限 |
|------|------|------|------|------|
| `/admin/models` | admin-models | ModelsView | Builder | 需管理员 |
| `/admin/prompts` | admin-prompts | PromptsView | Builder | 需管理员 |
| `/admin/users` | admin-users | UsersView | Builder | 需管理员 |
| `/admin/projects` | admin-projects | ProjectsView | Builder | 需管理员 |

### 错误路由

| 路径 | 页面 | 处置 |
|------|------|------|
| `/403` | ForbiddenPage | 新建 |
| `/404` | NotFoundPage | 新建 |
| `/:pathMatch(.*)*` | 重定向到 `/` | 保留 |

## 路由守卫

```typescript
router.beforeEach((to) => {
  const authStore = useAuthStore()
  // 未登录 → 登录页
  if (to.meta.requiresAuth !== false && !authStore.isLoggedIn) {
    return { name: 'login' }
  }
  // 已登录 → 不显示登录页
  if (to.name === 'login' && authStore.isLoggedIn) {
    return { name: 'dashboard' }
  }
  // 管理员路由校验
  if (to.meta.requiresAdmin && authStore.user?.role !== 'admin') {
    return { name: '403' }
  }
})
```

## 废除的路由

| 原路由 | 来源 | 说明 |
|--------|------|------|
| `#/` → HomeView | Builder (hash) | 合并到 Workbench Dashboard |
| `#/projects` | Builder (hash) | 合并到 Workbench `/projects` |

## 路由模式变更

- **从**：Builder 使用 `createWebHashHistory()` → `/#/login`、`/#/projects`
- **到**：统一 `createWebHistory()` → `/login`、`/projects`
- **兼容**：Nginx 配置 SPA History Fallback，404 返回 `index.html`
