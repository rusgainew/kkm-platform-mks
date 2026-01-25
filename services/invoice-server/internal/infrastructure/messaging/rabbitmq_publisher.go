package messaging

import (
	"context"
	"fmt"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

// RabbitMQPublisher публикует события в RabbitMQ
type RabbitMQPublisher struct {
	conn     *amqp091.Connection
	channel  *amqp091.Channel
	exchange string
	logger   *zap.Logger
}

// NewRabbitMQPublisher создает новый экземпляр издателя
func NewRabbitMQPublisher(amqpURL, exchange, exchangeType string, logger *zap.Logger) (*RabbitMQPublisher, error) {
	conn, err := amqp091.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Объявление exchange
	err = ch.ExchangeDeclare(
		exchange,     // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	logger.Info("RabbitMQ publisher initialized",
		zap.String("exchange", exchange),
		zap.String("type", exchangeType),
	)

	return &RabbitMQPublisher{
		conn:     conn,
		channel:  ch,
		exchange: exchange,
		logger:   logger,
	}, nil
}

// Publish публикует событие
func (p *RabbitMQPublisher) Publish(ctx context.Context, eventType string, payload []byte) error {
	routingKey := fmt.Sprintf("invoice.%s", eventType)

	err := p.channel.PublishWithContext(
		ctx,
		p.exchange, // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp091.Publishing{
			ContentType:  "application/json",
			Body:         payload,
			DeliveryMode: amqp091.Persistent,
			Timestamp:    time.Now(),
		},
	)

	if err != nil {
		p.logger.Error("Failed to publish event",
			zap.Error(err),
			zap.String("event_type", eventType),
		)
		return fmt.Errorf("failed to publish event: %w", err)
	}

	p.logger.Info("Event published",
		zap.String("event_type", eventType),
		zap.String("routing_key", routingKey),
	)

	return nil
}

// Close закрывает соединение
func (p *RabbitMQPublisher) Close() error {
	if p.channel != nil {
		p.channel.Close()
	}
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}
