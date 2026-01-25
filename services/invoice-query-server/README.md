# Invoice Query Server

CQRS Query service for reading invoice data with pagination.

## Features

- gRPC InvoiceQueryService
- ListInvoices(PageInfo) → APIResponse
- ListInvoiceDetails(PageInfo) → APIResponse
- PostgreSQL read-optimized queries
- Health checks
- Prometheus metrics (optional)

## Build

```bash
cd invoice-query-server
go build -o invoice-query-server ./cmd
```

## Run

```bash
# Local development
export INVOICE_QUERY_DB_URL="postgres://postgres:postgres@localhost:5432/invoicedb?sslmode=disable"
./invoice-query-server

# Docker
docker build -t invoice-query-server:latest .
docker run -p 50053:50053 invoice-query-server:latest
```

## Environment Variables

- `INVOICE_QUERY_GRPC_PORT=50053` - gRPC port
- `INVOICE_QUERY_DB_DRIVER=postgres` - Database driver
- `INVOICE_QUERY_DB_URL` - Database connection string

## Testing

```bash
grpcurl -plaintext localhost:50053 list
grpcurl -plaintext -d '{"page":1,"per_page":10}' \
  localhost:50053 api.InvoiceQueryService/ListInvoices
```
