#!/bin/bash

echo "一、清理环境、删除旧容器"
rm -rf application/app
rm -rf application/dist
docker rm -f blockchain-real-estate-application

echo "二、开始打包编译"
docker build -t togettoyou/blockchain-real-estate:v1 .

echo "三、运行编译容器"
docker run -it -d --name blockchain-real-estate-application togettoyou/blockchain-real-estate:v1

if [[ `uname` == 'Darwin' ]]; then
    echo "四、拷贝容器中编译后的静态资源"
    docker cp blockchain-real-estate-application:/root/application/dist ./application/dist
    echo "静态资源编译已完成，Mac OS 请自行在 application 下 go run main.go 或 go build 运行"
fi
if [[ `uname` == 'Linux' ]]; then
    echo "四、拷贝容器中编译后的文件并运行"
    docker cp blockchain-real-estate-application:/root/application/app ./application/
    docker cp blockchain-real-estate-application:/root/application/dist ./application/dist
    ./application/app
fi