# ⚡ Быстрый старт KKM Project MKS

**Время:** ~10 минут | **Сложность:** ⭐ Легко | **Обновлено:** 10 января 2026

---

## 🎯 Цель

Запустить полный production stack из 16 контейнеров (API, БД, кеш, мониторинг) и проверить работу API через Nginx Reverse Proxy.

---

## ✅ Требования

- ✔️ Docker 20.10+
- ✔️ Docker Compose 2.0+
- ✔️ 4GB RAM свободно
- ✔️ Порты 80, 443, 5432, 6379, 5672, 3000, 9091, 16686 свободны

---

## 🚀 Шаг 1: Запустить сервисы (2 минуты)

```bash
# 1. Перейдите в корень проекта
cd /path/to/kkm-project-mks

# 2. Запустите production stack с Nginx
docker-compose -f docker-compose.prod.yml up -d

# Вывод должен быть:
# ✔ Container postgres Started
# ✔ Container redis Started
# ✔ Container rabbitmq Started
# ... (еще 13 контейнеров)
# ✔ Container nginx-proxy Started
```

**Что запустилось:**

- 🔵 **6 компонентов инфраструктуры**: PostgreSQL, Redis, RabbitMQ, Prometheus, Jaeger, Grafana
- 🟢 **1 API Gateway**: Маршрутизирует HTTP → gRPC
- 🟣 **9 микросервисов**: User, Company, Catalog, Invoice, Bank Account, Foreign Company + их query версии
- 🔴 **Nginx Reverse Proxy**: Фасад для всех HTTP запросов

---

## 🔍 Шаг 2: Проверить здоровье (1 минута)

```bash
# Проверьте статус контейнеров
docker-compose -f docker-compose.prod.yml ps

# Все должны быть "Up" и "healthy" ✅
# Примеры вывода:
# SERVICE          STATUS              PORTS
# nginx-proxy      Up 30s (healthy)    0.0.0.0:80->80/tcp, 0.0.0.0:443->443/tcp
# api-gateway      Up 1m (healthy)     0.0.0.0:8080->8080/tcp
# postgres         Up 2m (healthy)     0.0.0.0:5432->5432/tcp
# redis            Up 2m (healthy)     0.0.0.0:6379->6379/tcp
```

### Тест здоровья API

```bash
# Проверьте что API отвечает через Nginx
curl http://localhost/api/v1/health

# Ожидаемый ответ:
# {"status":"healthy","time":1768056504}
```

❌ Если не работает:

```bash
# Посмотрите логи Nginx
docker-compose -f docker-compose.prod.yml logs nginx-proxy

# Посмотрите логи API Gateway
docker-compose -f docker-compose.prod.yml logs api-gateway
```

---

## 📚 Шаг 3: Открыть документацию (1 минута)

### Swagger UI (самый простой способ)

```bash
# Откройте в браузере
open http://localhost/swagger/index.html

# Или через curl
curl http://localhost/swagger/index.html
```

### Альтернативные endpoints

| Инструмент          | URL                                 | Описание                   |
| ------------------- | ----------------------------------- | -------------------------- |
| 🔵 **Swagger UI**   | http://localhost/swagger/index.html | Интерактивная документация |
| 📄 **OpenAPI JSON** | http://localhost/swagger/doc.json   | Полная спецификация API    |
| 🟡 **Jaeger UI**    | http://localhost:16686              | Трейсинг запросов          |
| 🔴 **Grafana**      | http://localhost:3000               | Мониторинг (admin/admin)   |
| 🟢 **Prometheus**   | http://localhost:9091               | Метрики                    |

---

## 🧪 Шаг 4: Проверить API (3 минуты)

### 4.1 Создать пользователя

```bash
# Регистрация
curl -X POST http://localhost/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "Test@12345"
  }'

# Ожидаемый ответ:
# {"id":"550e8400-e29b-41d4-a716-446655440000","username":"testuser","email":"test@example.com"}
```

### 4.2 Войти в систему

```bash
# Логин
RESPONSE=$(curl -s -X POST http://localhost/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "Test@12345"
  }')

# Извлеките токен
TOKEN=$(echo $RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)

echo "Ваш токен: $TOKEN"
```

### 4.3 Использовать защищённые endpoints

```bash
# Получить информацию о пользователе
curl http://localhost/api/v1/users/me \
  -H "Authorization: Bearer $TOKEN"

# Ответ должен содержать данные вашего пользователя
```

### 4.4 Обновление/аутентификация пользователя (новые endpoints)

```bash
# Обновить токены по refresh_token
curl -X POST http://localhost/api/v1/users/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "<REFRESH_TOKEN>"
  }'

# Выйти из системы (инвалидировать текущий токен)
curl -X POST http://localhost/api/v1/users/logout \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "revoke_all_sessions": false
  }'

# Обновить профиль пользователя (имя/фамилия)
curl -X PUT http://localhost/api/v1/users/profile \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "<USER_ID>",
    "first_name": "Ivan",
    "last_name": "Ivanov"
  }'

# Сменить пароль
curl -X PUT http://localhost/api/v1/users/password \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "<USER_ID>",
    "old_password": "Old@12345",
    "new_password": "New@12345"
  }'

# Список пользователей (только админ)
curl -X GET 'http://localhost/api/v1/users?page=0&size=20&status=active&role=user' \
  -H "Authorization: Bearer $ADMIN_TOKEN"

# Назначить роль пользователю (admin)
curl -X PUT http://localhost/api/v1/users/<USER_ID>/role \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "role": "manager"
  }'

# Забыл пароль — отправка письма
curl -X POST http://localhost/api/v1/users/forgot-password \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com"
  }'

# Сброс пароля по токену
curl -X POST http://localhost/api/v1/users/reset-password \
  -H "Content-Type: application/json" \
  -d '{
    "reset_token": "<RESET_TOKEN>",
    "new_password": "New@12345"
  }'
```

### 4.4 Создать компанию

```bash
# Требует авторизации
curl -X POST http://localhost/api/v1/companies \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My Company Inc",
    "tax_id": "12345678",
    "address": "123 Main St",
    "city": "New York"
  }'
```

---

## 📊 Шаг 5: Посмотреть метрики (2 минуты)

### Jaeger - Трейсинг запросов

```bash
# Откройте Jaeger UI
open http://localhost:16686

# В браузере:
# 1. Выберите сервис "api-gateway"
# 2. Нажмите "Find Traces"
# 3. Посмотрите полный путь вашего запроса со временем выполнения
```

**Что видите:**

- Полный trace запроса от Nginx до базы данных
- Время выполнения каждого компонента
- Все микросервисы которые были задействованы

### Grafana - Визуализация метрик

```bash
# Откройте Grafana
open http://localhost:3000

# Логин: admin / admin

# Смотрите готовые dashboard:
# - API Gateway Overview
# - Service Health
# - Database Performance
```

---

## 🛑 Остановить сервисы

```bash
# Остановить (контейнеры остаются)
docker-compose -f docker-compose.prod.yml stop

# Удалить контейнеры
docker-compose -f docker-compose.prod.yml down

# Удалить контейнеры и volumes (⚠️ удаляет данные!)
docker-compose -f docker-compose.prod.yml down -v
```

---

## 🔗 Основные endpoints

### Публичные (без авторизации)

```bash
# Health check
curl http://localhost/api/v1/health

# Регистрация
POST http://localhost/api/v1/users/register

# Логин
POST http://localhost/api/v1/users/login

# Документация
GET http://localhost/swagger/index.html
GET http://localhost/swagger/doc.json
```

### Защищённые (требуют токен)

```bash
# Получить текущего пользователя
GET http://localhost/api/v1/users/me
  -H "Authorization: Bearer TOKEN"

# Список компаний
GET http://localhost/api/v1/companies
  -H "Authorization: Bearer TOKEN"

# Создать компанию
POST http://localhost/api/v1/companies
  -H "Authorization: Bearer TOKEN"

# Список счётов
GET http://localhost/api/v1/invoices
  -H "Authorization: Bearer TOKEN"

# Список каталогов
GET http://localhost/api/v1/catalog
  -H "Authorization: Bearer TOKEN"

# Список банковских счётов
GET http://localhost/api/v1/bank-accounts
  -H "Authorization: Bearer TOKEN"
```

---

## ⚙️ Конфигурация

### Где находятся файлы

```
.
├── docker-compose.prod.yml     ← Production стек (используем)
├── docker-compose.yml          ← Основной стек
├── docker-compose.dev.yml      ← Development стек
├── services/
│   ├── api-gateway/            ← REST API маршрутизатор
│   ├── user-server/            ← Сервис пользователей
│   ├── company-server/         ← Сервис компаний
│   ├── invoice-server/         ← Сервис счётов
│   ├── catalog-server/         ← Сервис каталога
│   ├── bank-account-server/    ← Сервис банковских счётов
│   ├── foreign-company-server/ ← Сервис иностранных компаний
│   └── nginx-proxy/            ← Nginx reverse proxy ⭐ НОВОЕ
├── proto/                      ← gRPC протоколы
├── monitoring/                 ← Prometheus, Grafana конфиги
└── README.md                   ← Полная документация
```

### Переменные окружения

Все сервисы используют `.env.example` как шаблон:

```bash
# Посмотрите доступные переменные
cat services/api-gateway/.env.example

# Создайте свой .env (опционально)
cp services/api-gateway/.env.example services/api-gateway/.env
```

---

## 🐛 Частые проблемы

### ❌ "Port already in use"

```bash
# Найдите и убейте процесс на порту 80
lsof -i :80
kill -9 <PID>

# Или используйте другие порты в docker-compose.yml
```

### ❌ "Container exited with code 1"

```bash
# Посмотрите детальные логи
docker-compose -f docker-compose.prod.yml logs [SERVICE_NAME]

# Пересоберите контейнер
docker-compose -f docker-compose.prod.yml build --no-cache [SERVICE_NAME]
```

### ❌ "Permission denied" для /var/run/docker.sock

```bash
# Добавьте текущего пользователя в docker группу
sudo usermod -aG docker $USER

# Или используйте sudo перед docker-compose
sudo docker-compose -f docker-compose.prod.yml up -d
```

### ❌ API возвращает 404

```bash
# Проверьте что Nginx работает
curl -v http://localhost/api/v1/health

# Проверьте API Gateway
curl http://localhost:8080/api/v1/health

# Посмотрите логи обоих
docker-compose -f docker-compose.prod.yml logs nginx-proxy api-gateway
```

---

## 📚 Что дальше?

После успешного запуска:

1. **Изучите API** → [API_DOCUMENTATION.md](./API_DOCUMENTATION.md)
2. **Настройте HTTPS** → [services/nginx-proxy/SSL-SETUP.md](./services/nginx-proxy/SSL-SETUP.md)
3. **Напишите интеграционные тесты** → [services/invoice-server/INTEGRATION_TESTS_README.md](./services/invoice-server/INTEGRATION_TESTS_README.md)
4. **Настройте мониторинг** → [monitoring/](./monitoring/)
5. **Развернитесь в K8s** → [k8s/](./k8s/)

---

## 🆘 Нужна помощь?

- 📖 **Полная документация**: [README.md](./README.md)
- 🗺️ **Указатель всей документации**: [INDEX.md](./INDEX.md)
- 🔒 **Безопасность и JWT**: [SECURITY_BEST_PRACTICES.md](./SECURITY_BEST_PRACTICES.md)
- 📊 **Мониторинг и трейсинг**: [OPENTELEMETRY_TRACING_IMPLEMENTATION.md](./OPENTELEMETRY_TRACING_IMPLEMENTATION.md)
- ⚡ **Нагрузочное тестирование**: [LOAD_TESTING_GUIDE.md](./LOAD_TESTING_GUIDE.md)

---

**Успехов! 🚀 Если всё работает, значит вы готовы к production! 🎉**
