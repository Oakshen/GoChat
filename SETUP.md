# GoChat 快速配置指南

## 🚀 在本机运行 GoChat 的步骤

### 1. 准备 MySQL 数据库

首先确保 MySQL 已安装并运行：

```bash
# 检查 MySQL 是否运行
brew services list | grep mysql
# 或者
sudo systemctl status mysql

# 如果没有运行，启动 MySQL
brew services start mysql
# 或者
sudo systemctl start mysql
```

### 2. 创建数据库

连接到 MySQL 并创建数据库：

```bash
# 连接到 MySQL
mysql -u root -p

# 在 MySQL 命令行中执行：
CREATE DATABASE gochat CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
EXIT;
```

### 3. 配置环境变量

有两种方式配置：

#### 方式 1：设置环境变量
```bash
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=你的MySQL密码
export DB_NAME=gochat
export JWT_SECRET=gochat-jwt-secret-key-change-in-production-2024
export JWT_EXPIRE_HOURS=24
export SERVER_PORT=8080
export SERVER_MODE=debug
```

#### 方式 2：修改配置文件
编辑 `configs/config.yaml` 文件，修改以下内容：

```yaml
database:
  host: "localhost"
  port: 3306
  user: "root"
  password: "你的MySQL密码"  # 修改这里
  dbname: "gochat"
  timezone: "Asia/Shanghai"

jwt:
  secret: "gochat-jwt-secret-key-change-in-production-2024"
  expire_hours: 24
```

### 4. 安装依赖并运行

```bash
# 确保在项目根目录
cd /path/to/GoChat

# 安装依赖
go mod tidy

# 编译检查
go build cmd/server/main.go

# 运行服务器
go run cmd/server/main.go
```

### 5. 测试服务

如果一切正常，你应该看到类似这样的输出：
```
{"level":"info","msg":"Starting GoChat server...","time":"2024-01-01T10:00:00+08:00"}
{"level":"info","msg":"Database connected successfully","time":"2024-01-01T10:00:01+08:00"}
{"level":"info","msg":"Database migration completed","time":"2024-01-01T10:00:01+08:00"}
{"level":"info","msg":"Server starting on :8080","time":"2024-01-01T10:00:01+08:00"}
```

然后可以测试接口：

```bash
# 健康检查
curl http://localhost:8080/health

# 应该返回：
# {"message":"GoChat server is running","status":"ok"}

# API 测试
curl http://localhost:8080/api/ping

# 应该返回：
# {"message":"pong"}
```

## 🔧 JWT 配置说明

### JWT Secret 设置
- **开发环境**: 可以使用配置文件中的默认值
- **生产环境**: 必须使用强密码，建议使用环境变量

### 生成安全的 JWT Secret
```bash
# 使用 openssl 生成随机密钥
openssl rand -base64 32

# 或者使用 Go 生成
go run -c 'package main; import ("crypto/rand"; "encoding/base64"; "fmt"); func main() { bytes := make([]byte, 32); rand.Read(bytes); fmt.Println(base64.StdEncoding.EncodeToString(bytes)) }'
```

## 🐛 常见问题

### 1. 数据库连接失败
- 检查 MySQL 是否运行
- 检查用户名密码是否正确
- 检查数据库 `gochat` 是否存在

### 2. 编译错误
```bash
# 清理模块缓存
go clean -modcache
go mod tidy
```

### 3. 端口被占用
```bash
# 查看端口占用
lsof -i :8080

# 修改端口
export SERVER_PORT=8081
```

### 4. JWT Token 测试

启动服务后，可以通过以下方式测试 JWT：

```bash
# 注册用户（需要先实现注册接口）
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"123456"}'

# 登录获取 token（需要先实现登录接口）
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"123456"}'
```

## ✅ 下一步

如果服务成功启动，接下来可以：
1. 实现认证相关的 HTTP 接口
2. 添加 WebSocket 功能
3. 创建前端界面

有问题随时查看日志输出进行调试！ 