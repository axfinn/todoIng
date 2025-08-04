# 配置管理

系统支持多种配置选项，可以通过环境变量进行控制。

## 后端配置

### 基础配置
| 变量名 | 描述 | 默认值 | 示例 |
|--------|------|--------|------|
| `PORT` | 服务端口 | `5001` | `5001` |
| `MONGO_URI` | MongoDB连接字符串 | 无 | `mongodb://localhost:27017/todoing` |
| `JWT_SECRET` | JWT密钥 | 无 | `my_jwt_secret` |
| `NODE_ENV` | 运行环境 | `development` | `production` |

### 安全配置
| 变量名 | 描述 | 默认值 | 示例 |
|--------|------|--------|------|
| `DISABLE_REGISTRATION` | 是否禁用注册功能 | `false` | `true` |
| `ENABLE_CAPTCHA` | 是否启用登录验证码 | `false` | `true` |
| `ENABLE_EMAIL_VERIFICATION` | 是否启用邮箱验证码功能 | `false` | `true` |

### 邮件配置（当启用邮箱验证码功能时必需）
| 变量名 | 描述 | 默认值 | 示例 |
|--------|------|--------|------|
| `EMAIL_HOST` | SMTP服务器地址 | 无 | `smtp.gmail.com` |
| `EMAIL_PORT` | SMTP端口 | `587` | `465` |
| `EMAIL_SECURE` | 是否使用安全连接 | `false` | `true` |
| `EMAIL_USER` | 邮箱用户名 | 无 | `your_email@gmail.com` |
| `EMAIL_PASS` | 邮箱密码或应用专用密码 | 无 | `your_password` |
| `EMAIL_FROM` | 发件人邮箱 | 与EMAIL_USER相同 | `your_email@gmail.com` |

### 默认用户配置
| 变量名 | 描述 | 默认值 | 示例 |
|--------|------|--------|------|
| `DEFAULT_ADMIN_EMAIL` | 默认管理员邮箱 | 无 | `admin@example.com` |
| `DEFAULT_ADMIN_PASSWORD` | 默认管理员密码 | 无 | `admin123` |
| `DEFAULT_ADMIN_NAME` | 默认管理员姓名 | `Admin` | `Administrator` |

## 前端配置

### 基础配置
| 变量名 | 描述 | 默认值 | 示例 |
|--------|------|--------|------|
| `VITE_API_BASE_URL` | 后端API基础URL | `/api` | `http://localhost:5001/api` |

### 功能配置
| 变量名 | 描述 | 默认值 | 示例 |
|--------|------|--------|------|
| `VITE_DISABLE_REGISTRATION` | 是否禁用注册功能 | `false` | `true` |
| `VITE_ENABLE_CAPTCHA` | 是否启用登录验证码 | `false` | `true` |
| `VITE_ENABLE_EMAIL_VERIFICATION` | 是否启用邮箱验证码功能 | `false` | `true` |

## 配置示例

### 后端 .env 文件示例
```env
# 基础配置
PORT=5001
MONGO_URI=mongodb://localhost:27017/todoing
JWT_SECRET=my_jwt_secret

# 安全配置
DISABLE_REGISTRATION=false
ENABLE_CAPTCHA=true
ENABLE_EMAIL_VERIFICATION=true

# 邮件配置（当启用邮箱验证码功能时必需）
EMAIL_HOST=smtp.gmail.com
EMAIL_PORT=587
EMAIL_SECURE=false
EMAIL_USER=your_email@gmail.com
EMAIL_PASS=your_app_password
EMAIL_FROM=your_email@gmail.com

# 默认管理员用户
DEFAULT_ADMIN_EMAIL=admin@example.com
DEFAULT_ADMIN_PASSWORD=admin123
DEFAULT_ADMIN_NAME=Admin
```

### 前端 .env 文件示例
```env
# 基础配置
VITE_API_BASE_URL=http://localhost:5001/api

# 功能配置
VITE_DISABLE_REGISTRATION=false
VITE_ENABLE_CAPTCHA=true
VITE_ENABLE_EMAIL_VERIFICATION=true
```

## 配置说明

### 注册控制功能
当 `DISABLE_REGISTRATION=true` 时，注册接口将被禁用，新用户无法通过常规注册流程创建账户。

### 登录验证码功能
当 `ENABLE_CAPTCHA=true`（后端）和 `VITE_ENABLE_CAPTCHA=true`（前端）时，登录界面会显示验证码输入框和获取验证码按钮。

用户需要先点击"获取验证码"按钮，然后输入显示的验证码进行登录。

### 邮箱验证码功能
当 `ENABLE_EMAIL_VERIFICATION=true`（后端）和 `VITE_ENABLE_EMAIL_VERIFICATION=true`（前端）时，启用邮箱验证码功能：

1. **注册时邮箱验证码**：
   - 用户在注册时需要提供邮箱地址并获取验证码
   - 输入收到的验证码完成注册流程
   - 可与图片验证码同时使用以增强安全性

2. **登录时邮箱验证码**：
   - 用户可以选择使用邮箱验证码登录而无需密码
   - 点击"邮箱验证码登录"切换登录方式
   - 获取并输入验证码即可登录

注意：启用邮箱验证码功能需要正确配置邮件服务器参数。