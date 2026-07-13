# AI Developer Workbench Backend

Go backend for the AI Developer Workbench - a developer quality workbench for AI coding practitioners.

## Tech Stack

- **Language**: Go 1.24+
- **Web Framework**: Gin
- **ORM**: GORM with MySQL driver
- **Database**: MySQL 8.0+
- **AI Provider**: OpenAI-compatible Chat Completions API

## Features

### Five Core Tools

1. **UI Review** - AI-powered UI/UX review with screenshot or code analysis
2. **Project Doctor** - Comprehensive project health check and engineering quality analysis
3. **Agent Config Studio** - Generate AI agent configuration files (AGENTS.md, TASK_PLAN.md, etc.)
4. **API Doc Builder** - Auto-generate API documentation from code or descriptions
5. **DB Schema Review** - Review and optimize database schema with AI analysis

### Key Features

- Unified API response envelope: `{ code, message, data }`
- Mock Mode for development without AI API key
- Secure file upload with triple validation (extension, MIME, magic bytes)
- ZIP bomb and path traversal protection
- JSON fault tolerance with multiple parsing strategies
- Graceful shutdown with connection draining

## Quick Start

### Prerequisites

- Go 1.24+
- MySQL 8.0+

### Setup

1. Copy `.env.example` to `.env` and configure:

```bash
cp .env.example .env
```

2. Create MySQL database:

```sql
CREATE DATABASE ai_workbench;
CREATE USER 'workbench'@'localhost' IDENTIFIED BY 'your_password';
GRANT ALL PRIVILEGES ON ai_workbench.* TO 'workbench'@'localhost';
```

3. Run the server:

```bash
go run ./cmd/server
```

### Docker

```bash
# 在仓库根目录执行
docker compose up -d --build        # 构建并启动 MySQL + backend + frontend
docker compose logs -f backend      # 查看后端日志
docker compose ps                   # 查看健康状态
docker compose down                 # 停止
docker compose down -v              # 停止并清理数据卷（MySQL 数据丢失）
```

后端镜像健康检查命中 `GET /api/health`；MySQL 就绪后才启动后端，后端就绪后才启动前端。

## Configuration

Key configuration options in `.env`:

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_PORT` | Server port | 8080 |
| `DATABASE_HOST` | MySQL host | 127.0.0.1 |
| `DATABASE_PASSWORD` | MySQL password | (required) |
| `AI_API_KEY` | OpenAI API key | (empty = Mock Mode) |
| `AI_MOCK_MODE` | Force mock mode | true |
| `AI_BASE_URL` | AI provider URL | https://api.openai.com/v1 |
| `AI_MODEL` | Text model | gpt-4.1 |
| `AI_VISION_MODEL` | Vision model | gpt-4.1 |
| `MAX_UPLOAD_SIZE_MB` | Max upload size | 20 |
| `CORS_ALLOW_ORIGINS` | 允许的前端来源，逗号分隔 | http://localhost:5173 |

## 切换真实 AI

编辑仓库根目录 `.env`（从 `.env.example` 复制），将 `AI_MOCK_MODE` 设为 `false` 并填入 `AI_API_KEY`，然后 `docker compose up -d`。

## API Endpoints

### System

- `GET /api/health` - Health check
- `GET /api/system/status` - System status (no secrets)
- `GET /api/dashboard/stats` - Dashboard statistics
- `GET /api/tools` - Tool metadata

### Reports

- `GET /api/reports` - List reports (paginated)
- `GET /api/reports/:id` - Report detail
- `DELETE /api/reports/:id` - Delete report
- `GET /api/reports/:id/export?format=markdown` - Export Markdown
- `GET /api/reports/:id/files/:filename` - Download generated file

### Tool Execution

- `POST /api/tools/ui-review/run` - Run UI Review (multipart)
- `POST /api/tools/project-doctor/run` - Run Project Doctor (multipart)
- `POST /api/tools/agent-config/run` - Run Agent Config (JSON)
- `POST /api/tools/api-doc/run` - Run API Doc (multipart)
- `POST /api/tools/db-schema/run` - Run DB Schema (JSON)

## Mock Mode

When `AI_API_KEY` is empty or `AI_MOCK_MODE=true`, the backend uses mock responses that:

- Execute full business chain (validation, persistence, file generation)
- Produce realistic content for UI testing
- Pass through same DTO validation as real AI responses

## Security

- API key never returned to frontend
- File uploads validated by extension, MIME type, and magic bytes
- ZIP protection: path traversal, symlinks, size limits
- Secrets redacted before sending to AI
- User SQL never executed (text analysis only)
- CORS restricted to configured origins

## Development

```bash
# Build
go build ./cmd/server

# Run tests
go test ./...

# Run with hot reload (requires air)
air
```

## Project Structure

```
backend/
├─ cmd/server/main.go          # Entry point
├─ internal/
│  ├─ config/                  # Configuration loading
│  ├─ database/                # MySQL connection
│  ├─ model/                   # GORM models
│  ├─ dto/                     # API DTOs
│  ├─ handler/                 # HTTP handlers
│  ├─ service/                 # Business logic
│  │  ├─ ai_service.go         # AI interface
│  │  ├─ mock_ai_service.go    # Mock implementation
│  │  ├─ file_service.go       # File upload
│  │  ├─ zip_service.go        # ZIP extraction
│  │  └─ tools/                # Tool services
│  ├─ repository/              # Database access
│  ├─ middleware/              # HTTP middleware
│  ├─ prompts/                 # AI prompts
│  └─ util/                    # Utilities
├─ migrations/                 # SQL migrations
├─ uploads/                    # Upload storage
├─ temp/                       # Temp extraction
├─ .env.example
├─ Dockerfile
└─ README.md
```

## License

MIT