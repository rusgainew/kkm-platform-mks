# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Установка зависимостей для сборки
RUN apk add --no-cache git

# Копируем go.mod и go.sum
COPY go.mod go.sum ./

# Копируем все остальное
COPY . .

# Собираем приложение
RUN go build -o document-server ./cmd/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Установка ca-certificates для HTTPS
RUN apk --no-cache add ca-certificates

# Копируем бинарник из builder stage
COPY --from=builder /app/document-server .

EXPOSE 50054 9094

CMD ["./document-server"]
