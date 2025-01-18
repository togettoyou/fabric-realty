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

	// 索引前缀
	IDX_RE_STATUS_TIME = "RE_STATUS_TIME" // 房产状态时间索引
	IDX_TX_STATUS_TIME = "TX_STATUS_TIME" // 交易状态时间索引
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

// 组织 MSP ID 常量
const (
	REALTY_ORG_MSPID = "Org1MSP" // 不动产登记机构组织 MSP ID
	BANK_ORG_MSPID   = "Org2MSP" // 银行组织 MSP ID
	TRADE_ORG_MSPID  = "Org3MSP" // 交易平台组织 MSP ID
)

// GetClientIdentityMSPID 获取客户端身份信息
func (s *SmartContract) GetClientIdentityMSPID(ctx contractapi.TransactionContextInterface) (string, error) {
	clientID, err := cid.New(ctx.GetStub())
	if err != nil {
		return "", fmt.Errorf("获取客户端身份信息失败：%v", err)
	}
	return clientID.GetMSPID()
}

// CreateRealEstate 创建房产信息（仅不动产登记机构组织可以调用）
func (s *SmartContract) CreateRealEstate(ctx contractapi.TransactionContextInterface, id string, address string, area float64, owner string, createTime time.Time) error {
	// 检查调用者身份
	clientMSPID, err := s.GetClientIdentityMSPID(ctx)
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

	// 检查房产是否已存在
	key, err := ctx.GetStub().CreateCompositeKey(REAL_ESTATE, []string{id})
	if err != nil {
		return fmt.Errorf("创建复合键失败：%v", err)
	}

	realEstateBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		return fmt.Errorf("查询房产信息失败：%v", err)
	}
	if realEstateBytes != nil {
		return fmt.Errorf("房产ID %s 已存在", id)
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

	realEstateJSON, err := json.Marshal(realEstate)
	if err != nil {
		return fmt.Errorf("序列化房产信息失败：%v", err)
	}

	// 保存房产信息
	err = ctx.GetStub().PutState(key, realEstateJSON)
	if err != nil {
		return fmt.Errorf("保存房产信息失败：%v", err)
	}

	// 创建状态时间索引
	timeStr := createTime.Format("2006-01-02 15:04:05.000")
	statusTimeKey, err := ctx.GetStub().CreateCompositeKey(IDX_RE_STATUS_TIME, []string{string(realEstate.Status), timeStr, id})
	if err != nil {
		return fmt.Errorf("创建状态时间索引失败：%v", err)
	}

	// 保存空值，我们只需要键
	err = ctx.GetStub().PutState(statusTimeKey, []byte{0x00})
	if err != nil {
		return fmt.Errorf("保存状态时间索引失败：%v", err)
	}

	return nil
}

// CreateTransaction 创建交易（仅交易平台组织可以调用）
func (s *SmartContract) CreateTransaction(ctx contractapi.TransactionContextInterface, txID string, realEstateID string, seller string, buyer string, price float64, createTime time.Time) error {
	// 检查调用者身份
	clientMSPID, err := s.GetClientIdentityMSPID(ctx)
	if err != nil {
		return fmt.Errorf("获取调用者身份失败：%v", err)
	}

	// 验证是否是交易平台组织的成员
	if clientMSPID != TRADE_ORG_MSPID {
		return fmt.Errorf("只有交易平台组织成员才能创建交易")
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
	realEstateKey, err := ctx.GetStub().CreateCompositeKey(REAL_ESTATE, []string{realEstateID})
	if err != nil {
		return fmt.Errorf("创建房产复合键失败：%v", err)
	}

	realEstateBytes, err := ctx.GetStub().GetState(realEstateKey)
	if err != nil {
		return fmt.Errorf("查询房产信息失败：%v", err)
	}
	if realEstateBytes == nil {
		return fmt.Errorf("房产ID %s 不存在", realEstateID)
	}

	var realEstate RealEstate
	err = json.Unmarshal(realEstateBytes, &realEstate)
	if err != nil {
		return fmt.Errorf("解析房产信息失败：%v", err)
	}

	// 检查房产状态
	if realEstate.Status != NORMAL {
		return fmt.Errorf("房产状态不正确，当前状态：%s，需要状态：%s", realEstate.Status, NORMAL)
	}

	// 检查卖家是否是房产所有者
	if realEstate.CurrentOwner != seller {
		return fmt.Errorf("卖家不是房产所有者")
	}

	// 创建交易信息
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

	// 更新房产状态
	realEstate.Status = IN_TRANSACTION
	realEstate.UpdateTime = createTime
	realEstateJSON, err := json.Marshal(realEstate)
	if err != nil {
		return fmt.Errorf("序列化房产信息失败：%v", err)
	}

	// 保存状态
	err = ctx.GetStub().PutState(realEstateKey, realEstateJSON)
	if err != nil {
		return fmt.Errorf("更新房产信息失败：%v", err)
	}

	err = ctx.GetStub().PutState(txKey, transactionJSON)
	if err != nil {
		return fmt.Errorf("保存交易信息失败：%v", err)
	}

	// 创建状态时间索引
	timeStr := transaction.CreateTime.Format("2006-01-02 15:04:05.000")
	statusTimeKey, err := ctx.GetStub().CreateCompositeKey(IDX_TX_STATUS_TIME, []string{string(transaction.Status), timeStr, txID})
	if err != nil {
		return fmt.Errorf("创建状态时间索引失败：%v", err)
	}

	// 保存空值，我们只需要键
	err = ctx.GetStub().PutState(statusTimeKey, []byte{0x00})
	if err != nil {
		return fmt.Errorf("保存状态时间索引失败：%v", err)
	}

	return nil
}

// CompleteTransaction 完成交易（仅银行组织可以调用）
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

	// 查询交易信息
	txKey, err := ctx.GetStub().CreateCompositeKey(TRANSACTION, []string{txID})
	if err != nil {
		return fmt.Errorf("创建交易复合键失败：%v", err)
	}

	transactionBytes, err := ctx.GetStub().GetState(txKey)
	if err != nil {
		return fmt.Errorf("查询交易信息失败：%v", err)
	}
	if transactionBytes == nil {
		return fmt.Errorf("交易ID %s 不存在", txID)
	}

	var transaction Transaction
	err = json.Unmarshal(transactionBytes, &transaction)
	if err != nil {
		return fmt.Errorf("解析交易信息失败：%v", err)
	}

	// 检查交易状态
	if transaction.Status != PENDING {
		return fmt.Errorf("交易状态不正确，当前状态：%s，需要状态：%s", transaction.Status, PENDING)
	}

	// 查询房产信息
	realEstateKey, err := ctx.GetStub().CreateCompositeKey(REAL_ESTATE, []string{transaction.RealEstateID})
	if err != nil {
		return fmt.Errorf("创建房产复合键失败：%v", err)
	}

	realEstateBytes, err := ctx.GetStub().GetState(realEstateKey)
	if err != nil {
		return fmt.Errorf("查询房产信息失败：%v", err)
	}
	if realEstateBytes == nil {
		return fmt.Errorf("房产ID %s 不存在", transaction.RealEstateID)
	}

	var realEstate RealEstate
	err = json.Unmarshal(realEstateBytes, &realEstate)
	if err != nil {
		return fmt.Errorf("解析房产信息失败：%v", err)
	}

	// 更新房产信息
	realEstate.CurrentOwner = transaction.Buyer
	realEstate.Status = NORMAL
	realEstate.UpdateTime = updateTime

	realEstateJSON, err := json.Marshal(realEstate)
	if err != nil {
		return fmt.Errorf("序列化房产信息失败：%v", err)
	}

	// 更新交易信息
	transaction.Status = COMPLETED
	transaction.UpdateTime = updateTime

	transactionJSON, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("序列化交易信息失败：%v", err)
	}

	// 保存状态
	err = ctx.GetStub().PutState(realEstateKey, realEstateJSON)
	if err != nil {
		return fmt.Errorf("更新房产信息失败：%v", err)
	}

	err = ctx.GetStub().PutState(txKey, transactionJSON)
	if err != nil {
		return fmt.Errorf("更新交易信息失败：%v", err)
	}

	// 删除旧的状态时间索引
	oldTimeStr := transaction.CreateTime.Format("2006-01-02 15:04:05.000")
	oldStatusTimeKey, err := ctx.GetStub().CreateCompositeKey(IDX_TX_STATUS_TIME, []string{string(PENDING), oldTimeStr, txID})
	if err != nil {
		return fmt.Errorf("创建旧状态时间索引失败：%v", err)
	}
	err = ctx.GetStub().DelState(oldStatusTimeKey)
	if err != nil {
		return fmt.Errorf("删除旧状态时间索引失败：%v", err)
	}

	// 创建新的状态时间索引（使用原来的创建时间）
	newStatusTimeKey, err := ctx.GetStub().CreateCompositeKey(IDX_TX_STATUS_TIME, []string{string(COMPLETED), oldTimeStr, txID})
	if err != nil {
		return fmt.Errorf("创建新状态时间索引失败：%v", err)
	}
	err = ctx.GetStub().PutState(newStatusTimeKey, []byte{0x00})
	if err != nil {
		return fmt.Errorf("保存新状态时间索引失败：%v", err)
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

// QueryResult 分页查询结果
type QueryResult struct {
	Records             []interface{} `json:"records"`             // 记录列表
	RecordsCount        int32         `json:"recordsCount"`        // 本次返回的记录数
	Bookmark            string        `json:"bookmark"`            // 书签，用于下一页查询
	FetchedRecordsCount int32         `json:"fetchedRecordsCount"` // 总共获取的记录数
}

// QueryRealEstateList 分页查询房产列表
func (s *SmartContract) QueryRealEstateList(ctx contractapi.TransactionContextInterface, pageSize int32, bookmark string, status string) (*QueryResult, error) {
	// 使用状态时间索引查询
	var iterator shim.StateQueryIteratorInterface
	var metadata *peer.QueryResponseMetadata
	var err error

	if status != "" {
		// 如果指定了状态，使用状态前缀查询
		iterator, metadata, err = ctx.GetStub().GetStateByPartialCompositeKeyWithPagination(
			IDX_RE_STATUS_TIME,
			[]string{status},
			pageSize,
			bookmark,
		)
	} else {
		// 如果没有指定状态，查询所有记录
		iterator, metadata, err = ctx.GetStub().GetStateByPartialCompositeKeyWithPagination(
			IDX_RE_STATUS_TIME,
			[]string{},
			pageSize,
			bookmark,
		)
	}

	if err != nil {
		return nil, fmt.Errorf("查询房产列表失败：%v", err)
	}
	defer iterator.Close()

	realEstates := make([]interface{}, 0)
	for iterator.HasNext() {
		queryResponse, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("获取下一条记录失败：%v", err)
		}

		// 从复合键中解析出房产ID
		_, compositeKeyParts, err := ctx.GetStub().SplitCompositeKey(queryResponse.Key)
		if err != nil {
			return nil, fmt.Errorf("解析复合键失败：%v", err)
		}

		// 获取房产ID（复合键的最后一部分）
		realEstateID := compositeKeyParts[len(compositeKeyParts)-1]

		// 查询房产详细信息
		realEstateKey, err := ctx.GetStub().CreateCompositeKey(REAL_ESTATE, []string{realEstateID})
		if err != nil {
			return nil, fmt.Errorf("创建房产键失败：%v", err)
		}

		realEstateBytes, err := ctx.GetStub().GetState(realEstateKey)
		if err != nil {
			return nil, fmt.Errorf("查询房产信息失败：%v", err)
		}

		var realEstate RealEstate
		err = json.Unmarshal(realEstateBytes, &realEstate)
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
func (s *SmartContract) QueryTransactionList(ctx contractapi.TransactionContextInterface, pageSize int32, bookmark string, status string) (*QueryResult, error) {
	// 使用状态时间索引查询
	var iterator shim.StateQueryIteratorInterface
	var metadata *peer.QueryResponseMetadata
	var err error

	if status != "" {
		// 如果指定了状态，使用状态前缀查询
		iterator, metadata, err = ctx.GetStub().GetStateByPartialCompositeKeyWithPagination(
			IDX_TX_STATUS_TIME,
			[]string{status},
			pageSize,
			bookmark,
		)
	} else {
		// 如果没有指定状态，查询所有记录
		iterator, metadata, err = ctx.GetStub().GetStateByPartialCompositeKeyWithPagination(
			IDX_TX_STATUS_TIME,
			[]string{},
			pageSize,
			bookmark,
		)
	}

	if err != nil {
		return nil, fmt.Errorf("查询交易列表失败：%v", err)
	}
	defer iterator.Close()

	transactions := make([]interface{}, 0)
	for iterator.HasNext() {
		queryResponse, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("获取下一条记录失败：%v", err)
		}

		// 从复合键中解析出交易ID
		_, compositeKeyParts, err := ctx.GetStub().SplitCompositeKey(queryResponse.Key)
		if err != nil {
			return nil, fmt.Errorf("解析复合键失败：%v", err)
		}

		// 获取交易ID（复合键的最后一部分）
		txID := compositeKeyParts[len(compositeKeyParts)-1]

		// 查询交易详细信息
		txKey, err := ctx.GetStub().CreateCompositeKey(TRANSACTION, []string{txID})
		if err != nil {
			return nil, fmt.Errorf("创建交易键失败：%v", err)
		}

		txBytes, err := ctx.GetStub().GetState(txKey)
		if err != nil {
			return nil, fmt.Errorf("查询交易信息失败：%v", err)
		}

		var transaction Transaction
		err = json.Unmarshal(txBytes, &transaction)
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
