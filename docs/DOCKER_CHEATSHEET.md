# 🐳 Docker Compose - Шпаргалка команд

## ⚡ Быстрый старт

```bash
# Запустить все сервисы
docker compose up -d

# Или с помощью Makefile
make docker-up

# Или интерактивный менеджер
./docker-manager.sh
```

## 📋 Основные команды

### Запуск и остановка

```bash
# Запустить все сервисы (детач режим)
docker compose up -d

# Запустить с просмотром логов
docker compose up

# Остановить все сервисы
docker compose down

# Остановить и удалить volumes
docker compose down -v

# Перезапустить все сервисы
docker compose restart

# Перезапустить конкретный сервис
docker compose restart api-gateway
```

### Сборка

```bash
# Собрать все образы
docker compose build

# Собрать без кеша
docker compose build --no-cache

# Собрать конкретный сервис
docker compose build user-server

# Собрать и запустить
docker compose up -d --build
```

### Статус и мониторинг

```bash
# Показать статус контейнеров
docker compose ps

# Детальный статус
docker compose ps -a

# Потребление ресурсов
docker stats

# Проверка конфигурации
docker compose config
```

### Логи

```bash
# Все логи (следить в реальном времени)
docker compose logs -f

# Последние 100 строк
docker compose logs --tail=100

# Логи конкретного сервиса
docker compose logs -f api-gateway

# Логи с временными метками
docker compose logs -t -f
```

## 🎯 Полезные команды

### Управление отдельными группами

```bash
# Только инфраструктура
docker compose up -d postgres redis rabbitmq

# Только мониторинг
docker compose up -d prometheus grafana jaeger

# Только Command сервисы
docker compose up -d user-server company-server catalog-server invoice-server

# Только Query сервисы
docker compose up -d user-query-server catalog-query-server invoice-query-server
```

### Работа с контейнерами

```bash
# Войти в контейнер
docker compose exec postgres bash
docker compose exec api-gateway sh

# Выполнить команду в контейнере
docker compose exec postgres psql -U kkm_user -d kkm_db

# Подключиться к Redis
docker compose exec redis redis-cli

# Посмотреть переменные окружения
docker compose exec user-server env
```

### Очистка

```bash
# Остановить и удалить все
docker compose down

# Удалить volumes
docker compose down -v

# Полная очистка Docker
docker system prune -a
docker volume prune
```

## 🔧 Makefile команды

```bash
# Основные
make docker-up              # Запустить все
make docker-down            # Остановить все
make docker-logs            # Логи всех сервисов
make docker-ps              # Статус контейнеров
make docker-restart         # Перезапуск всех

# Production
make prod-up                # Запуск production stack
make prod-down              # Остановка production
make prod-logs              # Логи production

# Инфраструктура
make infra-up               # Только PostgreSQL, Redis, RabbitMQ
make monitoring-up          # Только Prometheus, Grafana, Jaeger

# Сборка
make docker-build-all       # Собрать все образы
make docker-build-no-cache  # Собрать без кеша

# Специальные
make docker-health          # Проверка здоровья
make docker-stats           # Статистика ресурсов
make docker-clean           # Полная очистка
```

## 🩺 Health Checks

```bash
# API Gateway
curl http://localhost:8080/api/v1/health

# Nginx Proxy
curl http://localhost/health

# Swagger UI
curl http://localhost/swagger/index.html

# Prometheus
curl http://localhost:9091/-/healthy

# Grafana
curl http://localhost:3000/api/health

# PostgreSQL
docker compose exec postgres pg_isready -U kkm_user

# Redis
docker compose exec redis redis-cli ping

# RabbitMQ
docker compose exec rabbitmq rabbitmq-diagnostics ping
```

## 📊 Доступ к сервисам

### API & Frontend

- **API Gateway HTTP:** http://localhost:8080
- **API Gateway gRPC:** localhost:9090
- **Nginx Proxy:** http://localhost или https://localhost
- **Swagger UI:** http://localhost/swagger/index.html
- **Frontend (Next.js):** http://localhost:4000

### Инфраструктура

- **PostgreSQL:** localhost:5432 (kkm_user/kkm_password)
- **Redis:** localhost:6379
- **RabbitMQ AMQP:** localhost:5672
- **RabbitMQ Management:** http://localhost:15672 (kkm_user/kkm_password)

### Мониторинг

- **Prometheus:** http://localhost:9091
- **Grafana:** http://localhost:3000 (admin/admin)
- **Jaeger:** http://localhost:16686

### Метрики микросервисов

- **user-server:** http://localhost:9101/metrics
- **user-query-server:** http://localhost:9102/metrics
- **company-server:** http://localhost:9103/metrics
- **catalog-server:** http://localhost:9105/metrics
- **catalog-query-server:** http://localhost:9106/metrics
- **invoice-server:** http://localhost:9107/metrics
- **invoice-query-server:** http://localhost:9108/metrics
- **bank-account-server:** http://localhost:9109/metrics
- **bank-account-query-server:** http://localhost:9110/metrics
- **foreign-company-server:** http://localhost:9111/metrics
- **foreign-company-query-server:** http://localhost:9112/metrics
- **document-server:** http://localhost:9113/metrics
- **document-query-server:** http://localhost:9114/metrics

## 🐛 Отладка

### Проблемы с запуском

```bash
# Проверить логи
docker compose logs service-name

# Проверить переменные окружения
docker compose config

# Проверить сеть
docker network inspect kkm-project-mks_kkm-network

# Пересоздать контейнеры
docker compose up -d --force-recreate
```

### Порты заняты

```bash
# Найти процесс на порту
sudo lsof -i :8080
# или
ss -tulpn | grep :8080

# Убить процесс
sudo kill -9 <PID>
```

### Проблемы с volumes

```bash
# Пересоздать volumes
docker compose down -v
docker compose up -d

# Проверить volumes
docker volume ls
docker volume inspect kkm-project-mks_postgres_data
```

### Health check не проходит

```bash
# Проверить статус
docker compose ps

# Проверить логи
docker compose logs service-name

# Войти в контейнер
docker compose exec service-name sh

# Проверить внутри контейнера
curl http://localhost:8080/health
```

## 💡 Полезные трюки

### Следить за конкретными сервисами

```bash
# Несколько сервисов одновременно
docker compose logs -f api-gateway nginx-proxy postgres
```

### Запуск с переменными окружения

```bash
# Указать .env файл
docker compose --env-file .env.production up -d

# Переопределить переменные
POSTGRES_PASSWORD=newpass docker compose up -d
```

### Масштабирование (для сервисов без container_name)

```bash
# Запустить 3 инстанса query сервиса
docker compose up -d --scale invoice-query-server=3
```

### Экспорт конфигурации

```bash
# Посмотреть результирующую конфигурацию
docker compose config > resolved-compose.yml
```

### Быстрая проверка

```bash
# Одной командой: запустить, проверить, показать логи
docker compose up -d && docker compose ps && docker compose logs --tail=50
```

## 📚 Дополнительные ресурсы

- [DOCKER_COMPOSE_GUIDE.md](./DOCKER_COMPOSE_GUIDE.md) - Полное руководство
- [DOCKER_SETUP_SUMMARY.md](./DOCKER_SETUP_SUMMARY.md) - Сводка реализации
- [README.md](./README.md) - Общая информация о проекте
- [Makefile](./Makefile) - Все доступные команды

## 🎬 Типичные сценарии

### Первый запуск

```bash
# 1. Создать .env файл
cp .env.example .env

# 2. Запустить все
docker compose up -d

# 3. Проверить статус
docker compose ps

# 4. Проверить health
curl http://localhost/api/v1/health
```

### Обновление кода

```bash
# 1. Пересобрать образы
docker compose build user-server

# 2. Перезапустить сервис
docker compose up -d user-server

# 3. Проверить логи
docker compose logs -f user-server
```

### Полный перезапуск

```bash
# Остановить все, удалить volumes, собрать заново, запустить
docker compose down -v && \
docker compose build && \
docker compose up -d
```

### Разработка (только инфраструктура)

```bash
# Запустить только БД и кеши
make infra-up

# Запустить свой сервис локально
cd services/user-server && go run cmd/main.go
```

---

**💡 Совет:** Используйте `./docker-manager.sh` для интерактивного управления всеми сервисами!
