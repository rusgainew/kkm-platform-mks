# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Установка зависимостей сборки
RUN apk add --no-cache git

# Копирование go.mod и go.sum
COPY services/foreign-company-server/go.mod services/foreign-company-server/go.sum ./
COPY proto-lib ../proto-lib/

# Загрузка зависимостей
RUN go mod download

# Копирование исходного кода
COPY services/foreign-company-server/ ./
RUN go mod tidy

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o foreign-company-server ./cmd/main.go

# Runtime stage
FROM alpine:3.18

RUN apk --no-cache add ca-certificates tzdata wget

# Create non-root user
RUN addgroup -S appuser && adduser -S appuser -G appuser

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/foreign-company-server .

# Copy migrations
COPY --from=builder /app/migrations ./migrations

# Install grpc_health_probe for health checks
RUN wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.4.19/grpc_health_probe-linux-amd64 && \
  chmod +x /bin/grpc_health_probe

# Set permissions
RUN chown -R appuser:appuser /root

USER appuser

# Expose ports
EXPOSE 50056 9096

# Health check
HEALTHCHECK --interval=10s --timeout=3s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:9096/metrics || exit 1

# Run
CMD ["./foreign-company-server"]
