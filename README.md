# 基于 Hyperledger Fabric 的房地产交易系统

本项目是一个基于 Hyperledger Fabric 的房地产交易系统，实现了房产登记和交易的业务流程。

系统采用联盟链技术，由不动产登记机构、交易平台和银行三个组织共同维护。

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
    - 负责创建交易信息
    - 维护两个 Peer 节点：peer0.org3 和 peer1.org3

### 智能合约（Chaincode）

智能合约实现了以下核心功能：

1. 房产信息管理
    - 创建房产（仅不动产登记机构可操作）
    - 查询房产信息
    - 分页查询房产列表

2. 交易管理
    - 创建交易（仅交易平台可操作）
    - 完成交易（仅银行可操作）
    - 查询交易信息
    - 分页查询交易列表

### 应用服务器（Application）

API 接口设计：

```
/api/realty-agency
  POST /realty/create         # 创建房产信息（使用不动产登记机构身份操作）

/api/trading-platform
  POST /transaction/create    # 创建交易（使用交易平台身份操作）

/api/bank
  POST /transaction/complete/:txId  # 完成交易（使用银行身份操作）

/api/query
  GET  /realty/:id           # 查询房产信息
  GET  /realty/list          # 分页查询房产列表
    - pageSize: 每页记录数
    - bookmark: 分页标记
    - status: 房产状态（可选，NORMAL-正常、IN_TRANSACTION-交易中）
  GET  /transaction/:txId    # 查询交易信息
  GET  /transaction/list     # 分页查询交易列表
    - pageSize: 每页记录数
    - bookmark: 分页标记
    - status: 交易状态（可选，PENDING-待付款、COMPLETED-已完成）
```
