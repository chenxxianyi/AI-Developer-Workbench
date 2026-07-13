# D-07 迁移校验工具

## 校验命令

```bash
# 1. 校验表行数
mysql -u root -p ai_workbench -e "
SELECT 'users' AS tbl, COUNT(*) AS cnt FROM users
UNION ALL SELECT 'projects', COUNT(*) FROM projects
UNION ALL SELECT 'project_requirements', COUNT(*) FROM project_requirements
UNION ALL SELECT 'project_blueprints', COUNT(*) FROM project_blueprints
UNION ALL SELECT 'tasks', COUNT(*) FROM jobs
UNION ALL SELECT 'project_files', COUNT(*) FROM project_files
UNION ALL SELECT 'reports', COUNT(*) FROM reports
UNION ALL SELECT 'ai_runs', COUNT(*) FROM ai_runs;
"

# 2. 校验孤儿记录（project 无对应 user）
mysql -u root -p ai_workbench -e "
SELECT p.id, p.name FROM projects p
LEFT JOIN users u ON u.id = p.user_id
WHERE u.id IS NULL AND p.user_id IS NOT NULL;
"

# 3. 校验 ID 映射完整性
mysql -u root -p ai_workbench -e "
SELECT COUNT(*) AS unmapped_users FROM users WHERE legacy_builder_id IS NOT NULL AND id IS NULL;
SELECT COUNT(*) AS unmapped_projects FROM projects WHERE legacy_source = 'builder' AND legacy_id IS NULL;
"

# 4. 校验 JSON 合法性（蓝图字段）
mysql -u root -p ai_workbench -e "
SELECT id FROM project_blueprints WHERE JSON_VALID(content) = 0;
SELECT id FROM project_requirements WHERE JSON_VALID(content) = 0;
"

# 5. 校验状态枚举
mysql -u root -p ai_workbench -e "
SELECT DISTINCT status FROM projects WHERE status NOT IN ('draft','analyzing','blueprint_pending','generating','building','completed','failed','archived');
SELECT DISTINCT status FROM jobs WHERE status NOT IN ('pending','running','success','failed','cancelled');
"
```

## 校验通过标准

- 所有表行数 ≥ 预期最小值
- 无孤儿记录
- 无未映射的旧系统 ID
- 所有 JSON 字段合法
- 所有状态枚举值合法

## 失败处理

校验失败时返回非零退出码，输出具体记录 ID 和失败原因。
