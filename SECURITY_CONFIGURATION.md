# 🔒 Security Configuration Guide

## Critical Security Requirements

### ⚠️ BEFORE DEPLOYMENT

All services now enforce **fail-fast** security validation. Services will **refuse to start** if required security configuration is missing.

## Required Environment Variables

### 1. JWT_SECRET (CRITICAL)

**Status:** ✅ Now REQUIRED (as of security fixes)

```bash
# Generate a secure JWT secret
openssl rand -base64 32

# Set in environment
export JWT_SECRET="<your-generated-secret>"
```

**Services requiring JWT_SECRET:**

- ✅ api-gateway
- ✅ catalog-query-server
- ✅ invoice-query-server
- ✅ foreign-company-query-server
- ✅ bank-account-query-server
- ✅ company-server (COMPANY_SERVER_JWT_SECRET or JWT_SECRET)
- ✅ document-server (DOCUMENT_SERVER_JWT_SECRET or JWT_SECRET)
- ✅ bank-account-server
- ✅ invoice-server
- ✅ foreign-company-server
- ✅ catalog-server

**What happens without it:**

```
FATAL: JWT_SECRET environment variable is required
Service will NOT start
```

## Quick Start - Development

### 1. Copy environment template

```bash
cp .env.example .env
```

### 2. Generate secure JWT secret

```bash
# Generate and set JWT secret
JWT_SECRET=$(openssl rand -base64 32)
echo "JWT_SECRET=$JWT_SECRET" >> .env
```

### 3. Configure database connections

Update `.env` with your database credentials:

```dotenv
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=<your-secure-password>
DB_NAME=kkm_user_db
```

### 4. Start services

```bash
# Using Docker Compose
docker-compose up -d

# Or individual services
cd services/api-gateway
go run cmd/main.go
```

## Production Deployment

### ✅ Security Checklist

- [ ] Generate strong JWT secret (minimum 32 bytes)
- [ ] Store secrets in secrets manager (AWS Secrets Manager, Vault, etc.)
- [ ] Never commit `.env` files to version control
- [ ] Use TLS/SSL for all connections
- [ ] Enable database SSL mode (`DB_SSL_MODE=require`)
- [ ] Set secure Redis password
- [ ] Configure RabbitMQ authentication
- [ ] Review all TODO comments in code
- [ ] Run security audit: `go run github.com/securego/gosec/v2/cmd/gosec@latest ./...`

### Environment Variable Hierarchy

Services check environment variables in this order:

1. **Service-specific**: `<SERVICE>_JWT_SECRET`
2. **Global**: `JWT_SECRET`
3. **None**: Service FAILS to start ❌

Example:

```bash
# Option 1: Service-specific (preferred)
export COMPANY_SERVER_JWT_SECRET="secret1"
export DOCUMENT_SERVER_JWT_SECRET="secret2"

# Option 2: Global (simpler for development)
export JWT_SECRET="shared-secret"
```

## Recent Security Fixes

### ✅ Fixed Issues (2026-01-27)

1. **SQL Injection** - Sort parameter validation (commit 7c8aeb5)
2. **Panic Crashes** - Graceful error handling (commit 4bde793)
3. **Goroutine Leaks** - Proper context cancellation (commit 86437f0)
4. **Context Propagation** - Lifecycle management (commit 7b10f65)
5. **Hardcoded Secrets** - Fail-fast validation (commit aeddbea)
6. **Connection Pool** - Database optimization (commit 238b8bc)
7. **Input Validation** - Comprehensive validation (commit bd4e8ed)

### Validation Added

**Catalog Service:**

- Max filter length: 255 characters
- Max search text: 500 characters
- Sort field whitelist
- Page size limit: 100 items

**Invoice Service:**

- Date format validation (YYYY-MM-DD)
- Amount range: 0 to 999,999,999,999.99
- Logical date range checks
- All catalog validations

## Monitoring & Alerts

### Prometheus Metrics

All services expose metrics on `/metrics` endpoint:

- Request rates and latencies
- Error rates by type
- Database connection pool status
- Cache hit/miss rates
- Custom business metrics

### Health Checks

```bash
# Liveness (is service running?)
curl http://localhost:<port>/health/live

# Readiness (is service ready to serve traffic?)
curl http://localhost:<port>/health/ready

# Overall health
curl http://localhost:<port>/health
```

## Troubleshooting

### Service won't start - JWT_SECRET missing

**Error:**

```
FATAL: JWT_SECRET environment variable is required
```

**Solution:**

```bash
export JWT_SECRET="$(openssl rand -base64 32)"
```

### Database connection failed

**Check:**

1. Database is running: `docker ps | grep postgres`
2. Credentials are correct in `.env`
3. Network connectivity: `telnet localhost 5432`

### Cache errors (Redis)

Redis is **optional**. Services will run without caching if Redis is unavailable.

**To disable Redis:**

```bash
# Remove or comment out in .env
# REDIS_HOST=
```

## Additional Resources

- [Full Code Analysis Report](./GO_CODEBASE_ANALYSIS_REPORT.md)
- [Implementation Summary](./GO_CODEBASE_ANALYSIS_SUMMARY.md)
- [Docker Setup Guide](./DOCKER_SETUP_COMPLETE.md)
- [Quick Start Guide](./QUICKSTART.md)

## Support

For security issues, please contact the security team immediately.
Do NOT create public GitHub issues for security vulnerabilities.
