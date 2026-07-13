# A-01 目标系统上下文和模块边界

## 目标架构

```text
┌─────────────────────────────────────────────────┐
│                   Vue 3 前端                      │
│  ┌──────────┬──────────┬──────────┬──────────┐  │
│  │ 认证模块  │ 项目工作区 │ AI 工具  │ 管理后台  │  │
│  │ 登录/注册 │ 需求/蓝图 │ UI Review│ 模型/Prompt│  │
│  │ Profile  │ 生成/预览 │ Doctor   │ 用户/项目  │  │
│  │          │ 文件/导出 │ API Doc  │           │  │
│  │          │ 报告/质量 │ DB Schema│           │  │
│  └──────────┴──────────┴──────────┴──────────┘  │
│         Pinia Store → API Client (Axios)          │
└──────────────────────┬──────────────────────────┘
                       │ /api/v1 (JWT Bearer)
┌──────────────────────┴──────────────────────────┐
│                   Go API 服务                     │
│  ┌──────────┬──────────┬──────────┬──────────┐  │
│  │ 认证模块  │ 项目模块  │ 生成模块  │ 工具模块  │  │
│  │ Auth     │ Project  │ Pipeline │ UI Review│  │
│  │ JWT      │ Require  │ SSE      │ Doctor   │  │
│  │          │ Blueprint│ Task     │ Export   │  │
│  │          │ File     │ Build    │          │  │
│  └──────────┴──────────┴──────────┴──────────┘  │
│         Service → Repository → GORM              │
└──────────────────────┬──────────────────────────┘
                       │
┌──────────────────────┴──────────────────────────┐
│               MySQL + Redis                      │
│  users / projects / requirements / blueprints    │
│  tasks / files / reports / ai_models / prompts   │
└─────────────────────────────────────────────────┘
```

## 模块边界

### 1. 认证模块 (auth)
- **职责**：用户注册、登录、JWT 签发/验证、Profile 管理、角色管理
- **接口**：`POST /auth/register`、`POST /auth/login`、`GET/PUT /auth/profile`
- **依赖**：User Repository
- **被依赖**：全部模块（通过 JWT 中间件）

### 2. 项目模块 (project)
- **职责**：项目 CRUD、需求管理、蓝图管理、文件浏览、项目状态机
- **接口**：`/projects`、`/projects/:id/requirements`、`/projects/:id/blueprint`、`/projects/:id/files`
- **依赖**：认证模块（权限校验）、文件模块
- **被依赖**：生成模块、工具模块、报告模块

### 3. 生成模块 (generation)
- **职责**：代码生成 Pipeline、任务管理、SSE 进度推送、构建、预览
- **接口**：`/tasks`、`/tasks/:id/stream`、`/projects/:id/build`、`/preview/:session/*`
- **依赖**：项目模块（蓝图/需求）、认证模块
- **被依赖**：报告模块

### 4. 工具模块 (tools)
- **职责**：UI Review、Project Doctor、Agent Config、API Doc、DB Schema
- **接口**：`/tools/:tool/run`
- **依赖**：项目模块（文件/上下文）、报告模块（写入结果）
- **被依赖**：无

### 5. 报告模块 (report)
- **职责**：报告 CRUD、报告对比、导出、Action Items
- **接口**：`/reports`、`/reports/:id`、`/reports/:id/compare/:target`
- **依赖**：项目模块、认证模块
- **被依赖**：Dashboard、工具模块

### 6. 文件模块 (file)
- **职责**：项目工作区管理、文件读写、路径安全校验、ZIP 导出
- **接口**：内部 Service（被项目/生成/工具模块调用）
- **依赖**：认证模块
- **被依赖**：项目模块、生成模块、工具模块

### 7. 管理后台模块 (admin)
- **职责**：AI 模型管理、Prompt 模板管理、用户管理、项目管理
- **接口**：`/admin/models`、`/admin/prompts`、`/admin/users`、`/admin/projects`
- **依赖**：认证模块（管理员角色校验）
- **被依赖**：无

## 模块间通信规则

1. **仅通过 Service/Repository 接口依赖**：不直接访问其他模块的数据库表
2. **共享类型通过 DTO**：使用统一 `dto` 包定义跨模块数据传输对象
3. **项目上下文通过 Context 传递**：JWT 中间件注入 `user_id`，项目权限中间件注入 `project_id`
4. **文件操作统一走 File Service**：不允许业务模块直接操作文件系统
