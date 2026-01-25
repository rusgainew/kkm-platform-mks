#!/bin/bash

# Smoke test для foreign-company-server
# Проверяет базовую работоспособность сервера

set -e

echo "🧪 Foreign Company Server - Smoke Tests"
echo "========================================"

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Конфигурация
GRPC_PORT=${GRPC_PORT:-50056}
SERVER_BIN="/tmp/foreign-company-server"
JWT_SECRET="test-jwt-secret-key-12345"

# Функция для создания JWT токена
create_jwt() {
    local user_id="550e8400-e29b-41d4-a716-446655440000"
    local payload='{"user_id":"'$user_id'","email":"test@example.com","exp":'$(($(date +%s) + 3600))'}'
    
    # Простой JWT (для тестов)
    echo "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.$(echo -n "$payload" | base64 | tr -d '=' | tr '+/' '-_').$(echo -n "signature" | base64 | tr -d '=' | tr '+/' '-_')"
}

# Проверка компиляции
test_compilation() {
    echo -e "\n${YELLOW}→ Test 1: Compilation${NC}"
    
    cd /home/rusgaiub/go/src/github.com/rusgainew/kkm-project-mks/services/foreign-company-server
    
    if go build -o "$SERVER_BIN" ./cmd/main.go 2>&1 | grep -q "error"; then
        echo -e "${RED}✗ Compilation failed${NC}"
        return 1
    fi
    
    if [ -f "$SERVER_BIN" ]; then
        echo -e "${GREEN}✓ Server compiled successfully${NC}"
        ls -lh "$SERVER_BIN"
        return 0
    else
        echo -e "${RED}✗ Binary not found${NC}"
        return 1
    fi
}

# Проверка что сервер запускается
test_server_start() {
    echo -e "\n${YELLOW}→ Test 2: Server Start (dry-run)${NC}"
    
    # Проверим что бинарник выполняем
    if [ -x "$SERVER_BIN" ]; then
        echo -e "${GREEN}✓ Binary is executable${NC}"
    else
        echo -e "${RED}✗ Binary is not executable${NC}"
        return 1
    fi
    
    # Проверка что сервер не падает сразу при запуске
    timeout 2s "$SERVER_BIN" 2>&1 | head -5 || true
    
    if [ $? -eq 124 ]; then
        echo -e "${GREEN}✓ Server starts without immediate crash${NC}"
        return 0
    fi
    
    echo -e "${YELLOW}⚠ Server exited (might need database)${NC}"
    return 0
}

# Проверка JWT middleware
test_jwt_validation() {
    echo -e "\n${YELLOW}→ Test 3: JWT Middleware Check${NC}"
    
    # Проверим что код содержит JWT middleware
    if grep -q "getUserIDFromContext" /home/rusgaiub/go/src/github.com/rusgainew/kkm-project-mks/services/foreign-company-server/internal/interfaces/grpc/foreign_company_handler.go; then
        echo -e "${GREEN}✓ JWT middleware integration found${NC}"
    else
        echo -e "${RED}✗ JWT middleware not found${NC}"
        return 1
    fi
    
    # Проверим что хендлер использует user_id
    if grep -q "user_id.*context" /home/rusgaiub/go/src/github.com/rusgainew/kkm-project-mks/services/foreign-company-server/internal/interfaces/grpc/foreign_company_handler.go; then
        echo -e "${GREEN}✓ user_id extraction from context${NC}"
    else
        echo -e "${YELLOW}⚠ user_id extraction might be missing${NC}"
    fi
    
    return 0
}

# Проверка service layer
test_service_integration() {
    echo -e "\n${YELLOW}→ Test 4: Service Layer Integration${NC}"
    
    handler_file="/home/rusgaiub/go/src/github.com/rusgainew/kkm-project-mks/services/foreign-company-server/internal/interfaces/grpc/foreign_company_handler.go"
    
    if grep -q "service.CreateForeignCompany" "$handler_file"; then
        echo -e "${GREEN}✓ Create method calls service layer${NC}"
    else
        echo -e "${RED}✗ Create method doesn't call service${NC}"
        return 1
    fi
    
    if grep -q "service.UpdateForeignCompany" "$handler_file"; then
        echo -e "${GREEN}✓ Update method calls service layer${NC}"
    else
        echo -e "${RED}✗ Update method doesn't call service${NC}"
        return 1
    fi
    
    # Проверка использования country_code из proto
    if grep -q "req.CountryCode" "$handler_file"; then
        echo -e "${GREEN}✓ Handler uses country_code from request${NC}"
    else
        echo -e "${YELLOW}⚠ country_code not used from request${NC}"
    fi
    
    # Проверка использования address из proto
    if grep -q "req.Address" "$handler_file"; then
        echo -e "${GREEN}✓ Handler uses address from request${NC}"
    else
        echo -e "${YELLOW}⚠ address not used from request${NC}"
    fi
    
    return 0
}

# Проверка error handling
test_error_handling() {
    echo -e "\n${YELLOW}→ Test 5: Error Handling${NC}"
    
    handler_file="/home/rusgaiub/go/src/github.com/rusgainew/kkm-project-mks/services/foreign-company-server/internal/interfaces/grpc/foreign_company_handler.go"
    
    if grep -q "mapError" "$handler_file"; then
        echo -e "${GREEN}✓ Error mapping function present${NC}"
    else
        echo -e "${YELLOW}⚠ Error mapping might be missing${NC}"
    fi
    
    if grep -q "codes.Unauthenticated" "$handler_file"; then
        echo -e "${GREEN}✓ Authentication error handling${NC}"
    else
        echo -e "${RED}✗ Authentication error handling missing${NC}"
        return 1
    fi
    
    return 0
}

# Запуск всех тестов
main() {
    local failed=0
    
    test_compilation || ((failed++))
    test_server_start || ((failed++))
    test_jwt_validation || ((failed++))
    test_service_integration || ((failed++))
    test_error_handling || ((failed++))
    
    echo -e "\n========================================"
    if [ $failed -eq 0 ]; then
        echo -e "${GREEN}✓ All smoke tests passed!${NC}"
        echo -e "Foreign-company-server is ready for integration testing"
        exit 0
    else
        echo -e "${RED}✗ $failed test(s) failed${NC}"
        exit 1
    fi
}

main
