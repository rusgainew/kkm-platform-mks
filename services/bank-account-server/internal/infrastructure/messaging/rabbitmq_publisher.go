package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
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

	// Объявление exchange для банковских счетов
	err = channel.ExchangeDeclare(
		"bank_account.events", // name
		"topic",               // type
		true,                  // durable
		false,                 // auto-deleted
		false,                 // internal
		false,                 // no-wait
		nil,                   // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	logger.Info("RabbitMQ publisher connected", zap.String("exchange", "bank_account.events"))

	return &RabbitMQPublisher{
		conn:    conn,
		channel: channel,
		logger:  logger,
	}, nil
}

// Publish публикует событие в RabbitMQ
func (p *RabbitMQPublisher) Publish(ctx context.Context, event interface{}) error {
	// Определение routing key на основе типа события
	routingKey := getRoutingKey(event)

	// Сериализация события
	body, err := serializeEvent(event)
	if err != nil {
		return fmt.Errorf("failed to serialize event: %w", err)
	}

	// Публикация
	err = p.channel.PublishWithContext(
		ctx,
		"bank_account.events", // exchange
		routingKey,            // routing key
		false,                 // mandatory
		false,                 // immediate
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

func getRoutingKey(event interface{}) string {
	data, _ := json.Marshal(event)
	var base struct {
		EventType string `json:"event_type"`
	}
	json.Unmarshal(data, &base)
	if base.EventType != "" {
		return base.EventType
	}
	return "bank_account.unknown"
}

func serializeEvent(event interface{}) ([]byte, error) {
	if e, ok := event.(interface{ ToJSON() ([]byte, error) }); ok {
		return e.ToJSON()
	}
	return json.Marshal(event)
}
