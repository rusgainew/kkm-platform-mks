package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"sync"
	"time"

	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/observability"
	"go.uber.org/zap"
)

// KeyRotationConfig конфигурация для ротации ключей
type KeyRotationConfig struct {
	Enabled      bool      // Включена ли ротация ключей
	RotationDays int       // Ротация каждые N дней
	MaxOldKeys   int       // Максимум старых ключей хранить
	KeyLength    int       // Длина генерируемого ключа в байтах
	StartTime    time.Time // Время начала ротации
	NextRotation time.Time // Время следующей ротации
}

// KeyVersion версия ключа с метаданными
type KeyVersion struct {
	ID        string    // Уникальный ID ключа
	Secret    string    // Сам ключ
	CreatedAt time.Time // Когда создан ключ
	ExpiresAt time.Time // Когда истекает ключ
	Active    bool      // Активный ключ (для новых токенов)
}

// KeyProvider управляет ротацией JWT ключей
type KeyProvider struct {
	mu          sync.RWMutex
	config      *KeyRotationConfig
	logger      *zap.Logger
	metrics     *observability.Metrics
	keyVersions map[string]*KeyVersion // key ID -> KeyVersion
	activeKeyID string                 // ID активного ключа
	keys        []string               // История ключей в порядке создания
	stopChan    chan struct{}
	doneChan    chan struct{}
}

// NewKeyProvider создает новый KeyProvider
func NewKeyProvider(config *KeyRotationConfig, logger *zap.Logger, metrics *observability.Metrics) (*KeyProvider, error) {
	if config == nil {
		config = &KeyRotationConfig{
			Enabled:      true,
			RotationDays: 30,
			MaxOldKeys:   5,
			KeyLength:    32,
			StartTime:    time.Now(),
		}
	}

	kp := &KeyProvider{
		config:      config,
		logger:      logger,
		metrics:     metrics,
		keyVersions: make(map[string]*KeyVersion),
		stopChan:    make(chan struct{}),
		doneChan:    make(chan struct{}),
	}

	// Генерируем первый ключ
	if err := kp.rotateKey(); err != nil {
		return nil, err
	}

	// Стартуем горутину для автоматической ротации
	if config.Enabled {
		go kp.rotationLoop()
	}

	return kp, nil
}

// rotateKey генерирует новый ключ и делает его активным
func (kp *KeyProvider) rotateKey() error {
	kp.mu.Lock()
	defer kp.mu.Unlock()

	// Генерируем новый ключ
	newSecret, err := generateRandomSecret(kp.config.KeyLength)
	if err != nil {
		kp.metrics.RecordKeyRotation("failed")
		return fmt.Errorf("failed to generate new key: %w", err)
	}

	keyID := generateKeyID()
	now := time.Now()
	expiresAt := now.AddDate(0, 0, kp.config.RotationDays*2) // Ключ валиден 2 периода

	newKeyVersion := &KeyVersion{
		ID:        keyID,
		Secret:    newSecret,
		CreatedAt: now,
		ExpiresAt: expiresAt,
		Active:    true,
	}

	// Отключаем старый активный ключ (остаем его валидировать старые токены)
	if kp.activeKeyID != "" {
		if oldKey, exists := kp.keyVersions[kp.activeKeyID]; exists {
			oldKey.Active = false
		}
	}

	// Добавляем новый ключ
	kp.keyVersions[keyID] = newKeyVersion
	kp.keys = append(kp.keys, keyID)
	kp.activeKeyID = keyID

	// Удаляем самые старые ключи если превышен лимит
	if len(kp.keys) > kp.config.MaxOldKeys {
		keysToDelete := len(kp.keys) - kp.config.MaxOldKeys
		for i := 0; i < keysToDelete; i++ {
			oldKeyID := kp.keys[i]
			delete(kp.keyVersions, oldKeyID)
		}
		kp.keys = kp.keys[keysToDelete:]
	}

	// Обновляем время следующей ротации
	kp.config.NextRotation = now.AddDate(0, 0, kp.config.RotationDays)

	// Записываем метрики
	kp.metrics.RecordKeyRotation("success")
	kp.metrics.SetActiveKeysCount(int64(len(kp.keyVersions)))

	kp.logger.Info("Key rotated successfully",
		zap.String("new_key_id", keyID),
		zap.String("active_key_id", kp.activeKeyID),
		zap.Int("total_keys", len(kp.keyVersions)),
		zap.Time("next_rotation", kp.config.NextRotation),
	)

	return nil
}

// GetActiveKey возвращает активный ключ для подписи новых токенов
func (kp *KeyProvider) GetActiveKey() (*KeyVersion, error) {
	start := time.Now()
	kp.mu.RLock()
	defer kp.mu.RUnlock()

	if kp.activeKeyID == "" {
		latency := time.Since(start).Seconds()
		kp.metrics.RecordKeyValidationLatency("error", latency)
		return nil, fmt.Errorf("no active key available")
	}

	keyVersion, exists := kp.keyVersions[kp.activeKeyID]
	if !exists {
		latency := time.Since(start).Seconds()
		kp.metrics.RecordKeyValidationLatency("error", latency)
		return nil, fmt.Errorf("active key not found: %s", kp.activeKeyID)
	}

	latency := time.Since(start).Seconds()
	kp.metrics.RecordKeyValidationLatency("success", latency)

	return keyVersion, nil
}

// GetKeyByID возвращает ключ по ID для проверки существующих токенов
func (kp *KeyProvider) GetKeyByID(keyID string) (*KeyVersion, error) {
	kp.mu.RLock()
	defer kp.mu.RUnlock()

	keyVersion, exists := kp.keyVersions[keyID]
	if !exists {
		return nil, fmt.Errorf("key not found: %s", keyID)
	}

	return keyVersion, nil
}

// GetValidKeys возвращает все валидные ключи для проверки токенов
func (kp *KeyProvider) GetValidKeys() []*KeyVersion {
	start := time.Now()
	kp.mu.RLock()
	defer kp.mu.RUnlock()

	now := time.Now()
	var validKeys []*KeyVersion

	for _, keyVersion := range kp.keyVersions {
		// Ключ валиден если еще не истек срок
		if keyVersion.ExpiresAt.After(now) {
			validKeys = append(validKeys, keyVersion)
		}
	}

	latency := time.Since(start).Seconds()
	kp.metrics.RecordKeyValidationLatency("get_valid_keys", latency)

	return validKeys
}

// GetAllKeys возвращает все ключи (включая истёкшие)
func (kp *KeyProvider) GetAllKeys() []*KeyVersion {
	kp.mu.RLock()
	defer kp.mu.RUnlock()

	var allKeys []*KeyVersion
	for _, keyVersion := range kp.keyVersions {
		allKeys = append(allKeys, keyVersion)
	}

	return allKeys
}

// rotationLoop запускает периодическую ротацию ключей
func (kp *KeyProvider) rotationLoop() {
	defer close(kp.doneChan)

	ticker := time.NewTicker(time.Duration(kp.config.RotationDays) * 24 * time.Hour)
	defer ticker.Stop()

	kp.logger.Info("Key rotation loop started",
		zap.Int("rotation_days", kp.config.RotationDays),
	)

	for {
		select {
		case <-kp.stopChan:
			kp.logger.Info("Key rotation loop stopped")
			return
		case <-ticker.C:
			if err := kp.rotateKey(); err != nil {
				kp.logger.Error("Failed to rotate key", zap.Error(err))
			}
		}
	}
}

// Stop останавливает ротацию ключей
func (kp *KeyProvider) Stop() {
	close(kp.stopChan)
	<-kp.doneChan
}

// generateRandomSecret генерирует случайный секрет
func generateRandomSecret(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

// generateKeyID генерирует уникальный ID ключа
func generateKeyID() string {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback в случае ошибки
		return fmt.Sprintf("key_%d", time.Now().UnixNano())
	}
	return "key_" + base64.URLEncoding.EncodeToString(bytes)
}

// IsKeyExpired проверяет истек ли ключ
func (kv *KeyVersion) IsKeyExpired() bool {
	return time.Now().After(kv.ExpiresAt)
}

// GetKeyStatus возвращает статус ключа
func (kv *KeyVersion) GetKeyStatus() string {
	if time.Now().After(kv.ExpiresAt) {
		return "expired"
	}
	if kv.Active {
		return "active"
	}
	return "valid"
}
