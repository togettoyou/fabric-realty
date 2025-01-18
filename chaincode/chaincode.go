package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/v2/pkg/cid"
	"github.com/hyperledger/fabric-chaincode-go/v2/shim"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
	"github.com/hyperledger/fabric-protos-go-apiv2/peer"
)

// SmartContract 提供房地产交易的功能
type SmartContract struct {
	contractapi.Contract
}

// 文档类型常量（用于创建复合键）
const (
	REAL_ESTATE = "RE" // 房产信息
	TRANSACTION = "TX" // 交易信息
)

// 索引前缀常量
const (
	IDX_RE_OWNER  = "IDX_RE_OWNER"  // 房产所有者索引
	IDX_RE_STATUS = "IDX_RE_STATUS" // 房产状态索引
	IDX_TX_SELLER = "IDX_TX_SELLER" // 交易卖家索引
	IDX_TX_BUYER  = "IDX_TX_BUYER"  // 交易买家索引
	IDX_TX_STATUS = "IDX_TX_STATUS" // 交易状态索引
)

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
	// 参数验证
	if len(id) == 0 {
		return fmt.Errorf("房产ID不能为空")
	}
	if len(address) == 0 {
		return fmt.Errorf("房产地址不能为空")
	}
	if area <= 0 {
		return fmt.Errorf("面积必须大于0")
	}
	if len(owner) == 0 {
		return fmt.Errorf("所有者不能为空")
	}
	if price <= 0 {
		return fmt.Errorf("价格必须大于0")
	}

	// 检查调用者身份
	clientMSPID, err := s.GetClientIdentityMSPID(ctx)
	if err != nil {
		return fmt.Errorf("获取调用者身份失败：%v", err)
	}

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

	// 创建主键
	key, err := ctx.GetStub().CreateCompositeKey(REAL_ESTATE, []string{id})
	if err != nil {
		return fmt.Errorf("创建复合键失败：%v", err)
	}

	// 创建所有者索引
	ownerKey, err := ctx.GetStub().CreateCompositeKey(IDX_RE_OWNER, []string{owner, id})
	if err != nil {
		return fmt.Errorf("创建所有者索引失败：%v", err)
	}

	// 创建状态索引
	statusKey, err := ctx.GetStub().CreateCompositeKey(IDX_RE_STATUS, []string{string(FOR_SALE), id})
	if err != nil {
		return fmt.Errorf("创建状态索引失败：%v", err)
	}

	// 写入状态
	err = ctx.GetStub().PutState(key, realEstateJSON)
	if err != nil {
		return fmt.Errorf("保存房产信息失败：%v", err)
	}

	err = ctx.GetStub().PutState(ownerKey, []byte{0x00})
	if err != nil {
		return fmt.Errorf("保存所有者索引失败：%v", err)
	}

	err = ctx.GetStub().PutState(statusKey, []byte{0x00})
	if err != nil {
		return fmt.Errorf("保存状态索引失败：%v", err)
	}

	return nil
}

// QueryRealEstate 查询房产信息
func (s *SmartContract) QueryRealEstate(ctx contractapi.TransactionContextInterface, id string) (*RealEstate, error) {
	key, err := ctx.GetStub().CreateCompositeKey(REAL_ESTATE, []string{id})
	if err != nil {
		return nil, fmt.Errorf("创建复合键失败：%v", err)
	}

	realEstateJSON, err := ctx.GetStub().GetState(key)
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
	// 参数验证
	if len(txID) == 0 {
		return fmt.Errorf("交易ID不能为空")
	}
	if len(realEstateID) == 0 {
		return fmt.Errorf("房产ID不能为空")
	}
	if len(seller) == 0 {
		return fmt.Errorf("卖家不能为空")
	}
	if len(buyer) == 0 {
		return fmt.Errorf("买家不能为空")
	}
	if seller == buyer {
		return fmt.Errorf("买家和卖家不能是同一人")
	}
	if price <= 0 {
		return fmt.Errorf("价格必须大于0")
	}

	realEstate, err := s.QueryRealEstate(ctx, realEstateID)
	if err != nil {
		return fmt.Errorf("查询房产信息失败：%v", err)
	}

	if realEstate.Status != FOR_SALE {
		return fmt.Errorf("房产状态不正确，当前状态：%s，需要状态：%s", realEstate.Status, FOR_SALE)
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

	// 创建交易主键
	txKey, err := ctx.GetStub().CreateCompositeKey(TRANSACTION, []string{txID})
	if err != nil {
		return fmt.Errorf("创建交易复合键失败：%v", err)
	}

	// 创建卖家索引
	sellerKey, err := ctx.GetStub().CreateCompositeKey(IDX_TX_SELLER, []string{seller, txID})
	if err != nil {
		return fmt.Errorf("创建卖家索引失败：%v", err)
	}

	// 创建买家索引
	buyerKey, err := ctx.GetStub().CreateCompositeKey(IDX_TX_BUYER, []string{buyer, txID})
	if err != nil {
		return fmt.Errorf("创建买家索引失败：%v", err)
	}

	// 创建状态索引
	statusKey, err := ctx.GetStub().CreateCompositeKey(IDX_TX_STATUS, []string{string(PENDING), txID})
	if err != nil {
		return fmt.Errorf("创建状态索引失败：%v", err)
	}

	// 更新房产状态
	err = s.updateRealEstateStatus(ctx, realEstateID, FOR_SALE, IN_TRANSACTION)
	if err != nil {
		return fmt.Errorf("更新房产状态失败：%v", err)
	}

	// 写入交易相关状态
	err = ctx.GetStub().PutState(txKey, transactionJSON)
	if err != nil {
		return fmt.Errorf("保存交易信息失败：%v", err)
	}

	err = ctx.GetStub().PutState(sellerKey, []byte{0x00})
	if err != nil {
		return fmt.Errorf("保存卖家索引失败：%v", err)
	}

	err = ctx.GetStub().PutState(buyerKey, []byte{0x00})
	if err != nil {
		return fmt.Errorf("保存买家索引失败：%v", err)
	}

	err = ctx.GetStub().PutState(statusKey, []byte{0x00})
	if err != nil {
		return fmt.Errorf("保存状态索引失败：%v", err)
	}

	return nil
}

// updateRealEstateStatus 更新房产状态（内部函数）
func (s *SmartContract) updateRealEstateStatus(ctx contractapi.TransactionContextInterface, realEstateID string, oldStatus RealEstateStatus, newStatus RealEstateStatus) error {
	// 删除旧的状态索引
	oldStatusKey, err := ctx.GetStub().CreateCompositeKey(IDX_RE_STATUS, []string{string(oldStatus), realEstateID})
	if err != nil {
		return fmt.Errorf("创建旧状态索引失败：%v", err)
	}

	// 创建新的状态索引
	newStatusKey, err := ctx.GetStub().CreateCompositeKey(IDX_RE_STATUS, []string{string(newStatus), realEstateID})
	if err != nil {
		return fmt.Errorf("创建新状态索引失败：%v", err)
	}

	// 更新房产信息
	realEstate, err := s.QueryRealEstate(ctx, realEstateID)
	if err != nil {
		return fmt.Errorf("查询房产信息失败：%v", err)
	}

	realEstate.Status = newStatus
	realEstate.UpdateTime = time.Now()

	realEstateJSON, err := json.Marshal(realEstate)
	if err != nil {
		return fmt.Errorf("序列化房产信息失败：%v", err)
	}

	// 创建房产主键
	reKey, err := ctx.GetStub().CreateCompositeKey(REAL_ESTATE, []string{realEstateID})
	if err != nil {
		return fmt.Errorf("创建房产复合键失败：%v", err)
	}

	// 更新状态
	err = ctx.GetStub().DelState(oldStatusKey)
	if err != nil {
		return fmt.Errorf("删除旧状态索引失败：%v", err)
	}

	err = ctx.GetStub().PutState(reKey, realEstateJSON)
	if err != nil {
		return fmt.Errorf("更新房产信息失败：%v", err)
	}

	err = ctx.GetStub().PutState(newStatusKey, []byte{0x00})
	if err != nil {
		return fmt.Errorf("保存新状态索引失败：%v", err)
	}

	return nil
}

// updateTransactionStatus 更新交易状态（内部函数）
func (s *SmartContract) updateTransactionStatus(ctx contractapi.TransactionContextInterface, txID string, oldStatus TransactionStatus, newStatus TransactionStatus) error {
	// 删除旧的状态索引
	oldStatusKey, err := ctx.GetStub().CreateCompositeKey(IDX_TX_STATUS, []string{string(oldStatus), txID})
	if err != nil {
		return fmt.Errorf("创建旧状态索引失败：%v", err)
	}

	// 创建新的状态索引
	newStatusKey, err := ctx.GetStub().CreateCompositeKey(IDX_TX_STATUS, []string{string(newStatus), txID})
	if err != nil {
		return fmt.Errorf("创建新状态索引失败：%v", err)
	}

	// 更新交易信息
	transaction, err := s.QueryTransaction(ctx, txID)
	if err != nil {
		return fmt.Errorf("查询交易信息失败：%v", err)
	}

	transaction.Status = newStatus
	transaction.UpdateTime = time.Now()

	transactionJSON, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("序列化交易信息失败：%v", err)
	}

	// 创建交易主键
	txKey, err := ctx.GetStub().CreateCompositeKey(TRANSACTION, []string{txID})
	if err != nil {
		return fmt.Errorf("创建交易复合键失败：%v", err)
	}

	// 更新状态
	err = ctx.GetStub().DelState(oldStatusKey)
	if err != nil {
		return fmt.Errorf("删除旧状态索引失败：%v", err)
	}

	err = ctx.GetStub().PutState(txKey, transactionJSON)
	if err != nil {
		return fmt.Errorf("更新交易信息失败：%v", err)
	}

	err = ctx.GetStub().PutState(newStatusKey, []byte{0x00})
	if err != nil {
		return fmt.Errorf("保存新状态索引失败：%v", err)
	}

	return nil
}

// ConfirmEscrow 确认资金托管
func (s *SmartContract) ConfirmEscrow(ctx contractapi.TransactionContextInterface, txID string, updateTime time.Time) error {
	// 检查调用者身份
	clientMSPID, err := s.GetClientIdentityMSPID(ctx)
	if err != nil {
		return fmt.Errorf("获取调用者身份失败：%v", err)
	}

	// 验证是否是银行组织的成员
	if clientMSPID != BANK_ORG_MSPID {
		return fmt.Errorf("只有银行组织成员才能确认资金托管")
	}

	// 创建交易的复合键
	txKey, err := ctx.GetStub().CreateCompositeKey(TRANSACTION, []string{txID})
	if err != nil {
		return fmt.Errorf("创建交易复合键失败：%v", err)
	}

	transactionJSON, err := ctx.GetStub().GetState(txKey)
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
		return fmt.Errorf("交易状态不正确，当前状态：%s，需要状态：%s", transaction.Status, PENDING)
	}

	transaction.Status = IN_ESCROW
	transaction.UpdateTime = updateTime

	transactionJSON, err = json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("序列化交易信息失败：%v", err)
	}

	return ctx.GetStub().PutState(txKey, transactionJSON)
}

// CompleteTransaction 完成交易
func (s *SmartContract) CompleteTransaction(ctx contractapi.TransactionContextInterface, txID string, updateTime time.Time) error {
	// 检查调用者身份
	clientMSPID, err := s.GetClientIdentityMSPID(ctx)
	if err != nil {
		return fmt.Errorf("获取调用者身份失败：%v", err)
	}

	// 验证是否是银行组织的成员
	if clientMSPID != BANK_ORG_MSPID {
		return fmt.Errorf("只有银行组织成员才能完成交易")
	}

	// 创建交易的复合键
	txKey, err := ctx.GetStub().CreateCompositeKey(TRANSACTION, []string{txID})
	if err != nil {
		return fmt.Errorf("创建交易复合键失败：%v", err)
	}

	transactionJSON, err := ctx.GetStub().GetState(txKey)
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
		return fmt.Errorf("交易状态不正确，当前状态：%s，需要状态：%s", transaction.Status, IN_ESCROW)
	}

	// 更新房产所有权
	realEstate, err := s.QueryRealEstate(ctx, transaction.RealEstateID)
	if err != nil {
		return fmt.Errorf("查询房产信息失败：%v", err)
	}

	// 删除旧的所有者索引
	oldOwnerKey, err := ctx.GetStub().CreateCompositeKey(IDX_RE_OWNER, []string{realEstate.CurrentOwner, realEstate.ID})
	if err != nil {
		return fmt.Errorf("创建旧所有者索引失败：%v", err)
	}

	// 删除旧的状态索引
	oldStatusKey, err := ctx.GetStub().CreateCompositeKey(IDX_RE_STATUS, []string{string(realEstate.Status), realEstate.ID})
	if err != nil {
		return fmt.Errorf("创建旧状态索引失败：%v", err)
	}

	// 更新房产信息
	realEstate.CurrentOwner = transaction.Buyer
	realEstate.Status = SOLD
	realEstate.UpdateTime = updateTime

	// 创建新的所有者索引
	newOwnerKey, err := ctx.GetStub().CreateCompositeKey(IDX_RE_OWNER, []string{realEstate.CurrentOwner, realEstate.ID})
	if err != nil {
		return fmt.Errorf("创建新所有者索引失败：%v", err)
	}

	// 创建新的状态索引
	newStatusKey, err := ctx.GetStub().CreateCompositeKey(IDX_RE_STATUS, []string{string(SOLD), realEstate.ID})
	if err != nil {
		return fmt.Errorf("创建新状态索引失败：%v", err)
	}

	realEstateJSON, err := json.Marshal(realEstate)
	if err != nil {
		return fmt.Errorf("序列化房产信息失败：%v", err)
	}

	// 创建房产的复合键
	reKey, err := ctx.GetStub().CreateCompositeKey(REAL_ESTATE, []string{transaction.RealEstateID})
	if err != nil {
		return fmt.Errorf("创建房产复合键失败：%v", err)
	}

	// 更新所有状态
	err = ctx.GetStub().DelState(oldOwnerKey)
	if err != nil {
		return fmt.Errorf("删除旧所有者索引失败：%v", err)
	}

	err = ctx.GetStub().DelState(oldStatusKey)
	if err != nil {
		return fmt.Errorf("删除旧状态索引失败：%v", err)
	}

	err = ctx.GetStub().PutState(reKey, realEstateJSON)
	if err != nil {
		return fmt.Errorf("更新房产信息失败：%v", err)
	}

	err = ctx.GetStub().PutState(newOwnerKey, []byte{0x00})
	if err != nil {
		return fmt.Errorf("保存新所有者索引失败：%v", err)
	}

	err = ctx.GetStub().PutState(newStatusKey, []byte{0x00})
	if err != nil {
		return fmt.Errorf("保存新状态索引失败：%v", err)
	}

	// 更新交易状态
	err = s.updateTransactionStatus(ctx, txID, IN_ESCROW, COMPLETED)
	if err != nil {
		return fmt.Errorf("更新交易状态失败：%v", err)
	}

	return nil
}

// QueryTransaction 查询交易信息
func (s *SmartContract) QueryTransaction(ctx contractapi.TransactionContextInterface, txID string) (*Transaction, error) {
	key, err := ctx.GetStub().CreateCompositeKey(TRANSACTION, []string{txID})
	if err != nil {
		return nil, fmt.Errorf("创建复合键失败：%v", err)
	}

	transactionJSON, err := ctx.GetStub().GetState(key)
	if err != nil {
		return nil, fmt.Errorf("读取交易信息失败：%v", err)
	}
	if transactionJSON == nil {
		return nil, fmt.Errorf("交易ID %s 不存在", txID)
	}

	var transaction Transaction
	err = json.Unmarshal(transactionJSON, &transaction)
	if err != nil {
		return nil, fmt.Errorf("解析交易信息失败：%v", err)
	}

	return &transaction, nil
}

// RealEstateExists 检查房产是否存在
func (s *SmartContract) RealEstateExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	key, err := ctx.GetStub().CreateCompositeKey(REAL_ESTATE, []string{id})
	if err != nil {
		return false, fmt.Errorf("创建复合键失败：%v", err)
	}

	realEstateJSON, err := ctx.GetStub().GetState(key)
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

// QueryRealEstateList 分页查询房产列表
func (s *SmartContract) QueryRealEstateList(ctx contractapi.TransactionContextInterface, pageSize int32, bookmark string) (*QueryResult, error) {
	// 使用部分复合键查询
	iterator, metadata, err := ctx.GetStub().GetStateByPartialCompositeKeyWithPagination(REAL_ESTATE, []string{}, pageSize, bookmark)
	if err != nil {
		return nil, fmt.Errorf("查询房产列表失败：%v", err)
	}
	defer iterator.Close()

	var realEstates []interface{}
	for iterator.HasNext() {
		queryResponse, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("获取下一条记录失败：%v", err)
		}

		var realEstate RealEstate
		err = json.Unmarshal(queryResponse.Value, &realEstate)
		if err != nil {
			return nil, fmt.Errorf("解析房产数据失败：%v", err)
		}
		realEstates = append(realEstates, realEstate)
	}

	return &QueryResult{
		Records:             realEstates,
		RecordsCount:        int32(len(realEstates)),
		Bookmark:            metadata.Bookmark,
		FetchedRecordsCount: metadata.FetchedRecordsCount,
	}, nil
}

// QueryTransactionList 分页查询交易列表
func (s *SmartContract) QueryTransactionList(ctx contractapi.TransactionContextInterface, pageSize int32, bookmark string) (*QueryResult, error) {
	// 使用部分复合键查询
	iterator, metadata, err := ctx.GetStub().GetStateByPartialCompositeKeyWithPagination(TRANSACTION, []string{}, pageSize, bookmark)
	if err != nil {
		return nil, fmt.Errorf("查询交易列表失败：%v", err)
	}
	defer iterator.Close()

	var transactions []interface{}
	for iterator.HasNext() {
		queryResponse, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("获取下一条记录失败：%v", err)
		}

		var transaction Transaction
		err = json.Unmarshal(queryResponse.Value, &transaction)
		if err != nil {
			return nil, fmt.Errorf("解析交易数据失败：%v", err)
		}
		transactions = append(transactions, transaction)
	}

	return &QueryResult{
		Records:             transactions,
		RecordsCount:        int32(len(transactions)),
		Bookmark:            metadata.Bookmark,
		FetchedRecordsCount: metadata.FetchedRecordsCount,
	}, nil
}

// QueryResult 查询结果
type QueryResult struct {
	Records             []interface{} `json:"records"`             // 记录列表
	RecordsCount        int32         `json:"recordsCount"`        // 本次返回的记录数
	Bookmark            string        `json:"bookmark"`            // 书签，用于下一页查询
	FetchedRecordsCount int32         `json:"fetchedRecordsCount"` // 总共获取的记录数
}

// QueryRealEstateByFilter 按条件查询房产列表
func (s *SmartContract) QueryRealEstateByFilter(ctx contractapi.TransactionContextInterface, owner string, status RealEstateStatus, pageSize int32, bookmark string) (*QueryResult, error) {
	// 验证状态是否有效（如果提供了状态）
	if len(status) > 0 {
		switch status {
		case FOR_SALE, IN_TRANSACTION, SOLD:
			// 有效状态
		default:
			return nil, fmt.Errorf("无效的状态：%s", status)
		}
	}

	var iterator shim.StateQueryIteratorInterface
	var metadata *peer.QueryResponseMetadata
	var err error

	if len(owner) > 0 {
		// 按所有者查询
		iterator, metadata, err = ctx.GetStub().GetStateByPartialCompositeKeyWithPagination(
			IDX_RE_OWNER,
			[]string{owner},
			pageSize,
			bookmark,
		)
	} else if len(status) > 0 {
		// 按状态查询
		iterator, metadata, err = ctx.GetStub().GetStateByPartialCompositeKeyWithPagination(
			IDX_RE_STATUS,
			[]string{string(status)},
			pageSize,
			bookmark,
		)
	} else {
		// 无筛选条件，查询所有房产
		iterator, metadata, err = ctx.GetStub().GetStateByPartialCompositeKeyWithPagination(
			REAL_ESTATE,
			[]string{},
			pageSize,
			bookmark,
		)
	}

	if err != nil {
		return nil, fmt.Errorf("查询房产列表失败：%v", err)
	}
	defer iterator.Close()

	var realEstates []interface{}
	for iterator.HasNext() {
		queryResponse, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("获取下一条记录失败：%v", err)
		}

		// 从复合键中提取房产ID
		_, compositeKeyParts, err := ctx.GetStub().SplitCompositeKey(queryResponse.Key)
		if err != nil {
			return nil, fmt.Errorf("解析复合键失败：%v", err)
		}

		var realEstateID string
		if len(owner) > 0 || len(status) > 0 {
			// 使用索引时，ID在第二个位置
			if len(compositeKeyParts) < 2 {
				return nil, fmt.Errorf("复合键格式错误")
			}
			realEstateID = compositeKeyParts[1]
		} else {
			// 使用主键时，ID在第一个位置
			if len(compositeKeyParts) < 1 {
				return nil, fmt.Errorf("复合键格式错误")
			}
			realEstateID = compositeKeyParts[0]
		}

		// 获取房产信息
		realEstate, err := s.QueryRealEstate(ctx, realEstateID)
		if err != nil {
			return nil, fmt.Errorf("查询房产信息失败：%v", err)
		}

		// 如果指定了状态，需要进行过滤
		if len(status) > 0 && realEstate.Status != status {
			continue
		}

		realEstates = append(realEstates, realEstate)
	}

	return &QueryResult{
		Records:             realEstates,
		RecordsCount:        int32(len(realEstates)),
		Bookmark:            metadata.Bookmark,
		FetchedRecordsCount: metadata.FetchedRecordsCount,
	}, nil
}

// QueryTransactionByFilter 按条件查询交易列表
func (s *SmartContract) QueryTransactionByFilter(ctx contractapi.TransactionContextInterface, seller string, buyer string, status TransactionStatus, pageSize int32, bookmark string) (*QueryResult, error) {
	// 验证状态是否有效（如果提供了状态）
	if len(status) > 0 {
		switch status {
		case PENDING, IN_ESCROW, COMPLETED:
			// 有效状态
		default:
			return nil, fmt.Errorf("无效的状态：%s", status)
		}
	}

	var iterator shim.StateQueryIteratorInterface
	var metadata *peer.QueryResponseMetadata
	var err error

	if len(seller) > 0 {
		// 按卖家查询
		iterator, metadata, err = ctx.GetStub().GetStateByPartialCompositeKeyWithPagination(
			IDX_TX_SELLER,
			[]string{seller},
			pageSize,
			bookmark,
		)
	} else if len(buyer) > 0 {
		// 按买家查询
		iterator, metadata, err = ctx.GetStub().GetStateByPartialCompositeKeyWithPagination(
			IDX_TX_BUYER,
			[]string{buyer},
			pageSize,
			bookmark,
		)
	} else if len(status) > 0 {
		// 按状态查询
		iterator, metadata, err = ctx.GetStub().GetStateByPartialCompositeKeyWithPagination(
			IDX_TX_STATUS,
			[]string{string(status)},
			pageSize,
			bookmark,
		)
	} else {
		// 无筛选条件，查询所有交易
		iterator, metadata, err = ctx.GetStub().GetStateByPartialCompositeKeyWithPagination(
			TRANSACTION,
			[]string{},
			pageSize,
			bookmark,
		)
	}

	if err != nil {
		return nil, fmt.Errorf("查询交易列表失败：%v", err)
	}
	defer iterator.Close()

	var transactions []interface{}
	for iterator.HasNext() {
		queryResponse, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("获取下一条记录失败：%v", err)
		}

		// 从复合键中提取交易ID
		_, compositeKeyParts, err := ctx.GetStub().SplitCompositeKey(queryResponse.Key)
		if err != nil {
			return nil, fmt.Errorf("解析复合键失败：%v", err)
		}

		var txID string
		if len(seller) > 0 || len(buyer) > 0 || len(status) > 0 {
			// 使用索引时，ID在第二个位置
			if len(compositeKeyParts) < 2 {
				return nil, fmt.Errorf("复合键格式错误")
			}
			txID = compositeKeyParts[1]
		} else {
			// 使用主键时，ID在第一个位置
			if len(compositeKeyParts) < 1 {
				return nil, fmt.Errorf("复合键格式错误")
			}
			txID = compositeKeyParts[0]
		}

		// 获取交易信息
		transaction, err := s.QueryTransaction(ctx, txID)
		if err != nil {
			return nil, fmt.Errorf("查询交易信息失败：%v", err)
		}

		// 如果指定了状态，需要进行过滤
		if len(status) > 0 && transaction.Status != status {
			continue
		}

		// 如果指定了卖家，需要进行过滤
		if len(seller) > 0 && transaction.Seller != seller {
			continue
		}

		// 如果指定了买家，需要进行过滤
		if len(buyer) > 0 && transaction.Buyer != buyer {
			continue
		}

		transactions = append(transactions, transaction)
	}

	return &QueryResult{
		Records:             transactions,
		RecordsCount:        int32(len(transactions)),
		Bookmark:            metadata.Bookmark,
		FetchedRecordsCount: metadata.FetchedRecordsCount,
	}, nil
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
