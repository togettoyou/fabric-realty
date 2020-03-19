# blockchain-real-estate

> 🚀基于区块链的房地产交易系统小模型。提供销售和捐赠功能。本项目使用Hyperledger Fabric构建区块链网络, go编写智能合约，应用层使用gin+fabric-sdk-go调用合约。前端展示使用vue+element。前后端分离。
>

注：本项目需放在 `$GOPATH/src/github.com/togettoyou/blockchain-real-estate` 下运行

## [在线体验地址](http://blockchain.togettoyou.com/web) 

## 技术栈

- Hyperledger Fabric
- Docker
- Go Gin
- Vue
- ElementUI

## 运行

> 默认已经安装Hyperledger Fabric环境，如果未安装，参考：https://juejin.im/post/5e5db4ebf265da57301bfba5
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

## 演示效果图

![Mar-19-2020_15-28-20](https://github.com/togettoyou/blockchain-real-estate/blob/master/screenshots/Mar-19-2020_15-28-20.gif)

## 感谢

- [go-gin-example](https://github.com/eddycjy/go-gin-example)
- [vue-admin-template](https://github.com/PanJiaChen/vue-admin-template)


