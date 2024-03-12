> 🚀 本项目使用 Hyperledger Fabric 构建底层区块链网络, go 编写智能合约，应用层使用 gin+fabric-sdk-go ，前端使用
> vue+element-ui

如果想要联系我，可以关注我的公众号【gopher云原生】

![gopher云原生](https://user-images.githubusercontent.com/55381228/221747734-13783ce6-1969-4c10-acd6-833f5046aa85.png)

## 教程

[万字长文，教你用go开发区块链应用](https://mp.weixin.qq.com/s/yDmGwfRjXxDJfgv1d0p3Ig)

> 🤔 有任何疑问，请先看完本篇文章。本项目涉及的知识点都有在文章中进行说明

## 手动部署

环境要求： 安装了 Docker 和 Docker Compose 的 Linux 或 Mac OS 环境

附 Linux Docker 安装教程：[点此跳转](Install.md)

> 🤔 Docker 和 Docker Compose 需要先自行学习。本项目的区块链网络搭建、链码部署、前后端编译/部署都是使用 Docker 和 Docker
> Compose 完成的。

1. 下载本项目放在任意目录下，例：`/root/fabric-realty`

2. 给予项目权限，执行 `sudo chmod -R +x /root/fabric-realty/`

3. 进入 `network` 目录，执行 `./start.sh` 部署区块链网络和智能合约

4. 进入 `application` 目录，执行 `./start.sh`
   启动前后端应用，然后就可使用浏览器访问前端页面 [http://localhost:8000](http://localhost:8000)
   ，其中后端接口地址为 [http://localhost:8888](http://localhost:8888)

5. （可选）进入 `network/explorer` 目录，执行 `./start.sh`
   启动区块链浏览器后，访问 [http://localhost:8080](http://localhost:8080)，用户名 admin，密码
   123456

## 完全清理环境

注意，该操作会将所有数据清空。按照该先后顺序：

1. （如果启动了区块链浏览器）进入 `network/explorer` 目录，执行 `./stop.sh` 关闭区块链浏览器

2. 进入 `application` 目录，执行 `./stop.sh` 关闭区块链应用

3. 最后进入 `network` 目录，执行 `./stop.sh` 关闭区块链网络并清理链码容器

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

