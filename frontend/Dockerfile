# 使用官方Node.js运行时作为基础镜像进行构建
FROM node:18-alpine AS build

# 设置工作目录
WORKDIR /app

# 复制package.json和package-lock.json（如果存在）
COPY package*.json ./

# 安装所有依赖
RUN npm ci

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