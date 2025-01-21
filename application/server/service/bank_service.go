package service

import (
	"application/pkg/fabric"
	"encoding/json"
	"fmt"
	"time"
)

type BankService struct{}

const BANK_ORG = "org2" // 银行组织

// CompleteTransaction 完成交易
func (s *BankService) CompleteTransaction(txID string) error {
	contract := fabric.GetContract(BANK_ORG)
	now := time.Now().Format(time.RFC3339)
	_, err := contract.SubmitTransaction("CompleteTransaction", txID, now)
	if err != nil {
		return fmt.Errorf("完成交易失败：%s", fabric.ExtractErrorMessage(err))
	}
	return nil
}

// QueryTransaction 查询交易信息
func (s *BankService) QueryTransaction(txID string) (map[string]interface{}, error) {
	contract := fabric.GetContract(BANK_ORG)
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
func (s *BankService) QueryTransactionList(pageSize int32, bookmark string, status string) (map[string]interface{}, error) {
	contract := fabric.GetContract(BANK_ORG)
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
func (s *BankService) QueryBlockList(pageSize int, pageNum int) (*fabric.BlockQueryResult, error) {
	result, err := fabric.GetBlockListener().GetBlocksByOrg(BANK_ORG, pageSize, pageNum)
	if err != nil {
		return nil, fmt.Errorf("查询区块列表失败：%v", err)
	}
	return result, nil
}
