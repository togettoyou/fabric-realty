package service

import (
	"application/utils"
	"encoding/json"
	"fmt"
)

type RealtyService struct{}

const (
	REALTY_ORG = "org1" // 房管局组织
	BANK_ORG   = "org2" // 银行组织
)

// CreateRealEstate 创建房产信息
func (s *RealtyService) CreateRealEstate(id, address string, area float64, owner string, price float64) error {
	// 使用房管局组织身份
	contract := utils.GetContract(REALTY_ORG)
	_, err := contract.SubmitTransaction("CreateRealEstate", id, address, fmt.Sprintf("%f", area), owner, fmt.Sprintf("%f", price))
	if err != nil {
		return fmt.Errorf("创建房产信息失败：%v", err)
	}
	return nil
}

// QueryRealEstate 查询房产信息
func (s *RealtyService) QueryRealEstate(id string) (map[string]interface{}, error) {
	// 查询操作可以使用任意组织身份
	contract := utils.GetContract(REALTY_ORG)
	result, err := contract.EvaluateTransaction("QueryRealEstate", id)
	if err != nil {
		return nil, fmt.Errorf("查询房产信息失败：%v", err)
	}

	var realEstate map[string]interface{}
	if err := json.Unmarshal(result, &realEstate); err != nil {
		return nil, fmt.Errorf("解析房产数据失败：%v", err)
	}

	return realEstate, nil
}

// CreateTransaction 创建交易
func (s *RealtyService) CreateTransaction(txID, realEstateID, seller, buyer string, price float64) error {
	// 创建交易可以使用任意组织身份
	contract := utils.GetContract(REALTY_ORG)
	_, err := contract.SubmitTransaction("CreateTransaction", txID, realEstateID, seller, buyer, fmt.Sprintf("%f", price))
	if err != nil {
		return fmt.Errorf("创建交易失败：%v", err)
	}
	return nil
}

// ConfirmEscrow 确认资金托管
func (s *RealtyService) ConfirmEscrow(txID string) error {
	// 使用银行组织身份
	contract := utils.GetContract(BANK_ORG)
	_, err := contract.SubmitTransaction("ConfirmEscrow", txID)
	if err != nil {
		return fmt.Errorf("确认资金托管失败：%v", err)
	}
	return nil
}

// CompleteTransaction 完成交易
func (s *RealtyService) CompleteTransaction(txID string) error {
	// 使用银行组织身份
	contract := utils.GetContract(BANK_ORG)
	_, err := contract.SubmitTransaction("CompleteTransaction", txID)
	if err != nil {
		return fmt.Errorf("完成交易失败：%v", err)
	}
	return nil
}
