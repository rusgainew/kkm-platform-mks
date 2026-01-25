package client

import (
	"fmt"
	"os"
)

// GRPCTLSConfig содержит конфигурацию TLS для gRPC соединений
type GRPCTLSConfig struct {
	Enabled            bool
	CertFile           string
	KeyFile            string
	CAFile             string
	InsecureSkipVerify bool
}

// GetTLSEnabled возвращает true если TLS включен
func (tc *GRPCTLSConfig) GetTLSEnabled() bool {
	return tc.Enabled
}

// Validate проверяет конфигурацию TLS
func (tc *GRPCTLSConfig) Validate() error {
	if !tc.Enabled {
		return nil
	}

	// Если TLS включен, проверяем файлы
	if tc.CAFile == "" {
		return fmt.Errorf("TLS_CA_FILE is required when TLS is enabled")
	}

	if _, err := os.Stat(tc.CAFile); err != nil {
		return fmt.Errorf("TLS_CA_FILE not found: %w", err)
	}

	// Если используем mTLS, проверяем клиентские файлы
	if tc.CertFile != "" || tc.KeyFile != "" {
		if tc.CertFile == "" || tc.KeyFile == "" {
			return fmt.Errorf("both TLS_CERT_FILE and TLS_KEY_FILE must be set for mTLS")
		}

		if _, err := os.Stat(tc.CertFile); err != nil {
			return fmt.Errorf("TLS_CERT_FILE not found: %w", err)
		}

		if _, err := os.Stat(tc.KeyFile); err != nil {
			return fmt.Errorf("TLS_KEY_FILE not found: %w", err)
		}
	}

	return nil
}

// IsMTLSEnabled возвращает true если mTLS включен
func (tc *GRPCTLSConfig) IsMTLSEnabled() bool {
	return tc.Enabled && tc.CertFile != "" && tc.KeyFile != ""
}
