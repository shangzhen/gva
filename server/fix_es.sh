#!/bin/bash

# Elasticsearch 集成快速修复脚本
# 用于 Go 1.18 环境

echo "====== Elasticsearch 集成修复脚本 ======"
echo ""

# 检查 Go 版本
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo "当前 Go 版本: $GO_VERSION"
echo ""

# 检查版本是否 >= 1.21
REQUIRED_VERSION="1.21.0"
if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" = "$REQUIRED_VERSION" ]; then
    echo "✅ Go 版本满足要求 (>= 1.21)"
    echo ""
    echo "执行编译..."
    go mod tidy
    go build -o gin-vue-admin main.go
    if [ $? -eq 0 ]; then
        echo "✅ 编译成功！"
        echo ""
        echo "可以运行: ./gin-vue-admin"
    else
        echo "❌ 编译失败"
        exit 1
    fi
else
    echo "⚠️  Go 版本过低 (当前: $GO_VERSION, 需要: >= 1.21)"
    echo ""
    echo "============ 解决方案 ============"
    echo ""
    echo "方案 1: 升级 Go 版本（推荐）"
    echo "  1. 访问 https://go.dev/dl/"
    echo "  2. 下载并安装 Go 1.21+ 版本"
    echo "  3. 重新运行此脚本"
    echo ""
    echo "方案 2: 使用 Docker 运行"
    echo "  docker run -it --rm -v \$(pwd):/app -w /app golang:1.21 bash"
    echo "  cd /app && go mod tidy && go build"
    echo ""
    echo "方案 3: 仅验证 ES 代码语法"
    echo "  运行: go run check_syntax.go"
    echo "  （这将验证 ES 代码本身没有语法错误）"
    echo ""

    # 运行语法检查
    echo "====== 运行 ES 代码语法检查 ======"
    if [ -f "check_syntax.go" ]; then
        go run check_syntax.go
        echo ""
    fi

    exit 1
fi
