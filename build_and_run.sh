#!/bin/bash
# 编译并启动名字生成API服务

# 设置端口
PORT=${1:-"8080"}

# 确保Go环境正确
echo "检查Go版本..."
go version

# 安装必要的依赖
echo "安装依赖..."
go get golang.org/x/text@v0.12.0

# 编译API
echo "编译API服务..."
go build -o namegen-api ./cmd/api

# 确保编译成功
if [ ! -f "./namegen-api" ]; then
    echo "编译失败，请检查错误信息"
    exit 1
fi

# 终止已存在的进程
echo "终止已存在的进程（如果有）..."
pkill -f namegen-api || true

# 使用nohup在后台启动服务
echo "启动API服务..."
nohup ./namegen-api -port "$PORT" > namegen-api.log 2>&1 &

# 获取进程ID
PID=$!

# 检查进程是否成功启动
sleep 2
if ps -p $PID > /dev/null; then
    echo "名字生成API服务已成功启动，进程ID：$PID"
    echo "API监听端口：$PORT"
    echo "日志文件：$(pwd)/namegen-api.log"
else
    echo "服务启动失败，请检查日志文件：$(pwd)/namegen-api.log"
    exit 1
fi

# 显示日志的最新内容
echo ""
echo "日志的最新内容："
tail -n 20 namegen-api.log

# 输出查看日志的命令
echo ""
echo "您可以使用以下命令查看日志："
echo "tail -f namegen-api.log" 