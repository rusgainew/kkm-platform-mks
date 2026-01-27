FROM golang:1.24-alpine AS builder

ENV GOPROXY=https://goproxy.io,https://proxy.golang.org,direct
ENV GOSUMDB=off

RUN apk add --no-cache git

WORKDIR /build

# Copy proto-lib go.mod first
COPY proto-lib/go.mod proto-lib/go.sum ../proto-lib/

# Copy only go.mod and go.sum
COPY services/document-query-server/go.mod services/document-query-server/go.sum ./document-query-server/

WORKDIR /build/document-query-server
RUN go mod download && go mod tidy

# Копируем остальной код
WORKDIR /build
COPY proto-lib/ ../proto-lib/
COPY services/document-query-server/ ./document-query-server/

WORKDIR /build/document-query-server
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
