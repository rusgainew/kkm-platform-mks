# Copilot instructions (KKM Project MKS)

## Big picture

- Микросервисная платформа с CQRS: command‑сервисы (запись) и query‑сервисы (чтение) общаются по gRPC; HTTP вход идет через Nginx → API Gateway → gRPC.
- Базы раздельные «DB per service»: см. отдельные `postgres-*` контейнеры в [docker-compose.yml](docker-compose.yml).
- Наблюдаемость встроена: Prometheus/Grafana/Jaeger + метрики и health endpoints.

## Ключевые части репозитория

- Go монорепозиторий с `go.work` (несколько модулей): сервисы в services/\*, общие библиотеки в lib/, proto‑код в proto-lib/.
- API Gateway на чистой архитектуре: см. [services/api-gateway/README.md](services/api-gateway/README.md) и структуру `internal/{domain,application,infrastructure,interfaces}`.
- gRPC контракты и генерация кода: см. [proto/README.md](proto/README.md).
- Frontend (Next.js) находится в kkm-platform/ (app/, package.json).

## Как запускать (локально)

- Полный стек: `docker compose up -d` или `make docker-up` (см. [README.md](README.md)).
- Production compose: `docker compose -f docker-compose.prod.yml up -d` (см. [QUICKSTART.md](QUICKSTART.md)).
- Health через Nginx: `curl http://localhost/api/v1/health`.

## Конвенции и паттерны

- Внешний HTTP трафик всегда идет через Nginx и API Gateway; прямой доступ к gRPC‑сервисам — только для локальной разработки.
- Для gRPC клиентов/серверов используйте сгенерированный код из proto-lib; не редактируйте сгенерированные файлы вручную.
- Добавляя новый сервис в gateway, следуйте шагам из [services/api-gateway/README.md](services/api-gateway/README.md) (service → handler → routes).

## Интеграции и зависимости

- RabbitMQ используется для событий/очередей; Redis — для кэша query‑сервисов; PostgreSQL — per‑service.
- Метрики доступны через `METRICS_PORT`, трейсинг отправляется в Jaeger (см. README по сервисам).

## Где искать примеры

- Docker‑оркестрация и сети: [docker-compose.yml](docker-compose.yml).
- Список сервисов и портов: [README.md](README.md).
- Примеры gRPC контрактов: [proto/README.md](proto/README.md).
