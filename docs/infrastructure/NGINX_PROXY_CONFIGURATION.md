# Конфигурация Nginx Proxy

## Обзор

Nginx Proxy настроен как единая точка входа для всего приложения:

- **Фронтенд (Next.js)**: Доступен на `http://localhost/`
- **Backend API**: Доступен на `http://localhost/api/`
- **Swagger UI**: Доступен на `http://localhost/swagger/`

## Архитектура

```
┌─────────────────────────────────────────────┐
│           Nginx Proxy (:80/:443)            │
│                                             │
│  ┌─────────────────────────────────────┐   │
│  │  Location Routing:                   │   │
│  │  • /api/*      → api-gateway:8080   │   │
│  │  • /swagger/*  → api-gateway:8080   │   │
│  │  • /docs/*     → api-gateway:8080   │   │
│  │  • /*          → kkm-platform:4000  │   │
│  └─────────────────────────────────────┘   │
└─────────────────────────────────────────────┘
```

## Маршрутизация

### API Endpoints

- **Path**: `/api/*`
- **Upstream**: `api-gateway:8080`
- **Rate Limit**: 1000 req/s (burst: 50)
- **Timeout**: 30s

### Swagger/Docs

- **Path**: `/swagger/*`, `/docs/*`, `/redoc/*`
- **Upstream**: `api-gateway:8080`
- **No Rate Limit**

### Frontend (Next.js)

- **Path**: `/*` (все остальное)
- **Upstream**: `kkm-platform:4000`
- **Rate Limit**: 5000 req/s (burst: 20)
- **Timeout**: 60s
- **WebSocket Support**: Да (для Hot Module Replacement)

### Static Assets

- **Next.js Static**: `/_next/static/*` (cache: 1 year)
- **Public Files**: `/public/*` (cache: 1 day)
- **API Static**: `*.js`, `*.css`, etc. от api-gateway (cache: 1 day)

## Конфигурация

### HTTP (Port 80)

Используется для разработки и тестирования:

```nginx
location /api/ {
    proxy_pass http://api_gateway;
    # ... headers, timeouts
}

location / {
    proxy_pass http://kkm_platform;
    # ... headers, timeouts, websocket support
}
```

### HTTPS (Port 443)

Используется для продакшена с SSL сертификатами:

- Сертификаты: `/etc/nginx/certs/server.crt` и `server.key`
- Протоколы: TLSv1.2, TLSv1.3
- HSTS включен

## Environment Variables

### kkm-platform

```yaml
NODE_ENV: production
NEXT_PUBLIC_API_URL: http://localhost/api/v1
```

API URL теперь указывает на nginx-proxy, который проксирует запросы к api-gateway.

## Docker Compose Dependencies

```yaml
kkm-platform:
  depends_on:
    api-gateway:
      condition: service_healthy

nginx-proxy:
  depends_on:
    api-gateway:
      condition: service_healthy
    kkm-platform:
      condition: service_started
```

## Endpoints

### Public Endpoints

- **Frontend**: http://localhost/
- **API Health**: http://localhost/api/v1/health
- **Swagger UI**: http://localhost/swagger/index.html
- **Nginx Health**: http://localhost/health

### Internal Endpoints

- **Metrics**: http://localhost/metrics (restricted to Docker network)
- **Nginx Status**: http://localhost/nginx_status (restricted)

## Performance Optimizations

1. **Keepalive Connections**: 32 connections pool к upstream
2. **Response Buffering**: Отключен для streaming
3. **Static File Caching**:
   - Next.js static: 1 год
   - Public files: 1 день
   - API assets: 1 день
4. **Gzip Compression**: Уровень 5, минимум 1KB
5. **Connection Pooling**: least_conn балансировка
6. **Rate Limiting**:
   - API: 1000 req/s
   - General: 5000 req/s

## CORS Configuration

CORS headers включены для всех запросов:

```nginx
add_header 'Access-Control-Allow-Origin' '*' always;
add_header 'Access-Control-Allow-Methods' 'GET, POST, PUT, DELETE, OPTIONS' always;
add_header 'Access-Control-Allow-Headers' 'Authorization, Content-Type, X-Requested-With' always;
```

## Тестирование

```bash
# Проверка фронтенда
curl http://localhost/

# Проверка API
curl http://localhost/api/v1/health

# Проверка Swagger
curl http://localhost/swagger/index.html

# Проверка healthcheck
curl http://localhost/health
```

## Troubleshooting

### Ошибка 502 Bad Gateway

- Проверьте что api-gateway и kkm-platform запущены
- Проверьте логи nginx: `docker logs nginx-proxy`

### Ошибка 404 Not Found

- Убедитесь что маршрут правильный
- Проверьте конфигурацию nginx: `docker exec nginx-proxy nginx -t`

### Медленная загрузка

- Проверьте логи upstream сервисов
- Проверьте настройки timeout и buffering

## Файлы конфигурации

- **Nginx Config**: `services/nginx-proxy/nginx.conf`
- **SSL Certs**: `services/nginx-proxy/certs/`
- **Docker Compose**: `docker-compose.yml`
