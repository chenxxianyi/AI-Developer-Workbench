-- 0003_projects.sql
-- Project profile table and reports.project_id lineage.
-- A project is a user-managed profile (name, repo, stacks, rules) that reports
-- can be associated with. Deleting a project SET NULLs reports.project_id so
-- historical reports survive (no cascade), matching the product rule.

CREATE TABLE IF NOT EXISTS projects (
  id CHAR(36) NOT NULL,
  name VARCHAR(128) NOT NULL,
  description TEXT NULL,
  repo_url VARCHAR(512) NULL,
  frontend_stack VARCHAR(256) NULL,
  backend_stack VARCHAR(256) NULL,
  database VARCHAR(128) NULL,
  ui_style VARCHAR(256) NULL,
  coding_rules TEXT NULL,
  created_at DATETIME(3) NOT NULL,
  updated_at DATETIME(3) NOT NULL,
  PRIMARY KEY (id),
  KEY idx_projects_name (name),
  KEY idx_projects_updated (updated_at)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_0900_ai_ci;

ALTER TABLE reports
  ADD COLUMN project_id CHAR(36) NULL AFTER parent_report_id,
  ADD KEY idx_reports_project (project_id),
  ADD CONSTRAINT fk_reports_project
    FOREIGN KEY (project_id) REFERENCES projects(id)
    ON DELETE SET NULL;
