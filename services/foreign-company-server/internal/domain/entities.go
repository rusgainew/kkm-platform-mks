// Файл foreign-company-server/internal/domain/entities.go содержит реализацию пакета domain.
package domain

import (
	"time"

	"github.com/google/uuid"
)

// ForeignCompany представляет иностранную компанию-контрагента
type ForeignCompany struct {
	ID          int64      `db:"id"`
	PIN         string     `db:"pin"`          // Идентификационный номер (Tax ID, VAT, TIN)
	FullName    string     `db:"full_name"`    // Полное наименование компании
	CountryCode string     `db:"country_code"` // Код страны ISO 3166-1 alpha-2 (US, DE, FR)
	Address     string     `db:"address"`      // Адрес компании (опционально)
	IsActive    bool       `db:"is_active"`    // Активна ли компания
	CreatedBy   uuid.UUID  `db:"created_by"`
	UpdatedBy   *uuid.UUID `db:"updated_by"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
}

// NewForeignCompany создает новую иностранную компанию
func NewForeignCompany(pin, fullName, countryCode string, createdBy uuid.UUID) *ForeignCompany {
	now := time.Now()
	return &ForeignCompany{
		PIN:         pin,
		FullName:    fullName,
		CountryCode: countryCode,
		IsActive:    true,
		CreatedBy:   createdBy,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Update обновляет данные компании
func (fc *ForeignCompany) Update(pin, fullName, countryCode, address string, updatedBy uuid.UUID) {
	fc.PIN = pin
	fc.FullName = fullName
	fc.CountryCode = countryCode
	fc.Address = address
	fc.UpdatedBy = &updatedBy
	fc.UpdatedAt = time.Now()
}

// Deactivate деактивирует компанию
func (fc *ForeignCompany) Deactivate(updatedBy uuid.UUID) {
	fc.IsActive = false
	fc.UpdatedBy = &updatedBy
	fc.UpdatedAt = time.Now()
}

// Activate активирует компанию
func (fc *ForeignCompany) Activate(updatedBy uuid.UUID) {
	fc.IsActive = true
	fc.UpdatedBy = &updatedBy
	fc.UpdatedAt = time.Now()
}

// Validate проверяет корректность данных
func (fc *ForeignCompany) Validate() error {
	if fc.PIN == "" {
		return ErrInvalidPIN
	}
	if fc.FullName == "" {
		return ErrInvalidFullName
	}
	if len(fc.PIN) < 5 || len(fc.PIN) > 50 {
		return ErrInvalidPIN
	}
	if len(fc.FullName) < 2 || len(fc.FullName) > 500 {
		return ErrInvalidFullName
	}
	if fc.CountryCode != "" && len(fc.CountryCode) != 2 {
		return ErrInvalidCountryCode
	}
	return nil
}
