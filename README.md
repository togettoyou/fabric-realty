> 🚀 本项目使用 Hyperledger Fabric 构建底层区块链网络, go 编写智能合约，应用层使用 gin+fabric-sdk-go ，前端使用 vue+element-ui

如果想要联系我，可以关注我的公众号【SuperGopher】

![微信公众号.png](https://user-images.githubusercontent.com/55381228/155444889-eacc0104-cd85-45c9-b7b7-9036e0c2334c.jpg)

## 教程

[万字长文，教你用go开发区块链应用](https://mp.weixin.qq.com/s/yDmGwfRjXxDJfgv1d0p3Ig)

## 环境要求

安装了 Docker 和 Docker Compose 的 Linux 环境

附 Docker 安装教程：[点此跳转](Install.md)

## 部署

1. 克隆本项目放在任意目录下，例：`/root/fabric-realty`


2. 给予项目权限，执行 `sudo chmod -R +x /root/fabric-realty/`


3. 进入 `network` 目录，执行 `./start.sh` 启动区块链网络以及部署智能合约


4. 进入 `application` 目录，执行 `./build.sh` 编译镜像，完成后继续执行 `./start.sh`
   启动应用，最后可使用浏览器访问 [http://localhost:8000/web](http://localhost:8000/web)


5. （可选）进入 `network/explorer` 目录，执行 `./start.sh` 启动区块链浏览器后，访问 [http://localhost:8080](http://localhost:8080)，用户名 admin，密码
   123456

## 停止或重启

注意，默认执行 `./start.sh` 脚本时都会调用 `./stop.sh` ，所以如果想持久化数据的情况下停止或重启本项目，请不要重新执行 `./start.sh` ，正确姿势参考：

1. （如果启动了区块链浏览器）进入 `network/explorer` 目录，执行 `docker-compose stop` 停止区块链浏览器，执行 `docker-compose start`
   启动区块链浏览器，执行 `docker-compose restart` 重启区块链浏览器

2. 进入 `network` 目录，执行 `docker-compose stop` 停止区块链网络，执行 `docker-compose start`
   启动区块链网络，执行 `docker-compose restart` 重启区块链网络

3. 进入 `application` 目录，区块链应用是无状态的，可以直接执行 `./stop.sh` 关闭区块链应用，执行 `./start.sh` 关闭区块链应用

## 完全清理环境

注意，该操作会将所有数据清空。按照该顺序：

1. （如果启动了区块链浏览器）进入 `network/explorer` 目录，执行 `./stop.sh` 关闭区块链浏览器

2. 进入 `network` 目录，执行 `./stop.sh` 关闭区块链网络并清理链码容器

3. 进入 `application` 目录，执行 `./stop.sh` 关闭区块链应用

## 目录结构

- `application/server` : `fabric-sdk-go` 调用链码（即智能合约），`gin` 提供外部访问接口（RESTful API）


- `application/web` : `vue` + `element-ui` 提供前端展示页面


- `chaincode` : go 编写的链码（即智能合约）


- `network` : Hyperledger Fabric 区块链网络配置

## 功能流程

管理员为用户业主创建房地产。

业主查看名下房产信息。

业主发起销售，所有人都可查看销售列表，购买者购买后进行扣款操作，并等待业主确认收款，交易完成后，更新房产持有人。在有效期期间可以随时取消交易，有效期到期后自动关闭交易。

业主发起捐赠，指定受赠人，受赠人确认接收受赠前，双方可取消捐赠/受赠。

## 演示效果

![login](https://user-images.githubusercontent.com/55381228/159389012-4d3d8617-2bd8-4d9c-bacf-452f97cc9bbc.png)

![addreal](https://user-images.githubusercontent.com/55381228/159389026-9ca119bd-fd5f-4b89-b003-a09907ce0cdf.png)

![info](https://user-images.githubusercontent.com/55381228/159389035-b84f2de1-18f9-48a7-93f5-db9dd20a5a4c.png)

![explorer](https://user-images.githubusercontent.com/55381228/159389002-0dbe329a-09aa-4aaf-aba8-4a98e4fdcc39.png)
