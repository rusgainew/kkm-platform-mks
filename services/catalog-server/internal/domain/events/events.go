package events

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// BaseEvent содержит общие поля для всех событий
type BaseEvent struct {
	EventID   uuid.UUID `json:"event_id"`
	EventType string    `json:"event_type"`
	Timestamp time.Time `json:"timestamp"`
}

// CatalogItemCreatedEvent событие создания элемента каталога
type CatalogItemCreatedEvent struct {
	BaseEvent
	ItemID         uuid.UUID `json:"item_id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Name           string    `json:"name"`
	Code           string    `json:"code"`
	Price          float64   `json:"price"`
	CreatedBy      uuid.UUID `json:"created_by"`
}

// NewCatalogItemCreatedEvent создает событие создания элемента
func NewCatalogItemCreatedEvent(itemID, organizationID, createdBy uuid.UUID, name, code string, price float64) *CatalogItemCreatedEvent {
	return &CatalogItemCreatedEvent{
		BaseEvent: BaseEvent{
			EventID:   uuid.New(),
			EventType: "catalog.item.created",
			Timestamp: time.Now(),
		},
		ItemID:         itemID,
		OrganizationID: organizationID,
		Name:           name,
		Code:           code,
		Price:          price,
		CreatedBy:      createdBy,
	}
}

// CatalogItemUpdatedEvent событие обновления элемента каталога
type CatalogItemUpdatedEvent struct {
	BaseEvent
	ItemID         uuid.UUID `json:"item_id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Name           string    `json:"name"`
	Price          float64   `json:"price"`
	UpdatedBy      uuid.UUID `json:"updated_by"`
}

// NewCatalogItemUpdatedEvent создает событие обновления элемента
func NewCatalogItemUpdatedEvent(itemID, organizationID, updatedBy uuid.UUID, name string, price float64) *CatalogItemUpdatedEvent {
	return &CatalogItemUpdatedEvent{
		BaseEvent: BaseEvent{
			EventID:   uuid.New(),
			EventType: "catalog.item.updated",
			Timestamp: time.Now(),
		},
		ItemID:         itemID,
		OrganizationID: organizationID,
		Name:           name,
		Price:          price,
		UpdatedBy:      updatedBy,
	}
}

// CatalogItemDeletedEvent событие удаления элемента каталога
type CatalogItemDeletedEvent struct {
	BaseEvent
	ItemID         uuid.UUID `json:"item_id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	DeletedBy      uuid.UUID `json:"deleted_by"`
}

// NewCatalogItemDeletedEvent создает событие удаления элемента
func NewCatalogItemDeletedEvent(itemID, organizationID, deletedBy uuid.UUID) *CatalogItemDeletedEvent {
	return &CatalogItemDeletedEvent{
		BaseEvent: BaseEvent{
			EventID:   uuid.New(),
			EventType: "catalog.item.deleted",
			Timestamp: time.Now(),
		},
		ItemID:         itemID,
		OrganizationID: organizationID,
		DeletedBy:      deletedBy,
	}
}

// ToJSON сериализует событие в JSON
func (e *CatalogItemCreatedEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

func (e *CatalogItemUpdatedEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

func (e *CatalogItemDeletedEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}
