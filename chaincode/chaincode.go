package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// SmartContract 提供房地产交易的功能
type SmartContract struct {
	contractapi.Contract
}

// RealEstateStatus 房产状态
type RealEstateStatus string

const (
	ForSale       RealEstateStatus = "ForSale"       // 待售
	InTransaction RealEstateStatus = "InTransaction" // 交易中
	Sold          RealEstateStatus = "Sold"          // 已售
)

// TransactionStatus 交易状态
type TransactionStatus string

const (
	Pending   TransactionStatus = "Pending"   // 待付款
	InEscrow  TransactionStatus = "InEscrow"  // 已托管
	Completed TransactionStatus = "Completed" // 已完成
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

// CreateRealEstate 创建房产信息
func (s *SmartContract) CreateRealEstate(ctx contractapi.TransactionContextInterface, id string, address string, area float64, owner string, price float64) error {
	exists, err := s.RealEstateExists(ctx, id)
	if err != nil {
		return err
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
		Status:          ForSale,
		CreateTime:      time.Now(),
		UpdateTime:      time.Now(),
	}

	realEstateJSON, err := json.Marshal(realEstate)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, realEstateJSON)
}

// QueryRealEstate 查询房产信息
func (s *SmartContract) QueryRealEstate(ctx contractapi.TransactionContextInterface, id string) (*RealEstate, error) {
	realEstateJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read real estate: %v", err)
	}
	if realEstateJSON == nil {
		return nil, fmt.Errorf("房产ID %s 不存在", id)
	}

	var realEstate RealEstate
	err = json.Unmarshal(realEstateJSON, &realEstate)
	if err != nil {
		return nil, err
	}

	return &realEstate, nil
}

// CreateTransaction 创建交易
func (s *SmartContract) CreateTransaction(ctx contractapi.TransactionContextInterface, txID string, realEstateID string, seller string, buyer string, price float64) error {
	realEstate, err := s.QueryRealEstate(ctx, realEstateID)
	if err != nil {
		return err
	}

	if realEstate.Status != ForSale {
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
		Status:       Pending,
		CreateTime:   time.Now(),
		UpdateTime:   time.Now(),
	}

	transactionJSON, err := json.Marshal(transaction)
	if err != nil {
		return err
	}

	// 更新房产状态
	realEstate.Status = InTransaction
	realEstate.UpdateTime = time.Now()
	realEstateJSON, err := json.Marshal(realEstate)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(realEstateID, realEstateJSON)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(txID, transactionJSON)
}

// ConfirmEscrow 确认资金托管
func (s *SmartContract) ConfirmEscrow(ctx contractapi.TransactionContextInterface, txID string) error {
	transactionJSON, err := ctx.GetStub().GetState(txID)
	if err != nil {
		return fmt.Errorf("failed to read transaction: %v", err)
	}
	if transactionJSON == nil {
		return fmt.Errorf("交易ID %s 不存在", txID)
	}

	var transaction Transaction
	err = json.Unmarshal(transactionJSON, &transaction)
	if err != nil {
		return err
	}

	if transaction.Status != Pending {
		return fmt.Errorf("交易状态不正确")
	}

	transaction.Status = InEscrow
	transaction.UpdateTime = time.Now()

	transactionJSON, err = json.Marshal(transaction)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(txID, transactionJSON)
}

// CompleteTransaction 完成交易
func (s *SmartContract) CompleteTransaction(ctx contractapi.TransactionContextInterface, txID string) error {
	transactionJSON, err := ctx.GetStub().GetState(txID)
	if err != nil {
		return fmt.Errorf("failed to read transaction: %v", err)
	}
	if transactionJSON == nil {
		return fmt.Errorf("交易ID %s 不存在", txID)
	}

	var transaction Transaction
	err = json.Unmarshal(transactionJSON, &transaction)
	if err != nil {
		return err
	}

	if transaction.Status != InEscrow {
		return fmt.Errorf("交易状态不正确")
	}

	// 更新房产所有权
	realEstate, err := s.QueryRealEstate(ctx, transaction.RealEstateID)
	if err != nil {
		return err
	}

	realEstate.CurrentOwner = transaction.Buyer
	realEstate.Status = Sold
	realEstate.UpdateTime = time.Now()

	realEstateJSON, err := json.Marshal(realEstate)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(transaction.RealEstateID, realEstateJSON)
	if err != nil {
		return err
	}

	// 更新交易状态
	transaction.Status = Completed
	transaction.UpdateTime = time.Now()

	transactionJSON, err = json.Marshal(transaction)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(txID, transactionJSON)
}

// RealEstateExists 检查房产是否存在
func (s *SmartContract) RealEstateExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	realEstateJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read real estate: %v", err)
	}
	return realEstateJSON != nil, nil
}

// Hello 用于验证
func (s *SmartContract) Hello(ctx contractapi.TransactionContextInterface) (string, error) {
	return "hello", nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("Error creating chaincode: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting chaincode: %v", err)
	}
}
