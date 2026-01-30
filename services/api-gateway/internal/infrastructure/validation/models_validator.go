// Файл api-gateway/internal/infrastructure/validation/models_validator.go содержит реализацию пакета validation.
package validation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/domain/models"
)

var (
	uuidPattern  = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	emailPattern = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

type Validator struct {
	v *validator.Validate
}

func NewValidator() *Validator {
	v := validator.New()
	return &Validator{v: v}
}

func (v *Validator) ValidateUUID(id string, fieldName string) error {
	if id == "" {
		return fmt.Errorf("%s is required", fieldName)
	}
	if !uuidPattern.MatchString(strings.ToLower(id)) {
		return fmt.Errorf("%s must be a valid UUID v4", fieldName)
	}
	return nil
}

func (v *Validator) ValidateString(value string, fieldName string, minLen, maxLen int) error {
	if value == "" && minLen > 0 {
		return fmt.Errorf("%s is required", fieldName)
	}
	if len(value) < minLen {
		return fmt.Errorf("%s must be at least %d characters", fieldName, minLen)
	}
	if maxLen > 0 && len(value) > maxLen {
		return fmt.Errorf("%s must not exceed %d characters", fieldName, maxLen)
	}
	return nil
}

func (v *Validator) ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email is required")
	}
	if !emailPattern.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

func (v *Validator) ValidatePagination(page, pageSize int) error {
	if page < 1 {
		return fmt.Errorf("page must be greater than 0")
	}
	if pageSize < 1 {
		return fmt.Errorf("page_size must be greater than 0")
	}
	if pageSize > 100 {
		return fmt.Errorf("page_size must not exceed 100 (DOS protection)")
	}
	return nil
}

func (v *Validator) ValidateEnum(value string, fieldName string, allowedValues []string) error {
	if value == "" {
		return nil
	}
	for _, allowed := range allowedValues {
		if value == allowed {
			return nil
		}
	}
	return fmt.Errorf("%s must be one of: %s", fieldName, strings.Join(allowedValues, ", "))
}

func (v *Validator) ValidateRequired(value string, fieldName string) error {
	if value == "" {
		return fmt.Errorf("%s is required", fieldName)
	}
	return nil
}

func (v *Validator) ValidateFloat(value float64, fieldName string, min, max float64) error {
	if value < min {
		return fmt.Errorf("%s must be at least %.2f", fieldName, min)
	}
	if max > 0 && value > max {
		return fmt.Errorf("%s must not exceed %.2f", fieldName, max)
	}
	return nil
}

func (v *Validator) ValidatePositive(value float64, fieldName string) error {
	if value <= 0 {
		return fmt.Errorf("%s must be positive", fieldName)
	}
	return nil
}

// ValidateCompany validates company data
func (v *Validator) ValidateCompany(company *models.Company) error {
	if err := v.ValidateString(company.Name, "name", 1, 255); err != nil {
		return err
	}
	if err := v.ValidateString(company.Description, "description", 0, 1000); err != nil {
		return err
	}
	if company.OwnerID != "" {
		if err := v.ValidateUUID(company.OwnerID, "owner_id"); err != nil {
			return err
		}
	}
	return nil
}

// ValidateInvoice validates invoice data
func (v *Validator) ValidateInvoice(invoice *models.Invoice) error {
	if err := v.ValidateString(invoice.InvoiceNumber, "invoice_number", 1, 50); err != nil {
		return err
	}
	if err := v.ValidatePositive(invoice.TotalAmount, "total_amount"); err != nil {
		return err
	}
	if err := v.ValidateString(invoice.InvoiceDate, "invoice_date", 1, 50); err != nil {
		return err
	}
	return nil
}

// ValidateDocument validates document data
func (v *Validator) ValidateDocument(doc *models.Document) error {
	if err := v.ValidateString(doc.OrganizationID, "organization_id", 1, 64); err != nil {
		return err
	}
	if err := v.ValidateUUID(doc.OrganizationID, "organization_id"); err != nil {
		return err
	}
	if err := v.ValidateString(doc.Title, "title", 1, 255); err != nil {
		return err
	}
	if err := v.ValidateString(doc.Content, "content", 1, 5000); err != nil {
		return err
	}
	return nil
}

// ValidateCatalog validates catalog data
func (v *Validator) ValidateCatalog(catalog *models.Catalog) error {
	if err := v.ValidateString(catalog.Name, "name", 1, 255); err != nil {
		return err
	}
	if err := v.ValidateString(catalog.Description, "description", 0, 1000); err != nil {
		return err
	}
	if catalog.Price > 0 {
		if err := v.ValidatePositive(catalog.Price, "price"); err != nil {
			return err
		}
	}
	if catalog.Currency != "" {
		if err := v.ValidateEnum(catalog.Currency, "currency", []string{"USD", "EUR", "RUB", "KZT"}); err != nil {
			return err
		}
	}
	return nil
}

// ValidateBankAccount validates bank account data
func (v *Validator) ValidateBankAccount(account *models.BankAccount) error {
	if err := v.ValidateString(account.AccountNumber, "account_number", 1, 50); err != nil {
		return err
	}
	if err := v.ValidateString(account.BankName, "bank_name", 1, 255); err != nil {
		return err
	}
	if err := v.ValidateString(account.BankCode, "bank_code", 1, 20); err != nil {
		return err
	}
	if account.Currency != "" {
		if err := v.ValidateEnum(account.Currency, "currency", []string{"USD", "EUR", "RUB", "KZT"}); err != nil {
			return err
		}
	}
	if account.OwnerID != "" {
		if err := v.ValidateUUID(account.OwnerID, "owner_id"); err != nil {
			return err
		}
	}
	return nil
}
