# Swagger/OpenAPI Documentation

## Overview

KKM API Gateway implements full Swagger/OpenAPI 2.0 documentation for all endpoints.

## Accessing Documentation

### Swagger UI (Interactive)

- **URL**: `http://localhost:8080/api/v1/docs/`
- **Alternative**: `http://localhost:8080/swagger/`
- **Access**: Available without authentication
- **Features**:
  - Try-it-out interface
  - Request/response examples
  - Parameter validation
  - Authentication with Bearer token

### Raw Specifications

- **JSON**: `/docs/swagger.json`
- **YAML**: `/docs/swagger.yaml`

## Documented Endpoints

### Authentication (Public)

- `POST /users/register` - User registration
- `POST /users/login` - User login

### Health Checks (Public)

- `GET /health` - Liveness probe
- `GET /ready` - Readiness probe with backend service checks

### Users (Protected)

- `GET /users/me` - Get current user info
- `GET /users/{id}` - Get user by ID
- `PUT /users/{id}` - Update user

### Companies (Protected)

- `POST /companies` - Create company
- `GET /companies` - List companies (paginated)
- `GET /companies/{id}` - Get company by ID
- `PUT /companies/{id}` - Update company
- `DELETE /companies/{id}` - Delete company

### Invoices (Protected)

- `POST /invoices` - Create invoice
- `PUT /invoices/{id}` - Update invoice
- `GET /invoices-query` - List invoices (paginated)
- `GET /invoices-query/search` - Search invoices

### Catalogs (Protected)

- `POST /catalogs` - Create catalog
- `PUT /catalogs/{id}` - Update catalog
- `GET /catalogs-query` - List catalogs (paginated)

### Bank Accounts (Protected)

- `POST /bank-accounts` - Create bank account
- `PUT /bank-accounts/{id}` - Update bank account
- `GET /bank-accounts-query` - List bank accounts (paginated)

### Foreign Companies (Protected)

- `POST /foreign-companies` - Create foreign company
- `PUT /foreign-companies/{id}` - Update foreign company

## Authentication

Protected endpoints require a Bearer token in the Authorization header:

```bash
curl -H "Authorization: Bearer <JWT_TOKEN>" http://localhost:8080/api/v1/companies
```

## Pagination

Query endpoints support pagination:

- `page` (default: 0) - 0-indexed page number
- `page_size` (default: 20) - Items per page (max: 100)

Example:

```bash
GET /api/v1/invoices-query?page=0&page_size=20
```

## Generation

Documentation is auto-generated from code comments using `swag`:

```bash
# Regenerate documentation after API changes
cd services/api-gateway
go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/api/main.go
```

### Code Comment Format

```go
// GetCompanies gets list of companies
// @Summary List companies
// @Description Get a paginated list of companies
// @Tags Companies
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Success 200 {object} models.APIResponse{data=[]models.Company}
// @Failure 500 {object} models.APIResponse
// @Router /companies [get]
// @Security BearerAuth
func (h *CompanyHandler) GetCompanies(c *gin.Context) { ... }
```

## Files

- `docs/swagger.json` - Generated OpenAPI 2.0 specification (JSON)
- `docs/swagger.yaml` - Generated OpenAPI 2.0 specification (YAML)
- `docs/docs.go` - Generated Go package containing embedded documentation
- `cmd/api/main.go` - Main package with API info annotations

## Next Steps

- Run locally: `docker-compose -f docker-compose.dev.yml up --build`
- Test endpoints via Swagger UI at `/api/v1/docs/`
- Load test with `ghz` tool
- Deploy to production
