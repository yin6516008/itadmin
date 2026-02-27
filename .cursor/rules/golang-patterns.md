---
description: Go 语言通用编码规范与最佳实践，适用于所有 Go 项目。
globs: **/*.go
alwaysApply: true
---

# Go 语言规范

Go 语言通用的惯用模式与最佳实践，适用于所有 Go 项目（后端 API、CLI 工具、库等）。

---

## 一、编码规范

### 核心原则

**简洁与清晰**：Go 崇尚简洁而非炫技。代码应当一目了然、易于阅读。

```go
// 好：清晰直接
func GetUser(id string) (*User, error) {
    user, err := db.FindUser(id)
    if err != nil {
        return nil, fmt.Errorf("get user %s: %w", id, err)
    }
    return user, nil
}

// 坏：过度炫技
func GetUser(id string) (*User, error) {
    return func() (*User, error) {
        if u, e := db.FindUser(id); e == nil {
            return u, nil
        } else {
            return nil, e
        }
    }()
}
```

**让零值可用**：设计类型时，确保其零值无需初始化即可直接使用。

```go
// 好：零值即可用
type Counter struct {
    mu    sync.Mutex
    count int // 零值为 0，可直接使用
}

func (c *Counter) Inc() {
    c.mu.Lock()
    c.count++
    c.mu.Unlock()
}

// 坏：需要初始化
type BadCounter struct {
    counts map[string]int // nil map 会 panic
}
```

**接收接口，返回结构体**：函数应接收接口类型参数，返回具体类型。

```go
// 好：接收接口，返回具体类型
func ProcessData(r io.Reader) (*Result, error) {
    data, err := io.ReadAll(r)
    if err != nil {
        return nil, err
    }
    return &Result{Data: data}, nil
}

// 坏：返回接口（不必要地隐藏了实现细节）
func ProcessData(r io.Reader) (io.Reader, error) {
    // ...
}
```

### 包命名

```go
// 好：简短、小写、不用下划线
package http
package json
package user

// 坏：冗长、混合大小写、或冗余
package httpHandler
package json_parser
package userService // 冗余的 'Service' 后缀
```

### 避免包级状态

```go
// 坏：全局可变状态
var db *sql.DB

func init() {
    db, _ = sql.Open("postgres", os.Getenv("DATABASE_URL"))
}

// 好：依赖注入
type Server struct {
    db *sql.DB
}

func NewServer(db *sql.DB) *Server {
    return &Server{db: db}
}
```

### 应避免的反模式

```go
// 坏：长函数中使用裸返回
func process() (result int, err error) {
    // ... 50 行 ...
    return // 返回的是什么？
}

// 坏：使用 panic 做流程控制
func GetUser(id string) *User {
    user, err := db.Find(id)
    if err != nil {
        panic(err) // 不要这样做
    }
    return user
}

// 坏：在结构体中传递 context
type Request struct {
    ctx context.Context // Context 应作为第一个参数
    ID  string
}

// 好：Context 作为第一个参数
func ProcessRequest(ctx context.Context, id string) error {
    // ...
}

// 坏：混用值接收者和指针接收者
type Counter struct{ n int }
func (c Counter) Value() int { return c.n }    // 值接收者
func (c *Counter) Increment() { c.n++ }        // 指针接收者
// 选择一种风格并保持一致
```

---

## 二、错误处理规范

### 携带上下文的错误包装

```go
func LoadConfig(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("load config %s: %w", path, err)
    }

    var cfg Config
    if err := json.Unmarshal(data, &cfg); err != nil {
        return nil, fmt.Errorf("parse config %s: %w", path, err)
    }

    return &cfg, nil
}
```

### 自定义错误类型

```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed on %s: %s", e.Field, e.Message)
}

var (
    ErrNotFound     = errors.New("resource not found")
    ErrUnauthorized = errors.New("unauthorized")
    ErrInvalidInput = errors.New("invalid input")
)
```

### 使用 errors.Is 和 errors.As 检查错误

```go
func HandleError(err error) {
    if errors.Is(err, sql.ErrNoRows) {
        log.Println("No records found")
        return
    }

    var validationErr *ValidationError
    if errors.As(err, &validationErr) {
        log.Printf("Validation error on field %s: %s",
            validationErr.Field, validationErr.Message)
        return
    }

    log.Printf("Unexpected error: %v", err)
}
```

### 永远不要忽略错误

```go
// 坏：用空白标识符忽略错误
result, _ := doSomething()

// 好：处理错误，或明确说明为何可以安全忽略
result, err := doSomething()
if err != nil {
    return err
}

// 可接受：错误确实无关紧要的场景（极少）
_ = writer.Close() // 尽力清理，错误已在其他地方记录
```

---

## 三、接口与类型设计规范

### 小而专注的接口

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

// 按需组合接口
type ReadWriteCloser interface {
    Reader
    Writer
    Closer
}
```

### 在使用方定义接口

```go
// 在消费者包中定义，而非提供者包
package service

type UserStore interface {
    GetUser(id string) (*User, error)
    SaveUser(user *User) error
}

type Service struct {
    store UserStore
}
```

### 通过类型断言实现可选行为

```go
type Flusher interface {
    Flush() error
}

func WriteAndFlush(w io.Writer, data []byte) error {
    if _, err := w.Write(data); err != nil {
        return err
    }
    if f, ok := w.(Flusher); ok {
        return f.Flush()
    }
    return nil
}
```

### 函数选项模式

```go
type Server struct {
    addr    string
    timeout time.Duration
    logger  *log.Logger
}

type Option func(*Server)

func WithTimeout(d time.Duration) Option {
    return func(s *Server) { s.timeout = d }
}

func WithLogger(l *log.Logger) Option {
    return func(s *Server) { s.logger = l }
}

func NewServer(addr string, opts ...Option) *Server {
    s := &Server{
        addr:    addr,
        timeout: 30 * time.Second,
        logger:  log.Default(),
    }
    for _, opt := range opts {
        opt(s)
    }
    return s
}
```

### 嵌入实现组合

```go
type Logger struct {
    prefix string
}

func (l *Logger) Log(msg string) {
    fmt.Printf("[%s] %s\n", l.prefix, msg)
}

type Server struct {
    *Logger // 嵌入 - Server 获得 Log 方法
    addr    string
}
```

---

## 四、并发规范

### Worker Pool（工作池）

```go
func WorkerPool(jobs <-chan Job, results chan<- Result, numWorkers int) {
    var wg sync.WaitGroup
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for job := range jobs {
                results <- process(job)
            }
        }()
    }
    wg.Wait()
    close(results)
}
```

### 使用 Context 实现取消与超时

```go
func FetchWithTimeout(ctx context.Context, url string) ([]byte, error) {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("create request: %w", err)
    }

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("fetch %s: %w", url, err)
    }
    defer resp.Body.Close()

    return io.ReadAll(resp.Body)
}
```

### 优雅关闭

```go
func GracefulShutdown(server *http.Server) {
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        log.Fatalf("Server forced to shutdown: %v", err)
    }
}
```

### 使用 errgroup 协调多个 Goroutine

```go
import "golang.org/x/sync/errgroup"

func FetchAll(ctx context.Context, urls []string) ([][]byte, error) {
    g, ctx := errgroup.WithContext(ctx)
    results := make([][]byte, len(urls))

    for i, url := range urls {
        i, url := i, url
        g.Go(func() error {
            data, err := FetchWithTimeout(ctx, url)
            if err != nil {
                return err
            }
            results[i] = data
            return nil
        })
    }

    if err := g.Wait(); err != nil {
        return nil, err
    }
    return results, nil
}
```

### 避免 Goroutine 泄漏

```go
// 坏：context 取消时 goroutine 泄漏
func leakyFetch(ctx context.Context, url string) <-chan []byte {
    ch := make(chan []byte)
    go func() {
        data, _ := fetch(url)
        ch <- data // 没有接收者时永远阻塞
    }()
    return ch
}

// 好：正确处理取消
func safeFetch(ctx context.Context, url string) <-chan []byte {
    ch := make(chan []byte, 1)
    go func() {
        data, err := fetch(url)
        if err != nil {
            return
        }
        select {
        case ch <- data:
        case <-ctx.Done():
        }
    }()
    return ch
}
```

---

## 五、性能规范

### 已知大小时预分配切片

```go
// 坏：切片多次扩容
var results []Result
for _, item := range items {
    results = append(results, process(item))
}

// 好：一次分配
results := make([]Result, 0, len(items))
for _, item := range items {
    results = append(results, process(item))
}
```

### 使用 sync.Pool 优化频繁分配

```go
var bufferPool = sync.Pool{
    New: func() interface{} { return new(bytes.Buffer) },
}

func ProcessRequest(data []byte) []byte {
    buf := bufferPool.Get().(*bytes.Buffer)
    defer func() {
        buf.Reset()
        bufferPool.Put(buf)
    }()
    buf.Write(data)
    return buf.Bytes()
}
```

### 避免在循环中拼接字符串

```go
// 坏：产生大量字符串分配
result += p + ","

// 好：使用 strings.Builder
var sb strings.Builder
sb.WriteString(p)

// 最佳：使用标准库
strings.Join(parts, ",")
```

---

## 六、工具链规范

### 常用命令

```bash
go build ./...          # 构建
go test ./...           # 测试
go test -race ./...     # 竞态检测
go test -cover ./...    # 覆盖率
go vet ./...            # 静态分析
golangci-lint run       # Lint
go mod tidy             # 模块清理
gofmt -w .              # 格式化
goimports -w .          # 格式化 + 整理 import
```

### 推荐的 Linter 配置（.golangci.yml）

```yaml
linters:
  enable:
    - errcheck
    - gosimple
    - govet
    - ineffassign
    - staticcheck
    - unused
    - gofmt
    - goimports
    - misspell
    - unconvert
    - unparam

linters-settings:
  errcheck:
    check-type-assertions: true
  govet:
    check-shadowing: true

issues:
  exclude-use-default: false
```

---

## 附录：速查表

| 惯用法 | 说明 |
|---|---|
| 接收接口，返回结构体 | 函数接收接口参数，返回具体类型 |
| 错误即值 | 将错误视为一等值，而非异常 |
| 不要通过共享内存来通信 | 使用 channel 在 goroutine 间协调 |
| 让零值可用 | 类型应无需显式初始化即可工作 |
| 少量复制优于少量依赖 | 避免不必要的外部依赖 |
| 清晰优于聪明 | 可读性优先于炫技 |
| gofmt 不是谁的最爱，但是所有人的朋友 | 始终使用 gofmt/goimports 格式化 |
| 提前返回 | 先处理错误，保持正常路径不缩进 |

**记住**：好的 Go 代码应当"无聊"——可预测、一致、易于理解。拿不准时，保持简单。
