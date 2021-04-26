#!/bin/bash

echo "删除之前的"
docker rm -f blockchain-real-estate-application

echo "打包application镜像"
docker build -t togettoyou/blockchain-real-estate:1.0 .

echo "运行application"
docker run -it -d --name blockchain-real-estate-application -p 8000:8000 togettoyou/blockchain-real-estate:1.0