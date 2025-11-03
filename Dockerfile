FROM golang:1.21-alpine AS builder

WORKDIR /app

# 复制依赖文件
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# 复制所有源代码
COPY . .

# 编译API服务
RUN CGO_ENABLED=0 GOOS=linux go build -o namegen-api ./cmd/api

# 创建最终镜像
FROM alpine:latest

# 安装基本工具和CA证书
RUN apk --no-cache add ca-certificates tzdata curl

WORKDIR /app

# 从构建阶段复制编译好的二进制文件
COPY --from=builder /app/namegen-api .

# 设置时区
ENV TZ=Asia/Shanghai

# 设置默认端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/api/v1/origins || exit 1

# 元数据标签
LABEL maintainer="Your Name <your.email@example.com>" \
      description="Name Generator API with support for multiple origins and languages" \
      version="1.0"

# 设置入口点
ENTRYPOINT ["./namegen-api"]

# 默认参数 - 可以在运行时被覆盖
CMD ["-port", "8080"]
