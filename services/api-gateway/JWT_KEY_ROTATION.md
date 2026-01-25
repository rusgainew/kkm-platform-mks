# JWT Key Rotation - Ротация ключей

## 📋 Обзор

Реализована **автоматическая ротация JWT ключей** для повышения безопасности. Система регулярно генерирует новые ключи, при этом сохраняя возможность валидировать старые токены с предыдущими ключами.

## 🏗️ Архитектура

### Компоненты

```
┌──────────────────────────────────────────┐
│        Auth Service                      │
│  ┌────────────────────────────────────┐  │
│  │  Key Provider (Rotation Manager)   │  │
│  │  ┌──────────────────────────────┐  │  │
│  │  │ Active Key      (for signing)│  │  │
│  │  │ Valid Keys      (for verify) │  │  │
│  │  │ Expired Keys    (archived)   │  │  │
│  │  └──────────────────────────────┘  │  │
│  └────────────────────────────────────┘  │
└──────────────────────────────────────────┘
         ↓                    ↓
    Generate Token      Validate Token
    (с активным)        (со всеми)
```

### Процесс ротации

```
Day 1-30: Use KEY-1 (Active)
         └─→ Generate tokens with KEY-1
         └─→ Can validate with KEY-1

Day 31-60: Use KEY-2 (Active)
          └─→ Generate tokens with KEY-2
          └─→ Can validate with KEY-1 and KEY-2 (grace period)

Day 61+:  KEY-1 Expired
         └─→ Can still validate old KEY-1 tokens
         └─→ Generate tokens only with KEY-2
```

## 🔧 Конфигурация

### Включить ротацию в .env

```bash
# JWT
JWT_SECRET=fallback-secret-key
JWT_EXPIRATION=24h

# Key Rotation
KEY_ROTATION_ENABLED=true
KEY_ROTATION_DAYS=30        # Ротация каждые 30 дней
KEY_ROTATION_MAX_OLD_KEYS=5 # Хранить 5 старых ключей
KEY_ROTATION_KEY_LENGTH=32  # 32 байта ключа
```

### В Docker Compose

```yaml
services:
  api-gateway:
    environment:
      KEY_ROTATION_ENABLED: "true"
      KEY_ROTATION_DAYS: "30"
      KEY_ROTATION_MAX_OLD_KEYS: "5"
```

## 📝 Использование

### Базовый способ (без ротации)

```go
// Старый способ - работает как раньше
authService := auth.NewService(jwtSecret, jwtExpiration, logger)

// Генерируем токен
token, _ := authService.GenerateToken(ctx, userID, claims)

// Валидируем токен
claims, _ := authService.ValidateToken(ctx, token)
```

### С ротацией ключей

```go
// Новый способ - с автоматической ротацией
rotationConfig := &auth.KeyRotationConfig{
    Enabled:      true,
    RotationDays: 30,
    MaxOldKeys:   5,
    KeyLength:    32,
}

authService, _ := auth.NewServiceWithKeyRotation(
    jwtSecret,
    jwtExpiration,
    logger,
    rotationConfig,
)

// Генерируем токен (с активным ключом)
token, _ := authService.GenerateToken(ctx, userID, claims)

// Валидируем токен (работает со всеми валидными ключами)
claims, _ := authService.ValidateToken(ctx, token)

// При shutdown
defer authService.Stop()
```

### В DI контейнере

```go
// internal/infrastructure/config/config.go
type KeyRotationConfig struct {
    Enabled       bool
    RotationDays  int
    MaxOldKeys    int
    KeyLength     int
}

func Load() (*Config, error) {
    // ...
    KeyRotation: KeyRotationConfig{
        Enabled:      getEnvAsBool("KEY_ROTATION_ENABLED", true),
        RotationDays: getEnvAsInt("KEY_ROTATION_DAYS", 30),
        MaxOldKeys:   getEnvAsInt("KEY_ROTATION_MAX_OLD_KEYS", 5),
        KeyLength:    getEnvAsInt("KEY_ROTATION_KEY_LENGTH", 32),
    },
}

// pkg/di_container/initializer.go
func (i *Initializer) InitializeInfrastructure() error {
    cfg := i.container.Config()
    logger := i.container.Logger()

    var authService *auth.Service
    var err error

    if cfg.KeyRotation.Enabled {
        authService, err = auth.NewServiceWithKeyRotation(
            cfg.JWT.Secret,
            cfg.JWT.Expiration,
            logger,
            &auth.KeyRotationConfig{
                Enabled:      cfg.KeyRotation.Enabled,
                RotationDays: cfg.KeyRotation.RotationDays,
                MaxOldKeys:   cfg.KeyRotation.MaxOldKeys,
                KeyLength:    cfg.KeyRotation.KeyLength,
            },
        )
    } else {
        authService = auth.NewService(cfg.JWT.Secret, cfg.JWT.Expiration, logger)
    }

    i.container.SetAuthService(authService)
    i.container.RegisterCleanup(func() error {
        authService.Stop()
        return nil
    })

    return nil
}
```

## 🔑 Стратегия ключей

### KeyVersion структура

```go
type KeyVersion struct {
    ID        string    // Уникальный ID (base64)
    Secret    string    // Секрет в base64
    CreatedAt time.Time // Когда создан
    ExpiresAt time.Time // Когда истекает
    Active    bool      // Активный ли (для подписи новых)
}
```

### Жизненный цикл ключа

```
1. ACTIVE (новый)
   ├─ Используется для подписи новых токенов
   └─ Валиден для проверки токенов

2. VALID (при ротации следующего)
   ├─ Не используется для подписи
   └─ Валиден для проверки старых токенов (grace period)

3. EXPIRED (после 2x RotationDays)
   ├─ Удаляется из памяти
   └─ Старые токены больше не валидны
```

## 🔄 Процесс ротации

### Автоматическая ротация

```go
// Ротация запускается в горутине каждые RotationDays
keyProvider.rotationLoop()

// Каждые 30 дней:
// 1. Генерируется новый ключ
// 2. Старый становится неактивным
// 3. История хранится (max MaxOldKeys)
```

### Ручная ротация

```go
keyProvider, _ := auth.NewKeyProvider(config, logger)

// Ротировать сейчас
keyProvider.rotateKey()

// Получить активный ключ
activeKey, _ := keyProvider.GetActiveKey()

// Получить все валидные ключи
validKeys := keyProvider.GetValidKeys()

// Получить конкретный ключ по ID
keyVersion, _ := keyProvider.GetKeyByID(keyID)
```

## 📊 Пример потока

```
Client Request: POST /api/v1/login
        ↓
Auth Handler
        ↓
authService.GenerateToken()
        ↓
KeyProvider.GetActiveKey()  ← KEY-123 (Active)
        ↓
Sign with KEY-123
        ↓
Return Token with KEY-123 in claims.ID
        ↓
Client stores token

---

Later Request: GET /api/v1/companies/123
Header: Authorization: Bearer <token>
        ↓
authService.ValidateToken()
        ↓
Parse token header (claims.ID = KEY-123)
        ↓
KeyProvider.GetValidKeys()
        ↓
Returns [KEY-123, KEY-124, ...] (если KEY-123 не истекла)
        ↓
Try validate with KEY-123 ✓ Success
        ↓
Return claims
```

## ✨ Особенности

✅ **Автоматическая ротация** — каждые N дней  
✅ **Backcompat** — старые токены остаются валидными  
✅ **Graceful degradation** — без KeyProvider работает как раньше  
✅ **Масштабируемость** — храним только последние N ключей  
✅ **Thread-safe** — используется sync.RWMutex  
✅ **Мониторируемо** — логирование всех операций  
✅ **Тестируемо** — с unit тестами

## 🧪 Тестирование

### Unit тесты

```bash
go test -v ./internal/infrastructure/auth/... -run TestKeyProvider
```

### Примеры тестов

```go
// Тест генерации нового ключа
TestKeyProvider_GenerateNewKey

// Тест ротации
TestKeyProvider_Rotate

// Тест валидных ключей
TestKeyProvider_GetValidKeys

// Тест лимита старых ключей
TestKeyProvider_MaxOldKeys

// Тест получения ключа по ID
TestKeyProvider_GetKeyByID

// Тест статуса ключа
TestKeyVersion_GetKeyStatus
```

## 📈 Мониторинг

### Логи

```
INFO: Key rotated successfully
  new_key_id=key_abc123
  active_key_id=key_abc123
  total_keys=3
  next_rotation=2026-02-10T10:30:00Z
```

### Метрики (рекомендуется)

```
api_gateway_key_rotation_total{status} # Всего ротаций
api_gateway_active_keys{status}        # Активные ключи
api_gateway_expired_keys{status}       # Истекшие ключи
```

## ⚠️ Обработка ошибок

### Ключ не найден

```go
// При валидации старого токена
// Если KEY-123 уже удален (истекла)
ValidateToken() → Error: "no valid key found"
```

### Fallback поведение

```go
// Если KeyProvider недоступен
s.keyProvider == nil
→ Используем jwtSecret как обычно
→ Сервис работает без ротации
```

## 🚀 Production рекомендации

1. **Хранение ключей**

   - ✅ Current: в памяти
   - ⏳ Future: в защищенном хранилище (Vault, Consul)

2. **Распределенная ротация**

   - ✅ Single instance: локальная ротация
   - ⏳ Multiple instances: синхронизированная ротация

3. **Аудит**

   - Логировать все ротации
   - Логировать все валидации с ключом
   - Отслеживать usage по ключам

4. **Frequency**
   - Для критичных систем: 7-14 дней
   - Для обычных: 30 дней
   - Для низкого risk: 90 дней

## 🔗 Связанные документы

- [service.go](./internal/infrastructure/auth/service.go) — Auth Service
- [key_provider.go](./internal/infrastructure/auth/key_provider.go) — Key Rotation
- [key_rotation_test.go](./internal/infrastructure/auth/key_rotation_test.go) — Tests
