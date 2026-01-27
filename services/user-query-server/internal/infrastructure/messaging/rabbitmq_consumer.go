package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

const (
	maxRetries        = 3
	retryCountHeader  = "x-retry-count"
	dlqExchangeSuffix = ".dlq"
)

// UserEvent represents a user domain event from RabbitMQ
type UserEvent struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Timestamp string                 `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

// EventHandler processes received events
type EventHandler interface {
	HandleUserRegistered(ctx context.Context, event UserEvent) error
	HandleUserLoggedIn(ctx context.Context, event UserEvent) error
	HandleTokenRefreshed(ctx context.Context, event UserEvent) error
}

// RabbitMQConsumer consumes user events from RabbitMQ
type RabbitMQConsumer struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	queue     string
	exchange  string
	logger    *zap.Logger
	handler   EventHandler
	validator *EventValidator
	closeMu   sync.Mutex
	closed    bool
	processed map[string]bool // In-memory tracking of processed events
	procMu    sync.RWMutex
}

// NewRabbitMQConsumer creates a new RabbitMQ consumer for user events
func NewRabbitMQConsumer(
	url, queueName, exchange string,
	logger *zap.Logger,
	handler EventHandler,
) (*RabbitMQConsumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Declare exchange (idempotent, no error if exists)
	err = channel.ExchangeDeclare(
		exchange, // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	// Declare queue (idempotent)
	q, err := channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	// Bind queue to exchange with routing keys for all user events
	routingKeys := []string{
		"user.registered",
		"user.logged_in",
		"token.refreshed",
	}

	for _, key := range routingKeys {
		err = channel.QueueBind(
			q.Name,   // queue name
			key,      // routing key (event type)
			exchange, // exchange
			false,
			nil,
		)
		if err != nil {
			channel.Close()
			conn.Close()
			return nil, fmt.Errorf("failed to bind queue: %w", err)
		}
	}

	// Declare DLQ (Dead Letter Queue) exchange and queue
	dlqExchange := exchange + dlqExchangeSuffix
	err = channel.ExchangeDeclare(
		dlqExchange, // name
		"topic",     // type
		true,        // durable
		false,       // auto-deleted
		false,       // internal
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare DLQ exchange: %w", err)
	}

	// Declare DLQ with expiration (7 days)
	dlqName := queueName + ".dlq"
	_, err = channel.QueueDeclare(
		dlqName, // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		amqp.Table{
			"x-message-ttl": int64(7 * 24 * 60 * 60 * 1000), // 7 days in milliseconds
		},
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare DLQ queue: %w", err)
	}

	// Bind DLQ to DLQ exchange for all routing keys
	for _, key := range routingKeys {
		err = channel.QueueBind(
			dlqName,     // queue name
			key,         // routing key
			dlqExchange, // exchange
			false,
			nil,
		)
		if err != nil {
			channel.Close()
			conn.Close()
			return nil, fmt.Errorf("failed to bind DLQ queue: %w", err)
		}
	}

	// Configure main queue with DLQ routing
	_, err = channel.QueueDelete(q.Name, false, false, false)
	if err != nil && err.Error() != "NOT_FOUND" {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to delete existing queue: %w", err)
	}

	q, err = channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		amqp.Table{
			"x-dead-letter-exchange": dlqExchange,
		},
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to redeclare queue with DLQ: %w", err)
	}

	// Re-bind main queue after recreation
	for _, key := range routingKeys {
		err = channel.QueueBind(
			q.Name,   // queue name
			key,      // routing key (event type)
			exchange, // exchange
			false,
			nil,
		)
		if err != nil {
			channel.Close()
			conn.Close()
			return nil, fmt.Errorf("failed to bind queue after DLQ setup: %w", err)
		}
	}

	// Set QoS (prefetch count for fair dispatch)
	err = channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to set QoS: %w", err)
	}

	consumer := &RabbitMQConsumer{
		conn:      conn,
		channel:   channel,
		queue:     q.Name,
		exchange:  exchange,
		logger:    logger,
		handler:   handler,
		validator: NewEventValidator(logger),
		processed: make(map[string]bool),
	}

	logger.Info("RabbitMQ consumer initialized",
		zap.String("queue", q.Name),
		zap.String("exchange", exchange),
		zap.Strings("routing_keys", routingKeys))

	return consumer, nil
}

// Start begins consuming messages from RabbitMQ with auto-reconnect on failure
func (c *RabbitMQConsumer) Start(ctx context.Context) error {
	c.closeMu.Lock()
	if c.closed {
		c.closeMu.Unlock()
		return fmt.Errorf("consumer already closed")
	}
	c.closeMu.Unlock()

	// Start reconnection loop in goroutine
	go c.startWithReconnect(ctx)

	return nil
}

// startWithReconnect handles reconnection logic with exponential backoff
func (c *RabbitMQConsumer) startWithReconnect(ctx context.Context) {
	reconnectAttempt := 0
	maxReconnectAttempts := 0 // 0 = unlimited (retry forever)

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("Consumer context cancelled, stopping")
			ConnectionStatus.Set(0)
			return
		default:
		}

		// Check if consumer is closed
		c.closeMu.Lock()
		if c.closed {
			c.closeMu.Unlock()
			c.logger.Info("Consumer closed, stopping reconnect loop")
			ConnectionStatus.Set(0)
			return
		}
		c.closeMu.Unlock()

		// Attempt to start consuming
		if err := c.startConsuming(ctx); err != nil {
			ConnectionStatus.Set(0)
			ReconnectionAttempts.Inc()
			reconnectAttempt++
			delay := c.calculateReconnectBackoff(reconnectAttempt)
			ReconnectionDelay.Set(delay.Seconds())

			if maxReconnectAttempts > 0 && reconnectAttempt >= maxReconnectAttempts {
				c.logger.Error("Max reconnect attempts exceeded, stopping consumer",
					zap.Int("attempt", reconnectAttempt),
					zap.Int("max_attempts", maxReconnectAttempts))
				return
			}

			c.logger.Warn("Reconnection failed, retrying",
				zap.Error(err),
				zap.Int("attempt", reconnectAttempt),
				zap.Duration("backoff", delay))

			// Wait before retrying
			select {
			case <-time.After(delay):
				continue
			case <-ctx.Done():
				return
			}
		} else {
			// Connection successful, reset attempt counter
			reconnectAttempt = 0
			ConnectionStatus.Set(1)
			ReconnectionDelay.Set(0)
			c.logger.Info("Consumer successfully connected and consuming")
		}
	}
}

// startConsuming starts the actual consuming process and blocks until connection closes
func (c *RabbitMQConsumer) startConsuming(ctx context.Context) error {
	msgs, err := c.channel.Consume(
		c.queue, // queue name
		"",      // consumer tag
		false,   // auto-ack (false for manual ack)
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	c.logger.Info("Consumer started, listening for events", zap.String("queue", c.queue))

	// Monitor connection close
	connCloseCh := c.conn.NotifyClose(make(chan *amqp.Error, 1))

	// Process messages in goroutine
	msgProcessingDone := make(chan struct{})
	go func() {
		c.processMessages(ctx, msgs)
		close(msgProcessingDone)
	}()

	// Wait for either context cancellation or connection close
	select {
	case <-ctx.Done():
		c.logger.Info("Context cancelled, stopping message processing")
		return ctx.Err()
	case connErr := <-connCloseCh:
		if connErr != nil {
			return fmt.Errorf("RabbitMQ connection closed: %w", connErr)
		}
		return fmt.Errorf("RabbitMQ connection closed")
	case <-msgProcessingDone:
		return fmt.Errorf("message processing ended unexpectedly")
	}
}

// calculateReconnectBackoff calculates exponential backoff delay (2^attempt, capped at 30s)
func (c *RabbitMQConsumer) calculateReconnectBackoff(attempt int) time.Duration {
	baseDelay := time.Second
	delay := baseDelay * time.Duration(1<<uint(attempt)) // 2^attempt

	// Cap at 30 seconds
	maxDelay := 30 * time.Second
	if delay > maxDelay {
		delay = maxDelay
	}

	return delay
}

// processMessages handles incoming messages
func (c *RabbitMQConsumer) processMessages(ctx context.Context, msgs <-chan amqp.Delivery) {
	for {
		select {
		case <-ctx.Done():
			c.logger.Info("Stopping message processor")
			return
		case msg, ok := <-msgs:
			if !ok {
				c.logger.Info("Message channel closed")
				return
			}

			// Get retry count from headers
			retryCount := 0
			if msg.Headers != nil {
				if val, exists := msg.Headers[retryCountHeader]; exists {
					if count, err := strconv.Atoi(val.(string)); err == nil {
						retryCount = count
					}
				}
			}

			// Process the message
			if err := c.handleMessage(ctx, msg); err != nil {
				c.logger.Error("Failed to handle message",
					zap.Error(err),
					zap.String("message_id", msg.MessageId),
					zap.String("routing_key", msg.RoutingKey),
					zap.Int("retry_count", retryCount))

				// Check if max retries exceeded
				if retryCount >= maxRetries {
					c.logger.Warn("Max retries exceeded, sending to DLQ",
						zap.String("message_id", msg.MessageId),
						zap.String("routing_key", msg.RoutingKey))

					// Send to DLQ by republishing with DLQ exchange
					if dlqErr := c.sendToDLQ(ctx, msg, retryCount); dlqErr != nil {
						c.logger.Error("Failed to send message to DLQ", zap.Error(dlqErr))
					} else {
						MessagesSentToDLQ.WithLabelValues(msg.RoutingKey).Inc()
					}

					// Acknowledge the original message (no requeue)
					msg.Ack(false)
				} else {
					// Increment retry count and requeue with exponential backoff
					delay := c.calculateBackoff(retryCount)
					c.republishWithRetry(ctx, msg, retryCount+1, delay)
					MessagesRetried.WithLabelValues(strconv.Itoa(retryCount + 1)).Inc()
					msg.Ack(false)
				}
			} else {
				// Positive acknowledge
				msg.Ack(false)
			}
		}
	}
}

// calculateBackoff calculates exponential backoff delay (2^retryCount, capped at 60s)
func (c *RabbitMQConsumer) calculateBackoff(retryCount int) time.Duration {
	baseDelay := time.Second
	delay := baseDelay * time.Duration(1<<uint(retryCount)) // 2^retryCount

	// Cap at 60 seconds
	maxDelay := 60 * time.Second
	if delay > maxDelay {
		delay = maxDelay
	}

	return delay
}

// republishWithRetry republishes message with incremented retry count and backoff
func (c *RabbitMQConsumer) republishWithRetry(ctx context.Context, msg amqp.Delivery, retryCount int, delay time.Duration) {
	// Create new publishing with updated headers
	headers := amqp.Table{}
	if msg.Headers != nil {
		for k, v := range msg.Headers {
			headers[k] = v
		}
	}
	headers[retryCountHeader] = strconv.Itoa(retryCount)

	newMsg := amqp.Publishing{
		ContentType:  msg.ContentType,
		Body:         msg.Body,
		Headers:      headers,
		Expiration:   fmt.Sprintf("%d", int(delay.Milliseconds())),
		DeliveryMode: amqp.Persistent,
	}

	// Use context with timeout from parent context
	publishCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := c.channel.PublishWithContext(
		publishCtx,
		c.exchange,     // exchange
		msg.RoutingKey, // routing key
		false,          // mandatory
		false,          // immediate
		newMsg,
	)
	if err != nil {
		c.logger.Error("Failed to republish message with retry",
			zap.Error(err),
			zap.String("message_id", msg.MessageId),
			zap.Int("retry_count", retryCount))
	} else {
		c.logger.Debug("Message republished with retry",
			zap.String("message_id", msg.MessageId),
			zap.Int("retry_count", retryCount),
			zap.Duration("backoff", delay))
	}
}

// sendToDLQ sends message to Dead Letter Queue
func (c *RabbitMQConsumer) sendToDLQ(ctx context.Context, msg amqp.Delivery, retryCount int) error {
	headers := amqp.Table{}
	if msg.Headers != nil {
		for k, v := range msg.Headers {
			headers[k] = v
		}
	}
	headers[retryCountHeader] = strconv.Itoa(retryCount)
	headers["x-dlq-timestamp"] = time.Now().UTC().String()

	dlqMsg := amqp.Publishing{
		ContentType:  msg.ContentType,
		Body:         msg.Body,
		Headers:      headers,
		DeliveryMode: amqp.Persistent,
	}

	// Use context with timeout from parent context
	publishCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	dlqExchange := c.exchange + dlqExchangeSuffix
	return c.channel.PublishWithContext(
		publishCtx,
		dlqExchange,    // DLQ exchange
		msg.RoutingKey, // same routing key
		false,
		false,
		dlqMsg,
	)
}

// isEventProcessed checks if event has already been processed (idempotency)
func (c *RabbitMQConsumer) isEventProcessed(ctx context.Context, eventID string) (bool, error) {
	c.procMu.RLock()
	defer c.procMu.RUnlock()

	return c.processed[eventID], nil
}

// markEventProcessed marks event as processed for idempotency
func (c *RabbitMQConsumer) markEventProcessed(ctx context.Context, eventID, eventType string) error {
	c.procMu.Lock()
	defer c.procMu.Unlock()

	c.processed[eventID] = true

	// Cleanup old entries if map gets too large (keep last 10000)
	if len(c.processed) > 10000 {
		// Simple cleanup: clear half of entries
		// In production, use LRU cache or time-based expiration
		for k := range c.processed {
			delete(c.processed, k)
			if len(c.processed) <= 5000 {
				break
			}
		}
	}

	return nil
}

// handleMessage processes a single message
func (c *RabbitMQConsumer) handleMessage(ctx context.Context, msg amqp.Delivery) error {
	var event UserEvent
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	// Start timing
	timer := prometheus.NewTimer(EventProcessingDuration.WithLabelValues(msg.RoutingKey))
	defer timer.ObserveDuration()

	// Validate event
	validationErrors := c.validator.ValidateUserEvent(event)
	if len(validationErrors) > 0 {
		c.validator.LogValidationErrors(event.ID, validationErrors)
		for _, valErr := range validationErrors {
			EventValidationErrors.WithLabelValues(event.Type, valErr.Field).Inc()
		}
		EventsProcessedTotal.WithLabelValues(event.Type, "validation_error").Inc()
		return fmt.Errorf("event validation failed: %d errors", len(validationErrors))
	}

	// Check for idempotency (skip if already processed)
	processed, err := c.isEventProcessed(ctx, event.ID)
	if err != nil {
		c.logger.Warn("Idempotency check failed, continuing with processing", zap.Error(err))
		IdempotencyChecks.WithLabelValues("error").Inc()
	} else if processed {
		IdempotencyChecks.WithLabelValues("duplicate").Inc()
	} else {
		IdempotencyChecks.WithLabelValues("new").Inc()
	}
	if processed {
		c.logger.Info("Event already processed, skipping",
			zap.String("event_id", event.ID),
			zap.String("event_type", event.Type))
		EventsProcessedTotal.WithLabelValues(event.Type, "duplicate").Inc()
		return nil
	}

	c.logger.Debug("Received event",
		zap.String("event_id", event.ID),
		zap.String("event_type", event.Type),
		zap.String("routing_key", msg.RoutingKey))

	// Process the event
	var processErr error
	switch event.Type {
	case "user.registered":
		processErr = c.handler.HandleUserRegistered(ctx, event)
	case "user.logged_in":
		processErr = c.handler.HandleUserLoggedIn(ctx, event)
	case "token.refreshed":
		processErr = c.handler.HandleTokenRefreshed(ctx, event)
	default:
		processErr = fmt.Errorf("unknown event type: %s", event.Type)
	}

	if processErr != nil {
		EventsProcessedTotal.WithLabelValues(event.Type, "error").Inc()
		return processErr
	}

	// Mark event as processed after successful handling
	if err := c.markEventProcessed(ctx, event.ID, event.Type); err != nil {
		c.logger.Error("Failed to mark event processed, but continuing", zap.Error(err))
		// Don't fail the entire message processing on idempotency marking failure
	}

	EventsProcessedTotal.WithLabelValues(event.Type, "success").Inc()
	return nil
}

// Close gracefully closes the consumer
func (c *RabbitMQConsumer) Close() error {
	c.closeMu.Lock()
	defer c.closeMu.Unlock()

	if c.closed {
		return nil
	}

	c.closed = true

	var errs []error

	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			c.logger.Error("Failed to close channel", zap.Error(err))
			errs = append(errs, err)
		}
	}

	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			c.logger.Error("Failed to close connection", zap.Error(err))
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to close consumer: %v", errs[0])
	}

	c.logger.Info("Consumer closed successfully")
	return nil
}
