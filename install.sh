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

# 检查命令是否存在
check_command() {
    if ! command -v $1 &> /dev/null; then
        log_error "$1 命令未找到，请确保已安装必要的依赖"
        exit 1
    fi
}

# 检查必要的命令
log_info "检查必要的依赖..."
check_command docker
check_command docker-compose

# 显示欢迎信息
echo -e "\n${GREEN}================================${NC}"
echo -e "${GREEN}   Fabric-Realty 一键安装脚本${NC}"
echo -e "${GREEN}================================${NC}\n"

# 部署区块链网络
log_info "开始部署区块链网络..."
cd network
if [ ! -f "./install.sh" ]; then
    log_error "network/install.sh 文件不存在！"
    exit 1
fi

log_info "执行 network/install.sh..."
./install.sh
if [ $? -ne 0 ]; then
    log_error "区块链网络部署失败！"
    exit 1
fi
log_success "区块链网络部署完成"

# 返回项目根目录
cd ..

# 启动应用服务
log_info "开始启动应用服务..."
cd application
if [ ! -f "docker-compose.yml" ]; then
    log_error "application/docker-compose.yml 文件不存在！"
    exit 1
fi

log_info "执行 docker-compose up -d..."
docker-compose up -d
if [ $? -ne 0 ]; then
    log_error "应用服务启动失败！"
    exit 1
fi
log_success "应用服务启动完成"

# 检查服务状态
log_info "检查服务状态..."
sleep 5
if [ "$(docker-compose ps -q | wc -l)" -gt 0 ]; then
    log_success "所有服务已成功启动"
else
    log_error "部分服务可能未正常启动，请检查 docker-compose logs"
    exit 1
fi

echo -e "\n${GREEN}================================${NC}"
echo -e "${GREEN}   安装部署完成！   ${NC}"
echo -e "${GREEN}================================${NC}\n"
