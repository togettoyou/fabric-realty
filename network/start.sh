#!/bin/bash

# 检查操作系统类型
if [[ `uname` == 'Darwin' ]]; then
  echo "当前操作系统是 Mac"
  export PATH=${PWD}/hyperledger-fabric-darwin-amd64-1.4.12/bin:$PATH
elif [[ `uname` == 'Linux' ]]; then
  echo "当前操作系统是 Linux"
  export PATH=${PWD}/hyperledger-fabric-linux-amd64-1.4.12/bin:$PATH
else
  echo "当前操作系统不是 Mac 或 Linux，脚本无法继续执行！"
  exit 1
fi

echo -e "注意：倘若您之前已经部署过了 network ，执行该脚本会丢失之前的数据！\n若只是想重启 network ，可以直接执行 docker-compose restart 命令。"
read -p "你确定要继续执行吗？请输入 Y 或 y 继续执行：" confirm

if [[ "$confirm" != "Y" && "$confirm" != "y" ]]; then
  echo "你取消了脚本的执行。"
  exit 1
fi

echo "一、清理环境"
./stop.sh

echo "二、生成证书和秘钥（ MSP 材料），生成结果将保存在 crypto-config 文件夹中"
cryptogen generate --config=./crypto-config.yaml

echo "三、创建排序通道创世区块"
configtxgen -profile TwoOrgsOrdererGenesis -outputBlock ./config/genesis.block -channelID firstchannel

echo "四、生成通道配置事务'appchannel.tx'"
configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ./config/appchannel.tx -channelID appchannel

echo "五、为 Taobao 定义锚节点"
configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./config/TaobaoAnchor.tx -channelID appchannel -asOrg Taobao

echo "六、为 JD 定义锚节点"
configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./config/JDAnchor.tx -channelID appchannel -asOrg JD

echo "区块链 ： 启动"
docker-compose up -d
echo "正在等待节点的启动完成，等待10秒"
sleep 10

TaobaoPeer0Cli="CORE_PEER_ADDRESS=peer0.taobao.com:7051 CORE_PEER_LOCALMSPID=TaobaoMSP CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/peer/taobao.com/users/Admin@taobao.com/msp"
TaobaoPeer1Cli="CORE_PEER_ADDRESS=peer1.taobao.com:7051 CORE_PEER_LOCALMSPID=TaobaoMSP CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/peer/taobao.com/users/Admin@taobao.com/msp"
JDPeer0Cli="CORE_PEER_ADDRESS=peer0.jd.com:7051 CORE_PEER_LOCALMSPID=JDMSP CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/peer/jd.com/users/Admin@jd.com/msp"
JDPeer1Cli="CORE_PEER_ADDRESS=peer1.jd.com:7051 CORE_PEER_LOCALMSPID=JDMSP CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/peer/jd.com/users/Admin@jd.com/msp"

echo "七、创建通道"
docker exec cli bash -c "$TaobaoPeer0Cli peer channel create -o orderer.qq.com:7050 -c appchannel -f /etc/hyperledger/config/appchannel.tx"

echo "八、将所有节点加入通道"
docker exec cli bash -c "$TaobaoPeer0Cli peer channel join -b appchannel.block"
docker exec cli bash -c "$TaobaoPeer1Cli peer channel join -b appchannel.block"
docker exec cli bash -c "$JDPeer0Cli peer channel join -b appchannel.block"
docker exec cli bash -c "$JDPeer1Cli peer channel join -b appchannel.block"

echo "九、更新锚节点"
docker exec cli bash -c "$TaobaoPeer0Cli peer channel update -o orderer.qq.com:7050 -c appchannel -f /etc/hyperledger/config/TaobaoAnchor.tx"
docker exec cli bash -c "$JDPeer0Cli peer channel update -o orderer.qq.com:7050 -c appchannel -f /etc/hyperledger/config/JDAnchor.tx"

# -n 链码名，可以自己随便设置
# -v 版本号
# -p 链码目录，在 /opt/gopath/src/ 目录下
echo "十、安装链码"
docker exec cli bash -c "$TaobaoPeer0Cli peer chaincode install -n fabric-realty -v 1.0.0 -l golang -p chaincode"
docker exec cli bash -c "$JDPeer0Cli peer chaincode install -n fabric-realty -v 1.0.0 -l golang -p chaincode"

# 只需要其中一个节点实例化
# -n 对应上一步安装链码的名字
# -v 版本号
# -C 是通道，在fabric的世界，一个通道就是一条不同的链
# -c 为传参，传入init参数
echo "十一、实例化链码"
docker exec cli bash -c "$TaobaoPeer0Cli peer chaincode instantiate -o orderer.qq.com:7050 -C appchannel -n fabric-realty -l golang -v 1.0.0 -c '{\"Args\":[\"init\"]}' -P \"AND ('TaobaoMSP.member','JDMSP.member')\""

echo "正在等待链码实例化完成，等待5秒"
sleep 5

# 进行链码交互，验证链码是否正确安装及区块链网络能否正常工作
echo "十二、验证链码"
docker exec cli bash -c "$TaobaoPeer0Cli peer chaincode invoke -C appchannel -n fabric-realty -c '{\"Args\":[\"hello\"]}'"

if docker exec cli bash -c "$JDPeer0Cli peer chaincode invoke -C appchannel -n fabric-realty -c '{\"Args\":[\"hello\"]}'" 2>&1 | grep "Chaincode invoke successful"; then
  echo "【恭喜您！】 network 部署成功，后续如需暂时停止运行，可以执行 docker-compose stop 命令（不会丢失数据）。"
  exit 0
fi

echo "【警告】network 未部署成功，请检查每一个步骤，定位具体问题。"