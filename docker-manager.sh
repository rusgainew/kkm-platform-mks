#!/bin/bash

# 🚀 Quick Docker Compose Commands
# Быстрые команды для работы с KKM Project MKS

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_header() {
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}========================================${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

# Check if docker compose is available
if ! command -v docker &> /dev/null; then
    print_error "Docker is not installed or not in PATH"
    exit 1
fi

# Main menu
show_menu() {
    clear
    print_header "KKM Project MKS - Docker Compose Manager"
    echo ""
    echo "1)  🚀 Start all services"
    echo "2)  🛑 Stop all services"
    echo "3)  🔄 Restart all services"
    echo "4)  🔨 Build all images"
    echo "5)  📊 Show service status"
    echo "6)  📝 Show logs (all services)"
    echo "7)  🔍 Show logs (specific service)"
    echo "8)  💾 Start only infrastructure (DB, Redis, RabbitMQ)"
    echo "9)  📈 Start only monitoring (Prometheus, Grafana, Jaeger)"
    echo "10) 🏥 Health check"
    echo "11) 🧹 Clean up (stop & remove volumes)"
    echo "12) 📊 Show resource usage"
    echo "13) 🔧 Quick troubleshooting"
    echo "14) 🚪 Exit"
    echo ""
}

start_all() {
    print_header "Starting all services..."
    docker compose up -d
    print_success "All services started!"
    echo ""
    echo "Access points:"
    echo "  - API Gateway: http://localhost:8080"
    echo "  - Nginx Proxy: http://localhost"
    echo "  - Swagger UI: http://localhost/swagger/index.html"
    echo "  - Grafana: http://localhost:3000 (admin/admin)"
    echo "  - Prometheus: http://localhost:9091"
    echo "  - Jaeger: http://localhost:16686"
    echo "  - RabbitMQ: http://localhost:15672 (kkm_user/kkm_password)"
}

stop_all() {
    print_header "Stopping all services..."
    docker compose down
    print_success "All services stopped!"
}

restart_all() {
    print_header "Restarting all services..."
    docker compose restart
    print_success "All services restarted!"
}

build_all() {
    print_header "Building all images..."
    print_warning "This may take several minutes..."
    docker compose build
    print_success "All images built!"
}

show_status() {
    print_header "Service Status"
    docker compose ps
}

show_logs() {
    print_header "Showing logs (Ctrl+C to exit)..."
    docker compose logs -f --tail=100
}

show_logs_service() {
    echo "Available services:"
    docker compose ps --services
    echo ""
    read -p "Enter service name: " service
    print_header "Showing logs for $service (Ctrl+C to exit)..."
    docker compose logs -f --tail=100 "$service"
}

start_infra() {
    print_header "Starting infrastructure services..."
    docker compose up -d postgres redis rabbitmq
    print_success "Infrastructure started!"
    echo ""
    echo "Services running:"
    echo "  - PostgreSQL: localhost:5432"
    echo "  - Redis: localhost:6379"
    echo "  - RabbitMQ: localhost:5672, Management: http://localhost:15672"
}

start_monitoring() {
    print_header "Starting monitoring services..."
    docker compose up -d prometheus grafana jaeger
    print_success "Monitoring started!"
    echo ""
    echo "Access points:"
    echo "  - Prometheus: http://localhost:9091"
    echo "  - Grafana: http://localhost:3000 (admin/admin)"
    echo "  - Jaeger: http://localhost:16686"
}

health_check() {
    print_header "Health Check"
    
    echo -n "Testing API Gateway... "
    if curl -sf http://localhost:8080/api/v1/health > /dev/null 2>&1; then
        print_success "OK"
    else
        print_error "FAILED"
    fi
    
    echo -n "Testing Nginx Proxy... "
    if curl -sf http://localhost/health > /dev/null 2>&1; then
        print_success "OK"
    else
        print_error "FAILED"
    fi
    
    echo -n "Testing PostgreSQL... "
    if docker compose exec -T postgres pg_isready -U kkm_user > /dev/null 2>&1; then
        print_success "OK"
    else
        print_error "FAILED"
    fi
    
    echo -n "Testing Redis... "
    if docker compose exec -T redis redis-cli ping > /dev/null 2>&1; then
        print_success "OK"
    else
        print_error "FAILED"
    fi
    
    echo -n "Testing Prometheus... "
    if curl -sf http://localhost:9091/-/healthy > /dev/null 2>&1; then
        print_success "OK"
    else
        print_error "FAILED"
    fi
    
    echo -n "Testing Grafana... "
    if curl -sf http://localhost:3000/api/health > /dev/null 2>&1; then
        print_success "OK"
    else
        print_error "FAILED"
    fi
}

cleanup() {
    print_warning "This will stop all services and remove volumes!"
    read -p "Are you sure? (yes/no): " confirm
    if [ "$confirm" = "yes" ]; then
        print_header "Cleaning up..."
        docker compose down -v
        print_success "Cleanup complete!"
    else
        print_warning "Cleanup cancelled"
    fi
}

show_resources() {
    print_header "Resource Usage"
    docker stats --no-stream --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.NetIO}}"
}

troubleshoot() {
    print_header "Quick Troubleshooting"
    echo ""
    echo "1. Check service status:"
    docker compose ps
    echo ""
    echo "2. Recent errors in logs:"
    docker compose logs --tail=20 | grep -i "error\|fatal\|panic" || echo "No recent errors found"
    echo ""
    echo "3. Resource usage:"
    docker stats --no-stream --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}"
    echo ""
    echo "4. Network connectivity:"
    docker network inspect kkm-project-mks_kkm-network | grep -A 5 "Containers" || echo "Network not found"
}

# Main loop
while true; do
    show_menu
    read -p "Choose option [1-14]: " choice
    echo ""
    
    case $choice in
        1) start_all ;;
        2) stop_all ;;
        3) restart_all ;;
        4) build_all ;;
        5) show_status ;;
        6) show_logs ;;
        7) show_logs_service ;;
        8) start_infra ;;
        9) start_monitoring ;;
        10) health_check ;;
        11) cleanup ;;
        12) show_resources ;;
        13) troubleshoot ;;
        14) 
            print_success "Goodbye!"
            exit 0
            ;;
        *)
            print_error "Invalid option"
            ;;
    esac
    
    echo ""
    read -p "Press Enter to continue..."
done
