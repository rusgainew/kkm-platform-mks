FROM golang:1.24-alpine AS builder

WORKDIR /build

# Copy proto-lib and service
COPY proto-lib/ ../proto-lib/
COPY services/bank-account-server/ ./bank-account-server/

WORKDIR /build/bank-account-server
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bank-account-server ./cmd/main.go

# Runtime
FROM alpine:3.18

RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -S appuser && adduser -S appuser -G appuser

WORKDIR /app

COPY --from=builder /build/bank-account-server/bank-account-server .
COPY --from=builder /build/bank-account-server/migrations ./migrations

# Set permissions
RUN chown -R appuser:appuser /app

USER appuser

EXPOSE 50054 9094

# Health check
HEALTHCHECK --interval=10s --timeout=3s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:9094/metrics || exit 1

CMD ["./bank-account-server"]
