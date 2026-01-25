#!/bin/bash
# HTTP Load Testing for API Gateway using Apache Bench (ab) or hey

set -euo pipefail

API_URL="${API_URL:-http://localhost:8080}"
API_ENDPOINT="/api/v1/health"

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}API Gateway Load Testing${NC}"
echo -e "${BLUE}========================================${NC}"

# Check if target is reachable
echo -e "\n${YELLOW}Checking API availability...${NC}"
if ! curl -s -f "${API_URL}${API_ENDPOINT}" > /dev/null; then
    echo -e "${RED}✗ API is not reachable at ${API_URL}${NC}"
    echo "Make sure api-gateway is running:"
    echo "  cd services/api-gateway && go run ./cmd/api/main.go"
    exit 1
fi
echo -e "${GREEN}✓ API is reachable${NC}"

# Test 1: Health Check (simple)
echo -e "\n${YELLOW}Test 1: Health Check Endpoint${NC}"
echo "URL: ${API_URL}${API_ENDPOINT}"
echo "Requests: 100, Concurrency: 10"
if command -v hey &> /dev/null; then
    hey -n 100 -c 10 "${API_URL}${API_ENDPOINT}" | tail -20
elif command -v ab &> /dev/null; then
    ab -n 100 -c 10 "${API_URL}${API_ENDPOINT}"
else
    echo -e "${YELLOW}Note: Install 'hey' or 'ab' for detailed load testing${NC}"
    # Fallback to curl
    for i in {1..10}; do
        curl -s "${API_URL}${API_ENDPOINT}" > /dev/null && echo -n "." || echo -n "X"
    done
    echo ""
fi

# Test 2: Circuit Breaker stress test
echo -e "\n${YELLOW}Test 2: Company List Endpoint (Protected)${NC}"
echo "Note: This test requires authentication. Generating JWT..."

# Generate JWT token
TOKEN=$(curl -s -X POST "${API_URL}/api/v1/users/login" \
    -H "Content-Type: application/json" \
    -d '{"email":"test@example.com","password":"password123"}' | grep -o '"token":"[^"]*' | cut -d'"' -f4) || TOKEN=""

if [ -z "$TOKEN" ]; then
    echo -e "${YELLOW}Could not generate token, testing without auth...${NC}"
    if command -v hey &> /dev/null; then
        hey -n 50 -c 5 "${API_URL}/api/v1/companies" 2>/dev/null | tail -10
    fi
else
    echo -e "${GREEN}✓ Token generated: ${TOKEN:0:20}...${NC}"
    if command -v hey &> /dev/null; then
        hey -n 50 -c 5 -H "Authorization: Bearer ${TOKEN}" "${API_URL}/api/v1/companies" | tail -20
    fi
fi

echo -e "\n${GREEN}========================================${NC}"
echo -e "${GREEN}Load Testing Complete${NC}"
echo -e "${GREEN}========================================${NC}"

# Print recommendations
echo -e "\n${BLUE}Recommendations:${NC}"
echo "1. Install 'hey' for detailed HTTP load testing:"
echo "   go install github.com/rakyll/hey@latest"
echo ""
echo "2. Monitor metrics during load test:"
echo "   curl http://localhost:9090/metrics | grep api_gateway"
echo ""
echo "3. Check Circuit Breaker metrics:"
echo "   curl http://localhost:9090/metrics | grep circuit_breaker"
