# Диагностика ошибки загрузки каталога

## Проблема

API возвращает ошибку 500 при попытке загрузить каталог: `GET http://localhost:4000/api/catalog? 500 (Internal Server Error)`

## Возможные причины

### 1. **Бэкэнд API Gateway недоступен**

- API Gateway должен работать на `http://localhost:8080`
- Это настраивается переменной `NEXT_PUBLIC_API_GATEWAY_URL=http://localhost:8080` в `.env.local`

**Решение:**

```bash
# Запустить docker-compose с API Gateway
docker-compose -f docker-compose.dev.yml up api-gateway

# Проверить, что API Gateway отвечает
curl http://localhost:8080/health
```

### 2. **Неправильная конфигурация бэкэнда**

- Убедитесь, что бэкэнд имеет эндпоинт `/api/v1/catalog`
- Проверьте логи бэкэнда на ошибки

**Решение:**

```bash
# Проверить логи API Gateway
docker-compose -f docker-compose.dev.yml logs api-gateway

# Проверить, доступен ли эндпоинт каталога
curl -H "Authorization: Bearer YOUR_TOKEN" http://localhost:8080/api/v1/catalog
```

### 3. **Проблема с авторизацией**

- Токен может истечь или быть недействительным
- Заголовок Authorization может быть неправильно передан

**Решение:**

```bash
# Проверить health эндпоинт Next.js
curl http://localhost:4000/api/catalog/health

# Это должно вернуть информацию о статусе подключения к бэкэнду
```

### 4. **Сетевые ограничения или CORS**

- Если запрос делается из браузера, может быть проблема с CORS
- Next.js API routes должны обойти эту проблему

**Решение:**

- Проверьте консоль браузера на CORS ошибки
- Убедитесь, что бэкэнд настроен правильно для CORS

## Диагностирование

### Шаг 1: Проверить статус Next.js

```bash
# Перейти в директорию фронтэнда
cd kkm-platform

# Убедитесь, что Next.js запущен (должен быть на 3000 или 4000)
pnpm dev
```

### Шаг 2: Проверить здоровье API

```bash
# Откройте браузер и перейдите на:
http://localhost:3000/api/catalog/health
# или
http://localhost:4000/api/catalog/health

# Это должно вернуть JSON с информацией о статусе подключения
```

### Шаг 3: Проверить логи

**Логи Next.js:**

```bash
# В терминале где запущен `pnpm dev`
# Ищите строки вида:
# [Catalog API] Fetching: http://localhost:8080/api/v1/catalog?...
# [Catalog API] Backend response status: 200
```

**Логи бэкэнда:**

```bash
# Если используется Docker
docker-compose -f docker-compose.dev.yml logs -f api-gateway
```

## Типичные ошибки и решения

### Ошибка: "Cannot connect to backend at http://localhost:8080"

**Причина:** API Gateway не запущен

**Решение:**

```bash
# Запустить все сервисы
make dev-up
# или
docker-compose -f docker-compose.dev.yml up

# Убедитесь, что контейнер api-gateway запущен
docker ps | grep api-gateway
```

### Ошибка: "Backend returned 404"

**Причина:** Эндпоинт `/api/v1/catalog` не существует на бэкэнде

**Решение:**

- Проверьте, что бэкэнд имеет это маршрут
- Проверьте версию API (может быть `/api/v2/catalog` или другое)
- Обновите `NEXT_PUBLIC_API_GATEWAY_URL` в `.env.local` если URL неправильный

### Ошибка: "Backend returned 401"

**Причина:** Токен авторизации отсутствует или недействителен

**Решение:**

- Перезагрузите страницу чтобы получить новый токен
- Очистите localStorage и авторизуйтесь заново

### Ошибка: "Backend returned 500"

**Причина:** Внутренняя ошибка на бэкэнде

**Решение:**

- Проверьте логи бэкэнда для деталей ошибки
- Убедитесь, что бэкэнд полностью инициализирован
- Проверьте, что база данных готова

## Как просмотреть детальные логи

### Во время разработки

Откройте браузер консоль (F12) и посмотрите:

1. **Network tab** - запрос к `/api/catalog`

- Статус ответа
- Заголовки запроса и ответа
- Тело ответа (JSON с ошибкой)

2. **Console tab** - ищите логи вида:

- `[Catalog API] Fetching: ...`
- `[Catalog API] Backend response status: ...`
- `Ошибка загрузки каталога: ...`

### Сборка и информация об окружении

```bash
# Проверьте переменные окружения
cat kkm-platform/.env.local | grep NEXT_PUBLIC_API_GATEWAY_URL

# Проверьте, что PORT правильный
lsof -i :3000  # или :4000 если используется другой порт
lsof -i :8080  # для бэкэнда
```

## Дополнительные команды для диагностики

```bash
# Проверить, что оба сервера работают
curl -I http://localhost:3000      # Next.js (3000) или (4000)
curl -I http://localhost:8080      # API Gateway

# Проверить эндпоинт каталога напрямую
curl http://localhost:8080/api/v1/catalog

# Проверить с токеном
curl -H "Authorization: Bearer YOUR_TOKEN" http://localhost:8080/api/v1/catalog

# Проверить health через Next.js proxy
curl http://localhost:3000/api/catalog/health
```

## Файлы для проверки

1. **kkm-platform/.env.local** - конфигурация окружения
2. **kkm-platform/app/api/catalog/route.ts** - API proxy
3. **kkm-platform/store/catalog/catalogStore.ts** - логика получения данных
4. **kkm-platform/next.config.ts** - конфигурация Next.js

## Решение 500 ошибки

Если вы видите 500 ошибку:

1. Проверьте `/api/catalog/health` endpoint
2. Посмотрите логи консоли браузера (F12 → Console)
3. Посмотрите терминал где запущен `pnpm dev`
4. Посмотрите логи бэкэнда `docker-compose logs api-gateway`
5. Убедитесь что `NEXT_PUBLIC_API_GATEWAY_URL` правильный

## Контрольный список

- [ ] API Gateway запущен на `http://localhost:8080`
- [ ] Next.js запущен (проверить порт в консоли)
- [ ] Токен валиден и не истек
- [ ] `.env.local` содержит `NEXT_PUBLIC_API_GATEWAY_URL=http://localhost:8080`
- [ ] Бэкэнд имеет эндпоинт `/api/v1/catalog`
- [ ] Нет CORS ошибок в консоли браузера
