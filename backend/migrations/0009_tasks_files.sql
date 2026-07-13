-- D-05: 统一任务表（合并 Workbench jobs + Builder tasks）
-- D-06: 文件表
-- 注：原 jobs 表已有 project_id/progress/retry_count/error_message 列，需 MODIFY 而非 ADD

ALTER TABLE jobs
    MODIFY COLUMN project_id VARCHAR(36) NULL COMMENT 'project_id (原 CHAR(36) 改为 VARCHAR)',
    ADD COLUMN user_id VARCHAR(36) NULL AFTER project_id,
    MODIFY COLUMN progress INT NOT NULL DEFAULT 0 COMMENT 'progress (原 TINYINT 改为 INT)',
    ADD COLUMN stage VARCHAR(100) NULL AFTER progress,
    ADD COLUMN error_code VARCHAR(50) NULL AFTER error_message,
    ADD COLUMN error_detail TEXT NULL AFTER error_code,
    MODIFY COLUMN retry_count INT NOT NULL DEFAULT 0 COMMENT 'retry_count (原 TINYINT 改为 INT)',
    ADD COLUMN max_retries INT NOT NULL DEFAULT 3,
    ADD COLUMN retryable TINYINT(1) NOT NULL DEFAULT 1,
    ADD COLUMN started_at DATETIME NULL,
    ADD COLUMN finished_at DATETIME NULL;

-- RENAME TABLE jobs TO tasks; -- 数据校验通过后手动执行

CREATE TABLE IF NOT EXISTS project_files (
    id VARCHAR(36) PRIMARY KEY,
    project_id VARCHAR(36) NOT NULL,
    name VARCHAR(500) NOT NULL,
    path VARCHAR(1000) NOT NULL,
    size BIGINT NOT NULL DEFAULT 0,
    mime_type VARCHAR(100) NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_file_project (project_id),
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
