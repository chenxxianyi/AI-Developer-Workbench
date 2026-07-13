-- D-02: 统一用户表迁移
-- 来源：Builder 001_init_schema.sql (users 表)
-- 变更：id uint → varchar(36) UUID, 增加 legacy_builder_id

CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'user',       -- user | admin
    status VARCHAR(20) NOT NULL DEFAULT 'active',    -- active | disabled
    legacy_builder_id INT UNSIGNED NULL,              -- 旧 Builder 系统 ID（迁移后保留）
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE INDEX idx_username (username),
    UNIQUE INDEX idx_email (email),
    INDEX idx_role (role),
    INDEX idx_legacy_builder_id (legacy_builder_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 迁移数据：从 Builder users 表导入
-- INSERT INTO users (id, username, email, password_hash, role, status, legacy_builder_id, created_at, updated_at)
-- SELECT UUID(), username, email, password_hash, role, status, id, created_at, updated_at
-- FROM ai_website_builder.users;

-- 回滚：DROP TABLE IF EXISTS users;
