# AI Developer Workbench 后端开发方案

> 版本：MVP 0.1.0  
> 技术栈：Go 1.22+ + Gin + GORM + MySQL 8.0+  
> 对应总方案：[DEVELOPMENT_PLAN.md](./DEVELOPMENT_PLAN.md)  
> 前端契约：[FRONTEND_DEVELOPMENT_PLAN.md](./FRONTEND_DEVELOPMENT_PLAN.md)

## 1. 后端目标

后端负责五个 AI 工具的业务编排、AI 调用、文件安全处理、报告持久化和统一 API。

五个核心工具：

1. UI Review
2. Project Doctor
3. Agent Config Studio
4. API Doc Builder
5. DB Schema Review

统一处理链路：

```text
接收请求
  -> 参数与文件校验
  -> 创建 processing 报告
  -> 构建工具上下文
  -> Mock 或真实 AI 分析
  -> JSON 容错与强类型校验
  -> 生成报告和附属文件
  -> MySQL 持久化
  -> 返回报告详情
```

后端必须保证：

- AI API Key 永远不返回前端。
- 没有 API Key 时自动启用 Mock Mode。
- AI 失败或返回非法 JSON 时服务不崩溃。
- 上传项目只做静态读取，不执行任何用户代码。
- 用户提交的 SQL 只做文本分析，不连接用户数据库、不执行 SQL。
- 报告、生成文件和上传资产具有统一生命周期。

## 2. 技术选型

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
- MySQL 8.0+

测试：

- Go 标准库 `testing`
- `httptest`
- Testcontainers 或独立测试 MySQL
- Fake AI Service

## 3. 分层架构

```text
Router
  -> Middleware
  -> Handler
  -> Service
     -> Tool Service
     -> AI Service
     -> File/ZIP Service
     -> Report/Export Service
  -> Repository
  -> MySQL / Local Storage / AI Provider
```

职责边界：

- Router：路由注册。
- Middleware：CORS、Recovery、Request ID、请求日志。
- Handler：请求绑定、基础校验、统一响应。
- Tool Service：五个工具业务编排。
- AI Service：Mock 和 OpenAI-compatible Provider。
- File/ZIP Service：上传、解压、静态扫描和清理。
- Report Service：报告事务、查询、导出和删除。
- Repository：只处理数据库访问。
- Prompt Builder：只负责构建 Prompt，不发请求、不访问数据库。

Handler 不得直接调用 GORM 或 AI Provider。

## 4. 推荐目录结构

```text
backend/
├─ cmd/
│  └─ server/
│     └─ main.go
├─ internal/
│  ├─ config/
│  │  └─ config.go
│  ├─ database/
│  │  ├─ database.go
│  │  └─ migrate.go
│  ├─ model/
│  │  ├─ report.go
│  │  ├─ generated_file.go
│  │  └─ report_asset.go
│  ├─ dto/
│  │  ├─ response.go
│  │  ├─ report_dto.go
│  │  ├─ tool_dto.go
│  │  ├─ ui_review_dto.go
│  │  ├─ project_doctor_dto.go
│  │  ├─ agent_config_dto.go
│  │  ├─ api_doc_dto.go
│  │  └─ db_schema_dto.go
│  ├─ handler/
│  │  ├─ health_handler.go
│  │  ├─ system_handler.go
│  │  ├─ dashboard_handler.go
│  │  ├─ tool_handler.go
│  │  └─ report_handler.go
│  ├─ service/
│  │  ├─ ai_service.go
│  │  ├─ openai_compatible_service.go
│  │  ├─ mock_ai_service.go
│  │  ├─ file_service.go
│  │  ├─ zip_service.go
│  │  ├─ report_service.go
│  │  ├─ export_service.go
│  │  ├─ tools/
│  │  │  ├─ ui_review_service.go
│  │  │  ├─ project_doctor_service.go
│  │  │  ├─ agent_config_service.go
│  │  │  ├─ api_doc_service.go
│  │  │  └─ db_schema_service.go
│  │  └─ mocks/
│  │     ├─ ui_review.go
│  │     ├─ project_doctor.go
│  │     ├─ agent_config.go
│  │     ├─ api_doc.go
│  │     └─ db_schema.go
│  ├─ prompts/
│  │  ├─ prompt_builder.go
│  │  ├─ ui_review_prompt.go
│  │  ├─ project_doctor_prompt.go
│  │  ├─ agent_config_prompt.go
│  │  ├─ api_doc_prompt.go
│  │  └─ db_schema_prompt.go
│  ├─ repository/
│  │  ├─ report_repository.go
│  │  ├─ generated_file_repository.go
│  │  └─ report_asset_repository.go
│  ├─ middleware/
│  │  ├─ cors.go
│  │  ├─ recovery.go
│  │  └─ request_id.go
│  └─ util/
│     ├─ json_parser.go
│     ├─ response.go
│     ├─ score.go
│     ├─ text_truncate.go
│     ├─ filename.go
│     └─ secret_redactor.go
├─ migrations/
│  └─ 0001_init.sql
├─ uploads/
│  └─ .gitkeep
├─ temp/
│  └─ .gitkeep
├─ .env.example
├─ Dockerfile
├─ go.mod
└─ README.md
```

## 5. 配置设计

`backend/.env.example`：

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

- 缺少必需数据库配置时启动失败。
- `AI_API_KEY` 为空时强制 Mock Mode。
- `AI_MOCK_MODE=true` 时禁止真实 AI 请求。
- CORS 来源按逗号拆分为精确 allowlist。
- 日志不得输出 `AI_API_KEY`、`DATABASE_PASSWORD` 和完整 DSN。
- Settings 接口只返回公开状态。

## 6. 启动流程

`main.go`：

1. 加载 `.env` 和环境变量。
2. 校验配置。
3. 初始化 slog。
4. 创建 uploads/temp 目录。
5. 连接 MySQL。
6. 配置连接池。
7. 按配置执行迁移。
8. 初始化 Repository。
9. 初始化公共 Service。
10. 初始化五个 Tool Service。
11. 注册 Middleware、Handler 和 Router。
12. 启动带超时的 HTTP Server。
13. 监听系统信号并优雅关闭。

HTTP Server 建议设置：

- Read Header Timeout。
- Read Timeout。
- Write Timeout 应高于 AI 请求超时。
- Idle Timeout。
- Max Header Bytes。

## 7. MySQL 数据库设计

### 7.1 基础规则

- MySQL 8.0+。
- InnoDB。
- 字符集 `utf8mb4`。
- 排序规则 `utf8mb4_0900_ai_ci`。
- 时间以 UTC 存储。
- API 时间使用 RFC 3339。
- UUID 使用 `CHAR(36)`。
- `input_json` 和 `report_json` 使用 MySQL `JSON`。

### 7.2 reports

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

状态：

- `processing`
- `succeeded`
- `fallback`
- `failed`

创建 processing 记录时将 `report_json` 写为 `{}`。

### 7.3 generated_files

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

### 7.4 report_assets

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

`asset_type`：

- `screenshot`
- `project_zip`
- `source_file`

### 7.5 迁移

- `migrations/0001_init.sql` 保存基线结构。
- 开发环境可使用 `DB_AUTO_MIGRATE=true`。
- 生产环境关闭 AutoMigrate，执行版本化 SQL。
- 迁移失败必须阻止服务启动。
- 数据库健康检查成功后才允许 API 对外提供服务。

## 8. Model 和 DTO

Model 对应数据库结构，DTO 对应 API 请求与响应，两者不得混用。

通用响应 DTO：

```go
type Response struct {
    Code      int    `json:"code"`
    Message   string `json:"message"`
    Data      any    `json:"data,omitempty"`
    Error     string `json:"error,omitempty"`
    RequestID string `json:"request_id,omitempty"`
}
```

报告详情 DTO 应返回：

- 基础报告字段。
- 已解析的 `input_data`。
- 已解析的 `report_data`。
- 生成文件元数据。

不得直接向前端返回：

- 本地绝对路径。
- AI API Key。
- 数据库密码。
- 完整 DSN。
- 未脱敏的上传项目敏感信息。
- 不必要的大段原始代码。

工具专属 DTO：

- UI Review Request/Result。
- Project Doctor Request/Result。
- Agent Config Request/Result。
- API Doc Request/Result。
- DB Schema Request/Result。

## 9. Repository 设计

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

要求：

- 所有方法接收 `context.Context`。
- 不在 Repository 中构建业务报告。
- 分页默认 page=1、page_size=10，最大 100。
- `tool_type` 使用白名单。
- 排序字段映射为固定 SQL，不能直接拼接用户输入。
- 详情查询预加载 generated files 和 assets。
- 数据库错误转换为可识别的内部错误，不向外暴露 SQL。

## 10. Report Service

职责：

- 创建 processing 报告。
- 更新 succeeded、fallback 或 failed 状态。
- 保存报告 JSON。
- 保存生成文件和上传资产。
- 查询报告列表和详情。
- 生成 Markdown 报告。
- 删除报告。

删除流程：

1. 查询需要删除的资产相对路径。
2. 开启 MySQL 事务。
3. 删除 Report，外键级联删除 generated files 和 assets。
4. 提交事务。
5. 删除磁盘文件和空目录。
6. 磁盘清理失败时记录告警，供后续补偿任务处理。

不要在数据库事务中执行耗时文件 I/O。

## 11. File Service

### 11.1 上传校验

截图：

- PNG
- JPEG
- WebP

项目文件：

- ZIP

校验维度：

- 请求体总大小。
- 文件扩展名。
- Content-Type。
- 文件头。
- 文件实际大小。

原则：

- 原始文件名只用于展示。
- 存储名使用 UUID。
- 所有路径必须验证位于配置的根目录内。
- 计算 SHA-256。
- 失败时删除不完整文件。

建议路径：

```text
uploads/{report_id}/source/{uuid}.{ext}
temp/{report_id}/extracted/
```

### 11.2 文件名安全

- 移除路径分隔符。
- 拒绝 `.`、`..` 和空文件名。
- 下载时只允许访问数据库记录中的文件。
- Agent 生成文件名使用服务端白名单或严格清洗。

## 12. ZIP Service

ZIP 只做静态分析。

必须防护：

- Zip Slip。
- ZIP Bomb。
- 符号链接。
- Device file。
- 绝对路径。
- Windows 盘符路径。
- `..` 路径。
- 文件数超限。
- 解压总大小超限。
- 单文件读取超限。
- 总文本读取预算超限。

忽略目录：

```text
node_modules
.git
dist
build
coverage
vendor
.idea
.vscode
```

优先读取：

```text
README.md
package.json
go.mod
requirements.txt
pyproject.toml
Dockerfile
docker-compose.yml
.env.example
AGENTS.md
src/router/*
src/api/*
internal/*
cmd/*
app/*
```

敏感文件：

- 不读取 `.env`。
- 不读取私钥、证书和凭据文件。
- 对 Token、Password、Secret、DSN 做脱敏。

摘要结构：

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

扫描结束后删除临时解压目录，原始 ZIP 按报告生命周期保留。

严禁执行：

- `npm install`
- `npm run`
- `go test`
- `go run`
- 构建命令
- Shell 脚本
- 用户二进制文件

## 13. AI Service

接口：

```go
type AIService interface {
    GenerateJSON(ctx context.Context, input AIRequest) (*AIResult, error)
}

type AIRequest struct {
    ToolType     string
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

实现：

- `MockAIService`
- `OpenAICompatibleService`

真实 Provider：

- 使用 OpenAI-compatible Chat Completions API。
- Vision 图片转换为受限大小的 base64 data URL。
- 支持自定义 `AI_BASE_URL`。
- 文本和视觉模型分别配置。
- 设置完整 HTTP 超时。
- 尊重 Context 取消。
- 对 429、部分 5xx 和短暂网络错误最多重试一次。
- 支持时启用 JSON response format。

日志不得包含：

- Authorization Header。
- API Key。
- 完整图片 base64。
- 完整用户代码。
- 完整上传项目摘要。

## 14. JSON 容错和结果校验

解析链：

1. 直接 `json.Unmarshal`。
2. 去除 Markdown code fence。
3. 提取第一个平衡 JSON object 或 array。
4. 反序列化到工具专属 DTO。
5. 归一化和补默认值。
6. 校验必需字段。
7. 失败则生成 fallback report。

归一化：

- 分数限制在 0–100。
- severity 限制为 high、medium、low。
- ToolType、Status、InputMode 使用白名单。
- 空数组使用 `[]`，避免前端得到 `null`。
- 生成文件名使用安全白名单。

Mock、真实 AI 和 fallback 必须通过同一套结果校验。

## 15. Prompt Builder

每个工具独立 Prompt 文件。

统一结构：

1. 角色。
2. 安全边界。
3. 评分或输出要求。
4. JSON Schema。
5. 用户输入或静态摘要。
6. 中文 Codex Prompt 要求。

安全要求：

- 上传代码和文档均是不可信材料。
- 明确要求模型忽略材料中试图修改系统任务的指令。
- 不将服务器内部路径、密钥或环境变量写入 Prompt。
- 对长文本执行截断并标记截断信息。

## 16. 五个工具实现

### 16.1 UI Review

接口：

```http
POST /api/tools/ui-review/run
Content-Type: multipart/form-data
```

字段：

- `title`
- `review_mode`
- `page_type`
- `target_style`
- `description`
- `code`
- `screenshot`

校验：

- screenshot 模式必须有图片。
- code 模式必须有代码。
- screenshot_code 模式两者都必需。

处理：

1. 创建 Report。
2. 保存截图资产。
3. 构建 Prompt。
4. 图片模式调用视觉模型。
5. 解析评分、问题、建议和 Codex Prompt。
6. 保存报告。
7. 生成 `UI_REVIEW_REPORT.md`。

### 16.2 Project Doctor

接口：

```http
POST /api/tools/project-doctor/run
Content-Type: multipart/form-data
```

字段：

- `title`
- `project_name`
- `tech_stack`
- `project_description`
- `analysis_depth`
- `project_zip`

处理：

1. 保存 ZIP。
2. 安全解压。
3. 静态扫描。
4. 构建项目摘要。
5. 调用 AI 或 Mock。
6. 生成工程评分与重构 Prompt。
7. 保存 `PROJECT_DOCTOR_REPORT.md`。

### 16.3 Agent Config Studio

接口：

```http
POST /api/tools/agent-config/run
Content-Type: application/json
```

输入：

- `title`
- `project_name`
- `project_type`
- `frontend_stack`
- `backend_stack`
- `database`
- `ui_style`
- `coding_preferences`
- `strict_rules`

至少生成：

- `AGENTS.md`
- `TASK_PLAN.md`

正常情况下同时生成：

- `CODING_RULES.md`
- `FRONTEND_STYLE_GUIDE.md`
- `BACKEND_ARCHITECTURE.md`
- `README_AGENT_CONTEXT.md`

所有文件保存到 `generated_files`。

### 16.4 API Doc Builder

接口：

```http
POST /api/tools/api-doc/run
Content-Type: multipart/form-data
```

输入：

- `title`
- `source_type`
- `backend_stack`
- `code`
- `api_description`
- `output_format`
- `project_zip`

模式：

- code：分析文本代码。
- project_zip：复用 ZIP Service，优先路由、Handler、DTO、Model。
- manual：分析手动 API 描述。

生成：

- `API_DOCUMENTATION.md`
- `openapi.json`，当格式为 openapi 或 both。

### 16.5 DB Schema Review

接口：

```http
POST /api/tools/db-schema/run
Content-Type: application/json
```

输入：

- `title`
- `schema_type`
- `database_type`
- `business_context`
- `schema_content`
- `target_goal`

处理：

- 只做纯文本分析。
- 不连接用户数据库。
- 不执行 SQL 或 migration。
- 支持 SQL、GORM、Prisma 和自然语言描述。
- 根据用户声明的数据库类型提供针对性建议。

生成：

- `DB_SCHEMA_REVIEW.md`
- 可选 `migration.sql`

`migration.sql` 必须注明执行前人工审查。

## 17. Mock Mode

Mock Mode 必须完成真实业务链路：

- 执行正常参数校验。
- 执行文件安全校验。
- 创建报告。
- 通过相同 DTO 校验。
- 写入 MySQL。
- 生成可下载文件。
- 支持查询和删除。

五个工具分别维护完整 Mock：

- UI Review Mock。
- Project Doctor Mock。
- Agent Config Mock。
- API Doc Mock。
- DB Schema Mock。

Mock 内容不能是几行占位文本，必须能够验证完整报告 UI。

## 18. API 设计

统一前缀：

```text
/api
```

### 18.1 基础接口

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/api/health` | 服务与数据库健康检查 |
| GET | `/api/tools` | 工具元数据 |
| GET | `/api/system/status` | 公开运行状态 |
| GET | `/api/dashboard/stats` | Dashboard 统计 |

`system/status` 只返回：

- Mock 状态。
- Provider 名称。
- 文本模型。
- 视觉模型。
- 上传和 ZIP 限制。

### 18.2 报告接口

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/api/reports` | 分页、筛选、排序 |
| GET | `/api/reports/:id` | 报告详情 |
| DELETE | `/api/reports/:id` | 删除报告及关联内容 |
| GET | `/api/reports/:id/export?format=markdown` | 导出 Markdown |
| GET | `/api/reports/:id/files/:filename` | 下载生成文件 |

列表参数：

- `page`
- `page_size`
- `tool_type`
- `sort=newest|oldest|score_desc|score_asc`

### 18.3 统一响应

成功：

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

失败：

```json
{
  "code": 40001,
  "message": "invalid request",
  "error": "title is required",
  "request_id": "..."
}
```

错误码：

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

HTTP 状态码必须与错误语义一致。

## 19. Middleware

### 19.1 Request ID

- 接受合法的上游 Request ID 或生成新 UUID。
- 写入 Context、响应头和错误响应。
- 所有日志包含 Request ID。

### 19.2 Recovery

- 捕获 panic。
- 记录内部堆栈。
- 对外返回统一 500。
- 不泄露绝对路径和堆栈。

### 19.3 CORS

- 精确 allowlist。
- 支持开发环境 `http://localhost:5173`。
- 生产环境不使用 `*`。
- 只开放实际使用的方法和请求头。

### 19.4 请求日志

记录：

- Request ID。
- 方法。
- 路径。
- 状态码。
- 耗时。

不记录：

- Authorization。
- Cookie。
- API Key。
- 上传文件内容。
- 完整 Prompt。

## 20. 安全方案

### 20.1 文件与 ZIP

- `http.MaxBytesReader` 限制请求体。
- 三重验证扩展名、MIME 和文件头。
- 路径必须处于受控目录。
- 防止 Zip Slip 和 ZIP Bomb。
- 拒绝符号链接和特殊文件。
- 不执行上传内容。

### 20.2 敏感信息

- 不读取 `.env`、私钥和凭据。
- 对摘要执行 Secret 脱敏。
- API 不返回本地绝对路径。
- 日志不打印 Key、密码和 DSN。

### 20.3 数据库和查询

- 所有数据库操作使用 GORM 参数绑定。
- 排序使用固定映射。
- 工具类型和枚举使用白名单。
- 删除使用事务和外键级联。
- 用户提交 SQL 不进入本项目 MySQL 执行。

### 20.4 部署边界

MVP 没有鉴权：

- 推荐仅本机或可信内网使用。
- 不可直接暴露公网。
- 公网部署必须先增加认证、限流、审计和安全反向代理。

## 21. Docker Compose 协作

根目录 Compose 包含：

- `mysql`
- `backend`
- `frontend`

后端依赖：

- MySQL `mysqladmin ping` 健康检查成功后启动。
- `mysql_data` 持久化数据库。
- `backend_uploads` 持久化上传文件。
- `backend_temp` 挂载临时目录。

后端容器：

- 监听 8080。
- 使用非 root 用户运行。
- 不将 `.env` 和 Secret 打包进镜像。
- uploads/temp 目录具有正确写权限。

## 22. 测试方案

### 22.1 单元测试

- 配置校验。
- JSON 容错提取。
- 分数和枚举归一化。
- 文件名清洗。
- 路径越界。
- ZIP Slip。
- ZIP 数量和大小限制。
- 忽略目录。
- Secret 脱敏。
- 五个 Prompt Builder。
- 五个 Mock 结构。
- Markdown 导出。

### 22.2 Repository 集成测试

使用测试 MySQL 验证：

- Migration。
- Create、Update、Get、List、Delete。
- 分页、筛选、排序。
- JSON 字段。
- 外键级联。
- Dashboard 统计。

### 22.3 Handler 和 Service 测试

- 使用 Fake AI Service。
- 参数错误。
- 文件过大和格式错误。
- 不安全 ZIP。
- AI 超时。
- AI 非法 JSON。
- fallback 报告。
- 文件下载。
- 报告删除与资产清理。

### 22.4 完整流程测试

Mock Mode 下：

1. 分别运行五个工具。
2. 验证报告写入 MySQL。
3. 获取列表和详情。
4. 下载报告和生成文件。
5. 删除报告。
6. 验证数据库级联和磁盘清理。

## 23. 开发顺序

### 阶段 1：工程和 MySQL

- 初始化 Go Module。
- 配置 Gin、Config、slog 和优雅关闭。
- 建立 MySQL 连接和连接池。
- 编写 `0001_init.sql`。
- 实现健康检查和统一响应。

### 阶段 2：报告基础能力

- Model、DTO 和 Repository。
- Report Service。
- 报告列表、详情、删除、导出。
- Dashboard stats 和 system status。

### 阶段 3：文件和 AI 公共能力

- File Service。
- ZIP Service。
- Secret 脱敏。
- AI Service。
- Mock Mode。
- JSON 容错与 fallback。

### 阶段 4：五个工具纵向交付

建议顺序：

1. Agent Config Studio
2. DB Schema Review
3. UI Review
4. Project Doctor
5. API Doc Builder

每完成一个工具同时完成：

- Request DTO。
- Result DTO。
- Prompt。
- Mock。
- Service。
- Handler。
- 生成文件。
- 测试。

### 阶段 5：部署与质量收口

- CORS、Recovery、Request ID。
- Dockerfile 和 Compose 联调。
- 性能和超时边界。
- 集成测试。
- 后端 README。
- 安全自检和最终验收。

## 24. 后端验收标准

- `go run ./cmd/server` 可以启动。
- 能连接 MySQL 8.0+。
- 首次启动可以完成数据库初始化。
- `/api/health` 正确报告服务和数据库状态。
- 五个工具接口全部可用。
- 无 API Key 时自动启用 Mock Mode。
- 配置 API Key 后可以调用 OpenAI-compatible Provider。
- AI 超时或非法 JSON 时生成 fallback 报告。
- 报告、生成文件和资产正确写入 MySQL。
- 报告支持分页、筛选、排序、详情和删除。
- Markdown 报告及生成文件可下载。
- 删除报告会级联删除数据库记录并清理磁盘资产。
- UI Review 支持 Vision 图片。
- Project Doctor 和 API Doc Builder 可安全扫描 ZIP。
- Agent Config 至少生成 `AGENTS.md` 和 `TASK_PLAN.md`。
- DB Schema Review 不连接用户数据库、不执行 SQL。
- 上传代码永远不会被执行。
- API Key、数据库密码和本地绝对路径不会泄露。
- 单元测试、MySQL 集成测试和 Mock 完整流程测试通过。

