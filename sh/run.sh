#!/bin/bash

echo "🚀 GoChat 分层架构重构版本启动"
echo "==============================================="

# 检查配置文件
if [ ! -f "configs/config.yaml" ]; then
    echo "❌ 配置文件不存在: configs/config.yaml"
    exit 1
fi

# 整理依赖
echo "📦 整理依赖包..."
go mod tidy

# 构建项目
echo "🔨 构建项目..."
go build -o main cmd/server/main.go

if [ $? -ne 0 ]; then
    echo "❌ 构建失败"
    exit 1
fi

echo "✅ 构建成功!"

# 运行服务
echo "🌟 启动GoChat服务器..."
echo "API接口:"
echo "  - 健康检查: GET http://localhost:8080/health"
echo "  - 用户注册: POST http://localhost:8080/api/auth/register"
echo "  - 用户登录: POST http://localhost:8080/api/auth/login"
echo "==============================================="

./main 