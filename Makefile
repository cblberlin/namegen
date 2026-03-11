.PHONY: help build build-cli build-api test clean run-cli run-api deps fmt lint

# 变量定义
BINARY_DIR=build
API_BINARY=$(BINARY_DIR)/namegen-api
CLI_BINARY=$(BINARY_DIR)/namegen
API_MAIN=cmd/api/main.go
CLI_MAIN=cmd/namegen/main.go

# 默认目标
help: ## 显示帮助信息
	@echo "可用命令:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

# 构建相关
build: build-cli build-api ## 构建CLI工具和API服务

build-cli: deps ## 构建命令行工具
	@echo "构建命令行工具..."
	@mkdir -p $(BINARY_DIR)
	go build -o $(CLI_BINARY) $(CLI_MAIN)
	@chmod +x $(CLI_BINARY)

build-api: deps ## 构建API服务
	@echo "构建API服务..."
	@mkdir -p $(BINARY_DIR)
	# 确保清理掉旧的 api 目录干扰
	@rm -rf api/ 
	go build -o $(API_BINARY) $(API_MAIN)
	@chmod +x $(API_BINARY)
	@echo "✅ 构建成功: $(API_BINARY)"

# 运行相关
run-api: build-api ## 运行API服务（前台模式）
	@echo "🚀 启动API服务..."
	./$(API_BINARY) $(ARGS)

run-api-bg: build-api ## 运行API服务（后台模式）
	@echo "启动API服务（后台模式）..."
	@PORT=$${PORT:-8080}; \
	nohup ./$(API_BINARY) -port $$PORT > namegen-api.log 2>&1 & \
	PID=$$!; \
	echo "API服务已启动, PID: $$PID, 端口: $$PORT"; \
	echo "停止服务: make stop-api"

stop-api: ## 停止后台运行的API服务
	@echo "停止API服务..."
	@pkill -f namegen-api || echo "没有找到运行中的API服务"

# 维护相关
deps: ## 整理并下载依赖
	@echo "整理依赖中..."
	go mod tidy
	go mod download

clean: ## 清理构建产物
	@echo "清理构建文件..."
	rm -rf $(BINARY_DIR)
	rm -f namegen-api.log

fmt: ## 格式化代码
	go fmt ./...

lint: ## 检查代码质量
	go vet ./...