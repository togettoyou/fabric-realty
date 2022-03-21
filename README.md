> 🚀 本项目使用 Hyperledger Fabric 构建底层区块链网络, go 编写智能合约，应用层使用 gin+fabric-sdk-go ，前端使用 vue+element-ui

如果想要联系我，可以关注我的公众号【SuperGopher】

![微信公众号.png](https://gitee.com/togettoyou/picture/raw/master/2022-2-9/1644374999459-weixin.jpg)

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

![image.png](https://cdn.nlark.com/yuque/0/2021/png/1077776/1619503587830-48d3d53d-92eb-4848-8a38-da2d07b5b119.png#align=left&display=inline&height=777&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1554&originWidth=2875&size=232911&status=done&style=none&width=1437.5#id=nUKaE&originHeight=1554&originWidth=2875&originalType=binary&ratio=1&status=done&style=none)
![image.png](https://cdn.nlark.com/yuque/0/2021/png/1077776/1619503608573-35bcf8ad-5738-4df8-bd7b-4824650c0e13.png#align=left&display=inline&height=778&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1555&originWidth=2880&size=255025&status=done&style=none&width=1440#id=aVYox&originHeight=1555&originWidth=2880&originalType=binary&ratio=1&status=done&style=none)
![image.png](https://cdn.nlark.com/yuque/0/2021/png/1077776/1619503660695-3596146f-a09c-4914-8667-f2f468e768a5.png#align=left&display=inline&height=779&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1558&originWidth=2880&size=232348&status=done&style=none&width=1440#id=tu55k&originHeight=1558&originWidth=2880&originalType=binary&ratio=1&status=done&style=none)
