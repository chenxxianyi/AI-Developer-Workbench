-- 0004_jobs.sql
-- Async job table.

CREATE TABLE IF NOT EXISTS jobs (
  id CHAR(36) NOT NULL,
  tool_type VARCHAR(32) NOT NULL,
  report_id CHAR(36) NOT NULL,
  project_id CHAR(36) NULL,
  status VARCHAR(24) NOT NULL DEFAULT 'queued',
  progress TINYINT UNSIGNED NOT NULL DEFAULT 0,
  phase VARCHAR(128) NOT NULL DEFAULT '',
  error_message TEXT NULL,
  retry_of_job_id CHAR(36) NULL,
  retry_count TINYINT UNSIGNED NOT NULL DEFAULT 0,
  created_at DATETIME(3) NOT NULL,
  updated_at DATETIME(3) NOT NULL,
  PRIMARY KEY (id),
  KEY idx_jobs_tool_type (tool_type),
  KEY idx_jobs_report_id (report_id),
  KEY idx_jobs_status (status),
  KEY idx_jobs_project_id (project_id),
  KEY idx_jobs_retry_of (retry_of_job_id),
  CONSTRAINT fk_jobs_report
    FOREIGN KEY (report_id) REFERENCES reports(id)
    ON DELETE CASCADE
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_0900_ai_ci;
