# 🐳 Docker Compose - Руководство по запуску

Полная конфигурация Docker Compose для KKM Project MKS со всеми микросервисами, инфраструктурой и мониторингом.

## 📦 Что включено

### Инфраструктура (3 сервиса)

- **PostgreSQL 15** - Основная БД с 7 схемами
- **Redis 7** - Кеширование для Query сервисов
- **RabbitMQ 3** - Message broker для event-driven communication

### Микросервисы (13 сервисов)

#### Command Services (Write - 7 сервисов)

- `user-server` (50051) - Управление пользователями
- `company-server` (50052) - Управление компаниями
- `catalog-server` (50053) - Управление каталогом товаров
- `invoice-server` (50054) - Управление счетами
- `bank-account-server` (50055) - Управление банковскими счетами
- `foreign-company-server` (50056) - Управление иностранными компаниями
- `document-server` (50057) - Управление документами

#### Query Services (Read - CQRS, 6 сервисов)

- `user-query-server` (50061) - Чтение пользователей
- `catalog-query-server` (50063) - Чтение каталога
- `invoice-query-server` (50064) - Чтение счетов
- `bank-account-query-server` (50065) - Чтение банковских счетов
- `foreign-company-query-server` (50066) - Чтение иностранных компаний
- `document-query-server` (50067) - Чтение документов

### API & Frontend (3 сервиса)

- `api-gateway` (8080, 9090) - HTTP/gRPC Gateway
- `nginx-proxy` (80, 443) - Reverse Proxy с SSL
- `kkm-platform` (4000) - Next.js фронтенд

### Мониторинг (3 сервиса)

- `prometheus` (9091) - Сбор метрик
- `grafana` (3000) - Визуализация метрик
- `jaeger` (16686) - Distributed tracing

## 🚀 Быстрый старт

### 1. Запуск всех сервисов

```bash
# Запустить все сервисы
docker compose up -d

# Следить за логами
docker compose logs -f

# Следить за конкретным сервисом
docker compose logs -f api-gateway
```

### 2. Проверка статуса

```bash
# Статус всех контейнеров
docker compose ps

# Health check API Gateway
curl http://localhost/api/v1/health

# Health check Prometheus
curl http://localhost:9091/-/healthy

# Health check Grafana
curl http://localhost:3000/api/health
```

### 3. Остановка

```bash
# Остановить все сервисы
docker compose down

# Остановить и удалить volumes
docker compose down -v
```

## 🔧 Управление отдельными сервисами

```bash
# Запустить только инфраструктуру
docker compose up -d postgres redis rabbitmq

# Запустить Command сервисы
docker compose up -d user-server company-server catalog-server invoice-server bank-account-server foreign-company-server document-server

# Запустить Query сервисы
docker compose up -d user-query-server catalog-query-server invoice-query-server bank-account-query-server foreign-company-query-server document-query-server

# Перезапустить конкретный сервис
docker compose restart api-gateway

# Пересобрать и запустить сервис
docker compose up -d --build user-server
```

## 📊 Мониторинг

### Prometheus

- URL: http://localhost:9091
- Метрики всех микросервисов доступны на портах 9101-9114

### Grafana

- URL: http://localhost:3000
- Логин: `admin`
- Пароль: `admin`
- Дашборды:
  - Overview Dashboard - общая статистика
  - Detailed Dashboard - детальные метрики сервисов
  - Nginx Dashboard - метрики прокси

### Jaeger

- URL: http://localhost:16686
- Трейсинг всех gRPC запросов между сервисами

## 🔐 Доступ к сервисам

### API через Nginx (рекомендуется)

```bash
# HTTP
curl http://localhost/api/v1/health

# Swagger UI
open http://localhost/swagger/index.html
```

### Прямой доступ к API Gateway (для разработки)

```bash
curl http://localhost:8080/api/v1/health
```

### RabbitMQ Management

- URL: http://localhost:15672
- Логин: `kkm_user`
- Пароль: `kkm_password`

### PostgreSQL

```bash
# Подключение через psql
docker compose exec postgres psql -U kkm_user -d kkm_db

# Или извне
psql -h localhost -p 5432 -U kkm_user -d kkm_db
```

### Redis

```bash
# Подключение через redis-cli
docker compose exec redis redis-cli

# Или извне
redis-cli -h localhost -p 6379
```

## 🔄 Обновление сервисов

```bash
# Пересобрать все образы
docker compose build

# Пересобрать без кеша
docker compose build --no-cache

# Пересобрать конкретный сервис
docker compose build user-server

# Применить изменения
docker compose up -d --build
```

## 🐛 Отладка

```bash
# Просмотр логов всех сервисов
docker compose logs

# Последние 100 строк логов
docker compose logs --tail=100

# Логи с временными метками
docker compose logs -t

# Интерактивный shell в контейнере
docker compose exec user-server sh

# Просмотр используемых ресурсов
docker compose stats

# Инспекция конкретного сервиса
docker compose exec postgres env
```

## 📁 Volumes и данные

Volumes создаются автоматически:

- `postgres_data` - Данные PostgreSQL
- `redis_data` - Данные Redis
- `rabbitmq_data` - Данные RabbitMQ
- `prometheus_data` - Метрики Prometheus
- `grafana_data` - Дашборды и настройки Grafana

```bash
# Список volumes
docker volume ls | grep kkm-project-mks

# Удалить все volumes (ОСТОРОЖНО!)
docker compose down -v
```

## 🌐 Сеть

Все сервисы работают в одной сети `kkm-network` (172.20.0.0/16):

- Сервисы могут обращаться друг к другу по именам контейнеров
- Например: `postgres:5432`, `redis:6379`, `user-server:50051`

```bash
# Проверка сети
docker network inspect kkm-project-mks_kkm-network
```

## 🔧 Переменные окружения

Создайте `.env` файл на основе `.env.example`:

```bash
cp .env.example .env
```

Основные переменные:

- `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`
- `REDIS_URL`
- `RABBITMQ_URL`
- `NEXT_PUBLIC_API_URL`

## 🎯 Production режим

Для production используйте отдельный compose файл:

```bash
# Создайте docker-compose.prod.yml с production настройками
docker compose -f docker-compose.prod.yml up -d
```

Рекомендации для production:

- Используйте внешние управляемые БД (AWS RDS, Google Cloud SQL)
- Настройте SSL сертификаты для Nginx
- Настройте автоматическое резервное копирование volumes
- Используйте Docker secrets для паролей
- Настройте ограничения ресурсов (CPU, Memory)
- Используйте health checks для всех сервисов

## 🧪 Тестирование

```bash
# Запустить сервисы для тестов
docker compose up -d postgres redis rabbitmq

# Подождать готовности
sleep 10

# Запустить тесты
make test

# Остановить
docker compose down
```

## 📝 Полезные команды

```bash
# Очистка неиспользуемых ресурсов Docker
docker system prune -a

# Пересоздать все сервисы
docker compose up -d --force-recreate

# Масштабирование сервиса (если не указан container_name)
docker compose up -d --scale invoice-query-server=3

# Экспорт конфигурации
docker compose config > docker-compose.resolved.yml
```

## 🆘 Устранение неполадок

### Порт уже занят

```bash
# Найти процесс на порту 5432
sudo lsof -i :5432

# Или через ss
ss -tulpn | grep :5432
```

### Сервис не стартует

```bash
# Проверить логи
docker compose logs service-name

# Проверить health check
docker inspect --format='{{json .State.Health}}' container-name | jq
```

### Проблемы с сетью

```bash
# Пересоздать сеть
docker compose down
docker network prune
docker compose up -d
```

## 📚 Дополнительные ресурсы

- [Docker Compose документация](https://docs.docker.com/compose/)
- [README.md проекта](./README.md)
- [QUICKSTART.md](./QUICKSTART.md)
