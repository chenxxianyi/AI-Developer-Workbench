# A-06 认证与权限矩阵

## 认证方案

- **方案**：JWT（JSON Web Token）
- **算法**：HS256
- **过期时间**：Access Token 24h，Refresh Token 7d
- **Payload**：`{ "sub": "<user_id>", "role": "user|admin", "exp": ..., "iat": ... }`
- **传递方式**：`Authorization: Bearer <token>`
- **Mock 模式**：开发环境允许 `X-Mock-User: <user_id>` 跳过认证（不可在生产环境使用）

## 角色定义

| 角色 | 说明 | 权限范围 |
|------|------|---------|
| `anonymous` | 未登录用户 | 仅登录/注册 |
| `user` | 普通用户 | 自己的项目、任务、报告 |
| `admin` | 管理员 | 所有项目、用户管理、模型/Prompt 管理 |

## 权限矩阵

### 认证接口 `/api/v1/auth/`

| 接口 | anonymous | user | admin |
|------|:---------:|:----:|:-----:|
| POST `/auth/register` | ✅ | ✅ | ✅ |
| POST `/auth/login` | ✅ | ✅ | ✅ |
| GET `/auth/profile` | ❌ | ✅ | ✅ |
| PUT `/auth/profile` | ❌ | ✅ | ✅ |
| PUT `/auth/password` | ❌ | ✅ | ✅ |

### 项目接口 `/api/v1/projects/`

| 接口 | anonymous | user（自己的） | user（他人的） | admin |
|------|:---------:|:-------------:|:-------------:|:-----:|
| GET `/projects` | ❌ | ✅（仅自己） | ❌ | ✅（全部） |
| POST `/projects` | ❌ | ✅ | — | ✅ |
| GET `/projects/:id` | ❌ | ✅ | ❌（403） | ✅ |
| PUT `/projects/:id` | ❌ | ✅ | ❌（403） | ✅ |
| DELETE `/projects/:id` | ❌ | ✅ | ❌（403） | ✅ |
| GET `/projects/:id/requirements` | ❌ | ✅ | ❌ | ✅ |
| PUT `/projects/:id/requirements` | ❌ | ✅ | ❌ | ✅ |
| GET `/projects/:id/blueprint` | ❌ | ✅ | ❌ | ✅ |
| POST `/projects/:id/blueprint/generate` | ❌ | ✅ | ❌ | ✅ |
| PUT `/projects/:id/blueprint` | ❌ | ✅ | ❌ | ✅ |
| GET `/projects/:id/files` | ❌ | ✅ | ❌ | ✅ |
| GET `/projects/:id/files/content` | ❌ | ✅ | ❌ | ✅ |
| POST `/projects/:id/build` | ❌ | ✅ | ❌ | ✅ |
| GET `/projects/:id/export` | ❌ | ✅ | ❌ | ✅ |

### 任务接口 `/api/v1/tasks/`

| 接口 | anonymous | user（自己的） | user（他人的） | admin |
|------|:---------:|:-------------:|:-------------:|:-----:|
| GET `/tasks/:id` | ❌ | ✅ | ❌ | ✅ |
| GET `/tasks/:id/stream` | ❌ | ✅（SSE） | ❌ | ✅ |
| POST `/tasks/:id/retry` | ❌ | ✅ | ❌ | ✅ |
| POST `/tasks/:id/cancel` | ❌ | ✅ | ❌ | ✅ |

### 预览接口 `/api/v1/preview/`

| 接口 | anonymous | user | admin |
|------|:---------:|:----:|:-----:|
| GET `/preview/:session/*` | ❌ | ✅（自己的会话） | ✅ |

### 工具接口 `/api/v1/tools/`

| 接口 | anonymous | user | admin |
|------|:---------:|:----:|:-----:|
| POST `/tools/:tool/run` | ❌ | ✅ | ✅ |

### 报告接口 `/api/v1/reports/`

| 接口 | anonymous | user（自己的） | user（他人的） | admin |
|------|:---------:|:-------------:|:-------------:|:-----:|
| GET `/reports` | ❌ | ✅（仅自己） | ❌ | ✅（全部） |
| GET `/reports/:id` | ❌ | ✅ | ❌ | ✅ |
| PUT `/reports/:id` | ❌ | ✅ | ❌ | ✅ |
| GET `/reports/:id/compare/:target` | ❌ | ✅ | ❌ | ✅ |

### 管理后台 `/api/v1/admin/`

| 接口 | anonymous | user | admin |
|------|:---------:|:----:|:-----:|
| 全部 `/admin/*` | ❌ | ❌（403） | ✅ |

## 资源所有权校验

项目资源（requirements、blueprint、files、build、preview、export）通过项目所有权校验：
- 解析路由参数中的 `:projectId`
- 查询 `project.user_id`
- 匹配 JWT 中的 `sub`
- 管理员跳过校验

## Mock 模式规则

- `X-Mock-User` 仅在 `APP_ENV=development` 且 `AI_MOCK_MODE=true` 时生效
- Mock 用户仍需通过权限矩阵校验（不能绕过 `user_id` 匹配）
