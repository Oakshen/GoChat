#!/bin/bash

# GoChat 连接测试脚本

echo "🔍 GoChat 连接测试"
echo "=================="

# 测试后端API连接
echo "📡 测试后端API连接 (localhost:8080)..."
if curl -s --connect-timeout 3 http://localhost:8080/health > /dev/null; then
    echo "✅ 后端API连接成功"
    curl -s http://localhost:8080/health | python3 -m json.tool 2>/dev/null || echo "API响应正常"
else
    echo "❌ 后端API连接失败"
    echo "💡 请确保后端服务正在运行在 http://localhost:8080"
    echo "💡 可以运行以下命令启动后端："
    echo "   cd /path/to/gochat && go run cmd/server/main.go"
    echo ""
fi

echo ""

# 检查端口占用
echo "🔍 检查端口占用情况..."
echo "端口 8080 (后端):"
if lsof -i :8080 > /dev/null 2>&1; then
    echo "✅ 端口 8080 已被占用 (后端应该在运行)"
    lsof -i :8080 | head -2
else
    echo "❌ 端口 8080 未被占用 (后端可能未启动)"
fi

echo ""
echo "端口 3000 (前端):"
if lsof -i :3000 > /dev/null 2>&1; then
    echo "⚠️  端口 3000 已被占用"
    lsof -i :3000 | head -2
    echo "💡 如果前端未启动，请先停止占用3000端口的进程"
else
    echo "✅ 端口 3000 可用"
fi

echo ""
echo "🚀 启动建议："
echo "1. 确保后端服务运行在 http://localhost:8080"
echo "2. 前端将运行在 http://localhost:3000"
echo "3. Vite代理会自动转发API请求到后端"
echo ""
echo "📝 启动顺序："
echo "1. 先启动后端: go run cmd/server/main.go"
echo "2. 再启动前端: cd frontend && npm run dev" 