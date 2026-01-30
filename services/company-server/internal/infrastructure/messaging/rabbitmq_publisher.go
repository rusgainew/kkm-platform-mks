// Файл company-server/internal/infrastructure/messaging/rabbitmq_publisher.go содержит реализацию пакета messaging.
package messaging

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

// RabbitMQPublisher реализует публикацию событий в RabbitMQ
type RabbitMQPublisher struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
	logger   *zap.Logger
}

// NewRabbitMQPublisher создает новый publisher для RabbitMQ
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

	// Объявление exchange
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

	logger.Info("RabbitMQ publisher initialized", zap.String("exchange", exchange))

	return &RabbitMQPublisher{
		conn:     conn,
		channel:  channel,
		exchange: exchange,
		logger:   logger,
	}, nil
}

// Publish публикует событие в RabbitMQ
func (p *RabbitMQPublisher) Publish(ctx context.Context, eventType string, payload []byte) error {
	err := p.channel.PublishWithContext(
		ctx,
		p.exchange, // exchange
		eventType,  // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        payload,
		},
	)

	if err != nil {
		p.logger.Error("Failed to publish event", zap.Error(err), zap.String("event_type", eventType))
		return err
	}

	p.logger.Debug("Event published", zap.String("event_type", eventType))
	return nil
}

// Close закрывает соединение с RabbitMQ
func (p *RabbitMQPublisher) Close() error {
	if p.channel != nil {
		if err := p.channel.Close(); err != nil {
			p.logger.Error("Failed to close channel", zap.Error(err))
		}
	}
	if p.conn != nil {
		if err := p.conn.Close(); err != nil {
			p.logger.Error("Failed to close connection", zap.Error(err))
			return err
		}
	}
	p.logger.Info("RabbitMQ publisher closed")
	return nil
}
