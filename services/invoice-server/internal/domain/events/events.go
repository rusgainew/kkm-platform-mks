// Файл invoice-server/internal/domain/events/events.go содержит реализацию пакета events.
package events

import (
	"encoding/json"
	"time"
)

// BaseEvent содержит общие поля для всех событий
type BaseEvent struct {
	EventID   string    `json:"event_id"`
	EventType string    `json:"event_type"`
	Timestamp time.Time `json:"timestamp"`
}

// ToJSON преобразует базовое событие в JSON
func (e *BaseEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// InvoiceCreated событие создания счета-фактуры
type InvoiceCreated struct {
	BaseEvent
	InvoiceID     string  `json:"invoice_id"`
	DocumentUUID  string  `json:"document_uuid"`
	InvoiceNumber string  `json:"invoice_number"`
	TotalAmount   float64 `json:"total_amount"`
	CreatedBy     string  `json:"created_by"`
}

// InvoiceUpdated событие обновления счета-фактуры
type InvoiceUpdated struct {
	BaseEvent
	InvoiceID     string  `json:"invoice_id"`
	DocumentUUID  string  `json:"document_uuid"`
	InvoiceNumber string  `json:"invoice_number"`
	TotalAmount   float64 `json:"total_amount"`
	CreatedBy     string  `json:"created_by"`
}

// InvoiceSigned событие подписания счета-фактуры
type InvoiceSigned struct {
	BaseEvent
	InvoiceID     string `json:"invoice_id"`
	DocumentUUID  string `json:"document_uuid"`
	InvoiceNumber string `json:"invoice_number"`
	SignedBy      string `json:"signed_by"`
}

// InvoiceAccepted событие принятия счета-фактуры
type InvoiceAccepted struct {
	BaseEvent
	InvoiceID     string `json:"invoice_id"`
	DocumentUUID  string `json:"document_uuid"`
	InvoiceNumber string `json:"invoice_number"`
}

// InvoiceRejected событие отклонения счета-фактуры
type InvoiceRejected struct {
	BaseEvent
	InvoiceID     string `json:"invoice_id"`
	DocumentUUID  string `json:"document_uuid"`
	InvoiceNumber string `json:"invoice_number"`
	Reason        string `json:"reason"`
}

// InvoiceRevoked событие отзыва счета-фактуры
type InvoiceRevoked struct {
	BaseEvent
	InvoiceID     string `json:"invoice_id"`
	DocumentUUID  string `json:"document_uuid"`
	InvoiceNumber string `json:"invoice_number"`
	Reason        string `json:"reason"`
}

// Event types
const (
	InvoiceCreatedEvent  = "invoice.created"
	InvoiceUpdatedEvent  = "invoice.updated"
	InvoiceSignedEvent   = "invoice.signed"
	InvoiceAcceptedEvent = "invoice.accepted"
	InvoiceRejectedEvent = "invoice.rejected"
	InvoiceRevokedEvent  = "invoice.revoked"
)
