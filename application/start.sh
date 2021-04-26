#!/bin/bash

echo "一、清理环境、删除旧容器"
rm -rf app
docker rm -f blockchain-real-estate-application


echo "二、开始打包编译"
if [[ `uname` == 'Darwin' ]]; then
  docker build -f  Dockerfile.Mac -t togettoyou/blockchain-real-estate:v1 .
fi
if [[ `uname` == 'Linux' ]]; then
  docker build -f  Dockerfile.Linux -t togettoyou/blockchain-real-estate:v1 .
fi

echo "三、运行编译容器"
docker run -it -d --name blockchain-real-estate-application togettoyou/blockchain-real-estate:v1

echo "四、拷贝容器中编译后的二进制可执行文件"
docker cp blockchain-real-estate-application:/root/application/app ./application/
sudo chmod +x ./app
./app