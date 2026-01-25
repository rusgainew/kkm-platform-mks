package ports

import (
	"context"
	"time"
)

// CacheLayer определяет интерфейс для кэширования документов
type CacheLayer interface {
	// Get получает значение из кэша по ключу
	// Возвращает nil если ключа нет или он истёк
	Get(ctx context.Context, key string) ([]byte, error)

	// Set устанавливает значение в кэше с TTL
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error

	// Delete удаляет значение из кэша
	Delete(ctx context.Context, key string) error

	// DeletePattern удаляет все ключи, совпадающие с паттерном
	// Например: document:* удалит все ключи документов
	DeletePattern(ctx context.Context, pattern string) error

	// Exists проверяет наличие ключа в кэше
	Exists(ctx context.Context, key string) (bool, error)

	// TTL получает время жизни ключа в секундах
	// Возвращает -2 если ключа нет, -1 если ключ без срока действия
	TTL(ctx context.Context, key string) (int64, error)

	// Close закрывает соединение с кэшем
	Close() error
}
