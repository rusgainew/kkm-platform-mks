// Файл invoice-server/internal/domain/ports/event_publisher.go содержит реализацию пакета ports.
package ports

import "context"

// EventPublisher определяет интерфейс для публикации событий
type EventPublisher interface {
	// Publish публикует событие
	Publish(ctx context.Context, eventType string, payload []byte) error

	// Close закрывает соединение с брокером сообщений
	Close() error
}
