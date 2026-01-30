// Файл catalog-server/internal/domain/ports/publisher.go содержит реализацию пакета ports.
package ports

import "context"

// EventPublisher определяет контракт для публикации доменных событий
type EventPublisher interface {
	// Publish публикует событие
	Publish(ctx context.Context, event interface{}) error

	// Close закрывает соединение с брокером сообщений
	Close() error
}
