# Bank Account Query Server

CQRS Query service for reading bank account data with pagination.

## Features

- gRPC BankAccountQueryService
- ListBankAccounts(PageInfo) → APIResponse
- PostgreSQL read-optimized queries
- Health checks

## Build

```bash
cd bank-account-query-server
go build -o bank-account-query-server ./cmd
```

## Run

```bash
export BANK_ACCOUNT_QUERY_DB_URL="postgres://postgres:postgres@localhost:5432/bankaccountdb?sslmode=disable"
./bank-account-query-server
```

## Environment Variables

- `BANK_ACCOUNT_QUERY_GRPC_PORT=50054`
- `BANK_ACCOUNT_QUERY_DB_DRIVER=postgres`
- `BANK_ACCOUNT_QUERY_DB_URL`

## Testing

```bash
grpcurl -plaintext -d '{"page":1,"per_page":10}' \
  localhost:50054 api.BankAccountQueryService/ListBankAccounts
```
