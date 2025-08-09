#!/bin/bash

# TodoIng Docker Compose 快速启动脚本
# 使用方法: ./deploy.sh [方案名] [操作] [选项]

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 打印带颜色的消息
print_message() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

# 显示帮助信息
show_help() {
    cat << EOF
TodoIng Docker Compose 部署工具

使用方法:
    $0 [方案名] [操作] [选项]

可用方案:
    nodejs      - Node.js 后端方案
    golang      - Golang 后端方案
    dev         - 全栈开发环境
    prod        - 生产环境
    micro       - 微服务方案

可用操作:
    up          - 启动服务
    down        - 停止服务
    restart     - 重启服务
    logs        - 查看日志
    ps          - 查看服务状态
    build       - 构建镜像
    pull        - 拉取镜像
    push        - 推送镜像到远程仓库
    exec        - 进入容器

可用选项:
    --profile   - 启用特定profile (如: grpc, monitoring, replica)
    --detach    - 后台运行 (默认)
    --build     - 强制重新构建
    --remove    - 停止时删除卷
    --follow    - 跟随日志输出

镜像使用说明:
    默认使用远程镜像 (axiu/todoing-go, axiu/todoing, axiu/todoing-frontend)
    使用 --build 选项或设置环境变量为空值可以本地构建
    
    镜像环境变量:
    GOLANG_BACKEND_IMAGE=axiu/todoing-go:latest
    NODEJS_BACKEND_IMAGE=axiu/todoing:latest  
    FRONTEND_IMAGE=axiu/todoing-frontend:latest

示例:
    $0 golang up                    # 启动 Golang 方案 (使用远程镜像)
    $0 golang up --build            # 启动 Golang 方案 (本地构建)
    $0 dev up --profile nodejs      # 启动开发环境的 Node.js 后端
    $0 prod up --profile replica    # 启动生产环境(包含副本)
    $0 golang logs backend-golang   # 查看 Golang 后端日志
    $0 micro exec auth-service bash # 进入认证服务容器
    $0 golang pull                  # 拉取最新镜像
    $0 golang build                 # 构建本地镜像

EOF
}

# 获取 docker-compose 文件路径
get_compose_file() {
    local scheme=$1
    case $scheme in
        nodejs)
            echo "docker-compose.nodejs.yml"
            ;;
        golang)
            echo "docker-compose.golang.yml"
            ;;
        dev)
            echo "docker-compose.dev-full.yml"
            ;;
        prod)
            echo "docker-compose.prod.yml"
            ;;
        micro)
            echo "docker-compose.microservices.yml"
            ;;
        *)
            print_message $RED "错误: 未知的方案 '$scheme'"
            show_help
            exit 1
            ;;
    esac
}

# 检查环境变量文件
check_env_file() {
    if [ ! -f ".env" ]; then
        print_message $YELLOW "警告: .env 文件不存在"
        if [ -f ".env.example" ]; then
            print_message $BLUE "正在从 .env.example 创建 .env 文件..."
            cp .env.example .env
            print_message $GREEN "已创建 .env 文件，请根据需要修改配置"
        else
            print_message $RED "错误: .env.example 文件也不存在"
            exit 1
        fi
    fi
}

# 检查 Docker 和 Docker Compose
check_docker() {
    if ! command -v docker &> /dev/null; then
        print_message $RED "错误: Docker 未安装"
        exit 1
    fi

    if ! command -v docker-compose &> /dev/null; then
        print_message $RED "错误: Docker Compose 未安装"
        exit 1
    fi
}

# 执行 docker-compose 命令
execute_compose() {
    local compose_file=$1
    local operation=$2
    shift 2
    local args=("$@")

    print_message $BLUE "执行: docker-compose -f $compose_file $operation ${args[*]}"
    docker-compose -f "$compose_file" "$operation" "${args[@]}"
}

# 主函数
main() {
    # 检查参数
    if [ $# -lt 2 ]; then
        show_help
        exit 1
    fi

    local scheme=$1
    local operation=$2
    shift 2

    # 解析选项
    local profile=""
    local compose_args=()
    local detach_flag="-d"
    local build_flag=""
    local remove_flag=""
    local follow_flag=""

    while [[ $# -gt 0 ]]; do
        case $1 in
            --profile)
                profile="$2"
                shift 2
                ;;
            --detach)
                detach_flag="-d"
                shift
                ;;
            --no-detach)
                detach_flag=""
                shift
                ;;
            --build)
                build_flag="--build"
                shift
                ;;
            --remove)
                remove_flag="-v"
                shift
                ;;
            --follow)
                follow_flag="-f"
                shift
                ;;
            -h|--help)
                show_help
                exit 0
                ;;
            *)
                compose_args+=("$1")
                shift
                ;;
        esac
    done

    # 检查环境
    check_docker
    check_env_file

    # 获取 compose 文件
    local compose_file
    compose_file=$(get_compose_file "$scheme")

    # 构建命令参数
    local cmd_args=()
    
    # 添加 profile
    if [ -n "$profile" ]; then
        cmd_args+=(--profile "$profile")
    fi

    case $operation in
        up)
            cmd_args+=("$detach_flag")
            if [ -n "$build_flag" ]; then
                cmd_args+=("$build_flag")
            fi
            cmd_args+=("${compose_args[@]}")
            print_message $GREEN "正在启动 $scheme 方案..."
            execute_compose "$compose_file" "up" "${cmd_args[@]}"
            print_message $GREEN "服务启动完成!"
            
            # 显示访问地址
            case $scheme in
                nodejs)
                    print_message $BLUE "访问地址:"
                    print_message $BLUE "  应用: http://localhost"
                    print_message $BLUE "  API: http://localhost:5001/api"
                    ;;
                golang)
                    print_message $BLUE "访问地址:"
                    print_message $BLUE "  应用: http://localhost"
                    print_message $BLUE "  API: http://localhost:5004/api"
                    print_message $BLUE "  Swagger: http://localhost:5004/swagger/"
                    ;;
                dev)
                    print_message $BLUE "访问地址:"
                    print_message $BLUE "  前端开发: http://localhost:3000"
                    print_message $BLUE "  Mongo Express: http://localhost:8081"
                    print_message $BLUE "  Redis Commander: http://localhost:8082"
                    print_message $BLUE "  MailHog: http://localhost:8025"
                    ;;
                prod)
                    print_message $BLUE "访问地址:"
                    print_message $BLUE "  应用: http://localhost"
                    print_message $BLUE "  Grafana: http://localhost:3001"
                    ;;
                micro)
                    print_message $BLUE "访问地址:"
                    print_message $BLUE "  API 网关: http://localhost:8000"
                    print_message $BLUE "  Consul: http://localhost:8500"
                    print_message $BLUE "  RabbitMQ: http://localhost:15672"
                    ;;
            esac
            ;;
        down)
            cmd_args+=("$remove_flag")
            cmd_args+=("${compose_args[@]}")
            print_message $YELLOW "正在停止 $scheme 方案..."
            execute_compose "$compose_file" "down" "${cmd_args[@]}"
            print_message $GREEN "服务已停止!"
            ;;
        restart)
            print_message $YELLOW "正在重启 $scheme 方案..."
            execute_compose "$compose_file" "restart" "${compose_args[@]}"
            print_message $GREEN "服务已重启!"
            ;;
        logs)
            cmd_args+=("$follow_flag")
            cmd_args+=("${compose_args[@]}")
            execute_compose "$compose_file" "logs" "${cmd_args[@]}"
            ;;
        ps)
            execute_compose "$compose_file" "ps" "${compose_args[@]}"
            ;;
        build)
            cmd_args+=("${compose_args[@]}")
            print_message $BLUE "正在构建镜像..."
            execute_compose "$compose_file" "build" "${cmd_args[@]}"
            print_message $GREEN "镜像构建完成!"
            ;;
        pull)
            cmd_args+=("${compose_args[@]}")
            print_message $BLUE "正在拉取镜像..."
            execute_compose "$compose_file" "pull" "${cmd_args[@]}"
            print_message $GREEN "镜像拉取完成!"
            ;;
        exec)
            if [ ${#compose_args[@]} -lt 2 ]; then
                print_message $RED "错误: exec 需要指定服务名和命令"
                print_message $BLUE "示例: $0 $scheme exec backend-golang bash"
                exit 1
            fi
            execute_compose "$compose_file" "exec" "${compose_args[@]}"
            ;;
        *)
            print_message $RED "错误: 未知的操作 '$operation'"
            show_help
            exit 1
            ;;
    esac
}

# 捕获 Ctrl+C
trap 'print_message $YELLOW "\n操作已取消"' INT

# 运行主函数
main "$@"
