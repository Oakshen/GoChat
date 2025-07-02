# GoChat - Go 语言实时聊天系统（仿微信）

一个基于 Go 语言开发的现代化实时聊天应用，支持单聊群聊、多种消息类型、音视频通话和后台管理功能。

## 📊 项目进度状态

### 🎯 整体完成度：30% 

#### ✅ 已完成功能

##### 1. 核心架构（100% ✅）
- ✅ 完整的分层架构设计（DDD领域驱动）
- ✅ CloudWeGo Hertz 高性能框架集成
- ✅ 配置管理系统（YAML + 环境变量）
- ✅ MySQL 数据库连接和自动迁移
- ✅ GORM ORM 集成

##### 2. 数据模型（100% ✅）
- ✅ User（用户）实体模型
- ✅ Room（聊天室）实体模型  
- ✅ Message（消息）实体模型
- ✅ RoomMember（成员关系）实体模型
- ✅ 完整的数据库关联关系
- ✅ MySQL 建表脚本

##### 3. 用户认证系统（100% ✅）
- ✅ 用户注册功能（用户名/邮箱验证）
- ✅ 用户登录功能（支持用户名/邮箱登录）
- ✅ JWT Token 生成和验证
- ✅ 认证中间件
- ✅ 密码安全加密（bcrypt）

##### 4. 用户管理模块（90% ✅）
- ✅ 用户信息 CRUD 操作
- ✅ 用户资料查询和更新
- ✅ 在线用户状态管理
- ✅ 用户删除功能
- 🔄 用户头像上传（待实现）

##### 5. HTTP API 接口（80% ✅）
- ✅ RESTful API 设计
- ✅ 统一响应格式
- ✅ 接口参数验证
- ✅ 错误处理机制
- 🔄 API 文档生成（待完善）

#### ❌ 待实现核心功能（70% 未完成）

##### 1. WebSocket 实时通信模块（0% ❌）- **最关键**
- ❌ WebSocket 连接管理器
- ❌ 客户端连接池
- ❌ 消息路由和广播
- ❌ 心跳检测机制
- ❌ 断线重连处理

##### 2. 聊天功能（0% ❌）
- ❌ 创建/删除聊天室
- ❌ 加入/退出聊天室
- ❌ 单聊/群聊逻辑
- ❌ 消息发送/接收
- ❌ 消息历史查询
- ❌ 未读消息统计

##### 3. 多媒体消息（0% ❌）
- ❌ 文本消息处理
- ❌ 图片上传和发送
- ❌ 文件上传和发送
- ❌ 语音消息录制和播放
- ❌ 视频消息支持

##### 4. 前端界面（0% ❌）
- ❌ 用户登录/注册页面
- ❌ 聊天主界面
- ❌ 联系人列表
- ❌ 聊天室列表
- ❌ 消息展示组件
- ❌ 文件上传组件

##### 5. 音视频通话（0% ❌）- **高级功能**
- ❌ WebRTC 集成
- ❌ 语音通话功能
- ❌ 视频通话功能
- ❌ 通话录制
- ❌ 屏幕共享

##### 6. 后台管理系统（0% ❌）
- ❌ 管理员登录
- ❌ 用户管理界面
- ❌ 聊天室管理
- ❌ 系统监控
- ❌ 数据统计

## 🚀 详细开发路线图

### 第一阶段：核心实时通信（3-4周）🔥

#### Week 1: WebSocket 基础架构
```
优先级：极高 🔥🔥🔥
```
- [ ] 集成 Gorilla WebSocket 库
- [ ] 实现 WebSocket 连接管理器
- [ ] 创建客户端连接池
- [ ] 实现心跳检测机制
- [ ] 添加断线重连逻辑

#### Week 2: 聊天室系统
```
优先级：极高 🔥🔥
```
- [ ] 聊天室创建/删除 API
- [ ] 用户加入/退出聊天室
- [ ] 聊天室成员管理
- [ ] 聊天室权限控制
- [ ] 私聊和群聊逻辑分离

#### Week 3: 消息系统
```
优先级：极高 🔥🔥
```
- [ ] 实时消息发送/接收
- [ ] 消息持久化存储
- [ ] 消息历史查询 API
- [ ] 未读消息统计
- [ ] 消息状态（已发送/已读）

#### Week 4: 基础前端界面
```
优先级：高 🔥
```
- [ ] HTML/CSS/JavaScript 聊天界面
- [ ] WebSocket 客户端连接
- [ ] 实时消息展示
- [ ] 基础聊天功能测试

### 第二阶段：多媒体消息（2-3周）

#### Week 5-6: 文件上传系统
```
优先级：高 🔥
```
- [ ] 文件上传接口
- [ ] 图片预览功能
- [ ] 文件下载接口
- [ ] 文件类型验证
- [ ] 文件大小限制

#### Week 7: 多媒体消息
```
优先级：中 📤
```
- [ ] 图片消息发送
- [ ] 文件消息发送
- [ ] 语音消息录制
- [ ] 语音消息播放

### 第三阶段：高级功能（4-5周）

#### Week 8-9: 音视频通话
```
优先级：中 📤
```
- [ ] WebRTC 集成
- [ ] 一对一语音通话
- [ ] 一对一视频通话
- [ ] 通话质量优化

#### Week 10-11: 后台管理系统
```
优先级：中 📤
```
- [ ] 管理员认证系统
- [ ] 用户管理界面
- [ ] 聊天记录管理
- [ ] 系统监控面板

#### Week 12: 系统优化
```
优先级：低 📋
```
- [ ] Redis 缓存集成
- [ ] 数据库查询优化
- [ ] 并发性能优化
- [ ] 部署和运维工具

## 🛠️ 技术栈

### 后端技术
- **框架**: CloudWeGo Hertz（高性能 HTTP 框架）
- **数据库**: MySQL 8.0 + GORM v2
- **实时通信**: Gorilla WebSocket
- **认证**: JWT Token
- **配置**: Viper（YAML + 环境变量）
- **日志**: Logrus
- **缓存**: Redis（计划中）

### 前端技术
- **基础**: HTML5 + CSS3 + JavaScript ES6+
- **WebSocket**: 原生 WebSocket API
- **多媒体**: MediaRecorder API、WebRTC
- **UI框架**: 待定（可选 Vue.js/React）

### 音视频技术
- **WebRTC**: 点对点通信
- **STUN/TURN**: NAT 穿透
- **媒体编码**: VP8/VP9、H.264

## 🚀 快速开始

### 1. 环境准备
```bash
# 确保已安装 Go 1.23.2+
go version

# 确保 MySQL 数据库运行
mysql --version

# 创建数据库
mysql -u root -p -e "CREATE DATABASE gochat CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
```

### 2. 项目配置
```bash
# 克隆项目
git clone <your-repo-url>
cd GoChat

# 复制配置文件
cp configs/config.yaml.example configs/config.yaml

# 编辑配置文件，设置数据库连接等信息
vim configs/config.yaml
```

### 3. 安装依赖并运行
```bash
# 安装依赖
go mod tidy

# 执行数据库迁移
go run cmd/server/main.go migrate

# 启动服务器
go run cmd/server/main.go
# 或使用脚本
chmod +x sh/run.sh
./sh/run.sh
```

### 4. 测试接口
```bash
# 健康检查
curl http://localhost:8080/health

# 用户注册
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@example.com","password":"123456"}'

# 用户登录
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"123456"}'
```

## 📋 近期开发重点

### 🔥 本周任务（优先级：极高）
1. **WebSocket 模块搭建**
   - 集成 Gorilla WebSocket 库
   - 实现连接管理器
   - 创建消息路由系统

2. **聊天室基础功能**
   - 聊天室创建 API
   - 用户加入聊天室逻辑
   - 基础消息发送测试

### 📤 下周计划（优先级：高）
1. **完善消息系统**
   - 消息持久化
   - 历史消息查询
   - 未读消息统计

2. **前端原型开发**
   - 基础聊天界面
   - WebSocket 客户端

## 📁 项目结构

```
GoChat/
├── cmd/
│   └── server/main.go          # 应用入口 ✅
├── internal/
│   ├── config/                 # 配置管理 ✅
│   ├── database/               # 数据库连接 ✅
│   ├── models/                 # 数据模型 ✅
│   │   ├── entities/           # 实体定义 ✅
│   │   ├── requests/           # 请求结构 ✅
│   │   └── responses/          # 响应结构 ✅
│   ├── handlers/               # HTTP 处理器 ✅
│   ├── services/               # 业务逻辑层 ✅
│   ├── dal/                    # 数据访问层 ✅
│   ├── middleware/             # 中间件 ✅
│   ├── router/                 # 路由配置 ✅
│   ├── websocket/              # WebSocket 模块 ❌ 待实现
│   └── chat/                   # 聊天业务逻辑 ❌ 待实现
├── pkg/
│   ├── logger/                 # 日志工具 ✅
│   ├── response/               # 响应工具 ✅
│   └── utils/                  # 工具函数 ✅
├── web/                        # 前端资源 ❌ 待实现
│   ├── static/                 # 静态文件
│   ├── templates/              # HTML 模板
│   └── assets/                 # 资源文件
├── configs/                    # 配置文件 ✅
├── sql/                        # SQL 脚本 ✅
├── docs/                       # 文档 ✅
└── sh/                         # 脚本文件 ✅
```

## 📖 API 文档

### 🟢 已实现接口

#### 系统接口
- `GET /health` - 健康检查

#### 认证接口
- `POST /api/auth/register` - 用户注册
- `POST /api/auth/login` - 用户登录
- `GET /api/auth/userinfo` - 获取当前用户信息

#### 用户接口
- `GET /api/users/:id` - 获取用户资料
- `PUT /api/users/:id` - 更新用户资料
- `GET /api/users/online` - 获取在线用户列表
- `GET /api/users/delete/:id` - 删除用户

### 🔄 计划中接口

#### 聊天室接口
- `POST /api/rooms` - 创建聊天室
- `GET /api/rooms` - 获取聊天室列表
- `GET /api/rooms/:id` - 获取聊天室详情
- `POST /api/rooms/:id/join` - 加入聊天室
- `POST /api/rooms/:id/leave` - 退出聊天室

#### 消息接口
- `GET /api/rooms/:id/messages` - 获取聊天记录
- `POST /api/messages/:id/read` - 标记消息已读
- `GET /api/messages/unread` - 获取未读消息

#### 文件接口
- `POST /api/upload/image` - 上传图片
- `POST /api/upload/file` - 上传文件
- `GET /api/files/:id` - 下载文件

#### WebSocket 接口
- `WebSocket /ws` - 建立 WebSocket 连接

## 🔧 开发指南

### WebSocket 消息协议设计
```json
{
  "type": "message|join|leave|typing",
  "room_id": "聊天室ID", 
  "user_id": "用户ID",
  "content": "消息内容",
  "timestamp": "时间戳"
}
```

### 添加新功能步骤
1. 在 `internal/models/entities/` 定义数据模型
2. 在 `internal/dal/` 实现数据访问层
3. 在 `internal/services/` 实现业务逻辑
4. 在 `internal/handlers/` 实现 HTTP 处理器
5. 在 `internal/router/router.go` 注册路由

## 🤝 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/new-feature`)
3. 提交更改 (`git commit -am 'Add new feature'`)
4. 推送到分支 (`git push origin feature/new-feature`)
5. 创建 Pull Request

## 📄 许可证

MIT License

---

## 🎯 开发里程碑

- [x] **里程碑 1**: 项目架构搭建 (已完成)
- [x] **里程碑 2**: 用户认证系统 (已完成)  
- [ ] **里程碑 3**: WebSocket 实时通信 (进行中)
- [ ] **里程碑 4**: 基础聊天功能
- [ ] **里程碑 5**: 多媒体消息
- [ ] **里程碑 6**: 音视频通话
- [ ] **里程碑 7**: 后台管理系统
- [ ] **里程碑 8**: 性能优化和部署

**当前重点：实现 WebSocket 实时通信模块，这是聊天系统的核心功能！** 