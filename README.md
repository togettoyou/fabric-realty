> 🚀基于区块链的房地产交易系统模型。提供销售和捐赠功能。本项目使用Hyperledger Fabric构建区块链网络, go编写智能合约，应用层使用gin+fabric-sdk-go调用合约。前端展示使用vue+element。前后端分离。

如果想要联系我，或者该项目确实帮助到了您，可以关注一下我的公众号【寻寻觅觅的Gopher】

![微信公众号.png](https://cdn.nlark.com/yuque/0/2021/png/1077776/1628483947581-9a649b2f-a0bb-4ef4-879d-92ab6e9fddde.png)

## 技术栈

- Hyperledger Fabric
- Go
- Vue

## 部署环境

- Docker
- Docker Compose

## 前提

Linux 或 Mac，要求安装了 Docker 和 Docker Compose

附 Linux 安装 Docker 和 Docker Compose 教程：[点此跳转](/Install.md)

我的测试环境：

![image.png](https://cdn.nlark.com/yuque/0/2021/png/1077776/1619497705974-f2cf0c33-5718-4b45-8bd8-aed870b86aa8.png#align=left&display=inline&height=160&margin=%5Bobject%20Object%5D&name=image.png&originHeight=319&originWidth=1116&size=40973&status=done&style=none&width=558#id=QpYhH&originHeight=319&originWidth=1116&originalType=binary&ratio=1&status=done&style=none)

## 运行

### 1、克隆本项目放在任意目录下，例：`/root/blockchain-real-estate`

![image.png](https://cdn.nlark.com/yuque/0/2021/png/1077776/1619497926959-136131db-40b9-4d9d-8949-9a24015e6b29.png#align=left&display=inline&height=139&margin=%5Bobject%20Object%5D&name=image.png&originHeight=278&originWidth=1345&size=29585&status=done&style=none&width=672.5#id=gMfwQ&originHeight=278&originWidth=1345&originalType=binary&ratio=1&status=done&style=none)

### 2、给予项目权限，执行 `sudo chmod -R +x /root/blockchain-real-estate/`

![image.png](https://cdn.nlark.com/yuque/0/2021/png/1077776/1619497967789-8728ba28-6026-4aca-bf6e-5847c9e6dee8.png#align=left&display=inline&height=128&margin=%5Bobject%20Object%5D&name=image.png&originHeight=255&originWidth=1422&size=32430&status=done&style=none&width=711#id=Oos1G&originHeight=255&originWidth=1422&originalType=binary&ratio=1&status=done&style=none)

### 3、进入 `deploy` 目录，执行 `./start.sh` 启动区块链网络

![image.png](https://cdn.nlark.com/yuque/0/2021/png/1077776/1619498040768-995d25af-fcd5-41e4-92b9-b0b1f5263c0e.png#align=left&display=inline&height=145&margin=%5Bobject%20Object%5D&name=image.png&originHeight=289&originWidth=1128&size=24879&status=done&style=none&width=564#id=RLedU&originHeight=289&originWidth=1128&originalType=binary&ratio=1&status=done&style=none)
![image.png](https://cdn.nlark.com/yuque/0/2021/png/1077776/1619503231479-0628da82-bb59-4cc2-8d6e-ec1b07b8d030.png#align=left&display=inline&height=698&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1396&originWidth=2361&size=417175&status=done&style=none&width=1180.5#id=nW5qo&originHeight=1396&originWidth=2361&originalType=binary&ratio=1&status=done&style=none)

### 4、进入 `vue` 目录，执行 `./build.sh` 编译前端

![image.png](https://cdn.nlark.com/yuque/0/2021/png/1077776/1619498139589-19c53edf-202c-429f-8cdd-381ef8083e66.png#align=left&display=inline&height=159&margin=%5Bobject%20Object%5D&name=image.png&originHeight=318&originWidth=1201&size=25754&status=done&style=none&width=600.5#id=BCV2I&originHeight=318&originWidth=1201&originalType=binary&ratio=1&status=done&style=none)
![image.png](https://cdn.nlark.com/yuque/0/2021/png/1077776/1619501158280-3028b279-eb14-47fc-9880-f5584df005c9.png#align=left&display=inline&height=500&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1000&originWidth=2361&size=167745&status=done&style=none&width=1180.5#id=n1sxZ&originHeight=1000&originWidth=2361&originalType=binary&ratio=1&status=done&style=none)

### 5、进入 `application` 目录，执行 `./build.sh` 编译后端

![image.png](https://cdn.nlark.com/yuque/0/2021/png/1077776/1619498187100-a82374b4-e985-439f-91d7-a3e9d3924dc4.png#align=left&display=inline&height=173&margin=%5Bobject%20Object%5D&name=image.png&originHeight=345&originWidth=1265&size=28209&status=done&style=none&width=632.5#id=Wy8vT&originHeight=345&originWidth=1265&originalType=binary&ratio=1&status=done&style=none)
![image.png](https://cdn.nlark.com/yuque/0/2021/png/1077776/1619503373258-82447169-cc83-4efe-ac32-98513b67bb29.png#align=left&display=inline&height=611&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1222&originWidth=1300&size=166511&status=done&style=none&width=650#id=WvGZ9&originHeight=1222&originWidth=1300&originalType=binary&ratio=1&status=done&style=none)

### 6、在 `application` 目录下，执行 `./start.sh` 启动应用

![image.png](https://cdn.nlark.com/yuque/0/2021/png/1077776/1619501464096-a543fd23-153e-4ddc-bd56-472698966940.png#align=left&display=inline&height=159&margin=%5Bobject%20Object%5D&name=image.png&originHeight=317&originWidth=1952&size=54818&status=done&style=none&width=976#id=vzis9&originHeight=317&originWidth=1952&originalType=binary&ratio=1&status=done&style=none)
![image.png](https://cdn.nlark.com/yuque/0/2021/png/1077776/1619501482450-7dc34559-6c39-4f8e-a7fe-177659517304.png#align=left&display=inline&height=698&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1396&originWidth=2361&size=367532&status=done&style=none&width=1180.5#id=vBRNT&originHeight=1396&originWidth=2361&originalType=binary&ratio=1&status=done&style=none)

### 7、浏览器访问 [http://localhost:8000/web](http://localhost:8000/web)

![](https://cdn.nlark.com/yuque/0/2021/png/1077776/1619503481607-d6dd7048-77aa-4461-817c-2fcf7507cf9d.png#id=gsVRB&originHeight=1568&originWidth=2874&originalType=binary&ratio=1&status=done&style=none)

### 8、（可选）进入 `deploy/explorer` 目录，执行 `./start.sh` 启动区块链浏览器

![image.png](https://cdn.nlark.com/yuque/0/2021/png/1077776/1623386161368-d06f0e91-a2be-43bf-83bc-d6921bc0dc3f.png#clientId=u7065799c-2510-4&from=paste&height=698&id=u5217fa8e&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1396&originWidth=2353&originalType=binary&ratio=2&size=177974&status=done&style=none&taskId=u842d45fc-0803-45be-ab6a-fc450905600&width=1176.5)

### 浏览器访问 [http://localhost:8080](http://localhost:8080)，用户名 admin，密码 123456

![image.png](https://cdn.nlark.com/yuque/0/2021/png/1077776/1623386227586-bc0f4deb-cf1e-4fae-9186-3c420ef7fd32.png#clientId=u7065799c-2510-4&from=paste&height=789&id=u50d0a26d&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1577&originWidth=2880&originalType=binary&ratio=2&size=133361&status=done&style=none&taskId=u759e0e20-65c1-43da-8cf5-a26b86b3643&width=1440)
![image.png](https://cdn.nlark.com/yuque/0/2021/png/1077776/1623386244686-58036523-b4d1-4054-9090-bf0156a53223.png#clientId=u7065799c-2510-4&from=paste&height=789&id=u3b2d7535&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1578&originWidth=2880&originalType=binary&ratio=2&size=300168&status=done&style=none&taskId=ua2921a32-db44-4b9f-bcbd-5cd5cd36a70&width=1440)

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

![image.png](https://cdn.nlark.com/yuque/0/2021/png/1077776/1619503587830-48d3d53d-92eb-4848-8a38-da2d07b5b119.png#align=left&display=inline&height=777&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1554&originWidth=2875&size=232911&status=done&style=none&width=1437.5#id=nUKaE&originHeight=1554&originWidth=2875&originalType=binary&ratio=1&status=done&style=none)
![image.png](https://cdn.nlark.com/yuque/0/2021/png/1077776/1619503608573-35bcf8ad-5738-4df8-bd7b-4824650c0e13.png#align=left&display=inline&height=778&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1555&originWidth=2880&size=255025&status=done&style=none&width=1440#id=aVYox&originHeight=1555&originWidth=2880&originalType=binary&ratio=1&status=done&style=none)
![image.png](https://cdn.nlark.com/yuque/0/2021/png/1077776/1619503660695-3596146f-a09c-4914-8667-f2f468e768a5.png#align=left&display=inline&height=779&margin=%5Bobject%20Object%5D&name=image.png&originHeight=1558&originWidth=2880&size=232348&status=done&style=none&width=1440#id=tu55k&originHeight=1558&originWidth=2880&originalType=binary&ratio=1&status=done&style=none)
