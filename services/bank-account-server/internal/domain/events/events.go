// Файл bank-account-server/internal/domain/events/events.go содержит реализацию пакета events.
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

// BankAccountCreatedEvent событие создания банковского счета
type BankAccountCreatedEvent struct {
	BaseEvent
	AccountID      uuid.UUID `json:"account_id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	AccountNumber  string    `json:"account_number"`
	Currency       string    `json:"currency"`
	BankName       string    `json:"bank_name"`
	CreatedBy      uuid.UUID `json:"created_by"`
}

// NewBankAccountCreatedEvent создает событие создания счета
func NewBankAccountCreatedEvent(accountID, organizationID, createdBy uuid.UUID, accountNumber, currency, bankName string) *BankAccountCreatedEvent {
	return &BankAccountCreatedEvent{
		BaseEvent: BaseEvent{
			EventID:   uuid.New(),
			EventType: "bank_account.created",
			Timestamp: time.Now(),
		},
		AccountID:      accountID,
		OrganizationID: organizationID,
		AccountNumber:  accountNumber,
		Currency:       currency,
		BankName:       bankName,
		CreatedBy:      createdBy,
	}
}

// BankAccountUpdatedEvent событие обновления банковского счета
type BankAccountUpdatedEvent struct {
	BaseEvent
	AccountID      uuid.UUID `json:"account_id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	AccountNumber  string    `json:"account_number"`
	Currency       string    `json:"currency"`
	UpdatedBy      uuid.UUID `json:"updated_by"`
}

// NewBankAccountUpdatedEvent создает событие обновления счета
func NewBankAccountUpdatedEvent(accountID, organizationID, updatedBy uuid.UUID, accountNumber, currency string) *BankAccountUpdatedEvent {
	return &BankAccountUpdatedEvent{
		BaseEvent: BaseEvent{
			EventID:   uuid.New(),
			EventType: "bank_account.updated",
			Timestamp: time.Now(),
		},
		AccountID:      accountID,
		OrganizationID: organizationID,
		AccountNumber:  accountNumber,
		Currency:       currency,
		UpdatedBy:      updatedBy,
	}
}

// BankAccountDeletedEvent событие удаления банковского счета
type BankAccountDeletedEvent struct {
	BaseEvent
	AccountID      uuid.UUID `json:"account_id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	DeletedBy      uuid.UUID `json:"deleted_by"`
}

// NewBankAccountDeletedEvent создает событие удаления счета
func NewBankAccountDeletedEvent(accountID, organizationID, deletedBy uuid.UUID) *BankAccountDeletedEvent {
	return &BankAccountDeletedEvent{
		BaseEvent: BaseEvent{
			EventID:   uuid.New(),
			EventType: "bank_account.deleted",
			Timestamp: time.Now(),
		},
		AccountID:      accountID,
		OrganizationID: organizationID,
		DeletedBy:      deletedBy,
	}
}

// BankAccountSetDefaultEvent событие установки счета по умолчанию
type BankAccountSetDefaultEvent struct {
	BaseEvent
	AccountID      uuid.UUID `json:"account_id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	UpdatedBy      uuid.UUID `json:"updated_by"`
}

// NewBankAccountSetDefaultEvent создает событие установки счета по умолчанию
func NewBankAccountSetDefaultEvent(accountID, organizationID, updatedBy uuid.UUID) *BankAccountSetDefaultEvent {
	return &BankAccountSetDefaultEvent{
		BaseEvent: BaseEvent{
			EventID:   uuid.New(),
			EventType: "bank_account.set_default",
			Timestamp: time.Now(),
		},
		AccountID:      accountID,
		OrganizationID: organizationID,
		UpdatedBy:      updatedBy,
	}
}

// ToJSON сериализует событие в JSON
func (e *BankAccountCreatedEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

func (e *BankAccountUpdatedEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

func (e *BankAccountDeletedEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

func (e *BankAccountSetDefaultEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}
