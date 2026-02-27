---
description: 基于 Gin 框架的后端 API 项目架构规范，包含目录结构、分层依赖、API 设计、日志等约定。
globs: backend/**
alwaysApply: true
---

# Gin 后端 API 规范

基于 Gin 框架构建 HTTP API 的项目架构与编码约定。Go 语言通用规范见 `golang-patterns.md`。

---

## 一、项目结构规范

### 标准项目目录结构

```text
myproject/
├── cmd/
│   └── myapp/
│       └── main.go           # 程序入口
├── internal/
│   ├── config/               # 配置加载（读取 YAML）
│   ├── errcode/              # 业务错误码
│   ├── model/                # 实体 + Request 结构体 + Response DTO
│   ├── handler/              # HTTP 处理器 + response + middleware
│   ├── service/              # 业务逻辑
│   ├── repository/           # 数据访问
│   ├── adapter/              # 第三方 SDK 封装（按需创建）
│   └── infra/                # 基础设施（DB、Redis、Logger 初始化 + 日志上下文工具）
├── pkg/
│   └── client/               # 对外公开的 API 客户端
├── config.yaml               # 应用配置文件
├── go.mod
└── Makefile
```

### handler 目录内部结构

handler 不再拆子包，全部扁平化为文件：

```text
internal/handler/
├── user.go              # UserHandler
├── order.go             # OrderHandler
├── response.go          # Success / Error 统一响应
└── middleware.go         # RequestLogger 等中间件
```

### 目录职责对比

| 目录 | 职责 | 典型内容 |
|---|---|---|
| `internal/model/` | 实体 + 所有 Request/Response 结构体 | User、CreateUserRequest、UserBriefResponse |
| `internal/handler/` | HTTP 路由、参数校验、统一响应、中间件 | handler、response.go、middleware.go |
| `internal/service/` | 业务逻辑编排，通过接口依赖 repository 和 adapter | 业务规则、流程协调 |
| `internal/repository/` | 数据访问层，使用 infra 提供的连接执行读写 | CRUD、缓存读写、查询构造 |
| `internal/adapter/` | 第三方业务服务 SDK 的封装 | 短信、OSS、支付网关、地图 API |
| `internal/infra/` | 基础设施初始化 + 日志工具 | DB、Redis 连接池，Logger 初始化，日志上下文工具 |
| `internal/errcode/` | 业务错误码，底层包 | AppError、错误码常量 |
| `pkg/client/` | 本项目对外暴露的客户端库 | 调用本项目 API 的 SDK |

### 配置加载规范

使用 **YAML 文件** 管理应用配置，禁止通过环境变量分散配置。配置文件默认为项目根目录的 `config.yaml`。

```yaml
# config.yaml
app:
  env: development

server:
  addr: ":8080"

database:
  dsn: "host=localhost user=postgres password=postgres dbname=myproject port=5432 sslmode=disable"
```

```go
// internal/config/config.go
package config

// Config 应用配置根结构，字段与 config.yaml 一一对应。
type Config struct {
    App      AppConfig      `yaml:"app"`
    Server   ServerConfig   `yaml:"server"`
    Database DatabaseConfig `yaml:"database"`
}

type AppConfig struct {
    Env string `yaml:"env"`
}

type ServerConfig struct {
    Addr string `yaml:"addr"`
}

type DatabaseConfig struct {
    DSN string `yaml:"dsn"`
}

// Load 从指定路径读取 YAML 配置文件并解析。
func Load(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("read config file: %w", err)
    }
    var cfg Config
    if err := yaml.Unmarshal(data, &cfg); err != nil {
        return nil, fmt.Errorf("parse config file: %w", err)
    }
    return &cfg, nil
}
```

**规则**：
- 新增配置项时，先在 `config.yaml` 加字段，再在 `Config` 结构体加对应字段和 `yaml` tag
- 敏感信息（密码、密钥）不提交到仓库，使用 `config.local.yaml`（加入 `.gitignore`）
- YAML 依赖使用 `gopkg.in/yaml.v3`

简而言之：**infra 管"连接"，repository 管"读写"，adapter 管"调别人"，client 管"别人调我"**。

---

## 二、分层依赖规范

### 依赖方向图

```text
                    ┌──────────┐
                    │   cmd/   │  程序入口，组装所有依赖
                    └────┬─────┘
                         │
                         ▼
                  ┌──────────────┐
                  │   handler/   │  HTTP 路由 & 参数校验
                  └──────┬───────┘
                         │ 依赖 service
                         ▼
        ┌────────────────────────────────┐
        │           service/             │  业务逻辑编排
        └──┬─────────────┬──────────────┘
           │             │
           ▼             ▼
    ┌────────────┐ ┌──────────┐
    │ repository/│ │ adapter/ │
    │ (数据访问) │ │(第三方SDK)│
    └─────┬──────┘ └──────────┘
          │
          ▼
    ┌──────────────────────────────────────────────┐
    │  infra/  │  model/  │  errcode/             │  ← 底层包，所有层都可依赖
    │ DB/Redis │  实体/DTO │  业务错误码           │
    │ Logger   │          │                       │
    └──────────────────────────────────────────────┘
```

### 依赖规则

| 规则 | 说明 |
|---|---|
| handler → service | handler 直接依赖 service，不直接访问 repository |
| service → repository, adapter | service 通过接口调用 repository 和 adapter |
| repository → infra | repository 使用 infra 提供的 DB/Redis 连接 |
| model、errcode、infra 被所有层依赖 | 底层包，不依赖其他层 |
| **禁止** handler → repository | handler 不能跳过 service 直接访问数据 |
| **禁止** repository → service | 下层不能反向依赖上层 |

### 通过接口解耦

service 在自身包内定义接口，repository 实现该接口。

```go
// internal/service/user.go
package service

type UserRepo interface {
    Save(ctx context.Context, user *model.User) error
    FindByID(ctx context.Context, id string) (*model.User, error)
}

type UserService struct {
    repo UserRepo
}

func NewUserService(repo UserRepo) *UserService {
    return &UserService{repo: repo}
}
```

```go
// internal/repository/user.go
package repository

type UserRepo struct {
    db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
    return &UserRepo{db: db}
}

func (r *UserRepo) Save(ctx context.Context, user *model.User) error {
    return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepo) FindByID(ctx context.Context, id string) (*model.User, error) {
    var user model.User
    if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
        return nil, err
    }
    return &user, nil
}
```

### 在 cmd/main.go 中组装所有依赖

```go
func main() {
    configPath := flag.String("c", "config.yaml", "配置文件路径")
    flag.Parse()

    cfg, err := config.Load(*configPath)
    if err != nil {
        slog.Error("failed to load config", slog.String("error", err.Error()))
        os.Exit(1)
    }
    infra.InitLogger(cfg.App.Env)

    db, err := infra.NewDatabase(cfg.Database.DSN)
    if err != nil {
        slog.Error("failed to connect database", slog.String("error", err.Error()))
        os.Exit(1)
    }

    userRepo := repository.NewUserRepo(db)
    smsSender := sms.New(cfg.SMSAccessKey, cfg.SMSSecret)

    userService := service.NewUserService(userRepo)
    notificationService := service.NewNotificationService(smsSender)

    userHandler := handler.NewUserHandler(userService)

    router := gin.Default()
    router.Use(handler.RequestLogger())

    v1 := router.Group("/api/v1")
    userHandler.RegisterRoutes(v1)

    router.Run(cfg.Server.Addr)
}
```

---

## 三、API 规范

### model 包结构体规范

每个领域一个文件，实体 + Request + Response DTO 放在一起：

```go
// internal/model/user.go
package model

// ——— 实体（使用 GORM tag 定义字段约束，软删除使用 gorm.DeletedAt）———

type User struct {
    ID           string         `json:"id" gorm:"primaryKey;size:36"`
    Name         string         `json:"name" gorm:"size:50;not null"`
    Email        string         `json:"email" gorm:"size:100;uniqueIndex;not null"`
    Phone        string         `json:"phone" gorm:"size:20"`
    CreatedAt    time.Time      `json:"created_at"`
    UpdatedAt    time.Time      `json:"updated_at"`
    PasswordHash string         `json:"-" gorm:"size:255"`
    DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// ——— Request（所有参数结构体统一命名 XxxRequest）———

type CreateUserRequest struct {
    Name  string `json:"name" binding:"required,max=50"`
    Email string `json:"email" binding:"required,email"`
}

type UpdateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

type ListUsersRequest struct {
    Page int `form:"page" binding:"min=1"`
    Size int `form:"size" binding:"min=1,max=100"`
}

// ——— 分页响应 ———

type PageResult struct {
    List  any   `json:"list"`
    Total int64 `json:"total"`
    Page  int   `json:"page"`
    Size  int   `json:"size"`
}

// ——— Response DTO（仅在与实体差异大时定义）———

type UserBriefResponse struct {
    ID     string `json:"id"`
    Name   string `json:"name"`
    Avatar string `json:"avatar"`
}

func NewUserBriefResponse(u *User) UserBriefResponse {
    return UserBriefResponse{ID: u.ID, Name: u.Name, Avatar: u.Avatar}
}
```

**命名规则**：

| 类型 | 命名 | 示例 |
|---|---|---|
| 参数结构体 | `XxxRequest` | `CreateUserRequest`、`ListUsersRequest` |
| 响应结构体 | `XxxResponse` | `UserBriefResponse` |
| 实体 | 领域名词 | `User`、`Order` |

**Response 策略**：字段与实体高度重合时直接返回实体（敏感字段用 `json:"-"`），差异大时定义 `XxxResponse`。

**禁止使用 copier 等反射库做自动字段映射**。

### 参数来源规范

| 参数来源 | 绑定方式 | 是否需要 Request 结构体 |
|---|---|---|
| **Path** `/users/:id` | `c.Param("id")` | 不需要 |
| **Query** `?page=1&size=20` | `c.ShouldBindQuery(&req)`，tag 用 `form` | 参数 >2 个时定义 |
| **Body** JSON | `c.ShouldBindJSON(&req)`，tag 用 `json` + `binding` | 始终定义 |

### Service 层参数规范

handler 和 service **共用 `model.XxxRequest`**，handler 绑定后直接传给 service，无需转换。

```go
// service 方法签名
func (s *UserService) Create(ctx context.Context, req *model.CreateUserRequest) (*model.User, error)
func (s *UserService) Update(ctx context.Context, id string, req *model.UpdateUserRequest) error
func (s *UserService) List(ctx context.Context, req *model.ListUsersRequest) ([]*model.User, int64, error)

// 简单 ID 查询不需要 Request 结构体，直接传字段
func (s *UserService) GetByID(ctx context.Context, id string) (*model.User, error)
func (s *UserService) Delete(ctx context.Context, id string) error
```

### 统一响应与错误码

```go
// internal/errcode/errcode.go
type AppError struct {
    HTTPStatus int    `json:"-"`
    Code       int    `json:"code"`
    Message    string `json:"message"`
}

func (e *AppError) Error() string { return e.Message }

var (
    ErrInvalidParam = &AppError{HTTPStatus: 400, Code: 40001, Message: "参数校验失败"}
    ErrUnauthorized = &AppError{HTTPStatus: 401, Code: 40100, Message: "未登录"}
    ErrForbidden    = &AppError{HTTPStatus: 403, Code: 40300, Message: "无权限"}
    ErrNotFound     = &AppError{HTTPStatus: 404, Code: 40400, Message: "资源不存在"}
    ErrInternal     = &AppError{HTTPStatus: 500, Code: 50000, Message: "服务器内部错误"}
)
```

```go
// internal/handler/response.go
package handler

type responseBody struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Data    any    `json:"data,omitempty"`
}

func Success(c *gin.Context, data any) {
    c.JSON(http.StatusOK, responseBody{Code: 0, Message: "ok", Data: data})
}

func Error(c *gin.Context, err error) {
    var appErr *errcode.AppError
    if errors.As(err, &appErr) {
        c.JSON(appErr.HTTPStatus, responseBody{Code: appErr.Code, Message: appErr.Message})
    } else {
        c.JSON(500, responseBody{Code: 50000, Message: "服务器内部错误"})
    }
    c.Abort()
}
```

service 直接返回 `*errcode.AppError`，handler 透传给 `Error`，不做错误类型判断。

### Handler 标准写法

所有 handler 遵循**三段式**：绑定参数 → 调 service → 返回结果。

```go
// internal/handler/user.go
package handler

type UserHandler struct {
    service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
    return &UserHandler{service: s}
}

func (h *UserHandler) Create(c *gin.Context) {
    var req model.CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        Error(c, errcode.ErrInvalidParam)
        return
    }

    user, err := h.service.Create(c.Request.Context(), &req)
    if err != nil {
        Error(c, err)
        return
    }

    Success(c, user)
}

func (h *UserHandler) GetByID(c *gin.Context) {
    id := c.Param("id")

    user, err := h.service.GetByID(c.Request.Context(), id)
    if err != nil {
        Error(c, err)
        return
    }

    Success(c, user)
}

func (h *UserHandler) List(c *gin.Context) {
    var req model.ListUsersRequest
    if err := c.ShouldBindQuery(&req); err != nil {
        Error(c, errcode.ErrInvalidParam)
        return
    }

    users, total, err := h.service.List(c.Request.Context(), &req)
    if err != nil {
        Error(c, err)
        return
    }

    Success(c, model.PageResult{
        List: users, Total: total,
        Page: req.Page, Size: req.Size,
    })
}
```

---

## 四、日志规范

使用标准库 `log/slog`，不引入第三方日志库。

### 初始化

```go
// internal/infra/logger.go
func InitLogger(env string) {
    opts := &slog.HandlerOptions{AddSource: true}
    if env == "production" {
        opts.Level = slog.LevelInfo
        slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, opts)))
    } else {
        opts.Level = slog.LevelDebug
        slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, opts)))
    }
}
```

### 日志上下文工具

`infra/logger.go` 同时提供日志初始化和 context 传递工具：

```go
// internal/infra/logger.go（接 InitLogger 之后）

type logCtxKey struct{}

// WithLogger 将带请求级字段的 logger 存入 context。
func WithLogger(ctx context.Context, l *slog.Logger) context.Context {
    return context.WithValue(ctx, logCtxKey{}, l)
}

// Logger 从 context 取出 logger，不存在时返回 slog.Default()。
func Logger(ctx context.Context) *slog.Logger {
    if l, ok := ctx.Value(logCtxKey{}).(*slog.Logger); ok {
        return l
    }
    return slog.Default()
}
```

```go
// internal/handler/middleware.go
package handler

func RequestLogger() gin.HandlerFunc {
    return func(c *gin.Context) {
        traceID := c.GetHeader("X-Trace-ID")
        if traceID == "" {
            traceID = uuid.NewString()
        }

        logger := slog.Default().With(
            slog.String("trace_id", traceID),
            slog.String("method", c.Request.Method),
            slog.String("path", c.Request.URL.Path),
        )

        if userID, exists := c.Get("user_id"); exists {
            logger = logger.With(slog.String("user_id", userID.(string)))
        }

        ctx := infra.WithLogger(c.Request.Context(), logger)
        c.Request = c.Request.WithContext(ctx)
        c.Next()
    }
}
```

### 业务代码使用

```go
func (s *UserService) Create(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
    infra.Logger(ctx).Info("creating user", slog.String("email", req.Email))

    user, err := s.repo.Save(ctx, &model.User{Name: req.Name, Email: req.Email})
    if err != nil {
        infra.Logger(ctx).Error("failed to create user", slog.String("error", err.Error()))
        return nil, errcode.ErrInternal
    }
    return user, nil
}
```

### 日志级别

| 级别 | 使用场景 |
|---|---|
| **Debug** | 开发调试，生产不输出 |
| **Info** | 关键业务节点：注册、下单、支付 |
| **Warn** | 异常但不影响主流程：重试成功、降级 |
| **Error** | 需要关注的错误：DB 连接失败、第三方超时 |

### 禁止事项

- 禁止用 `fmt.Println` / `log.Println` 打日志
- 禁止记录密码、token、身份证号等敏感信息
- 禁止用 Error 级别记录业务预期错误（如"用户不存在"应为 Info 或 Warn）
- Error 日志必须带上下文（操作名、资源标识、错误信息）

---

## 五、注释规范

遵循 Go 官方 [godoc 约定](https://go.dev/blog/godoc)，注释既是文档也是 API 合约。

### 包注释

每个 package 必须在一个文件顶部写 **package-level doc comment**，说明该包的职责：

```go
// Package errcode 定义统一的业务错误码，供 handler 和 service 层使用。
package errcode
```

```go
// Package handler 实现 HTTP 路由注册、参数校验和统一响应。
package handler
```

一般放在包内最核心的文件中（如 `errcode/errcode.go`、`handler/response.go`）。

### 导出符号注释

所有导出的类型、函数、方法、常量、变量 **必须** 有 godoc 注释，以符号名称开头：

```go
// AppError 表示可直接返回给前端的业务错误，包含 HTTP 状态码和业务码。
type AppError struct { ... }

// NewUserHandler 创建 UserHandler，依赖 UserService。
func NewUserHandler(s *service.UserService) *UserHandler { ... }

// RegisterRoutes 将用户相关路由注册到指定路由组。
func (h *UserHandler) RegisterRoutes(rg *gin.RouterGroup) { ... }

// Success 返回统一的成功 JSON 响应，code=0。
func Success(c *gin.Context, data any) { ... }

// ErrNotFound 表示请求的资源不存在。
var ErrNotFound = &AppError{...}
```

### 接口注释

接口本身和每个方法都需要注释：

```go
// UserRepo 定义用户数据访问层的接口，由 repository.UserRepo 实现。
type UserRepo interface {
    // Save 创建一条用户记录。
    Save(ctx context.Context, user *model.User) error
    // FindByID 根据 ID 查询用户，未找到时返回 gorm.ErrRecordNotFound。
    FindByID(ctx context.Context, id string) (*model.User, error)
}
```

### 不需要注释的情况

| 场景 | 示例 | 说明 |
|---|---|---|
| 未导出符号 | `type ctxKey struct{}` | 包内私有，无需 godoc |
| 意图自明的简单函数 | `func (e *AppError) Error() string` | 实现标准接口，签名即文档 |
| 结构体字段 | `Name string` | 字段名 + tag 已说明，无需逐字段注释 |

### 代码内分隔注释

同一文件内按逻辑分区时，使用 `// ——— 区域名 ———` 格式：

```go
// ——— 实体 ———

type User struct { ... }

// ——— Request ———

type CreateUserRequest struct { ... }
```

### 禁止事项

- **禁止**复述代码的注释（如 `// 返回错误`、`// 遍历用户列表`）
- **禁止**注释掉的废弃代码留在仓库中，用版本控制代替
- **禁止**用 `/* */` 块注释写 godoc（用 `//` 行注释）
- **禁止**在注释中包含 TODO/FIXME 而不附带负责人和时间（格式：`// TODO(zhangsan): 2025-03 优化分页查询`）

---

## 六、各层代码示例（含注释）

### infra 层

```go
// internal/infra/database.go
package infra

import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func NewDatabase(dsn string) (*gorm.DB, error) {
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, fmt.Errorf("open database: %w", err)
    }
    sqlDB, err := db.DB()
    if err != nil {
        return nil, fmt.Errorf("get underlying sql.DB: %w", err)
    }
    sqlDB.SetMaxOpenConns(25)
    sqlDB.SetMaxIdleConns(5)
    return db, nil
}
```

```go
// internal/infra/redis.go
package infra

func NewRedis(addr, password string, db int) *redis.Client {
    return redis.NewClient(&redis.Options{Addr: addr, Password: password, DB: db})
}
```

### adapter 层

```go
// internal/adapter/sms/sms.go
package sms

type Sender struct {
    accessKey string
    secret    string
}

func New(accessKey, secret string) *Sender {
    return &Sender{accessKey: accessKey, secret: secret}
}

func (s *Sender) Send(phone, content string) error {
    // 调用第三方短信 SDK
    return nil
}
```

### service 层

```go
// internal/service/notification.go
package service

type SMSSender interface {
    Send(phone, content string) error
}

type NotificationService struct {
    sms SMSSender
}

func NewNotificationService(sms SMSSender) *NotificationService {
    return &NotificationService{sms: sms}
}

func (n *NotificationService) NotifyUser(phone, msg string) error {
    return n.sms.Send(phone, msg)
}
```

---

## 七、数据库规范（GORM + PostgreSQL）

### ORM 选型

使用 **GORM** 作为 ORM，搭配 **PostgreSQL** 数据库。禁止在 repository 层直接拼写原生 SQL（特殊复杂查询除外，需注释说明原因）。

### Model 与 GORM Tag

实体结构体同时承载 JSON 序列化和 GORM 映射，通过 tag 区分：

```go
type User struct {
    ID        string         `json:"id" gorm:"primaryKey;size:36"`
    Name      string         `json:"name" gorm:"size:50;not null"`
    Email     string         `json:"email" gorm:"size:100;uniqueIndex;not null"`
    Phone     string         `json:"phone" gorm:"size:20"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
```

| GORM Tag | 说明 |
|---|---|
| `primaryKey` | 主键 |
| `size:N` | 字段长度 |
| `not null` | 非空约束 |
| `uniqueIndex` | 唯一索引 |
| `index` | 普通索引 |
| `default:value` | 默认值 |
| `column:name` | 指定列名（仅在字段名不符合 GORM 自动映射时使用） |

### 软删除

统一使用 GORM 内置软删除，实体包含 `DeletedAt gorm.DeletedAt` 字段：

- `db.Delete(&user)` 自动设置 `deleted_at` 时间戳
- 所有查询自动过滤已删除记录，无需手动加 `WHERE deleted_at IS NULL`
- 查询包含已删除记录：`db.Unscoped().Find(&users)`

### Repository 层 GORM 使用规范

```go
type UserRepo struct {
    db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
    return &UserRepo{db: db}
}
```

| 规则 | 说明 |
|---|---|
| 始终使用 `WithContext(ctx)` | 每个数据库操作都必须传递 context，确保链路追踪和超时控制 |
| 错误直接返回 | repository 不判断错误类型，由 service 层处理 `gorm.ErrRecordNotFound` 等 |
| 分页使用 `Offset` + `Limit` | 配合 `Count` 查总数 |
| 禁止在 repository 中使用 `AutoMigrate` | 迁移统一在 `cmd/main.go` 启动时执行 |

```go
// 创建
func (r *UserRepo) Save(ctx context.Context, user *model.User) error {
    return r.db.WithContext(ctx).Create(user).Error
}

// 查询单个（未找到时返回 gorm.ErrRecordNotFound）
func (r *UserRepo) FindByID(ctx context.Context, id string) (*model.User, error) {
    var user model.User
    if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
        return nil, err
    }
    return &user, nil
}

// 更新
func (r *UserRepo) Update(ctx context.Context, user *model.User) error {
    return r.db.WithContext(ctx).Save(user).Error
}

// 删除（软删除）
func (r *UserRepo) Delete(ctx context.Context, id string) error {
    return r.db.WithContext(ctx).Delete(&model.User{}, "id = ?", id).Error
}

// 分页列表
func (r *UserRepo) List(ctx context.Context, offset, limit int) ([]*model.User, int64, error) {
    var total int64
    if err := r.db.WithContext(ctx).Model(&model.User{}).Count(&total).Error; err != nil {
        return nil, 0, err
    }
    var users []*model.User
    if err := r.db.WithContext(ctx).Order("created_at DESC").Offset(offset).Limit(limit).Find(&users).Error; err != nil {
        return nil, 0, err
    }
    return users, total, nil
}
```

### Service 层错误处理

service 使用 `errors.Is(err, gorm.ErrRecordNotFound)` 判断"未找到"，其他错误统一返回 `errcode.ErrInternal`：

```go
func (s *UserService) GetByID(ctx context.Context, id string) (*model.User, error) {
    user, err := s.repo.FindByID(ctx, id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errcode.ErrNotFound
        }
        infra.Logger(ctx).Error("failed to get user", slog.String("id", id), slog.String("error", err.Error()))
        return nil, errcode.ErrInternal
    }
    return user, nil
}
```

### 数据库迁移

开发阶段使用 `AutoMigrate` 自动同步表结构，在 `cmd/main.go` 启动时执行：

```go
if err := db.AutoMigrate(&model.User{}, &model.Order{}); err != nil {
    slog.Error("failed to migrate database", slog.String("error", err.Error()))
    os.Exit(1)
}
```

生产环境建议使用 SQL 迁移文件（`migrations/` 目录），由 CI/CD 流程管理。

### 禁止事项

- 禁止在 handler 或 service 层直接操作 `*gorm.DB`，必须通过 repository
- 禁止使用 `db.Exec` 拼接 SQL 字符串（防注入），使用 GORM 的参数占位符
- 禁止在循环中执行单条查询（N+1 问题），使用 `Preload` 或批量查询
- 禁止省略 `WithContext(ctx)`
