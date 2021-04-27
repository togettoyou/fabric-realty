# 全新极速运行部署方式

---

> 🚀基于区块链的房地产交易系统小模型。提供销售和捐赠功能。本项目使用Hyperledger Fabric构建区块链网络, go编写智能合约，应用层使用gin+fabric-sdk-go调用合约。前端展示使用vue+element。前后端分离。

## 技术栈


- Hyperledger Fabric
- Go
- Vue
- Docker
## 前提
Linux 或者 Mac，连接网络，要求安装了 Docker 和 Docker Compose

附Linux安装 Docker 和 Docker Compose 教程：点此跳转

## 运行


1、克隆本项目放在任意目录下，例：`/root/blockchain-real-estate`

2、给予项目权限，执行 `sudo chmod -R +x /root/blockchain-real-estate/`

3、进入 `deploy` 目录，执行 `./start.sh` 启动区块链网络

4、进入 `vue` 目录，执行 `./build.sh` 编译前端

5、进入 `application` 目录，执行 `./build.sh` 编译后端

6、在 `application` 目录下，执行 `./start.sh` 启动应用

7、浏览器访问 [http://localhost:8000/web](http://localhost:8000/web)

## 目录结构


`application` : go gin + fabric-sdk-go 调用链码，提供外部访问接口，前端编译后静态资源放在`dist`目录下


`chaincode` : go 编写的智能合约


`deploy` : 区块链网络配置


`vue` : vue + element的前端展示页面


## 功能流程


管理员为用户业主创建房地产。


业主查看名下房产信息。


业主发起销售，所有人都可查看销售列表，购买者购买后进行扣款操作，并等待业主确认收款，交易完成，更新房产持有人。在有效期期间可以随时取消交易，有效期到期后自动关闭交易。


业主发起捐赠，指定受赠人，受赠人确认接收受赠前，双方可取消捐赠/受赠。


## 演示效果图




