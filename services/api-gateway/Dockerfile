# Stage 1: Build
FROM golang:1.24-alpine AS builder

# Install dependencies
RUN apk add --no-cache git make ca-certificates

WORKDIR /app

# Copy local module proto-lib (used via replace ../../proto-lib)
COPY proto-lib/ ./proto-lib

# Copy local module pkg (used via replace ../pkg)
COPY services/pkg/ ./services/pkg

# Copy go.mod/go.sum for api-gateway
COPY services/api-gateway/go.mod services/api-gateway/go.sum ./services/api-gateway/

# Download dependencies
WORKDIR /app/services/api-gateway
RUN go mod download

# Copy source code
COPY services/api-gateway/ .

# Ensure go.sum up to date
RUN go mod tidy

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o api-gateway ./cmd/api

# Stage 2: Run
FROM alpine:3.18

RUN apk --no-cache add ca-certificates curl

# Create non-root user
RUN addgroup -S api-gateway && adduser -S api-gateway -G api-gateway

WORKDIR /home/api-gateway

# Copy binary from builder
COPY --from=builder /app/services/api-gateway/api-gateway .

# Copy documentation files for Swagger
COPY --from=builder /app/services/api-gateway/docs ./docs

# Set permissions
RUN chown -R api-gateway:api-gateway /home/api-gateway

USER api-gateway

EXPOSE 8080 9090

HEALTHCHECK --interval=10s --timeout=5s --retries=3 --start-period=30s \
  CMD curl -f http://localhost:8080/health || exit 1

CMD ["./api-gateway"]
