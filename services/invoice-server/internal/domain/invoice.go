// Файл invoice-server/internal/domain/invoice.go содержит реализацию пакета domain.
package domain

import (
	"time"
)

// InvoiceStatus представляет статус счета-фактуры
type InvoiceStatus string

const (
	StatusDraft    InvoiceStatus = "draft"    // Черновик
	StatusSent     InvoiceStatus = "sent"     // Отправлен
	StatusAccepted InvoiceStatus = "accepted" // Принят
	StatusRejected InvoiceStatus = "rejected" // Отклонен
	StatusSigned   InvoiceStatus = "signed"   // Подписан
	StatusRevoked  InvoiceStatus = "revoked"  // Отозван
)

// Invoice представляет счет-фактуру
type Invoice struct {
	ID                           string        `json:"id" db:"id"`
	DocumentUUID                 string        `json:"document_uuid" db:"document_uuid"`
	InvoiceNumber                string        `json:"invoice_number" db:"invoice_number"`
	Number                       string        `json:"number" db:"number"`
	CorrectedReceiptUUID         string        `json:"corrected_receipt_uuid,omitempty" db:"corrected_receipt_uuid"`
	InvoiceDate                  *time.Time    `json:"invoice_date,omitempty" db:"invoice_date"`
	CreatedDate                  time.Time     `json:"created_date" db:"created_date"`
	DeliveryDate                 *time.Time    `json:"delivery_date,omitempty" db:"delivery_date"`
	CorrectedReceiptCreationDate *time.Time    `json:"corrected_receipt_creation_date,omitempty" db:"corrected_receipt_creation_date"`
	TotalAmount                  float64       `json:"total_amount" db:"total_amount"`
	IsResident                   bool          `json:"is_resident" db:"is_resident"`
	Note                         string        `json:"note,omitempty" db:"note"`
	Status                       InvoiceStatus `json:"status" db:"status"`

	// Party information (stored as JSON in DB)
	LegalPersonID string `json:"legal_person_id" db:"legal_person_id"`
	ContractorID  string `json:"contractor_id" db:"contractor_id"`

	// Metadata
	CreatedBy string     `json:"created_by" db:"created_by"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	SignedAt  *time.Time `json:"signed_at,omitempty" db:"signed_at"`
	SignedBy  *string    `json:"signed_by,omitempty" db:"signed_by"`
}

// Validate проверяет валидность счета-фактуры
func (i *Invoice) Validate() error {
	if i.InvoiceNumber == "" {
		return ErrInvoiceNumberRequired
	}
	if i.TotalAmount < 0 {
		return ErrInvalidTotalAmount
	}
	if i.LegalPersonID == "" {
		return ErrLegalPersonRequired
	}
	if i.ContractorID == "" {
		return ErrContractorRequired
	}
	return nil
}

// CanSign проверяет, может ли счет быть подписан
func (i *Invoice) CanSign() bool {
	return i.Status == StatusDraft || i.Status == StatusSent
}

// CanAccept проверяет, может ли счет быть принят
func (i *Invoice) CanAccept() bool {
	return i.Status == StatusSigned
}

// CanReject проверяет, может ли счет быть отклонен
func (i *Invoice) CanReject() bool {
	return i.Status == StatusSigned
}

// CanRevoke проверяет, может ли счет быть отозван
func (i *Invoice) CanRevoke() bool {
	return i.Status == StatusSigned || i.Status == StatusAccepted
}

// Sign подписывает счет
func (i *Invoice) Sign(signedBy string) error {
	if !i.CanSign() {
		return ErrCannotSign
	}
	now := time.Now()
	i.Status = StatusSigned
	i.SignedAt = &now
	i.SignedBy = &signedBy
	i.UpdatedAt = now
	return nil
}

// Accept принимает счет
func (i *Invoice) Accept() error {
	if !i.CanAccept() {
		return ErrCannotAccept
	}
	i.Status = StatusAccepted
	i.UpdatedAt = time.Now()
	return nil
}

// Reject отклоняет счет
func (i *Invoice) Reject() error {
	if !i.CanReject() {
		return ErrCannotReject
	}
	i.Status = StatusRejected
	i.UpdatedAt = time.Now()
	return nil
}

// Revoke отзывает счет
func (i *Invoice) Revoke() error {
	if !i.CanRevoke() {
		return ErrCannotRevoke
	}
	i.Status = StatusRevoked
	i.UpdatedAt = time.Now()
	return nil
}
