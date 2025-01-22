#!/bin/bash

###########################################
# Hyperledger Fabric 网络部署脚本
# 版本: 1.0
# 描述: 自动部署三组织六节点的Fabric网络
# 依赖:
#   - docker & docker-compose
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

# 时间统计函数
time_elapsed() {
    local start_time=$1
    local end_time=$(date +%s)
    local elapsed=$((end_time - start_time))
    local hours=$((elapsed / 3600))
    local minutes=$(((elapsed % 3600) / 60))
    local seconds=$((elapsed % 60))

    if [ $hours -gt 0 ]; then
        printf "%d小时%d分钟%d秒" $hours $minutes $seconds
    elif [ $minutes -gt 0 ]; then
        printf "%d分钟%d秒" $minutes $seconds
    else
        printf "%d秒" $seconds
    fi
}

# 步骤执行时间跟踪函数
execute_with_timer() {
    local step_name=$1
    local command=$2
    local start_time=$(date +%s)

    echo -e "${BLUE}[开始] $step_name...${NC}"
    eval "$command"
    local result=$?

    if [ $result -eq 0 ]; then
        echo -e "${GREEN}[完成] $step_name (耗时: $(time_elapsed $start_time))${NC}"
        return 0
    else
        echo -e "${RED}[失败] $step_name (耗时: $(time_elapsed $start_time))${NC}"
        return 1
    fi
}

# 等待操作完成函数
wait_for_completion() {
    local operation=$1
    local wait_time=$2
    local start_time=$(date +%s)

    echo -e "${BLUE}[等待] $operation...${NC}"
    sleep $wait_time
    echo -e "${GREEN}[完成] $operation (耗时: $(time_elapsed $start_time))${NC}"
}

# 进度显示函数
show_progress() {
    local current_step=$1
    local total_steps=16
    local step_name=$2
    local start_time=${3:-}  # 如果第三个参数未定义，则设为空

    # 定义步骤标签
    local step_tags=(
        ""                          # 占位，使索引从1开始
        "🔧 [环境]"                 # 步骤1
        "🧹 [清理]"                 # 步骤2
        "🛠️ [工具]"                 # 步骤3
        "🔑 [证书]"                 # 步骤4
        "📦 [创世]"                 # 步骤5
        "⚙️ [配置]"                 # 步骤6
        "⚓ [锚点]"                 # 步骤7
        "🚀 [启动]"                 # 步骤8
        "📝 [通道]"                 # 步骤9
        "🔗 [加入]"                 # 步骤10
        "📌 [更新]"                 # 步骤11
        "📦 [打包]"                 # 步骤12
        "💾 [安装]"                 # 步骤13
        "✅ [批准]"                 # 步骤14
        "📤 [提交]"                 # 步骤15
        "🔍 [验证]"                 # 步骤16
    )

    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    log_info "${step_tags[$current_step]} [步骤 $current_step/$total_steps] $step_name"
    if [ ! -z "${start_time}" ]; then
        echo -e "${BLUE}已耗时: $(time_elapsed $start_time)${NC}"
    fi
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
}

# 错误处理函数
handle_error() {
    local exit_code=$?
    local step_name=$1
    log_error "步骤失败: $step_name"
    log_error "错误代码: $exit_code"
    exit $exit_code
}

# 健康检查函数
check_prerequisites() {
    local prerequisites=("docker" "docker-compose")

    for cmd in "${prerequisites[@]}"; do
        if ! command -v $cmd &> /dev/null; then
            log_error "命令 '$cmd' 未找到。请确保已安装所有必需的组件。"
            exit 1
        fi
    done
    log_success "前置条件检查通过"
}

# 检查docker服务状态
check_docker_service() {
    if ! docker info &> /dev/null; then
        log_error "Docker 服务未运行，请先启动 Docker"
        exit 1
    fi
    log_success "Docker 服务运行正常"
}

# 清理环境
clean_environment() {
    local start_time=${1:-$(date +%s)}
    show_progress 2 "清理环境" $start_time
    ./uninstall.sh || handle_error "清理环境"
}

###########################################
# 配置参数
###########################################

# 等待时间配置（秒）
NETWORK_STARTUP_WAIT=10
CHAINCODE_INIT_WAIT=5

# 域名配置
DOMAIN="togettoyou.com"
ORG1_DOMAIN="org1.${DOMAIN}"
ORG2_DOMAIN="org2.${DOMAIN}"
ORG3_DOMAIN="org3.${DOMAIN}"
CLI_CONTAINER="cli.${DOMAIN}"

# CLI命令前缀
CLI_CMD="docker exec ${CLI_CONTAINER} bash -c"

# 基础路径配置
HYPERLEDGER_PATH="/etc/hyperledger"
CONFIG_PATH="${HYPERLEDGER_PATH}/config"
CRYPTO_PATH="${HYPERLEDGER_PATH}/crypto-config"

# 通道和链码配置
ChannelName="mychannel"
ChainCodeName="mychaincode"
Version="1.0.0"
Sequence="1"
CHAINCODE_PATH="/opt/gopath/src/chaincode"
CHAINCODE_PACKAGE="${CHAINCODE_PATH}/chaincode_${Version}.tar.gz"

# Order 配置
ORDERER1_ADDRESS="orderer1.${DOMAIN}:7050"
ORDERER_CA="${CRYPTO_PATH}/ordererOrganizations/${DOMAIN}/orderers/orderer1.${DOMAIN}/msp/tlscacerts/tlsca.${DOMAIN}-cert.pem"

# Org 配置
PEER_ORGS_MSP_PATH="${CRYPTO_PATH}/peerOrganizations"
CORE_PEER_TLS_ENABLED=true

# 生成节点配置函数
generate_peer_config() {
    local org=$1    # 组织编号
    local peer=$2   # 节点编号
    local org_domain="org${org}.${DOMAIN}"
    local peer_name="peer${peer}.${org_domain}"

    # 设置环境变量
    eval "ORG${org}_PEER${peer}_ADDRESS=\"${peer_name}:7051\""
    eval "ORG${org}_PEER${peer}_LOCALMSPID=\"Org${org}MSP\""
    eval "ORG${org}_PEER${peer}_MSPCONFIGPATH=\"${PEER_ORGS_MSP_PATH}/${org_domain}/users/Admin@${org_domain}/msp\""
    eval "ORG${org}_PEER${peer}_TLS_ROOTCERT_FILE=\"${PEER_ORGS_MSP_PATH}/${org_domain}/peers/${peer_name}/tls/ca.crt\""
    eval "ORG${org}_PEER${peer}_TLS_CERT_FILE=\"${PEER_ORGS_MSP_PATH}/${org_domain}/peers/${peer_name}/tls/server.crt\""
    eval "ORG${org}_PEER${peer}_TLS_KEY_FILE=\"${PEER_ORGS_MSP_PATH}/${org_domain}/peers/${peer_name}/tls/server.key\""
}

# 生成CLI配置函数
generate_cli_config() {
    local org=$1    # 组织编号
    local peer=$2   # 节点编号

    eval "Org${org}Peer${peer}Cli=\"CORE_PEER_ADDRESS=\${ORG${org}_PEER${peer}_ADDRESS} \\
CORE_PEER_LOCALMSPID=\${ORG${org}_PEER${peer}_LOCALMSPID} \\
CORE_PEER_MSPCONFIGPATH=\${ORG${org}_PEER${peer}_MSPCONFIGPATH} \\
CORE_PEER_TLS_ENABLED=\${CORE_PEER_TLS_ENABLED} \\
CORE_PEER_TLS_ROOTCERT_FILE=\${ORG${org}_PEER${peer}_TLS_ROOTCERT_FILE} \\
CORE_PEER_TLS_CERT_FILE=\${ORG${org}_PEER${peer}_TLS_CERT_FILE} \\
CORE_PEER_TLS_KEY_FILE=\${ORG${org}_PEER${peer}_TLS_KEY_FILE}\""
}

# 生成所有节点配置
for org in 1 2 3; do
    for peer in 0 1; do
        generate_peer_config $org $peer
        generate_cli_config $org $peer
    done
done

###########################################
# 主程序
###########################################

main() {
    # 记录开始时间
    local start_time=$(date +%s)

    # 显示脚本信息
    log_info "Hyperledger Fabric 网络部署脚本启动"

    # 执行前置检查
    show_progress 1 "检查环境依赖" $start_time
    execute_with_timer "检查前置条件" "check_prerequisites"
    execute_with_timer "检查Docker服务" "check_docker_service"

    # 确认执行
    echo -e "${RED}注意：倘若您之前已经部署过了 network ，执行该脚本会丢失之前的数据！${NC}"
    read -p "你确定要继续执行吗？请输入 Y 或 y 继续执行：" confirm

    if [[ "$confirm" != "Y" && "$confirm" != "y" ]]; then
        log_info "用户取消执行"
        exit 2
    fi

    # 清理环境
    show_progress 2 "清理环境" $start_time
    execute_with_timer "清理环境" "./uninstall.sh"
    mkdir config crypto-config data

    # 启动工具容器
    show_progress 3 "部署工具容器" $start_time
    execute_with_timer "部署工具容器" "docker-compose up -d ${CLI_CONTAINER}"
    log_success "工具容器部署完成"

    # 生成证书和密钥
    show_progress 4 "生成证书和密钥（MSP 材料）" $start_time
    execute_with_timer "生成证书和密钥" "$CLI_CMD \"cryptogen generate --config=${HYPERLEDGER_PATH}/crypto-config.yaml --output=${CRYPTO_PATH}\""

    # 创建排序通道创世区块
    show_progress 5 "创建排序通道创世区块" $start_time
    execute_with_timer "创建创世区块" "$CLI_CMD \"configtxgen -configPath ${HYPERLEDGER_PATH} -profile SampleGenesis -outputBlock ${CONFIG_PATH}/genesis.block -channelID firstchannel\""

    # 生成通道配置事务
    show_progress 6 "生成通道配置事务" $start_time
    execute_with_timer "生成通道配置" "$CLI_CMD \"configtxgen -configPath ${HYPERLEDGER_PATH} -profile SampleChannel -outputCreateChannelTx ${CONFIG_PATH}/$ChannelName.tx -channelID $ChannelName\""

    # 定义组织锚节点
    show_progress 7 "定义组织锚节点" $start_time
    execute_with_timer "定义Org1锚节点" "$CLI_CMD \"configtxgen -configPath ${HYPERLEDGER_PATH} -profile SampleChannel -outputAnchorPeersUpdate ${CONFIG_PATH}/Org1Anchor.tx -channelID $ChannelName -asOrg Org1\""
    execute_with_timer "定义Org2锚节点" "$CLI_CMD \"configtxgen -configPath ${HYPERLEDGER_PATH} -profile SampleChannel -outputAnchorPeersUpdate ${CONFIG_PATH}/Org2Anchor.tx -channelID $ChannelName -asOrg Org2\""
    execute_with_timer "定义Org3锚节点" "$CLI_CMD \"configtxgen -configPath ${HYPERLEDGER_PATH} -profile SampleChannel -outputAnchorPeersUpdate ${CONFIG_PATH}/Org3Anchor.tx -channelID $ChannelName -asOrg Org3\""

    # 启动所有节点
    show_progress 8 "启动所有节点" $start_time
    execute_with_timer "启动节点" "docker-compose up -d"
    wait_for_completion "等待节点启动（${NETWORK_STARTUP_WAIT}秒）" $NETWORK_STARTUP_WAIT

    # 创建通道
    show_progress 9 "创建通道" $start_time
    execute_with_timer "创建通道" "$CLI_CMD \"$Org1Peer0Cli peer channel create --outputBlock ${CONFIG_PATH}/$ChannelName.block -o $ORDERER1_ADDRESS -c $ChannelName -f ${CONFIG_PATH}/$ChannelName.tx --tls --cafile $ORDERER_CA\""

    # 节点加入通道
    show_progress 10 "节点加入通道" $start_time
    execute_with_timer "Org1Peer0加入通道" "$CLI_CMD \"$Org1Peer0Cli peer channel join -b ${CONFIG_PATH}/$ChannelName.block\""
    execute_with_timer "Org1Peer1加入通道" "$CLI_CMD \"$Org1Peer1Cli peer channel join -b ${CONFIG_PATH}/$ChannelName.block\""
    execute_with_timer "Org2Peer0加入通道" "$CLI_CMD \"$Org2Peer0Cli peer channel join -b ${CONFIG_PATH}/$ChannelName.block\""
    execute_with_timer "Org2Peer1加入通道" "$CLI_CMD \"$Org2Peer1Cli peer channel join -b ${CONFIG_PATH}/$ChannelName.block\""
    execute_with_timer "Org3Peer0加入通道" "$CLI_CMD \"$Org3Peer0Cli peer channel join -b ${CONFIG_PATH}/$ChannelName.block\""
    execute_with_timer "Org3Peer1加入通道" "$CLI_CMD \"$Org3Peer1Cli peer channel join -b ${CONFIG_PATH}/$ChannelName.block\""

    # 更新锚节点
    show_progress 11 "更新锚节点" $start_time
    execute_with_timer "更新Org1锚节点" "$CLI_CMD \"$Org1Peer0Cli peer channel update -o $ORDERER1_ADDRESS -c $ChannelName -f ${CONFIG_PATH}/Org1Anchor.tx --tls --cafile $ORDERER_CA\""
    execute_with_timer "更新Org2锚节点" "$CLI_CMD \"$Org2Peer0Cli peer channel update -o $ORDERER1_ADDRESS -c $ChannelName -f ${CONFIG_PATH}/Org2Anchor.tx --tls --cafile $ORDERER_CA\""
    execute_with_timer "更新Org3锚节点" "$CLI_CMD \"$Org3Peer0Cli peer channel update -o $ORDERER1_ADDRESS -c $ChannelName -f ${CONFIG_PATH}/Org3Anchor.tx --tls --cafile $ORDERER_CA\""

    # 打包链码
    show_progress 12 "打包链码" $start_time
    execute_with_timer "打包链码" "$CLI_CMD \"peer lifecycle chaincode package ${CHAINCODE_PACKAGE} --path ${CHAINCODE_PATH} --lang golang --label chaincode_${Version}\""

    # 安装链码
    show_progress 13 "安装链码" $start_time
    execute_with_timer "Org1Peer0安装链码" "$CLI_CMD \"$Org1Peer0Cli peer lifecycle chaincode install ${CHAINCODE_PACKAGE}\""
    execute_with_timer "Org1Peer1安装链码" "$CLI_CMD \"$Org1Peer1Cli peer lifecycle chaincode install ${CHAINCODE_PACKAGE}\""
    execute_with_timer "Org2Peer0安装链码" "$CLI_CMD \"$Org2Peer0Cli peer lifecycle chaincode install ${CHAINCODE_PACKAGE}\""
    execute_with_timer "Org2Peer1安装链码" "$CLI_CMD \"$Org2Peer1Cli peer lifecycle chaincode install ${CHAINCODE_PACKAGE}\""
    execute_with_timer "Org3Peer0安装链码" "$CLI_CMD \"$Org3Peer0Cli peer lifecycle chaincode install ${CHAINCODE_PACKAGE}\""
    execute_with_timer "Org3Peer1安装链码" "$CLI_CMD \"$Org3Peer1Cli peer lifecycle chaincode install ${CHAINCODE_PACKAGE}\""

    # 批准链码
    show_progress 14 "批准链码" $start_time
    PackageID=$($CLI_CMD "$Org1Peer0Cli peer lifecycle chaincode calculatepackageid ${CHAINCODE_PACKAGE}")
    execute_with_timer "Org1批准链码" "$CLI_CMD \"$Org1Peer0Cli peer lifecycle chaincode approveformyorg -o $ORDERER1_ADDRESS --channelID $ChannelName --name $ChainCodeName --version $Version --package-id $PackageID --sequence $Sequence --tls --cafile $ORDERER_CA\""
    execute_with_timer "Org2批准链码" "$CLI_CMD \"$Org2Peer0Cli peer lifecycle chaincode approveformyorg -o $ORDERER1_ADDRESS --channelID $ChannelName --name $ChainCodeName --version $Version --package-id $PackageID --sequence $Sequence --tls --cafile $ORDERER_CA\""
    execute_with_timer "Org3批准链码" "$CLI_CMD \"$Org3Peer0Cli peer lifecycle chaincode approveformyorg -o $ORDERER1_ADDRESS --channelID $ChannelName --name $ChainCodeName --version $Version --package-id $PackageID --sequence $Sequence --tls --cafile $ORDERER_CA\""

    # 提交链码
    show_progress 15 "提交链码" $start_time
    execute_with_timer "提交链码定义" "$CLI_CMD \"$Org1Peer0Cli peer lifecycle chaincode commit -o $ORDERER1_ADDRESS --channelID $ChannelName --name $ChainCodeName --version $Version --sequence $Sequence --tls --cafile $ORDERER_CA --peerAddresses $ORG1_PEER0_ADDRESS --tlsRootCertFiles $ORG1_PEER0_TLS_ROOTCERT_FILE --peerAddresses $ORG2_PEER0_ADDRESS --tlsRootCertFiles $ORG2_PEER0_TLS_ROOTCERT_FILE --peerAddresses $ORG3_PEER0_ADDRESS --tlsRootCertFiles $ORG3_PEER0_TLS_ROOTCERT_FILE\""

    # 初始化并验证
    show_progress 16 "初始化并验证" $start_time
    execute_with_timer "初始化链码" "$CLI_CMD \"$Org1Peer0Cli peer chaincode invoke -o $ORDERER1_ADDRESS -C $ChannelName -n $ChainCodeName -c '{\\\"function\\\":\\\"InitLedger\\\",\\\"Args\\\":[]}' --tls --cafile $ORDERER_CA --peerAddresses $ORG1_PEER0_ADDRESS --tlsRootCertFiles $ORG1_PEER0_TLS_ROOTCERT_FILE --peerAddresses $ORG2_PEER0_ADDRESS --tlsRootCertFiles $ORG2_PEER0_TLS_ROOTCERT_FILE --peerAddresses $ORG3_PEER0_ADDRESS --tlsRootCertFiles $ORG3_PEER0_TLS_ROOTCERT_FILE\""

    wait_for_completion "等待链码初始化（${CHAINCODE_INIT_WAIT}秒）" $CHAINCODE_INIT_WAIT

    if $CLI_CMD "$Org1Peer0Cli peer chaincode query -C $ChannelName -n $ChainCodeName -c '{\"Args\":[\"Hello\"]}'" 2>&1 | grep "hello"; then
        log_success "【恭喜您！】network 部署成功 (总耗时: $(time_elapsed $start_time))"
        exit 0
    fi

    log_error "【警告】network 未部署成功，请检查日志定位具体问题。(总耗时: $(time_elapsed $start_time))"
    exit 1
}

# 执行主程序
main "$@"
