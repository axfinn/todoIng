# 使用Node.js 18作为构建环境
FROM node:18-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制package.json和package-lock.json
COPY package*.json ./

# 安装所有依赖
RUN npm install

# 安装esbuild-wasm作为跨平台兼容的替代方案
RUN npm install esbuild-wasm

# 复制源代码
COPY . .

# 允许通过构建参数控制 Vite 环境变量 (这些变量只在构建期生效)
ARG VITE_API_BASE_URL=/api
ARG VITE_ENABLE_CAPTCHA=false
ARG VITE_ENABLE_EMAIL_VERIFICATION=false
ARG VITE_DISABLE_REGISTRATION=false

# 导出为环境变量以便 Vite 打包时内联
ENV VITE_API_BASE_URL=${VITE_API_BASE_URL} \
	VITE_ENABLE_CAPTCHA=${VITE_ENABLE_CAPTCHA} \
	VITE_ENABLE_EMAIL_VERIFICATION=${VITE_ENABLE_EMAIL_VERIFICATION} \
	VITE_DISABLE_REGISTRATION=${VITE_DISABLE_REGISTRATION}

# 构建应用
RUN npm run build

# 使用nginx作为生产环境
FROM nginx:alpine

# 复制构建产物到nginx目录
COPY --from=builder /app/dist /usr/share/nginx/html

# 复制nginx配置
COPY nginx.conf /etc/nginx/nginx.conf

# 暴露80端口
EXPOSE 80

# 启动nginx
CMD ["nginx", "-g", "daemon off;"]