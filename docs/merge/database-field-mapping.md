# D-01 数据库 Schema 字段映射

> 基于两个项目的迁移 SQL 文件分析

## 1. 表清单

| 表 | Workbench | Builder | 目标 | 备注 |
|----|:--------:|:------:|:----:|------|
| users | ❌ | ✅ | **新建** | 从 Builder 迁移 |
| projects | ✅ | ✅ | **合并** | 字段差异大 |
| project_requirements | ❌ | ✅ | **迁移** | |
| project_blueprints | ❌ | ✅ | **迁移** | |
| tasks | ❌ (jobs) | ✅ | **合并** | jobs → tasks |
| task_progress | ❌ | ✅ | **迁移** | |
| files | ❌ | ✅ | **迁移** | |
| ai_models | ❌ | ✅ | **迁移** | |
| prompt_templates | ❌ | ✅ | **迁移** | |
| reports | ✅ | ❌ | **保留** | |
| report_assets | ✅ | ❌ | **保留** | |
| generated_files | ✅ | ❌ | **保留** | |
| ai_runs | ✅ | ❌ | **保留** | |
| report_lineage | ✅ | ❌ | **保留** | |

## 2. 关键字段映射

### users 表
| Builder 字段 | 目标字段 | 类型变更 | 转换规则 |
|-------------|---------|---------|---------|
| id (uint, PK) | id (varchar(36), PK) | uint → UUID | 生成新 UUID，保留 legacy_id |
| username | username | 一致 | 直接迁移 |
| email | email | 一致 | 直接迁移 |
| password_hash | password_hash | 一致 | 直接迁移 |
| role | role | 一致 | 直接迁移 |
| status | status | 一致 | 直接迁移 |
| created_at | created_at | 一致 | 直接迁移 |
| updated_at | updated_at | 一致 | 直接迁移 |
| — | legacy_builder_id (uint) | 新增 | 记录原始 Builder ID |

### projects 表
| Workbench 字段 | Builder 字段 | 目标字段 | 说明 |
|--------------|-------------|---------|------|
| id (varchar UUID) | id (uint) | id (varchar UUID) | UUID 统一 |
| name | name | name | 合并 |
| description | description | description | 合并 |
| status (简单) | status (详细) | status (扩展) | 8 状态 |
| — | type | type | + |
| — | source_type | source_type | + |
| — | quality_score | quality_score | + |
| — | user_id (uint) | user_id (varchar) | UUID |
| — | blueprint_id | blueprint_id | + |
| created_at | created_at | created_at | 合并 |
| updated_at | updated_at | updated_at | 合并 |
| deleted_at (soft) | — | deleted_at | 保留软删除 |
| — | — | legacy_source | 新增 (workbench/builder) |
| — | — | legacy_id | 新增 (旧系统 ID) |

### tasks (jobs → tasks)
| Workbench jobs | Builder tasks | 目标 tasks | 说明 |
|--------------|-------------|---------|------|
| id | id | id (varchar UUID) | UUID |
| type | type | type (扩展) | generation/build/tool_run/report |
| status | status | status | pending/running/success/failed/cancelled |
| — | project_id | project_id | + |
| — | progress | progress | + |
| — | stage | stage | + |
| — | error_code | error_code | + |
| — | retry_count | retry_count | + |
