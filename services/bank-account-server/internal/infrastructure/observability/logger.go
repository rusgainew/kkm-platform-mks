// Файл bank-account-server/internal/infrastructure/observability/logger.go содержит реализацию пакета observability.
package observability

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger создает настроенный логгер
func NewLogger(level, environment string) (*zap.Logger, error) {
	var config zap.Config

	if environment == "production" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Установка уровня логирования
	zapLevel, err := zapcore.ParseLevel(level)
	if err != nil {
		zapLevel = zapcore.InfoLevel
	}
	config.Level = zap.NewAtomicLevelAt(zapLevel)

	return config.Build()
}
