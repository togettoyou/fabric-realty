#!/bin/bash

priv_sk_path=$(ls ../crypto-config/peerOrganizations/taobao.com/users/Admin\@taobao.com/msp/keystore/)

cp -rf ./connection-profile/network_temp.json ./connection-profile/network.json

if [[ `uname` == 'Darwin' ]]; then
  echo "当前操作系统是 Mac"
  sed -i '' "s/priv_sk/$priv_sk_path/" ./connection-profile/network.json
elif [[ `uname` == 'Linux' ]]; then
  echo "当前操作系统是 Linux"
  sed -i "s/priv_sk/$priv_sk_path/" ./connection-profile/network.json
else
  echo "当前操作系统不是 Mac 或 Linux，脚本无法继续执行！"
  exit 1
fi

docker-compose down -v
docker-compose up -d
