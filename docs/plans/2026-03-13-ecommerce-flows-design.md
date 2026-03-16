# 电商全流程设计文档

## 一、订单全流程时序图

```mermaid
sequenceDiagram
    participant Client as 客户端
    participant OrderService as 订单服务
    participant InventoryService as 库存服务
    participant PaymentService as 支付服务
    participant LogisticsService as 物流服务
    participant PointsService as 积分服务
    participant DB as 数据库

    Note over Client,DB: 1. 订单创建流程
    Client->>OrderService: 创建订单(商品ID, 数量)
    OrderService->>DB: 开启事务
    OrderService->>InventoryService: 检查库存
    InventoryService-->>OrderService: 返回可用库存
    alt 库存不足
        OrderService-->>Client: 返回库存不足错误
    else 库存充足
        OrderService->>DB: 创建订单记录(状态:待支付)
        OrderService->>InventoryService: 锁定库存
        InventoryService->>DB: 扣减可用库存<br/>增加锁定库存
        DB-->>OrderService: 提交事务
        OrderService-->>Client: 返回订单号
    end

    Note over Client,DB: 2. 订单支付流程
    Client->>PaymentService: 发起支付(订单号)
    PaymentService->>PaymentService: 处理支付
    alt 支付失败
        PaymentService-->>Client: 支付失败
        Client->>OrderService: 取消订单
        OrderService->>DB: 更新订单状态(已取消)
        OrderService->>InventoryService: 释放锁定库存
        InventoryService->>DB: 锁定库存→可用库存
    else 支付成功
        PaymentService-->>Client: 支付成功
        PaymentService->>OrderService: 支付回调
        OrderService->>DB: 更新订单状态(待发货)
    end

    Note over Client,DB: 3. 订单发货流程
    Client->>LogisticsService: 查看订单详情
    LogisticsService->>OrderService: 查询订单状态
    OrderService-->>LogisticsService: 返回订单信息
    OrderService->>DB: 更新订单状态(待签收)
    OrderService->>LogisticsService: 创建物流单
    LogisticsService->>DB: 创建物流记录(状态:已发货)
    LogisticsService-->>Client: 返回物流单号

    Note over Client,DB: 4. 订单签收流程
    Client->>LogisticsService: 确认收货
    LogisticsService->>DB: 更新物流状态(已签收)
    LogisticsService->>OrderService: 签收通知
    OrderService->>DB: 更新订单状态(已完成)
    OrderService->>PointsService: 触发积分发放
```

## 二、库存全流程时序图

```mermaid
sequenceDiagram
    participant Client as 客户端
    participant OrderService as 订单服务
    participant InventoryService as 库存服务
    participant DB as 数据库
    participant Cron as 定时任务

    Note over Client,DB: 1. 库存检查流程
    Client->>OrderService: 创建订单请求
    OrderService->>InventoryService: 检查库存(商品ID, 数量)
    InventoryService->>DB: 查询商品库存
    DB-->>InventoryService: 返回库存信息
    alt 可用库存 >= 请求数量
        InventoryService-->>OrderService: 返回充足
    else 可用库存 < 请求数量
        InventoryService-->>OrderService: 返回不足
        OrderService-->>Client: 返回库存不足
    end

    Note over Client,DB: 2. 库存锁定流程
    OrderService->>DB: 开启事务
    OrderService->>InventoryService: 锁定库存(商品ID, 数量)
    InventoryService->>DB: 更新可用库存(-数量)
    InventoryService->>DB: 更新锁定库存(+数量)
    DB-->>InventoryService: 锁定成功
    InventoryService-->>OrderService: 锁定成功
    OrderService->>DB: 创建订单记录
    DB-->>OrderService: 提交事务

    Note over Client,DB: 3. 库存扣减流程(支付成功)
    OrderService->>InventoryService: 扣减库存(订单号)
    InventoryService->>DB: 开启事务
    InventoryService->>DB: 更新锁定库存(-数量)
    InventoryService->>DB: 更新商品销量(+数量)
    DB-->>InventoryService: 扣减成功
    InventoryService-->>OrderService: 扣减成功

    Note over Client,DB: 4. 库存释放流程(订单取消/超时)
    OrderService->>InventoryService: 释放库存(订单号)
    InventoryService->>DB: 开启事务
    InventoryService->>DB: 更新锁定库存(-数量)
    InventoryService->>DB: 更新可用库存(+数量)
    DB-->>InventoryService: 释放成功
    InventoryService-->>OrderService: 释放成功

    Note over Client,DB: 5. 超时订单库存自动释放
    Cron->>OrderService: 查询超时未支付订单
    OrderService->>DB: 查询订单列表
    DB-->>OrderService: 返回订单列表
    loop 遍历超时订单
        OrderService->>DB: 更新订单状态(已取消)
        OrderService->>InventoryService: 释放锁定库存
        InventoryService->>DB: 锁定库存→可用库存
    end
```

## 三、积分全流程时序图

```mermaid
sequenceDiagram
    participant Client as 客户端
    participant OrderService as 订单服务
    member MemberService as 会员服务
    participant PointsService as 积分服务
    participant DB as 数据库

    Note over Client,DB: 1. 订单完成触发积分流程
    Client->>OrderService: 确认收货
    OrderService->>DB: 更新订单状态(已完成)
    OrderService->>PointsService: 计算并发放积分(订单号)

    PointsService->>DB: 查询订单信息
    DB-->>PointsService: 返回订单详情
    PointsService->>DB: 查询会员信息
    DB-->>PointsService: 返回会员等级

    PointsService->>PointsService: 计算积分<br/>积分 = 订单金额 * 会员等级积分倍率

    PointsService->>DB: 开启事务
    PointsService->>DB: 插入积分记录<br/>(类型:获得, 余额:当前+新增)
    PointsService->>MemberService: 更新会员总积分
    MemberService->>DB: 更新会员积分余额
    DB-->>PointsService: 提交事务

    PointsService->>PointsService: 检查会员升级
    alt 积分达到下一等级
        PointsService->>DB: 查询下一等级配置
        DB-->>PointsService: 返回等级信息
        PointsService->>MemberService: 升级会员等级
        MemberService->>DB: 更新会员等级ID
    end

    PointsService-->>Client: 返回积分发放结果
```

## 四、订单数据流程图

```mermaid
flowchart TD
    A[客户端发起创建订单] --> B{商品库存检查}
    B -->|不足| C[返回库存不足]
    B -->|充足| D[开启数据库事务]

    D --> E[创建订单记录<br/>状态:待支付<br/>订单号, 商品ID, 数量, 金额]
    E --> F[锁定库存<br/>可用库存 - 数量<br/>锁定库存 + 数量]
    F --> G{事务提交}
    G -->|失败| H[回滚所有操作]
    G -->|成功| I[返回订单号]

    I --> J[用户发起支付]
    J --> K{支付结果}
    K -->|失败| L[更新订单状态:已取消]
    L --> M[释放锁定库存<br/>锁定库存 → 可用库存]
    K -->|成功| N[更新订单状态:待发货]

    N --> O[商家发货]
    O --> P[创建物流单<br/>物流单号, 物流公司<br/>收发货人信息]
    P --> Q[更新订单状态:待签收]

    Q --> R[用户确认收货]
    R --> S[更新物流状态:已签收]
    S --> T[更新订单状态:已完成]

    T --> U[触发积分发放<br/>计算订单金额 × 积分倍率]
    U --> V[创建积分记录<br/>类型:获得, 关联订单号]
    V --> W[更新会员积分余额]
    W --> X{检查会员升级}
    X -->|达到下一等级| Y[升级会员等级]
    X -->|未升级| Z[流程结束]
    Y --> Z

    style A fill:#e1f5fe
    style C fill:#ffebee
    style I fill:#e8f5e9
    style L fill:#ffebee
    style N fill:#e8f5e9
    style T fill:#e8f5e9
    style Y fill:#fff3e0
```

## 五、库存数据流程图

```mermaid
flowchart TD
    A[订单创建请求] --> B[检查商品可用库存]
    B --> C{可用库存 >= 请求数量?}
    C -->|否| D[返回库存不足]
    C -->|是| E[开启事务]

    E --> F[更新可用库存<br/>可用库存 = 可用库存 - 请求数量]
    F --> G[更新锁定库存<br/>锁定库存 = 锁定库存 + 请求数量]
    G --> H{事务提交成功?}
    H -->|否| I[回滚库存操作]
    H -->|是| J[创建订单成功]

    J --> K[用户支付订单]
    K --> L{支付结果}
    L -->|支付失败| M[订单取消]
    M --> N[释放锁定库存<br/>锁定库存 = 锁定库存 - 请求数量<br/>可用库存 = 可用库存 + 请求数量]
    L -->|支付成功| O[扣减锁定库存<br/>锁定库存 = 锁定库存 - 请求数量<br/>商品销量 = 商品销量 + 请求数量]

    N --> P[订单已取消]
    O --> Q[订单已支付]

    R[定时任务<br/>每分钟执行] --> S[查询超时未支付订单<br/>创建时间 > 30分钟<br/>状态 = 待支付]
    S --> T{是否有超时订单?}
    T -->|否| R
    T -->|是| U[遍历超时订单]
    U --> V[更新订单状态:已取消]
    V --> W[释放锁定库存<br/>锁定库存 → 可用库存]
    W --> R

    style A fill:#e1f5fe
    style D fill:#ffebee
    style J fill:#e8f5e9
    style M fill:#ffebee
    style N fill:#fff3e0
    style O fill:#e8f5e9
    style V fill:#fff3e0
```

## 六、积分数据流程图

```mermaid
flowchart TD
    A[订单完成<br/>用户确认收货] --> B[查询订单信息<br/>订单号, 订单金额, 会员ID]
    B --> C[查询会员等级信息<br/>等级ID, 积分倍率]

    C --> D[计算应得积分<br/>积分 = 订单金额 × 会员等级积分倍率]
    D --> E[开启事务]

    E --> F[查询当前会员积分余额]
    F --> G[计算新积分余额<br/>新余额 = 当前余额 + 应得积分]
    G --> H[创建积分记录<br/>会员ID, 积分数量, 新余额<br/>类型:获得, 关联订单号]

    H --> I[更新会员积分余额<br/>总积分 = 总积分 + 应得积分<br/>可用积分 = 可用积分 + 应得积分]
    I --> J{事务提交成功?}
    J -->|否| K[回滚所有操作]
    J -->|是| L[查询会员等级配置]

    L --> M{当前积分是否达到下一等级?}
    M -->|否| N[积分发放完成]
    M -->|是| O[升级会员等级<br/>会员等级ID = 下一等级ID]

    O --> P[流程结束]
    N --> P

    Q[积分消费流程] --> R[用户使用积分抵扣]
    R --> S[计算抵扣金额<br/>抵扣金额 = 使用积分 / 积分兑换比例]
    S --> T[检查可用积分余额]

    T --> U{可用积分 >= 使用积分?}
    U -->|否| V[返回积分不足]
    U -->|是| W[开启事务]

    W --> X[创建积分记录<br/>类型:消费, 余额:当前-消费]
    X --> Y[更新会员积分<br/>可用积分 = 可用积分 - 使用积分]
    Y --> Z{事务提交成功?}
    Z -->|否| AA[回滚所有操作]
    Z -->|是| AB[订单金额减去抵扣金额]

    AA --> AC[订单创建失败]
    AB --> AD[订单创建成功]

    style A fill:#e1f5fe
    style K fill:#ffebee
    style O fill:#fff3e0
    style V fill:#ffebee
    style AA fill:#ffebee
    style AB fill:#e8f5e9
```

## 七、数据库表关系说明

### 核心表关联关系

```mermaid
erDiagram
    Member ||--o{ Points : "拥有"
    Member ||--|| MemberLevel : "属于"
    Member ||--o{ Order : "创建"
    Product ||--o{ Inventory : "拥有"
    Product ||--o{ OrderItem : "关联"
    Product ||--o{ Points : "关联消费"
    Order ||--|| Logistics : "配送"
    Order ||--o{ OrderItem : "包含"
    Inventory ||--o{ Order : "锁定"

    Member {
        int ID PK
        string MemberNo UK
        string Username UK
        string Mobile UK
        int MemberLevelID FK
        int TotalPoints
        int AvailablePoints
    }

    Product {
        int ID PK
        string ProductNo UK
        string ProductName
        float SalePrice
        int Stock
        int Sales
    }

    Order {
        int ID PK
        string OrderNo UK
        int MemberID FK
        string OrderStatus
        float TotalAmount
        int Status
    }

    Points {
        int ID PK
        int MemberID FK
        int Points
        int Balance
        int Type
        string RelatedType
        int RelatedID
    }

    Inventory {
        int ID PK
        int ProductID FK
        string SKU UK
        int Stock
        int LockedStock
        int AvailableStock
    }

    Logistics {
        int ID PK
        string OrderNo FK
        string LogisticsNo UK
        string LogisticsCompany
        int Status
    }

    MemberLevel {
        int ID PK
        string LevelName UK
        int LevelNo UK
        int RequiredPoints
        float DiscountRate
        float PointsRate
    }
```

## 八、关键业务规则

### 8.1 订单业务规则

1. **订单号生成规则**: 采用雪花算法，保证全局唯一
2. **订单状态流转**:
   - 待支付 (0) → 待发货 (1) → 待签收 (2) → 已完成 (3)
   - 待支付 (0) → 已取消 (4)
3. **订单超时自动取消**: 下单后30分钟未支付自动取消
4. **订单取消条件**: 仅待支付状态的订单可取消

### 8.2 库存业务规则

1. **库存预扣机制**: 下单时先锁定库存，支付成功后扣减
2. **库存释放机制**: 订单取消或超时未支付自动释放锁定库存
3. **库存预警**: 可用库存低于预警值时触发告警
4. **库存防超卖**: 使用数据库行级锁保证库存扣减原子性

### 8.3 积分业务规则

1. **积分发放时机**: 订单完成后发放
2. **积分计算公式**: `积分 = 订单金额 × 会员等级积分倍率`
3. **积分过期规则**: 积分有效期为1年，过期自动失效
4. **会员升级规则**: 积分达到下一等级要求时自动升级
5. **积分抵扣规则**: 100积分 = 1元，可抵扣订单金额

### 8.4 物流业务规则

1. **物流单号生成**: 采用雪花算法 + 物流公司前缀
2. **物流状态流转**:
   - 待发货 (1) → 已发货 (2) → 运输中 (3) → 已签收 (4)
   - 已发货 (2) → 已拒收 (5)
   - 待发货 (1) → 已取消 (6)
3. **物流时效要求**: 下单后24小时内发货

## 九、事务一致性保证

### 9.1 订单创建与库存锁定

```go
// 伪代码示例
func CreateOrder(req *CreateOrderReq) error {
    return db.Transaction(func(tx *gorm.DB) error {
        // 1. 检查库存（使用悲观锁）
        var inventory Inventory
        err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
            Where("product_id = ? AND available_stock >= ?", req.ProductID, req.Quantity).
            First(&inventory).Error
        if err != nil {
            return err // 库存不足
        }

        // 2. 锁定库存
        err = tx.Model(&inventory).
            Updates(map[string]interface{}{
                "available_stock": gorm.Expr("available_stock - ?", req.Quantity),
                "locked_stock":    gorm.Expr("locked_stock + ?", req.Quantity),
            }).Error
        if err != nil {
            return err
        }

        // 3. 创建订单
        order := &Order{
            OrderNo:     generateOrderNo(),
            ProductID:   req.ProductID,
            Quantity:    req.Quantity,
            TotalAmount: req.Price * float64(req.Quantity),
            Status:      OrderStatusPendingPay,
        }
        err = tx.Create(order).Error
        if err != nil {
            return err
        }

        return nil
    })
}
```

### 9.2 订单支付与库存扣减

```go
func PayOrder(orderNo string) error {
    return db.Transaction(func(tx *gorm.DB) error {
        // 1. 更新订单状态
        err := tx.Model(&Order{}).
            Where("order_no = ? AND status = ?", orderNo, OrderStatusPendingPay).
            Update("status", OrderStatusPendingShip).Error
        if err != nil {
            return err
        }

        // 2. 获取订单信息
        var order Order
        err = tx.Where("order_no = ?", orderNo).First(&order).Error
        if err != nil {
            return err
        }

        // 3. 扣减锁定库存
        err = tx.Model(&Inventory{}).
            Where("product_id = ?", order.ProductID).
            Updates(map[string]interface{}{
                "locked_stock": gorm.Expr("locked_stock - ?", order.Quantity),
            }).Error
        if err != nil {
            return err
        }

        // 4. 增加商品销量
        err = tx.Model(&Product{}).
            Where("id = ?", order.ProductID).
            Update("sales", gorm.Expr("sales + ?", order.Quantity)).Error
        if err != nil {
            return err
        }

        return nil
    })
}
```

### 9.3 订单取消与库存释放

```go
func CancelOrder(orderNo string) error {
    return db.Transaction(func(tx *gorm.DB) error {
        // 1. 获取订单信息
        var order Order
        err := tx.Where("order_no = ?", orderNo).First(&order).Error
        if err != nil {
            return err
        }

        // 2. 更新订单状态
        err = tx.Model(&Order{}).
            Where("order_no = ?", orderNo).
            Update("status", OrderStatusCancelled).Error
        if err != nil {
            return err
        }

        // 3. 释放锁定库存
        err = tx.Model(&Inventory{}).
            Where("product_id = ?", order.ProductID).
            Updates(map[string]interface{}{
                "locked_stock":    gorm.Expr("locked_stock - ?", order.Quantity),
                "available_stock": gorm.Expr("available_stock + ?", order.Quantity),
            }).Error
        if err != nil {
            return err
        }

        return nil
    })
}
```

## 十、异常处理与补偿机制

### 10.1 库存异常处理

1. **库存不足**: 返回明确错误，提示用户
2. **库存锁定失败**: 回滚订单创建，释放资源
3. **库存释放失败**: 记录异常日志，人工介入处理

### 10.2 支付异常处理

1. **支付超时**: 订单自动取消，释放锁定库存
2. **支付失败**: 订单状态保持待支付，用户可重试
3. **支付回调重复幂等**: 使用唯一订单号+支付流水号保证幂等性

### 10.3 积分异常处理

1. **积分发放失败**: 记录失败日志，支持手动补发
2. **会员升级失败**: 记录异常，不影响积分发放
3. **积分重复发放**: 使用订单号+会员ID作为幂等键

## 十一、性能优化建议

1. **数据库优化**:
   - 为高频查询字段添加索引（order_no, member_id, product_id等）
   - 读写分离，查询操作走从库
   - 分库分表，按会员ID或订单号分片

2. **缓存策略**:
   - 商品库存信息缓存到Redis，设置过期时间
   - 会员等级信息缓存，减少数据库查询
   - 使用Redis分布式锁防止超卖

3. **异步处理**:
   - 积分发放、会员升级异步处理
   - 物流信息更新异步处理
   - 使用消息队列解耦服务

4. **定时任务**:
   - 超时订单自动取消（每分钟执行）
   - 积分过期处理（每日执行）
   - 库存预警检查（每小时执行）