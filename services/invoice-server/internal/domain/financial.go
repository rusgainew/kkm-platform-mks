// Файл invoice-server/internal/domain/financial.go содержит реализацию пакета domain.
package domain

// FinancialData содержит финансовые данные счета
type FinancialData struct {
	InvoiceUUID                 string  `json:"invoice_uuid" db:"invoice_uuid"`
	TotalAmount                 float64 `json:"total_amount" db:"total_amount"`
	OpeningBalances             float64 `json:"opening_balances" db:"opening_balances"`
	AssessedContributionsAmount float64 `json:"assessed_contributions_amount" db:"assessed_contributions_amount"`
	PaidAmount                  float64 `json:"paid_amount" db:"paid_amount"`
	PenaltiesAmount             float64 `json:"penalties_amount" db:"penalties_amount"`
	FinesAmount                 float64 `json:"fines_amount" db:"fines_amount"`
	ClosingBalances             float64 `json:"closing_balances" db:"closing_balances"`
	AmountToBePaid              float64 `json:"amount_to_be_paid" db:"amount_to_be_paid"`
	PersonalAccountNumber       string  `json:"personal_account_number,omitempty" db:"personal_account_number"`
	LegalPersonBankAccount      string  `json:"legal_person_bank_account,omitempty" db:"legal_person_bank_account"`
	ContractorBankAccount       string  `json:"contractor_bank_account,omitempty" db:"contractor_bank_account"`
}

// InvoiceContract содержит контрактные данные счета
type InvoiceContract struct {
	InvoiceUUID    string `json:"invoice_uuid" db:"invoice_uuid"`
	ContractNumber string `json:"contract_number,omitempty" db:"contract_number"`
	ContractDate   string `json:"contract_date,omitempty" db:"contract_date"`
	PaymentTerms   string `json:"payment_terms,omitempty" db:"payment_terms"`
	DeliveryTerms  string `json:"delivery_terms,omitempty" db:"delivery_terms"`
}
