// Файл catalog-server/internal/domain/entities.go содержит реализацию пакета domain.
package domain

import (
	"time"

	"github.com/google/uuid"
)

// CatalogItem представляет товар или услугу в каталоге
type CatalogItem struct {
	ID             uuid.UUID  `db:"id"`
	OrganizationID uuid.UUID  `db:"organization_id"`
	Name           string     `db:"name"`
	Code           string     `db:"code"`        // Артикул/код товара
	Description    string     `db:"description"` // Описание
	UnitType       string     `db:"unit_type"`   // Единица измерения (шт, кг, л, м и т.д.)
	Price          float64    `db:"price"`       // Цена за единицу
	VATRate        float64    `db:"vat_rate"`    // Ставка НДС (0, 12, 20%)
	Category       string     `db:"category"`    // Категория товара/услуги
	IsActive       bool       `db:"is_active"`   // Активен ли товар
	TNVED          string     `db:"tnved"`       // Код ТНВЭД (для товаров)
	GKED           string     `db:"gked"`        // Код ГКЭД (для услуг)
	Barcode        string     `db:"barcode"`     // Штрих-код
	CreatedBy      uuid.UUID  `db:"created_by"`
	UpdatedBy      *uuid.UUID `db:"updated_by"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at"`
}

// NewCatalogItem создает новый элемент каталога
func NewCatalogItem(organizationID, createdBy uuid.UUID, name, code, unitType string, price, vatRate float64) *CatalogItem {
	now := time.Now()
	return &CatalogItem{
		ID:             uuid.New(),
		OrganizationID: organizationID,
		Name:           name,
		Code:           code,
		UnitType:       unitType,
		Price:          price,
		VATRate:        vatRate,
		IsActive:       true,
		CreatedBy:      createdBy,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// Update обновляет данные товара
func (c *CatalogItem) Update(name, description, unitType string, price, vatRate float64, updatedBy uuid.UUID) {
	c.Name = name
	c.Description = description
	c.UnitType = unitType
	c.Price = price
	c.VATRate = vatRate
	c.UpdatedBy = &updatedBy
	c.UpdatedAt = time.Now()
}

// Deactivate деактивирует товар
func (c *CatalogItem) Deactivate(updatedBy uuid.UUID) {
	c.IsActive = false
	c.UpdatedBy = &updatedBy
	c.UpdatedAt = time.Now()
}

// Activate активирует товар
func (c *CatalogItem) Activate(updatedBy uuid.UUID) {
	c.IsActive = true
	c.UpdatedBy = &updatedBy
	c.UpdatedAt = time.Now()
}
