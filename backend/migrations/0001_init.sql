-- 0001_init.sql
-- AI Developer Workbench baseline schema
-- MySQL 8.0+, InnoDB, utf8mb4

CREATE TABLE IF NOT EXISTS reports (
  id CHAR(36) NOT NULL,
  tool_type VARCHAR(32) NOT NULL,
  title VARCHAR(255) NOT NULL,
  input_mode VARCHAR(32) NOT NULL DEFAULT '',
  status VARCHAR(24) NOT NULL DEFAULT 'processing',
  summary TEXT NULL,
  total_score SMALLINT UNSIGNED NULL,
  grade VARCHAR(64) NULL,
  input_json JSON NULL,
  report_json JSON NOT NULL,
  file_path VARCHAR(1024) NULL,
  file_url VARCHAR(1024) NULL,
  error_message TEXT NULL,
  created_at DATETIME(3) NOT NULL,
  updated_at DATETIME(3) NOT NULL,
  PRIMARY KEY (id),
  KEY idx_reports_tool_created (tool_type, created_at),
  KEY idx_reports_status_created (status, created_at),
  KEY idx_reports_score_created (total_score, created_at),
  CONSTRAINT chk_reports_score
    CHECK (total_score IS NULL OR total_score BETWEEN 0 AND 100)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS generated_files (
  id CHAR(36) NOT NULL,
  report_id CHAR(36) NOT NULL,
  filename VARCHAR(255) NOT NULL,
  language VARCHAR(32) NULL,
  mime_type VARCHAR(100) NOT NULL DEFAULT 'text/markdown',
  content LONGTEXT NOT NULL,
  size_bytes BIGINT UNSIGNED NOT NULL DEFAULT 0,
  created_at DATETIME(3) NOT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uk_generated_file_report_name (report_id, filename),
  KEY idx_generated_files_report (report_id),
  CONSTRAINT fk_generated_files_report
    FOREIGN KEY (report_id) REFERENCES reports(id)
    ON DELETE CASCADE
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS report_assets (
  id CHAR(36) NOT NULL,
  report_id CHAR(36) NOT NULL,
  asset_type VARCHAR(32) NOT NULL,
  original_name VARCHAR(255) NOT NULL,
  stored_name VARCHAR(255) NOT NULL,
  relative_path VARCHAR(1024) NOT NULL,
  mime_type VARCHAR(100) NULL,
  size_bytes BIGINT UNSIGNED NOT NULL,
  sha256 CHAR(64) NULL,
  created_at DATETIME(3) NOT NULL,
  PRIMARY KEY (id),
  KEY idx_report_assets_report (report_id),
  CONSTRAINT fk_report_assets_report
    FOREIGN KEY (report_id) REFERENCES reports(id)
    ON DELETE CASCADE
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_0900_ai_ci;
