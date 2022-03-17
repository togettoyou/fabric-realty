#!/bin/bash

# 根据需求保留，这里相当于使用fabric-samples_v1.4.12中的bin
if [[ `uname` == 'Darwin' ]]; then
    echo "Mac OS"
fi
if [[ `uname` == 'Linux' ]]; then
    echo "Linux"
    export PATH=${PWD}/hyperledger-fabric-linux-amd64-1.4.12/bin:$PATH
fi

echo "一、清理环境"
mkdir -p config
mkdir -p crypto-config
rm -rf config/*
rm -rf crypto-config/*
docker-compose down -v
echo "清理完毕"

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

TaobaoCli="CORE_PEER_ADDRESS=peer0.taobao.com:7051 CORE_PEER_LOCALMSPID=TaobaoMSP CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/peer/taobao.com/users/Admin@taobao.com/msp"
JDCli="CORE_PEER_ADDRESS=peer0.jd.com:7051 CORE_PEER_LOCALMSPID=JDMSP CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/peer/jd.com/users/Admin@jd.com/msp"

echo "七、创建通道"
docker exec cli bash -c "$TaobaoCli peer channel create -o orderer.qq.com:7050 -c appchannel -f /etc/hyperledger/config/appchannel.tx"

echo "八、节点加入通道"
docker exec cli bash -c "$TaobaoCli peer channel join -b appchannel.block"
docker exec cli bash -c "$JDCli peer channel join -b appchannel.block"

echo "九、更新锚节点"
docker exec cli bash -c "$TaobaoCli peer channel update -o orderer.qq.com:7050 -c appchannel -f /etc/hyperledger/config/TaobaoAnchor.tx"
docker exec cli bash -c "$JDCli peer channel update -o orderer.qq.com:7050 -c appchannel -f /etc/hyperledger/config/JDAnchor.tx"

# -n 链码名，可以自己随便设置
# -v 版本号
# -p 链码目录，在 /opt/gopath/src/ 目录下
echo "十、安装链码"
docker exec cli bash -c "$TaobaoCli peer chaincode install -n fabric-realty -v 1.0.0 -l golang -p chaincode"
docker exec cli bash -c "$JDCli peer chaincode install -n fabric-realty -v 1.0.0 -l golang -p chaincode"

# -n 对应上一步安装链码的名字
# -v 版本号
# -C 是通道，在fabric的世界，一个通道就是一条不同的链
# -c 为传参，传入init参数
echo "十一、实例化链码"
docker exec cli bash -c "$TaobaoCli peer chaincode instantiate -o orderer.qq.com:7050 -C appchannel -n fabric-realty -l golang -v 1.0.0 -c '{\"Args\":[\"init\"]}' -P \"OR ('TaobaoMSP.member','JDMSP.member')\""

echo "正在等待链码实例化完成，等待5秒"
sleep 5

# 进行链码交互，验证链码是否正确安装及区块链网络能否正常工作
echo "十二、验证链码"
docker exec cli bash -c "$TaobaoCli peer chaincode invoke -C appchannel -n fabric-realty -c '{\"Args\":[\"hello\"]}'"
docker exec cli bash -c "$JDCli peer chaincode invoke -C appchannel -n fabric-realty -c '{\"Args\":[\"hello\"]}'"