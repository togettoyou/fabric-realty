#!/bin/bash

echo "一、删除旧容器"
docker rm -f blockchain-real-estate-application

echo "二、开始打包编译"
docker build -t togettoyou/blockchain-real-estate:v1 .

echo "三、运行编译容器"
docker run -it -d --name blockchain-real-estate-application togettoyou/blockchain-real-estate:v1

echo "四、拷贝容器中编译后的文件"
docker cp blockchain-real-estate-application:/root/application/app ./application/
docker cp blockchain-real-estate-application:/root/application/dist ./application/dist

if [[ `uname` == 'Darwin' ]]; then
    echo "Mac OS 请自行在 application 下 go run main.go 或 go build"
fi
if [[ `uname` == 'Linux' ]]; then
    echo "五、运行"
    ./application/app
fi