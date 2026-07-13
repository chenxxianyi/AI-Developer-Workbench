# D-08 数据库回滚和重跑策略

## 迁移原则

1. **幂等**：所有迁移使用 `IF NOT EXISTS` / `ADD COLUMN IF NOT EXISTS`
2. **可回滚**：每个迁移文件底部注释回滚 SQL
3. **可重跑**：中断后可从断点继续

## 回滚边界

| 阶段 | 回滚方式 |
|------|---------|
| 仅结构迁移完成 | 执行回滚 SQL |
| 数据已导入 | 从备份恢复 |
| 新系统已上线 | 双读期间切回旧系统 |

## 中断恢复

```bash
# 查看已执行的迁移
mysql -u root -p ai_workbench -e "SELECT * FROM schema_migrations;"

# 从中断点继续（GORM AutoMigrate 自动跳过已存在的列/表）
./server  # AutoMigrate=true 时自动处理
```

## 回滚脚本模板

```sql
-- 回滚 D-02: users
DROP TABLE IF EXISTS users;

-- 回滚 D-03: projects 字段
ALTER TABLE projects
    DROP COLUMN IF EXISTS type,
    DROP COLUMN IF EXISTS source_type,
    DROP COLUMN IF EXISTS user_id;

-- 回滚 D-04: requirements/blueprints
DROP TABLE IF EXISTS project_blueprints;
DROP TABLE IF EXISTS project_requirements;

-- 回滚 D-05: tasks/files
ALTER TABLE jobs
    DROP COLUMN IF EXISTS project_id,
    DROP COLUMN IF EXISTS progress;
DROP TABLE IF EXISTS project_files;
```

## 备份恢复

```bash
# 恢复完整数据库
mysql -u root -p ai_workbench < backup/pre_migration_dump.sql

# 恢复单个表
mysql -u root -p ai_workbench < backup/users_table.sql
```
