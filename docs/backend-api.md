# 后端 API 接口文档

本文档详细描述了 todoIng 项目的所有后端 API 接口，供前端开发人员使用。

## 基础配置

### 基础URL
```
http://localhost:5000/api
```

### 认证方式
除注册和登录接口外，所有 API 请求都需要在 Header 中包含有效的 JWT Token：
```
Authorization: Bearer <token>
```

### 响应格式
成功响应格式：
```json
{
  // 对于认证接口，返回对象
  // 对于其他接口，直接返回数据
}
```

错误响应格式：
```json
{
  "msg": "错误信息"
}
或
{
  "errors": [
    {
      "msg": "验证错误信息",
      "param": "字段名",
      "location": "body|query|params"
    }
  ]
}
```

## API 接口详情

### 1. 认证相关接口

#### 用户注册
- **URL**: `POST /auth/register`
- **描述**: 创建新用户账户（支持邮箱验证码注册）
- **请求参数**:
  | 参数名 | 类型 | 必填 | 说明 |
  |--------|------|------|------|
  | username | string | 是 | 用户名 |
  | email | string | 是 | 邮箱地址 |
  | password | string | 是 | 密码，至少6位 |
  | emailCode | string | 否 | 邮箱验证码（当启用邮箱验证时必填） |
  | emailCodeId | string | 否 | 邮箱验证码ID（当启用邮箱验证时必填） |

- **请求示例**:
  ```json
  {
    "username": "john_doe",
    "email": "john@example.com",
    "password": "123456",
    "emailCode": "ABC123",
    "emailCodeId": "uuid-string"
  }
  ```

- **成功响应**:
  ```json
  {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
  ```

- **可能的错误**:
  - `400`: 用户名、邮箱或密码验证失败
  - `400`: 用户已存在
  - `400`: 邮箱验证码错误或已过期
  - `500`: 服务器内部错误

#### 用户登录
- **URL**: `POST /auth/login`
- **描述**: 用户登录获取访问令牌（支持密码登录和邮箱验证码登录）
- **请求参数**:
  | 参数名 | 类型 | 必填 | 说明 |
  |--------|------|------|------|
  | email | string | 是 | 邮箱地址 |
  | password | string | 否 | 密码（密码登录时必填） |
  | emailCode | string | 否 | 邮箱验证码（邮箱验证码登录时必填） |
  | emailCodeId | string | 否 | 邮箱验证码ID（邮箱验证码登录时必填） |
  | captcha | string | 否 | 图片验证码（当启用验证码且使用密码登录时必填） |
  | captchaId | string | 否 | 图片验证码ID（当启用验证码且使用密码登录时必填） |

- **请求示例（密码登录）**:
  ```json
  {
    "email": "john@example.com",
    "password": "123456",
    "captcha": "XYZ789",
    "captchaId": "uuid-string"
  }
  ```

- **请求示例（邮箱验证码登录）**:
  ```json
  {
    "email": "john@example.com",
    "emailCode": "ABC123",
    "emailCodeId": "uuid-string"
  }
  ```

- **成功响应**:
  ```json
  {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
  ```

- **可能的错误**:
  - `400`: 邮箱或密码验证失败
  - `400`: 无效的凭证
  - `400`: 邮箱验证码错误或已过期
  - `400`: 图片验证码错误或已过期
  - `500`: 服务器内部错误

#### 发送邮箱验证码
- **URL**: `POST /auth/send-email-code`
- **描述**: 发送邮箱验证码（用于登录或注册）
- **请求参数**:
  | 参数名 | 类型 | 必填 | 说明 |
  |--------|------|------|------|
  | email | string | 是 | 邮箱地址 |

- **请求示例**:
  ```json
  {
    "email": "john@example.com"
  }
  ```

- **成功响应**:
  ```json
  {
    "msg": "Verification code sent successfully",
    "id": "uuid-string"
  }
  ```

- **可能的错误**:
  - `400`: 邮箱格式错误
  - `400`: 用户不存在（登录时）
  - `400`: 用户已存在（注册时）
  - `500`: 邮件发送失败

#### 获取图片验证码
- **URL**: `GET /auth/captcha`
- **描述**: 获取图片验证码（当启用验证码功能时）
- **请求参数**: 无

- **成功响应**:
  ```json
  {
    "image": "data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMTUwIiBoZWlnaHQ9IjUwIiB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPjxyZWN0IHdpZHRoPSIxMDAlIiBoZWlnaHQ9IjEwMCUiIGZpbGw9IiNmMGYwZjAiLz48bGluZSB4MT0iNzMuMDQ5MTMzMzg0MDQ4OTQiIHkxPSIyNy4wMjI0NzU4MzA2NzQxNyIgeDI9IjcwLjM5MTc5MjM0NTM5MjY4NSIgeTI9IjI0LjUxMjg1MzA0NTQ2MzU2MiIgc3Ryb2tlPSIjZjRlY2Y1IiBzdHJva2Utd2lkdGg9IjEiLz48dGV4dCB4PSIzNi43NjQ3MDU4ODIzNTI5NDkiIHk9IjMwLjY2NjY2NjY2NjY2NjY2NSIgZm9udC1zaXplPSIyNCIgZmlsbD0iI2U1Y2M1NiIgdHJhbnNmb3JtPSJyb3RhdGUoLTI0LjUxMzMxMzYwOTQ2NzQ1NiA1Mi4yMzUyOTQxMTc2NDcwNiAzMC42NjY2NjY2NjY2NjY2NjUpIiBmb250LWZhbWlseT0iQXJpYWwiPjA8L3RleHQ+PC9zdmc+",
    "id": "uuid-string"
  }
  ```

- **可能的错误**:
  - `400`: 验证码功能未启用
  - `500`: 服务器内部错误

### 2. 任务管理接口

#### 创建任务
- **URL**: `POST /tasks`
- **描述**: 创建新任务
- **请求参数**:
  | 参数名 | 类型 | 必填 | 说明 |
  |--------|------|------|------|
  | title | string | 是 | 任务标题 |
  | description | string | 否 | 任务描述 |
  | status | string | 否 | 任务状态: 'To Do', 'In Progress', 'Done' (默认: 'To Do') |
  | priority | string | 否 | 优先级: 'Low', 'Medium', 'High' (默认: 'Medium') |
  | assignee | string | 否 | 分配给用户的ID |
  | dueDate | string | 否 | 截止日期 (ISO 8601格式) |
  | scheduledDate | string | 否 | 计划日期 (ISO 8601格式) |

- **请求示例**:
  ```json
  {
    "title": "完成项目文档",
    "description": "编写并完善项目技术文档",
    "status": "In Progress",
    "priority": "High",
    "dueDate": "2025-08-10T10:00:00Z",
    "scheduledDate": "2025-08-05T09:00:00Z"
  }
  ```

- **成功响应**:
  ```json
  {
    "_id": "5f9d3b3b9d9b3b0017f1b3f1",
    "title": "完成项目文档",
    "description": "编写并完善项目技术文档",
    "status": "In Progress",
    "priority": "High",
    "assignee": "5f9d3b3b9d9b3b0017f1b3f2",
    "dueDate": "2025-08-10T10:00:00.000Z",
    "scheduledDate": "2025-08-05T09:00:00.000Z",
    "createdAt": "2025-08-04T10:00:00.000Z",
    "updatedAt": "2025-08-04T10:00:00.000Z",
    "history": []
  }
  ```

- **可能的错误**:
  - `400`: 任务数据验证失败
  - `401`: 未认证
  - `500`: 服务器内部错误

#### 获取任务列表
- **URL**: `GET /tasks`
- **描述**: 获取当前用户创建的所有任务，按创建时间倒序排列
- **查询参数**: 无

- **成功响应**:
  ```json
  [
    {
      "_id": "5f8d0d5f4e67c4001c8b4567",
      "title": "完成项目文档",
      "description": "编写并完成项目技术文档",
      "status": "In Progress",
      "priority": "High",
      "createdBy": "5f8d0d5f4e67c4001c8b4568",
      "createdAt": "2023-01-01T00:00:00.000Z",
      "updatedAt": "2023-01-01T00:00:00.000Z"
    }
  ]
  ```

- **可能的错误**:
  - `401`: 未授权访问
  - `500`: 服务器内部错误

#### 获取任务详情
- **URL**: `GET /tasks/:id`
- **描述**: 获取指定任务的详细信息
- **路径参数**:
  | 参数名 | 类型 | 必填 | 说明 |
  |--------|------|------|------|
  | id | string | 是 | 任务ID |

- **成功响应**:
  ```json
  {
    "_id": "5f8d0d5f4e67c4001c8b4567",
    "title": "完成项目文档",
    "description": "编写并完成项目技术文档",
    "status": "In Progress",
    "priority": "High",
    "createdBy": "5f8d0d5f4e67c4001c8b4568",
    "createdAt": "2023-01-01T00:00:00.000Z",
    "updatedAt": "2023-01-01T00:00:00.000Z"
  }
  ```

- **可能的错误**:
  - `404`: 任务未找到
  - `401`: 未授权访问
  - `500`: 服务器内部错误

#### 更新任务
- **URL**: `PUT /tasks/:id`
- **描述**: 更新指定任务的信息
- **路径参数**:
  | 参数名 | 类型 | 必填 | 说明 |
  |--------|------|------|------|
  | id | string | 是 | 任务ID |

- **请求参数**:
  | 参数名 | 类型 | 必填 | 说明 |
  |--------|------|------|------|
  | title | string | 否 | 任务标题 |
  | description | string | 否 | 任务描述 |
  | status | string | 否 | 任务状态: 'To Do', 'In Progress', 'Done' |
  | priority | string | 否 | 优先级: 'Low', 'Medium', 'High' |
  | assignee | string | 否 | 分配给用户的ID |

- **请求示例**:
  ```json
  {
    "title": "完成并提交项目文档",
    "status": "Done"
  }
  ```

- **成功响应**:
  ```json
  {
    "_id": "5f8d0d5f4e67c4001c8b4567",
    "title": "完成并提交项目文档",
    "description": "编写并完成项目技术文档",
    "status": "Done",
    "priority": "High",
    "createdBy": "5f8d0d5f4e67c4001c8b4568",
    "createdAt": "2023-01-01T00:00:00.000Z",
    "updatedAt": "2023-01-02T00:00:00.000Z"
  }
  ```

- **可能的错误**:
  - `400`: 验证失败
  - `404`: 任务未找到
  - `401`: 未授权访问
  - `500`: 服务器内部错误

#### 删除任务
- **URL**: `DELETE /tasks/:id`
- **描述**: 删除指定任务
- **路径参数**:
  | 参数名 | 类型 | 必填 | 说明 |
  |--------|------|------|------|
  | id | string | 是 | 任务ID |

- **成功响应**:
  ```json
  {
    "msg": "Task removed"
  }
  ```

- **可能的错误**:
  - `404`: 任务未找到
  - `401`: 未授权访问
  - `500`: 服务器内部错误