# Stage 1: Build
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git make ca-certificates

WORKDIR /app

# Copy proto-lib
COPY proto-lib/ ./proto-lib

# Copy go.mod and go.sum
COPY services/analytics-server/go.mod services/analytics-server/go.sum ./services/analytics-server/

WORKDIR /app/services/analytics-server
RUN go mod download

# Copy source code
COPY services/analytics-server/ .

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o analytics-server ./cmd

# Stage 2: Run
FROM alpine:3.18

RUN apk --no-cache add ca-certificates tzdata && \
  addgroup -g 1000 analytics && \
  adduser -D -u 1000 -G analytics analytics

WORKDIR /home/analytics

# Copy binary
COPY --from=builder /app/services/analytics-server/analytics-server .

# Copy migrations
COPY --from=builder /app/services/analytics-server/internal/infrastructure/repository/migrations ./migrations

RUN chown -R analytics:analytics /home/analytics

USER analytics

EXPOSE 50070 9115

CMD ["./analytics-server"]
