FROM golang:1.24-alpine AS builder

ENV GOPROXY=https://goproxy.io,https://proxy.golang.org,direct
ENV GOSUMDB=off

RUN apk add --no-cache git ca-certificates

WORKDIR /build

# Copy proto-lib go.mod first
COPY proto-lib/go.mod proto-lib/go.sum ../proto-lib/

# Copy only go.mod and go.sum
COPY services/bank-account-query-server/go.mod services/bank-account-query-server/go.sum ./bank-account-query-server/

WORKDIR /build/bank-account-query-server
RUN go mod download

# Copy rest of code
WORKDIR /build
COPY proto-lib/ ../proto-lib/
COPY services/bank-account-query-server/ ./bank-account-query-server/

WORKDIR /build/bank-account-query-server
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bank-account-query-server ./cmd

FROM alpine:3.18
RUN apk add --no-cache ca-certificates wget

# Create non-root user
RUN addgroup -S appuser && adduser -S appuser -G appuser

WORKDIR /root/
COPY --from=builder /build/bank-account-query-server/bank-account-query-server .

# Set permissions
RUN chown -R appuser:appuser /root

USER appuser

EXPOSE 50059 9099

# Health check
HEALTHCHECK --interval=10s --timeout=3s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:9099/metrics || exit 1

CMD ["./bank-account-query-server"]
