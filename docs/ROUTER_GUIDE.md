# GoChat 路由结构指南

## 📁 路由模块结构

路由逻辑已从 `cmd/server/main.go` 分离到独立的模块中，以提高代码的模块化和可维护性。

### 文件结构
```
internal/router/
└── router.go          # 路由配置主文件
```

## 🔧 路由模块设计

### 主要函数

#### `SetupRouter(addr string) *server.Hertz`
- **功能**: 初始化和配置所有路由
- **参数**: addr - 服务器监听地址
- **返回**: 配置好的 Hertz 引擎实例
- **调用位置**: `cmd/server/main.go`

#### 路由分组函数
- `setupBaseRoutes()` - 基础路由（健康检查等）
- `setupAPIRoutes()` - API 路由总入口，包含所有业务路由

## 🗺️ 当前路由映射

### 基础路由
| 方法 | 路径 | 功能 | 状态 |
|------|------|------|------|
| GET | `/health` | 健康检查 | ✅ 已实现 |

### API 路由组 (`/api`)
| 方法 | 路径 | 功能 | 状态 |
|------|------|------|------|
| GET | `/api/ping` | 连通性测试 | ✅ 已实现 |

### 认证路由组 (`/api/auth`)
| 方法 | 路径 | 功能 | 状态 |
|------|------|------|------|
| GET | `/api/auth/info` | 临时占位接口 | ✅ 占位 |
| POST | `/api/auth/register` | 用户注册 | 🔄 待实现 |
| POST | `/api/auth/login` | 用户登录 | 🔄 待实现 |
| POST | `/api/auth/logout` | 用户登出 | 🔄 待实现 |
| POST | `/api/auth/refresh` | 刷新 Token | 🔄 待实现 |

### 用户路由组 (`/api/users`)
| 方法 | 路径 | 功能 | 状态 |
|------|------|------|------|
| GET | `/api/users/info` | 临时占位接口 | ✅ 占位 |
| GET | `/api/users/profile` | 获取用户信息 | 🔄 待实现 |
| PUT | `/api/users/profile` | 更新用户信息 | 🔄 待实现 |
| GET | `/api/users/online` | 获取在线用户 | 🔄 待实现 |

### 聊天路由组 (`/api/chat`)
| 方法 | 路径 | 功能 | 状态 |
|------|------|------|------|
| GET | `/api/chat/info` | 临时占位接口 | ✅ 占位 |
| GET | `/api/chat/rooms` | 获取聊天室列表 | 🔄 待实现 |
| POST | `/api/chat/rooms` | 创建聊天室 | 🔄 待实现 |
| GET | `/api/chat/rooms/:id/messages` | 获取聊天记录 | 🔄 待实现 |
| GET | `/api/chat/ws` | WebSocket 连接 | 🔄 待实现 |

## 🚀 快速测试

### 启动服务器
```bash
go run cmd/server/main.go
```

### 手动测试
```bash
# 基础测试
curl http://localhost:8080/health
curl http://localhost:8080/api/ping

# 模块测试
curl http://localhost:8080/api/auth/info
curl http://localhost:8080/api/users/info
curl http://localhost:8080/api/chat/info

# CORS 测试
curl -I -X OPTIONS http://localhost:8080/api/ping
```

## 📝 添加新路由的步骤

### 1. 在 setupAPIRoutes 函数中添加路由
```go
// 例如：在认证路由组中添加注册接口
auth.POST("/register", handlers.Register)
```

### 2. 创建对应的 handler 函数
```go
// 在 internal/handlers/ 目录下创建处理函数
func Register(ctx context.Context, c *app.RequestContext) {
    // 实现注册逻辑
}
```

### 3. 添加必要的中间件
```go
// 需要认证的路由添加认证中间件
auth.POST("/logout", middleware.AuthMiddleware(jwtSecret), handlers.Logout)
```

## 🏗️ 下一步开发计划

1. **创建 handlers 包**: 实现具体的路由处理逻辑
2. **完善认证路由**: 实现注册、登录、登出功能
3. **实现用户管理**: 用户信息的增删改查
4. **开发聊天功能**: WebSocket 实时通信
5. **添加中间件**: 认证、限流、日志等中间件

## 🔍 中间件使用

当前已集成的中间件：
- `recovery.Recovery()` - 异常恢复
- CORS 中间件 - 跨域处理

计划添加的中间件：
- `AuthMiddleware()` - JWT 认证
- `RateLimitMiddleware()` - 请求限流
- `ValidatorMiddleware()` - 参数验证

## 📚 技术栈

- **框架**: CloudWeGo Hertz - 字节跳动开源的高性能 HTTP 框架
- **数据库**: MySQL + GORM
- **认证**: JWT Token
- **日志**: Logrus
- **实时通信**: WebSocket (待实现)

## 🔄 从 Gin 迁移

本项目已完成从 Gin 到 Hertz 的迁移，主要改动：

1. **处理函数签名**:
   - Gin: `func(c *gin.Context)`
   - Hertz: `func(ctx context.Context, c *app.RequestContext)`

2. **框架初始化**:
   - Gin: `gin.New()` → `r.Run(":8080")`
   - Hertz: `server.Default(server.WithHostPorts(":8080"))` → `h.Spin()`

3. **响应方法**:
   - Gin: `c.JSON(200, gin.H{...})`
   - Hertz: `c.JSON(consts.StatusOK, utils.H{...})`

## 📚 参考资源

- [Hertz 官方文档](https://www.cloudwego.io/zh/docs/hertz/)
- [Hertz 迁移指南](https://www.cloudwego.io/zh/docs/hertz/tutorials/service-migration/)
- [Go 项目结构最佳实践](https://github.com/golang-standards/project-layout)
- [RESTful API 设计规范](https://restfulapi.net/) 