# A-03 统一 API 规范

## 基础约定

| 项目 | 规范 |
|------|------|
| 前缀 | `/api/v1` |
| 认证 | `Authorization: Bearer <JWT>` |
| Content-Type | `application/json` (默认)、`multipart/form-data` (上传) |
| 字符编码 | UTF-8 |
| 时间格式 | ISO 8601 (`2026-07-13T10:30:00Z`) |
| ID 格式 | UUID v4 字符串 |
| Request ID | 每个响应包含 `X-Request-ID` 头 |

## 成功响应

```json
{
  "data": { ... },
  "meta": {
    "page": 1,
    "page_size": 20,
    "total": 150
  }
}
```

- `data`：单个对象或数组
- `meta`：分页信息（仅列表接口返回）

## 错误响应

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "项目名称不能为空",
    "details": [
      {
        "field": "name",
        "reason": "required"
      }
    ]
  },
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

### 标准错误码

| HTTP 状态 | 错误码 | 说明 |
|-----------|--------|------|
| 400 | `VALIDATION_ERROR` | 请求参数校验失败 |
| 401 | `UNAUTHORIZED` | 未认证或 Token 过期 |
| 403 | `FORBIDDEN` | 无权限访问资源 |
| 404 | `NOT_FOUND` | 资源不存在 |
| 409 | `CONFLICT` | 资源冲突（如重复创建） |
| 413 | `PAYLOAD_TOO_LARGE` | 请求体/文件过大 |
| 422 | `BUSINESS_ERROR` | 业务逻辑错误（如非法状态转换） |
| 429 | `RATE_LIMITED` | 请求频率超限 |
| 500 | `INTERNAL_ERROR` | 服务端内部错误 |
| 503 | `SERVICE_UNAVAILABLE` | 服务暂不可用 |

## 分页规范

请求参数：
```
GET /api/v1/projects?page=1&page_size=20&sort=created_at&order=desc
```

| 参数 | 类型 | 默认值 | 最大值 |
|------|------|--------|--------|
| `page` | int | 1 | — |
| `page_size` | int | 20 | 100 |
| `sort` | string | `created_at` | — |
| `order` | string | `desc` | — |

## 文件上传

```
POST /api/v1/projects/:id/upload
Content-Type: multipart/form-data

字段：file (binary)
最大：50MB（可配置）
```

成功响应：
```json
{
  "data": {
    "id": "uuid",
    "name": "screenshot.png",
    "size": 204800,
    "mime_type": "image/png",
    "url": "/api/v1/files/uuid"
  }
}
```

## 文件下载

```
GET /api/v1/projects/:id/export
Response: application/zip
Content-Disposition: attachment; filename="project-name.zip"
```

## SSE 事件流

```
GET /api/v1/tasks/:id/stream
Authorization: Bearer <JWT>
Accept: text/event-stream
```

Nginx 配置要求：
```nginx
proxy_buffering off;
proxy_cache off;
proxy_read_timeout 3600s;
```

事件格式：
```
event: progress
data: {"stage":"building","progress":75,"message":"Building frontend..."}

event: complete
data: {"status":"success"}

event: error
data: {"code":"BUILD_FAILED","message":"npm install failed"}
```

心跳：每 30 秒发送 `: heartbeat`

## 预览会话

```
POST /api/v1/projects/:id/build
→ { "data": { "preview_url": "/api/v1/preview/<session-id>/", "expires_at": "..." } }
```

- 会话有效期：2 小时
- Cookie 鉴权：`preview_token=<jwt>`
- 支持 SPA History Fallback

## 旧 API 映射

| 旧 Workbench `/api` | 新 `/api/v1` |
|---------------------|--------------|
| `/api/health` | `/api/v1/health` |
| `/api/dashboard` | `/api/v1/dashboard` |
| `/api/projects` | `/api/v1/projects` |
| `/api/reports` | `/api/v1/reports` |
| `/api/tools/:tool/run` | `/api/v1/tools/:tool/run` |

| 旧 Builder `/api/v1` | 新 `/api/v1` |
|----------------------|--------------|
| `/api/v1/auth/*` | 保留不变 |
| `/api/v1/projects/*` | 保留不变 |
| `/api/v1/tasks/*` | 保留不变 |
| `/api/v1/preview/*` | 保留不变 |
| `/api/v1/admin/*` | 保留不变 |
