// Файл bank-account-server/internal/domain/entities.go содержит реализацию пакета domain.
package domain

import (
	"time"

	"github.com/google/uuid"
)

// BankAccount представляет банковский счет организации
type BankAccount struct {
	ID             uuid.UUID  `db:"id"`
	OrganizationID uuid.UUID  `db:"organization_id"`
	BankID         uuid.UUID  `db:"bank_id"`        // ID банка из справочника
	AccountNumber  string     `db:"account_number"` // Номер счета
	IBAN           string     `db:"iban"`           // Международный номер счета
	Currency       string     `db:"currency"`       // Валюта счета (KZT, USD, EUR и т.д.)
	IsActive       bool       `db:"is_active"`      // Активен ли счет
	IsDefault      bool       `db:"is_default"`     // Счет по умолчанию для организации
	BIC            string     `db:"bic"`            // БИК банка
	BankName       string     `db:"bank_name"`      // Название банка (денормализация для быстрого доступа)
	CreatedBy      uuid.UUID  `db:"created_by"`
	UpdatedBy      *uuid.UUID `db:"updated_by"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at"`
}

// NewBankAccount создает новый банковский счет
func NewBankAccount(organizationID, bankID, createdBy uuid.UUID, accountNumber, iban, currency, bic, bankName string) *BankAccount {
	now := time.Now()
	return &BankAccount{
		ID:             uuid.New(),
		OrganizationID: organizationID,
		BankID:         bankID,
		AccountNumber:  accountNumber,
		IBAN:           iban,
		Currency:       currency,
		BIC:            bic,
		BankName:       bankName,
		IsActive:       true,
		IsDefault:      false,
		CreatedBy:      createdBy,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// Update обновляет данные счета
func (ba *BankAccount) Update(accountNumber, iban, currency, bic, bankName string, updatedBy uuid.UUID) {
	ba.AccountNumber = accountNumber
	ba.IBAN = iban
	ba.Currency = currency
	ba.BIC = bic
	ba.BankName = bankName
	ba.UpdatedBy = &updatedBy
	ba.UpdatedAt = time.Now()
}

// Deactivate деактивирует счет
func (ba *BankAccount) Deactivate(updatedBy uuid.UUID) {
	ba.IsActive = false
	ba.UpdatedBy = &updatedBy
	ba.UpdatedAt = time.Now()
}

// Activate активирует счет
func (ba *BankAccount) Activate(updatedBy uuid.UUID) {
	ba.IsActive = true
	ba.UpdatedBy = &updatedBy
	ba.UpdatedAt = time.Now()
}

// SetAsDefault устанавливает счет как счет по умолчанию
func (ba *BankAccount) SetAsDefault(updatedBy uuid.UUID) {
	ba.IsDefault = true
	ba.UpdatedBy = &updatedBy
	ba.UpdatedAt = time.Now()
}

// UnsetAsDefault снимает флаг счета по умолчанию
func (ba *BankAccount) UnsetAsDefault(updatedBy uuid.UUID) {
	ba.IsDefault = false
	ba.UpdatedBy = &updatedBy
	ba.UpdatedAt = time.Now()
}
