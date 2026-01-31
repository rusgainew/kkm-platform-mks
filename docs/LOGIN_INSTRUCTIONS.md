# Инструкция по входу в систему

## Доступ к приложению

**URL**: http://localhost/

## Учетные данные

Пользователь уже создан в системе:

- **Email**: `admin@example.com`
- **Пароль**: `Admin@2024Sec#secure99`
- **Роль**: `user`
- **User ID**: `06c7f417-85dc-405b-97bc-1bbe371ad799`

## Использование API

### Логин

```bash
curl -X POST http://localhost/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "Admin@2024Sec#secure99"
  }'
```

**Ответ**:

```json
{
  "access_token": "eyJ...",
  "refresh_token": "eyJ...",
  "expires_in": 3600,
  "refresh_expires_in": 86400,
  "token_type": "Bearer",
  "user": {
    "user_id": "06c7f417-85dc-405b-97bc-1bbe371ad799",
    "email": "admin@example.com",
    "first_name": "Admin",
    "last_name": "User",
    "role": "user"
  }
}
```

### Регистрация нового пользователя

```bash
curl -X POST http://localhost/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newuser@example.com",
    "password": "SecurePass123!",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

### Использование токена

```bash
curl -X GET http://localhost/api/v1/users/profile \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## Swagger UI

**URL**: http://localhost/swagger/index.html

Swagger UI содержит полную документацию по всем доступным API endpoints.

## Endpoints

### Пользователи

- `POST /api/v1/users/register` - Регистрация
- `POST /api/v1/users/login` - Логин
- `GET /api/v1/users/profile` - Получить профиль (требует токен)
- `POST /api/v1/users/refresh` - Обновить токен

### Health Check

- `GET /api/v1/health` - Проверка состояния API
- `GET /health` - Проверка состояния Nginx

## Токены

- **Access Token**: Срок действия 1 час (3600 секунд)
- **Refresh Token**: Срок действия 24 часа (86400 секунд)

Когда access token истечет, используйте refresh token для получения нового access token.

## Предупреждения браузера

### Font Preload Warnings

Предупреждения о preload шрифтов (woff2) не являются критическими:

```
The resource was preloaded using link preload but not used within a few seconds...
```

Это нормальное поведение Next.js при оптимизации загрузки шрифтов. Шрифты загружаются корректно.

## Troubleshooting

### 404 Not Found при логине

**Причина**: Пользователь не создан в БД  
**Решение**: Зарегистрируйте пользователя через `/api/v1/users/register`

### CORS errors

Все CORS headers настроены в nginx-proxy. Если видите CORS ошибки:

- Проверьте логи: `docker logs nginx-proxy`
- Убедитесь что запрос идёт через http://localhost/ (не http://localhost:4000)

### Токен истёк

**Ответ API**: `401 Unauthorized`  
**Решение**: Используйте refresh token для получения нового access token

## Дополнительно

Файл с учетными данными: `/user_password.txt`

```
admin@example.com
Admin@2024Sec#secure99
```
