# 旅行者助手 (Traveler Agent)

基于Go语言和Gin框架开发的旅游行程规划与推荐系统后端服务。

## 项目架构

```
traveler_agent/
├── configs/            # 配置文件管理（生产、开发、测试环境）
├── controllers/        # HTTP请求入口，参数校验、调用业务服务
├── services/           # 业务逻辑层，处理具体的业务规则
├── repositories/       # 数据访问层，处理数据库操作
├── models/             # 数据模型定义
├── middlewares/        # 中间件，如JWT认证、日志、限流等
├── utils/              # 工具函数（日期、文件、字符串处理等）
├── docs/               # Swagger API文档
├── routers/            # 路由管理
└── main.go             # 程序入口
```

## 技术栈

- **后端框架**: Gin
- **数据库**: PostgreSQL
- **缓存**: Redis
- **API文档**: Swagger
- **部署**: Docker & Kubernetes

## 功能特性

- 用户注册与身份验证
- 景点/酒店/交通数据管理
- 行程规划算法
- 个性化推荐系统
- 实时数据集成（如天气信息）

## 开发环境配置

### 前提条件

- Go 1.16+
- PostgreSQL
- Redis

### 安装依赖

```bash
go mod tidy
```

### 运行服务

```bash
go run main.go
```

服务将在 http://localhost:8080 上运行。

## API 端点示例

- `GET /api/v1/health`: 健康检查
- `GET /api/v1/ping`: 测试端点

## Docker 部署

```bash
docker build -t traveler-agent .
docker run -p 8080:8080 traveler-agent
```

## 许可证

[MIT License](LICENSE) 