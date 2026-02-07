#!/bin/bash

# Fish Music 部署脚本
# 使用方法: ./scripts/deploy.sh [environment]
# environment: dev (默认) | prod

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查必需的命令
check_requirements() {
    log_info "检查系统依赖..."

    if ! command -v docker &> /dev/null; then
        log_error "Docker 未安装，请先安装 Docker"
        exit 1
    fi

    if ! command -v docker-compose &> /dev/null; then
        log_error "Docker Compose 未安装，请先安装 Docker Compose"
        exit 1
    fi

    log_info "系统依赖检查通过"
}

# 检查配置文件
check_config() {
    log_info "检查配置文件..."

    if [ ! -f "config.yaml" ]; then
        log_warn "config.yaml 不存在，从模板创建..."
        cp config.yaml.example config.yaml
        log_warn "请编辑 config.yaml 填入正确的配置信息"
        log_warn "特别是 Bot Token 和数据库密码"
        exit 1
    fi

    log_info "配置文件检查通过"
}

# 创建临时目录
create_temp_dir() {
    log_info "创建临时目录..."
    mkdir -p tmp
    log_info "临时目录已创建: ./tmp"
}

# 构建镜像
build_images() {
    log_info "构建 Docker 镜像..."
    docker-compose build
    log_info "Docker 镜像构建完成"
}

# 启动服务
start_services() {
    log_info "启动服务..."
    docker-compose up -d

    # 等待服务启动
    log_info "等待服务启动..."
    sleep 5

    # 检查服务状态
    if docker-compose ps | grep -q "Exit"; then
        log_error "服务启动失败，请查看日志: make docker-logs"
        exit 1
    fi

    log_info "服务启动成功!"
}

# 显示服务信息
show_info() {
    log_info "============================================"
    log_info "Fish Music 部署完成!"
    log_info "============================================"
    echo ""
    log_info "服务地址:"
    echo "  - Web 管理端: http://localhost:9999"
    echo "  - Bot Token: 请查看 config.yaml"
    echo ""
    log_info "常用命令:"
    echo "  - 查看日志: make docker-logs"
    echo "  - 停止服务: make docker-down"
    echo "  - 重启服务: docker-compose restart"
    echo ""
    log_info "============================================"
}

# 主函数
main() {
    log_info "开始部署 Fish Music..."
    echo ""

    check_requirements
    check_config
    create_temp_dir
    build_images
    start_services
    show_info
}

# 运行主函数
main
