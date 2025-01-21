package service

import (
	"application/pkg/fabric"
	"encoding/json"
	"fmt"
	"time"
)

type RealtyAgencyService struct{}

const REALTY_ORG = "org1" // 不动产登记机构组织

// CreateRealEstate 创建房产信息
func (s *RealtyAgencyService) CreateRealEstate(id, address string, area float64, owner string) error {
	contract := fabric.GetContract(REALTY_ORG)
	now := time.Now().Format(time.RFC3339)
	_, err := contract.SubmitTransaction("CreateRealEstate", id, address, fmt.Sprintf("%f", area), owner, now)
	if err != nil {
		return fmt.Errorf("创建房产信息失败：%s", fabric.ExtractErrorMessage(err))
	}
	return nil
}

// QueryRealEstate 查询房产信息
func (s *RealtyAgencyService) QueryRealEstate(id string) (map[string]interface{}, error) {
	contract := fabric.GetContract(REALTY_ORG)
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

// QueryRealEstateList 分页查询房产列表
func (s *RealtyAgencyService) QueryRealEstateList(pageSize int32, bookmark string, status string) (map[string]interface{}, error) {
	contract := fabric.GetContract(REALTY_ORG)
	result, err := contract.EvaluateTransaction("QueryRealEstateList", fmt.Sprintf("%d", pageSize), bookmark, status)
	if err != nil {
		return nil, fmt.Errorf("查询房产列表失败：%s", fabric.ExtractErrorMessage(err))
	}

	var queryResult map[string]interface{}
	if err := json.Unmarshal(result, &queryResult); err != nil {
		return nil, fmt.Errorf("解析查询结果失败：%v", err)
	}

	return queryResult, nil
}

// QueryBlockList 分页查询区块列表
func (s *RealtyAgencyService) QueryBlockList(pageSize int, pageNum int) (*fabric.BlockQueryResult, error) {
	result, err := fabric.GetBlockListener().GetBlocksByOrg(REALTY_ORG, pageSize, pageNum)
	if err != nil {
		return nil, fmt.Errorf("查询区块列表失败：%v", err)
	}
	return result, nil
}
