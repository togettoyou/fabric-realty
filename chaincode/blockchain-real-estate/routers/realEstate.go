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
	snowflake "github.com/togettoyou/blockchain-real-estate/chaincode/blockchain-real-estate/pkg"
	"github.com/togettoyou/blockchain-real-estate/chaincode/blockchain-real-estate/utils"
	"strconv"
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
	results, err := utils.GetStateByPartialCompositeKeys(stub, lib.AccountKey, []string{accountId})
	if err != nil || len(results) != 1 {
		return shim.Error(fmt.Sprintf("操作人权限验证失败%s", err))
	}
	var account lib.Account
	if err = json.Unmarshal(results[0], &account); err != nil {
		return shim.Error(fmt.Sprintf("查询操作人信息-反序列化出错: %s", err))
	}
	if account.UserName != "管理员" {
		return shim.Error(fmt.Sprintf("操作人权限不足%s", err))
	}
	//生成唯一ID
	node, err := snowflake.NewNode(1)
	if err != nil {
		fmt.Println(err)
		return shim.Error(fmt.Sprintf("生成唯一ID出错: %s", err))
	}
	id := node.Generate()
	realEstate := &lib.RealEstate{
		RealEstateID: id.Base64(),
		Proprietor:   proprietor,
		Encumbrance:  false,
		TotalArea:    formattedTotalArea,
		LivingSpace:  formattedLivingSpace,
	}
	// 写入账本
	if err := utils.WriteLedger(realEstate, stub, lib.RealEstateKey, []string{realEstate.RealEstateID, realEstate.Proprietor}); err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	// 成功返回
	return shim.Success(nil)
}
