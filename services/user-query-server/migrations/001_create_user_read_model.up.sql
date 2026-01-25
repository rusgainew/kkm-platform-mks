-- User Query Server Read Model
-- Оптимизирована для быстрого чтения и фильтрации

CREATE TABLE IF NOT EXISTS user_read_model (
    id UUID PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    full_name VARCHAR(511) GENERATED ALWAYS AS (
        COALESCE(first_name, '') || ' ' || COALESCE(last_name, '')
    ) STORED,
    phone VARCHAR(20),
    status VARCHAR(50) NOT NULL,
    role VARCHAR(50) NOT NULL,
    last_login_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

-- Индексы для быстрого поиска
CREATE INDEX IF NOT EXISTS idx_user_email ON user_read_model(email);
CREATE INDEX IF NOT EXISTS idx_user_status ON user_read_model(status);
CREATE INDEX IF NOT EXISTS idx_user_role ON user_read_model(role);
CREATE INDEX IF NOT EXISTS idx_user_created_at ON user_read_model(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_user_full_name ON user_read_model USING GIN(to_tsvector('russian', full_name));

-- View для активных пользователей
CREATE OR REPLACE VIEW active_users AS
SELECT * FROM user_read_model
WHERE deleted_at IS NULL
ORDER BY created_at DESC;

-- Комментарии
COMMENT ON TABLE user_read_model IS 'Read-side model для user-query-server (CQRS)';
COMMENT ON COLUMN user_read_model.full_name IS 'Генерируется автоматически из first_name и last_name';
