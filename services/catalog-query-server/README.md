# Catalog Query Server

CQRS Query service for reading catalog data with pagination.

## Features

- gRPC CatalogQueryService
- ListCatalogs(PageInfo) → APIResponse
- PostgreSQL read-optimized queries
- Health checks

## Build

```bash
cd catalog-query-server
go build -o catalog-query-server ./cmd
```

## Run

```bash
export CATALOG_QUERY_DB_URL="postgres://postgres:postgres@localhost:5432/catalogdb?sslmode=disable"
./catalog-query-server
```

## Environment Variables

- `CATALOG_QUERY_GRPC_PORT=50055`
- `CATALOG_QUERY_DB_DRIVER=postgres`
- `CATALOG_QUERY_DB_URL`

## Testing

```bash
grpcurl -plaintext -d '{"page":1,"per_page":10}' \
  localhost:50055 api.CatalogQueryService/ListCatalogs
```
