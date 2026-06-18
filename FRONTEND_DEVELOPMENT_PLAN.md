# AI Developer Workbench 前端开发方案

> 版本：MVP 0.1.0  
> 技术栈：Vue 3 + Vite + TypeScript + Tailwind CSS  
> 对应总方案：[DEVELOPMENT_PLAN.md](./DEVELOPMENT_PLAN.md)  
> 后端契约：[BACKEND_DEVELOPMENT_PLAN.md](./BACKEND_DEVELOPMENT_PLAN.md)

## 1. 前端目标

前端负责为五个 AI 开发工具提供统一、清晰、可响应的操作界面，并完成从用户输入到报告查看、复制和下载的完整交互闭环。

五个核心工具：

1. UI Review
2. Project Doctor
3. Agent Config Studio
4. API Doc Builder
5. DB Schema Review

统一用户流程：

```text
选择工具
  -> 填写表单或上传文件
  -> 前端校验
  -> 提交后端分析
  -> 显示运行状态
  -> 跳转报告详情
  -> 查看评分、问题和建议
  -> 复制 Codex Prompt
  -> 下载报告或生成文件
```

MVP 不实现登录注册、权限、支付、团队空间、在线代码执行和前端 API Key 配置。

## 2. 前端职责边界

前端负责：

- 页面路由和整体布局。
- 五个工具的输入表单与客户端校验。
- 图片和 ZIP 文件选择、格式提示与上传。
- 调用后端 API 并统一处理加载、成功和失败状态。
- Dashboard、报告列表、报告详情和 Settings 展示。
- Markdown 安全渲染。
- Codex Prompt 与生成文件的复制、下载。
- 响应式、可访问性和基础前端测试。

前端不负责：

- 连接 MySQL。
- 调用真实 AI Provider。
- 保存或展示 AI API Key。
- 解压 ZIP、扫描项目或执行 SQL。
- 信任浏览器端文件校验作为最终安全边界。
- 自行解析数据库原始 JSON 字符串。

## 3. 技术选型

- Vue 3
- Vite
- TypeScript
- Vue Router
- Pinia
- Axios
- Tailwind CSS
- lucide-vue-next
- markdown-it
- Vitest
- Vue Test Utils
- Playwright

约束：

- 使用 Composition API 和 `<script setup lang="ts">`。
- 不引入重型 UI 框架和后台模板。
- 不引入大型动画、图表或低代码表单库。
- 通用交互组件自行封装，业务组件保持明确职责。
- TypeScript 开启严格模式，避免滥用 `any`。

## 4. 推荐目录结构

```text
frontend/
├─ src/
│  ├─ api/
│  │  ├─ client.ts
│  │  ├─ tools.ts
│  │  ├─ reports.ts
│  │  └─ system.ts
│  ├─ components/
│  │  ├─ common/
│  │  │  ├─ AppButton.vue
│  │  │  ├─ AppDialog.vue
│  │  │  ├─ AppSelect.vue
│  │  │  ├─ AppTextarea.vue
│  │  │  ├─ Badge.vue
│  │  │  ├─ EmptyState.vue
│  │  │  ├─ ErrorState.vue
│  │  │  ├─ LoadingState.vue
│  │  │  └─ SectionHeader.vue
│  │  ├─ layout/
│  │  │  ├─ AppShell.vue
│  │  │  ├─ AppHeader.vue
│  │  │  └─ Sidebar.vue
│  │  ├─ dashboard/
│  │  │  ├─ ToolCard.vue
│  │  │  ├─ RecentReports.vue
│  │  │  └─ WorkbenchStats.vue
│  │  ├─ tools/
│  │  │  ├─ ToolHeader.vue
│  │  │  ├─ InputSection.vue
│  │  │  ├─ FileUpload.vue
│  │  │  ├─ CodeInput.vue
│  │  │  ├─ OptionSelect.vue
│  │  │  └─ RunToolButton.vue
│  │  └─ report/
│  │     ├─ ReportHeader.vue
│  │     ├─ ScoreOverview.vue
│  │     ├─ ScoreBreakdown.vue
│  │     ├─ IssueList.vue
│  │     ├─ RecommendationPanel.vue
│  │     ├─ GeneratedFilesPanel.vue
│  │     ├─ CodexPromptBox.vue
│  │     ├─ MarkdownViewer.vue
│  │     └─ ReportHistoryItem.vue
│  ├─ pages/
│  │  ├─ tools/
│  │  │  ├─ UIReviewPage.vue
│  │  │  ├─ ProjectDoctorPage.vue
│  │  │  ├─ AgentConfigPage.vue
│  │  │  ├─ APIDocPage.vue
│  │  │  └─ DBSchemaPage.vue
│  │  ├─ LandingPage.vue
│  │  ├─ DashboardPage.vue
│  │  ├─ ReportsPage.vue
│  │  ├─ ReportDetailPage.vue
│  │  └─ SettingsPage.vue
│  ├─ router/
│  │  └─ index.ts
│  ├─ stores/
│  │  ├─ toolStore.ts
│  │  ├─ reportStore.ts
│  │  └─ systemStore.ts
│  ├─ types/
│  │  ├─ api.ts
│  │  ├─ tool.ts
│  │  ├─ report.ts
│  │  └─ system.ts
│  ├─ utils/
│  │  ├─ clipboard.ts
│  │  ├─ download.ts
│  │  ├─ format.ts
│  │  └─ validation.ts
│  ├─ styles/
│  │  └─ globals.css
│  ├─ App.vue
│  └─ main.ts
├─ tests/
│  └─ e2e/
├─ .env.example
├─ Dockerfile
├─ index.html
├─ package.json
├─ tailwind.config.js
├─ tsconfig.json
└─ vite.config.ts
```

## 5. 路由设计

```text
/
/dashboard
/tools/ui-review
/tools/project-doctor
/tools/agent-config
/tools/api-doc
/tools/db-schema
/reports
/reports/:id
/settings
```

路由规则：

- `/` 使用独立 Landing 布局。
- Dashboard、工具页、报告页和 Settings 使用 `AppShell`。
- `/reports/:id` 在报告不存在时展示明确的 404 状态，而不是空白页。
- 路由切换后主内容区回到顶部。
- MVP 无鉴权，不设置登录守卫。
- 未知路径跳转到 Landing 或专用 Not Found 页面。

## 6. API 契约

统一前缀由环境变量控制：

```env
VITE_API_BASE_URL=http://localhost:8080/api
```

统一成功响应：

```ts
interface ApiResponse<T> {
  code: 0
  message: string
  data: T
}
```

统一失败响应：

```ts
interface ApiErrorResponse {
  code: number
  message: string
  error?: string
  request_id?: string
}
```

前端调用的接口：

| 方法 | 路径 | 用途 |
| --- | --- | --- |
| GET | `/health` | 后端连接状态 |
| GET | `/tools` | 工具元数据 |
| GET | `/system/status` | Mock、模型和上传限制 |
| GET | `/dashboard/stats` | Dashboard 统计与最近报告 |
| POST | `/tools/ui-review/run` | 运行 UI Review |
| POST | `/tools/project-doctor/run` | 运行 Project Doctor |
| POST | `/tools/agent-config/run` | 运行 Agent Config |
| POST | `/tools/api-doc/run` | 运行 API Doc |
| POST | `/tools/db-schema/run` | 运行 DB Schema |
| GET | `/reports` | 分页获取报告 |
| GET | `/reports/:id` | 获取报告详情 |
| DELETE | `/reports/:id` | 删除报告 |
| GET | `/reports/:id/export?format=markdown` | 下载 Markdown 报告 |
| GET | `/reports/:id/files/:filename` | 下载生成文件 |

## 7. API 客户端设计

### 7.1 Axios 实例

`src/api/client.ts` 负责：

- 使用 `VITE_API_BASE_URL`。
- 默认请求超时建议 100 秒，略高于后端 AI 超时。
- 统一解析业务错误。
- 将网络错误、超时和后端错误转换为一致的前端错误对象。
- 不在请求或日志中处理 AI API Key。
- 不为 multipart 请求手工设置 boundary。

### 7.2 API 模块

`tools.ts`：

```ts
runUIReview(formData: FormData)
runProjectDoctor(formData: FormData)
runAgentConfig(payload: AgentConfigInput)
runAPIDoc(formData: FormData)
runDBSchema(payload: DBSchemaInput)
getTools()
```

`reports.ts`：

```ts
listReports(params: ReportListParams)
getReport(id: string)
deleteReport(id: string)
downloadReport(id: string)
downloadGeneratedFile(reportId: string, filename: string)
```

`system.ts`：

```ts
getHealth()
getSystemStatus()
getDashboardStats()
```

### 7.3 文件下载

- 请求使用 `responseType: 'blob'`。
- 优先读取后端 `Content-Disposition` 中的文件名。
- 使用临时 Object URL 触发浏览器下载。
- 下载完成后释放 Object URL。
- 服务端返回 JSON 错误时，不把错误响应保存成文件。

## 8. TypeScript 类型设计

### 8.1 通用类型

```ts
export type ToolType =
  | 'ui_review'
  | 'project_doctor'
  | 'agent_config'
  | 'api_doc'
  | 'db_schema'

export type ReportStatus =
  | 'processing'
  | 'succeeded'
  | 'fallback'
  | 'failed'

export interface GeneratedFileMeta {
  id: string
  filename: string
  language?: string
  mime_type: string
  size_bytes: number
}

export interface Report<T = Record<string, unknown>> {
  id: string
  tool_type: ToolType
  title: string
  input_mode: string
  status: ReportStatus
  summary: string
  total_score: number | null
  grade: string | null
  input_data: Record<string, unknown>
  report_data: T
  generated_files: GeneratedFileMeta[]
  created_at: string
  updated_at: string
}
```

### 8.2 问题与摘要

```ts
export interface ReportSummary {
  one_sentence: string
  strengths: string[]
  weaknesses: string[]
  top_priorities: string[]
}

export interface Issue {
  title: string
  severity: 'high' | 'medium' | 'low'
  category: string
  problem: string
  suggestion: string
  action: string
}
```

### 8.3 工具专属类型

每个工具声明独立的输入和结果类型：

- `UIReviewInput`、`UIReviewResult`
- `ProjectDoctorInput`、`ProjectDoctorResult`
- `AgentConfigInput`、`AgentConfigResult`
- `APIDocInput`、`APIDocResult`
- `DBSchemaInput`、`DBSchemaResult`

Agent Config 和 API Doc 的 `total_score`、`grade` 可以为 `null`，前端不能人为补造评分。

## 9. 状态管理

Pinia 只保存跨页面共享状态。

### 9.1 toolStore

职责：

- 加载五个工具元数据。
- 提供按 `tool_type` 查询工具信息的方法。
- 缓存加载状态和错误状态。

### 9.2 reportStore

职责：

- 报告列表、总数和分页条件。
- 工具筛选和排序条件。
- 最近报告。
- 报告详情缓存。
- 删除报告后同步清理本地状态。

不保存工具表单草稿，避免 Store 膨胀。

### 9.3 systemStore

职责：

- 后端健康状态。
- Mock Mode 状态。
- 当前 Provider 和模型名称。
- 上传大小、文件数量等限制。

Settings 页面只能展示后端返回的公开配置，不提供 API Key 输入框。

## 10. 页面开发方案

### 10.1 Landing Page

内容：

- 顶部导航。
- Hero。
- 五个工具介绍。
- 工作流说明。
- 示例报告预览。
- Dashboard CTA。
- Footer。

文案：

```text
Build better AI-generated projects.
```

```text
Review UI quality, inspect project structure, generate AGENTS.md,
build API docs, and improve database schemas in one developer workbench.
```

设计重点：

- 真实开发者产品感。
- 避免模板化 SaaS Hero。
- 不使用大面积渐变、发光和玻璃拟态。
- 示例报告应表现实际产品能力，而不是纯装饰插图。

### 10.2 Dashboard

数据来源：`GET /dashboard/stats` 和 `GET /tools`。

展示：

- 五个工具入口。
- 报告总数。
- 各工具使用数量。
- 可计算时展示平均分。
- 最近报告。
- 快捷操作。

接口不可用时：

- 页面骨架仍可正常呈现。
- 工具入口可使用本地静态兜底元数据。
- 统计区域展示可恢复错误状态，不能整页白屏。

### 10.3 UI Review 页面

字段：

- `title`
- `review_mode`
- `page_type`
- `target_style`
- `description`
- `code`
- `screenshot`

联动校验：

- screenshot：必须上传图片。
- code：必须填写代码。
- screenshot_code：图片和代码都必填。
- 图片仅接受 PNG、JPEG、WebP。
- 页面显示后端返回的上传大小限制。

提交：

- 使用 `FormData`。
- 禁用重复提交。
- 成功后跳转 `/reports/:id`。
- 失败时保留输入和已选择文件。

### 10.4 Project Doctor 页面

字段：

- `title`
- `project_name`
- `tech_stack`
- `project_description`
- `analysis_depth`
- `project_zip`

交互：

- 仅允许 ZIP。
- 清楚提示“只做静态分析，不执行代码”。
- 展示文件名和大小。
- basic、standard、deep 均显示，MVP 可提示 deep 与 standard 接近。

### 10.5 Agent Config 页面

字段：

- `title`
- `project_name`
- `project_type`
- `frontend_stack`
- `backend_stack`
- `database`
- `ui_style`
- `coding_preferences`
- `strict_rules`

交互：

- JSON 请求。
- 长文本区域显示字符数。
- 提交成功后，报告页展示并支持下载 `AGENTS.md`、`TASK_PLAN.md` 等文件。

### 10.6 API Doc 页面

字段：

- `title`
- `source_type`
- `backend_stack`
- `code`
- `api_description`
- `output_format`
- `project_zip`

模式联动：

- code：代码必填。
- project_zip：ZIP 必填。
- manual：API 描述必填。
- output_format 支持 markdown、openapi、both。

### 10.7 DB Schema 页面

字段：

- `title`
- `schema_type`
- `database_type`
- `business_context`
- `schema_content`
- `target_goal`

交互：

- `schema_content` 必填。
- 提示用户 SQL 只会作为文本分析，不会执行。
- `generate_migration` 结果必须展示“执行前人工审查”提示。

### 10.8 Reports 页面

功能：

- 分页。
- 按工具筛选。
- 按 newest、oldest、score_desc、score_asc 排序。
- 查看详情。
- 删除确认。
- 空状态和错误重试。

注意：

- 无评分报告在分数排序中按后端返回顺序展示。
- 删除成功后更新当前页；当前页变空时回退上一页。
- 删除操作应明确说明会同时删除关联上传文件和生成文件。

### 10.9 Report Detail 页面

通用区域：

- 报告标题、工具类型和创建时间。
- 状态、总分和等级。
- Summary。
- Scores。
- Issues。
- Recommendations。
- Generated Files。
- Codex Prompt。
- 脱敏后的原始输入摘要。

专属渲染：

- UI Review：AI 模板风险和视觉评分。
- Project Doctor：项目类型与工程评分。
- Agent Config：生成文件为主要内容，不强制显示分数。
- API Doc：模块、端点、Markdown 和 OpenAPI。
- DB Schema：评分、优化 Schema 和迁移建议。

状态处理：

- `fallback`：显示温和提示，报告仍可正常使用。
- `failed`：显示失败原因和返回工具页入口。
- `processing`：显示加载状态并提供刷新。

### 10.10 Settings 页面

只读展示：

- 后端服务状态。
- 是否 Mock Mode。
- AI Provider。
- 文本模型和视觉模型。
- 上传大小限制。
- ZIP 文件数和读取限制。
- API Key 配置说明。

禁止：

- 在前端输入、保存或回显 API Key。
- 展示数据库密码或完整 DSN。

## 11. 通用组件设计

### 11.1 AppShell

- 桌面端左侧 Sidebar。
- 顶部 Header。
- 主内容区。
- 移动端抽屉式 Sidebar。
- 主内容提供一致的最大宽度和边距。

### 11.2 FileUpload

- 支持图片和 ZIP 两种配置模式。
- 显示文件名、大小、格式提示。
- 支持重新选择和清空。
- 拖拽上传必须有键盘可用的等价操作。
- 客户端校验失败不发送请求。

### 11.3 CodeInput

- textarea。
- 固定最小高度。
- 字符数统计。
- 清空按钮。
- 支持代码、SQL 和 Schema。
- 等宽字体。

### 11.4 ScoreOverview

- 允许 `total_score` 为 `null`。
- 分数颜色遵循统一阈值。
- 等级文字由后端返回，不在多个组件中重复推导。

### 11.5 IssueList

- severity 文本、图标和颜色同时表达。
- 支持展开问题、建议和行动项。
- 默认优先显示 high severity。

### 11.6 GeneratedFilesPanel

- 文件名、语言、大小。
- 复制文本内容。
- 下载文件。
- 对二进制或不返回文本内容的文件仅提供下载。

### 11.7 CodexPromptBox

- Markdown 安全展示。
- 一键复制。
- 下载 `.md`。
- 复制成功使用非阻塞提示。

## 12. UI 视觉规范

颜色变量：

```css
:root {
  --background: #f7f7f5;
  --surface: #ffffff;
  --surface-muted: #f1f1ef;
  --border: #e5e5e0;
  --text-primary: #18181b;
  --text-secondary: #52525b;
  --text-muted: #71717a;
  --accent: #2563eb;
  --accent-soft: #dbeafe;
  --danger: #dc2626;
  --warning: #d97706;
  --success: #16a34a;
}
```

布局：

- 常规内容最大宽度 1180px。
- Dashboard 最大宽度 1280px。
- 桌面端左右 padding 32px。
- 移动端左右 padding 16px。
- 使用留白、排版、分隔线和轻量面板表达层级。
- 常规区域不使用厚重阴影。
- 动效控制在 150–200ms。

禁止：

- 大面积渐变。
- Neon glow。
- 过度玻璃拟态。
- 所有内容卡片化。
- 花哨入场动画。
- 模板感强的 SaaS 首页。

## 13. 响应式和可访问性

响应式断点至少覆盖：

- 手机：单列，Sidebar 收起，操作按钮可全宽。
- 平板：表单保持单列或有限双列。
- 桌面：工具输入与说明可分栏，报告评分可网格展示。

可访问性要求：

- 所有输入具有可见 label。
- 错误信息与字段关联。
- 所有按钮可键盘操作。
- 焦点状态清晰。
- Dialog 支持焦点锁定和 Escape 关闭。
- 图标按钮提供 `aria-label`。
- 不只使用颜色表示状态。
- 正文和交互文本满足基本对比度。

## 14. 安全要求

- Markdown 禁止直接渲染原始 HTML，或经过严格净化。
- 不使用 `v-html` 渲染未经净化的 AI 内容。
- 文件扩展名校验只是用户体验，不能代替后端校验。
- 不在 LocalStorage、SessionStorage 或 Pinia 持久化中保存用户代码、ZIP 内容和敏感配置。
- 不记录完整 API 响应中的代码和 Prompt。
- 下载文件名以服务端响应为准并进行基本清洗。
- Settings 不展示任何 Secret。

## 15. 错误与加载体验

统一处理：

- 网络不可用。
- 后端超时。
- 参数错误。
- 文件过大。
- 不支持的文件格式。
- 不安全 ZIP。
- 报告不存在。
- AI fallback。
- 服务端内部错误。

原则：

- 表单提交失败保留用户输入。
- 可恢复错误提供重试入口。
- 错误信息使用用户可理解的中文。
- `request_id` 可作为问题反馈信息展示，但不突出技术细节。
- 长任务显示明确状态，不使用无限无文案旋转图标。

## 16. 测试方案

### 16.1 单元与组件测试

- 工具表单模式联动与必填校验。
- FileUpload 文件类型和大小提示。
- 报告状态和空分数渲染。
- IssueList 展开。
- Copy 和 Download 工具函数。
- API 错误转换。
- Markdown 安全配置。

### 16.2 页面测试

- Dashboard 正常、空数据和后端失败状态。
- Reports 筛选、排序、分页和删除。
- 五种 Report Detail 渲染。
- Settings Mock 状态。
- 工具提交成功跳转和失败保留输入。

### 16.3 E2E 冒烟测试

Mock Mode 下完成：

1. 访问 Landing 和 Dashboard。
2. 分别运行五个工具。
3. 查看每份报告。
4. 复制 Codex Prompt。
5. 下载 Markdown 和生成文件。
6. 筛选历史报告。
7. 删除报告。
8. 验证移动端关键页面无横向溢出。

## 17. 开发顺序

### 阶段 1：工程基础

- 初始化 Vue 3、Vite、TypeScript。
- 配置 Tailwind、Router、Pinia、Axios。
- 建立类型、API Client 和全局样式。
- 实现 Landing 基础页面。

### 阶段 2：产品框架

- AppShell、Sidebar、Header。
- 通用 Button、Badge、Loading、Error、Empty。
- Dashboard、Settings。
- 后端健康和公开状态接入。

### 阶段 3：报告中心

- Report 类型体系。
- Reports 页面。
- Report Detail 通用结构。
- 五类专属报告区块。
- 复制和下载。

### 阶段 4：五个工具页面

建议顺序：

1. Agent Config Studio
2. DB Schema Review
3. UI Review
4. Project Doctor
5. API Doc Builder

该顺序先验证普通 JSON、生成文件，再接入图片和 ZIP 上传。

### 阶段 5：质量收口

- 响应式。
- 可访问性。
- 错误与边界状态。
- Vitest 和 Playwright。
- Dockerfile、环境变量示例和前端 README。

## 18. 前端验收标准

- 所有路由可正常访问。
- Landing 与 Dashboard 具有真实、克制的开发者产品感。
- 五个工具表单字段、模式联动和校验正确。
- 图片和 ZIP 上传具有明确格式及大小提示。
- 提交期间不能重复运行。
- 成功后跳转对应报告详情。
- Reports 支持筛选、排序、分页和删除。
- 五类报告均可正确展示，空分数不会导致异常。
- fallback 报告有清晰提示且仍可使用。
- Codex Prompt 可以复制。
- Markdown 报告和生成文件可以下载。
- Settings 只展示公开配置，不出现 API Key。
- 后端暂时不可用时页面不会白屏。
- Markdown 内容不会引入明显 XSS 风险。
- 手机、平板和桌面端均无关键布局问题。
- 核心组件测试和 Mock Mode E2E 通过。

