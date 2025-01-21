package service

import (
	"application/pkg/fabric"
	"encoding/json"
	"fmt"
	"time"
)

type TradingPlatformService struct{}

const TRADE_ORG = "org3" // 交易平台组织

// CreateTransaction 生成交易
func (s *TradingPlatformService) CreateTransaction(txID, realEstateID, seller, buyer string, price float64) error {
	contract := fabric.GetContract(TRADE_ORG)
	now := time.Now().Format(time.RFC3339)
	_, err := contract.SubmitTransaction("CreateTransaction", txID, realEstateID, seller, buyer, fmt.Sprintf("%f", price), now)
	if err != nil {
		return fmt.Errorf("生成交易失败：%s", fabric.ExtractErrorMessage(err))
	}
	return nil
}

// QueryRealEstate 查询房产信息
func (s *TradingPlatformService) QueryRealEstate(id string) (map[string]interface{}, error) {
	contract := fabric.GetContract(TRADE_ORG)
	result, err := contract.EvaluateTransaction("QueryRealEstate", id)
	if err != nil {
		return nil, fmt.Errorf("查询房产信息失败：%s", fabric.ExtractErrorMessage(err))
	}

	var realEstate map[string]interface{}
	if err := json.Unmarshal(result, &realEstate); err != nil {
		return nil, fmt.Errorf("解析房产数据失败：%v", err)
	}

	return realEstate, nil
}

// QueryTransaction 查询交易信息
func (s *TradingPlatformService) QueryTransaction(txID string) (map[string]interface{}, error) {
	contract := fabric.GetContract(TRADE_ORG)
	result, err := contract.EvaluateTransaction("QueryTransaction", txID)
	if err != nil {
		return nil, fmt.Errorf("查询交易信息失败：%s", fabric.ExtractErrorMessage(err))
	}

	var transaction map[string]interface{}
	if err := json.Unmarshal(result, &transaction); err != nil {
		return nil, fmt.Errorf("解析交易数据失败：%v", err)
	}

	return transaction, nil
}

// QueryTransactionList 分页查询交易列表
func (s *TradingPlatformService) QueryTransactionList(pageSize int32, bookmark string, status string) (map[string]interface{}, error) {
	contract := fabric.GetContract(TRADE_ORG)
	result, err := contract.EvaluateTransaction("QueryTransactionList", fmt.Sprintf("%d", pageSize), bookmark, status)
	if err != nil {
		return nil, fmt.Errorf("查询交易列表失败：%s", fabric.ExtractErrorMessage(err))
	}

	var queryResult map[string]interface{}
	if err := json.Unmarshal(result, &queryResult); err != nil {
		return nil, fmt.Errorf("解析查询结果失败：%v", err)
	}

	return queryResult, nil
}

// QueryBlockList 分页查询区块列表
func (s *TradingPlatformService) QueryBlockList(pageSize int, pageNum int) (*fabric.BlockQueryResult, error) {
	result, err := fabric.GetBlockListener().GetBlocksByOrg(TRADE_ORG, pageSize, pageNum)
	if err != nil {
		return nil, fmt.Errorf("查询区块列表失败：%v", err)
	}
	return result, nil
}
