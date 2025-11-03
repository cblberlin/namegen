# Docker部署指南 - 名字生成API服务

本文档提供了如何使用Docker和Docker Compose部署名字生成API服务的说明。

## 前提条件

- 已安装Docker（推荐19.03或更高版本）
- 已安装Docker Compose（推荐1.27.0或更高版本）
- 基本的命令行操作知识

## 快速开始

1. 克隆代码库到您的服务器
   ```bash
   git clone <你的代码库URL> namegen-api
   cd namegen-api
   ```

2. 启动服务
   ```bash
   # 构建并启动容器
   docker-compose up -d
   ```

3. 验证服务是否正常运行
   ```bash
   curl http://localhost:8080/api/v1/origins
   ```

## 配置选项

### 自定义端口

如果需要更改默认端口（8080），可以编辑`docker-compose.yml`文件：

```yaml
ports:
  - "自定义端口:8080"
```

### 持久化数据

如有需要持久化数据的需求，取消`docker-compose.yml`中相关卷映射的注释：

```yaml
volumes:
  - ./data:/app/data
```

## API使用说明

### 端点

1. 生成名字：
   ```
   GET /api/v1/names
   ```

2. 获取支持的名字起源：
   ```
   GET /api/v1/origins
   ```

### 参数

- `origin`: 名字的起源/国家（如：english, chinese, russian等）
- `gender`: 性别（male, female, both(默认)）
- `count`: 返回的名字数量（默认为1，最大100）
- `mode`: 名字生成模式（full(完整名字), firstname(仅名), lastname(仅姓)）
- `normalize`: 是否将特殊字符标准化为基本拉丁字母（false, true(默认)）

### 示例

```bash
# 生成5个法语男性名字
curl "http://localhost:8080/api/v1/names?origin=french&gender=male&count=5"

# 生成一个俄语女性名字（保留特殊字符）
curl "http://localhost:8080/api/v1/names?origin=russian&gender=female&normalize=false"
```

## 健康检查

容器配置了自动健康检查，每30秒会检查一次API是否正常响应。您可以通过以下命令查看容器健康状态：

```bash
docker ps --format "{{.Names}}: {{.Status}}"
```

## 日志

容器日志配置了自动轮转，以防止日志占用过多磁盘空间。您可以通过以下命令查看容器日志：

```bash
docker logs namegen-api
```

## 关闭服务

```bash
docker-compose down
```

## 升级服务

当有新版本时，可以使用以下步骤升级：

```bash
# 拉取最新代码
git pull

# 重新构建并启动容器
docker-compose up -d --build
``` 