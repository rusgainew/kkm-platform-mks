-- Синхронизация данных из users в user_read_model
INSERT INTO user_read_model (
    id, email, first_name, last_name, status, role, 
    created_at, updated_at
)
SELECT 
    id::uuid, 
    email, 
    first_name, 
    last_name, 
    status, 
    role,
    to_timestamp(created_at),
    to_timestamp(updated_at)
FROM users
ON CONFLICT (email) DO NOTHING;

-- Показать результат
SELECT COUNT(*) as total_synced FROM user_read_model;
