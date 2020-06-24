# blockchain-real-estate

2020.6.24 更新详细运行步骤

1. 确保你的项目目录为`$GOPATH/src/github.com/togettoyou/blockchain-real-estate`
2. 项目由于未使用mod管理，请先将go mod环境设置为auto： `go env -w GO111MODULE=auto`
3. 首先测试chaincode是否正常调用，运行`chaincode/blockchain-real-estate/chaincode_test.go`测试用例
![image](https://user-images.githubusercontent.com/55381228/85498013-8200a100-b611-11ea-938f-9ac1d3ad5b89.png)

4. 在deploy目录下运行`./start.sh`,观察有无报错提示。运行成功后在终端执行`docker exec cli peer chaincode invoke -C assetschannel -n blockchain-real-estate -c '{"Args":["queryAccountList"]}'` 等cli命令，Args可以替换为Invoke中的funcName，先验证链码是否正确安装及区块链网络能否正常工作。建议`./start.sh`之前可以先运行`./stop.sh`清理一下环境。
![image](https://user-images.githubusercontent.com/55381228/85497727-0141a500-b611-11ea-8d10-deacb8bd627e.png)

5. 如果以上都成功，说明区块链网络是没有问题的。接下来同样先执行`application/sdk_test.go`单元测试，看是否可以成功使用sdk调用链码(此步骤前提你区块链网络即以上步骤已成功启动)
![image](https://user-images.githubusercontent.com/55381228/85497628-d7887e00-b610-11ea-9749-0006ad0df814.png)

6. 运行application，`go run main.go` 

我的本机测试环境：
![image](https://user-images.githubusercontent.com/55381228/85497883-4960c780-b611-11ea-93b0-4a2ec69b8142.png)

***
分割线
***

> 🚀基于区块链的房地产交易系统小模型。提供销售和捐赠功能。本项目使用Hyperledger Fabric构建区块链网络, go编写智能合约，应用层使用gin+fabric-sdk-go调用合约。前端展示使用vue+element。前后端分离。

注：本项目需放在 `$GOPATH/src/github.com/togettoyou/blockchain-real-estate` 下运行

## [在线体验地址](http://blockchain.togettoyou.com/web) 

## 技术栈

- Hyperledger Fabric
- Docker
- Go Gin
- Vue
- ElementUI

## 运行

> 默认已经安装Hyperledger Fabric环境，如果未安装，参考：https://www.yuque.com/togettoyou/blog/his57f
>
> 我的本机环境参考：
>
> ![Snipaste_2020-03-19_14-52-13](https://github.com/togettoyou/blockchain-real-estate/blob/master/screenshots/Snipaste_2020-03-19_14-52-13.png)



1、克隆本项目放在 `$GOPATH/src/github.com/togettoyou/blockchain-real-estate` 下

2、进入deploy目录，执行`start.sh`脚本

```shell
# 赋予权限
sudo chmod +x *.sh
# 启动区块链网络
./start.sh
# 停止区块链网络
./stop.sh
# 如果启动失败，可能是环境清理不干净，可以尝试先./stop.sh清理环境再./start.sh启动
```

3、进入application目录，启动应用程序

```shell
# 编译
go build
# 赋予权限
sudo chmod +x application
# 启动
./application
```

4、浏览器访问 http://localhost:8000/web

## 目录结构

`application` : go gin + fabric-sdk-go 调用链码，提供外部访问接口，前端静态资源放在`dist`目录下

`chaincode` : go 编写的智能合约

`deploy` : 区块链网络的配置以及启动停止脚本

`vendor` : 项目所需依赖包，防止网络原因下载失败

`vue` : vue + element的前端展示页面

```shell
# 如果需要修改前端页面，在vue目录下执行
yarn install
# 启动
yarn dev
# 重新打包生成dist资源，将dist放到application目录下覆盖
yarn build:prod
```

`screenshots` : 截图

## 功能流程

管理员为用户业主创建房地产。

业主查看名下房产信息。

业主发起销售，所有人都可查看销售列表，购买者购买后进行扣款操作，并等待业主确认收款，交易完成，更新房产持有人。在有效期期间可以随时取消交易，有效期到期后自动关闭交易。

业主发起捐赠，指定受赠人，受赠人确认接收受赠前，双方可取消捐赠/受赠。

## 演示效果图

![Mar-19-2020_15-28-20](https://github.com/togettoyou/blockchain-real-estate/blob/master/screenshots/Mar-19-2020_15-28-20.gif)

## 感谢

- [go-gin-example](https://github.com/eddycjy/go-gin-example)
- [vue-admin-template](https://github.com/PanJiaChen/vue-admin-template)

