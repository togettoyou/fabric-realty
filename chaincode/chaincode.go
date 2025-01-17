package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/v2/pkg/cid"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// SmartContract 提供房地产交易的功能
type SmartContract struct {
	contractapi.Contract
}

// RealEstateStatus 房产状态
type RealEstateStatus string

const (
	FOR_SALE       RealEstateStatus = "FOR_SALE"       // 待售
	IN_TRANSACTION RealEstateStatus = "IN_TRANSACTION" // 交易中
	SOLD           RealEstateStatus = "SOLD"           // 已售出
)

// TransactionStatus 交易状态
type TransactionStatus string

const (
	PENDING   TransactionStatus = "PENDING"   // 待付款
	IN_ESCROW TransactionStatus = "IN_ESCROW" // 已托管
	COMPLETED TransactionStatus = "COMPLETED" // 已完成
)

// RealEstate 房产信息
type RealEstate struct {
	ID              string           `json:"id"`              // 房产ID
	PropertyAddress string           `json:"propertyAddress"` // 房产地址
	Area            float64          `json:"area"`            // 面积
	CurrentOwner    string           `json:"currentOwner"`    // 当前所有者
	Price           float64          `json:"price"`           // 价格
	Status          RealEstateStatus `json:"status"`          // 状态
	CreateTime      time.Time        `json:"createTime"`      // 创建时间
	UpdateTime      time.Time        `json:"updateTime"`      // 更新时间
}

// Transaction 交易信息
type Transaction struct {
	ID           string            `json:"id"`           // 交易ID
	RealEstateID string            `json:"realEstateId"` // 房产ID
	Seller       string            `json:"seller"`       // 卖家
	Buyer        string            `json:"buyer"`        // 买家
	Price        float64           `json:"price"`        // 成交价格
	Status       TransactionStatus `json:"status"`       // 状态
	CreateTime   time.Time         `json:"createTime"`   // 创建时间
	UpdateTime   time.Time         `json:"updateTime"`   // 更新时间
}

// 组织 MSP ID 常量
const (
	REALTY_ORG_MSPID = "Org1MSP" // 房管局组织 MSP ID
	BANK_ORG_MSPID   = "Org2MSP" // 银行组织 MSP ID
)

// GetClientIdentityMSPID 获取客户端身份信息
func (s *SmartContract) GetClientIdentityMSPID(ctx contractapi.TransactionContextInterface) (string, error) {
	clientID, err := cid.New(ctx.GetStub())
	if err != nil {
		return "", fmt.Errorf("获取客户端身份信息失败：%v", err)
	}
	return clientID.GetMSPID()
}

// CreateRealEstate 创建房产信息
func (s *SmartContract) CreateRealEstate(ctx contractapi.TransactionContextInterface, id string, address string, area float64, owner string, price float64, createTime time.Time) error {
	// 检查调用者身份
	clientMSPID, err := s.GetClientIdentityMSPID(ctx)
	if err != nil {
		return fmt.Errorf("获取调用者身份失败：%v", err)
	}

	log.Println("clientMSPID", clientMSPID, REALTY_ORG_MSPID)

	// 验证是否是房管局组织的成员
	if clientMSPID != REALTY_ORG_MSPID {
		return fmt.Errorf("只有房管局组织成员才能创建房产信息")
	}

	exists, err := s.RealEstateExists(ctx, id)
	if err != nil {
		return fmt.Errorf("检查房产是否存在时发生错误：%v", err)
	}
	if exists {
		return fmt.Errorf("房产ID %s 已存在", id)
	}

	realEstate := RealEstate{
		ID:              id,
		PropertyAddress: address,
		Area:            area,
		CurrentOwner:    owner,
		Price:           price,
		Status:          FOR_SALE,
		CreateTime:      createTime,
		UpdateTime:      createTime,
	}

	realEstateJSON, err := json.Marshal(realEstate)
	if err != nil {
		return fmt.Errorf("序列化房产信息失败：%v", err)
	}

	return ctx.GetStub().PutState(id, realEstateJSON)
}

// QueryRealEstate 查询房产信息
func (s *SmartContract) QueryRealEstate(ctx contractapi.TransactionContextInterface, id string) (*RealEstate, error) {
	realEstateJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("读取房产信息失败：%v", err)
	}
	if realEstateJSON == nil {
		return nil, fmt.Errorf("房产ID %s 不存在", id)
	}

	var realEstate RealEstate
	err = json.Unmarshal(realEstateJSON, &realEstate)
	if err != nil {
		return nil, fmt.Errorf("解析房产信息失败：%v", err)
	}

	return &realEstate, nil
}

// CreateTransaction 创建交易
func (s *SmartContract) CreateTransaction(ctx contractapi.TransactionContextInterface, txID string, realEstateID string, seller string, buyer string, price float64, createTime time.Time) error {
	realEstate, err := s.QueryRealEstate(ctx, realEstateID)
	if err != nil {
		return fmt.Errorf("查询房产信息失败：%v", err)
	}

	if realEstate.Status != FOR_SALE {
		return fmt.Errorf("房产不在可售状态")
	}

	if realEstate.CurrentOwner != seller {
		return fmt.Errorf("卖家不是房产所有者")
	}

	transaction := Transaction{
		ID:           txID,
		RealEstateID: realEstateID,
		Seller:       seller,
		Buyer:        buyer,
		Price:        price,
		Status:       PENDING,
		CreateTime:   createTime,
		UpdateTime:   createTime,
	}

	transactionJSON, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("序列化交易信息失败：%v", err)
	}

	// 更新房产状态
	realEstate.Status = IN_TRANSACTION
	realEstate.UpdateTime = createTime
	realEstateJSON, err := json.Marshal(realEstate)
	if err != nil {
		return fmt.Errorf("序列化房产信息失败：%v", err)
	}

	err = ctx.GetStub().PutState(realEstateID, realEstateJSON)
	if err != nil {
		return fmt.Errorf("更新房产状态失败：%v", err)
	}

	return ctx.GetStub().PutState(txID, transactionJSON)
}

// ConfirmEscrow 确认资金托管
func (s *SmartContract) ConfirmEscrow(ctx contractapi.TransactionContextInterface, txID string, updateTime time.Time) error {
	// 检查调用者身份
	clientMSPID, err := s.GetClientIdentityMSPID(ctx)
	if err != nil {
		return fmt.Errorf("获取调用者身份失败：%v", err)
	}

	log.Println("clientMSPID", clientMSPID, BANK_ORG_MSPID)

	// 验证是否是银行组织的成员
	if clientMSPID != BANK_ORG_MSPID {
		return fmt.Errorf("只有银行组织成员才能确认资金托管")
	}

	transactionJSON, err := ctx.GetStub().GetState(txID)
	if err != nil {
		return fmt.Errorf("读取交易信息失败：%v", err)
	}
	if transactionJSON == nil {
		return fmt.Errorf("交易ID %s 不存在", txID)
	}

	var transaction Transaction
	err = json.Unmarshal(transactionJSON, &transaction)
	if err != nil {
		return fmt.Errorf("解析交易信息失败：%v", err)
	}

	if transaction.Status != PENDING {
		return fmt.Errorf("交易状态不正确，当前状态：%s，需要状态：待付款", transaction.Status)
	}

	transaction.Status = IN_ESCROW
	transaction.UpdateTime = updateTime

	transactionJSON, err = json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("序列化交易信息失败：%v", err)
	}

	return ctx.GetStub().PutState(txID, transactionJSON)
}

// CompleteTransaction 完成交易
func (s *SmartContract) CompleteTransaction(ctx contractapi.TransactionContextInterface, txID string, updateTime time.Time) error {
	// 检查调用者身份
	clientMSPID, err := s.GetClientIdentityMSPID(ctx)
	if err != nil {
		return fmt.Errorf("获取调用者身份失败：%v", err)
	}

	log.Println("clientMSPID", clientMSPID, BANK_ORG_MSPID)

	// 验证是否是银行组织的成员
	if clientMSPID != BANK_ORG_MSPID {
		return fmt.Errorf("只有银行组织成员才能完成交易")
	}

	transactionJSON, err := ctx.GetStub().GetState(txID)
	if err != nil {
		return fmt.Errorf("读取交易信息失败：%v", err)
	}
	if transactionJSON == nil {
		return fmt.Errorf("交易ID %s 不存在", txID)
	}

	var transaction Transaction
	err = json.Unmarshal(transactionJSON, &transaction)
	if err != nil {
		return fmt.Errorf("解析交易信息失败：%v", err)
	}

	if transaction.Status != IN_ESCROW {
		return fmt.Errorf("交易状态不正确，当前状态：%s，需要状态：已托管", transaction.Status)
	}

	// 更新房产所有权
	realEstate, err := s.QueryRealEstate(ctx, transaction.RealEstateID)
	if err != nil {
		return fmt.Errorf("查询房产信息失败：%v", err)
	}

	realEstate.CurrentOwner = transaction.Buyer
	realEstate.Status = SOLD
	realEstate.UpdateTime = updateTime

	realEstateJSON, err := json.Marshal(realEstate)
	if err != nil {
		return fmt.Errorf("序列化房产信息失败：%v", err)
	}

	err = ctx.GetStub().PutState(transaction.RealEstateID, realEstateJSON)
	if err != nil {
		return fmt.Errorf("更新房产信息失败：%v", err)
	}

	// 更新交易状态
	transaction.Status = COMPLETED
	transaction.UpdateTime = updateTime

	transactionJSON, err = json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("序列化交易信息失败：%v", err)
	}

	return ctx.GetStub().PutState(txID, transactionJSON)
}

// RealEstateExists 检查房产是否存在
func (s *SmartContract) RealEstateExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	realEstateJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("读取房产信息失败：%v", err)
	}
	return realEstateJSON != nil, nil
}

// Hello 用于验证
func (s *SmartContract) Hello(ctx contractapi.TransactionContextInterface) (string, error) {
	return "hello", nil
}

// InitLedger 初始化账本
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	log.Println("InitLedger")
	return nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("创建智能合约失败：%v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("启动智能合约失败：%v", err)
	}
}
