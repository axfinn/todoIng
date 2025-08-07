# Backend 和 Backend-go API 路由对比表

## 总览
- **Backend (Node.js)**: 基于 Express.js 的 REST API
- **Backend-go (Go)**: 基于 Gin 框架的 REST API

## API 路由对比表

### 认证相关 (Auth Routes)

| 端点 | HTTP方法 | Backend (Node.js) | Backend-go | 说明 |
|-----|---------|------------------|------------|------|
| `/api/auth/me` | GET | ✅ 实现 | ✅ 实现 | 获取当前用户信息 |
| `/api/auth/register` | POST | ✅ 实现 | ✅ 实现 | 用户注册 |
| `/api/auth/login` | POST | ✅ 实现 (支持密码+验证码/邮箱验证码) | ✅ 实现 (仅支持密码) | 用户登录 |
| `/api/auth/captcha` | GET | ✅ 实现 (SVG) | ✅ 实现 (PNG) | 生成图形验证码 |
| `/api/auth/send-email-code` | POST | ✅ 实现 | ✅ 实现 (未完整实现) | 发送注册邮箱验证码 |
| `/api/auth/send-login-email-code` | POST | ✅ 实现 | ✅ 实现 (未完整实现) | 发送登录邮箱验证码 |
| `/api/auth/verify-captcha` | POST | ✅ 实现 | ❌ 缺失 | 验证图形验证码 |

### 任务相关 (Task Routes)

| 端点 | HTTP方法 | Backend (Node.js) | Backend-go | 说明 |
|-----|---------|------------------|------------|------|
| `/api/tasks` | GET | ✅ 实现 | ✅ 实现 | 获取任务列表 |
| `/api/tasks/:id` | GET | ✅ 实现 | ✅ 实现 | 获取单个任务 |
| `/api/tasks` | POST | ✅ 实现 | ✅ 实现 | 创建任务 |
| `/api/tasks/:id` | PUT | ✅ 实现 (支持评论) | ✅ 实现 (不支持评论) | 更新任务 |
| `/api/tasks/:id` | DELETE | ✅ 实现 | ✅ 实现 | 删除任务 |
| `/api/tasks/:id/assign` | POST | ❌ 不存在 | ✅ 实现 | 分配任务（Go版独有） |
| `/api/tasks/:id/comments` | POST | ❌ 不存在 | ✅ 实现 | 添加评论（Go版独立接口） |
| `/api/tasks/export/all` | GET | ✅ 实现 | ❌ 缺失 | 导出所有任务 |
| `/api/tasks/import` | POST | ✅ 实现 | ❌ 缺失 | 导入任务 |

### 报告相关 (Report Routes)

| 端点 | HTTP方法 | Backend (Node.js) | Backend-go | 说明 |
|-----|---------|------------------|------------|------|
| `/api/reports` | GET | ✅ 实现 | ✅ 实现 (模拟) | 获取报告列表 |
| `/api/reports/:id` | GET | ✅ 实现 (带任务详情) | ✅ 实现 (模拟) | 获取单个报告 |
| `/api/reports/generate` | POST | ✅ 实现 | ❌ 缺失 | 生成报告 |
| `/api/reports/:id/polish` | POST | ✅ 实现 (OpenAI/自定义AI) | ❌ 缺失 | AI润色报告 |
| `/api/reports/:id` | DELETE | ✅ 实现 | ✅ 实现 (模拟) | 删除报告 |
| `/api/reports/:id/export/:format` | GET | ✅ 实现 | ❌ 缺失 | 导出报告 |

### 其他端点

| 端点 | HTTP方法 | Backend (Node.js) | Backend-go | 说明 |
|-----|---------|------------------|------------|------|
| `/swagger/*` | GET | ❌ 不存在 | ✅ 实现 | Swagger文档（Go版独有） |

## 主要差异分析

### 1. 认证功能差异
- **Node.js版**：
  - 支持密码登录和邮箱验证码登录两种方式
  - 支持图形验证码验证（独立接口）
  - 完整的邮箱验证码功能（发送和验证）
  - 支持注册功能开关（DISABLE_REGISTRATION环境变量）

- **Go版**：
  - 仅支持密码登录
  - 验证码生成但无独立验证接口
  - 邮箱验证码接口存在但未完整实现（仅返回模拟响应）
  - 缺少注册功能开关

### 2. 任务功能差异
- **Node.js版**：
  - 评论功能集成在任务更新接口中
  - 支持任务导入导出功能
  - 支持批量操作

- **Go版**：
  - 评论功能独立为单独接口
  - 新增任务分配独立接口
  - 缺少导入导出功能

### 3. 报告功能差异
- **Node.js版**：
  - 完整的报告生成功能（支持日报、周报、月报）
  - AI润色功能支持多种AI服务商
  - 支持多格式导出（MD、TXT）
  - 报告包含详细的任务时间线和评论

- **Go版**：
  - 基本的CRUD操作已定义但多为模拟实现
  - 缺少实际的报告生成逻辑
  - 缺少AI润色功能
  - 缺少导出功能

### 4. 技术实现差异
- **Node.js版**：
  - 使用SVG生成验证码
  - 使用内存存储验证码和邮箱验证码
  - 支持环境变量配置各种功能开关

- **Go版**：
  - 使用PNG生成验证码（7段显示数字）
  - 集成Swagger文档
  - 使用Protocol Buffers定义数据结构

## 建议优先补充的功能

### 高优先级
1. **认证相关**：
   - 实现 `/api/auth/verify-captcha` 接口
   - 完善邮箱验证码功能的实际实现
   - 支持邮箱验证码登录

2. **任务相关**：
   - 实现 `/api/tasks/export/all` 导出功能
   - 实现 `/api/tasks/import` 导入功能

3. **报告相关**：
   - 实现 `/api/reports/generate` 报告生成功能
   - 实现实际的数据库操作替代模拟响应

### 中优先级
1. **报告相关**：
   - 实现 `/api/reports/:id/polish` AI润色功能
   - 实现 `/api/reports/:id/export/:format` 导出功能

2. **认证相关**：
   - 添加注册功能开关支持

### 低优先级
1. 统一两个版本的API响应格式
2. 统一错误处理机制
3. 完善单元测试覆盖

## 数据结构差异

### 任务状态映射
- Node.js: `"To Do"`, `"In Progress"`, `"Done"`
- Go: `"TO_DO"`, `"IN_PROGRESS"`, `"DONE"`

### 任务优先级映射
- Node.js: `"Low"`, `"Medium"`, `"High"`
- Go: `"LOW"`, `"MEDIUM"`, `"HIGH"`

### 响应格式差异
- Node.js: 通常返回 `{ msg: "message" }` 或 `{ errors: [...] }`
- Go: 通常返回 `{ message: "message" }` 或 `{ error: "error" }`