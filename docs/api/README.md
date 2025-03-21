# 旅行者助手 API 文档

本文档详细说明了旅行者助手系统的API接口，供前端开发人员参考。

## 基础信息

- **基础URL**: `http://localhost:8080/api/v1`
- **认证方式**: JWT Token (Authorization: Bearer {token})
- **响应格式**: JSON

## 通用响应格式

```json
{
  "status": "success|error",
  "message": "操作结果描述",
  "data": {}, // 数据对象，成功时返回
  "error": {} // 错误对象，失败时返回
}
```

## API 端点

### 1. 用户认证

#### 1.1 用户注册

- **URL**: `/auth/register`
- **方法**: `POST`
- **描述**: 注册新用户
- **请求体**:
  ```json
  {
    "username": "string",
    "email": "string",
    "password": "string"
  }
  ```
- **响应**:
  ```json
  {
    "status": "success",
    "message": "注册成功",
    "data": {
      "user_id": "uuid",
      "username": "string",
      "email": "string"
    }
  }
  ```

#### 1.2 用户登录

- **URL**: `/auth/login`
- **方法**: `POST`
- **描述**: 用户登录获取Token
- **请求体**:
  ```json
  {
    "email": "string",
    "password": "string"
  }
  ```
- **响应**:
  ```json
  {
    "status": "success",
    "message": "登录成功",
    "data": {
      "token": "string",
      "expires_at": "timestamp",
      "user": {
        "user_id": "uuid",
        "username": "string",
        "email": "string"
      }
    }
  }
  ```

### 2. 景点管理

#### 2.1 获取景点列表

- **URL**: `/attractions`
- **方法**: `GET`
- **描述**: 获取景点列表
- **参数**:
  - `page`: 页码 (默认: 1)
  - `limit`: 每页数量 (默认: 10)
  - `city`: 城市筛选 (可选)
  - `category`: 类别筛选 (可选)
  - `sort`: 排序方式 (可选，支持rating、popularity)
- **响应**:
  ```json
  {
    "status": "success",
    "data": {
      "total": 100,
      "page": 1,
      "limit": 10,
      "attractions": [
        {
          "id": "uuid",
          "name": "景点名称",
          "description": "景点描述",
          "location": "具体位置",
          "city": "所在城市",
          "country": "所在国家",
          "category": "类别",
          "rating": 4.5,
          "images": ["url1", "url2"],
          "visiting_time": 120
        }
      ]
    }
  }
  ```

#### 2.2 获取单个景点详情

- **URL**: `/attractions/{id}`
- **方法**: `GET`
- **描述**: 获取单个景点详细信息
- **响应**:
  ```json
  {
    "status": "success",
    "data": {
      "id": "uuid",
      "name": "景点名称",
      "description": "景点描述",
      "location": "具体位置",
      "city": "所在城市",
      "country": "所在国家",
      "category": "类别",
      "rating": 4.5,
      "images": ["url1", "url2"],
      "visiting_time": 120,
      "open_hours": {
        "monday": "9:00-17:00",
        "tuesday": "9:00-17:00"
      },
      "ticket_info": {
        "adult": 100,
        "child": 50
      },
      "reviews": [
        {
          "user_id": "uuid",
          "username": "用户名",
          "rating": 5,
          "comment": "评论内容",
          "created_at": "timestamp"
        }
      ]
    }
  }
  ```

### 3. 行程管理

#### 3.1 创建行程

- **URL**: `/itineraries`
- **方法**: `POST`
- **描述**: 创建新行程
- **认证**: 需要
- **请求体**:
  ```json
  {
    "title": "行程标题",
    "start_date": "YYYY-MM-DD",
    "end_date": "YYYY-MM-DD",
    "destination": "目的地",
    "description": "行程描述",
    "is_public": false
  }
  ```
- **响应**:
  ```json
  {
    "status": "success",
    "message": "行程创建成功",
    "data": {
      "id": "uuid",
      "title": "行程标题",
      "start_date": "YYYY-MM-DD",
      "end_date": "YYYY-MM-DD",
      "destination": "目的地",
      "description": "行程描述",
      "is_public": false,
      "created_at": "timestamp"
    }
  }
  ```

#### 3.2 获取用户行程列表

- **URL**: `/itineraries`
- **方法**: `GET`
- **描述**: 获取当前用户的行程列表
- **认证**: 需要
- **响应**:
  ```json
  {
    "status": "success",
    "data": {
      "itineraries": [
        {
          "id": "uuid",
          "title": "行程标题",
          "start_date": "YYYY-MM-DD",
          "end_date": "YYYY-MM-DD",
          "destination": "目的地",
          "description": "行程描述",
          "is_public": false,
          "created_at": "timestamp"
        }
      ]
    }
  }
  ```

#### 3.3 获取行程详情

- **URL**: `/itineraries/{id}`
- **方法**: `GET`
- **描述**: 获取行程详细信息
- **认证**: 需要（如果是私有行程）
- **响应**:
  ```json
  {
    "status": "success",
    "data": {
      "id": "uuid",
      "title": "行程标题",
      "start_date": "YYYY-MM-DD",
      "end_date": "YYYY-MM-DD",
      "destination": "目的地",
      "description": "行程描述",
      "is_public": false,
      "created_at": "timestamp",
      "days": [
        {
          "day": 1,
          "date": "YYYY-MM-DD",
          "items": [
            {
              "id": "uuid",
              "type": "attraction",
              "attraction_id": "uuid",
              "name": "景点名称",
              "start_time": "HH:MM",
              "end_time": "HH:MM",
              "note": "备注"
            }
          ]
        }
      ]
    }
  }
  ```

#### 3.4 添加行程项

- **URL**: `/itineraries/{id}/items`
- **方法**: `POST`
- **描述**: 添加行程项目
- **认证**: 需要
- **请求体**:
  ```json
  {
    "day": 1,
    "type": "attraction",
    "attraction_id": "uuid",
    "start_time": "HH:MM",
    "end_time": "HH:MM",
    "note": "备注"
  }
  ```
- **响应**:
  ```json
  {
    "status": "success",
    "message": "行程项目添加成功",
    "data": {
      "id": "uuid",
      "day": 1,
      "type": "attraction",
      "attraction_id": "uuid",
      "name": "景点名称",
      "start_time": "HH:MM",
      "end_time": "HH:MM",
      "note": "备注"
    }
  }
  ```

### 4. 系统相关

#### 4.1 健康检查

- **URL**: `/health`
- **方法**: `GET`
- **描述**: 系统健康状态检查
- **响应**:
  ```json
  {
    "status": "ok",
    "message": "服务运行正常"
  }
  ```

## 错误码说明

| 错误码 | 描述 |
| ----- | ---- |
| 400 | 请求参数错误 |
| 401 | 未认证或认证失败 |
| 403 | 权限不足 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

## 开发环境测试账户

- **邮箱**: `test@example.com`
- **密码**: `password123` 