#!/usr/bin/env python3
"""
Скрипт для генерации 10000 пользователей с реальными именами и фамилиями
"""

import uuid
import bcrypt
import time
import random
from faker import Faker
import psycopg2
from psycopg2.extras import execute_batch

# Конфигурация подключения к БД
DB_CONFIG = {
    'host': 'localhost',
    'port': 5432,
    'database': 'user_db',
    'user': 'user_svc',
    'password': 'user_pass_2026'
}

# Параметры генерации
NUM_USERS = 10000
BATCH_SIZE = 1000
DEFAULT_PASSWORD = 'Password123!'  # Один пароль для всех для простоты

# Статусы и роли
STATUSES = ['active', 'inactive', 'suspended']
ROLES = ['user', 'admin', 'manager', 'moderator']

# Вероятности для статусов и ролей
STATUS_WEIGHTS = [0.85, 0.10, 0.05]  # 85% active, 10% inactive, 5% suspended
ROLE_WEIGHTS = [0.90, 0.02, 0.05, 0.03]  # 90% user, 2% admin, 5% manager, 3% moderator


def generate_password_hash(password: str) -> str:
    """Генерирует bcrypt хэш пароля"""
    salt = bcrypt.gensalt(rounds=10)
    return bcrypt.hashpw(password.encode('utf-8'), salt).decode('utf-8')


def generate_users(num_users: int, faker: Faker) -> list:
    """Генерирует список пользователей"""
    users = []
    password_hash = generate_password_hash(DEFAULT_PASSWORD)
    
    # Используем время создания с небольшими интервалами
    base_time = int(time.time() * 1000) - (num_users * 1000)  # Начинаем с прошлого времени
    
    print(f"Генерируем {num_users} пользователей...")
    
    for i in range(num_users):
        user_id = str(uuid.uuid4())
        
        # Генерируем имя и фамилию (50% кириллица, 50% латиница)
        if random.random() < 0.5:
            first_name = faker['ru'].first_name()
            last_name = faker['ru'].last_name()
        else:
            first_name = faker['en'].first_name()
            last_name = faker['en'].last_name()
        
        # Генерируем email на основе имени и фамилии
        email_base = f"{first_name.lower()}.{last_name.lower()}"
        # Убираем спецсимволы для email
        email_base = ''.join(c for c in email_base if c.isalnum() or c == '.')
        email = f"{email_base}{i}@example.com"
        
        # Выбираем статус и роль с учетом весов
        status = random.choices(STATUSES, weights=STATUS_WEIGHTS)[0]
        role = random.choices(ROLES, weights=ROLE_WEIGHTS)[0]
        
        # Время создания с интервалом ~1 секунда между пользователями
        created_at = base_time + (i * 1000)
        updated_at = created_at
        
        users.append({
            'id': user_id,
            'email': email,
            'password_hash': password_hash,
            'first_name': first_name,
            'last_name': last_name,
            'status': status,
            'role': role,
            'created_at': created_at,
            'updated_at': updated_at
        })
        
        if (i + 1) % 1000 == 0:
            print(f"  Сгенерировано {i + 1}/{num_users} пользователей...")
    
    print(f"✓ Генерация завершена: {len(users)} пользователей")
    return users


def insert_users_to_db(users: list, batch_size: int = 1000):
    """Вставляет пользователей в БД пакетами"""
    print(f"\nПодключаемся к БД {DB_CONFIG['database']}...")
    
    try:
        conn = psycopg2.connect(**DB_CONFIG)
        cursor = conn.cursor()
        
        print(f"✓ Подключение установлено")
        
        # SQL для вставки
        insert_query = """
        INSERT INTO users (id, email, password_hash, first_name, last_name, status, role, created_at, updated_at)
        VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s)
        ON CONFLICT (email) DO NOTHING;
        """
        
        # Подготавливаем данные для вставки
        user_tuples = [
            (
                u['id'], u['email'], u['password_hash'], 
                u['first_name'], u['last_name'], 
                u['status'], u['role'], 
                u['created_at'], u['updated_at']
            )
            for u in users
        ]
        
        # Вставляем пакетами
        total_batches = (len(user_tuples) + batch_size - 1) // batch_size
        print(f"\nВставляем {len(user_tuples)} пользователей пакетами по {batch_size}...")
        
        for i in range(0, len(user_tuples), batch_size):
            batch = user_tuples[i:i + batch_size]
            batch_num = i // batch_size + 1
            
            execute_batch(cursor, insert_query, batch, page_size=batch_size)
            conn.commit()
            
            print(f"  Пакет {batch_num}/{total_batches}: вставлено {len(batch)} записей")
        
        # Проверяем количество вставленных записей
        cursor.execute("SELECT COUNT(*) FROM users;")
        total_count = cursor.fetchone()[0]
        
        print(f"\n✓ Вставка завершена")
        print(f"✓ Всего пользователей в БД: {total_count}")
        
        cursor.close()
        conn.close()
        
    except psycopg2.Error as e:
        print(f"\n✗ Ошибка БД: {e}")
        raise
    except Exception as e:
        print(f"\n✗ Ошибка: {e}")
        raise


def print_statistics(users: list):
    """Выводит статистику по сгенерированным пользователям"""
    print("\n" + "="*60)
    print("СТАТИСТИКА СГЕНЕРИРОВАННЫХ ПОЛЬЗОВАТЕЛЕЙ")
    print("="*60)
    
    # Статистика по статусам
    status_counts = {}
    for status in STATUSES:
        count = sum(1 for u in users if u['status'] == status)
        status_counts[status] = count
    
    print("\nПо статусам:")
    for status, count in status_counts.items():
        percent = (count / len(users)) * 100
        print(f"  {status:12s}: {count:6d} ({percent:5.2f}%)")
    
    # Статистика по ролям
    role_counts = {}
    for role in ROLES:
        count = sum(1 for u in users if u['role'] == role)
        role_counts[role] = count
    
    print("\nПо ролям:")
    for role, count in role_counts.items():
        percent = (count / len(users)) * 100
        print(f"  {role:12s}: {count:6d} ({percent:5.2f}%)")
    
    print("\nПримеры пользователей:")
    for i, user in enumerate(users[:5], 1):
        print(f"\n  {i}. {user['first_name']} {user['last_name']}")
        print(f"     Email: {user['email']}")
        print(f"     Role: {user['role']}, Status: {user['status']}")
    
    print("\n" + "="*60)
    print(f"ПАРОЛЬ ДЛЯ ВСЕХ ПОЛЬЗОВАТЕЛЕЙ: {DEFAULT_PASSWORD}")
    print("="*60 + "\n")


def main():
    """Главная функция"""
    print("="*60)
    print("ГЕНЕРАТОР ПОЛЬЗОВАТЕЛЕЙ")
    print("="*60)
    print(f"Количество: {NUM_USERS}")
    print(f"База данных: {DB_CONFIG['host']}:{DB_CONFIG['port']}/{DB_CONFIG['database']}")
    print(f"Пароль для всех: {DEFAULT_PASSWORD}")
    print("="*60 + "\n")
    
    # Инициализируем Faker для русского и английского языков
    faker_dict = {
        'ru': Faker('ru_RU'),
        'en': Faker('en_US')
    }
    
    try:
        # Генерируем пользователей
        users = generate_users(NUM_USERS, faker_dict)
        
        # Выводим статистику
        print_statistics(users)
        
        # Спрашиваем подтверждение
        response = input("Вставить пользователей в БД? (yes/no): ")
        if response.lower() in ['yes', 'y', 'да', 'д']:
            insert_users_to_db(users, batch_size=BATCH_SIZE)
            print("\n✓ Готово!")
        else:
            print("\n✗ Отменено")
            
    except KeyboardInterrupt:
        print("\n\n✗ Прервано пользователем")
    except Exception as e:
        print(f"\n✗ Ошибка: {e}")
        raise


if __name__ == '__main__':
    main()
