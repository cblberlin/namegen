#!/bin/bash
# 名字生成API服务部署脚本

set -e

# 显示帮助信息
show_help() {
  echo "名字生成API服务 - 部署脚本"
  echo "用法: ./deploy.sh [选项]"
  echo "选项:"
  echo "  -h, --help        显示帮助信息"
  echo "  -k, --key KEY     设置API密钥 (如不指定，将使用.env文件中的API_KEY或默认值)"
  echo "  -p, --port PORT   设置服务端口 (默认: 8080)"
  echo "  -d, --down        停止服务并移除容器"
  echo "  -r, --restart     重启服务"
  echo "  -l, --logs        查看日志"
}

# 默认参数
API_KEY=""
ACTION="deploy"

# 解析命令行参数
while [[ $# -gt 0 ]]; do
  key="$1"
  case $key in
    -h|--help)
      show_help
      exit 0
      ;;
    -k|--key)
      API_KEY="$2"
      shift 2
      ;;
    -d|--down)
      ACTION="down"
      shift
      ;;
    -r|--restart)
      ACTION="restart"
      shift
      ;;
    -l|--logs)
      ACTION="logs"
      shift
      ;;
    *)
      echo "未知选项: $1"
      show_help
      exit 1
      ;;
  esac
done

# 如果.env文件存在，加载它
if [ -f .env ]; then
  source .env
fi

# 如果命令行参数中指定了API_KEY，则覆盖.env中的设置
if [ -n "$API_KEY" ]; then
  export API_KEY
fi

# 如果仍然没有设置API_KEY，则使用默认值
if [ -z "$API_KEY" ]; then
  echo "警告: 未设置API_KEY，将使用默认值 'mysecretkey'"
  echo "为提高安全性，建议通过.env文件或命令行参数设置自定义密钥"
  export API_KEY="mysecretkey"
fi

# 执行相应操作
case $ACTION in
  deploy)
    echo "开始部署名字生成API服务..."
    docker-compose up -d --build
    echo "服务已启动，API正在监听端口8080"
    echo "API密钥: $API_KEY"
    ;;
  down)
    echo "停止服务并移除容器..."
    docker-compose down
    echo "服务已停止"
    ;;
  restart)
    echo "重启服务..."
    docker-compose restart
    echo "服务已重启"
    ;;
  logs)
    echo "显示服务日志..."
    docker-compose logs -f
    ;;
esac

exit 0 