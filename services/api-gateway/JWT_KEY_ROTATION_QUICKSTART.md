# JWT Key Rotation: Быстрый старт

## 🚀 Включить ротацию в 3 шага

### 1. Добавить переменные окружения в .env

```bash
# JWT Authentication
JWT_SECRET=your-secure-secret-key-change-in-production
JWT_EXPIRATION=24h

# Key Rotation (опционально, включено по умолчанию)
KEY_ROTATION_ENABLED=true
KEY_ROTATION_DAYS=30
KEY_ROTATION_MAX_OLD_KEYS=5
KEY_ROTATION_KEY_LENGTH=32
```

### 2. В initializer.go - обновить инициализацию auth

```go
// pkg/di_container/initializer.go

func (i *Initializer) InitializeInfrastructure() error {
    cfg := i.container.Config()
    logger := i.container.Logger()

    // Инициализация Auth Service с rotation
    authService, err := auth.NewServiceWithKeyRotation(
        cfg.JWT.Secret,
        cfg.JWT.Expiration,
        logger,
        &auth.KeyRotationConfig{
            Enabled:      true,
            RotationDays: 30,
            MaxOldKeys:   5,
            KeyLength:    32,
        },
    )
    if err != nil {
        logger.Error("Failed to initialize auth service", zap.Error(err))
        return err
    }

    i.container.SetAuthService(authService)

    // Важно: Cleanup при shutdown
    i.container.RegisterCleanup(func() error {
        authService.Stop()
        return nil
    })

    return nil
}
```

### 3. Запустить

```bash
go run cmd/api/main.go
# или
docker-compose -f docker-compose.dev.yml up
```

## ✅ Проверить что работает

### Логи

```bash
# Должны увидеть
INFO: Key rotation loop started rotation_days=30

# При первой генерации токена
DEBUG: Token generated successfully user_id=uuid expires_at=2026-01-11...

# Каждые 30 дней
INFO: Key rotated successfully
  new_key_id=key_abc123
  total_keys=2
  next_rotation=2026-02-10...
```

### Unit тесты

```bash
# Запустить все тесты auth
go test -v ./internal/infrastructure/auth/...

# Запустить тесты key provider
go test -v ./internal/infrastructure/auth/... -run TestKeyProvider

# Запустить тесты key version
go test -v ./internal/infrastructure/auth/... -run TestKeyVersion
```

## 📊 Как это работает

### Генерация токена

```go
// Каждый новый токен подписан АКТИВНЫМ ключом
token, _ := authService.GenerateToken(ctx, userID, claims)

// В токене сохранен ID ключа (для audit trail)
// Header: { alg: HS256, typ: JWT }
// Payload: { user_id: "...", iat: 1234567890, exp: 1234671490, jti: "key_abc123" }
```

### Валидация токена

```go
// Валидация работает со ВСЕМИ валидными ключами
// Старые токены остаются валидными 2 периода (60 дней)
claims, _ := authService.ValidateToken(ctx, token)

// Процесс:
// 1. Парсим токен
// 2. Пытаемся валидировать с активным ключом
// 3. Если не подходит, пробуем старые ключи (до истечения)
// 4. Возвращаем claims если найден подходящий ключ
```

### Ротация ключей

```
День 1-30:   Key-1 Active
             ├─ Новые токены: подписаны Key-1
             └─ Валидация: Key-1

День 31-60:  Key-2 Active (ротация произошла)
             ├─ Новые токены: подписаны Key-2
             └─ Валидация: Key-1 и Key-2 (grace period)

День 61+:    Key-3 Active (Key-1 expired)
             ├─ Новые токены: подписаны Key-3
             └─ Валидация: Key-2 и Key-3
             └─ Key-1 токены: ❌ Невалидны (истекли)
```

## 🔧 Конфигурация параметры

| Параметр                    | Значение | Описание                       |
| --------------------------- | -------- | ------------------------------ |
| `KEY_ROTATION_ENABLED`      | `true`   | Включить авторотацию           |
| `KEY_ROTATION_DAYS`         | `30`     | Ротация каждые N дней          |
| `KEY_ROTATION_MAX_OLD_KEYS` | `5`      | Хранить 5 старых ключей        |
| `KEY_ROTATION_KEY_LENGTH`   | `32`     | Длина ключа в байтах (256 бит) |

### Рекомендуемые значения

**Development:**

```bash
KEY_ROTATION_ENABLED=false      # Отключено для быстрого тестирования
```

**Staging:**

```bash
KEY_ROTATION_ENABLED=true
KEY_ROTATION_DAYS=7             # Более частая ротация для тестирования
KEY_ROTATION_MAX_OLD_KEYS=5
```

**Production:**

```bash
KEY_ROTATION_ENABLED=true
KEY_ROTATION_DAYS=30            # Стандартная ротация (месячная)
KEY_ROTATION_MAX_OLD_KEYS=5
```

## 🔑 API для управления ключами

```go
// Получить KeyProvider из auth service
keyProvider := authService.GetKeyProvider()

// Получить активный ключ
activeKey, _ := keyProvider.GetActiveKey()
fmt.Println("Active Key ID:", activeKey.ID)

// Получить все валидные ключи
validKeys := keyProvider.GetValidKeys()
fmt.Println("Valid Keys:", len(validKeys))

// Получить конкретный ключ
keyVersion, _ := keyProvider.GetKeyByID("key_abc123")
fmt.Println("Key Status:", keyVersion.GetKeyStatus())

// Ручная ротация (обычно не нужна)
keyProvider.rotateKey()
```

## 📈 Мониторинг

### Логирование ротации

```bash
# Просмотреть логи ротации
docker logs api-gateway | grep "Key rotated"

# Посмотреть время следующей ротации
docker logs api-gateway | grep "next_rotation"
```

### Метрики (рекомендуется добавить)

```go
// Добавить в observability/metrics.go
keyRotationsTotal := promauto.NewCounterVec(
    prometheus.CounterOpts{
        Name: "api_gateway_key_rotations_total",
        Help: "Total number of key rotations",
    },
    []string{"status"},
)

activeKeysCount := promauto.NewGaugeVec(
    prometheus.GaugeOpts{
        Name: "api_gateway_active_keys",
        Help: "Number of active keys",
    },
    []string{"status"},
)
```

## ⚠️ Важно

1. **Не отключайте на production** — это критично для безопасности
2. **Сохраняйте старые ключи** — для валидации старых токенов (минимум 2x RotationDays)
3. **Логируйте ротации** — для аудита и отладки
4. **Тестируйте в staging** — перед использованием на production
5. **Backup секретов** — сохраняйте текущие секреты в safe place

## 🐛 Отладка

### Ключи не ротируются?

```bash
# 1. Проверить что ротация включена
echo $KEY_ROTATION_ENABLED

# 2. Проверить логи
docker logs api-gateway | grep "Key rotation"

# 3. Проверить конфиг
# KEY_ROTATION_ENABLED должен быть true

# 4. Перезагрузить контейнер
docker restart api-gateway
```

### Токены не валидируются?

```bash
# 1. Проверить что ключ еще не истек
# День создания + (RotationDays * 2) = expiration

# 2. Проверить что используется правильный auth service
# С NewServiceWithKeyRotation или обычный NewService?

# 3. Посмотреть логи ошибок
docker logs api-gateway | grep "Failed to validate"
```

## 📚 Документация

- [JWT_KEY_ROTATION.md](./JWT_KEY_ROTATION.md) — Полная документация
- [key_provider.go](./internal/infrastructure/auth/key_provider.go) — Реализация
- [service.go](./internal/infrastructure/auth/service.go) — Auth Service
- [key_rotation_test.go](./internal/infrastructure/auth/key_rotation_test.go) — Tests

## 🎯 Next Steps

1. ✅ Базовая ротация (DONE)
2. ⏳ Добавить мониторинг (metrics)
3. ⏳ Добавить alert на failed rotation
4. ⏳ Vault интеграция для хранения ключей
5. ⏳ Multi-server синхронизация
