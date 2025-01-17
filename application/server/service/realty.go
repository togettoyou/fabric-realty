package service

import (
	"application/utils"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/grpc/status"
)

type RealtyService struct{}

const (
	REALTY_ORG = "org1" // 房管局组织
	BANK_ORG   = "org2" // 银行组织
)

// extractErrorMessage 从错误中提取详细信息
func extractErrorMessage(err error) string {
	if err == nil {
		return ""
	}
	// 尝试获取 gRPC 状态
	if st, ok := status.FromError(err); ok {
		// 获取详细信息
		msg := st.Message()
		details := st.Details()
		code := st.Code()

		// 构建完整的错误信息
		fullError := fmt.Sprintf("错误码: %v, 消息: %v", code, msg)
		if len(details) > 0 {
			fullError += fmt.Sprintf(", 详情: %+v", details)
		}
		return fullError
	}
	return err.Error()
}

// CreateRealEstate 创建房产信息
func (s *RealtyService) CreateRealEstate(id, address string, area float64, owner string, price float64) error {
	// 使用房管局组织身份
	contract := utils.GetContract(REALTY_ORG)
	now := time.Now().Format(time.RFC3339)
	_, err := contract.SubmitTransaction("CreateRealEstate", id, address, fmt.Sprintf("%f", area), owner, fmt.Sprintf("%f", price), now)
	if err != nil {
		return fmt.Errorf("创建房产信息失败：%s", extractErrorMessage(err))
	}
	return nil
}

// QueryRealEstate 查询房产信息
func (s *RealtyService) QueryRealEstate(id string) (map[string]interface{}, error) {
	// 查询操作可以使用任意组织身份
	contract := utils.GetContract(REALTY_ORG)
	result, err := contract.EvaluateTransaction("QueryRealEstate", id)
	if err != nil {
		return nil, fmt.Errorf("查询房产信息失败：%s", extractErrorMessage(err))
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
	now := time.Now().Format(time.RFC3339)
	_, err := contract.SubmitTransaction("CreateTransaction", txID, realEstateID, seller, buyer, fmt.Sprintf("%f", price), now)
	if err != nil {
		return fmt.Errorf("创建交易失败：%s", extractErrorMessage(err))
	}
	return nil
}

// ConfirmEscrow 确认资金托管
func (s *RealtyService) ConfirmEscrow(txID string) error {
	// 使用银行组织身份
	contract := utils.GetContract(BANK_ORG)
	now := time.Now().Format(time.RFC3339)
	_, err := contract.SubmitTransaction("ConfirmEscrow", txID, now)
	if err != nil {
		return fmt.Errorf("确认资金托管失败：%s", extractErrorMessage(err))
	}
	return nil
}

// CompleteTransaction 完成交易
func (s *RealtyService) CompleteTransaction(txID string) error {
	// 使用银行组织身份
	contract := utils.GetContract(BANK_ORG)
	now := time.Now().Format(time.RFC3339)
	_, err := contract.SubmitTransaction("CompleteTransaction", txID, now)
	if err != nil {
		return fmt.Errorf("完成交易失败：%s", extractErrorMessage(err))
	}
	return nil
}
