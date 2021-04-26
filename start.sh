#!/bin/bash

echo "删除旧容器"
docker rm -f blockchain-real-estate-application

echo "开始打包编译"
docker build -t togettoyou/blockchain-real-estate:1.0 .

echo "运行容器"
docker run -it -d --name blockchain-real-estate-application togettoyou/blockchain-real-estate:1.0

echo "拷贝容器中编译后的文件"
docker cp blockchain-real-estate-application:/root/blockchain-real-estate/application/app ./application/
docker cp blockchain-real-estate-application:/root/blockchain-real-estate/application/vue/dist ./application/dist

echo "运行"
./application/app