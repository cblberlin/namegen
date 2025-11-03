#!/bin/bash
# 名字生成API服务启动脚本（使用nohup）

# 设置端口
PORT=${1:-"8080"}

# 确保有执行权限
chmod +x namegen-api

# 终止已存在的进程（如果有）
pkill -f namegen-api || true

# 使用nohup在后台启动服务
nohup ./namegen-api -port "$PORT" > namegen-api.log 2>&1 &

# 获取进程ID
PID=$!

# 检查进程是否成功启动
if ps -p $PID > /dev/null; then
    echo "名字生成API服务已成功启动，进程ID：$PID"
    echo "API监听端口：$PORT"
    echo "日志文件：$(pwd)/namegen-api.log"
else
    echo "服务启动失败，请检查日志文件：$(pwd)/namegen-api.log"
    exit 1
fi

# 输出查看日志的命令
echo ""
echo "您可以使用以下命令查看日志："
echo "tail -f namegen-api.log" 