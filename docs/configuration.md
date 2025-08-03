# 系统配置管理

## 概述

本文档描述了系统的配置选项，包括注册控制、默认用户和验证码功能的开关配置。

## 环境变量配置

系统通过环境变量来控制各种功能开关，以下是可以配置的环境变量：

| 变量名 | 描述 | 默认值 | 可选值 |
|--------|------|--------|--------|
| `DISABLE_REGISTRATION` | 是否禁用注册功能 | `false` | `true`/`false` |
| `ENABLE_CAPTCHA` | 是否启用登录验证码 | `false` | `true`/`false` |
| `DEFAULT_USERNAME` | 默认用户名 | 无 | 字符串 |
| `DEFAULT_EMAIL` | 默认用户邮箱 | 无 | 字符串 |
| `DEFAULT_PASSWORD` | 默认用户密码 | 无 | 字符串 |

## 功能说明

### 1. 禁用注册功能

当 `DISABLE_REGISTRATION=true` 时，系统将关闭用户注册接口，新用户无法通过常规注册流程创建账户。

### 2. 默认用户

当配置了默认用户信息时，系统会在启动时检查是否存在该用户，如果不存在则自动创建。这允许在禁用注册的情况下，通过默认用户访问系统。

### 3. 登录验证码

当 `ENABLE_CAPTCHA=true` 时，系统会在登录界面显示验证码，用户需要正确输入验证码才能登录。

## 使用示例

在 `.env` 文件中添加以下配置：

```
# 禁用注册功能
DISABLE_REGISTRATION=true

# 启用登录验证码
ENABLE_CAPTCHA=true

# 默认用户配置
DEFAULT_USERNAME=admin
DEFAULT_EMAIL=admin@example.com
DEFAULT_PASSWORD=securepassword
```