FROM golang:1.24-alpine AS builder

WORKDIR /build

# Копируем только необходимые части репозитория
COPY proto-lib/ ../proto-lib/
COPY services/document-query-server/ ./document-query-server/

WORKDIR /build/document-query-server
RUN apk add --no-cache git
RUN go mod download && go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o document-query-server ./cmd/main.go

FROM alpine:3.18
RUN apk --no-cache add ca-certificates wget

WORKDIR /root/
COPY --from=builder /build/document-query-server/document-query-server .

EXPOSE 50062 9103

# Health check-style parity with other services
HEALTHCHECK --interval=10s --timeout=3s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:9103/metrics || exit 1

CMD ["./document-query-server"]
