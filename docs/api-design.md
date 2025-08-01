# API 设计

## 概述

本文档描述了 todoIng 任务管理系统的 API 接口设计，特别关注任务生命周期管理和变更历史追踪功能。

## 基础配置

### 基础URL
```
https://api.todoing.example.com/v1
```

### 认证方式
所有 API 请求都需要在 Header 中包含有效的 JWT Token：
```
Authorization: Bearer <token>
```

### 响应格式
所有响应都采用 JSON 格式：
```json
{
  "success": true,
  "data": {},
  "message": "Success message"
}
```

错误响应格式：
```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Error description"
  }
}
```

## API 详细设计

### 1. 认证相关接口

#### 用户注册
- **URL**: `POST /auth/register`
- **描述**: 创建新用户账户
- **请求参数**:
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
    "success": true,
    "data": {
      "user": {
        "_id": "ObjectId",
        "username": "string",
        "email": "string",
        "createdAt": "date"
      },
      "token": "jwt_token"
    }
  }
  ```

#### 用户登录
- **URL**: `POST /auth/login`
- **描述**: 用户登录获取访问令牌
- **请求参数**:
  ```json
  {
    "email": "string",
    "password": "string"
  }
  ```
- **响应**:
  ```json
  {
    "success": true,
    "data": {
      "user": {
        "_id": "ObjectId",
        "username": "string",
        "email": "string"
      },
      "token": "jwt_token"
    }
  }
  ```

### 2. 任务管理接口

#### 获取任务列表
- **URL**: `GET /tasks`
- **描述**: 获取用户任务列表，支持筛选和分页
- **查询参数**:
  - `status`: 任务状态筛选 (created, in-progress, paused, completed, cancelled)
  - `priority`: 优先级筛选 (low, medium, high)
  - `assignee`: 分配人ID
  - `team`: 团队ID
  - `page`: 页码 (默认: 1)
  - `limit`: 每页数量 (默认: 10, 最大: 100)
- **响应**:
  ```json
  {
    "success": true,
    "data": {
      "tasks": [
        {
          "_id": "ObjectId",
          "title": "string",
          "description": "string",
          "status": "string",
          "priority": "string",
          "assignee": "ObjectId",
          "team": "ObjectId",
          "dueDate": "date",
          "createdAt": "date",
          "updatedAt": "date"
        }
      ],
      "pagination": {
        "page": 1,
        "limit": 10,
        "total": 100,
        "pages": 10
      }
    }
  }
  ```

#### 创建任务
- **URL**: `POST /tasks`
- **描述**: 创建新任务
- **请求参数**:
  ```json
  {
    "title": "string",
    "description": "string",
    "status": "string",
    "priority": "string",
    "assignee": "ObjectId",
    "team": "ObjectId",
    "dueDate": "date",
    "startDate": "date",
    "estimatedHours": "number",
    "tags": ["string"]
  }
  ```
- **响应**:
  ```json
  {
    "success": true,
    "data": {
      "task": {
        "_id": "ObjectId",
        "title": "string",
        "description": "string",
        "status": "string",
        "priority": "string",
        "assignee": "ObjectId",
        "team": "ObjectId",
        "dueDate": "date",
        "startDate": "date",
        "estimatedHours": "number",
        "tags": ["string"],
        "createdBy": "ObjectId",
        "createdAt": "date",
        "updatedAt": "date"
      }
    }
  }
  ```

#### 获取任务详情
- **URL**: `GET /tasks/{id}`
- **描述**: 获取特定任务的详细信息
- **路径参数**:
  - `id`: 任务ID
- **响应**:
  ```json
  {
    "success": true,
    "data": {
      "task": {
        "_id": "ObjectId",
        "title": "string",
        "description": "string",
        "status": "string",
        "priority": "string",
        "assignee": {
          "_id": "ObjectId",
          "username": "string"
        },
        "team": {
          "_id": "ObjectId",
          "name": "string"
        },
        "dueDate": "date",
        "startDate": "date",
        "completedDate": "date",
        "estimatedHours": "number",
        "actualHours": "number",
        "tags": ["string"],
        "attachments": [
          {
            "name": "string",
            "url": "string",
            "size": "number",
            "type": "string"
          }
        ],
        "createdBy": {
          "_id": "ObjectId",
          "username": "string"
        },
        "createdAt": "date",
        "updatedAt": "date"
      }
    }
  }
  ```

#### 更新任务
- **URL**: `PUT /tasks/{id}`
- **描述**: 更新任务信息，会自动记录变更历史
- **路径参数**:
  - `id`: 任务ID
- **请求参数**:
  ```json
  {
    "title": "string",
    "description": "string",
    "status": "string",
    "priority": "string",
    "assignee": "ObjectId",
    "team": "ObjectId",
    "dueDate": "date",
    "startDate": "date",
    "estimatedHours": "number",
    "tags": ["string"]
  }
  ```
- **响应**:
  ```json
  {
    "success": true,
    "data": {
      "task": {
        "_id": "ObjectId",
        "title": "string",
        "description": "string",
        "status": "string",
        "priority": "string",
        "assignee": "ObjectId",
        "team": "ObjectId",
        "dueDate": "date",
        "startDate": "date",
        "estimatedHours": "number",
        "tags": ["string"],
        "createdBy": "ObjectId",
        "createdAt": "date",
        "updatedAt": "date"
      }
    }
  }
  ```

#### 删除任务
- **URL**: `DELETE /tasks/{id}`
- **描述**: 删除指定任务
- **路径参数**:
  - `id`: 任务ID
- **响应**:
  ```json
  {
    "success": true,
    "message": "Task deleted successfully"
  }
  ```

### 3. 任务历史记录接口

#### 获取任务变更历史
- **URL**: `GET /tasks/{id}/history`
- **描述**: 获取特定任务的所有变更历史记录
- **路径参数**:
  - `id`: 任务ID
- **查询参数**:
  - `type`: 变更类型筛选
  - `page`: 页码 (默认: 1)
  - `limit`: 每页数量 (默认: 10, 最大: 100)
- **响应**:
  ```json
  {
    "success": true,
    "data": {
      "history": [
        {
          "_id": "ObjectId",
          "taskId": "ObjectId",
          "field": "string",
          "oldValue": "any",
          "newValue": "any",
          "changedBy": {
            "_id": "ObjectId",
            "username": "string"
          },
          "changeType": "string",
          "comment": "string",
          "timestamp": "date"
        }
      ],
      "pagination": {
        "page": 1,
        "limit": 10,
        "total": 50,
        "pages": 5
      }
    }
  }
  ```

#### 获取所有历史记录
- **URL**: `GET /history`
- **描述**: 获取当前用户相关的所有任务历史记录
- **查询参数**:
  - `taskId`: 任务ID筛选
  - `changeType`: 变更类型筛选
  - `userId`: 用户ID筛选
  - `page`: 页码 (默认: 1)
  - `limit`: 每页数量 (默认: 10, 最大: 100)
- **响应**:
  ```json
  {
    "success": true,
    "data": {
      "history": [
        {
          "_id": "ObjectId",
          "taskId": {
            "_id": "ObjectId",
            "title": "string"
          },
          "field": "string",
          "oldValue": "any",
          "newValue": "any",
          "changedBy": {
            "_id": "ObjectId",
            "username": "string"
          },
          "changeType": "string",
          "comment": "string",
          "timestamp": "date"
        }
      ],
      "pagination": {
        "page": 1,
        "limit": 10,
        "total": 200,
        "pages": 20
      }
    }
  }
  ```

#### 获取历史记录详情
- **URL**: `GET /history/{id}`
- **描述**: 获取特定历史记录的详细信息
- **路径参数**:
  - `id`: 历史记录ID
- **响应**:
  ```json
  {
    "success": true,
    "data": {
      "history": {
        "_id": "ObjectId",
        "taskId": {
          "_id": "ObjectId",
          "title": "string",
          "description": "string",
          "status": "string",
          "priority": "string"
        },
        "field": "string",
        "oldValue": "any",
        "newValue": "any",
        "changedBy": {
          "_id": "ObjectId",
          "username": "string"
        },
        "changeType": "string",
        "comment": "string",
        "timestamp": "date"
      }
    }
  }
  ```

#### 获取历史时间线
- **URL**: `GET /history/timeline`
- **描述**: 获取用户任务变更历史时间线
- **查询参数**:
  - `startDate`: 开始日期
  - `endDate`: 结束日期
  - `limit`: 数量限制 (默认: 50, 最大: 200)
- **响应**:
  ```json
  {
    "success": true,
    "data": {
      "timeline": [
        {
          "_id": "ObjectId",
          "taskId": "ObjectId",
          "taskTitle": "string",
          "changeType": "string",
          "description": "string",
          "changedBy": {
            "_id": "ObjectId",
            "username": "string"
          },
          "timestamp": "date"
        }
      ]
    }
  }
  ```

### 4. 团队管理接口

#### 获取团队列表
- **URL**: `GET /teams`
- **描述**: 获取用户所属的团队列表
- **查询参数**:
  - `page`: 页码 (默认: 1)
  - `limit`: 每页数量 (默认: 10, 最大: 100)
- **响应**:
  ```json
  {
    "success": true,
    "data": {
      "teams": [
        {
          "_id": "ObjectId",
          "name": "string",
          "description": "string",
          "avatar": "string",
          "membersCount": "number",
          "createdAt": "date"
        }
      ],
      "pagination": {
        "page": 1,
        "limit": 10,
        "total": 5,
        "pages": 1
      }
    }
  }
  ```

#### 创建团队
- **URL**: `POST /teams`
- **描述**: 创建新团队
- **请求参数**:
  ```json
  {
    "name": "string",
    "description": "string",
    "avatar": "string"
  }
  ```
- **响应**:
  ```json
  {
    "success": true,
    "data": {
      "team": {
        "_id": "ObjectId",
        "name": "string",
        "description": "string",
        "avatar": "string",
        "members": [
          {
            "user": "ObjectId",
            "role": "owner",
            "joinedAt": "date"
          }
        ],
        "createdAt": "date",
        "updatedAt": "date"
      }
    }
  }
  ```

#### 获取团队详情
- **URL**: `GET /teams/{id}`
- **描述**: 获取特定团队的详细信息
- **路径参数**:
  - `id`: 团队ID
- **响应**:
  ```json
  {
    "success": true,
    "data": {
      "team": {
        "_id": "ObjectId",
        "name": "string",
        "description": "string",
        "avatar": "string",
        "members": [
          {
            "user": {
              "_id": "ObjectId",
              "username": "string",
              "avatar": "string"
            },
            "role": "string",
            "joinedAt": "date"
          }
        ],
        "createdAt": "date",
        "updatedAt": "date"
      }
    }
  }
  ```

#### 更新团队信息
- **URL**: `PUT /teams/{id}`
- **描述**: 更新团队信息
- **路径参数**:
  - `id`: 团队ID
- **请求参数**:
  ```json
  {
    "name": "string",
    "description": "string",
    "avatar": "string"
  }
  ```
- **响应**:
  ```json
  {
    "success": true,
    "data": {
      "team": {
        "_id": "ObjectId",
        "name": "string",
        "description": "string",
        "avatar": "string",
        "members": "array",
        "createdAt": "date",
        "updatedAt": "date"
      }
    }
  }
  ```

#### 删除团队
- **URL**: `DELETE /teams/{id}`
- **描述**: 删除指定团队（需要owner权限）
- **路径参数**:
  - `id`: 团队ID
- **响应**:
  ```json
  {
    "success": true,
    "message": "Team deleted successfully"
  }
  ```

#### 添加团队成员
- **URL**: `POST /teams/{id}/members`
- **描述**: 添加团队成员（需要admin或owner权限）
- **路径参数**:
  - `id`: 团队ID
- **请求参数**:
  ```json
  {
    "userId": "ObjectId",
    "role": "string"  // owner, admin, member
  }
  ```
- **响应**:
  ```json
  {
    "success": true,
    "message": "Member added successfully"
  }
  ```

#### 移除团队成员
- **URL**: `DELETE /teams/{id}/members/{userId}`
- **描述**: 移除团队成员（需要admin或owner权限）
- **路径参数**:
  - `id`: 团队ID
  - `userId`: 用户ID
- **响应**:
  ```json
  {
    "success": true,
    "message": "Member removed successfully"
  }
  ```

### 5. 用户管理接口

#### 获取用户详情
- **URL**: `GET /users/{id}`
- **描述**: 获取用户公开信息
- **路径参数**:
  - `id`: 用户ID
- **响应**:
  ```json
  {
    "success": true,
    "data": {
      "user": {
        "_id": "ObjectId",
        "username": "string",
        "avatar": "string",
        "bio": "string",
        "skills": ["string"],
        "location": "string",
        "createdAt": "date"
      }
    }
  }
  ```

#### 更新用户信息
- **URL**: `PUT /users/{id}`
- **描述**: 更新当前用户信息
- **路径参数**:
  - `id`: 用户ID (必须等于当前用户ID)
- **请求参数**:
  ```json
  {
    "username": "string",
    "avatar": "string",
    "firstName": "string",
    "lastName": "string",
    "bio": "string",
    "skills": ["string"],
    "location": "string",
    "preferences": {
      "theme": "string",
      "notifications": {
        "email": "boolean",
        "push": "boolean"
      },
      "language": "string"
    }
  }
  ```
- **响应**:
  ```json
  {
    "success": true,
    "data": {
      "user": {
        "_id": "ObjectId",
        "username": "string",
        "email": "string",
        "profile": {
          "avatar": "string",
          "firstName": "string",
          "lastName": "string",
          "bio": "string",
          "skills": ["string"],
          "location": "string"
        },
        "preferences": {
          "theme": "string",
          "notifications": {
            "email": "boolean",
            "push": "boolean"
          },
          "language": "string"
        },
        "createdAt": "date",
        "updatedAt": "date"
      }
    }
  }
  ```

#### 获取用户统计信息
- **URL**: `GET /users/{id}/statistics`
- **描述**: 获取用户任务统计信息
- **路径参数**:
  - `id`: 用户ID
- **响应**:
  ```json
  {
    "success": true,
    "data": {
      "statistics": {
        "_id": "ObjectId",
        "tasksCompleted": "number",
        "tasksCreated": "number",
        "completionRate": "number",
        "activeDays": ["date"],
        "lastActive": "date",
        "achievements": [
          {
            "id": "string",
            "name": "string",
            "description": "string",
            "earnedAt": "date"
          }
        ],
        "weeklyStats": [
          {
            "week": "date",
            "tasksCompleted": "number",
            "tasksCreated": "number"
          }
        ],
        "createdAt": "date",
        "updatedAt": "date"
      }
    }
  }
  ```

### 6. 统计和总结接口

#### 获取任务分析数据
- **URL**: `GET /analytics/tasks`
- **描述**: 获取任务分析数据
- **查询参数**:
  - `period`: 时间周期 (day, week, month)
  - `startDate`: 开始日期
  - `endDate`: 结束日期
- **响应**:
  ```json
  {
    "success": true,
    "data": {
      "analytics": {
        "totalTasks": "number",
        "completedTasks": "number",
        "completionRate": "number",
        "byStatus": {
          "created": "number",
          "in-progress": "number",
          "paused": "number",
          "completed": "number",
          "cancelled": "number"
        },
        "byPriority": {
          "low": "number",
          "medium": "number",
          "high": "number"
        },
        "dailyData": [
          {
            "date": "date",
            "created": "number",
            "completed": "number"
          }
        ]
      }
    }
  }
  ```

#### 获取任务总结
- **URL**: `GET /summaries/tasks`
- **描述**: 获取任务总结
- **查询参数**:
  - `period`: 时间周期 (daily, weekly, monthly)
  - `startDate`: 开始日期
  - `endDate`: 结束日期
- **响应**:
  ```json
  {
    "success": true,
    "data": {
      "summary": {
        "_id": "ObjectId",
        "period": "string",
        "startDate": "date",
        "endDate": "date",
        "tasksCompleted": "number",
        "tasksCreated": "number",
        "tasksCancelled": "number",
        "averageCompletionTime": "number",
        "completionRate": "number",
        "insights": [
          {
            "type": "string",
            "message": "string",
            "severity": "string"
          }
        ],
        "generatedAt": "date"
      }
    }
  }
  ```

#### 获取用户任务总结
- **URL**: `GET /users/{id}/summaries`
- **描述**: 获取指定用户的任务总结
- **路径参数**:
  - `id`: 用户ID
- **查询参数**:
  - `period`: 时间周期 (daily, weekly, monthly)
  - `limit`: 数量限制 (默认: 10)
- **响应**:
  ```json
  {
    "success": true,
    "data": {
      "summaries": [
        {
          "_id": "ObjectId",
          "period": "string",
          "startDate": "date",
          "endDate": "date",
          "tasksCompleted": "number",
          "tasksCreated": "number",
          "completionRate": "number",
          "insights": [
            {
              "type": "string",
              "message": "string",
              "severity": "string"
            }
          ],
          "generatedAt": "date"
        }
      ]
    }
  }
  ```

## WebSocket 实时更新

系统提供 WebSocket 连接以支持实时更新：

### 连接地址
```
wss://api.todoing.example.com/ws
```

### 事件类型
- `task.created`: 新任务创建
- `task.updated`: 任务更新
- `task.deleted`: 任务删除
- `history.created`: 新的历史记录创建
- `team.created`: 新团队创建
- `team.updated`: 团队更新
- `team.deleted`: 团队删除
- `team.member.added`: 团队成员添加
- `team.member.removed`: 团队成员移除

## 错误码定义

| 错误码 | 描述 |
|--------|------|
| VALIDATION_ERROR | 输入参数验证失败 |
| UNAUTHORIZED | 未授权访问 |
| FORBIDDEN | 权限不足 |
| NOT_FOUND | 资源未找到 |
| INTERNAL_ERROR | 服务器内部错误 |
| TASK_NOT_FOUND | 任务未找到 |
| TEAM_NOT_FOUND | 团队未找到 |
| USER_NOT_FOUND | 用户未找到 |
| INVALID_STATUS_TRANSITION | 无效的状态转换 |
| TEAM_NAME_EXISTS | 团队名称已存在 |
| ALREADY_TEAM_MEMBER | 已经是团队成员 |

## 限流策略

- 每个用户每分钟最多 1000 次 API 请求
- 每个 IP 地址每分钟最多 5000 次 API 请求