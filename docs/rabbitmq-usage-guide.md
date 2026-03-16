# RabbitMQ 消息队列使用指南

## 一、快速开始

### 1.1 配置文件

在 `config.yml` 中添加 RabbitMQ 配置：

```yaml
RabbitMQ:
  Host: "115.190.43.83"
  Port: 5672
  User: "guest"
  Password: "guest"
  Vhost: "/"
```

### 1.2 依赖安装

```bash
go get github.com/streadway/amqp
go get github.com/go-redis/redis/v8
```

### 1.3 初始化

在 `init.go` 中已自动初始化 RabbitMQ 和 Redis：

```go
func init() {
    ViperInit()
    MysqlInit()
    RedisInit()    // Redis初始化
    RabbitMQInit() // RabbitMQ初始化
}
```

## 二、API 说明

### 2.1 消息发送

#### 发送简单字符串消息

```go
import lianxiMQ "lianxi/srv/dasic/mq"

err := lianxiMQ.SendMsg("exchange", "routing.key", "message content")
if err != nil {
    log.Printf("发送失败: %v", err)
}
```

#### 发送 JSON 消息

```go
type OrderMessage struct {
    OrderNo   string  `json:"order_no"`
    ProductID int     `json:"product_id"`
    Quantity  int     `json:"quantity"`
}

msg := OrderMessage{
    OrderNo:   "20260313123456",
    ProductID: 1,
    Quantity:  2,
}

err := lianxiMQ.SendJSONMsg("order.exchange", "order.created", msg)
```

#### Direct 模式：直接发送到队列

```go
err := lianxiMQ.SendToQueue("queue.name", "message content")
```

### 2.2 消息订阅

#### Topic 模式：按路由键订阅

```go
err := lianxiMQ.SubscribeMsg("exchange", "order.*", func(msg string) error {
    log.Printf("收到消息: %s", msg)

    // 处理消息逻辑
    // ...

    return nil // 返回 nil 表示处理成功
})
if err != nil {
    log.Printf("订阅失败: %v", err)
}
```

#### Direct 模式：订阅队列

```go
err := lianxiMQ.SubscribeQueue("queue.name", func(msg string) error {
    log.Printf("收到消息: %s", msg)
    // 处理消息逻辑
    return nil
})
```

## 三、幂等性保证

系统自动使用 Redis SETNX 实现消息幂等性检查：

- **机制**: 每条消息处理前检查 Redis 中是否已存在该 MessageId
- **Key 格式**: `mq:message:idem:{MessageId}`
- **过期时间**: 24 小时
- **处理逻辑**:
  - 首次处理: SETNX 返回 true，标记已处理，执行业务逻辑
  - 重复消息: SETNX 返回 false，跳过处理，直接确认消息

```go
// 消息消费者示例
err := lianxiMQ.SubscribeMsg("order.exchange", "order.created", func(msg string) error {
    // 自动幂等性检查
    // 如果消息重复，会自动跳过，不会执行到这里

    // 业务处理逻辑
    // ...

    return nil
})
```

## 四、电商场景示例

### 4.1 订单创建流程

#### 生产者：订单服务

```go
type OrderCreatedEvent struct {
    OrderNo     string  `json:"order_no"`
    MemberID    int     `json:"member_id"`
    ProductID   int     `json:"product_id"`
    Quantity    int     `json:"quantity"`
    TotalAmount float64 `json:"total_amount"`
}

// 订单创建成功后发布事件
event := OrderCreatedEvent{
    OrderNo:     "ORD001",
    MemberID:    1001,
    ProductID:   1,
    Quantity:    2,
    TotalAmount: 299.00,
}

err := lianxiMQ.SendJSONMsg("ecommerce.exchange", "order.created", event)
```

#### 消费者：库存服务

```go
// 订阅订单创建事件，锁定库存
lianxiMQ.SubscribeMsg("ecommerce.exchange", "order.created", func(msg string) error {
    var event OrderCreatedEvent
    json.Unmarshal([]byte(msg), &event)

    // 锁定库存
    success := LockInventory(event.ProductID, event.Quantity)

    // 发布锁定结果
    lockEvent := InventoryLockedEvent{
        OrderNo: event.OrderNo,
        Success: success,
    }
    lianxiMQ.SendJSONMsg("ecommerce.exchange", "inventory.locked", lockEvent)

    return nil
})
```

### 4.2 订单支付流程

#### 生产者：订单服务

```go
type OrderPaidEvent struct {
    OrderNo   string `json:"order_no"`
    MemberID  int    `json:"member_id"`
    Address   string `json:"address"`
}

// 支付成功后发布事件
event := OrderPaidEvent{
    OrderNo:  "ORD001",
    MemberID: 1001,
    Address:  "北京市朝阳区",
}

err := lianxiMQ.SendJSONMsg("ecommerce.exchange", "order.paid", event)
```

#### 消费者：物流服务

```go
// 订阅订单支付事件，创建物流单
lianxiMQ.SubscribeMsg("ecommerce.exchange", "order.paid", func(msg string) error {
    var event OrderPaidEvent
    json.Unmarshal([]byte(msg), &event)

    // 创建物流单
    CreateLogistics(event.OrderNo, event.Address)

    return nil
})
```

### 4.3 订单完成流程

#### 生产者：订单服务

```go
type OrderCompletedEvent struct {
    OrderNo     string  `json:"order_no"`
    MemberID    int     `json:"member_id"`
    TotalAmount float64 `json:"total_amount"`
}

// 用户确认收货后发布事件
event := OrderCompletedEvent{
    OrderNo:     "ORD001",
    MemberID:    1001,
    TotalAmount: 299.00,
}

err := lianxiMQ.SendJSONMsg("ecommerce.exchange", "order.completed", event)
```

#### 消费者：积分服务

```go
// 订阅订单完成事件，发放积分
lianxiMQ.SubscribeMsg("ecommerce.exchange", "order.completed", func(msg string) error {
    var event OrderCompletedEvent
    json.Unmarshal([]byte(msg), &event)

    // 计算积分
    points := event.TotalAmount * 1.0 // 积分倍率 1.0

    // 发放积分
    GrantPoints(event.MemberID, int(points))

    // 检查会员升级
    CheckMemberUpgrade(event.MemberID)

    return nil
})
```

## 五、启动所有消费者

在 `main.go` 中启动所有消费者：

```go
package main

import (
    "lianxi/srv/handler/service"
)

func main() {
    // ... 其他初始化

    // 启动所有消息消费者
    service.StartAllConsumers()

    // ... 启动 gRPC 服务
}
```

## 六、路由键规范

### 6.1 订单相关

| 事件 | 路由键 | 说明 |
|------|--------|------|
| 订单创建 | `order.created` | 订单创建成功 |
| 订单支付 | `order.paid` | 订单支付成功 |
| 订单完成 | `order.completed` | 用户确认收货 |
| 订单取消 | `order.cancelled` | 订单已取消 |

### 6.2 库存相关

| 事件 | 路由键 | 说明 |
|------|--------|------|
| 库存锁定 | `inventory.locked` | 库存锁定完成 |
| 库存扣减 | `inventory.deducted` | 库存扣减成功 |
| 库存释放 | `inventory.released` | 库存释放成功 |

### 6.3 积分相关

| 事件 | 路由键 | 说明 |
|------|--------|------|
| 积分发放 | `points.granted` | 积分发放成功 |
| 积分消费 | `points.consumed` | 积分消费成功 |

## 七、错误处理

### 7.1 消息处理失败

当处理函数返回错误时：

```go
lianxiMQ.SubscribeMsg("exchange", "key", func(msg string) error {
    // 处理失败
    return fmt.Errorf("处理失败")
})
```

系统会自动：
1. 记录错误日志
2. 拒绝消息（Nack）
3. 重新入队（requeue=true）
4. 重试处理

### 7.2 连接异常

RabbitMQ 连接断开后，系统会自动重连（在应用重启时）。

## 八、性能优化

### 8.1 QoS 设置

已设置 prefetch count 为 100，限制每个消费者未确认消息数量：

```go
channel.Qos(100, 0, false)
```

### 8.2 消息持久化

所有消息默认开启持久化：

```go
amqp.Publishing{
    DeliveryMode: amqp.Persistent,
}
```

### 8.3 交换机持久化

所有交换机声明为持久化：

```go
channel.ExchangeDeclare(
    "exchange",
    "topic",
    true, // durable
    ...
)
```

## 九、监控与调试

### 9.1 查看日志

所有消息发送和消费都有日志输出：

```
2026/03/13 15:30:00 发送消息成功 - Topic: order.exchange, RoutingKey: order.created, Msg: {"order_no":"ORD001"}
2026/03/13 15:30:01 收到消息 - Topic: order.exchange, RoutingKey: order.created, Msg: {"order_no":"ORD001"}
2026/03/13 15:30:01 消息处理成功 - MessageId: xxx
```

### 9.2 Redis 查看已处理消息

```bash
redis-cli
> KEYS mq:message:idem:*
> TTL mq:message:idem:{message_id}
```

### 9.3 RabbitMQ 管理界面

访问 RabbitMQ 管理界面（http://host:15672）查看：
- 交换机状态
- 队列状态
- 消息堆积情况
- 消费者连接状态

## 十、常见问题

### Q1: 消息重复处理怎么办？

A: 系统已自动使用 Redis 实现幂等性，无需手动处理。

### Q2: 消息处理失败会重试吗？

A: 会，消息会重新入队，无限重试。建议在处理函数中实现重试次数限制。

### Q3: 如何保证消息不丢失？

A: 已开启：
- 消息持久化（DeliveryMode = Persistent）
- 交换机持久化（durable = true）
- 手动确认模式（auto-ack = false）

### Q4: 如何处理消息堆积？

A:
- 增加消费者实例
- 调整 QoS 参数
- 优化处理函数性能
- 考虑使用死信队列