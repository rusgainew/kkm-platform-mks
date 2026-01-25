package masking

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rusgainew/kkm-project-mks/api-gateway/internal/infrastructure/observability"
)

func init() {
	observability.ResetForTesting()
}

func TestMaskPassword(t *testing.T) {
	metrics := observability.NewMetrics()
	masker := NewSensitiveFieldMasker(metrics)
	input := `{"username":"john","password":"my-secret-password"}`
	output := masker.MaskSensitiveFields(input)

	assert.NotContains(t, output, "my-secret-password")
	assert.Contains(t, output, "***")
	assert.Contains(t, output, "john")
}

func TestMaskJWT(t *testing.T) {
	metrics := observability.NewMetrics()
	masker := NewSensitiveFieldMasker(metrics)
	input := `Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9`
	output := masker.Mask(input)

	assert.NotContains(t, output, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9")
	assert.Contains(t, output, "***")
}

func TestMaskAPIKey(t *testing.T) {
	metrics := observability.NewMetrics()
	masker := NewSensitiveFieldMasker(metrics)
	input := `{"api_key":"sk-1234567890abcdefghij"}`
	output := masker.MaskSensitiveFields(input)

	assert.NotContains(t, output, "1234567890abcdefghij")
	assert.Contains(t, output, "***")
}

func TestMaskCreditCard(t *testing.T) {
	metrics := observability.NewMetrics()
	masker := NewSensitiveFieldMasker(metrics)
	input := "Credit card: 4532-1234-5678-9010"
	output := masker.Mask(input)

	assert.NotContains(t, output, "4532-1234-5678-9010")
	assert.Contains(t, output, "***")
}

func TestMaskEmail(t *testing.T) {
	metrics := observability.NewMetrics()
	masker := NewSensitiveFieldMasker(metrics)
	input := "User email: admin@example.com and backup@example.com"
	output := masker.Mask(input)

	assert.NotContains(t, output, "admin@example.com")
	assert.NotContains(t, output, "backup@example.com")
	assert.Contains(t, output, "***")
}

func TestMaskSSN(t *testing.T) {
	metrics := observability.NewMetrics()
	masker := NewSensitiveFieldMasker(metrics)
	input := "SSN: 123-45-6789"
	output := masker.Mask(input)

	assert.NotContains(t, output, "123-45-6789")
	assert.Contains(t, output, "***")
}

func TestMaskMultipleFields(t *testing.T) {
	metrics := observability.NewMetrics()
	masker := NewSensitiveFieldMasker(metrics)
	input := `{
		"username":"john",
		"password":"secret123",
		"api_key":"key-abc-123",
		"email":"john@example.com"
	}`
	output := masker.MaskSensitiveFields(input)

	assert.NotContains(t, output, "secret123")
	assert.NotContains(t, output, "key-abc-123")
	assert.NotContains(t, output, "john@example.com")
	assert.Contains(t, output, "john")
}

func TestMaskHttpHeader(t *testing.T) {
	metrics := observability.NewMetrics()
	masker := NewSensitiveFieldMasker(metrics)
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
	masked := masker.MaskInHttpHeader("Authorization", token)

	assert.NotEqual(t, token, masked)
	assert.Contains(t, masked, "***")
}
