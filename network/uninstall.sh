#!/bin/bash

ChainCodeName="chaincodetogettoyou"

# 清除链码容器
function clearContainers() {
  CONTAINER_IDS=$(docker ps -a | awk '($2 ~ /dev-peer.*'$ChainCodeName'.*/) {print $1}')
  if [ -z "$CONTAINER_IDS" -o "$CONTAINER_IDS" == " " ]; then
    echo "---- No containers available for deletion ----"
  else
    docker rm -f $CONTAINER_IDS
  fi
}

# 清除链码镜像
function removeUnwantedImages() {
  DOCKER_IMAGE_IDS=$(docker images | awk '($1 ~ /dev-peer.*'$ChainCodeName'.*/) {print $3}')
  if [ -z "$DOCKER_IMAGE_IDS" -o "$DOCKER_IMAGE_IDS" == " " ]; then
    echo "---- No images available for deletion ----"
  else
    docker rmi -f $DOCKER_IMAGE_IDS
  fi
}

read -p "你确定要卸载吗？请输入 Y 或 y 继续执行：" confirm

if [[ "$confirm" != "Y" && "$confirm" != "y" ]]; then
  echo "你取消了脚本的执行。"
  exit 1
fi

echo "清理环境"
docker-compose -f explorer/docker-compose.yaml down -v
docker-compose down -v
mkdir -p config
mkdir -p crypto-config
mkdir -p data
rm -rf config/*
rm -rf crypto-config/*
rm -rf data/*
clearContainers
removeUnwantedImages
echo "清理完毕"
