# GoChat - Go 语言 Web 聊天软件

一个基于 Go 语言开发的现代化 web 实时聊天应用。

## 🎯 项目状态

### ✅ 已完成的工作

#### 1. 项目结构搭建
- ✅ 完整的目录结构按照 Go 项目最佳实践创建
- ✅ `go.mod` 文件配置，包含所有必要依赖
- ✅ 配置管理系统（支持环境变量和配置文件）

#### 2. 数据库层
- ✅ GORM 数据模型定义（User、Room、Message、RoomMember）
- ✅ 数据库连接管理
- ✅ 自动迁移功能

#### 3. 核心组件
- ✅ 日志模块（基于 logrus）
- ✅ 统一响应格式
- ✅ JWT 工具函数
- ✅ 密码加密工具

#### 4. 业务服务
- ✅ 用户认证服务（注册/登录）
- ✅ 用户管理服务
- ✅ JWT 认证中间件

#### 5. 应用入口
- ✅ 主程序入口
- ✅ 基础 HTTP 服务器设置
- ✅ CORS 中间件
- ✅ 健康检查接口

#### 6. 框架迁移
- ✅ 从 Gin 完整迁移到 CloudWeGo Hertz 框架
- ✅ 所有路由和中间件适配 Hertz API
- ✅ 响应工具和认证中间件迁移完成

## 🚀 快速开始

### 1. 环境准备
```bash
# 确保已安装 Go 1.23.2+
go version

# 确保 MySQL 数据库运行
# 创建数据库 'gochat'
```

### 2. 配置环境变量
创建 `.env` 文件（或设置环境变量）：
```env
# Server
SERVER_PORT=8080
SERVER_MODE=debug

# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=你的密码
DB_NAME=gochat

# JWT
JWT_SECRET=你的密钥
JWT_EXPIRE_HOURS=24
```

### 3. 安装依赖并运行
```bash
# 安装依赖
go mod tidy

# 运行服务器
go run cmd/server/main.go
```

### 4. 测试接口
```bash
# 健康检查
curl http://localhost:8080/health

# API测试
curl http://localhost:8080/api/ping

# 认证模块测试
curl http://localhost:8080/api/auth/info

# 用户模块测试
curl http://localhost:8080/api/users/info

# 聊天模块测试
curl http://localhost:8080/api/chat/info
```

## 📋 下一步开发计划

### 🔄 当前需要完成的任务

#### 1. API 路由完善（优先级：高）
- [ ] 添加认证相关路由（注册/登录）
- [ ] 添加用户管理路由
- [ ] 创建路由处理器（handlers）

#### 2. WebSocket 模块（优先级：高）
- [ ] WebSocket 连接管理器
- [ ] 聊天室管理
- [ ] 消息广播系统

#### 3. 前端界面（优先级：中）
- [ ] 基础 HTML 模板
- [ ] 登录/注册页面
- [ ] 聊天室界面
- [ ] WebSocket 客户端

#### 4. 高级功能（优先级：低）
- [ ] Redis 集成
- [ ] 文件上传
- [ ] 消息历史
- [ ] 用户头像

## 🛠️ 技术栈

- **后端**: Go + CloudWeGo Hertz Framework
- **数据库**: MySQL + GORM
- **认证**: JWT
- **实时通信**: WebSocket
- **日志**: Logrus
- **配置**: 环境变量 + 配置文件

## 📁 项目结构

```
GoChat/
├── cmd/server/           # 应用入口
├── internal/            # 内部业务逻辑
│   ├── auth/           # 认证模块
│   ├── user/           # 用户管理
│   ├── config/         # 配置管理
│   ├── database/       # 数据库相关
│   ├── middleware/     # 中间件
│   ├── router/         # 路由配置
│   └── (待添加更多模块)
├── pkg/                # 可重用包
│   ├── logger/         # 日志
│   ├── response/       # 响应格式
│   └── utils/          # 工具函数
├── configs/            # 配置文件
├── web/               # 前端资源（待添加）
└── docs/              # 文档
```

## 🔧 开发指南

### 添加新的 API 接口
1. 在 `internal/router/router.go` 的 `setupAPIRoutes` 函数中添加路由
2. 创建对应的 handler 处理 HTTP 请求
3. 使用 Hertz 的处理函数签名：`func(ctx context.Context, c *app.RequestContext)`

### 数据库操作
- 使用 `database.DB` 全局实例
- 遵循 GORM 最佳实践
- 所有模型已定义在 `internal/database/models.go`

### 认证流程
1. 用户注册/登录获取 JWT token
2. 后续请求在 Header 中携带 `Authorization: Bearer <token>`
3. 使用 `AuthMiddleware` 保护需要认证的路由

## 📖 API 文档

### 当前可用接口

- `GET /health` - 健康检查
- `GET /api/ping` - API 测试
- `GET /api/auth/info` - 认证模块信息（占位）
- `GET /api/users/info` - 用户模块信息（占位）
- `GET /api/chat/info` - 聊天模块信息（占位）

### 计划中的接口

- `POST /api/auth/register` - 用户注册
- `POST /api/auth/login` - 用户登录  
- `GET /api/user/profile` - 获取用户信息
- `PUT /api/user/profile` - 更新用户信息
- `GET /api/users/online` - 获取在线用户
- `WebSocket /ws` - WebSocket 连接

## 🔄 Gin 到 Hertz 迁移完成

本项目已成功从 Gin 迁移到字节跳动的 CloudWeGo Hertz 框架，获得了更好的性能和更丰富的功能：

### 迁移亮点
- ✅ 高性能：Hertz 基于 Netpoll 网络库，性能更优
- ✅ 兼容性：API 设计与 Gin 类似，迁移成本低
- ✅ 扩展性：更好的中间件支持和插件生态
- ✅ 云原生：字节跳动内部大规模生产环境验证

### 主要改动
1. **框架依赖**：`gin-gonic/gin` → `cloudwego/hertz`
2. **处理函数**：增加 `context.Context` 参数
3. **启动方式**：`r.Run()` → `h.Spin()`
4. **响应方法**：`gin.H` → `utils.H`

## 🤝 贡献指南

1. 遵循 Go 代码规范
2. 添加适当的注释和文档
3. 确保代码通过测试
4. 提交前运行 `go fmt` 和 `go vet`

## �� 许可证

MIT License 