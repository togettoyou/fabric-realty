# 基于 Hyperledger Fabric 的房地产交易系统

本项目是一个基于 Hyperledger Fabric 的房地产交易系统，实现了房产登记和交易的业务流程。

系统采用联盟链技术，由不动产登记机构、交易平台和银行三个组织共同维护。

> 🎓 提供项目教学及问题解答服务，欢迎通过以下方式联系：

<img src="https://github.com/user-attachments/assets/ea93572c-6c05-4751-bde7-35a58fe083f1" width="520" alt="gopher云原生公众号二维码">

👆 扫码或搜索关注公众号：**gopher云原生**

## 本地开发

参考：[本地开发指南](dev.md)

推荐首次使用时选择快速部署方式，以便快速体验系统功能，往下看 👇

## 快速部署

### 环境要求

- Docker
- Docker Compose

### 部署步骤

1. 拉取项目（或手动下载）

   ```bash
   git clone --depth 1 https://github.com/togettoyou/fabric-realty.git
   ```

2. 设置脚本权限

   ```bash
   cd fabric-realty
   find . -name "*.sh" -exec chmod +x {} \;
   ```

3. 一键部署

   ```bash
   ./install.sh
   ```

4. 一键卸载

   ```bash
   ./uninstall.sh
   ```

### 访问服务

http://localhost:8000

### 系统演示

系统包含三个组织身份，每个组织都有独立的操作界面和权限

<img width="1337" alt="1" src="https://github.com/user-attachments/assets/185492e0-ac3f-419c-a64f-b17421046bc8" />

#### 业务操作流程

1. 房产登记上链
    - 不动产登记机构登录系统
    - 点击"登记新房产"，填写房产信息
    - 提交后，房产信息将上链保存

<img width="1337" alt="2" src="https://github.com/user-attachments/assets/e7474b46-f2f5-4561-91db-ed6f27ba858d" />

2. 发起房产交易
    - 交易平台登录系统
    - 点击"生成新交易"，填写交易信息
    - 提交后，交易信息将上链保存

<img width="1337" alt="3" src="https://github.com/user-attachments/assets/c3977cb0-cd48-495a-ab3b-6d244a81b6e0" />

3. 银行确认交易
    - 银行登录系统
    - 核实双方交易信息和资金状态
    - 点击"完成交易"，完成交易
    - 交易完成后，房产所有权将自动变更

<img width="1337" alt="4" src="https://github.com/user-attachments/assets/600cc2c2-52e5-4472-9e50-d18cabb27cf2" />

4. 区块链浏览
    - 所有组织都可以查看区块信息
    - 确保信息公开透明且不可篡改

<img width="1337" alt="5" src="https://github.com/user-attachments/assets/c1e9088a-b9dd-422d-95a1-5534243e471e" />

> 💡 系统特点：
> - 基于区块链的分布式账本，确保数据不可篡改
> - 基于智能合约的权限控制，确保操作安全可控
> - 每个组织只能执行自己权限范围内的操作
> - 所有操作都会记录在区块链上，可追溯、可审计

## 系统架构

### 网络架构（Network）

系统由三个组织构成的联盟链网络：

1. 不动产登记机构（Org1）
    - 负责房产信息的登记
    - 维护两个 Peer 节点：peer0.org1 和 peer1.org1

2. 银行（Org2）
    - 负责交易的完成确认
    - 维护两个 Peer 节点：peer0.org2 和 peer1.org2

3. 交易平台（Org3）
    - 负责生成交易信息
    - 维护两个 Peer 节点：peer0.org3 和 peer1.org3

### 智能合约（Chaincode）

智能合约实现了以下核心功能：

1. 房产信息管理
    - 创建房产（仅不动产登记机构可操作）
    - 查询房产信息
    - 分页查询房产列表

2. 交易管理
    - 生成交易（仅交易平台可操作）
    - 完成交易（仅银行可操作）
    - 查询交易信息
    - 分页查询交易列表

### 应用服务器（Application）

API 接口设计：

```
/api/realty-agency
  POST /realty/create         # 创建房产信息
  GET  /realty/:id           # 查询房产信息
  GET  /realty/list          # 分页查询房产列表
    - pageSize: 每页记录数
    - bookmark: 分页标记
    - status: 房产状态（可选，NORMAL-正常、IN_TRANSACTION-交易中）
  GET  /block/list           # 分页查询区块列表
    - pageSize: 每页记录数，默认10
    - pageNum: 页码，默认1

/api/trading-platform
  POST /transaction/create    # 生成交易
  GET  /realty/:id           # 查询房产信息
  GET  /transaction/:txId    # 查询交易信息
  GET  /transaction/list     # 分页查询交易列表
    - pageSize: 每页记录数
    - bookmark: 分页标记
    - status: 交易状态（可选，PENDING-待付款、COMPLETED-已完成）
  GET  /block/list           # 分页查询区块列表
    - pageSize: 每页记录数，默认10
    - pageNum: 页码，默认1

/api/bank
  POST /transaction/complete/:txId  # 完成交易
  GET  /transaction/:txId    # 查询交易信息
  GET  /transaction/list     # 分页查询交易列表
    - pageSize: 每页记录数
    - bookmark: 分页标记
    - status: 交易状态（可选，PENDING-待付款、COMPLETED-已完成）
  GET  /block/list           # 分页查询区块列表
    - pageSize: 每页记录数，默认10
    - pageNum: 页码，默认1
```

## 技术栈功能说明

### 区块链层

- Hyperledger Fabric v2.5.10
    - 分布式账本存储房产和交易数据
    - 智能合约实现业务逻辑和权限控制
    - 多组织（不动产登记机构、银行、交易平台）的联盟链网络

### 后端（Go）

- Gin v1.10.0
    - RESTful API 接口框架
- fabric-gateway v1.7.0
    - 提供区块链网络交互接口
    - 处理链码调用和查询
- fabric-protos-go-apiv2 v0.3.4
    - 处理区块链数据的序列化和反序列化
    - 提供区块链交互的协议支持

### 前端（Vue）

- Vue v3.3.8
    - 前端主框架
    - 响应式数据处理
- TypeScript v5.0.2
    - 类型检查
    - 代码提示和重构支持
- Vite v4.5.0
    - 开发服务器
    - 构建工具
- Ant Design Vue v3.2.20
    - UI 组件库
    - 提供完整的设计体系
