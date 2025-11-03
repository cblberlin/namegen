#!/bin/bash
# 名字生成API服务部署脚本

set -e

# 显示帮助信息
show_help() {
  echo "名字生成API服务 - 部署脚本"
  echo "用法: ./deploy.sh [选项]"
  echo "选项:"
  echo "  -h, --help        显示帮助信息"
  echo "  -p, --port PORT   设置服务端口 (默认: 8080)"
  echo "  -d, --down        停止服务并移除容器"
  echo "  -r, --restart     重启服务"
  echo "  -l, --logs        查看日志"
}

# 默认参数
ACTION="deploy"

# 解析命令行参数
while [[ $# -gt 0 ]]; do
  key="$1"
  case $key in
    -h|--help)
      show_help
      exit 0
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

# 执行相应操作
case $ACTION in
  deploy)
    echo "开始部署名字生成API服务..."
    docker-compose up -d --build
    echo "服务已启动，API正在监听端口8080"
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