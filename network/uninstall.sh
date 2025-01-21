#!/bin/bash

###########################################
# Hyperledger Fabric 网络清理脚本
# 版本: 1.0
# 描述: 清理Fabric网络环境，包括容器、镜像、数据卷等
###########################################

set -e  # 遇到错误立即退出
set -u  # 使用未定义的变量时报错

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO] $1${NC}"
}

log_success() {
    echo -e "${GREEN}[SUCCESS] $1${NC}"
}

log_error() {
    echo -e "${RED}[ERROR] $1${NC}"
}

# 错误处理函数
handle_error() {
    local exit_code=$?
    local step_name=$1
    log_error "步骤失败: $step_name"
    log_error "错误代码: $exit_code"
    exit $exit_code
}

# 检查docker服务状态
check_docker_service() {
    if ! docker info &> /dev/null; then
        log_error "Docker 服务未运行，请先启动 Docker"
        exit 1
    fi
    log_success "Docker 服务运行正常"
}

# 清理docker容器
clean_containers() {
    log_info "清理相关Docker容器..."
    docker-compose down --volumes --remove-orphans || handle_error "停止并删除容器"
    docker rm -f $(docker ps -a | grep "dev-peer*" | awk '{print $1}') 2>/dev/null || true
    log_success "Docker容器清理完成"
}

# 清理链码容器和镜像
clean_chaincode() {
    log_info "清理链码相关容器和镜像..."
    docker rmi -f $(docker images -a | grep "dev-peer*" | awk '{print $3}') 2>/dev/null || true
    log_success "链码清理完成"
}

# 清理数据文件
clean_files() {
    log_info "清理数据文件..."
    rm -rf config crypto-config data
    log_success "数据文件清理完成"
}

# 主程序
main() {
    log_info "开始清理Fabric网络环境..."

    # 检查Docker服务
    check_docker_service

    # 执行清理步骤
    clean_containers
    clean_chaincode
    clean_files

    log_success "Fabric网络环境清理完成！"
}

# 执行主程序
main "$@"
