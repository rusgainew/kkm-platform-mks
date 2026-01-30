// Файл document-server/internal/infrastructure/messaging/rabbitmq.go содержит реализацию пакета messaging.
package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type RabbitMQPublisher struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
	logger   *zap.Logger
}

// NewRabbitMQPublisher создает новый издатель событий
func NewRabbitMQPublisher(url, exchange, exchangeType string, logger *zap.Logger) (*RabbitMQPublisher, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Объявляем exchange
	err = channel.ExchangeDeclare(
		exchange,     // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	return &RabbitMQPublisher{
		conn:     conn,
		channel:  channel,
		exchange: exchange,
		logger:   logger,
	}, nil
}

// Publish публикует событие в RabbitMQ
func (p *RabbitMQPublisher) Publish(ctx context.Context, event string, data interface{}) error {
	body, err := json.Marshal(data)
	if err != nil {
		p.logger.Error("Failed to marshal event", zap.Error(err), zap.String("event", event))
		return err
	}

	err = p.channel.PublishWithContext(
		ctx,
		p.exchange, // exchange
		event,      // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		p.logger.Error("Failed to publish event", zap.Error(err), zap.String("event", event))
		return err
	}

	p.logger.Debug("Event published", zap.String("event", event))
	return nil
}

// Close закрывает соединение
func (p *RabbitMQPublisher) Close() error {
	if p.channel != nil {
		p.channel.Close()
	}
	if p.conn != nil {
		p.conn.Close()
	}
	return nil
}

// NoOpPublisher издатель, который ничего не делает (для dev/test)
type NoOpPublisher struct {
	logger *zap.Logger
}

// NewNoOpPublisher создает no-op издатель
func NewNoOpPublisher(logger *zap.Logger) *NoOpPublisher {
	return &NoOpPublisher{logger: logger}
}

// Publish ничего не делает
func (p *NoOpPublisher) Publish(ctx context.Context, event string, data interface{}) error {
	p.logger.Debug("Event (no-op)", zap.String("event", event))
	return nil
}

// Close ничего не делает
func (p *NoOpPublisher) Close() error {
	return nil
}
