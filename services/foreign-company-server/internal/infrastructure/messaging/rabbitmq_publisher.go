package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rusgainew/kkm-project-mks/foreign-company-server/internal/domain/events"
	"go.uber.org/zap"
)

// RabbitMQPublisher реализует публикацию событий в RabbitMQ
type RabbitMQPublisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	logger  *zap.Logger
}

// NewRabbitMQPublisher создает новый издатель для RabbitMQ
func NewRabbitMQPublisher(url string, logger *zap.Logger) (*RabbitMQPublisher, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Объявление exchange для иностранных компаний
	err = channel.ExchangeDeclare(
		"foreign_company.events", // name
		"topic",                  // type
		true,                     // durable
		false,                    // auto-deleted
		false,                    // internal
		false,                    // no-wait
		nil,                      // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	logger.Info("RabbitMQ publisher connected", zap.String("exchange", "foreign_company.events"))

	return &RabbitMQPublisher{
		conn:    conn,
		channel: channel,
		logger:  logger,
	}, nil
}

// Publish публикует событие в RabbitMQ
func (p *RabbitMQPublisher) Publish(ctx context.Context, event interface{}) error {
	if p == nil || p.channel == nil {
		return fmt.Errorf("publisher not initialized")
	}

	// Определение routing key на основе типа события
	var routingKey string
	switch event.(type) {
	case *events.ForeignCompanyCreatedEvent:
		routingKey = "foreign_company.created"
	case *events.ForeignCompanyUpdatedEvent:
		routingKey = "foreign_company.updated"
	case *events.ForeignCompanyDeletedEvent:
		routingKey = "foreign_company.deleted"
	case *events.ForeignCompanyDeactivatedEvent:
		routingKey = "foreign_company.deactivated"
	default:
		routingKey = "foreign_company.unknown"
		p.logger.Warn("unknown event type", zap.String("type", fmt.Sprintf("%T", event)))
	}

	// Сериализация события
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to serialize event: %w", err)
	}

	// Публикация
	err = p.channel.PublishWithContext(
		ctx,
		"foreign_company.events", // exchange
		routingKey,               // routing key
		false,                    // mandatory
		false,                    // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	p.logger.Debug("event published",
		zap.String("routing_key", routingKey),
		zap.Int("body_size", len(body)),
	)

	return nil
}

// Close закрывает соединение с RabbitMQ
func (p *RabbitMQPublisher) Close() error {
	if p.channel != nil {
		p.channel.Close()
	}
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}
