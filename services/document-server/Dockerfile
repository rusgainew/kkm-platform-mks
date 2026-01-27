# Build stage
FROM golang:1.24-alpine AS builder

ENV GOPROXY=https://goproxy.io,https://proxy.golang.org,direct
ENV GOSUMDB=off

WORKDIR /build

# Установка зависимостей для сборки
RUN apk add --no-cache git

# Copy proto-lib go.mod first
COPY proto-lib/go.mod proto-lib/go.sum ../proto-lib/

# Copy only go.mod and go.sum
COPY services/document-server/go.mod services/document-server/go.sum ./document-server/

WORKDIR /build/document-server
RUN go mod download

# Copy rest of code
WORKDIR /build
COPY proto-lib/ ../proto-lib/
COPY services/document-server/ ./document-server/

WORKDIR /build/document-server
# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o document-server ./cmd/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Установка ca-certificates для HTTPS
RUN apk --no-cache add ca-certificates

# Копируем бинарник из builder stage
COPY --from=builder /build/document-server/document-server .

EXPOSE 50054 9094

CMD ["./document-server"]
