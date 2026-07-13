-- 0002_report_lineage.sql
-- Add optional parent_report_id to support "re-run from this report" lineage.
-- A child report points at its parent; deleting the parent SET NULLs the FK so
-- the child report survives (no cascade), matching the product rule that
-- re-running never deletes the original report.

ALTER TABLE reports
  ADD COLUMN parent_report_id CHAR(36) NULL AFTER error_message,
  ADD KEY idx_reports_parent (parent_report_id),
  ADD CONSTRAINT fk_reports_parent
    FOREIGN KEY (parent_report_id) REFERENCES reports(id)
    ON DELETE SET NULL;
