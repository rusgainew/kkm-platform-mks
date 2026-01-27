# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make ca-certificates

# Set Go proxy to direct mode to avoid DNS issues
ENV GOPROXY=https://goproxy.io,https://proxy.golang.org,direct
ENV GOSUMDB=off

WORKDIR /build

# Copy proto-lib go.mod first for better caching
COPY proto-lib/go.mod proto-lib/go.sum ../proto-lib/

# Copy only go.mod and go.sum for dependency caching
COPY services/invoice-server/go.mod services/invoice-server/go.sum ./invoice-server/

WORKDIR /build/invoice-server
RUN go mod download

# Now copy the rest of proto-lib and service code
WORKDIR /build
COPY proto-lib/ ../proto-lib/
COPY services/invoice-server/ ./invoice-server/

# Build the application
WORKDIR /build/invoice-server
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o invoice-server ./cmd/main.go

# Final stage
FROM alpine:3.18

RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -S appuser && adduser -S appuser -G appuser

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/invoice-server/invoice-server .
COPY --from=builder /build/invoice-server/migrations ./migrations

# Set permissions
RUN chown -R appuser:appuser /app

USER appuser

# Expose ports
EXPOSE 50055 9095

# Health check
HEALTHCHECK --interval=10s --timeout=3s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:9095/metrics || exit 1

# Run the application
CMD ["./invoice-server"]
