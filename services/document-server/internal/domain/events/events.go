// Файл document-server/internal/domain/events/events.go содержит реализацию пакета events.
package events

// DocumentCreatedEvent события создания документа
type DocumentCreatedEvent struct {
	DocumentID     string
	OrganizationID string
	Title          string
	CreatedBy      string
	CreatedAt      int64
}

// DocumentUpdatedEvent события обновления документа
type DocumentUpdatedEvent struct {
	DocumentID string
	Title      string
	UpdatedAt  int64
}

// DocumentSentEvent события отправки документа на согласование
type DocumentSentEvent struct {
	DocumentID  string
	RecipientID string
	Status      string
	SentAt      int64
}

// DocumentApprovedEvent события одобрения документа
type DocumentApprovedEvent struct {
	DocumentID string
	ApprovedBy string
	ApprovedAt int64
}

// DocumentRejectedEvent события отклонения документа
type DocumentRejectedEvent struct {
	DocumentID string
	RejectedBy string
	Reason     string
	RejectedAt int64
}

// DocumentArchivedEvent события архивирования документа
type DocumentArchivedEvent struct {
	DocumentID string
	ArchivedAt int64
}
