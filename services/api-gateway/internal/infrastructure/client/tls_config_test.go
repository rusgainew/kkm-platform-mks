// Файл api-gateway/internal/infrastructure/client/tls_config_test.go содержит реализацию пакета client.
package client

import (
	"os"
	"testing"
)

func TestGRPCTLSConfig_Validate_TLSDisabled(t *testing.T) {
	config := &GRPCTLSConfig{
		Enabled: false,
	}

	err := config.Validate()
	if err != nil {
		t.Errorf("Validate() failed for disabled TLS: %v", err)
	}
}

func TestGRPCTLSConfig_Validate_TLSEnabledWithCA(t *testing.T) {
	// Создаем временный файл для CA
	tmpFile, err := os.CreateTemp("", "ca-*.crt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	config := &GRPCTLSConfig{
		Enabled:            true,
		CAFile:             tmpFile.Name(),
		InsecureSkipVerify: false,
	}

	err = config.Validate()
	if err != nil {
		t.Errorf("Validate() failed for TLS with CA: %v", err)
	}
}

func TestGRPCTLSConfig_Validate_mTLS(t *testing.T) {
	// Создаем временные файлы
	caFile, _ := os.CreateTemp("", "ca-*.crt")
	certFile, _ := os.CreateTemp("", "cert-*.crt")
	keyFile, _ := os.CreateTemp("", "key-*.key")

	defer os.Remove(caFile.Name())
	defer os.Remove(certFile.Name())
	defer os.Remove(keyFile.Name())

	caFile.Close()
	certFile.Close()
	keyFile.Close()

	config := &GRPCTLSConfig{
		Enabled:            true,
		CAFile:             caFile.Name(),
		CertFile:           certFile.Name(),
		KeyFile:            keyFile.Name(),
		InsecureSkipVerify: false,
	}

	err := config.Validate()
	if err != nil {
		t.Errorf("Validate() failed for mTLS: %v", err)
	}
}

func TestGRPCTLSConfig_Validate_NoCAFile(t *testing.T) {
	config := &GRPCTLSConfig{
		Enabled: true,
		CAFile:  "",
	}

	err := config.Validate()
	if err == nil {
		t.Error("Expected error for missing CA file, got nil")
	}

	expectedErr := "TLS_CA_FILE is required when TLS is enabled"
	if err.Error() != expectedErr {
		t.Errorf("Expected error '%s', got '%s'", expectedErr, err.Error())
	}
}

func TestGRPCTLSConfig_Validate_CertWithoutKey(t *testing.T) {
	config := &GRPCTLSConfig{
		Enabled:  true,
		CAFile:   "/fake/ca.crt",
		CertFile: "/fake/client.crt",
		KeyFile:  "",
	}

	err := config.Validate()
	if err == nil {
		t.Error("Expected error for cert without key, got nil")
	}

	// Может быть либо ошибка о недостающем файле, либо о необходимости обоих файлов
	if err.Error() != "both TLS_CERT_FILE and TLS_KEY_FILE must be set for mTLS" &&
		!contains(err.Error(), "not found") {
		t.Errorf("Expected mTLS validation error, got '%s'", err.Error())
	}
}

func TestGRPCTLSConfig_Validate_KeyWithoutCert(t *testing.T) {
	config := &GRPCTLSConfig{
		Enabled:  true,
		CAFile:   "/fake/ca.crt",
		CertFile: "",
		KeyFile:  "/fake/client.key",
	}

	err := config.Validate()
	if err == nil {
		t.Error("Expected error for key without cert, got nil")
	}

	// Может быть либо ошибка о недостающем файле, либо о необходимости обоих файлов
	if err.Error() != "both TLS_CERT_FILE and TLS_KEY_FILE must be set for mTLS" &&
		!contains(err.Error(), "not found") {
		t.Errorf("Expected mTLS validation error, got '%s'", err.Error())
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[0:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			findSubstr(s, substr)))
}

func findSubstr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestGRPCTLSConfig_IsMTLSEnabled_True(t *testing.T) {
	config := &GRPCTLSConfig{
		Enabled:  true,
		CAFile:   "/path/to/ca.crt",
		CertFile: "/path/to/client.crt",
		KeyFile:  "/path/to/client.key",
	}

	if !config.IsMTLSEnabled() {
		t.Error("Expected mTLS to be enabled, got false")
	}
}

func TestGRPCTLSConfig_IsMTLSEnabled_False_NoCerts(t *testing.T) {
	config := &GRPCTLSConfig{
		Enabled: true,
		CAFile:  "/path/to/ca.crt",
	}

	if config.IsMTLSEnabled() {
		t.Error("Expected mTLS to be disabled, got true")
	}
}

func TestGRPCTLSConfig_IsMTLSEnabled_False_TLSDisabled(t *testing.T) {
	config := &GRPCTLSConfig{
		Enabled: false,
	}

	if config.IsMTLSEnabled() {
		t.Error("Expected mTLS to be disabled when TLS disabled, got true")
	}
}

func TestGRPCTLSConfig_IsMTLSEnabled_False_OnlyCert(t *testing.T) {
	config := &GRPCTLSConfig{
		Enabled:  true,
		CAFile:   "/path/to/ca.crt",
		CertFile: "/path/to/client.crt",
	}

	if config.IsMTLSEnabled() {
		t.Error("Expected mTLS to be disabled with only cert, got true")
	}
}

func TestGRPCTLSConfig_IsMTLSEnabled_False_OnlyKey(t *testing.T) {
	config := &GRPCTLSConfig{
		Enabled: true,
		CAFile:  "/path/to/ca.crt",
		KeyFile: "/path/to/client.key",
	}

	if config.IsMTLSEnabled() {
		t.Error("Expected mTLS to be disabled with only key, got true")
	}
}

func TestGRPCTLSConfig_GetTLSEnabled_True(t *testing.T) {
	config := &GRPCTLSConfig{
		Enabled: true,
	}

	if !config.GetTLSEnabled() {
		t.Error("Expected TLS enabled, got false")
	}
}

func TestGRPCTLSConfig_GetTLSEnabled_False(t *testing.T) {
	config := &GRPCTLSConfig{
		Enabled: false,
	}

	if config.GetTLSEnabled() {
		t.Error("Expected TLS disabled, got true")
	}
}
