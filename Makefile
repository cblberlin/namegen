.PHONY: help build build-cli build-api test clean run-cli run-api

# 默认目标
help: ## 显示帮助信息
	@echo "可用命令:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

# 构建相关
build: build-cli build-api ## 构建CLI工具和API服务

build-cli: ## 构建命令行工具
	@echo "构建命令行工具..."
	@mkdir -p build
	go build -o build/namegen cmd/namegen/main.go

build-api: ## 构建API服务
	@echo "构建API服务..."
	@mkdir -p build
	go build -o build/namegen-api cmd/api/main.go

# 测试相关
test: ## 运行所有测试
	@echo "运行测试..."
	go test ./...

# 清理
clean: ## 清理构建产物
	@echo "清理构建文件..."
	rm -rf build/
	rm -f namegen-api.log

# 运行相关
run-cli: build-cli ## 运行命令行工具
	./build/namegen $(ARGS)

run-api: build-api ## 运行API服务（前台模式）
	@echo "启动API服务..."
	./build/namegen-api $(ARGS)

run-api-bg: build-api ## 运行API服务（后台模式）
	@echo "启动API服务（后台模式）..."
	@PORT=$${PORT:-8080}; \
	nohup ./build/namegen-api -port $$PORT > namegen-api.log 2>&1 & \
	PID=$$!; \
	echo "API服务已启动，PID: $$PID，端口: $$PORT"; \
	echo "日志文件: namegen-api.log"; \
	echo "停止服务: kill $$PID 或者 make stop-api"

stop-api: ## 停止后台运行的API服务
	@echo "停止API服务..."
	@pkill -f namegen-api || echo "没有找到运行中的API服务"

# 开发相关
deps: ## 下载依赖
	@echo "下载依赖..."
	go mod download
	go mod tidy

fmt: ## 格式化代码
	@echo "格式化代码..."
	go fmt ./...

lint: ## 检查代码质量
	@echo "检查代码..."
	go vet ./...

# CI/CD相关
ci: deps fmt lint test build ## 持续集成流程