# 旅行者助手部署指南

本文档提供了旅行者助手系统的部署指南，包括开发环境和生产环境的部署方法。

## 开发环境部署

### 前提条件

- Go 1.16+
- PostgreSQL 12+
- Redis 6+
- Git

### 步骤

1. **克隆代码库**

   ```bash
   git clone https://github.com/JasonRobertDestiny/Travel_Planning_Assistant.git
   cd Travel_Planning_Assistant
   ```

2. **安装依赖**

   ```bash
   go mod tidy
   ```

3. **配置环境**

   复制示例配置文件并根据需要修改：

   ```bash
   cp configs/config.example.json configs/config.json
   ```

   修改 `configs/config.json` 中的数据库和Redis连接信息。

4. **初始化数据库**

   确保PostgreSQL服务已启动，并创建相应的数据库：

   ```bash
   createdb traveler
   ```

   运行迁移脚本：

   ```bash
   go run utils/migrate.go up
   ```

5. **启动服务**

   ```bash
   go run main.go
   ```

   服务将在 http://localhost:8080 上运行。

## 生产环境部署

### 使用Docker部署

1. **构建Docker镜像**

   ```bash
   docker build -t traveler-agent:latest .
   ```

2. **运行Docker容器**

   ```bash
   docker run -d \
     --name traveler-agent \
     -p 8080:8080 \
     -e DB_HOST=your-db-host \
     -e DB_PORT=your-db-port \
     -e DB_USER=your-db-user \
     -e DB_PASSWORD=your-db-password \
     -e DB_NAME=your-db-name \
     -e REDIS_HOST=your-redis-host \
     -e REDIS_PORT=your-redis-port \
     -e REDIS_PASSWORD=your-redis-password \
     -e JWT_SECRET=your-jwt-secret \
     traveler-agent:latest
   ```

### 使用Docker Compose部署

1. **创建docker-compose.yml文件**

   ```yaml
   version: '3'
   services:
     app:
       image: traveler-agent:latest
       build: .
       ports:
         - "8080:8080"
       environment:
         - DB_HOST=postgres
         - DB_PORT=5432
         - DB_USER=postgres
         - DB_PASSWORD=postgres
         - DB_NAME=traveler
         - REDIS_HOST=redis
         - REDIS_PORT=6379
         - REDIS_PASSWORD=
         - JWT_SECRET=your-jwt-secret
       depends_on:
         - postgres
         - redis
     
     postgres:
       image: postgres:12
       ports:
         - "5432:5432"
       environment:
         - POSTGRES_USER=postgres
         - POSTGRES_PASSWORD=postgres
         - POSTGRES_DB=traveler
       volumes:
         - postgres-data:/var/lib/postgresql/data
     
     redis:
       image: redis:6
       ports:
         - "6379:6379"
       volumes:
         - redis-data:/data
   
   volumes:
     postgres-data:
     redis-data:
   ```

2. **启动服务**

   ```bash
   docker-compose up -d
   ```

### 使用Kubernetes部署

1. **应用Kubernetes配置**

   ```bash
   kubectl apply -f kubernetes/
   ```

   这将部署所有必要的资源，包括Deployment、Service、ConfigMap和Secret。

## 环境变量

以下是系统支持的环境变量列表，可用于覆盖配置文件中的设置：

| 环境变量 | 描述 | 默认值 |
|---------|------|-------|
| SERVER_PORT | 服务器端口 | 8080 |
| SERVER_ENV | 运行环境 (development/production) | development |
| DB_HOST | 数据库主机 | localhost |
| DB_PORT | 数据库端口 | 5432 |
| DB_USER | 数据库用户名 | postgres |
| DB_PASSWORD | 数据库密码 | postgres |
| DB_NAME | 数据库名称 | traveler |
| REDIS_HOST | Redis主机 | localhost |
| REDIS_PORT | Redis端口 | 6379 |
| REDIS_PASSWORD | Redis密码 | |
| REDIS_DB | Redis数据库索引 | 0 |
| JWT_SECRET | JWT密钥 | your-secret-key |
| JWT_EXPIRES_IN | JWT过期时间（小时） | 24 |

## 数据备份与恢复

### 数据库备份

```bash
pg_dump -U postgres -d traveler > backup.sql
```

### 数据库恢复

```bash
psql -U postgres -d traveler < backup.sql
```

## 日志管理

系统日志存储在 `logs/` 目录下，按日期分割。

在生产环境中，建议配置日志转发到集中式日志管理系统，如ELK Stack。

## 监控

服务暴露了Prometheus监控指标，可通过 `/metrics` 端点访问。

## 故障排除

1. **服务无法启动**
   - 检查配置文件是否正确
   - 验证数据库和Redis连接是否正常
   - 查看日志文件获取详细错误信息

2. **API响应慢**
   - 检查数据库性能
   - 验证Redis缓存是否正常工作
   - 考虑增加服务实例数量

3. **数据库连接错误**
   - 验证数据库凭据
   - 确保数据库服务器正常运行
   - 检查防火墙设置 