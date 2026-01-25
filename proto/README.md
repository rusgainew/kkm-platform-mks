# Proto Files Organization

## 📁 Структура папки `api/proto/`

Все gRPC контракты и protobuf сообщения хранятся в одном месте.

```
api/proto/
├── auth_service.proto      # gRPC AuthService definition
├── auth.proto              # User, Token, Credential messages
├── company_service.proto   # gRPC CompanyService definition
├── company.proto           # Organization, Employee messages
├── document_service.proto  # gRPC DocumentService definition
├── document.proto          # Document, DocumentEntry messages
├── common.proto            # Shared messages: Error, Empty, Pagination
├── Makefile                # Build script for compilation
├── generate_protos.sh      # (optional) Helper script
└── README.md               # This file
```

## 🔧 Компиляция

### Установка инструментов

```bash
cd api/proto
make proto-install-tools
```

### Компиляция всех proto файлов (в отдельную библиотеку)

Генерация Go-кода теперь складывается в модуль `proto-lib/`, который можно импортировать из сервисов.

```bash
cd api/proto
make proto
```

### Очистка сгенерированных файлов

```bash
cd api/proto
make proto-clean
```

## 📋 Proto Files Description

### Service Definitions (синхронные RPC вызовы)

#### `auth_service.proto`

```protobuf
service AuthService {
  rpc Register(RegisterRequest) returns (RegisterResponse);
  rpc Login(LoginRequest) returns (LoginResponse);
  rpc ValidateToken(ValidateTokenRequest) returns (ValidateTokenResponse);
  rpc GetUser(GetUserRequest) returns (User);
  rpc Logout(LogoutRequest) returns (Empty);
}
```

**Используется:**

- auth-service (реализует сервис)
- company-service, document-service (клиенты для проверки JWT)

#### `company_service.proto`

```protobuf
service CompanyService {
  rpc GetOrganization(GetOrgRequest) returns (Organization);
  rpc CreateOrganization(CreateOrgRequest) returns (Organization);
  rpc UpdateOrganization(UpdateOrgRequest) returns (Organization);
  rpc ListOrganizations(ListOrgRequest) returns (ListOrgResponse);
}
```

**Используется:**

- company-service (реализует сервис)
- document-service, api-gateway (клиенты)

#### `document_service.proto`

```protobuf
service DocumentService {
  rpc GetDocument(GetDocRequest) returns (Document);
  rpc CreateDocument(CreateDocRequest) returns (Document);
  rpc SendDocument(SendDocRequest) returns (Document);
  rpc ApproveDocument(ApproveDocRequest) returns (Document);
  rpc ListDocuments(ListDocRequest) returns (ListDocResponse);
}
```

**Используется:**

- document-service (реализует сервис)
- api-gateway (клиент)

### Message Definitions (типы данных)

#### `auth.proto`

```protobuf
message User {
  string id = 1;
  string email = 2;
  string first_name = 3;
  string last_name = 4;
  int64 created_at = 5;
}

message Token {
  string access_token = 1;
  int64 expires_in = 2;
  string token_type = 3;  // "Bearer"
}

message Credential {
  string email = 1;
  string password = 2;
}
```

#### `company.proto`

```protobuf
message Organization {
  string id = 1;
  string name = 2;
  string description = 3;
  string owner_id = 4;
  int64 created_at = 5;
  int64 updated_at = 6;
}

message Employee {
  string id = 1;
  string organization_id = 2;
  string user_id = 3;
  string role = 4;
  int64 created_at = 5;
}
```

#### `document.proto`

```protobuf
message Document {
  string id = 1;
  string organization_id = 2;
  string title = 3;
  string content = 4;
  string status = 5;  // "draft", "sent", "approved"
  string created_by = 6;
  int64 created_at = 7;
  int64 updated_at = 8;
  repeated DocumentEntry entries = 9;
}

message DocumentEntry {
  string id = 1;
  string document_id = 2;
  string key = 3;
  string value = 4;
}
```

#### `common.proto`

```protobuf
message Error {
  int32 code = 1;
  string message = 2;
  string details = 3;
}

message Empty {}

message PageInfo {
  int32 page = 1;
  int32 per_page = 2;
  int64 total = 3;
}
```

## 🔄 Generated Go Code

После компиляции `make proto` генерируются файлы:

### Для service definitions:

- `auth_service.pb.go` - Структуры сообщений
- `auth_service_grpc.pb.go` - gRPC client и server interfaces

### Для message definitions:

- `auth.pb.go` - Структуры сообщений

## 📦 Integration with Microservices

### In each microservice (e.g., auth-service):

```go
// internal/interfaces/grpc/handlers/auth_handler.go

package handlers

import (
    pb "github.com/rusgainew/tunduck/api/proto"  // Generated proto code
    "context"
)

type AuthHandler struct {
    pb.UnimplementedAuthServiceServer
    service AuthService
}

func (h *AuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
    // Implementation
}
```

### In clients (e.g., company-service calling auth-service):

```go
// internal/infrastructure/grpc/client/auth_client.go

import (
    pb "github.com/rusgainew/tunduck/api/proto"
)

func (c *AuthClient) ValidateToken(ctx context.Context, token string) (*pb.User, error) {
    req := &pb.ValidateTokenRequest{Token: token}
    return c.client.ValidateToken(ctx, req)
}
```

## 🚀 Best Practices

1. **Versioning**: Всегда добавляйте комментарии в proto файлы

   ```protobuf
   // Added in v1.2.0
   string new_field = 10;
   ```

2. **Backward Compatibility**: Никогда не удаляйте номера полей (tags)

   ```protobuf
   // BAD: Удаление поля сломает клиентов
   // string old_field = 5;  // REMOVED

   // GOOD: Пометьте как deprecated
   string old_field = 5 [deprecated = true];
   ```

3. **Documentation**: Документируйте услуги и сообщения

   ```protobuf
   // User aggregate root
   // Contains basic user information
   message User {
     // Unique user identifier
     string id = 1;
   }
   ```

4. **Package Names**: Используйте понятные имена пакетов

   ```protobuf
   syntax = "proto3";
   package github.rusgainew.tunduck.auth;
   ```

5. **Go Package Names**: Явно указывайте go_package
   ```protobuf
   option go_package = "github.com/rusgainew/tunduck/gen/proto/auth";
   ```

## � Services and Dependencies

### All gRPC Services Inventory

Complete listing of all gRPC services, RPC methods, and their parameters:

- **[SERVICES_AND_DEPENDENCIES.md](./SERVICES_AND_DEPENDENCIES.md)** - Complete services list with RPC methods

### Proto Import Graph

Visual dependency graph showing how proto files import from each other:

![Proto Dependency Graph](./SERVICES_DEP_GRAPH.png)

_Generated from [SERVICES_DEP_GRAPH.dot](./SERVICES_DEP_GRAPH.dot) using Graphviz. To regenerate:_

```bash
dot -Tpng SERVICES_DEP_GRAPH.dot -o SERVICES_DEP_GRAPH.png
```

## �🔗 References

- [Protocol Buffers Documentation](https://developers.google.com/protocol-buffers)
- [gRPC Go Quickstart](https://grpc.io/docs/languages/go/quickstart/)
- [Protocol Buffers Go Generated Code](https://developers.google.com/protocol-buffers/docs/reference/go-generated)
