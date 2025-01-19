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

// RealEstateStatus 房产状态
type RealEstateStatus string

const (
	NORMAL         RealEstateStatus = "NORMAL"         // 正常
	IN_TRANSACTION RealEstateStatus = "IN_TRANSACTION" // 交易中
)

// TransactionStatus 交易状态
type TransactionStatus string

const (
	PENDING   TransactionStatus = "PENDING"   // 待付款
	COMPLETED TransactionStatus = "COMPLETED" // 已完成
)

// RealEstate 房产信息
type RealEstate struct {
	ID              string           `json:"id"`              // 房产ID
	PropertyAddress string           `json:"propertyAddress"` // 房产地址
	Area            float64          `json:"area"`            // 面积
	CurrentOwner    string           `json:"currentOwner"`    // 当前所有者
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

// QueryResult 分页查询结果
type QueryResult struct {
	Records             []interface{} `json:"records"`             // 记录列表
	RecordsCount        int32         `json:"recordsCount"`        // 本次返回的记录数
	Bookmark            string        `json:"bookmark"`            // 书签，用于下一页查询
	FetchedRecordsCount int32         `json:"fetchedRecordsCount"` // 总共获取的记录数
}

// 组织 MSP ID 常量
const (
	REALTY_ORG_MSPID = "Org1MSP" // 不动产登记机构组织 MSP ID
	BANK_ORG_MSPID   = "Org2MSP" // 银行组织 MSP ID
	TRADE_ORG_MSPID  = "Org3MSP" // 交易平台组织 MSP ID
)

// 通用方法: 获取客户端身份信息
func (s *SmartContract) getClientIdentityMSPID(ctx contractapi.TransactionContextInterface) (string, error) {
	clientID, err := cid.New(ctx.GetStub())
	if err != nil {
		return "", fmt.Errorf("获取客户端身份信息失败：%v", err)
	}
	return clientID.GetMSPID()
}

// 通用方法：创建和获取复合键
func (s *SmartContract) getCompositeKey(ctx contractapi.TransactionContextInterface, objectType string, attributes []string) (string, error) {
	key, err := ctx.GetStub().CreateCompositeKey(objectType, attributes)
	if err != nil {
		return "", fmt.Errorf("创建复合键失败：%v", err)
	}
	return key, nil
}

// 通用方法：获取状态
func (s *SmartContract) getState(ctx contractapi.TransactionContextInterface, key string, value interface{}) error {
	bytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		return fmt.Errorf("读取状态失败：%v", err)
	}
	if bytes == nil {
		return fmt.Errorf("键 %s 不存在", key)
	}

	err = json.Unmarshal(bytes, value)
	if err != nil {
		return fmt.Errorf("解析数据失败：%v", err)
	}
	return nil
}

// 通用方法：保存状态
func (s *SmartContract) putState(ctx contractapi.TransactionContextInterface, key string, value interface{}) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("序列化数据失败：%v", err)
	}

	err = ctx.GetStub().PutState(key, bytes)
	if err != nil {
		return fmt.Errorf("保存状态失败：%v", err)
	}
	return nil
}

// CreateRealEstate 创建房产信息（仅不动产登记机构组织可以调用）
func (s *SmartContract) CreateRealEstate(ctx contractapi.TransactionContextInterface, id string, address string, area float64, owner string, createTime time.Time) error {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return fmt.Errorf("获取调用者身份失败：%v", err)
	}

	// 验证是否是不动产登记机构组织的成员
	if clientMSPID != REALTY_ORG_MSPID {
		return fmt.Errorf("只有不动产登记机构组织成员才能创建房产信息")
	}

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

	// 检查房产是否已存在（检查所有可能的状态）
	for _, status := range []RealEstateStatus{NORMAL, IN_TRANSACTION} {
		key, err := s.getCompositeKey(ctx, REAL_ESTATE, []string{string(status), id})
		if err != nil {
			return fmt.Errorf("创建复合键失败：%v", err)
		}

		exists, err := ctx.GetStub().GetState(key)
		if err != nil {
			return fmt.Errorf("查询房产信息失败：%v", err)
		}
		if exists != nil {
			return fmt.Errorf("房产ID %s 已存在", id)
		}
	}

	// 创建房产信息
	realEstate := RealEstate{
		ID:              id,
		PropertyAddress: address,
		Area:            area,
		CurrentOwner:    owner,
		Status:          NORMAL,
		CreateTime:      createTime,
		UpdateTime:      createTime,
	}

	// 保存房产信息（复合键：类型_状态_ID）
	key, err := s.getCompositeKey(ctx, REAL_ESTATE, []string{string(NORMAL), id})
	if err != nil {
		return err
	}

	err = s.putState(ctx, key, realEstate)
	if err != nil {
		return err
	}

	return nil
}

// CreateTransaction 生成交易（仅交易平台组织可以调用）
func (s *SmartContract) CreateTransaction(ctx contractapi.TransactionContextInterface, txID string, realEstateID string, seller string, buyer string, price float64, createTime time.Time) error {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return fmt.Errorf("获取调用者身份失败：%v", err)
	}

	// 验证是否是交易平台组织的成员
	if clientMSPID != TRADE_ORG_MSPID {
		return fmt.Errorf("只有交易平台组织成员才能生成交易")
	}

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

	// 查询房产信息
	realEstateKey, err := s.getCompositeKey(ctx, REAL_ESTATE, []string{string(NORMAL), realEstateID})
	if err != nil {
		return err
	}

	var realEstate RealEstate
	err = s.getState(ctx, realEstateKey, &realEstate)
	if err != nil {
		return err
	}

	// 检查卖家是否是房产所有者
	if realEstate.CurrentOwner != seller {
		return fmt.Errorf("卖家不是房产所有者")
	}

	// 生成交易信息
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

	// 更新房产状态
	realEstate.Status = IN_TRANSACTION
	realEstate.UpdateTime = createTime

	// 保存状态
	txKey, err := s.getCompositeKey(ctx, TRANSACTION, []string{string(PENDING), txID})
	if err != nil {
		return err
	}

	// 删除旧的房产记录
	err = ctx.GetStub().DelState(realEstateKey)
	if err != nil {
		return fmt.Errorf("删除旧的房产记录失败：%v", err)
	}

	// 创建新的房产记录（使用新状态）
	newRealEstateKey, err := s.getCompositeKey(ctx, REAL_ESTATE, []string{string(IN_TRANSACTION), realEstateID})
	if err != nil {
		return err
	}

	err = s.putState(ctx, txKey, transaction)
	if err != nil {
		return err
	}

	err = s.putState(ctx, newRealEstateKey, realEstate)
	if err != nil {
		return err
	}

	return nil
}

// CompleteTransaction 完成交易（仅银行组织可以调用）
func (s *SmartContract) CompleteTransaction(ctx contractapi.TransactionContextInterface, txID string, updateTime time.Time) error {
	// 检查调用者身份
	clientMSPID, err := s.getClientIdentityMSPID(ctx)
	if err != nil {
		return fmt.Errorf("获取调用者身份失败：%v", err)
	}

	// 验证是否是银行组织的成员
	if clientMSPID != BANK_ORG_MSPID {
		return fmt.Errorf("只有银行组织成员才能完成交易")
	}

	// 查询交易信息
	txKey, err := s.getCompositeKey(ctx, TRANSACTION, []string{string(PENDING), txID})
	if err != nil {
		return err
	}

	var transaction Transaction
	err = s.getState(ctx, txKey, &transaction)
	if err != nil {
		return err
	}

	// 查询房产信息
	realEstateKey, err := s.getCompositeKey(ctx, REAL_ESTATE, []string{string(IN_TRANSACTION), transaction.RealEstateID})
	if err != nil {
		return err
	}

	var realEstate RealEstate
	err = s.getState(ctx, realEstateKey, &realEstate)
	if err != nil {
		return err
	}

	// 更新状态
	realEstate.CurrentOwner = transaction.Buyer
	realEstate.Status = NORMAL
	realEstate.UpdateTime = updateTime

	transaction.Status = COMPLETED
	transaction.UpdateTime = updateTime

	// 删除旧记录
	err = ctx.GetStub().DelState(txKey)
	if err != nil {
		return fmt.Errorf("删除旧的交易记录失败：%v", err)
	}

	err = ctx.GetStub().DelState(realEstateKey)
	if err != nil {
		return fmt.Errorf("删除旧的房产记录失败：%v", err)
	}

	// 创建新记录
	newTxKey, err := s.getCompositeKey(ctx, TRANSACTION, []string{string(COMPLETED), txID})
	if err != nil {
		return err
	}

	newRealEstateKey, err := s.getCompositeKey(ctx, REAL_ESTATE, []string{string(NORMAL), transaction.RealEstateID})
	if err != nil {
		return err
	}

	err = s.putState(ctx, newTxKey, transaction)
	if err != nil {
		return err
	}

	err = s.putState(ctx, newRealEstateKey, realEstate)
	if err != nil {
		return err
	}

	return nil
}

// QueryRealEstate 查询房产信息
func (s *SmartContract) QueryRealEstate(ctx contractapi.TransactionContextInterface, id string) (*RealEstate, error) {
	// 遍历所有可能的状态查询房产
	for _, status := range []RealEstateStatus{NORMAL, IN_TRANSACTION} {
		key, err := s.getCompositeKey(ctx, REAL_ESTATE, []string{string(status), id})
		if err != nil {
			return nil, fmt.Errorf("创建复合键失败：%v", err)
		}

		bytes, err := ctx.GetStub().GetState(key)
		if err != nil {
			return nil, fmt.Errorf("查询房产信息失败：%v", err)
		}
		if bytes != nil {
			var realEstate RealEstate
			err = json.Unmarshal(bytes, &realEstate)
			if err != nil {
				return nil, fmt.Errorf("解析房产信息失败：%v", err)
			}
			return &realEstate, nil
		}
	}

	return nil, fmt.Errorf("房产ID %s 不存在", id)
}

// QueryTransaction 查询交易信息
func (s *SmartContract) QueryTransaction(ctx contractapi.TransactionContextInterface, txID string) (*Transaction, error) {
	// 遍历所有可能的状态查询交易
	for _, status := range []TransactionStatus{PENDING, COMPLETED} {
		key, err := s.getCompositeKey(ctx, TRANSACTION, []string{string(status), txID})
		if err != nil {
			return nil, fmt.Errorf("创建复合键失败：%v", err)
		}

		bytes, err := ctx.GetStub().GetState(key)
		if err != nil {
			return nil, fmt.Errorf("查询交易信息失败：%v", err)
		}
		if bytes != nil {
			var transaction Transaction
			err = json.Unmarshal(bytes, &transaction)
			if err != nil {
				return nil, fmt.Errorf("解析交易信息失败：%v", err)
			}
			return &transaction, nil
		}
	}

	return nil, fmt.Errorf("交易ID %s 不存在", txID)
}

// QueryRealEstateList 分页查询房产列表
func (s *SmartContract) QueryRealEstateList(ctx contractapi.TransactionContextInterface, pageSize int32, bookmark string, status string) (*QueryResult, error) {
	var iterator shim.StateQueryIteratorInterface
	var metadata *peer.QueryResponseMetadata
	var err error

	if status != "" {
		iterator, metadata, err = ctx.GetStub().GetStateByPartialCompositeKeyWithPagination(
			REAL_ESTATE,
			[]string{status},
			pageSize,
			bookmark,
		)
	} else {
		iterator, metadata, err = ctx.GetStub().GetStateByPartialCompositeKeyWithPagination(
			REAL_ESTATE,
			[]string{},
			pageSize,
			bookmark,
		)
	}

	if err != nil {
		return nil, fmt.Errorf("查询列表失败：%v", err)
	}
	defer iterator.Close()

	records := make([]interface{}, 0)
	for iterator.HasNext() {
		queryResponse, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("获取下一条记录失败：%v", err)
		}

		var realEstate RealEstate
		err = json.Unmarshal(queryResponse.Value, &realEstate)
		if err != nil {
			return nil, fmt.Errorf("解析房产信息失败：%v", err)
		}

		records = append(records, realEstate)
	}

	return &QueryResult{
		Records:             records,
		RecordsCount:        int32(len(records)),
		Bookmark:            metadata.Bookmark,
		FetchedRecordsCount: metadata.FetchedRecordsCount,
	}, nil
}

// QueryTransactionList 分页查询交易列表
func (s *SmartContract) QueryTransactionList(ctx contractapi.TransactionContextInterface, pageSize int32, bookmark string, status string) (*QueryResult, error) {
	var iterator shim.StateQueryIteratorInterface
	var metadata *peer.QueryResponseMetadata
	var err error

	if status != "" {
		iterator, metadata, err = ctx.GetStub().GetStateByPartialCompositeKeyWithPagination(
			TRANSACTION,
			[]string{status},
			pageSize,
			bookmark,
		)
	} else {
		iterator, metadata, err = ctx.GetStub().GetStateByPartialCompositeKeyWithPagination(
			TRANSACTION,
			[]string{},
			pageSize,
			bookmark,
		)
	}

	if err != nil {
		return nil, fmt.Errorf("查询列表失败：%v", err)
	}
	defer iterator.Close()

	records := make([]interface{}, 0)
	for iterator.HasNext() {
		queryResponse, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("获取下一条记录失败：%v", err)
		}

		var transaction Transaction
		err = json.Unmarshal(queryResponse.Value, &transaction)
		if err != nil {
			return nil, fmt.Errorf("解析交易信息失败：%v", err)
		}

		records = append(records, transaction)
	}

	return &QueryResult{
		Records:             records,
		RecordsCount:        int32(len(records)),
		Bookmark:            metadata.Bookmark,
		FetchedRecordsCount: metadata.FetchedRecordsCount,
	}, nil
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
