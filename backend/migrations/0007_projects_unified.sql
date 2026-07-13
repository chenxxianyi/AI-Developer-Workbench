-- D-03: 统一项目表（合并 Workbench projects + Builder projects）
-- 变更：统一字段，增加 type/source_type/user_id/blueprint_id/quality_score
-- 注：迁移系统已追踪版本，每条 SQL 仅执行一次，无需 IF NOT EXISTS

ALTER TABLE projects
    ADD COLUMN type VARCHAR(30) NOT NULL DEFAULT 'website' AFTER description,
    ADD COLUMN source_type VARCHAR(30) NOT NULL DEFAULT 'generated' AFTER type,
    ADD COLUMN source_url VARCHAR(500) NULL AFTER source_type,
    ADD COLUMN user_id VARCHAR(36) NULL AFTER source_url,
    ADD COLUMN blueprint_id VARCHAR(36) NULL AFTER user_id,
    ADD COLUMN quality_score DECIMAL(5,2) NULL AFTER blueprint_id,
    ADD COLUMN legacy_source VARCHAR(20) NULL COMMENT 'workbench|builder',
    ADD COLUMN legacy_id VARCHAR(100) NULL COMMENT '旧系统原始ID',
    ADD COLUMN status VARCHAR(30) NOT NULL DEFAULT 'draft'
        COMMENT 'draft|analyzing|blueprint_pending|generating|building|completed|failed|archived';

CREATE INDEX idx_projects_user_id ON projects(user_id);
CREATE INDEX idx_projects_type ON projects(type);
CREATE INDEX idx_projects_status ON projects(status);
