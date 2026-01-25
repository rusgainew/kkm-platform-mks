# 🚀 KKM Project MKS - Микросервисная архитектура

**Версия:** 1.0 Production Ready  
**Дата обновления:** 10 января 2026  
**Статус:** ✅ Полностью развёрнута с Nginx Reverse Proxy

---

## 📋 Содержание

- [Архитектура](#архитектура)
- [Быстрый старт](#быстрый-старт)
- [Порты и сервисы](#порты-и-сервисы)
- [API документация](#api-документация)
- [Мониторинг](#мониторинг)
- [Безопасность](#безопасность)
- [Развёртывание](#развёртывание)

---

## 🏗️ Архитектура

### Общая схема

```
┌─────────────────────────────────────────────────────────────────┐
│                       Internet / Клиент                          │
└────────────────────────┬────────────────────────────────────────┘
                         │
                    HTTP/HTTPS (порт 80/443)
                         │
        ┌────────────────▼────────────────┐
        │    Nginx Reverse Proxy          │
        │  (Rate Limiting, Gzip, SSL)     │
        │         :80/:443                │
        └────────────────┬────────────────┘
                         │
                   (внутренняя сеть)
                         │
        ┌────────────────▼────────────────┐
        │      API Gateway                │
        │    (HTTP/gRPC Router)           │
        │        :8080/:9090              │
        └────────────────┬────────────────┘
                         │
         ┌───────────────┼───────────────┐
         │               │               │
    ┌────▼────┐  ┌──────▼──────┐  ┌─────▼────┐
    │ Command │  │Query (CQRS) │  │Monitoring│
    │Services │  │  Services   │  │Services  │
    └────────┘  └─────────────┘  └──────────┘
         │
    ┌────▼────────────────────┐
    │  Data Layer             │
    │ - PostgreSQL (5432)     │
    │ - Redis (6379)          │
    │ - RabbitMQ (5672)       │
    └─────────────────────────┘
```

### Компоненты

| Компонент         | Порт        | Описание                                                                      |
| ----------------- | ----------- | ----------------------------------------------------------------------------- |
| **Nginx Proxy**   | 80, 443     | Reverse proxy, rate limiting, SSL/TLS                                         |
| **API Gateway**   | 8080        | HTTP маршрутизатор → microservices                                            |
| **Microservices** | 50051-50059 | gRPC сервисы (user, company, invoice, catalog, bank-account, foreign-company) |
| **PostgreSQL**    | 5432        | Основная БД (6 схем)                                                          |
| **Redis**         | 6379        | Кеш и сессии                                                                  |
| **RabbitMQ**      | 5672        | Message broker (AMQP)                                                         |
| **Prometheus**    | 9091        | Метрики                                                                       |
| **Jaeger**        | 16686       | Трейсинг запросов                                                             |
| **Grafana**       | 3000        | Визуализация метрик                                                           |

---

## ⚡ Быстрый старт

### Требования

- Docker 20.10+
- Docker Compose 2.0+
- Go 1.24+ (для локальной разработки)
- curl или grpcurl (для тестирования)

### 1️⃣ Запустить весь стек

```bash
# Перейти в корень проекта
cd /path/to/kkm-project-mks

# Запустить production stack с Nginx
docker-compose -f docker-compose.prod.yml up -d

# Проверить статус контейнеров
docker-compose -f docker-compose.prod.yml ps
```

### 2️⃣ Проверить здоровье системы

```bash
# Health check через Nginx
curl http://localhost/api/v1/health

# Ответ должен быть:
# {"status":"healthy","time":1768056504}
```

### 3️⃣ Доступ к сервисам

**Через Nginx (рекомендуется для production):**

```bash
# HTTP
curl http://localhost/api/v1/users/register
curl http://localhost/api/v1/catalog/:id
curl http://localhost/api/v1/invoices

# HTTPS (после настройки сертификатов)
curl https://localhost/api/v1/health
```

**Прямой доступ (только для локальной разработки):**

```bash
# API Gateway напрямую
curl http://localhost:8080/api/v1/health

# Prometheus metrics
curl http://localhost:9091/metrics

# Jaeger UI
open http://localhost:16686
```

### 4️⃣ Документация API

```bash
# Swagger UI
open http://localhost/swagger/index.html

# OpenAPI JSON
curl http://localhost/swagger/doc.json | jq

# Alternative docs endpoint
open http://localhost/api/v1/docs/index.html
```

---

## 🔌 Порты и сервисы

### Публичные порты (через Nginx)

| Сервис          | Порт | URL                  | Описание                  |
| --------------- | ---- | -------------------- | ------------------------- |
| **Nginx HTTP**  | 80   | `http://localhost/`  | Все API запросы идут сюда |
| **Nginx HTTPS** | 443  | `https://localhost/` | SSL/TLS (настраивается)   |

### Внутренние порты (Docker сеть, недоступны снаружи)

| Сервис                  | Порт        | Доступ                                  |
| ----------------------- | ----------- | --------------------------------------- |
| API Gateway             | 8080        | Только из Nginx                         |
| User Service            | 50051       | Только из API Gateway                   |
| Company Service         | 50052       | Только из API Gateway                   |
| Catalog Service         | 50053       | Только из API Gateway                   |
| Bank Account Service    | 50054       | Только из API Gateway                   |
| Foreign Company Service | 50056       | Только из API Gateway                   |
| Invoice Service         | 50055       | Только из API Gateway                   |
| Query Services          | 50057-50059 | Только из API Gateway                   |
| PostgreSQL              | 5432        | Только из микросервисов                 |
| Redis                   | 6379        | Только из микросервисов                 |
| RabbitMQ                | 5672        | Только из микросервисов                 |
| Prometheus              | 9091        | Только из Grafana                       |
| Jaeger                  | 16686       | `http://localhost:16686` (вебинтерфейс) |
| Grafana                 | 3000        | `http://localhost:3000`                 |

### Примеры вызовов через Nginx

```bash
# ✅ Health check
curl http://localhost/api/v1/health

# ✅ Swagger документация
curl http://localhost/swagger/doc.json

# ✅ Метрики Prometheus (вместе с логин/пароль при необходимости)
curl http://localhost/metrics

# ✅ Создать пользователя
curl -X POST http://localhost/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@example.com","password":"pass123"}'

# ✅ Получить компанию (требует токен)
curl http://localhost/api/v1/companies/550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## 📚 API документация

### Доступные endpoints

- **Authentication**: `/api/v1/users/register`, `/api/v1/users/login`
- **Companies**: `/api/v1/companies` (CRUD)
- **Invoices**: `/api/v1/invoices` (CRUD + sign, revoke, accept/reject)
- **Catalog**: `/api/v1/catalog` (CRUD + query)
- **Bank Accounts**: `/api/v1/bank-accounts` (CRUD + query)
- **Foreign Companies**: `/api/v1/foreign-companies` (CRUD)
- **Health**: `/api/v1/health`, `/api/v1/ready`

### Полная документация

Смотрите [API_DOCUMENTATION.md](./API_DOCUMENTATION.md) для:

- Полного списка всех endpoints
- Описания параметров запроса/ответа
- Примеров использования gRPC и REST
- Кодов ошибок и обработки исключений

---

## 📊 Мониторинг

### Jaeger - Распределённый трейсинг

```bash
# Открыть Jaeger UI
open http://localhost:16686

# Доступны сервисы:
# - api-gateway
# - user-server
# - company-server
# - catalog-server
# - invoice-server
# - и другие...
```

**Использование:**

1. Выберите сервис из dropdown
2. Отправьте запрос через Nginx
3. Посмотрите полный trace с временем выполнения каждого компонента

### Prometheus - Сбор метрик

```bash
# Доступ к метрикам API Gateway
curl http://localhost/metrics

# Основные метрики:
# - api_gateway_http_requests_total
# - api_gateway_http_request_duration_seconds
# - api_gateway_active_connections
# - process_resident_memory_bytes
```

### Grafana - Визуализация

```bash
# Открыть Grafana
open http://localhost:3000

# Credentials (по умолчанию):
# Username: admin
# Password: admin

# Доступные dashboards:
# - API Gateway Overview
# - Service Health
# - Database Performance
# - Message Queue Stats
```

---

## 🔒 Безопасность

### Nginx - Rate Limiting

Nginx настроен с rate limiting для защиты от DDoS:

```
- General API: 200 req/sec (burst 20)
- API endpoints: 100 req/sec
- Auth endpoints: 5 req/sec (встроенное в API Gateway)
```

При превышении лимита возвращается **HTTP 429 (Too Many Requests)**.

### SSL/TLS (HTTPS)

Nginx готов к HTTPS. Для включения:

1. Поместите сертификат и ключ в `/etc/nginx/certs/`:

   ```bash
   cp /path/to/cert.pem /etc/nginx/certs/
   cp /path/to/key.pem /etc/nginx/certs/
   ```

2. Раскомментируйте HTTPS блок в `services/nginx-proxy/nginx.conf`

3. Перезагрузите Nginx:
   ```bash
   docker-compose -f docker-compose.prod.yml restart nginx-proxy
   ```

Смотрите [services/nginx-proxy/SSL-SETUP.md](./services/nginx-proxy/SSL-SETUP.md) для подробной инструкции по Let's Encrypt.

### JWT Authentication

- **Жизненный цикл**: 1 час (configurable)
- **Ротация ключей**: Автоматическая каждые 7 дней
- **Формат**: Bearer token в заголовке `Authorization`

```bash
# Пример с авторизацией
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  http://localhost/api/v1/companies
```

Смотрите [SECURITY_BEST_PRACTICES.md](./SECURITY_BEST_PRACTICES.md) для полной информации.

---

## 🚀 Развёртывание

### Production Docker Compose

```bash
# Запуск в production режиме
docker-compose -f docker-compose.prod.yml up -d

# Просмотр логов
docker-compose -f docker-compose.prod.yml logs -f [SERVICE_NAME]

# Масштабирование
docker-compose -f docker-compose.prod.yml up -d --scale user-server=3

# Обновление образов
docker-compose -f docker-compose.prod.yml build --no-cache

# Остановка и удаление
docker-compose -f docker-compose.prod.yml down
```

### Development стек

```bash
# Быстрая разработка с hot-reload
docker-compose -f docker-compose.dev.yml up -d

# Или запустить одну услугу локально
cd services/company-server
go run cmd/main.go
```

### Kubernetes (готовится)

K8s manifests находятся в `k8s/` директории. Используются для production deployment с:

- Horizontal Pod Autoscaling
- Service Mesh (Istio)
- Ingress с SSL/TLS
- ConfigMaps и Secrets для конфигурации

---

## 📖 Дополнительная документация

| Документ                                                                 | Описание                      |
| ------------------------------------------------------------------------ | ----------------------------- |
| [QUICKSTART.md](./QUICKSTART.md)                                         | 10-минутный гайд для новичков |
| [API_DOCUMENTATION.md](./API_DOCUMENTATION.md)                           | Полная API справка            |
| [SECURITY_BEST_PRACTICES.md](./SECURITY_BEST_PRACTICES.md)               | Безопасность и JWT            |
| [PRODUCTION_DEPLOYMENT.md](./PRODUCTION_DEPLOYMENT.md)                   | Production deployment guide   |
| [LOAD_TESTING_GUIDE.md](./LOAD_TESTING_GUIDE.md)                         | Нагрузочное тестирование      |
| [services/nginx-proxy/README.md](./services/nginx-proxy/README.md)       | Nginx proxy конфигурация      |
| [services/nginx-proxy/SSL-SETUP.md](./services/nginx-proxy/SSL-SETUP.md) | Настройка HTTPS               |
| [INDEX.md](./INDEX.md)                                                   | Полный указатель документации |

---

## 🛠️ Trouble Shooting

### Контейнер не стартует

```bash
# Проверьте логи
docker-compose -f docker-compose.prod.yml logs nginx-proxy

# Перестройте образ
docker-compose -f docker-compose.prod.yml build --no-cache nginx-proxy

# Проверьте конфиг nginx
docker exec nginx-proxy nginx -t
```

### API недоступен

```bash
# Проверьте health check
curl -v http://localhost/api/v1/health

# Проверьте API Gateway
docker-compose -f docker-compose.prod.yml logs api-gateway

# Проверьте сетевое соединение
docker network inspect kkm-project-mks_app-network
```

### Rate limiting срабатывает неправильно

```bash
# Посмотрите логи nginx
docker-compose -f docker-compose.prod.yml logs nginx-proxy | grep "limiting requests"

# Проверьте конфиг лимитов
curl -s http://localhost/swagger/doc.json | grep -i rate
```

### Swagger показывает пустую страницу

```bash
# Проверьте доступность документации
curl -v http://localhost/swagger/index.html

# Проверьте статические файлы
docker exec api-gateway ls -la /home/api-gateway/docs/
```

---

## 📞 Контакты и поддержка

- **Issues**: GitHub Issues (если используется GitHub)
- **Email**: support@kkm-project.com
- **Documentation**: [INDEX.md](./INDEX.md)
- **Slack**: #kkm-project (если используется)

---

## 📝 История изменений

### Версия 1.0 (10 января 2026)

- ✅ Nginx Reverse Proxy добавлен (rate limiting, SSL templates)
- ✅ API Gateway документация в Swagger UI
- ✅ Production Docker Compose с 16 сервисами
- ✅ Полный мониторинг (Jaeger, Prometheus, Grafana)
- ✅ JWT авторизация и ротация ключей
- ✅ Безопасность (HTTPS ready, security headers)

---

**Последнее обновление:** 10 января 2026  
**Автор:** KKM Project Team  
**Лицензия:** Apache 2.0
