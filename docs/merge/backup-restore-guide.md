# M0 数据库和文件资产备份方案

## 1. 数据库资产

### Workbench（主仓库）

| 文件 | 路径 | 说明 |
|------|------|------|
| 0001_init.sql | `backend/migrations/0001_init.sql` | 初始表结构（reports, report_assets 等） |
| 0002_report_lineage.sql | `backend/migrations/0002_report_lineage.sql` | 报告谱系 |
| 0003_projects.sql | `backend/migrations/0003_projects.sql` | 项目表 |
| 0003_z_uuid_collation.sql | `backend/migrations/0003_z_uuid_collation.sql` | UUID 排序规则 |
| 0004_jobs.sql | `backend/migrations/0004_jobs.sql` | 任务表 |
| 0005_ai_runs.sql | `backend/migrations/0005_ai_runs.sql` | AI 运行记录 |

### Builder（能力来源项目）

| 文件 | 路径 | 说明 |
|------|------|------|
| 001_init_schema.sql | `database/migrations/001_init_schema.sql` | 完整初始化（users/projects/requirements/blueprints/tasks/files/models/prompts） |

---

## 2. 文件资产目录

### Workbench

| 目录 | 路径 | 说明 | 风险 |
|------|------|------|------|
| 上传目录 | `backend/uploads/` | 用户上传文件（ZIP、截图等） | 可能含敏感文件 |
| 临时目录 | `backend/temp/` | 临时解压和处理文件 | 可清理 |
| 构建产物 | `frontend/dist/` | 前端构建输出 | 可重建 |
| 后端二进制 | `backend/server.exe` | 编译后的服务端 | 可重建 |

### Builder

| 目录 | 路径 | 说明 | 风险 |
|------|------|------|------|
| 工作区 | `apps/api/cmd/server/workspace/` | 生成的网站项目（含 node_modules） | **最大**，可能含密钥 |
| 构建产物 | `apps/web/dist/` | 前端构建输出 | 可重建 |

---

## 3. 备份命令

### 导出数据库结构（SQLite → SQL 文本）

```bash
# Workbench（如使用 SQLite）
cp backend/*.db backend/backup_$(date +%Y%m%d).db
```

### 导出数据库结构（MySQL）

```bash
# Workbench
mysqldump -u root -p workbench --no-data > backup/workbench_schema_$(date +%Y%m%d).sql

# Builder
mysqldump -u root -p website_builder --no-data > backup/builder_schema_$(date +%Y%m%d).sql
```

### 文件资产备份（排除 node_modules、dist、.git）

```bash
# Workbench
tar -czf backup/workbench_files_$(date +%Y%m%d).tar.gz \
  --exclude='node_modules' --exclude='dist' --exclude='.git' \
  --exclude='*.exe' --exclude='temp/*' \
  backend/uploads/ backend/migrations/ frontend/src/

# Builder
tar -czf backup/builder_files_$(date +%Y%m%d).tar.gz \
  --exclude='node_modules' --exclude='dist' --exclude='.git' \
  --exclude='.npm-cache' --exclude='*.exe' \
  database/migrations/ apps/web/src/ apps/api/internal/
```

---

## 4. 恢复命令

```bash
# 恢复数据库结构
mysql -u root -p workbench < backup/workbench_schema_YYYYMMDD.sql

# 恢复文件资产
tar -xzf backup/workbench_files_YYYYMMDD.tar.gz -C restore/
```

---

## 5. 校验命令

```bash
# 校验表结构一致性
diff <(mysql -u root -p -e "SHOW TABLES" workbench | sort) \
     <(grep -i "CREATE TABLE" backup/workbench_schema.sql | sed 's/.*`\(.*\)`.*/\1/' | sort)

# 校验迁移文件完整性
sha256sum backend/migrations/*.sql > backup/migration_checksums.txt
sha256sum -c backup/migration_checksums.txt
```

---

## 6. 安全规则

- **禁止提交**：`.env`、数据库密码、JWT Secret、API 密钥、用户上传文件
- **禁止备份**：`node_modules/`、`.git/`、`dist/`、`*.exe`、`.npm-cache/`
- **加密存储**：如果备份含敏感数据，使用 `gpg --encrypt`
- **验证后删除**：恢复演练完成后立即删除测试恢复数据
