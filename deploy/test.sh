#!/bin/bash

echo "验证查询账户信息"
docker exec cli peer chaincode invoke -C assetschannel -n blockchain-real-estate -c '{"Args":["queryAccountList"]}'