# A-04 统一 Project 领域模型

## 统一 Project 结构

```go
type Project struct {
    ID          string     `gorm:"primaryKey;type:varchar(36)" json:"id"`
    Name        string     `gorm:"not null" json:"name"`
    Description string     `json:"description"`
    Slug        string     `gorm:"uniqueIndex" json:"slug"`

    // 项目类型和来源
    Type        string     `gorm:"not null;default:'website'" json:"type"`      // website | analysis | import
    SourceType  string     `gorm:"default:'generated'" json:"source_type"`       // generated | imported | manual
    SourceURL   string     `json:"source_url,omitempty"`                         // GitHub URL 等

    // 状态
    Status      string     `gorm:"not null;default:'draft'" json:"status"`       // draft | analyzing | blueprint_pending | generating | building | completed | failed | archived

    // 质量分
    QualityScore *float64  `json:"quality_score,omitempty"`

    // 用户归属
    UserID      string     `gorm:"not null;index" json:"user_id"`

    // 蓝图 ID（生成项目）
    BlueprintID *string    `json:"blueprint_id,omitempty"`

    // 旧系统映射
    LegacyBuilderID *uint  `json:"legacy_builder_id,omitempty"`
    LegacySource    string `json:"legacy_source,omitempty"`                      // "workbench" | "builder"

    // 时间戳
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
    DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}
```

## 项目类型

| 类型 | 说明 | 关联 |
|------|------|------|
| `website` | 生成网站的完整流程 | 需求 → 蓝图 → 生成 → 构建 → 预览 |
| `analysis` | 分析已有项目 | 文件 → 工具 → 报告 |
| `import` | 导入 ZIP/GitHub 项目 | 文件 → 工具 → 报告 |

## 来源类型

| 来源 | 说明 |
|------|------|
| `generated` | AI 生成的项目 |
| `imported` | 用户导入的项目（ZIP/GitHub） |
| `manual` | 手动创建的项目（仅分析） |

## 状态机

```text
draft → analyzing → blueprint_pending → generating → building → completed
  ↓         ↓              ↓                ↓           ↓
  └─────────┴──────────────┴────────────────┴───────────→ failed
                                                            ↓
                                                        archived
```

状态转换规则：
- `draft` → `analyzing`：提交分析
- `draft` → `blueprint_pending`：跳过分析直接生成
- `blueprint_pending` → `generating`：确认蓝图
- `generating` → `building`：生成完成
- `building` → `completed`：构建成功
- 任意状态 → `failed`：操作失败
- `completed`/`failed` → `archived`：归档

## JSON 表达

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "我的企业网站",
  "slug": "my-company-site",
  "type": "website",
  "source_type": "generated",
  "status": "generating",
  "quality_score": 85.5,
  "user_id": "user-uuid",
  "blueprint_id": "bp-uuid",
  "created_at": "2026-07-13T10:30:00Z",
  "updated_at": "2026-07-13T11:00:00Z"
}
```

## 字段映射：旧 → 新

| Workbench 字段 | Builder 字段 | 统一字段 | 说明 |
|---------------|-------------|---------|------|
| `id` (string UUID) | `id` (uint) | `id` (string UUID) | 统一为 UUID |
| `name` | `name` | `name` | 一致 |
| `description` | `description` | `description` | 一致 |
| `status` (简单) | `status` (详细) | `status` | 扩展状态枚举 |
| — | `type` | `type` | Builder 独有 |
| — | `source_type` | `source_type` | Builder 独有 |
| — | `quality_score` | `quality_score` | Builder 独有 |
| — | `user_id` | `user_id` | Builder 独有 |
| — | `blueprint_id` | `blueprint_id` | Builder 独有 |
| `created_at` | `created_at` | `created_at` | 一致 |
