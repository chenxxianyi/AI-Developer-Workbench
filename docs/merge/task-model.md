# A-05 统一 Task 领域模型

## 统一 Task 结构

```go
type Task struct {
    ID          string     `gorm:"primaryKey;type:varchar(36)" json:"id"`
    ProjectID   string     `gorm:"not null;index" json:"project_id"`
    UserID      string     `gorm:"not null;index" json:"user_id"`

    // 任务类型
    Type        string     `gorm:"not null" json:"type"`       // generation | build | tool_run | report

    // 状态
    Status      string     `gorm:"not null;default:'pending'" json:"status"` // pending | running | success | failed | cancelled

    // 进度
    Progress    int        `gorm:"default:0" json:"progress"`  // 0-100
    Stage       string     `json:"stage,omitempty"`            // 当前阶段名称
    Message     string     `json:"message,omitempty"`          // 当前阶段描述

    // 结果
    Result      string     `gorm:"type:text" json:"result,omitempty"`       // JSON 格式结果
    ErrorCode   string     `json:"error_code,omitempty"`
    ErrorDetail string     `gorm:"type:text" json:"error_detail,omitempty"`

    // 重试
    RetryCount  int        `gorm:"default:0" json:"retry_count"`
    MaxRetries  int        `gorm:"default:3" json:"max_retries"`
    Retryable   *bool      `gorm:"default:true" json:"retryable"`

    // 时间
    StartedAt   *time.Time `json:"started_at,omitempty"`
    FinishedAt  *time.Time `json:"finished_at,omitempty"`
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
}
```

## 任务类型

| 类型 | 说明 | 关联数据 |
|------|------|---------|
| `generation` | 代码生成 Pipeline | 蓝图 ID |
| `build` | 项目构建 | 构建日志 |
| `tool_run` | AI 工具运行 | 工具名称、输入参数 |
| `report` | 报告生成 | 报告 ID |

## 状态机

```text
                    ┌──────────┐
                    │  pending │
                    └────┬─────┘
                         │ start
                    ┌────▼─────┐
              ┌─────│ running  │─────┐
              │     └────┬─────┘     │
              │ cancel   │  finish   │ error
         ┌────▼─────┐   │      ┌────▼─────┐
         │ cancelled│   │      │  failed   │
         └──────────┘   │      └────┬─────┘
                         │           │ retry (if retryable && retry_count < max_retries)
                    ┌────▼─────┐     │
                    │ success   │     │
                    └──────────┘     │
                               ┌────▼─────┐
                               │  pending  │ (重置后重新开始)
                               └──────────┘
```

## 非法状态转换

以下转换应被后端拒绝：

| 当前状态 | 非法目标 | 原因 |
|---------|---------|------|
| `success` | `running` | 已完成任务不可重新运行 |
| `cancelled` | `running` | 已取消任务不可重新运行 |
| `running` | `pending` | 运行中不可回退 |
| `failed` | `running` | 应通过 retry 创建新 task |

## JSON 表达

```json
{
  "id": "task-uuid",
  "project_id": "project-uuid",
  "user_id": "user-uuid",
  "type": "generation",
  "status": "running",
  "progress": 65,
  "stage": "frontend_generation",
  "message": "正在生成前端页面 (3/5)...",
  "retry_count": 0,
  "max_retries": 3,
  "retryable": true,
  "started_at": "2026-07-13T10:30:00Z",
  "created_at": "2026-07-13T10:29:00Z"
}
```

## 与 Workbench Job 模型的映射

| Workbench Job | 统一 Task |
|--------------|-----------|
| `id` | `id` |
| `status` | `status` |
| `type` | `type` (扩展枚举) |
| — | `project_id` (新增) |
| — | `stage` (新增) |
| — | `progress` (新增) |
| — | `retry_count` / `max_retries` (新增) |
