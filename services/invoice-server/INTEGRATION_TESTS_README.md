# Invoice Server Integration Tests

Smoke tests для all invoice-server gRPC handlers используя grpcurl.

## Prerequisites

- `grpcurl` установлен и доступен в PATH
- Invoice-server работает на `localhost:50052`
- Proto-definitions скомпилированы и доступны

## Running Tests

```bash
./integration_tests.sh
```

## Test Cases

| #   | Handler                | Input                      | Expected                            |
| --- | ---------------------- | -------------------------- | ----------------------------------- |
| 1   | CreateInvoice          | Valid invoice data         | Success (invoice_uuid returned)     |
| 2   | ListInvoices           | PageInfo (page=0, size=10) | List of invoices + pagination       |
| 3   | ListInvoicesWithFilter | InvoiceFilterRequest       | Filtered list + pagination          |
| 4   | SearchInvoices         | SearchRequest              | Matching invoices + pagination      |
| 5   | GetInvoiceByNumber     | InvoiceNumberRequest       | Single invoice or NotFound          |
| 6   | GetInvoicesByDateRange | DateRangeRequest           | Invoices in date range + pagination |
| 7   | ListInvoiceDetails     | InvoiceDetailsRequest      | Detail list + pagination            |

## Test Data

### CreateInvoice Request

```json
{
  "invoice_number": "INV-TEST-001",
  "invoice_date": "2024-01-15",
  "details": [
    {
      "goods_name": "Test Goods",
      "base_count": 1,
      "price": 100,
      "amount": 100
    }
  ],
  "legal_person": {
    "pin": "12345678901234"
  },
  "contractor": {
    "pin": "98765432109876"
  }
}
```

## Notes

- Test 7 now uses `InvoiceDetailsRequest` with `invoice_uuid` captured from CreateInvoice response

## TODO

- Add UpdateInvoice test (requires invoice ID from CreateInvoice)
- Add SignInvoice test (requires signed document)
- Add AcceptOrRejectInvoice test (requires signed invoice)
- Add RevokeInvoice test (requires revocable invoice)
- Add test database setup/teardown
- Add concurrent load testing
- Add test result aggregation and reporting
