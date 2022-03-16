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

echo "三、生成创世区块文件，通道ID为 firstchannel"
configtxgen -profile TwoOrgsOrdererGenesis -outputBlock ./config/genesis.block -channelID firstchannel

echo "四、生成通道交易文件，通道ID为 appchannel"
configtxgen -profile TwoOrgChannel -outputCreateChannelTx ./config/appchannel.tx -channelID appchannel

echo "区块链 ： 启动"
docker-compose up -d
echo "正在等待节点的启动完成，等待10秒"
sleep 10

echo "五、创建通道"
docker exec cli peer channel create -o orderer.qq.com:7050 -c appchannel -f /etc/hyperledger/config/appchannel.tx

echo "六、节点加入通道"
docker exec cli peer channel join -b appchannel.block

# -n 链码名，可以自己随便设置
# -v 版本号
# -p 链码目录，在 /opt/gopath/src/ 目录下
echo "七、链码安装"
docker exec cli peer chaincode install -n fabric-realty -v 1.0.0 -l golang -p chaincode

# -n 对应上一步安装链码的名字
# -v 版本号
# -C 是通道，在fabric的世界，一个通道就是一条不同的链
# -c 为传参，传入init参数
echo "八、实例化链码"
docker exec cli peer chaincode instantiate -o orderer.qq.com:7050 -C appchannel -n fabric-realty -l golang -v 1.0.0 -c '{"Args":["init"]}'

echo "正在等待链码实例化完成，等待5秒"
sleep 5

# 进行链码交互，验证链码是否正确安装及区块链网络能否正常工作
echo "九、验证链码"
docker exec cli peer chaincode invoke -C appchannel -n fabric-realty -c '{"Args":["hello"]}'