-- 0005_ai_runs.sql
-- AI provider call observability table.

CREATE TABLE IF NOT EXISTS ai_runs (
  id CHAR(36) NOT NULL,
  report_id CHAR(36) NOT NULL,
  job_id CHAR(36) NULL,
  tool_type VARCHAR(32) NOT NULL,
  provider VARCHAR(64) NOT NULL,
  model VARCHAR(128) NOT NULL,
  is_mock TINYINT(1) NOT NULL DEFAULT 0,
  duration_ms BIGINT NOT NULL DEFAULT 0,
  retry_count TINYINT UNSIGNED NOT NULL DEFAULT 0,
  parse_success TINYINT(1) NOT NULL DEFAULT 1,
  fallback_used TINYINT(1) NOT NULL DEFAULT 0,
  error_type VARCHAR(64) NOT NULL DEFAULT '',
  created_at DATETIME(3) NOT NULL,
  PRIMARY KEY (id),
  KEY idx_ai_runs_report_id (report_id),
  KEY idx_ai_runs_job_id (job_id),
  KEY idx_ai_runs_tool_type (tool_type),
  KEY idx_ai_runs_created (created_at),
  CONSTRAINT fk_ai_runs_report
    FOREIGN KEY (report_id) REFERENCES reports(id)
    ON DELETE CASCADE,
  CONSTRAINT fk_ai_runs_job
    FOREIGN KEY (job_id) REFERENCES jobs(id)
    ON DELETE SET NULL
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_0900_ai_ci;
