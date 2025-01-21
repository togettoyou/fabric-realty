#!/bin/bash

# 设置错误时立即退出
set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

echo -e "\n${RED}================================${NC}"
echo -e "${RED}   Fabric-Realty 一键卸载脚本${NC}"
echo -e "${RED}================================${NC}\n"

# 确认卸载
read -p "$(echo -e ${YELLOW}"警告：此操作将清除所有相关的容器、网络和数据。确定要继续吗？[y/N] "${NC})" confirm
if [[ $confirm != [yY] ]]; then
    log_info "操作已取消"
    exit 0
fi

# 停止并删除应用服务容器
log_info "清理应用服务..."
if [ -d "application" ]; then
    cd application
    if [ -f "docker-compose.yml" ]; then
        log_info "停止并删除应用服务容器..."
        docker-compose down -v || log_warning "应用服务清理过程中出现错误，继续执行..."
    fi
    cd ..
fi

# 调用网络卸载脚本
log_info "清理区块链网络..."
if [ -d "network" ]; then
    cd network
    if [ ! -f "./uninstall.sh" ]; then
        log_error "network/uninstall.sh 文件不存在！"
        exit 1
    fi

    log_info "执行 network/uninstall.sh..."
    ./uninstall.sh
    if [ $? -ne 0 ]; then
        log_error "区块链网络清理失败！"
        exit 1
    fi
    cd ..
fi

echo -e "\n${GREEN}================================${NC}"
echo -e "${GREEN}   卸载完成！   ${NC}"
echo -e "${GREEN}================================${NC}\n"

log_success "所有组件已成功清理"
log_info "如果您想重新安装，请运行 ./install.sh"
