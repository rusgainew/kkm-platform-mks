package messaging

import (
	"context"

	"go.uber.org/zap"
)

// NoOpPublisher издатель-заглушка для случаев, когда события отключены
type NoOpPublisher struct {
	logger *zap.Logger
}

// NewNoOpPublisher создает новый экземпляр издателя-заглушки
func NewNoOpPublisher(logger *zap.Logger) *NoOpPublisher {
	logger.Info("Event publishing is disabled, using NoOp publisher")
	return &NoOpPublisher{logger: logger}
}

// Publish логирует событие без публикации
func (p *NoOpPublisher) Publish(ctx context.Context, eventType string, payload []byte) error {
	p.logger.Debug("Event would be published (NoOp mode)",
		zap.String("event_type", eventType),
	)
	return nil
}

// Close закрывает издателя (нет операций)
func (p *NoOpPublisher) Close() error {
	return nil
}
