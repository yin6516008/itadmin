# IT 技能管理后台

基于 Go + React 的全栈后台管理系统。

## 技术栈

**后端**

- Go 1.23 / Gin / GORM
- PostgreSQL
- 分层架构：Handler → Service → Repository

**前端**

- React 19 / TypeScript / Vite
- shadcn/ui / Tailwind CSS
- TanStack Query + TanStack Table
- React Router / Zustand / React Hook Form + Zod

## 项目结构

```
├── backend/                 # Go 后端
│   ├── cmd/myapp/           # 入口
│   ├── internal/
│   │   ├── config/          # 配置加载
│   │   ├── errcode/         # 错误码
│   │   ├── handler/         # HTTP 处理器
│   │   ├── infra/           # 数据库、日志等基础设施
│   │   ├── model/           # 数据模型与 DTO
│   │   ├── repository/      # 数据访问层
│   │   └── service/         # 业务逻辑层
│   ├── config.yaml          # 配置文件（需自建）
│   ├── Makefile
│   └── go.mod
├── frontend/                # React 前端
│   ├── src/
│   │   ├── components/      # 布局、共享组件、UI 组件
│   │   ├── features/        # 业务模块（用户管理等）
│   │   ├── hooks/           # 通用 Hooks
│   │   ├── lib/             # 工具函数、API 客户端、常量
│   │   ├── pages/           # 页面组件
│   │   ├── routes/          # 路由与菜单配置
│   │   ├── stores/          # 状态管理
│   │   └── types/           # 类型定义
│   ├── package.json
│   └── vite.config.ts
└── .cursor/rules/           # AI 编码规范
```

## 快速开始

### 前置要求

- Go >= 1.23
- Node.js >= 18
- PostgreSQL

### 后端

```bash
cd backend

# 创建配置文件
cp config.yaml.example config.yaml
# 编辑 config.yaml，填入数据库连接信息

# 安装依赖并启动
go mod tidy
make run
```

后端默认监听 `http://localhost:8080`。

### 前端

```bash
cd frontend

npm install
npm run dev
```

前端默认运行在 `http://localhost:5173`，开发模式下 API 请求自动代理到后端。

### 配置文件

后端配置文件 `backend/config.yaml`（不纳入版本控制）：

```yaml
app:
  env: development

server:
  addr: ":8080"

database:
  dsn: "host=localhost port=5432 user=postgres password=YOUR_PASSWORD dbname=postgres sslmode=disable"
```

## 功能模块

- 登录认证（开发环境 Mock）
- 用户管理：增删改查、搜索、启用/禁用
- 侧边栏导航（可折叠）
- 分页数据表格

## 开发命令

| 命令 | 说明 |
|------|------|
| `cd backend && make run` | 启动后端 |
| `cd backend && make build` | 编译后端 |
| `cd backend && make test` | 运行后端测试 |
| `cd backend && make lint` | 后端代码检查 |
| `cd frontend && npm run dev` | 启动前端开发服务器 |
| `cd frontend && npm run build` | 构建前端生产包 |
| `cd frontend && npm run lint` | 前端代码检查 |

## License

MIT
