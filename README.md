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
├── docs/               # 文档（API文档、部署文档等）
├── routers/            # 路由管理
├── Dockerfile          # Docker镜像构建文件
├── docker-compose.yml  # Docker Compose配置文件
└── main.go             # 程序入口
```

## 技术栈

- **后端框架**: [Gin](https://github.com/gin-gonic/gin)
- **数据库**: [PostgreSQL](https://www.postgresql.org/)
- **ORM**: [pgx](https://github.com/jackc/pgx)
- **缓存**: [Redis](https://redis.io/)
- **认证**: JWT (JSON Web Token)
- **部署**: Docker & Docker Compose
- **日志**: logrus

## 功能特性

- 用户注册与身份验证
  - 用户注册、登录、密码重置
  - JWT认证
  - 用户偏好设置

- 景点/酒店/交通数据管理
  - 景点信息查询
  - 城市和国家筛选
  - 评分和分类过滤

- 行程规划功能
  - 创建、修改、共享行程
  - 行程日历视图
  - 添加景点到行程
  - 设置参观时间和备注

- 个性化推荐系统
  - 基于用户历史行为的推荐
  - 基于目的地的推荐
  - 季节性推荐

- 实时数据集成
  - 天气信息
  - 汇率转换
  - 当地活动

## 开发环境配置

### 前提条件

- Go 1.16+
- PostgreSQL 12+
- Redis 6+
- Git

### 安装依赖

```bash
go mod tidy
```

### 配置环境

复制示例配置文件并根据需要修改：

```bash
cp configs/config.example.json configs/config.json
```

### 运行服务

```bash
go run main.go
```

服务将在 http://localhost:8080 上运行。

## 使用Docker运行

### 使用Docker Compose

```bash
docker-compose up -d
```

这将启动以下服务：
- 旅行者助手API (http://localhost:8080)
- PostgreSQL数据库
- Redis缓存
- pgAdmin管理界面 (http://localhost:5050)

### 仅构建API服务

```bash
docker build -t traveler-agent .
docker run -p 8080:8080 traveler-agent
```

## API文档

详细的API文档可以在以下位置找到：

- [API文档](docs/api/README.md) - 提供所有API端点的详细说明
- [部署文档](docs/deployment.md) - 提供部署相关信息

## 常用API端点

- `GET /api/v1/health`: 健康检查
- `POST /api/v1/auth/register`: 用户注册
- `POST /api/v1/auth/login`: 用户登录
- `GET /api/v1/attractions`: 获取景点列表
- `GET /api/v1/attractions/{id}`: 获取景点详情
- `POST /api/v1/itineraries`: 创建行程
- `GET /api/v1/itineraries`: 获取用户行程列表
- `GET /api/v1/itineraries/{id}`: 获取行程详情

## 与前端集成

旅行者助手后端服务设计为与前端框架（如React、Vue或Angular）集成。API遵循RESTful设计原则，方便前端集成。

主要集成点：
1. 用户认证 - 通过JWT令牌进行认证
2. 数据获取 - 使用标准的HTTP方法获取和提交数据
3. 实时更新 - 提供Webhook支持实时通知

## 贡献指南

欢迎贡献代码、报告问题或提出新功能建议。

1. Fork仓库
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 提交Pull Request

## 许可证

[MIT License](LICENSE) 