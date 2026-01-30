#!/usr/bin/env python3
"""
Синхронизация данных из users в user_read_model
Копирует всех пользователей из write-модели в read-модель
"""

import psycopg2
from psycopg2.extras import execute_batch

# Подключение к БД
conn = psycopg2.connect(
    host="localhost",
    port=5432,
    database="user_db",
    user="useruser",
    password="userpass"
)

cursor = conn.cursor()

try:
    print("Начало синхронизации данных...")
    
    # Проверяем количество записей в исходной таблице
    cursor.execute("SELECT COUNT(*) FROM users WHERE deleted_at IS NULL")
    source_count = cursor.fetchone()[0]
    print(f"Пользователей в таблице users: {source_count}")
    
    # Проверяем количество записей в целевой таблице
    cursor.execute("SELECT COUNT(*) FROM user_read_model WHERE deleted_at IS NULL")
    target_count = cursor.fetchone()[0]
    print(f"Пользователей в таблице user_read_model (до синхронизации): {target_count}")
    
    # Копируем данные (ON CONFLICT DO UPDATE для обновления существующих)
    sync_query = """
    INSERT INTO user_read_model (
        id, email, first_name, last_name, phone, status, role, 
        last_login_at, created_at, updated_at, deleted_at
    )
    SELECT 
        id, email, first_name, last_name, phone, status, role,
        last_login_at, created_at, updated_at, deleted_at
    FROM users
    ON CONFLICT (id) 
    DO UPDATE SET
        email = EXCLUDED.email,
        first_name = EXCLUDED.first_name,
        last_name = EXCLUDED.last_name,
        phone = EXCLUDED.phone,
        status = EXCLUDED.status,
        role = EXCLUDED.role,
        last_login_at = EXCLUDED.last_login_at,
        updated_at = EXCLUDED.updated_at,
        deleted_at = EXCLUDED.deleted_at
    """
    
    cursor.execute(sync_query)
    rows_affected = cursor.rowcount
    conn.commit()
    
    print(f"✓ Синхронизировано записей: {rows_affected}")
    
    # Проверяем результат
    cursor.execute("SELECT COUNT(*) FROM user_read_model WHERE deleted_at IS NULL")
    final_count = cursor.fetchone()[0]
    print(f"Пользователей в таблице user_read_model (после синхронизации): {final_count}")
    
    # Показываем примеры данных
    cursor.execute("""
        SELECT id, email, first_name, last_name, status, role 
        FROM user_read_model 
        WHERE deleted_at IS NULL 
        ORDER BY created_at DESC 
        LIMIT 5
    """)
    
    print("\nПримеры синхронизированных пользователей:")
    for row in cursor.fetchall():
        print(f"  - {row[2]} {row[3]} ({row[1]}) - {row[5]} [{row[4]}]")
    
    print(f"\n✓ Синхронизация завершена успешно!")
    print(f"  Всего пользователей в read-модели: {final_count}")

except Exception as e:
    print(f"✗ Ошибка при синхронизации: {e}")
    conn.rollback()
    raise

finally:
    cursor.close()
    conn.close()
