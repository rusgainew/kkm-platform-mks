-- Синхронизация первой 1000 записей для теста
INSERT INTO user_read_model (
    id, email, first_name, last_name, status, role, 
    created_at, updated_at
)
SELECT 
    id, email, first_name, last_name, status, role,
    created_at, updated_at
FROM users
LIMIT 1000
ON CONFLICT (email) DO NOTHING;

SELECT COUNT(*) as total FROM user_read_model;
