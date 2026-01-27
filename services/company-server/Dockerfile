# Stage 1: Build
FROM golang:1.24-alpine AS builder

ENV GOPROXY=https://goproxy.io,https://proxy.golang.org,direct
ENV GOSUMDB=off

# Install dependencies
RUN apk add --no-cache git make

WORKDIR /app

# Copy proto-lib go.mod first
COPY proto-lib/go.mod proto-lib/go.sum ../proto-lib/

# Copy go mod files
COPY services/company-server/go.* ./

# Download dependencies
RUN go mod download

# Copy proto-lib and source code
COPY proto-lib/ ../proto-lib/
COPY services/company-server/ ./
WORKDIR /app
RUN go mod tidy

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o company-server ./cmd/main.go

# Stage 2: Run
FROM alpine:3.18

RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -S appuser && adduser -S appuser -G appuser

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/company-server .

# Copy migrations
COPY --from=builder /app/migrations ./migrations

# Set permissions
RUN chown -R appuser:appuser /root

USER appuser

# Expose ports
EXPOSE 50052 9092

# Health check
HEALTHCHECK --interval=10s --timeout=3s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:9092/metrics || exit 1

# Run the application
CMD ["./company-server"]
