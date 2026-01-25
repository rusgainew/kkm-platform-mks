package auth

import (
	"testing"
	"time"

	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/observability"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func init() {
	observability.ResetForTesting()
}

func TestKeyProvider_GenerateNewKey(t *testing.T) {
	logger := zap.NewNop()
	metrics := observability.NewMetrics()
	config := &KeyRotationConfig{
		Enabled:      false, // Отключаем авторотацию для теста
		RotationDays: 30,
		MaxOldKeys:   5,
		KeyLength:    32,
	}

	kp, err := NewKeyProvider(config, logger, metrics)
	assert.NoError(t, err)
	assert.NotNil(t, kp)

	// Проверяем что активный ключ создан
	activeKey, err := kp.GetActiveKey()
	assert.NoError(t, err)
	assert.NotNil(t, activeKey)
	assert.NotEmpty(t, activeKey.Secret)
	assert.NotEmpty(t, activeKey.ID)
	assert.True(t, activeKey.Active)
}

func TestKeyProvider_Rotate(t *testing.T) {
	logger := zap.NewNop()
	metrics := observability.NewMetrics()
	config := &KeyRotationConfig{
		Enabled:      false,
		RotationDays: 30,
		MaxOldKeys:   5,
		KeyLength:    32,
	}

	kp, err := NewKeyProvider(config, logger, metrics)
	assert.NoError(t, err)

	// Получаем первый ключ
	key1, err := kp.GetActiveKey()
	assert.NoError(t, err)
	firstSecret := key1.Secret

	// Ротируем ключ
	err = kp.rotateKey()
	assert.NoError(t, err)

	// Получаем новый ключ
	key2, err := kp.GetActiveKey()
	assert.NoError(t, err)

	// Проверяем что ключи разные
	assert.NotEqual(t, firstSecret, key2.Secret)
	assert.NotEqual(t, key1.ID, key2.ID)

	// Проверяем что старый ключ больше не активный
	assert.False(t, key1.Active)

	// Но старый ключ все еще доступен (валиден для старых токенов)
	oldKey, err := kp.GetKeyByID(key1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, oldKey)
}

func TestKeyProvider_GetValidKeys(t *testing.T) {
	logger := zap.NewNop()
	metrics := observability.NewMetrics()
	config := &KeyRotationConfig{
		Enabled:      false,
		RotationDays: 30,
		MaxOldKeys:   5,
		KeyLength:    32,
	}

	kp, err := NewKeyProvider(config, logger, metrics)
	assert.NoError(t, err)

	// Ротируем несколько раз
	for i := 0; i < 3; i++ {
		err = kp.rotateKey()
		assert.NoError(t, err)
	}

	// Получаем валидные ключи
	validKeys := kp.GetValidKeys()
	assert.NotEmpty(t, validKeys)

	// Все валидные ключи должны быть в будущем
	now := time.Now()
	for _, key := range validKeys {
		assert.True(t, key.ExpiresAt.After(now))
	}
}

func TestKeyProvider_MaxOldKeys(t *testing.T) {
	logger := zap.NewNop()
	metrics := observability.NewMetrics()
	config := &KeyRotationConfig{
		Enabled:      false,
		RotationDays: 30,
		MaxOldKeys:   3,
		KeyLength:    32,
	}

	kp, err := NewKeyProvider(config, logger, metrics)
	assert.NoError(t, err)

	// Ротируем больше чем MaxOldKeys
	for i := 0; i < 5; i++ {
		err = kp.rotateKey()
		assert.NoError(t, err)
	}

	// Проверяем что ключей не больше MaxOldKeys + 1 (активный)
	allKeys := kp.GetAllKeys()
	assert.LessOrEqual(t, len(allKeys), config.MaxOldKeys+1)
}

func TestKeyProvider_GetKeyByID(t *testing.T) {
	logger := zap.NewNop()
	metrics := observability.NewMetrics()
	config := &KeyRotationConfig{
		Enabled:      false,
		RotationDays: 30,
		MaxOldKeys:   5,
		KeyLength:    32,
	}

	kp, err := NewKeyProvider(config, logger, metrics)
	assert.NoError(t, err)

	activeKey, _ := kp.GetActiveKey()
	keyID := activeKey.ID

	// Должны найти ключ по ID
	foundKey, err := kp.GetKeyByID(keyID)
	assert.NoError(t, err)
	assert.Equal(t, keyID, foundKey.ID)
	assert.Equal(t, activeKey.Secret, foundKey.Secret)

	// Несуществующий ключ должен вернуть ошибку
	_, err = kp.GetKeyByID("nonexistent")
	assert.Error(t, err)
}

func TestKeyVersion_IsKeyExpired(t *testing.T) {
	now := time.Now()

	// Истекший ключ
	expiredKey := &KeyVersion{
		ID:        "key1",
		ExpiresAt: now.Add(-1 * time.Hour),
	}
	assert.True(t, expiredKey.IsKeyExpired())

	// Валидный ключ
	validKey := &KeyVersion{
		ID:        "key2",
		ExpiresAt: now.Add(1 * time.Hour),
	}
	assert.False(t, validKey.IsKeyExpired())
}

func TestKeyVersion_GetKeyStatus(t *testing.T) {
	now := time.Now()

	// Активный ключ
	activeKey := &KeyVersion{
		ID:        "key1",
		Active:    true,
		ExpiresAt: now.Add(1 * time.Hour),
	}
	assert.Equal(t, "active", activeKey.GetKeyStatus())

	// Неактивный но валидный ключ
	validKey := &KeyVersion{
		ID:        "key2",
		Active:    false,
		ExpiresAt: now.Add(1 * time.Hour),
	}
	assert.Equal(t, "valid", validKey.GetKeyStatus())

	// Истекший ключ
	expiredKey := &KeyVersion{
		ID:        "key3",
		Active:    false,
		ExpiresAt: now.Add(-1 * time.Hour),
	}
	assert.Equal(t, "expired", expiredKey.GetKeyStatus())
}
