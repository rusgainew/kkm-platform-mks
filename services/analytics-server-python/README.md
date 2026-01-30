# analytics-server-python

Python gRPC сервис аналитики. Повторяет контракт из [proto/analytics/analytics.proto](../../proto/analytics/analytics.proto) и логику Go‑сервиса.

## Быстрый старт

1. Создайте venv и установите зависимости:

- `python -m venv .venv`
- `source .venv/bin/activate`
- `pip install -r requirements.txt`

2. Сгенерируйте gRPC код из proto:

- `python scripts/generate_proto.py`

3. Запустите сервис:

- `python -m analytics_server.server`

## Переменные окружения

- `GRPC_PORT` (по умолчанию 50070)
- `METRICS_PORT` (не используется, оставлено для совместимости)
- `DB_HOST` (по умолчанию localhost)
- `DB_PORT` (по умолчанию 5432)
- `DB_USER` (по умолчанию analytics_svc)
- `DB_PASSWORD` (по умолчанию analytics_pass_2026)
- `DB_NAME` (по умолчанию analytics_db)
- `DB_SSLMODE` (по умолчанию disable)
- `DB_MAX_OPEN_CONNS` (по умолчанию 25)
- `DB_MAX_IDLE_CONNS` (по умолчанию 5)
- `DB_CONN_MAX_LIFETIME` (например `300s`, `5m`)
- `REDIS_ADDR` (по умолчанию localhost:6379)
- `REDIS_PASSWORD` (по умолчанию пусто)
- `REDIS_DB` (по умолчанию 0)
- `REDIS_TTL` (по умолчанию 300 секунд)
- `LOG_LEVEL` (по умолчанию INFO)

## Примечания

- Для запуска нужен доступ к базе с таблицей `invoices` и полями, используемыми в запросах.
- Если Redis недоступен, кеширование отключается автоматически.
