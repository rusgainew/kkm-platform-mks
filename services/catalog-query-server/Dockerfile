FROM golang:1.24-alpine AS builder

WORKDIR /build
# Copy proto-lib and service sources
COPY proto-lib/ ../proto-lib/
COPY services/catalog-query-server/ ./catalog-query-server/

WORKDIR /build/catalog-query-server
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o catalog-query-server ./cmd

FROM alpine:3.18
RUN apk add --no-cache ca-certificates wget

# Create non-root user
RUN addgroup -S appuser && adduser -S appuser -G appuser

WORKDIR /root/
COPY --from=builder /build/catalog-query-server/catalog-query-server .

# Set permissions
RUN chown -R appuser:appuser /root

USER appuser

EXPOSE 50058 9098

# Health check
HEALTHCHECK --interval=10s --timeout=3s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:9098/metrics || exit 1

CMD ["./catalog-query-server"]
