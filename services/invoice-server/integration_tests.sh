#!/bin/bash

# Integration tests for invoice-server using grpcurl
# Tests all command and query handlers with smoke tests

set -e

# Allow override via env var; default to local 50053
INVOICE_SERVER="${INVOICE_SERVER:-localhost:50053}"

# Build grpcurl options as an array to handle headers safely
GRPCURL_OPTS=( -plaintext )
if [ -n "$AUTH_TOKEN" ]; then
  GRPCURL_OPTS+=( -H "authorization: Bearer $AUTH_TOKEN" )
fi
TIMEOUT="10s"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

if ! command -v jq >/dev/null 2>&1; then
  echo -e "${RED}jq is required for parsing grpcurl responses. Install jq and rerun.${NC}"
  exit 1
fi

echo -e "${YELLOW}Invoice Server Integration Tests${NC}"
echo "======================================="
echo "Target: $INVOICE_SERVER"
echo ""

CREATE_PAYLOAD='{
  "invoice_number": "INV-TEST-001",
  "invoice_date": "2024-01-15",
  "details": [
    { "goods_name": "Item A", "base_count": 1, "price": 100, "amount": 100 },
    { "goods_name": "Item B", "base_count": 1, "price": 200, "amount": 200 },
    { "goods_name": "Item C", "base_count": 1, "price": 300, "amount": 300 },
    { "goods_name": "Item D", "base_count": 1, "price": 400, "amount": 400 },
    { "goods_name": "Item E", "base_count": 1, "price": 500, "amount": 500 }
  ],
  "legal_person": { "pin": "12345678901234" },
  "contractor":  { "pin": "98765432109876" }
}'

# Test 1: CreateInvoice
echo -e "${YELLOW}Test 1: CreateInvoice${NC}"
CREATE_RESP=$(grpcurl "${GRPCURL_OPTS[@]}" -d "$CREATE_PAYLOAD" $INVOICE_SERVER api.InvoiceCommandService/CreateInvoice)
if [ $? -eq 0 ]; then
  echo -e "${GREEN}✓ CreateInvoice passed${NC}"
else
  echo -e "${RED}✗ CreateInvoice failed${NC}"
  echo "$CREATE_RESP"
  exit 1
fi
echo ""

INVOICE_UUID=$(echo "$CREATE_RESP" | jq -r '.documentUuid // .document_uuid // .invoice.documentUuid // .invoice.document_uuid')
if [ -z "$INVOICE_UUID" ] || [ "$INVOICE_UUID" = "null" ]; then
  echo -e "${RED}Failed to extract invoice UUID from CreateInvoice response${NC}"
  echo "$CREATE_RESP"
  exit 1
fi
echo -e "${GREEN}Captured invoice_uuid: $INVOICE_UUID${NC}"
echo ""

# Test 2: ListInvoices
echo -e "${YELLOW}Test 2: ListInvoices${NC}"
grpcurl "${GRPCURL_OPTS[@]}" -d '{"page": 0, "size": 10}' $INVOICE_SERVER api.InvoiceQueryService/ListInvoices
if [ $? -eq 0 ]; then
  echo -e "${GREEN}✓ ListInvoices passed${NC}"
else
  echo -e "${RED}✗ ListInvoices failed${NC}"
fi
echo ""

# Test 3: ListInvoicesWithFilter
echo -e "${YELLOW}Test 3: ListInvoicesWithFilter${NC}"
grpcurl "${GRPCURL_OPTS[@]}" -d @ $INVOICE_SERVER api.InvoiceQueryService/ListInvoicesWithFilter <<EOF
{
  "search_text": "INV",
  "page": 0,
  "size": 10
}
EOF
if [ $? -eq 0 ]; then
  echo -e "${GREEN}✓ ListInvoicesWithFilter passed${NC}"
else
  echo -e "${RED}✗ ListInvoicesWithFilter failed${NC}"
fi
echo ""

# Test 4: SearchInvoices
echo -e "${YELLOW}Test 4: SearchInvoices${NC}"
grpcurl "${GRPCURL_OPTS[@]}" -d @ $INVOICE_SERVER api.InvoiceQueryService/SearchInvoices <<EOF
{
  "search_text": "test",
  "page": 0,
  "size": 10
}
EOF
if [ $? -eq 0 ]; then
  echo -e "${GREEN}✓ SearchInvoices passed${NC}"
else
  echo -e "${RED}✗ SearchInvoices failed${NC}"
fi
echo ""

# Test 5: GetInvoiceByNumber
echo -e "${YELLOW}Test 5: GetInvoiceByNumber${NC}"
grpcurl "${GRPCURL_OPTS[@]}" -d '{"invoice_number": "INV-TEST-001"}' $INVOICE_SERVER api.InvoiceQueryService/GetInvoiceByNumber
if [ $? -eq 0 ]; then
  echo -e "${GREEN}✓ GetInvoiceByNumber passed${NC}"
else
  echo -e "${RED}✗ GetInvoiceByNumber failed${NC}"
fi
echo ""

# Test 6: GetInvoicesByDateRange
echo -e "${YELLOW}Test 6: GetInvoicesByDateRange${NC}"
grpcurl "${GRPCURL_OPTS[@]}" -d @ $INVOICE_SERVER api.InvoiceQueryService/GetInvoicesByDateRange <<EOF
{
  "date_from": "2024-01-01",
  "date_to": "2024-12-31",
  "page": 0,
  "size": 10
}
EOF
if [ $? -eq 0 ]; then
  echo -e "${GREEN}✓ GetInvoicesByDateRange passed${NC}"
else
  echo -e "${RED}✗ GetInvoicesByDateRange failed${NC}"
fi
echo ""

# Test 7: ListInvoiceDetails (now implemented)
echo -e "${YELLOW}Test 7: ListInvoiceDetails (pagination)${NC}"
DETAILS_RESP=$(grpcurl "${GRPCURL_OPTS[@]}" -d @ $INVOICE_SERVER api.InvoiceQueryService/ListInvoiceDetails <<EOF
{
  "invoice_uuid": "$INVOICE_UUID",
  "page": 1,
  "size": 2
}
EOF
)
if [ $? -ne 0 ]; then
  echo -e "${RED}✗ ListInvoiceDetails failed${NC}"
  echo "$DETAILS_RESP"
else
  PAGE_LEN=$(echo "$DETAILS_RESP" | jq -r '.detail_list.details | length')
  TOTAL_ELEM=$(echo "$DETAILS_RESP" | jq -r '.detail_list.total_elements')
  TOTAL_PAGE=$(echo "$DETAILS_RESP" | jq -r '.detail_list.total_page')
  if [ "$PAGE_LEN" -eq 2 ] && [ "$TOTAL_ELEM" -ge 5 ] && [ "$TOTAL_PAGE" -ge 3 ]; then
    echo -e "${GREEN}✓ ListInvoiceDetails pagination passed (len=$PAGE_LEN, total=$TOTAL_ELEM, total_page=$TOTAL_PAGE)${NC}"
  else
    echo -e "${RED}✗ ListInvoiceDetails unexpected pagination (len=$PAGE_LEN, total=$TOTAL_ELEM, total_page=$TOTAL_PAGE)${NC}"
    echo "$DETAILS_RESP"
  fi
fi
echo ""

echo -e "${YELLOW}=======================================${NC}"
echo -e "${GREEN}All integration tests completed!${NC}"
