## 🔍 Swagger接口扫描问题分析报告

### 问题现象
Swagger文档只显示了3个认证相关接口，缺少任务管理、报表和验证码接口。

### 根本原因
1. **swag扫描机制限制**：swag主要从main.go开始扫描，但我们的handler函数分布在不同的包中
2. **注释格式问题**：某些swagger注释可能不符合swag的解析规则
3. **包导入问题**：swag无法发现没有被明确引用的handler函数

### 当前状态
- ✅ 认证接口：完整显示（登录、注册、v2版本）
- ❌ 任务管理：5个接口未显示（增删改查+导入导出）
- ❌ 报表管理：2个接口未显示（列表+详情）
- ❌ 验证码：2个接口未显示（生成+验证）

### 解决方案
有几种方法可以解决这个问题：

#### 方案1：在main.go中明确引用所有handler（推荐）
```go
// 在main.go中添加
import (
    _ "github.com/axfinn/todoIng/backend-go/internal/api" // 强制导入所有handlers
)

// 在init函数中引用类型，确保swag能发现
func init() {
    _ = api.TaskDeps{}
    _ = api.ReportDeps{}
    _ = api.CaptchaDeps{}
}
```

#### 方案2：使用swag的注册机制
在每个handler文件中添加专门的文档注册函数。

#### 方案3：合并所有swagger注释到一个文件
创建一个专门的docs.go文件包含所有接口定义。

### 当前接口完整列表

#### 认证模块 (✅ 已显示)
- POST `/api/auth/register` - 用户注册
- POST `/api/auth/login` - 用户登录
- POST `/api/v2/auth/login` - v2登录

#### 任务管理模块 (❌ 未显示)
- POST `/api/tasks` - 创建任务
- GET `/api/tasks` - 获取任务列表
- GET `/api/tasks/{id}` - 获取任务详情
- PUT `/api/tasks/{id}` - 更新任务
- DELETE `/api/tasks/{id}` - 删除任务

#### 报表管理模块 (❌ 未显示)
- GET `/api/reports` - 获取报表列表
- GET `/api/reports/{id}` - 获取报表详情

#### 验证码模块 (❌ 未显示)
- GET `/api/auth/captcha` - 生成验证码
- POST `/api/auth/verify-captcha` - 验证验证码

### 立即解决方案
我现在实施方案1，确保所有接口都能被swagger正确识别和显示。
