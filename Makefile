.PHONY: help build run stop clean dev-up dev-down prod-up prod-down logs test docker-build docker-push docker-up docker-down docker-logs docker-ps docker-restart docker-clean

# Variables
IMAGE_NAME := user-server
IMAGE_TAG := latest
DOCKER_REGISTRY := # Set your registry here, e.g., ghcr.io/username

help: ## Show this help
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

# ===================================
# Docker Compose Commands
# ===================================

docker-up: ## Start all services with docker-compose
	docker compose up -d

docker-down: ## Stop all services
	docker compose down

docker-down-volumes: ## Stop all services and remove volumes
	docker compose down -v

docker-build-all: ## Build all Docker images
	docker compose build

docker-build-no-cache: ## Build all Docker images without cache
	docker compose build --no-cache

docker-ps: ## Show running containers status
	docker compose ps

docker-logs: ## Show logs from all services
	docker compose logs -f

docker-logs-api: ## Show API Gateway logs
	docker compose logs -f api-gateway

docker-logs-nginx: ## Show Nginx logs
	docker compose logs -f nginx-proxy

docker-logs-postgres: ## Show PostgreSQL logs
	docker compose logs -f postgres

docker-logs-redis: ## Show Redis logs
	docker compose logs -f redis

docker-logs-rabbitmq: ## Show RabbitMQ logs
	docker compose logs -f rabbitmq

docker-restart: ## Restart all services
	docker compose restart

docker-restart-service: ## Restart specific service (usage: make docker-restart-service SERVICE=api-gateway)
	docker compose restart $(SERVICE)

docker-clean: ## Remove all stopped containers, networks, and volumes
	docker compose down -v
	docker system prune -f

docker-health: ## Check health of all services
	@echo "=== Service Health Status ==="
	@docker compose ps --format "table {{.Service}}\t{{.Status}}\t{{.State}}"

docker-stats: ## Show resource usage statistics
	docker stats --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.NetIO}}"

# ===================================
# Production Compose
# ===================================

prod-up: ## Start production stack
	docker compose -f docker-compose.prod.yml up -d

prod-down: ## Stop production stack
	docker compose -f docker-compose.prod.yml down

prod-logs: ## Show production logs
	docker compose -f docker-compose.prod.yml logs -f

prod-build: ## Build production images
	docker compose -f docker-compose.prod.yml build

prod-restart: ## Restart production stack
	docker compose -f docker-compose.prod.yml restart

# ===================================
# Infrastructure Only
# ===================================

infra-up: ## Start only infrastructure services (postgres, redis, rabbitmq)
	docker compose up -d postgres redis rabbitmq

infra-down: ## Stop infrastructure services
	docker compose stop postgres redis rabbitmq

# ===================================
# Monitoring
# ===================================

monitoring-up: ## Start monitoring services (prometheus, grafana, jaeger)
	docker compose up -d prometheus grafana jaeger

monitoring-down: ## Stop monitoring services
	docker compose stop prometheus grafana jaeger

# ===================================
# Individual Services
# ===================================

start-user-server: ## Start user-server
	docker compose up -d user-server

start-catalog-server: ## Start catalog-server
	docker compose up -d catalog-server

start-invoice-server: ## Start invoice-server
	docker compose up -d invoice-server

start-api-gateway: ## Start api-gateway
	docker compose up -d api-gateway

# ===================================
# Local Build & Run
# ===================================

build: ## Build the Go binary locally
	go build -v -o bin/user-server ./cmd/main.go

run: ## Run the server locally
	go run ./cmd/main.go

test: ## Run tests
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean: ## Clean build artifacts
	rm -f bin/user-server main coverage.out coverage.html

# Docker targets
docker-build: ## Build Docker image
	cd .. && docker build -f user-server/Dockerfile -t $(IMAGE_NAME):$(IMAGE_TAG) .

docker-run: ## Run Docker container standalone
	docker run -d --name $(IMAGE_NAME) \
		-p 50051:50051 \
		-e USER_SERVER_DB_DRIVER=memory \
		$(IMAGE_NAME):$(IMAGE_TAG)

docker-stop: ## Stop Docker container
	docker stop $(IMAGE_NAME) || true
	docker rm $(IMAGE_NAME) || true

docker-push: docker-build ## Push Docker image to registry
	@if [ -z "$(DOCKER_REGISTRY)" ]; then \
		echo "Error: DOCKER_REGISTRY is not set"; \
		exit 1; \
	fi
	docker tag $(IMAGE_NAME):$(IMAGE_TAG) $(DOCKER_REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)
	docker push $(DOCKER_REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)

# Development compose
dev-up: ## Start development infrastructure (postgres, redis, rabbitmq)
	cd .. && docker-compose -f docker-compose.dev.yml up -d
	@echo ""
	@echo "Development infrastructure started!"
	@echo "PostgreSQL: localhost:5432 (postgres/postgres)"
	@echo "Redis: localhost:6379"
	@echo "RabbitMQ: localhost:5672 (admin/admin123)"
	@echo "RabbitMQ Management: http://localhost:15672"
	@echo ""
	@echo "Now run the server with: make run"

dev-down: ## Stop development infrastructure
	cd .. && docker-compose -f docker-compose.dev.yml down

dev-clean: ## Stop and remove development volumes
	cd .. && docker-compose -f docker-compose.dev.yml down -v

# Production compose
prod-up: ## Start full production stack (all services)
	docker-compose up -d --build

prod-down: ## Stop production stack
	docker-compose down

prod-clean: ## Stop and remove production volumes
	docker-compose down -v

prod-restart: ## Restart production stack
	docker-compose restart

prod-restart-server: ## Restart specific server
	@read -p "Enter server name (user-server/company-server/catalog-server/bank-account-server): " server; \
	docker-compose restart $$server

# Individual services
start-catalog: ## Start catalog-server
	docker-compose up -d catalog-server

start-bank-account: ## Start bank-account-server
	docker-compose up -d bank-account-server

start-foreign-company: ## Start foreign-company-server
	docker-compose up -d foreign-company-server

start-api-gateway: ## Start api-gateway
	docker-compose up -d api-gateway

restart-catalog: ## Restart catalog-server
	docker-compose restart catalog-server

restart-bank-account: ## Restart bank-account-server
	docker-compose restart bank-account-server

restart-foreign-company: ## Restart foreign-company-server
	docker-compose restart foreign-company-server

restart-api-gateway: ## Restart api-gateway
	docker-compose restart api-gateway

restart-invoice-query: ## Restart invoice-query-server
	docker-compose restart invoice-query-server

restart-bank-account-query: ## Restart bank-account-query-server
	docker-compose restart bank-account-query-server

restart-catalog-query: ## Restart catalog-query-server
	docker-compose restart catalog-query-server

build-catalog: ## Build catalog-server
	docker-compose build catalog-server

build-bank-account: ## Build bank-account-server
	docker-compose build bank-account-server

build-foreign-company: ## Build foreign-company-server
	docker-compose build foreign-company-server

build-api-gateway: ## Build api-gateway
	docker-compose build api-gateway

build-invoice-query: ## Build invoice-query-server
	docker-compose build invoice-query-server

build-bank-account-query: ## Build bank-account-query-server
	docker-compose build bank-account-query-server

build-catalog-query: ## Build catalog-query-server
	docker-compose build catalog-query-server

# Logs
logs: ## Show logs from all services
	docker-compose logs -f

logs-server: ## Show user-server logs only
	docker-compose logs -f user-server

logs-catalog: ## Show catalog-server logs
	docker-compose logs -f catalog-server

logs-bank-account: ## Show bank-account-server logs
	docker-compose logs -f bank-account-server

logs-api-gateway: ## Show api-gateway logs
	docker-compose logs -f api-gateway

logs-foreign-company: ## Show foreign-company-server logs
	docker-compose logs -f foreign-company-server

logs-invoice-query: ## Show invoice-query-server logs
	docker-compose logs -f invoice-query-server

logs-bank-account-query: ## Show bank-account-query-server logs
	docker-compose logs -f bank-account-query-server

logs-catalog-query: ## Show catalog-query-server logs
	docker-compose logs -f catalog-query-server

logs-postgres: ## Show postgres logs
	docker-compose logs -f postgres

logs-redis: ## Show redis logs
	docker-compose logs -f redis

logs-rabbitmq: ## Show rabbitmq logs
	docker-compose logs -f rabbitmq

# Status
status: ## Show status of all services
	docker-compose ps

health-check: ## Check health of all services
	@echo "=== User Server ==="
	@grpcurl -plaintext localhost:50051 grpc.health.v1.Health/Check || echo "❌ Not responding"
	@echo "\n=== Company Server ==="
	@grpcurl -plaintext localhost:50052 grpc.health.v1.Health/Check || echo "❌ Not responding"
	@echo "\n=== Invoice Query Server ==="
	@curl -sf http://localhost:8053/health || echo "❌ Not responding"
	@echo "\n=== Bank Account Query Server ==="
	@curl -sf http://localhost:8054/health || echo "❌ Not responding"
	@echo "\n=== Catalog Query Server ==="
	@curl -sf http://localhost:8055/health || echo "❌ Not responding"
	@echo "\n=== Foreign Company Server ==="
	@grpcurl -plaintext localhost:50056 grpc.health.v1.Health/Check || echo "❌ Not responding"
	@echo "\n=== API Gateway ==="
	@curl -sf http://localhost:8080/health || echo "❌ Not responding"

# Database
db-shell: ## Connect to postgres database
	docker-compose exec postgres psql -U postgres -d userdb

db-shell-catalog: ## Connect to catalog database
	docker-compose exec postgres psql -U postgres -d catalogdb

db-shell-bank: ## Connect to bank account database
	docker-compose exec postgres psql -U postgres -d bankaccountdb

db-shell-foreign: ## Connect to foreign company database
	docker-compose exec postgres psql -U postgres -d foreigncompanydb

# Redis
redis-cli: ## Connect to redis CLI
	cd .. && docker-compose exec redis redis-cli

# RabbitMQ
rabbitmq-status: ## Show RabbitMQ status
	cd .. && docker-compose exec rabbitmq rabbitmqctl status

# Integration test with docker
test-integration: dev-up ## Run integration tests with docker infrastructure
	@echo "Waiting for services to be ready..."
	@sleep 5
	go test -v -tags=integration ./internal/interfaces/grpc/
	$(MAKE) dev-down
