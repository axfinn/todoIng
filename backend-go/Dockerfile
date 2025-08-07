# 使用官方Go镜像作为基础镜像
FROM golang:1.23-alpine AS builder

# 安装必要的工具
RUN apk add --no-cache \
    protobuf-dev

# 设置工作目录
WORKDIR /app

# 复制go mod和sum文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 安装protobuf生成工具
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 生成protobuf代码
RUN mkdir -p proto/gen && \
    protoc -I./proto --go_out=./proto/gen --go-grpc_out=./proto/gen proto/*.proto

# 构建应用
RUN go build -o bin/todoing ./src/main.go

# 使用更小的基础镜像
FROM alpine:latest

# 安装ca-certificates以支持HTTPS请求
RUN apk --no-cache add ca-certificates

# 设置工作目录
WORKDIR /root/

# 从builder阶段复制编译好的二进制文件
COPY --from=builder /app/bin/todoing .

# 复制proto文件（如果需要在运行时访问）
COPY --from=builder /app/proto ./proto

# 暴露端口
EXPOSE 5001


# 运行应用
CMD ["./todoing"]