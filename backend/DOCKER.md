# 后端 Docker 配置

## Docker 镜像概述

后端服务使用 Node.js 18 Alpine 作为基础镜像，构建轻量级 Docker 镜像用于生产环境部署。

## Dockerfile 说明

```dockerfile
# 使用官方Node.js运行时作为基础镜像
FROM node:18-alpine

# 设置工作目录
WORKDIR /app

# 复制package.json和package-lock.json（如果存在）
COPY package*.json ./

# 安装依赖
RUN npm install

# 复制应用源代码
COPY . .

# 暴露端口
EXPOSE 5000

# 启动应用
CMD ["npm", "start"]
```

### 构建阶段

1. 使用 `node:18-alpine` 作为基础镜像，这是一个轻量级的 Node.js 运行环境
2. 设置工作目录为 `/app`
3. 复制 `package.json` 和 `package-lock.json` 到容器中
4. 运行 `npm install` 安装所有依赖
5. 复制所有源代码到容器中
6. 暴露端口 5000（应用默认运行端口）
7. 设置容器启动命令为 `npm start`

## 环境变量

后端容器需要以下环境变量：

- `MONGODB_URI` - MongoDB 连接字符串
- `JWT_SECRET` - JWT 密钥
- `PORT` - 应用运行端口（默认 5000）

## 构建镜像

在 backend 目录下运行以下命令构建镜像：

```bash
docker build -t todoing-backend .
```

## 运行容器

使用以下命令运行容器：

```bash
docker run -d \
  --name todoing-backend \
  -p 5001:5000 \
  -e MONGODB_URI=mongodb://mongodb:27017/todoing \
  -e JWT_SECRET=your_jwt_secret_here \
  -e PORT=5000 \
  todoing-backend
```

注意：在生产环境中，应该使用 Docker Compose 或 Kubernetes 来管理容器，确保 MongoDB 服务可用。

## 在 Docker Compose 中使用

后端服务通常与 MongoDB 和前端服务一起在 Docker Compose 中运行：

```yaml
services:
  backend:
    build: ./backend
    ports:
      - "5001:5000"
    environment:
      - MONGODB_URI=mongodb://mongodb:27017/todoing
      - JWT_SECRET=your_jwt_secret_here
      - PORT=5000
    depends_on:
      - mongodb
```

## 故障排除

### 常见问题

1. **数据库连接失败**：
   - 检查 `MONGODB_URI` 环境变量是否正确
   - 确保 MongoDB 服务正在运行且可访问
   - 验证网络连接和防火墙设置

2. **端口冲突**：
   - 修改 `-p` 参数中的主机端口，例如 `-p 5002:5000`

3. **依赖安装失败**：
   - 清理构建缓存后重新构建
   - 检查网络连接是否正常

### 日志查看

使用以下命令查看容器日志：

```bash
docker logs todoing-backend
```