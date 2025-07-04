# GoChat - Go 语言实时聊天系统

一个基于 Go 语言和 CloudWeGo Hertz 框架开发的现代化实时聊天应用，配备完整的 Vue.js 前端界面，支持多用户实时聊天、聊天室管理、消息持久化等功能。

![Example](docs/images/img1.png)
![Example](docs/images/img2.png)

## 🎯 开发里程碑

- [x] **里程碑 1**: 项目架构搭建 (已完成) ✅
- [x] **里程碑 2**: 用户认证系统 (已完成) ✅
- [x] **里程碑 3**: WebSocket 实时通信 (已完成) ✅
- [x] **里程碑 4**: 基础聊天功能 (已完成) ✅
- [x] **里程碑 5**: 多用户实时聊天 (已完成) ✅
- [x] **里程碑 6**: Vue.js 前端应用 (已完成) ✅ **NEW!**
- [ ] **里程碑 7**: 多媒体消息 (计划中) 📋
- [ ] **里程碑 8**: 音视频通话 (长期计划) 📋
- [ ] **里程碑 9**: 后台管理系统 (长期计划) 📋
- [ ] **里程碑 10**: 性能优化和部署 (长期计划) 📋

## 📊 项目进度状态

### 🎯 整体完成度：80% ⬆️ **大幅提升!**

#### ✅ 已完成功能

##### 1. 核心架构（100% ✅）
- ✅ 完整的分层架构设计（Models → DAL → Services → Handlers）
- ✅ CloudWeGo Hertz 高性能框架集成
- ✅ 配置管理系统（YAML + 环境变量）
- ✅ MySQL 数据库连接和自动迁移
- ✅ GORM ORM 集成
- ✅ 统一错误处理和响应格式

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
- ✅ WebSocket 连接 JWT 认证

##### 4. WebSocket 实时通信模块（100% ✅）- **核心功能已完成**
- ✅ WebSocket 连接管理器（Hub）
- ✅ 客户端连接池管理
- ✅ 消息路由和广播系统
- ✅ **心跳检测机制（Ping/Pong）** 🆕 **最新修复！**
- ✅ 异步消息处理
- ✅ 连接断开自动清理
- ✅ 多用户并发连接支持
- ✅ **WebSocket错误处理优化** 🆕

##### 5. 聊天功能（95% ✅）⬆️ **提升！**
- ✅ 创建/删除聊天室
- ✅ 加入/退出聊天室
- ✅ **智能加入聊天室（支持ID和名称搜索）** 🆕 **新功能！**
- ✅ **防重复加入逻辑** 🆕 **问题修复！**
- ✅ 单聊/群聊逻辑
- ✅ 实时消息发送/接收
- ✅ 消息持久化存储
- ✅ 用户在线状态管理
- ✅ 聊天室成员管理
- ✅ 系统消息（用户加入/离开提醒）
- ✅ 消息历史查询
- 🔄 未读消息统计（待优化）

##### 6. Vue.js 前端应用（95% ✅）🆕 **重大更新！**
- ✅ **完整的Vue 3 + Vite项目架构**
- ✅ **Element Plus UI组件库集成**
- ✅ **Pinia状态管理**
- ✅ **Vue Router路由管理**
- ✅ **用户认证界面（登录/注册）**
- ✅ **现代化聊天界面（仿微信设计）**
- ✅ **实时WebSocket连接**
- ✅ **聊天室管理（创建/加入/切换）**
- ✅ **消息发送接收界面**
- ✅ **在线用户列表**
- ✅ **响应式设计（移动端适配）**
- ✅ **错误处理和用户提示**
- ✅ **API请求封装**
- 🔄 UI细节优化（待完善）

##### 7. HTTP API 接口（95% ✅）⬆️
- ✅ RESTful API 设计
- ✅ 统一响应格式
- ✅ 接口参数验证
- ✅ 错误处理机制
- ✅ 聊天室管理 API
- ✅ **幂等性保证（重复操作处理）** 🆕
- 🔄 API 文档生成（待完善）

#### 🔄 正在优化的功能

##### 1. 消息系统增强（10% 待完成）⬇️
- 🔄 消息状态（已发送/已读）
- 🔄 消息搜索功能
- 🔄 消息撤回功能
- 🔄 @提及功能

##### 2. 用户体验优化（10% 待完成）⬇️
- 🔄 用户头像上传
- 🔄 用户在线状态显示优化
- 🔄 打字状态提示
- 🔄 消息通知

#### ❌ 待实现高级功能（20% 未完成）⬇️

##### 1. 多媒体消息（0% ❌）
- ❌ 图片上传和发送
- ❌ 文件上传和发送
- ❌ 语音消息录制和播放
- ❌ 视频消息支持

##### 2. 音视频通话（0% ❌）- **高级功能**
- ❌ WebRTC 集成
- ❌ 语音通话功能
- ❌ 视频通话功能
- ❌ 通话录制
- ❌ 屏幕共享

##### 3. 后台管理系统（0% ❌）
- ❌ 管理员登录
- ❌ 用户管理界面
- ❌ 聊天室管理
- ❌ 系统监控
- ❌ 数据统计

##### 4. 性能优化（0% ❌）
- ❌ Redis 缓存集成
- ❌ 消息队列
- ❌ 数据库查询优化
- ❌ 负载均衡

## 🆕 最新更新记录

### v0.8.0 (最新版本) - 2024-01-01
#### 🎉 重大更新
- ✅ **完整Vue.js前端应用上线**
- ✅ **现代化聊天界面（三栏布局）**
- ✅ **智能聊天室加入功能**
  - 支持ID和名称搜索
  - 实时搜索结果展示
  - 防重复加入逻辑

#### 🐛 问题修复
- ✅ **WebSocket心跳问题修复**
  - 修复`{"type":"ping"}`消息"不支持的消息类型"错误
  - 实现完整的ping/pong心跳机制
- ✅ **重复加入聊天室问题修复**
  - 修复`{"error":"duplicated key not allowed"}`错误
  - 前端智能检测用户是否已在聊天室
  - 后端幂等性处理

#### 🔧 技术改进
- ✅ 前端状态管理优化（Pinia）
- ✅ WebSocket连接稳定性提升
- ✅ API请求错误处理完善
- ✅ 用户体验优化

## 🚀 快速开始

### 1. 环境准备
```bash
# 确保已安装 Go 1.23.2+
go version

# 确保已安装 Node.js 18.0+
node --version
npm --version

# 确保 MySQL 数据库运行
mysql --version

# 创建数据库
mysql -u root -p -e "CREATE DATABASE gochat CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
```

### 2. 后端启动
```bash
# 克隆项目
git clone <your-repo-url>
cd GoChat

# 安装Go依赖
go mod tidy

# 编辑配置文件，设置数据库连接等信息
# 配置文件位置: configs/config.yaml

# 启动后端服务器
go run cmd/server/main.go
# 或使用脚本
chmod +x sh/run.sh && ./sh/run.sh
```

### 3. 前端启动 🆕
```bash
# 进入前端目录
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 或者构建生产版本
npm run build
```

### 4. 访问应用
- **前端界面**: http://localhost:3000
- **后端API**: http://localhost:8080
- **WebSocket**: ws://localhost:3000/ws

### 5. 测试聊天功能 🆕

#### 推荐方式：Vue.js前端界面
1. 打开浏览器访问：`http://localhost:3000`
2. 注册两个测试用户（如：user1, user2）
3. 分别在不同浏览器窗口登录
4. 创建聊天室或加入已有聊天室
5. 开始实时聊天！

#### 备用方式：原生Web界面
1. 访问：`http://localhost:8080/web/index.html`
2. 进行基础聊天测试

#### API测试脚本
```bash
# 运行多用户聊天测试脚本
chmod +x sh/test_multi_user_chat.sh
./sh/test_multi_user_chat.sh
```

## 🎮 详细使用说明

### Vue.js前端功能 🆕

#### 1. 用户认证
- **注册**: 用户名、邮箱、密码注册
- **登录**: 支持用户名或邮箱登录
- **自动登录**: JWT Token本地存储，刷新页面保持登录状态

#### 2. 聊天室管理
- **创建聊天室**: 设置名称、描述、公开/私密
- **加入聊天室**: 通过ID或名称搜索加入
- **聊天室列表**: 显示已加入的聊天室
- **实时成员**: 显示在线成员数量

#### 3. 实时聊天
- **消息发送**: 支持文本消息发送
- **实时接收**: WebSocket实时消息推送
- **消息历史**: 自动加载聊天记录
- **系统消息**: 用户加入/离开提醒

#### 4. 界面特性
- **响应式设计**: 适配桌面和移动端
- **现代化UI**: Element Plus组件库
- **三栏布局**: 仿微信聊天界面
- **操作流畅**: 无刷新页面切换

### WebSocket 连接 🔧
```javascript
// Vue.js前端自动处理WebSocket连接
// 手动连接示例：
const ws = new WebSocket(`ws://localhost:3000/ws?token=${jwt_token}`);

// 加入聊天室
ws.send(JSON.stringify({
    type: 'join',
    room_id: 1
}));

// 发送消息
ws.send(JSON.stringify({
    type: 'text',
    room_id: 1,
    content: '你好，世界！'
}));
```

### WebSocket 消息协议
```json
{
  "type": "text|join|leave|typing|system|userlist|error|ping|pong",
  "room_id": 1,
  "user_id": 2,
  "username": "user1",
  "content": "消息内容",
  "timestamp": "2025-07-03T15:30:00Z",
  "message_id": 123
}
```

### 支持的消息类型
- `text`: 文本消息
- `join`: 加入聊天室
- `leave`: 离开聊天室
- `typing`: 正在输入状态
- `system`: 系统消息
- `userlist`: 用户列表更新
- `error`: 错误消息
- `ping`: 心跳检测 🆕
- `pong`: 心跳响应 🆕

## 🛠️ 技术栈

### 后端技术
- **框架**: CloudWeGo Hertz v0.8.0（高性能 HTTP 框架）
- **WebSocket**: hertz-contrib/websocket（官方WebSocket支持）
- **数据库**: MySQL 8.0 + GORM v2
- **认证**: JWT Token
- **配置**: Viper（YAML + 环境变量）
- **日志**: 自定义日志系统
- **密码加密**: bcrypt

### 前端技术 🆕
- **框架**: Vue.js 3.4+ (Composition API)
- **构建工具**: Vite 5.0+
- **UI组件**: Element Plus
- **状态管理**: Pinia
- **路由**: Vue Router 4
- **HTTP客户端**: Axios
- **WebSocket**: 原生 WebSocket API + 封装
- **样式**: CSS3 + 响应式设计

### 架构设计
- **分层架构**: Models → DAL → Services → Handlers
- **前后端分离**: Vue.js SPA + Go API
- **并发安全**: Goroutine + Channel + Mutex
- **消息处理**: 异步广播机制
- **连接管理**: Hub模式连接池
- **状态管理**: Vuex/Pinia响应式状态

## 🚀 详细开发路线图

### ✅ 第一阶段：核心实时通信（已完成）

#### Week 1-2: WebSocket 基础架构 ✅
- ✅ 集成 hertz-contrib/websocket 库
- ✅ 实现 WebSocket 连接管理器（Hub）
- ✅ 创建客户端连接池
- ✅ 实现心跳检测机制
- ✅ 添加异步消息处理

#### Week 3: 聊天室系统 ✅
- ✅ 聊天室创建/删除 API
- ✅ 用户加入/退出聊天室
- ✅ 聊天室成员管理
- ✅ 聊天室权限控制
- ✅ 私聊和群聊逻辑

#### Week 4: 消息系统 ✅
- ✅ 实时消息发送/接收
- ✅ 消息持久化存储
- ✅ 消息历史查询 API
- ✅ 多用户并发聊天支持
- ✅ 系统消息推送

#### Week 5: 前端界面 ✅
- ✅ HTML/CSS/JavaScript 聊天界面
- ✅ WebSocket 客户端连接
- ✅ 实时消息展示
- ✅ 完整聊天功能测试

### ✅ 第二阶段：Vue.js前端应用（已完成）🆕

#### Week 6-7: Vue.js项目搭建 ✅
- ✅ Vue 3 + Vite 项目初始化
- ✅ Element Plus UI组件库集成
- ✅ Pinia状态管理配置
- ✅ Vue Router路由配置
- ✅ Axios HTTP客户端配置

#### Week 8: 前端核心功能 ✅
- ✅ 用户认证页面（登录/注册）
- ✅ 聊天主界面设计
- ✅ WebSocket客户端封装
- ✅ 状态管理实现

#### Week 9: 聊天功能实现 ✅
- ✅ 聊天室列表和管理
- ✅ 实时消息收发
- ✅ 加入聊天室功能
- ✅ 用户界面优化

#### Week 10: 问题修复和优化 ✅
- ✅ WebSocket心跳机制修复
- ✅ 重复加入聊天室问题修复
- ✅ 前端错误处理完善
- ✅ 用户体验优化

### 🔄 第三阶段：功能增强（进行中）

#### Week 11-12: 消息系统增强 🔄
- 🔄 消息状态（已发送/已读）
- 🔄 未读消息统计优化
- 🔄 消息搜索功能
- 🔄 @提及功能

#### Week 13: 用户体验优化 📋
- [ ] 用户头像上传
- [ ] 打字状态提示
- [ ] 消息通知
- [ ] UI/UX 改进

### 📤 第四阶段：多媒体消息（计划中）

#### Week 14-15: 文件上传系统 📋
- [ ] 文件上传接口
- [ ] 图片预览功能
- [ ] 文件下载接口
- [ ] 文件类型验证

#### Week 16: 多媒体消息 📋
- [ ] 图片消息发送
- [ ] 文件消息发送
- [ ] 语音消息录制
- [ ] 语音消息播放

### 📋 第五阶段：高级功能（长期计划）

#### 音视频通话 📋
- [ ] WebRTC 集成
- [ ] 一对一语音通话
- [ ] 一对一视频通话
- [ ] 通话质量优化

#### 后台管理系统 📋
- [ ] 管理员认证系统
- [ ] 用户管理界面
- [ ] 聊天记录管理
- [ ] 系统监控面板

#### 系统优化 📋
- [ ] Redis 缓存集成
- [ ] 数据库查询优化
- [ ] 并发性能优化
- [ ] 部署和运维工具

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
│   │   ├── auth_handler.go     # 认证处理器 ✅
│   │   ├── room_handler.go     # 聊天室处理器 ✅
│   │   ├── user_handler.go     # 用户处理器 ✅
│   │   └── websocket_handler.go # WebSocket处理器 ✅
│   ├── services/               # 业务逻辑层 ✅
│   │   ├── auth_service.go     # 认证服务 ✅
│   │   ├── message_service.go  # 消息服务 ✅
│   │   ├── room_service.go     # 聊天室服务 ✅
│   │   └── user_service.go     # 用户服务 ✅
│   ├── dal/                    # 数据访问层 ✅
│   │   ├── message_dal.go      # 消息数据访问 ✅
│   │   └── user_dal.go         # 用户数据访问 ✅
│   ├── websocket/              # WebSocket 模块 ✅
│   │   ├── client.go          # 客户端管理 ✅
│   │   ├── hub.go             # 连接管理中心 ✅
│   │   └── message.go         # 消息协议 ✅
│   ├── middleware/             # 中间件 ✅
│   └── router/                 # 路由配置 ✅
├── pkg/
│   ├── logger/                 # 日志工具 ✅
│   ├── response/               # 响应工具 ✅
│   └── utils/                  # 工具函数 ✅
├── frontend/                   # Vue.js前端应用 ✅ 🆕
│   ├── src/
│   │   ├── components/        # Vue组件 ✅
│   │   ├── views/             # 页面视图 ✅
│   │   ├── stores/            # Pinia状态管理 ✅
│   │   ├── api/               # API接口封装 ✅
│   │   ├── utils/             # 工具函数 ✅
│   │   └── router/            # 路由配置 ✅
│   ├── public/                # 静态资源 ✅
│   ├── package.json           # 项目配置 ✅
│   ├── vite.config.js         # Vite配置 ✅
│   └── README.md              # 前端说明 ✅
├── web/                        # 原生前端资源 ✅
│   └── index.html             # 原生聊天测试界面 ✅
├── configs/                    # 配置文件 ✅
├── sql/                        # SQL 脚本 ✅
├── docs/                       # 文档 ✅
└── sh/                         # 脚本文件 ✅
    ├── run.sh                 # 启动脚本 ✅
    ├── test_api.sh            # API测试脚本 ✅
    └── test_multi_user_chat.sh # 多用户聊天测试 ✅
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

#### 聊天室接口
- `POST /api/rooms` - 创建聊天室
- `GET /api/rooms` - 获取聊天室列表
- `GET /api/rooms/:id` - 获取聊天室详情
- `POST /api/rooms/:id/join` - 加入聊天室 🔧 **已优化幂等性**
- `POST /api/rooms/:id/leave` - 退出聊天室
- `GET /api/rooms/:id/members` - 获取聊天室成员
- `GET /api/rooms/:id/messages` - 获取聊天记录

#### WebSocket 接口
- `WebSocket /ws` - 建立 WebSocket 连接（需JWT认证）
- `GET /api/ws/stats` - 获取WebSocket连接统计
- `POST /api/ws/broadcast/:roomId` - 向聊天室广播消息

### 🔄 计划中接口

#### 消息接口
- `POST /api/messages/:id/read` - 标记消息已读
- `GET /api/messages/unread` - 获取未读消息

#### 文件接口
- `POST /api/upload/image` - 上传图片
- `POST /api/upload/file` - 上传文件
- `GET /api/files/:id` - 下载文件

## 🔧 开发指南

### 后端开发
#### 添加新功能步骤
1. 在 `internal/models/entities/` 定义数据模型
2. 在 `internal/dal/` 实现数据访问层
3. 在 `internal/services/` 实现业务逻辑
4. 在 `internal/handlers/` 实现 HTTP 处理器
5. 在 `internal/router/router.go` 注册路由

#### WebSocket 开发指南
1. 消息类型在 `internal/websocket/message.go` 中定义
2. 在 `internal/websocket/hub.go` 中添加处理逻辑
3. 客户端连接管理在 `internal/websocket/client.go`

### 前端开发 🆕
#### Vue.js开发规范
1. 组件开发在 `frontend/src/components/`
2. 页面视图在 `frontend/src/views/`
3. 状态管理在 `frontend/src/stores/`
4. API接口在 `frontend/src/api/`

#### 开发命令
```bash
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 构建生产版本
npm run build

# 代码检查
npm run lint
```

## 🎯 贡献指南

### 开发环境设置
1. Fork 项目到您的GitHub账户
2. 克隆项目到本地
3. 设置开发环境（Go + Node.js + MySQL）
4. 创建功能分支进行开发
5. 提交Pull Request

### 代码规范
- **Go代码**: 遵循 Go 标准规范，使用 `gofmt` 格式化
- **Vue.js代码**: 遵循 Vue.js 官方风格指南
- **提交信息**: 使用清晰的提交信息，包含功能描述

### 问题报告
如果您发现bug或有功能建议，请：
1. 在GitHub Issues中创建新issue
2. 详细描述问题或需求
3. 提供复现步骤（如果是bug）

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 📧 联系方式

- 项目维护者: [您的姓名]
- 邮箱: [您的邮箱]
- 项目地址: [GitHub项目链接]

---

⭐ 如果这个项目对您有帮助，请给个Star支持一下！

📈 **项目状态**: 积极开发中，欢迎贡献代码和提出建议！ 