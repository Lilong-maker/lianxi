# AGENTS.md - Repository Guidelines for Agentic Coding

This document provides essential information for agentic coding agents working in this repository.

## Project Overview

**Technology Stack:**
- Language: Go 1.25.8
- RPC Framework: gRPC (google.golang.org/grpc)
- HTTP Framework: Gin (github.com/gin-gonic/gin) - BFF layer
- Serialization: Protocol Buffers (google.golang.org/protobuf)
- Database ORM: GORM (gorm.io/gorm) with MySQL driver
- Configuration: Viper (spf13/viper)
- Service Discovery: Consul (hashicorp/consul)
- Config Center: Nacos (nacos-group/nacos-sdk-go)
- Message Queue: RabbitMQ (streadway/amqp)
- Payment: Alipay (smartwalle/alipay)

**Architecture:**
```
lianxi/
├── srv/                          # gRPC Services Layer
│   ├── dasic/cmd/main.go         # gRPC server entry (port 50051)
│   ├── dasic/config/             # Global config (DB, AppConfig)
│   ├── dasic/inits/              # Init functions (MySQL, Redis, Consul, Nacos)
│   ├── handler/model/            # GORM data models
│   ├── handler/service/          # gRPC service implementations
│   └── proto/                    # Proto definitions & generated code
├── bff/                          # Backend for Frontend Layer
│   ├── dasic/cmd/main.go         # HTTP server entry (port 8888)
│   ├── router/router.go          # Gin router setup
│   ├── handler/                  # HTTP handlers
│   └── api/                      # API utilities
├── RabbitMQ/                     # RabbitMQ wrapper for async messaging
├── pkg/                          # Shared utilities (Alipay)
└── config.yml                    # Application configuration
```

## Build and Development Commands

**Build gRPC Server:**
```bash
go build -o bin/grpc-server ./srv/dasic/cmd
```

**Build BFF HTTP Server:**
```bash
go build -o bin/bff-server ./bff/dasic/cmd
```

**Run gRPC Server:**
```bash
go run ./srv/dasic/cmd/main.go
# or with custom port
./bin/grpc-server -port=50051
```

**Run BFF HTTP Server:**
```bash
go run ./bff/dasic/cmd/main.go
```

**Protocol Buffer Generation:**
```bash
# Generate Goods proto
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       srv/proto/goods/user.proto

# Generate Order proto
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       srv/proto/order/order.proto
```

**Dependencies:**
```bash
go mod tidy
```

## Testing

**Run all tests:**
```bash
go test ./...
```

**Run tests for a specific package:**
```bash
go test -v ./RabbitMQ/...
go test -v ./srv/handler/model/...
```

**Run a specific test function:**
```bash
go test -v -run TestSendMsg ./RabbitMQ/...
go test -v -run TestStockDeduct ./RabbitMQ/...
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
- **Proto imports:** Use numbered aliases (`__`, `__2`) for proto packages to avoid conflicts

```go
import (
    "context"
    "fmt"
    "lianxi/srv/dasic/config"
    "lianxi/srv/handler/model"
    __ "lianxi/srv/proto/goods"
    __2 "lianxi/srv/proto/order"
    "github.com/google/uuid"
    "gorm.io/gorm"
)
```

### Naming Conventions
- **Files:** snake_case (`goods.go`, `order.go`, `mq_test.go`)
- **Packages:** lowercase, single word preferred (`model`, `goods`, `order`)
- **Structs/Types:** PascalCase (`Goods`, `Order`, `AppConfig`, `RabbitMQ`)
- **Functions (exported):** PascalCase (`GoodsAdd`, `MysqlInit`, `Router`)
- **Variables:** camelCase (`orderNo`, `totalPrice`, `pageSize`)
- **Constants/Globals:** UPPERCASE or exported PascalCase (`DB`, `Gen`, `MQURL`)
- **Proto message suffix:** `Req` for requests, `Resp` for responses (`GoodsAddReq`, `GoodsAddResp`)

### Type Definitions
- **GORM Models:** Always embed `gorm.Model` for auto ID, timestamps, soft delete
- **Struct Fields:** Use GORM tags with Chinese comments for database schema

```go
type Order struct {
    gorm.Model
    OrderNo    string  `gorm:"type:varchar(32);"`
    UserID     int     `gorm:"not null;comment:用户ID"`
    TotalPrice float64 `gorm:"type:decimal(10,2);not null;comment:订单总金额"`
    PayStatus  int     `gorm:"default:0;comment:支付状态 0未支付 1已支付"`
}
```

- **Model Methods:** Attach CRUD methods to struct pointers, accept `*gorm.DB` as parameter

```go
func (o *Goods) GoodsAdd(db *gorm.DB) error {
    return db.Create(&o).Error
}

func GoodsList(db *gorm.DB, page, pageSize int) ([]Goods, int64, error) {
    // Package-level function for list queries
}
```

### Error Handling
- **Startup failures:** Use `panic()` or `log.Fatalf()` for initialization errors
- **Non-fatal errors:** Use `fmt.Printf()` for migration/non-critical failures
- **gRPC methods:** Return error message in response body, `nil` error for flow control

```go
// Startup error
panic(fmt.Sprintf("数据库连接失败: %v", err))
log.Fatalf("Consul初始化失败: %v", err)

// gRPC service pattern - error in response, not as Go error
if err != nil {
    return &__.GoodsAddResp{
        Msg:  "商品不存在",
        Code: 404,
    }, nil
}
```

### Comments
- **Model tables:** Chinese comments describing purpose: `// Order 订单表`
- **GORM tags:** Chinese field descriptions: `comment:商品名称`
- **Test functions:** Chinese comments explaining test purpose: `// 测试发送消息`

### gRPC Service Pattern
- Embed `UnimplementedXxxServer` for forward compatibility
- Use blank context `_` when context is unused
- Method signature: `func (s *Server) MethodName(_ context.Context, in *Request) (*Response, error)`

```go
type Server struct {
    __.UnimplementedGoodsServer
}

func (s *Server) GoodsAdd(_ context.Context, in *__.GoodsAddReq) (*__.GoodsAddResp, error) {
    // Implementation
    return &__.GoodsAddResp{Msg: "成功", Code: 200}, nil
}
```

### Gin HTTP Handler Pattern (BFF Layer)
- Handler functions accept `*gin.Context`
- Return JSON with `c.JSON()` or string with `c.String()`

```go
func GoodsAdd(c *gin.Context) {
    // Handler implementation
    c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success"})
}
```

### Init Pattern
- Use `init()` function for automatic initialization
- Call init functions explicitly for clarity

```go
func init() {
    ViperInit()
    MysqlInit()
    RedisInit()
    NacosInit()
}
```

## Important Notes

- **Config Security:** `config.yml` contains credentials - use environment variables in production
- **Ports:** gRPC on 50051, BFF HTTP on 8888
- **Proto Files:** Must regenerate Go code after modifying `.proto` files
- **Database Models:** Auto-migration runs on startup via `config.DB.AutoMigrate()`
- **Async Processing:** Stock deduction uses RabbitMQ with goroutines for non-blocking operations
- **Service Discovery:** Consul registration in `inits/consul.go`
- **Graceful Shutdown:** BFF server implements graceful shutdown with 5s timeout