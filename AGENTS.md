# AGENTS.md - Repository Guidelines for Agentic Coding

This document provides essential information for agentic coding agents working in this repository.

## Project Overview

**Technology Stack:**
- Language: Go 1.25.8
- RPC Framework: gRPC (google.golang.org/grpc)
- Serialization: Protocol Buffers (google.golang.org/protobuf)
- Database ORM: GORM (gorm.io/gorm) with MySQL driver
- Configuration: Viper (spf13/viper)
- Database: MySQL (115.190.43.83:3306)
- Cache: Redis (115.190.43.83:6379)

**Architecture:**
```
lianxi/
├── srv/dasic/cmd/     # Main entry point (gRPC server on port 50051)
├── srv/dasic/config/  # Configuration loading with Viper
├── srv/dasic/init/    # Database and Redis initialization
├── srv/handler/model/ # GORM data models (User, Product, Member, etc.)
├── srv/handler/service/ # gRPC service implementations
├── proto/             # Protocol Buffer definitions and generated code
└── config.yml         # Application configuration (DB credentials, etc.)
```

## Build and Development Commands

**Build:**
```bash
go build -o bin/server ./srv/dasic/cmd
```

**Run:**
```bash
go run ./srv/dasic/cmd/main.go
# or
./bin/server -port=50051
```

**Protocol Buffer Generation (if modifying .proto files):**
```bash
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       proto/user.proto
```

**Dependencies:**
```bash
go mod tidy
```

## Testing

**Note:** No test files currently exist in this repository. When adding tests, use standard Go conventions:

**Run all tests:**
```bash
go test ./...
```

**Run tests for a specific package:**
```bash
go test -v ./srv/handler/model/...
```

**Run a specific test function:**
```bash
go test -v -run TestFunctionName ./srv/handler/model/...
```

**Run tests with coverage:**
```bash
go test -cover ./...
```

## Linting and Formatting

**Format code:**
```bash
go fmt ./...
```

**Vet:**
```bash
go vet ./...
```

**Lint (requires golangci-lint):**
```bash
golangci-lint run
```

## Code Style Guidelines

### Import Statements
- **Order:** Standard library → internal packages → third-party
- **Pattern:** Absolute imports with `lianxi/` prefix for internal packages
- **Proto imports:** Use blank identifier alias: `__ "lianxi/proto"`

```go
import (
    "fmt"
    "lianxi/srv/dasic/config"
    "lianxi/srv/handler/model"
    "github.com/spf13/viper"
    "gorm.io/gorm"
)
```

### Naming Conventions
- **Files:** snake_case (`user.go`, `config.go`, `init.go`)
- **Structs/Types:** PascalCase (`User`, `Product`, `AppConfig`)
- **Functions (exported):** PascalCase (`OrderAdd`, `MysqlInit`)
- **Variables:** camelCase (`port`, `mysqlConfig`, `dsn`)
- **Constants/Globals:** UPPERCASE (`DB`, `Gen`)

### Type Definitions
- **GORM Models:** Always embed `gorm.Model` for auto ID, timestamps, soft delete
- **Struct Fields:** Use GORM tags with Chinese comments for descriptions

```go
type User struct {
    gorm.Model
    Name     string `gorm:"type:varchar(30);comment:用户名"`
    Password string `gorm:"type:varchar(32);comment:密码"`
}
```

### Error Handling
- **Fatal errors:** Use `panic()` or `log.Fatalf()` for startup failures (DB connection)
- **Non-fatal errors:** Use `fmt.Printf()` for migration failures
- **gRPC methods:** Return `nil` error on success

```go
panic(fmt.Sprintf("数据库连接失败: %v", err))
return &__.OrderAddResp{}, nil
```

### Comments
- **Model tables:** Chinese comments describing purpose: `// User 用户表`
- **GORM tags:** Chinese field descriptions: `comment:商品名称`

### gRPC Service Pattern
- Embed `UnimplementedOrderServer` for forward compatibility
- Method signature: `func (s *Server) MethodName(ctx context.Context, in *Request) (*Response, error)`

## Important Notes

- **Config Security:** `config.yml` contains MySQL/Redis credentials - consider environment variables
- **Port:** gRPC server runs on port 50051 by default
- **Proto Files:** Must regenerate Go code if `.proto` files are modified
- **Database Models:** Use GORM tags for schema definition and migrations