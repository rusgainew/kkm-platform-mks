// Файл company-server/internal/infrastructure/messaging/noop_publisher.go содержит реализацию пакета messaging.
package messaging

import (
	"context"

	"go.uber.org/zap"
)

// NoOpPublisher реализация publisher, которая ничего не делает (для тестов и отключенного режима)
type NoOpPublisher struct {
	logger *zap.Logger
}

// NewNoOpPublisher создает новый NoOp publisher
func NewNoOpPublisher(logger *zap.Logger) *NoOpPublisher {
	logger.Info("NoOp event publisher initialized (events disabled)")
	return &NoOpPublisher{logger: logger}
}

// Publish ничего не делает
func (p *NoOpPublisher) Publish(ctx context.Context, eventType string, payload []byte) error {
	p.logger.Debug("Event would be published (NoOp mode)", zap.String("event_type", eventType))
	return nil
}

// Close ничего не делает
func (p *NoOpPublisher) Close() error {
	return nil
}
