package messaging

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// EventsProcessedTotal tracks total number of events processed
	EventsProcessedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "user_events_processed_total",
			Help: "Total number of user events processed",
		},
		[]string{"event_type", "status"}, // status: success, error, duplicate
	)

	// EventProcessingDuration tracks event processing time
	EventProcessingDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "user_event_processing_duration_seconds",
			Help:    "Duration of event processing in seconds",
			Buckets: prometheus.DefBuckets, // Default buckets: 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10
		},
		[]string{"event_type"},
	)

	// EventValidationErrors tracks validation errors
	EventValidationErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "user_event_validation_errors_total",
			Help: "Total number of event validation errors",
		},
		[]string{"event_type", "error_type"},
	)

	// ReconnectionAttempts tracks RabbitMQ reconnection attempts
	ReconnectionAttempts = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "user_consumer_reconnection_attempts_total",
			Help: "Total number of RabbitMQ reconnection attempts",
		},
	)

	// ReconnectionDelay tracks current reconnection delay
	ReconnectionDelay = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "user_consumer_reconnection_delay_seconds",
			Help: "Current reconnection delay in seconds",
		},
	)

	// ConnectionStatus tracks current connection status (1 = connected, 0 = disconnected)
	ConnectionStatus = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "user_consumer_connection_status",
			Help: "Current RabbitMQ connection status (1 = connected, 0 = disconnected)",
		},
	)

	// MessagesRetried tracks messages sent for retry
	MessagesRetried = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "user_messages_retried_total",
			Help: "Total number of messages retried",
		},
		[]string{"retry_count"},
	)

	// MessagesSentToDLQ tracks messages sent to Dead Letter Queue
	MessagesSentToDLQ = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "user_messages_dlq_total",
			Help: "Total number of messages sent to DLQ",
		},
		[]string{"event_type"},
	)

	// IdempotencyChecks tracks idempotency check results
	IdempotencyChecks = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "user_idempotency_checks_total",
			Help: "Total number of idempotency checks",
		},
		[]string{"result"}, // result: duplicate, new, error
	)
)
