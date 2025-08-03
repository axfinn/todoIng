# 前端 Docker 配置

## Docker 镜像概述

前端服务使用多阶段构建过程创建 Docker 镜像。首先使用 Node.js 环境构建生产版本的应用，然后将构建产物复制到 Nginx 服务器中运行。

## Dockerfile 说明

```dockerfile
# 使用官方Node.js运行时作为基础镜像进行构建
FROM node:18-alpine AS build

# 设置工作目录
WORKDIR /app

# 复制package.json和package-lock.json（如果存在）
COPY package*.json ./

# 安装所有依赖
RUN npm install

# 复制应用源代码
COPY . .

# 构建生产版本
RUN npm run build

# 使用nginx作为生产服务器
FROM nginx:alpine

# 复制构建产物到nginx服务器
COPY --from=build /app/dist /usr/share/nginx/html

# 复制nginx配置文件
COPY nginx.conf /etc/nginx/nginx.conf

# 暴露端口
EXPOSE 80

# 启动nginx
CMD ["nginx", "-g", "daemon off;"]
```

### 构建阶段

1. **构建阶段**：
   - 使用 `node:18-alpine` 作为构建环境
   - 设置工作目录为 `/app`
   - 复制 `package.json` 和 `package-lock.json` 到容器中
   - 运行 `npm install` 安装所有依赖
   - 复制所有源代码到容器中
   - 运行 `npm run build` 构建生产版本

2. **运行阶段**：
   - 使用 `nginx:alpine` 作为运行环境
   - 从构建阶段复制构建产物到 Nginx 静态文件目录
   - 复制自定义 Nginx 配置文件
   - 暴露端口 80
   - 设置容器启动命令为 Nginx

## Nginx 配置

前端使用 Nginx 作为 Web 服务器，并配置了 API 代理：

```nginx
# API代理配置
location /api {
    proxy_pass http://backend:5000;
    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection 'upgrade';
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
    proxy_cache_bypass $http_upgrade;
}
```

Nginx 配置支持：
- 静态文件缓存
- SPA 路由支持
- API 请求代理到后端服务

## 构建镜像

在 frontend 目录下运行以下命令构建镜像：

```bash
docker build -t todoing-frontend .
```

## 运行容器

使用以下命令运行容器：

```bash
docker run -d \
  --name todoing-frontend \
  -p 80:80 \
  todoing-frontend
```

注意：前端需要与后端服务通信，因此在使用 Docker Compose 时需要确保网络配置正确。

## 在 Docker Compose 中使用

前端服务通常与后端服务和 MongoDB 一起在 Docker Compose 中运行：

```yaml
services:
  frontend:
    build: ./frontend
    ports:
      - "80:80"
    depends_on:
      - backend
```

## 自定义配置

### 环境变量

前端在构建时使用 Vite 环境变量。创建 `.env` 文件来配置环境变量：

```env
VITE_API_URL=/api
```

### Nginx 配置

可以通过修改 [nginx.conf](file:///Volumes/M20/code/docs/todoIng/frontend/nginx.conf) 文件来自定义 Nginx 配置。

## 故障排除

### 常见问题

1. **API 请求失败**：
   - 检查 Nginx 配置中的代理设置
   - 确保后端服务正在运行且可访问
   - 验证 Docker 网络配置

2. **页面无法加载**：
   - 检查构建是否成功完成
   - 查看浏览器开发者工具中的错误信息
   - 确认 Nginx 配置中的静态文件路径

3. **静态资源缓存问题**：
   - 清除浏览器缓存
   - 修改 nginx.conf 中的缓存配置

### 日志查看

使用以下命令查看容器日志：

```bash
docker logs todoing-frontend
```

查看 Nginx 访问日志：

```bash
docker exec todoing-frontend cat /var/log/nginx/access.log
```

查看 Nginx 错误日志：

```bash
docker exec todoing-frontend cat /var/log/nginx/error.log
```