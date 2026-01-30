// Файл company-server/internal/domain/events/events.go содержит реализацию пакета events.
package events

import (
	"encoding/json"
	"time"
)

// Event types
const (
	OrganizationCreatedEvent = "organization.created"
	OrganizationUpdatedEvent = "organization.updated"
	OrganizationDeletedEvent = "organization.deleted"
	EmployeeAddedEvent       = "employee.added"
	EmployeeRemovedEvent     = "employee.removed"
)

// BaseEvent содержит общие поля для всех событий
type BaseEvent struct {
	EventID   string    `json:"event_id"`
	EventType string    `json:"event_type"`
	Timestamp time.Time `json:"timestamp"`
}

// OrganizationCreated событие создания организации
type OrganizationCreated struct {
	BaseEvent
	OrganizationID string `json:"organization_id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	OwnerID        string `json:"owner_id"`
}

// OrganizationUpdated событие обновления организации
type OrganizationUpdated struct {
	BaseEvent
	OrganizationID string `json:"organization_id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
}

// OrganizationDeleted событие удаления организации
type OrganizationDeleted struct {
	BaseEvent
	OrganizationID string `json:"organization_id"`
}

// EmployeeAdded событие добавления участника
type EmployeeAdded struct {
	BaseEvent
	EmployeeID     string `json:"employee_id"`
	OrganizationID string `json:"organization_id"`
	UserID         string `json:"user_id"`
	Role           string `json:"role"`
}

// EmployeeRemoved событие удаления участника
type EmployeeRemoved struct {
	BaseEvent
	EmployeeID     string `json:"employee_id"`
	OrganizationID string `json:"organization_id"`
	UserID         string `json:"user_id"`
}

// ToJSON преобразует событие в JSON
func (e *BaseEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}
