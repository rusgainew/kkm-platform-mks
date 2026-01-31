# Goroutine Leak Fix Report

## Issue #3: Goroutine Leaks in RabbitMQ Consumers

### Problem Description

RabbitMQ consumers were creating goroutines that could not be properly terminated on service shutdown due to using `context.Background()` instead of the parent context. This caused memory leaks over time.

### Root Cause

The functions `republishWithRetry` and `sendToDLQ` created new contexts from `context.Background()`:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
```

This meant that when the parent context was cancelled (e.g., during service shutdown), these goroutines would continue to run until their own timeout expired, preventing graceful shutdown and causing resource leaks.

### Solution

Modified the function signatures to accept a parent `context.Context` parameter and create child contexts from it:

```go
func (c *RabbitMQConsumer) republishWithRetry(ctx context.Context, msg amqp.Delivery, retryCount int, delay time.Duration)
func (c *RabbitMQConsumer) sendToDLQ(ctx context.Context, msg amqp.Delivery, retryCount int) error
```

Now uses:

```go
publishCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
```

This ensures proper context hierarchy and allows parent cancellation to propagate down to child goroutines.

### Files Modified

1. **services/user-query-server/internal/infrastructure/messaging/rabbitmq_consumer.go**
   - Modified `republishWithRetry` signature (line ~452)
   - Modified `sendToDLQ` signature (line ~495)
   - Updated call sites in `processMessages` (lines ~410, 418)

2. **services/document-query-server/internal/infrastructure/messaging/rabbitmq_consumer.go**
   - Modified `republishWithRetry` signature (line ~454)
   - Modified `sendToDLQ` signature (line ~497)
   - Updated call sites in `processMessages` (lines ~410, 418)

### Testing

- [x] user-query-server compiles successfully
- [x] document-query-server compiles successfully
- [x] Context hierarchy properly maintained
- [x] Graceful shutdown now possible

### Impact

- **Memory Safety**: Eliminates goroutine leaks during service shutdown
- **Resource Cleanup**: Proper context cancellation ensures all operations terminate
- **Production Readiness**: Services can now shut down gracefully without orphaned goroutines
- **Observability**: Better control over goroutine lifecycle

### Time Spent

Approximately 15 minutes

### Next Steps

Proceed to Issue #4: Missing context deadline propagation in 20+ handler locations
