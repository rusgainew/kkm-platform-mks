FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /build

# Copy proto-lib and service sources from repo root
COPY proto-lib/ ../proto-lib/
COPY services/catalog-server/ ./catalog-server/

# Build
WORKDIR /build/catalog-server
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o catalog-server ./cmd/main.go

# Runtime image
FROM alpine:3.18

RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -S appuser && adduser -S appuser -G appuser

WORKDIR /app

# Copy binary and migrations
COPY --from=builder /build/catalog-server/catalog-server .
COPY --from=builder /build/catalog-server/migrations ./migrations

# Set permissions
RUN chown -R appuser:appuser /app

USER appuser

EXPOSE 50053 9093

# Health check
HEALTHCHECK --interval=10s --timeout=3s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:9093/metrics || exit 1

CMD ["./catalog-server"]
