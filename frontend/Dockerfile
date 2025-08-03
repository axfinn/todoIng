# 使用官方Node.js运行时作为构建环境
FROM node:18-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制package.json和package-lock.json（如果存在）
COPY package*.json ./

# 安装依赖
RUN npm install

# 复制应用源代码
COPY . .

# 构建应用
RUN npm run build

# 使用nginx作为生产环境服务器
FROM nginx:alpine

# 复制构建产物到nginx服务器
COPY --from=builder /app/dist /usr/share/nginx/html

# 复制自定义nginx配置文件
COPY nginx.conf /etc/nginx/nginx.conf

# 暴露端口
EXPOSE 80

# 启动nginx服务器
CMD ["nginx", "-g", "daemon off;"]