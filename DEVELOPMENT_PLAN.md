# AI Developer Workbench 开发方案

> 版本：MVP 0.1.0  
> 技术路线：Vue 3 + TypeScript + Golang + Gin + GORM + MySQL  
> 项目形态：前后端分离的单用户开发者工作台  
> 方案原则：先打通五个工具的完整闭环，再扩展账号、团队和商业化能力

## 1. 项目目标

AI Developer Workbench 是一个面向独立开发者、前后端工程师和 AI Coding 用户的开发质量工作台。MVP 围绕以下五个工具建立统一工作流：

1. UI Review：根据截图和前端代码审查 UI 质量。
2. Project Doctor：静态检查项目 ZIP 的结构与工程质量。
3. Agent Config Studio：生成 `AGENTS.md`、`TASK_PLAN.md` 等 Agent 配置文件。
4. API Doc Builder：根据代码或项目 ZIP 生成 Markdown/OpenAPI 文档。
5. DB Schema Review：审查 SQL、GORM、Prisma 等数据库结构。

所有工具统一遵循以下闭环：

```text
用户输入
  -> 输入校验与安全处理
  -> 构建工具上下文
  -> AI 或 Mock 分析
  -> JSON 容错解析与结构校验
  -> 生成统一报告及附属文件
  -> MySQL 持久化
  -> 报告查看、复制、下载和删除
```

## 2. 对原始 Prompt 的关键调整

原始 Prompt 指定 SQLite，本方案根据项目要求全面改为 MySQL，其余核心产品范围保持不变。

| 项目 | 原始要求 | 本方案 |
| --- | --- | --- |
| 数据库 | SQLite | MySQL 8.0+ |
| 数据库驱动 | GORM SQLite Driver | `gorm.io/driver/mysql` |
| JSON 保存 | TEXT 字符串 | MySQL `JSON` 类型 |
| 数据初始化 | SQLite 自动建库 | MySQL 容器建库 + 版本化迁移 |
| 文件关联 | Report 单个路径 | Report 字段兼容 + `report_assets` 管理上传文件 |
| 报告删除 | 分步删除 | MySQL 事务 + 外键级联 + 文件清理补偿 |
| Dashboard 数据 | 原 Prompt 未定义 API | 增加统计接口 |
| Settings 数据 | 原 Prompt 未定义 API | 增加公开运行状态接口 |
| Markdown 报告下载 | 有产品要求，无明确 API | 增加报告导出接口 |

MVP 仍然不实现登录注册、支付、团队协作、GitHub OAuth、在线代码执行和复杂权限系统。由于没有鉴权，MVP 默认定位为本地或可信网络中的单用户工具，不应直接裸露到公网。

## 3. 总体架构

```text
Browser
  |
  | HTTP/JSON、multipart/form-data
  v
Vue 3 SPA
  |
  | /api
  v
Gin API
  |-- Handler：参数绑定、校验、统一响应
  |-- Tool Service：五个工具的业务编排
  |-- AI Service：Mock/OpenAI-compatible Provider
  |-- File/ZIP Service：上传、解压、静态扫描
  |-- Report Service：报告持久化、导出、删除
  |-- Repository：GORM 数据访问
  |
  +--> MySQL：报告、生成文件、上传资产
  +--> Local Storage：截图、ZIP、临时解压目录
  +--> AI Provider：文本/视觉分析
```

核心架构约束：

- 前端不保存、不接收 AI API Key。
- Handler 不直接访问数据库或 AI Provider。
- 五个工具使用独立输入 DTO、Prompt Builder、结果 DTO 和 Mock 数据。
- 数据库保存完整结构化 JSON，同时保存摘要、分数等高频查询字段。
- 上传项目只允许静态读取，任何情况下都不执行其中的脚本、二进制文件或构建命令。
- 文件、数据库记录和生成内容必须具备同一份 Report 生命周期。
- AI 失败不能导致服务崩溃；可生成有明确标识的 fallback 报告。

## 4. 技术选型

### 4.1 前端

- Vue 3
- Vite
- TypeScript
- Vue Router
- Pinia
- Axios
- Tailwind CSS
- lucide-vue-next
- markdown-it
- Vitest + Vue Test Utils
- Playwright，用于关键流程冒烟测试

前端不引入重型 UI 框架。通用按钮、输入框、上传区、Badge、空状态等组件自行封装，保持视觉系统轻量统一。

### 4.2 后端

- Go 1.22+
- Gin
- GORM
- `gorm.io/driver/mysql`
- godotenv
- google/uuid
- slog
- net/http
- archive/zip
- multipart/form-data

### 4.3 数据和存储

- MySQL 8.0+
- 默认字符集：`utf8mb4`
- 默认排序规则：`utf8mb4_0900_ai_ci`
- 时间统一存储为 UTC，API 使用 RFC 3339 输出
- 上传文件使用本地目录，数据库只保存受控相对路径和元数据
- Docker 使用命名卷持久化 MySQL、uploads 和 temp

## 5. 推荐目录结构

```text
AI Developer Workbench/
├─ backend/
│  ├─ cmd/
│  │  └─ server/
│  │     └─ main.go
│  ├─ internal/
│  │  ├─ config/
│  │  │  └─ config.go
│  │  ├─ database/
│  │  │  ├─ database.go
│  │  │  └─ migrate.go
│  │  ├─ model/
│  │  │  ├─ report.go
│  │  │  ├─ generated_file.go
│  │  │  └─ report_asset.go
│  │  ├─ dto/
│  │  │  ├─ response.go
│  │  │  ├─ report_dto.go
│  │  │  ├─ tool_dto.go
│  │  │  ├─ ui_review_dto.go
│  │  │  ├─ project_doctor_dto.go
│  │  │  ├─ agent_config_dto.go
│  │  │  ├─ api_doc_dto.go
│  │  │  └─ db_schema_dto.go
│  │  ├─ handler/
│  │  │  ├─ health_handler.go
│  │  │  ├─ system_handler.go
│  │  │  ├─ dashboard_handler.go
│  │  │  ├─ tool_handler.go
│  │  │  └─ report_handler.go
│  │  ├─ service/
│  │  │  ├─ ai_service.go
│  │  │  ├─ openai_compatible_service.go
│  │  │  ├─ mock_ai_service.go
│  │  │  ├─ file_service.go
│  │  │  ├─ zip_service.go
│  │  │  ├─ report_service.go
│  │  │  ├─ export_service.go
│  │  │  ├─ tools/
│  │  │  │  ├─ ui_review_service.go
│  │  │  │  ├─ project_doctor_service.go
│  │  │  │  ├─ agent_config_service.go
│  │  │  │  ├─ api_doc_service.go
│  │  │  │  └─ db_schema_service.go
│  │  │  └─ mocks/
│  │  │     ├─ ui_review.go
│  │  │     ├─ project_doctor.go
│  │  │     ├─ agent_config.go
│  │  │     ├─ api_doc.go
│  │  │     └─ db_schema.go
│  │  ├─ prompts/
│  │  │  ├─ prompt_builder.go
│  │  │  ├─ ui_review_prompt.go
│  │  │  ├─ project_doctor_prompt.go
│  │  │  ├─ agent_config_prompt.go
│  │  │  ├─ api_doc_prompt.go
│  │  │  └─ db_schema_prompt.go
│  │  ├─ repository/
│  │  │  ├─ report_repository.go
│  │  │  ├─ generated_file_repository.go
│  │  │  └─ report_asset_repository.go
│  │  ├─ middleware/
│  │  │  ├─ cors.go
│  │  │  ├─ recovery.go
│  │  │  └─ request_id.go
│  │  └─ util/
│  │     ├─ json_parser.go
│  │     ├─ response.go
│  │     ├─ score.go
│  │     ├─ text_truncate.go
│  │     ├─ filename.go
│  │     └─ secret_redactor.go
│  ├─ migrations/
│  │  └─ 0001_init.sql
│  ├─ uploads/
│  │  └─ .gitkeep
│  ├─ temp/
│  │  └─ .gitkeep
│  ├─ .env.example
│  ├─ Dockerfile
│  ├─ go.mod
│  └─ README.md
├─ frontend/
│  ├─ src/
│  │  ├─ api/
│  │  │  ├─ client.ts
│  │  │  ├─ tools.ts
│  │  │  ├─ reports.ts
│  │  │  └─ system.ts
│  │  ├─ components/
│  │  │  ├─ common/
│  │  │  ├─ dashboard/
│  │  │  ├─ layout/
│  │  │  ├─ report/
│  │  │  └─ tools/
│  │  ├─ pages/
│  │  │  ├─ tools/
│  │  │  ├─ DashboardPage.vue
│  │  │  ├─ LandingPage.vue
│  │  │  ├─ ReportDetailPage.vue
│  │  │  ├─ ReportsPage.vue
│  │  │  └─ SettingsPage.vue
│  │  ├─ router/
│  │  ├─ stores/
│  │  ├─ types/
│  │  ├─ utils/
│  │  ├─ styles/
│  │  ├─ App.vue
│  │  └─ main.ts
│  ├─ .env.example
│  ├─ Dockerfile
│  ├─ package.json
│  └─ vite.config.ts
├─ docker-compose.yml
├─ .env.example
├─ .gitignore
├─ DEVELOPMENT_PLAN.md
└─ README.md
```

## 6. MySQL 数据库设计

### 6.1 设计原则

- UUID 使用 `CHAR(36)`，优先保证开发可读性和与 API 的一致性。
- `input_json` 和 `report_json` 使用 MySQL 原生 `JSON` 类型。
- 报告列表所需字段独立列出，避免列表页反复解析大 JSON。
- 生成文件内容保存在数据库，上传的原始文件保存在受控磁盘目录。
- 外键使用 `ON DELETE CASCADE`，数据库记录删除与磁盘清理由 Service 事务协调。
- 不在数据库中保存 AI API Key。

### 6.2 reports

```sql
CREATE TABLE reports (
  id CHAR(36) NOT NULL,
  tool_type VARCHAR(32) NOT NULL,
  title VARCHAR(255) NOT NULL,
  input_mode VARCHAR(32) NOT NULL DEFAULT '',
  status VARCHAR(24) NOT NULL DEFAULT 'processing',
  summary TEXT NULL,
  total_score SMALLINT UNSIGNED NULL,
  grade VARCHAR(64) NULL,
  input_json JSON NULL,
  report_json JSON NOT NULL,
  file_path VARCHAR(1024) NULL,
  file_url VARCHAR(1024) NULL,
  error_message TEXT NULL,
  created_at DATETIME(3) NOT NULL,
  updated_at DATETIME(3) NOT NULL,
  PRIMARY KEY (id),
  KEY idx_reports_tool_created (tool_type, created_at),
  KEY idx_reports_status_created (status, created_at),
  KEY idx_reports_score_created (total_score, created_at),
  CONSTRAINT chk_reports_score
    CHECK (total_score IS NULL OR total_score BETWEEN 0 AND 100)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_0900_ai_ci;
```

建议状态值：

- `processing`：已创建，正在分析。
- `succeeded`：AI 或 Mock 正常生成。
- `fallback`：AI 调用或解析失败，已生成降级报告。
- `failed`：输入或持久化等不可恢复错误。

创建 `processing` 记录时，`report_json` 先写入空对象 `{}`，分析完成后再更新为最终报告，保证非空约束和任务状态一致。

### 6.3 generated_files

```sql
CREATE TABLE generated_files (
  id CHAR(36) NOT NULL,
  report_id CHAR(36) NOT NULL,
  filename VARCHAR(255) NOT NULL,
  language VARCHAR(32) NULL,
  mime_type VARCHAR(100) NOT NULL DEFAULT 'text/markdown',
  content LONGTEXT NOT NULL,
  size_bytes BIGINT UNSIGNED NOT NULL DEFAULT 0,
  created_at DATETIME(3) NOT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uk_generated_file_report_name (report_id, filename),
  KEY idx_generated_files_report (report_id),
  CONSTRAINT fk_generated_files_report
    FOREIGN KEY (report_id) REFERENCES reports(id)
    ON DELETE CASCADE
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_0900_ai_ci;
```

典型内容包括：

- Agent Config Studio 生成的 `AGENTS.md`、`TASK_PLAN.md`。
- API Doc Builder 生成的 `API_DOCUMENTATION.md`、`openapi.json`。
- 每份报告统一导出的 Markdown 文档。

### 6.4 report_assets

该表用于管理截图、项目 ZIP 等上传资产，解决单个 `file_path` 无法表达多个文件及文件生命周期的问题。

```sql
CREATE TABLE report_assets (
  id CHAR(36) NOT NULL,
  report_id CHAR(36) NOT NULL,
  asset_type VARCHAR(32) NOT NULL,
  original_name VARCHAR(255) NOT NULL,
  stored_name VARCHAR(255) NOT NULL,
  relative_path VARCHAR(1024) NOT NULL,
  mime_type VARCHAR(100) NULL,
  size_bytes BIGINT UNSIGNED NOT NULL,
  sha256 CHAR(64) NULL,
  created_at DATETIME(3) NOT NULL,
  PRIMARY KEY (id),
  KEY idx_report_assets_report (report_id),
  CONSTRAINT fk_report_assets_report
    FOREIGN KEY (report_id) REFERENCES reports(id)
    ON DELETE CASCADE
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_0900_ai_ci;
```

`asset_type` 建议值：

- `screenshot`
- `project_zip`
- `source_file`

### 6.5 迁移策略

- `backend/migrations/0001_init.sql` 保存可审查的基线结构。
- 开发环境可通过 `DB_AUTO_MIGRATE=true` 在启动时执行 GORM AutoMigrate，方便首次运行。
- Docker Compose 通过 MySQL 健康检查后再启动后端。
- 生产环境关闭 AutoMigrate，使用显式、可回滚的版本化迁移。
- 初始化失败时后端直接退出并输出不包含密码的错误信息，不能带病启动。

## 7. 后端模块设计

### 7.1 启动流程

`main.go` 按以下顺序完成初始化：

1. 读取并校验环境变量。
2. 初始化结构化日志。
3. 创建 uploads/temp 目录。
4. 连接 MySQL 并配置连接池。
5. 按配置执行迁移。
6. 初始化 Repository。
7. 初始化 File、ZIP、AI、Report 和五个 Tool Service。
8. 注册中间件和路由。
9. 启动带超时配置的 HTTP Server。
10. 监听终止信号并优雅关闭 HTTP 和数据库连接。

### 7.2 配置模型

```env
APP_ENV=development
APP_PORT=8080
APP_VERSION=0.1.0

DATABASE_DRIVER=mysql
DATABASE_HOST=127.0.0.1
DATABASE_PORT=3306
DATABASE_NAME=ai_workbench
DATABASE_USER=workbench
DATABASE_PASSWORD=change_me
DATABASE_CHARSET=utf8mb4
DATABASE_LOC=UTC
DATABASE_MAX_OPEN_CONNS=20
DATABASE_MAX_IDLE_CONNS=10
DATABASE_CONN_MAX_LIFETIME_MINUTES=30
DB_AUTO_MIGRATE=true

UPLOAD_DIR=./uploads
TEMP_DIR=./temp
MAX_UPLOAD_SIZE_MB=20
MAX_PROJECT_FILES=120
MAX_FILE_READ_BYTES=12000
MAX_PROJECT_TOTAL_READ_BYTES=300000
MAX_ZIP_UNCOMPRESSED_MB=100

AI_PROVIDER=openai
AI_BASE_URL=https://api.openai.com/v1
AI_API_KEY=
AI_MODEL=gpt-4.1
AI_VISION_MODEL=gpt-4.1
AI_MOCK_MODE=true
AI_TIMEOUT_SECONDS=90
AI_MAX_RETRIES=1

CORS_ALLOW_ORIGINS=http://localhost:5173
```

配置规则：

- `AI_API_KEY` 为空时强制使用 Mock Mode。
- Settings 接口只返回 provider、模型、Mock 状态和上传限制，不返回 Key、密码或 DSN。
- 日志禁止输出 `AI_API_KEY`、`DATABASE_PASSWORD` 和完整请求头。

### 7.3 Repository

Repository 只负责数据库操作，不包含 AI 或文件逻辑：

```go
type ReportRepository interface {
    Create(ctx context.Context, report *model.Report) error
    Update(ctx context.Context, report *model.Report) error
    GetByID(ctx context.Context, id string) (*model.Report, error)
    List(ctx context.Context, query ListReportsQuery) ([]model.Report, int64, error)
    Delete(ctx context.Context, tx *gorm.DB, id string) error
    GetDashboardStats(ctx context.Context) (*DashboardStats, error)
}
```

列表查询需要：

- `page` 默认 1。
- `page_size` 默认 10，最大 100。
- `tool_type` 使用白名单。
- `sort` 映射到固定 SQL 片段，禁止直接拼接用户输入。

### 7.4 Report Service

Report Service 负责：

- 创建 processing 报告。
- 保存结构化结果。
- 保存 generated files 和 assets。
- 查询列表和详情。
- 生成统一 Markdown 报告。
- 在事务内删除数据库记录。
- 数据库删除成功后删除磁盘文件。
- 磁盘删除失败时记录告警，不能回滚已经成功的数据库事务；后续可由清理任务补偿。

### 7.5 File Service

上传策略：

- 截图仅允许 PNG、JPEG、WebP。
- 项目文件仅允许 ZIP。
- 同时校验扩展名、Content-Type 和文件头，不能只相信浏览器传入的 MIME。
- 原始文件名只用于显示；磁盘文件名使用 UUID。
- 所有保存路径都通过 `filepath.Clean`、`filepath.Rel` 校验必须位于配置目录内。
- 计算 SHA-256 便于审计和后续去重。
- 上传失败时立即清理不完整文件。

目录建议：

```text
uploads/{report_id}/source/{uuid}.{ext}
temp/{report_id}/extracted/
```

### 7.6 ZIP Service

ZIP 扫描必须实现以下限制：

- 拒绝绝对路径、`..` 路径、盘符路径和解压后越界路径，防止 Zip Slip。
- 拒绝符号链接和其他非常规文件。
- 限制压缩包大小、文件数量、单文件读取字节数和总解压大小。
- 可增加压缩比阈值以识别 ZIP Bomb。
- 忽略 `.git`、`node_modules`、`dist`、`build`、`coverage`、`vendor`、`.idea`、`.vscode`。
- 忽略二进制、媒体、数据库、锁文件和超大文件。
- 优先读取 README、依赖清单、入口、路由、API、配置、测试和 Agent 上下文文件。
- 不读取 `.env`、私钥、证书和常见凭据文件。
- 对文本中的疑似 Token、密码和连接串执行脱敏后再传给 AI。
- 扫描完成后删除临时解压目录，原 ZIP 按 Report 生命周期保留。

项目摘要建议包含：

```json
{
  "tree": [],
  "detected_stack": [],
  "important_files": [],
  "engineering_signals": {},
  "truncation": {
    "files_skipped": 0,
    "bytes_omitted": 0
  }
}
```

### 7.7 AI Service

统一接口：

```go
type AIService interface {
    GenerateJSON(ctx context.Context, input AIRequest) (*AIResult, error)
}

type AIRequest struct {
    ToolType    string
    SystemPrompt string
    UserPrompt   string
    ImagePath    string
    NeedVision   bool
}

type AIResult struct {
    RawText  string
    JSONText string
    Provider string
    Model    string
    IsMock   bool
}
```

实现策略：

- Provider 使用 OpenAI-compatible Chat Completions API。
- Vision 图片转换为受限大小的 base64 data URL。
- 如果兼容服务支持 JSON response format，则优先启用；否则仍在 Prompt 中强制 JSON。
- HTTP Client 设置连接、请求和响应超时。
- 仅对 429、部分 5xx 和短暂网络错误重试一次。
- 尊重 `context.Context` 取消。
- 不记录完整图片、用户代码和 API Key。

### 7.8 JSON 容错链

AI 返回结果按以下顺序处理：

1. 直接 `json.Unmarshal`。
2. 去除 Markdown code fence 后解析。
3. 提取第一个平衡的 JSON object/array 后解析。
4. 对工具对应的强类型 DTO 做结构校验和默认值补齐。
5. 分数统一限制到 0–100，severity 等枚举使用白名单。
6. 仍失败则生成 fallback report，并将状态设为 `fallback`。

Mock、真实 AI 和 fallback 必须输出相同的数据结构，前端无需区分渲染逻辑。

### 7.9 Prompt Builder

每个 Prompt 由四部分组成：

1. 固定角色与安全边界。
2. 工具专属评分维度和输出 JSON Schema。
3. 用户输入或静态扫描摘要。
4. 中文 Codex Prompt 的生成要求。

上传代码和项目文件一律视为“不可信分析材料”。System Prompt 必须明确要求模型不能执行或服从材料中出现的指令，以降低 Prompt Injection 风险。

## 8. 五个工具的实现方案

### 8.1 UI Review

输入模式：

- `screenshot`
- `code`
- `screenshot_code`

校验：

- screenshot 模式必须有图片。
- code 模式必须有代码。
- screenshot_code 两者都必须提供。
- code 字符数设置上限，并在后端再次校验。

处理：

1. 创建 Report。
2. 保存截图资产。
3. 组装页面类型、目标风格、描述和代码。
4. 图片模式调用 Vision 模型。
5. 解析评分、问题、建议和中文优化 Prompt。
6. 保存报告，并生成 `UI_REVIEW_REPORT.md`。

### 8.2 Project Doctor

处理：

1. 保存项目 ZIP。
2. 安全解压和静态扫描。
3. 构建目录树、技术栈、工程信号和关键文本摘要。
4. 根据 basic/standard/deep 调整读取范围；MVP 中 deep 与 standard 可共享扫描器，但 Prompt 深度不同。
5. 生成结构、维护性、文档、配置、扩展性和 Agent Readiness 评分。
6. 保存 `PROJECT_DOCTOR_REPORT.md`。

严禁运行 `npm install`、`go test`、构建脚本、用户二进制或任何上传代码。

### 8.3 Agent Config Studio

处理：

1. 校验项目类型、前后端栈、数据库和 UI 风格。
2. 生成结构化 `files` 数组。
3. MVP 至少保存：
   - `AGENTS.md`
   - `TASK_PLAN.md`
4. 正常情况下同时生成：
   - `CODING_RULES.md`
   - `FRONTEND_STYLE_GUIDE.md`
   - `BACKEND_ARCHITECTURE.md`
   - `README_AGENT_CONTEXT.md`
5. 每个文件保存为 `generated_files`，支持独立复制和下载。

文件名必须使用服务端白名单，避免模型输出危险路径。

### 8.4 API Doc Builder

输入模式：

- `code`
- `project_zip`
- `manual`

处理：

- code：直接分析受限长度代码。
- project_zip：复用 ZIP Service，但优先提取路由、Handler、DTO、Model 和中间件。
- manual：根据 API 描述生成文档，作为 MVP 的轻量支持。

输出：

- 模块和端点结构。
- Markdown API 文档。
- 可选 OpenAPI JSON。
- 补全前端调用或后端文档的 Codex Prompt。

生成文件：

- `API_DOCUMENTATION.md`
- `openapi.json`，当 output_format 为 openapi 或 both 时生成。

### 8.5 DB Schema Review

输入支持 SQL、GORM、Prisma、自然语言表说明。

处理：

- 将用户声明的数据库类型作为分析上下文，而不是假设所有输入都是本项目 MySQL。
- 评分命名、范式、索引、关系、扩展性和数据完整性。
- 给出优化 Schema 和迁移建议。
- 当目标数据库是 MySQL 时，重点检查字符集、字段类型、索引长度、唯一约束、外键、时间字段和 JSON 使用方式。
- 所有 SQL 和迁移内容只按纯文本分析，服务端不连接用户数据库，也不执行任何 SQL。

生成 `DB_SCHEMA_REVIEW.md`；若 `target_goal=generate_migration`，可额外生成 `migration.sql`，但必须明确标记为“建议脚本，执行前人工审查”。

## 9. API 设计

统一前缀：`/api`

统一成功响应：

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

统一失败响应：

```json
{
  "code": 40001,
  "message": "invalid request",
  "error": "title is required",
  "request_id": "..."
}
```

### 9.1 基础接口

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/api/health` | 服务和数据库健康检查 |
| GET | `/api/tools` | 五个工具元数据 |
| GET | `/api/system/status` | Mock、Provider、模型、上传限制等公开状态 |
| GET | `/api/dashboard/stats` | 报告总数、各工具数量、最近报告、平均分 |

### 9.2 工具接口

| 方法 | 路径 | Content-Type |
| --- | --- | --- |
| POST | `/api/tools/ui-review/run` | multipart/form-data |
| POST | `/api/tools/project-doctor/run` | multipart/form-data |
| POST | `/api/tools/agent-config/run` | application/json |
| POST | `/api/tools/api-doc/run` | multipart/form-data |
| POST | `/api/tools/db-schema/run` | application/json |

MVP 可采用同步请求。前端需显示长任务加载状态，后端 AI 请求超时建议 90 秒。后续任务耗时明显增加时，再演进为异步 Job，不在首版提前引入消息队列。

### 9.3 报告接口

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/api/reports` | 分页、筛选和排序 |
| GET | `/api/reports/:id` | 报告详情和生成文件列表 |
| DELETE | `/api/reports/:id` | 删除报告、生成文件和上传资产 |
| GET | `/api/reports/:id/export?format=markdown` | 下载统一 Markdown 报告 |
| GET | `/api/reports/:id/files/:filename` | 下载某个生成文件 |

详情接口建议直接返回解析后的 `input_data` 和 `report_data` 对象，不让前端自行解析数据库 JSON 字符串。原始输入返回前必须去除磁盘路径、敏感字段和不必要的大段源代码。

### 9.4 错误码

```text
0       success
40001   invalid request
40002   upload failed
40003   unsupported file type
40004   invalid tool type
40005   file too large
40006   unsafe archive
40401   report not found
40901   report state conflict
50001   internal server error
50002   ai provider error
50003   database error
```

HTTP 状态码和业务错误码同时正确使用，例如参数错误返回 HTTP 400，报告不存在返回 HTTP 404。

## 10. 前端方案

### 10.1 路由

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

### 10.2 页面结构

- Landing Page：产品定位、五工具、工作流、示例报告和 CTA。
- Dashboard：五个工具入口、统计、最近报告和快捷操作。
- Tool Pages：统一 ToolHeader + 表单 + 输入区域 + 运行状态。
- Reports：工具筛选、排序、分页、删除和空状态。
- Report Detail：统一头部和摘要，按 tool_type 渲染专属评分与内容。
- Settings：只读展示后端 AI 状态、模型和上传限制。

### 10.3 状态管理

Pinia 只管理跨页面状态：

- `toolStore`：工具元数据。
- `reportStore`：报告列表条件、最近报告和缓存详情。
- `systemStore`：Mock 状态和运行限制。

工具表单局部状态保留在页面组件，不把所有输入堆到全局 Store。

### 10.4 报告类型

使用一个通用 Report 外壳和五类专属 `report_data` 类型：

```ts
interface Report<T = Record<string, unknown>> {
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

Agent Config 和 API Doc 不一定天然具有总分，前端必须允许 `total_score` 和 `grade` 为空，不能人为制造无意义评分。

### 10.5 UI 规范

- 背景、边框、文字和强调色沿用原 Prompt 的克制中性色。
- 内容最大宽度 1180px，Dashboard 可到 1280px。
- 主要信息依靠排版、分隔线和留白组织，不将每段内容都包成卡片。
- 阴影只用于浮层，常规面板以边框和背景层级为主。
- 动效限制为 150–200ms 的状态反馈。
- 移动端侧栏折叠，表单、报告评分和操作区改为单列。
- 所有表单字段具有 label、错误信息、键盘焦点和禁用状态。
- 颜色不是 severity 的唯一表达方式，同时显示文本和图标。

### 10.6 API 客户端

- `VITE_API_BASE_URL` 控制 API 地址。
- Axios 统一处理超时和错误响应。
- 上传请求由浏览器生成 multipart boundary，不手工写死。
- 文件下载读取 `Content-Disposition`，并通过 Blob 触发保存。
- 运行工具期间禁用重复提交。
- 失败时保留用户输入，便于修改后重试。

## 11. Mock Mode

Mock Mode 是完整功能，不是简短占位数据：

- 五个工具各有独立、稳定、内容充足的 mock report。
- Mock 输入仍需执行正常参数校验和文件安全校验。
- Mock 结果经过与真实 AI 相同的 DTO 校验和持久化流程。
- UI Review Mock 可以根据有无截图调整摘要文案。
- Agent Config Mock 必须真实生成可下载文件。
- API Doc Mock 必须包含多个模块和端点示例。
- Settings 页面显著但克制地展示当前为 Mock Mode。

这样可以在没有 API Key 时完成全产品验收，而不是只验证页面。

## 12. 安全方案

### 12.1 上传和 ZIP 安全

- Gin 层使用 `http.MaxBytesReader` 限制请求体。
- 保存前验证文件头、扩展名和 MIME。
- 文件路径使用服务端生成值，不接受用户传入目录。
- 解压时逐项校验目标路径必须位于临时目录。
- 限制解压文件数、总大小、目录深度、压缩比和读取总字节数。
- 拒绝 symlink、device file 和异常 ZIP entry。
- 永不执行上传内容。

### 12.2 敏感信息

- 不读取项目中的 `.env`、私钥和凭据文件。
- 对摘要中的 token、password、secret、DSN 做脱敏。
- 前端永远拿不到 AI Key 和数据库密码。
- 日志不记录 Authorization、Cookie、完整代码和图片 base64。

### 12.3 Web 和 API

- CORS 使用精确 allowlist，不使用生产环境通配符。
- 下载文件名经过清洗并设置安全的 `Content-Disposition`。
- Markdown 渲染关闭原始 HTML或使用严格净化，避免 XSS。
- 所有排序、枚举和工具类型均使用白名单。
- Gin Recovery 返回统一错误，不向前端暴露堆栈和绝对路径。

### 12.4 部署边界

MVP 没有用户鉴权。部署说明必须明确：

- 推荐仅本机使用，或置于带认证的反向代理之后。
- 不可直接公开到公网。
- 若进入公网或多人使用阶段，登录、权限、限流、审计和对象存储必须先于新增业务功能实施。

## 13. Docker Compose 方案

服务组成：

- `mysql`：MySQL 数据库。
- `backend`：Gin API，8080 端口。
- `frontend`：Vue 应用，5173 端口。

持久化卷：

- `mysql_data:/var/lib/mysql`
- `backend_uploads:/app/uploads`
- `backend_temp:/app/temp`

依赖关系：

- MySQL 配置 `mysqladmin ping` 健康检查。
- backend 在 MySQL healthy 后启动。
- frontend 依赖 backend 可访问，但页面自身不应因后端暂时不可用而白屏。

根目录 `.env.example` 提供 Compose 所需的数据库名、用户和密码示例，真实 `.env` 加入 `.gitignore`。

## 14. 测试策略

### 14.1 后端单元测试

重点覆盖：

- JSON code fence 和异常文本提取。
- 分数及枚举归一化。
- 文件名清洗和路径越界防护。
- ZIP Slip、文件数量、解压大小和忽略目录。
- 敏感信息脱敏。
- 五个 Prompt Builder 的必要字段。
- 五个 Mock 报告的结构校验。

### 14.2 后端集成测试

- 使用测试 MySQL 验证迁移、CRUD、分页、排序和级联删除。
- Handler 使用 fake AI Service 验证统一响应。
- 测试一次完整的工具运行、报告持久化、文件下载和删除流程。
- 验证 AI 超时、非法 JSON 和数据库异常时的降级行为。

### 14.3 前端测试

- 表单模式切换与必填校验。
- 上传文件提示和错误状态。
- 报告列表筛选、分页和删除确认。
- 五类报告详情渲染。
- Copy 和 Download 操作。
- Mock 状态和后端不可用状态。

### 14.4 E2E 冒烟测试

在 Mock Mode 下至少跑通：

1. 打开 Landing 和 Dashboard。
2. 分别运行五个工具。
3. 进入生成的报告详情。
4. 复制 Codex Prompt。
5. 下载 Markdown 或生成文件。
6. 在历史报告中找到并删除报告。

## 15. 开发阶段与交付物

### 阶段 1：工程骨架和 MySQL 基础

交付：

- 根目录、backend、frontend 基础结构。
- Gin 启动、配置加载、日志和优雅关闭。
- MySQL 连接、连接池、健康检查和 `0001_init.sql`。
- 统一响应、错误码和 request ID。
- Vue/Vite/Tailwind/Router/Pinia 基础工程。

验收：

- MySQL、backend、frontend 可通过 Docker Compose 启动。
- `/api/health` 能验证服务和数据库状态。

### 阶段 2：报告、文件和 AI 公共能力

交付：

- Report、GeneratedFile、ReportAsset Model 和 Repository。
- 报告列表、详情、删除、导出。
- File Service、ZIP Service 和安全限制。
- AI Service、Mock Mode、JSON 容错和 fallback。
- Dashboard stats 和 system status。

验收：

- 可通过测试数据完成报告 CRUD。
- 恶意 ZIP 和超限文件被拒绝。
- 无 API Key 时自动进入 Mock Mode。

### 阶段 3：五个后端工具纵向打通

按以下顺序实现：

1. Agent Config Studio：无文件/视觉依赖，优先验证生成文件链路。
2. DB Schema Review：验证 JSON 报告和 Markdown 导出。
3. UI Review：接入截图和 Vision。
4. Project Doctor：接入 ZIP 扫描。
5. API Doc Builder：复用 ZIP 扫描和生成文件。

每完成一个工具即补齐 DTO、Prompt、Mock、Service、Handler 和测试，不先堆完所有空壳。

### 阶段 4：前端产品框架和报告中心

交付：

- Landing、AppShell、Sidebar、Header、Dashboard。
- Reports、Report Detail、Settings。
- API Client、状态管理、通用加载/错误/空状态。
- 报告复制、Markdown 渲染和下载。

### 阶段 5：五个工具页面

交付：

- 五个工具输入页面。
- 文件上传、CodeInput、OptionSelect 等通用组件。
- 完整校验、提交状态和成功跳转。
- 最近报告入口。

### 阶段 6：质量收口

交付：

- 响应式和可访问性优化。
- Mock Mode E2E。
- Dockerfile、docker-compose.yml。
- 根 README、后端 README、前端 README。
- `.env.example` 和安全部署说明。
- 最终验收记录。

## 16. MVP 完成定义

项目只有在以下条件全部满足时才算完成：

- 后端可以连接 MySQL 并自动完成首次初始化。
- 前端和后端可本地启动，也可通过 Docker Compose 启动。
- 首页、Dashboard、Reports、Report Detail、Settings 正常工作。
- 五个工具全部能在 Mock Mode 下完成真实流程。
- 配置 API Key 后可以使用 OpenAI-compatible Provider。
- UI Review 支持截图、代码和组合模式。
- Project Doctor 和 API Doc Builder 可以安全分析 ZIP。
- Agent Config Studio 至少生成 `AGENTS.md` 和 `TASK_PLAN.md`。
- DB Schema Review 可以审查 MySQL SQL，也支持其他声明的数据库类型。
- 报告写入 MySQL，支持分页、筛选、排序、详情和删除。
- Codex Prompt 可以复制。
- Markdown 报告和生成文件可以下载。
- 删除报告会清理数据库关联记录和上传文件。
- AI 非法 JSON、超时或 Provider 错误时能够返回 fallback 报告。
- API Key、数据库密码和上传项目中的敏感信息不会泄漏到前端或日志。
- ZIP 分析永不执行用户代码。
- 关键单元测试、集成测试和 Mock E2E 通过。
- README 包含项目介绍、技术栈、启动方式、MySQL 配置、Mock Mode、API、上传限制、安全边界和 Roadmap。

## 17. 主要风险和控制方式

| 风险 | 影响 | 控制方式 |
| --- | --- | --- |
| AI 返回结构不稳定 | 报告页无法渲染 | 强约束 Prompt、JSON 容错、强类型校验、fallback |
| Vision/大项目输入过长 | 超时或费用过高 | 图片大小限制、文本预算、摘要化、单次重试 |
| 恶意或异常 ZIP | 路径穿越、磁盘耗尽 | Zip Slip 防护、数量/大小/压缩比限制、临时目录隔离 |
| MySQL 与磁盘删除不一致 | 产生孤儿文件 | 数据库事务、删除补偿日志、后续清理任务接口 |
| 无鉴权服务公开部署 | 数据泄露或资源滥用 | 明确本地定位，公网必须加认证代理；后续补鉴权和限流 |
| 通用报告模型过度抽象 | 各工具信息丢失 | 通用外壳 + 工具专属 JSON DTO，不强行统一内部字段 |
| 前端自行解析原始 JSON | 类型混乱、兼容困难 | 后端详情 DTO 返回解析后的对象 |

## 18. 后续演进

MVP 稳定后，推荐按以下顺序扩展：

1. 异步任务、进度查询和任务取消。
2. 登录、项目空间、权限和审计日志。
3. S3 兼容对象存储，替换本地上传目录。
4. GitHub 仓库连接、PR Review 和 Issue 生成。
5. 自定义规则模板和多模型对比。
6. 报告趋势、PDF 导出和团队协作。
7. VS Code、浏览器和 Figma 插件。
8. 套餐、用量计费和订阅。

## 19. 推荐实施结论

本项目应采用“公共能力先行、工具纵向交付”的方式实施。第一优先级不是一次性铺满所有页面，而是先打通：

```text
Agent Config 输入
  -> Mock/AI
  -> 结构化结果
  -> MySQL
  -> 报告详情
  -> 文件下载
```

这条链路稳定后，DB Schema、UI Review、Project Doctor 和 API Doc Builder 都可以复用同一套报告、AI、文件和前端展示能力。该顺序能够尽早验证架构中风险最高的部分，同时保证每个开发阶段都有可运行、可验收的产品增量。
