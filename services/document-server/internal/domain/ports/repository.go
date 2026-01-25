package ports

import "context"

// DocumentRepository определяет интерфейс для работы с документами
type DocumentRepository interface {
	// Create создает новый документ
	Create(ctx context.Context, id, organizationID, title, content, createdBy string, createdAt int64) error

	// Get получает документ по ID
	Get(ctx context.Context, documentID string) (map[string]interface{}, error)

	// GetWithVersion получает документ с информацией о версии (для optimistic locking)
	GetWithVersion(ctx context.Context, documentID string) (map[string]interface{}, error)

	// Update обновляет документ
	Update(ctx context.Context, documentID, title, content string, updatedAt int64) error

	// UpdateWithVersion обновляет документ проверяя версию (optimistic locking)
	// Возвращает ErrVersionConflict если версия не совпадает
	UpdateWithVersion(ctx context.Context, documentID, title, content string, expectedVersion int) error

	// UpdateStatus обновляет статус документа
	UpdateStatus(ctx context.Context, documentID, status string, statusChangedAt int64) error

	// List получает список документов по организации с пагинацией
	List(ctx context.Context, organizationID string, status string, page, perPage int) ([]map[string]interface{}, int64, error)

	// Delete удаляет документ
	Delete(ctx context.Context, documentID string) error
}

// EventPublisher определяет интерфейс для публикации событий
type EventPublisher interface {
	// Publish публикует событие
	Publish(ctx context.Context, event string, data interface{}) error

	// Close закрывает соединение
	Close() error
}
