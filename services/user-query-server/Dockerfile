FROM golang:1.24-alpine AS builder

ENV GOPROXY=https://goproxy.io,https://proxy.golang.org,direct
ENV GOSUMDB=off

WORKDIR /build

RUN apk add --no-cache git

# Copy proto-lib go.mod first
COPY proto-lib/go.mod proto-lib/go.sum ../proto-lib/

# Copy only go.mod and go.sum for caching
COPY services/user-query-server/go.mod services/user-query-server/go.sum ./user-query-server/

WORKDIR /build/user-query-server
RUN go mod download

# Copy rest of code
WORKDIR /build
COPY proto-lib/ ../proto-lib/
COPY services/user-query-server/ ./user-query-server/

WORKDIR /build/user-query-server
RUN go build -o user-query-server ./cmd/main.go

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /build/user-query-server/user-query-server .

EXPOSE 50061 9102

CMD ["./user-query-server"]
