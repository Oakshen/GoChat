#!/bin/bash

# GoChat Vue.js 前端启动脚本

echo "🚀 GoChat Vue.js 前端启动中..."

# 检查 Node.js 是否安装
if ! command -v node &> /dev/null; then
    echo "❌ 错误: 请先安装 Node.js"
    exit 1
fi

# 检查 npm 是否安装
if ! command -v npm &> /dev/null; then
    echo "❌ 错误: 请先安装 npm"
    exit 1
fi

# 进入前端目录
cd "$(dirname "$0")/.."

# 检查 package.json 是否存在
if [ ! -f "package.json" ]; then
    echo "❌ 错误: 未找到 package.json 文件"
    exit 1
fi

# 检查 node_modules 是否存在
if [ ! -d "node_modules" ]; then
    echo "📦 安装依赖中..."
    npm install
    if [ $? -ne 0 ]; then
        echo "❌ 错误: 依赖安装失败"
        exit 1
    fi
fi

# 启动开发服务器
echo "🌐 启动开发服务器..."
echo "📍 前端地址: http://localhost:3000"
echo "🔗 后端地址: http://localhost:8080"
echo ""
echo "💡 提示: 请确保后端服务正在运行"
echo "💡 提示: 按 Ctrl+C 停止服务"
echo ""

npm run dev 