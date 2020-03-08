#!/bin/bash

# 清除链码容器 正则：/dev-peer.*.blockchain-real-estate.*/ 匹配上的都会被删除，其中blockchain-real-estate是链码名称，在安装和实例化的时候会指定
function clearContainers() {
  CONTAINER_IDS=$(docker ps -a | awk '($2 ~ /dev-peer.*.blockchain-real-estate.*/) {print $1}')
  if [ -z "$CONTAINER_IDS" -o "$CONTAINER_IDS" == " " ]; then
    echo "---- No containers available for deletion ----"
  else
    docker rm -f $CONTAINER_IDS
  fi
}

# 清除不需要的链码镜像 正则：/dev-peer.*.blockchain-real-estate.*/ 匹配上的都会被删除，其中blockchain-real-estate是链码名称，在安装和实例化的时候会指定
function removeUnwantedImages() {
  DOCKER_IMAGE_IDS=$(docker images | awk '($1 ~ /dev-peer.*.blockchain-real-estate.*/) {print $3}')
  if [ -z "$DOCKER_IMAGE_IDS" -o "$DOCKER_IMAGE_IDS" == " " ]; then
    echo "---- No images available for deletion ----"
  else
    docker rmi -f $DOCKER_IMAGE_IDS
  fi
}

echo "区块链 ： 关闭"

echo "开始删除链码生成的docker镜像"
docker-compose down --volumes --remove-orphans

# 调用函数清除链码容器
clearContainers

# 调用函数清除不需要的链码镜像
removeUnwantedImages




