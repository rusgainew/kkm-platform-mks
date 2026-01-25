# Улучшения Swagger Документации API Gateway

## 📋 Обзор

Проведена доработка Swagger/OpenAPI документации для API Gateway с улучшением аннотаций, добавлением недостающих endpoint'ов и обновлением конфигурации.

## ✅ Выполненные улучшения

### 1. **Добавлена Security Definition для JWT**

**Файл:** `cmd/api/main.go`

Добавлены аннотации для JWT аутентификации:

```go
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
```

Теперь Swagger UI корректно отображает поле для ввода JWT токена.

### 2. **Добавлена документация для SignInvoice**

**Файл:** `internal/interfaces/http/invoice_handler.go`

Добавлена полная swagger аннотация для метода подписи счета:

```go
// @Summary Sign invoice
// @Description Sign an existing invoice with digital signature
// @Tags Invoices
// @Accept json
// @Produce json
// @Param id path string true "Invoice ID (UUID)"
// @Param body body object{signature_data=string} true "Signature data"
// @Success 200 {object} models.APIResponse{data=models.Invoice} "Invoice signed successfully"
// @Failure 400 {object} models.APIResponse "Invalid input"
// @Failure 500 {object} models.APIResponse "Internal server error"
// @Router /invoices/{id}/sign [post]
// @Security BearerAuth
```

### 3. **Улучшены описания endpoints**

#### Users

- **POST /users/register** - добавлено описание про минимальную длину пароля (8 символов)
- **POST /users/login** - улучшено описание про возврат JWT токенов
- **GET /users/me** - добавлено описание про извлечение пользователя из JWT
- **GET /users/{id}** - стандартизировано описание, исправлен путь роутера

#### Companies

- **POST /companies** - добавлено указание требуемых ролей (admin/manager)
- Все endpoints компаний теперь содержат четкие описания требований безопасности

#### Invoices

- **POST /invoices** - добавлено указание требуемых ролей
- **POST /invoices/{id}/sign** - полностью документирован endpoint подписи
- Улучшены descriptions для всех операций с счетами

#### Foreign Companies

- **POST /foreign-companies** - добавлено описание формата country_code (ISO 3166-1 alpha-2)
- **GET /foreign-companies/{id}** - стандартизировано описание
- Исправлены пути роутеров с `/api/v1/...` на относительные пути

### 4. **Обновлена общая информация API**

**Файл:** `cmd/api/main.go`

```go
// @title KKM API Gateway
// @version 1.0
// @description Microservices API Gateway with command and query services. Supports CQRS pattern with separate command and query endpoints.
// @contact.email support@kkm-project.com
// @schemes http https
```

Добавлен:

- Email для контактов
- Поддержка HTTPS схемы
- Расширенное описание с упоминанием CQRS

### 5. **Регенерирована документация**

Выполнена команда:

```bash
swag init -g cmd/api/main.go -o docs
```

Сгенерированы файлы:

- `docs/docs.go` - Go код для встраивания
- `docs/swagger.json` - JSON спецификация
- `docs/swagger.yaml` - YAML спецификация

## 📊 Статистика

### Покрытие Endpoints

| Группа                  | Endpoints | Задокументировано | Статус   |
| ----------------------- | --------- | ----------------- | -------- |
| Users                   | 4         | 4                 | ✅       |
| Companies               | 4         | 4                 | ✅       |
| Invoices (Command)      | 5         | 5                 | ✅       |
| Invoices (Query)        | 4         | 4                 | ✅       |
| Catalogs (Command)      | 4         | 4                 | ✅       |
| Catalogs (Query)        | 2         | 2                 | ✅       |
| Bank Accounts (Command) | 4         | 4                 | ✅       |
| Bank Accounts (Query)   | 2         | 2                 | ✅       |
| Foreign Companies       | 4         | 4                 | ✅       |
| Health                  | 2         | 2                 | ✅       |
| **ИТОГО**               | **35**    | **35**            | **100%** |

### Модели данных

Задокументировано моделей: **12**

- `APIResponse` - общий формат ответа
- `APIError` - структура ошибок
- `MetaData` - метаданные пагинации
- `User` - пользователь
- `AuthResponse` - ответ аутентификации
- `Company` - компания
- `Invoice` - счет-фактура
- `Catalog` - элемент каталога
- `BankAccount` - банковский счет
- `ForeignCompany` - иностранная компания
- `RegisterRequest` - запрос регистрации
- `LoginRequest` - запрос входа

## 🚀 Использование

### Просмотр документации

После запуска API Gateway документация доступна по адресам:

```
http://localhost:8080/api/v1/docs/index.html
http://localhost:8080/swagger/index.html
```

### Пример использования с JWT

1. Получите токен через `/users/login`
2. В Swagger UI нажмите кнопку "Authorize"
3. Введите: `Bearer <ваш_токен>`
4. Теперь можете тестировать защищенные endpoints

### Регенерация документации

При изменении swagger аннотаций выполните:

```bash
cd services/api-gateway
swag init -g cmd/api/main.go -o docs
```

## 📝 Лучшие практики

### Структура аннотаций

```go
// @Summary Краткое описание (3-5 слов)
// @Description Подробное описание с деталями
// @Tags Название группы
// @Accept json
// @Produce json
// @Param name type datatype required "description"
// @Success 200 {object} ResponseType "Success description"
// @Failure 400 {object} ErrorType "Error description"
// @Router /path [method]
// @Security BearerAuth
```

### Требования к ролям

Указывайте в description требуемые роли:

```go
// @Description Create a new company. Requires admin or manager role.
```

### Валидация параметров

Указывайте формат и ограничения:

```go
// @Param id path string true "Invoice ID (UUID)"
// @Param password body string true "Password (min 8 chars)"
```

## 🔍 Проверка

### Компиляция

```bash
cd services/api-gateway
go build -o api-gateway cmd/api/main.go
```

✅ Компиляция успешна без ошибок

### Swagger спецификация

Проверьте корректность:

- ✅ SecurityDefinitions присутствует
- ✅ Все endpoints имеют теги
- ✅ Все protected endpoints имеют `@Security BearerAuth`
- ✅ Все модели корректно определены
- ✅ Paths не содержат дублирующих `/api/v1/`

## 📖 Дополнительные ресурсы

- [Swagger/OpenAPI 2.0 Spec](https://swagger.io/specification/v2/)
- [swaggo/swag documentation](https://github.com/swaggo/swag)
- [Gin-Swagger integration](https://github.com/swaggo/gin-swagger)

## 🎯 Следующие шаги

Рекомендуемые улучшения:

1. **Добавить примеры запросов/ответов** - использовать `@Example` для сложных структур
2. **Версионирование API** - подготовка к v2 с сохранением обратной совместимости
3. **Документация ошибок** - расширенное описание error codes
4. **Rate limiting** - документировать лимиты в descriptions
5. **Webhooks** - если планируется, добавить документацию

---

**Дата обновления:** 10 января 2026  
**Версия API:** 1.0  
**Статус:** ✅ Завершено
