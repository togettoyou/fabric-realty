/**
 * @Author: 夜央 Oh oh oh oh oh oh (https://github.com/togettoyou)
 * @Email: zoujh99@qq.com
 * @Date: 2020/3/10 12:33 上午
 * @Description: 房地产相关合约路由
 */
package routers

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/togettoyou/blockchain-real-estate/chaincode/blockchain-real-estate/lib"
	"github.com/togettoyou/blockchain-real-estate/chaincode/blockchain-real-estate/utils"
	"strconv"
	"time"
)

//新建房地产(管理员)
func CreateRealEstate(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 验证参数
	if len(args) != 4 {
		return shim.Error("参数个数不满足")
	}
	accountId := args[0] //accountId用于验证是否为管理员
	proprietor := args[1]
	totalArea := args[2]
	livingSpace := args[3]
	if accountId == "" || proprietor == "" || totalArea == "" || livingSpace == "" {
		return shim.Error("参数存在空值")
	}
	if accountId == proprietor {
		return shim.Error("操作人应为管理员且与所有人不能相同")
	}
	// 参数数据格式转换
	var formattedTotalArea float64
	if val, err := strconv.ParseFloat(totalArea, 64); err != nil {
		return shim.Error(fmt.Sprintf("totalArea参数格式转换出错: %s", err))
	} else {
		formattedTotalArea = val
	}
	var formattedLivingSpace float64
	if val, err := strconv.ParseFloat(livingSpace, 64); err != nil {
		return shim.Error(fmt.Sprintf("livingSpace参数格式转换出错: %s", err))
	} else {
		formattedLivingSpace = val
	}
	//判断是否管理员操作
	resultsAccount, err := utils.GetStateByPartialCompositeKeys(stub, lib.AccountKey, []string{accountId})
	if err != nil || len(resultsAccount) != 1 {
		return shim.Error(fmt.Sprintf("操作人权限验证失败%s", err))
	}
	var account lib.Account
	if err = json.Unmarshal(resultsAccount[0], &account); err != nil {
		return shim.Error(fmt.Sprintf("查询操作人信息-反序列化出错: %s", err))
	}
	if account.UserName != "管理员" {
		return shim.Error(fmt.Sprintf("操作人权限不足%s", err))
	}
	//判断业主是否存在
	resultsProprietor, err := utils.GetStateByPartialCompositeKeys(stub, lib.AccountKey, []string{proprietor})
	if err != nil || len(resultsProprietor) != 1 {
		return shim.Error(fmt.Sprintf("业主proprietor信息验证失败%s", err))
	}
	realEstate := &lib.RealEstate{
		RealEstateID: fmt.Sprintf("%d", time.Now().Local().UnixNano()),
		Proprietor:   proprietor,
		Encumbrance:  false,
		TotalArea:    formattedTotalArea,
		LivingSpace:  formattedLivingSpace,
	}
	// 写入账本
	if err := utils.WriteLedger(realEstate, stub, lib.RealEstateKey, []string{realEstate.Proprietor, realEstate.RealEstateID}); err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	//将成功创建的信息返回
	realEstateByte, err := json.Marshal(realEstate)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化成功创建的信息出错: %s", err))
	}
	// 成功返回
	return shim.Success(realEstateByte)
}

//查询房地产(可查询所有，也可根据所有人查询名下房产)
func QueryRealEstateList(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var realEstateList []lib.RealEstate
	results, err := utils.GetStateByPartialCompositeKeys2(stub, lib.RealEstateKey, args)
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	for _, v := range results {
		if v != nil {
			var realEstate lib.RealEstate
			err := json.Unmarshal(v, &realEstate)
			if err != nil {
				return shim.Error(fmt.Sprintf("QueryRealEstateList-反序列化出错: %s", err))
			}
			realEstateList = append(realEstateList, realEstate)
		}
	}
	realEstateListByte, err := json.Marshal(realEstateList)
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryRealEstateList-序列化出错: %s", err))
	}
	return shim.Success(realEstateListByte)
}
