// Файл foreign-company-server/internal/domain/events/events.go содержит реализацию пакета events.
package events

import (
	"time"

	"github.com/google/uuid"
)

// ForeignCompanyCreatedEvent событие создания иностранной компании
type ForeignCompanyCreatedEvent struct {
	ID          int64     `json:"id"`
	PIN         string    `json:"pin"`
	FullName    string    `json:"full_name"`
	CountryCode string    `json:"country_code"`
	CreatedBy   uuid.UUID `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}

// ForeignCompanyUpdatedEvent событие обновления иностранной компании
type ForeignCompanyUpdatedEvent struct {
	ID          int64     `json:"id"`
	PIN         string    `json:"pin"`
	FullName    string    `json:"full_name"`
	CountryCode string    `json:"country_code"`
	Address     string    `json:"address"`
	UpdatedBy   uuid.UUID `json:"updated_by"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ForeignCompanyDeletedEvent событие удаления иностранной компании
type ForeignCompanyDeletedEvent struct {
	ID        int64     `json:"id"`
	DeletedBy uuid.UUID `json:"deleted_by"`
	DeletedAt time.Time `json:"deleted_at"`
}

// ForeignCompanyDeactivatedEvent событие деактивации иностранной компании
type ForeignCompanyDeactivatedEvent struct {
	ID            int64     `json:"id"`
	DeactivatedBy uuid.UUID `json:"deactivated_by"`
	DeactivatedAt time.Time `json:"deactivated_at"`
}
