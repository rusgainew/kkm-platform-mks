package domain

// InvoiceDetail представляет позицию (строку) в счете-фактуре
type InvoiceDetail struct {
	ID               int64   `json:"id" db:"id"`
	InvoiceUUID      string  `json:"invoice_uuid" db:"invoice_uuid"`
	BaseCount        float64 `json:"base_count" db:"base_count"`
	Price            float64 `json:"price" db:"price"`
	Amount           float64 `json:"amount" db:"amount"`
	AmountWithoutVAT float64 `json:"amount_without_vat" db:"amount_without_vat"`
	AmountVAT        float64 `json:"amount_vat" db:"amount_vat"`
	AmountST         float64 `json:"amount_st" db:"amount_st"`
	GoodsName        string  `json:"goods_name" db:"goods_name"`
	TNVEDCode        string  `json:"tnved_code,omitempty" db:"tnved_code"`
	GKEDCode         string  `json:"gked_code,omitempty" db:"gked_code"`
	FCDNumber        string  `json:"fcd_number,omitempty" db:"fcd_number"`

	// Catalog references
	CatalogID  string `json:"catalog_id,omitempty" db:"catalog_id"`
	UnitTypeID string `json:"unit_type_id,omitempty" db:"unit_type_id"`
}

// Validate проверяет валидность позиции счета
func (d *InvoiceDetail) Validate() error {
	if d.InvoiceUUID == "" {
		return ErrInvalidInvoiceID
	}
	if d.BaseCount <= 0 {
		return ErrInvalidQuantity
	}
	if d.Price < 0 {
		return ErrInvalidPrice
	}
	if d.Amount < 0 {
		return ErrInvalidAmount
	}
	if d.GoodsName == "" {
		return ErrMissingRequiredField
	}
	return nil
}

// CalculateTotalWithTaxes вычисляет общую сумму с учетом налогов
func (d *InvoiceDetail) CalculateTotalWithTaxes() float64 {
	return d.Amount + d.AmountVAT + d.AmountST
}
