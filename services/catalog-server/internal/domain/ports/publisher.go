package ports

import "context"

// EventPublisher определяет контракт для публикации доменных событий
type EventPublisher interface {
	// Publish публикует событие
	Publish(ctx context.Context, event interface{}) error

	// Close закрывает соединение с брокером сообщений
	Close() error
}
