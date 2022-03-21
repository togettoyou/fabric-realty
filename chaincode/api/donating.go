package api

import (
	"chaincode/model"
	"chaincode/pkg/utils"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// CreateDonating 发起捐赠
func CreateDonating(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 验证参数
	if len(args) != 3 {
		return shim.Error("参数个数不满足")
	}
	objectOfDonating := args[0]
	donor := args[1]
	grantee := args[2]
	if objectOfDonating == "" || donor == "" || grantee == "" {
		return shim.Error("参数存在空值")
	}
	if donor == grantee {
		return shim.Error("捐赠人和受赠人不能同一人")
	}
	//判断objectOfDonating是否属于donor
	resultsRealEstate, err := utils.GetStateByPartialCompositeKeys2(stub, model.RealEstateKey, []string{donor, objectOfDonating})
	if err != nil || len(resultsRealEstate) != 1 {
		return shim.Error(fmt.Sprintf("验证%s属于%s失败: %s", objectOfDonating, donor, err))
	}
	var realEstate model.RealEstate
	if err = json.Unmarshal(resultsRealEstate[0], &realEstate); err != nil {
		return shim.Error(fmt.Sprintf("CreateDonating-反序列化出错: %s", err))
	}
	//根据grantee获取受赠人信息
	resultsAccount, err := utils.GetStateByPartialCompositeKeys(stub, model.AccountKey, []string{grantee})
	if err != nil || len(resultsAccount) != 1 {
		return shim.Error(fmt.Sprintf("grantee受赠人信息验证失败%s", err))
	}
	var accountGrantee model.Account
	if err = json.Unmarshal(resultsAccount[0], &accountGrantee); err != nil {
		return shim.Error(fmt.Sprintf("查询操作人信息-反序列化出错: %s", err))
	}
	if accountGrantee.UserName == "管理员" {
		return shim.Error(fmt.Sprintf("不能捐赠给管理员%s", err))
	}
	//判断记录是否已存在，不能重复发起捐赠
	//若Encumbrance为true即说明此房产已经正在担保状态
	if realEstate.Encumbrance {
		return shim.Error("此房地产已经作为担保状态，不能再发起捐赠")
	}
	createTime, _ := stub.GetTxTimestamp()
	donating := &model.Donating{
		ObjectOfDonating: objectOfDonating,
		Donor:            donor,
		Grantee:          grantee,
		CreateTime:       time.Unix(int64(createTime.GetSeconds()), int64(createTime.GetNanos())).Local().Format("2006-01-02 15:04:05"),
		DonatingStatus:   model.DonatingStatusConstant()["donatingStart"],
	}
	// 写入账本
	if err := utils.WriteLedger(donating, stub, model.DonatingKey, []string{donating.Donor, donating.ObjectOfDonating, donating.Grantee}); err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	//将房子状态设置为正在担保状态
	realEstate.Encumbrance = true
	if err := utils.WriteLedger(realEstate, stub, model.RealEstateKey, []string{realEstate.Proprietor, realEstate.RealEstateID}); err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	//将本次购买交易写入账本,可供受赠人查询
	donatingGrantee := &model.DonatingGrantee{
		Grantee:    grantee,
		CreateTime: time.Unix(int64(createTime.GetSeconds()), int64(createTime.GetNanos())).Local().Format("2006-01-02 15:04:05"),
		Donating:   *donating,
	}
	if err := utils.WriteLedger(donatingGrantee, stub, model.DonatingGranteeKey, []string{donatingGrantee.Grantee, donatingGrantee.CreateTime}); err != nil {
		return shim.Error(fmt.Sprintf("将本次捐赠交易写入账本失败%s", err))
	}
	donatingGranteeByte, err := json.Marshal(donatingGrantee)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化成功创建的信息出错: %s", err))
	}
	// 成功返回
	return shim.Success(donatingGranteeByte)
}

// QueryDonatingList 查询捐赠列表(可查询所有，也可根据发起捐赠人查询)(发起的)(供捐赠人查询)
func QueryDonatingList(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var donatingList []model.Donating
	results, err := utils.GetStateByPartialCompositeKeys2(stub, model.DonatingKey, args)
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	for _, v := range results {
		if v != nil {
			var donating model.Donating
			err := json.Unmarshal(v, &donating)
			if err != nil {
				return shim.Error(fmt.Sprintf("QueryDonatingList-反序列化出错: %s", err))
			}
			donatingList = append(donatingList, donating)
		}
	}
	donatingListByte, err := json.Marshal(donatingList)
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryDonatingList-序列化出错: %s", err))
	}
	return shim.Success(donatingListByte)
}

// QueryDonatingListByGrantee 根据受赠人(受赠人AccountId)查询捐赠(受赠的)(供受赠人查询)
func QueryDonatingListByGrantee(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error(fmt.Sprintf("必须指定受赠人AccountId查询"))
	}
	var donatingGranteeList []model.DonatingGrantee
	results, err := utils.GetStateByPartialCompositeKeys2(stub, model.DonatingGranteeKey, args)
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	for _, v := range results {
		if v != nil {
			var donatingGrantee model.DonatingGrantee
			err := json.Unmarshal(v, &donatingGrantee)
			if err != nil {
				return shim.Error(fmt.Sprintf("QueryDonatingListByGrantee-反序列化出错: %s", err))
			}
			donatingGranteeList = append(donatingGranteeList, donatingGrantee)
		}
	}
	donatingGranteeListByte, err := json.Marshal(donatingGranteeList)
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryDonatingListByGrantee-序列化出错: %s", err))
	}
	return shim.Success(donatingGranteeListByte)
}

// UpdateDonating 更新捐赠状态（确认受赠、取消）
func UpdateDonating(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 验证参数
	if len(args) != 4 {
		return shim.Error("参数个数不满足")
	}
	objectOfDonating := args[0]
	donor := args[1]
	grantee := args[2]
	status := args[3]
	if objectOfDonating == "" || donor == "" || grantee == "" || status == "" {
		return shim.Error("参数存在空值")
	}
	if donor == grantee {
		return shim.Error("捐赠人和受赠人不能同一人")
	}
	//根据objectOfDonating和donor获取想要购买的房产信息，确认存在该房产
	resultsRealEstate, err := utils.GetStateByPartialCompositeKeys2(stub, model.RealEstateKey, []string{donor, objectOfDonating})
	if err != nil || len(resultsRealEstate) != 1 {
		return shim.Error(fmt.Sprintf("根据%s和%s获取想要购买的房产信息失败: %s", objectOfDonating, donor, err))
	}
	var realEstate model.RealEstate
	if err = json.Unmarshal(resultsRealEstate[0], &realEstate); err != nil {
		return shim.Error(fmt.Sprintf("UpdateDonating-反序列化出错: %s", err))
	}
	//根据grantee获取受赠人
	resultsGranteeAccount, err := utils.GetStateByPartialCompositeKeys(stub, model.AccountKey, []string{grantee})
	if err != nil || len(resultsGranteeAccount) != 1 {
		return shim.Error(fmt.Sprintf("grantee受赠人信息验证失败%s", err))
	}
	var accountGrantee model.Account
	if err = json.Unmarshal(resultsGranteeAccount[0], &accountGrantee); err != nil {
		return shim.Error(fmt.Sprintf("查询grantee受赠人信息-反序列化出错: %s", err))
	}
	//根据objectOfDonating和donor和grantee获取捐赠信息
	resultsDonating, err := utils.GetStateByPartialCompositeKeys2(stub, model.DonatingKey, []string{donor, objectOfDonating, grantee})
	if err != nil || len(resultsDonating) != 1 {
		return shim.Error(fmt.Sprintf("根据%s和%s和%s获取销售信息失败: %s", objectOfDonating, donor, grantee, err))
	}
	var donating model.Donating
	if err = json.Unmarshal(resultsDonating[0], &donating); err != nil {
		return shim.Error(fmt.Sprintf("UpdateDonating-反序列化出错: %s", err))
	}
	//不管完成还是取消操作,必须确保捐赠处于捐赠中状态
	if donating.DonatingStatus != model.DonatingStatusConstant()["donatingStart"] {
		return shim.Error("此交易并不处于捐赠中，确认/取消捐赠失败")
	}
	//根据grantee获取买家购买信息donatingGrantee
	var donatingGrantee model.DonatingGrantee
	resultsDonatingGrantee, err := utils.GetStateByPartialCompositeKeys2(stub, model.DonatingGranteeKey, []string{grantee})
	if err != nil || len(resultsDonatingGrantee) == 0 {
		return shim.Error(fmt.Sprintf("根据%s获取受赠人信息失败: %s", grantee, err))
	}
	for _, v := range resultsDonatingGrantee {
		if v != nil {
			var s model.DonatingGrantee
			err := json.Unmarshal(v, &s)
			if err != nil {
				return shim.Error(fmt.Sprintf("UpdateDonating-反序列化出错: %s", err))
			}
			if s.Donating.ObjectOfDonating == objectOfDonating && s.Donating.Donor == donor && s.Grantee == grantee {
				//还必须判断状态必须为交付中,防止房子已经交易过，只是被取消了
				if s.Donating.DonatingStatus == model.DonatingStatusConstant()["donatingStart"] {
					donatingGrantee = s
					break
				}
			}
		}
	}
	var data []byte
	//判断捐赠状态
	switch status {
	case "done":
		//将房产信息转入受赠人，并重置担保状态
		realEstate.Proprietor = grantee
		realEstate.Encumbrance = false
		//realEstate.RealEstateID = stub.GetTxID() //重新更新房产ID
		if err := utils.WriteLedger(realEstate, stub, model.RealEstateKey, []string{realEstate.Proprietor, realEstate.RealEstateID}); err != nil {
			return shim.Error(fmt.Sprintf("%s", err))
		}
		//清除原来的房产信息
		if err := utils.DelLedger(stub, model.RealEstateKey, []string{donor, objectOfDonating}); err != nil {
			return shim.Error(fmt.Sprintf("%s", err))
		}
		//捐赠状态设置为完成，写入账本
		donating.DonatingStatus = model.DonatingStatusConstant()["done"]
		donating.ObjectOfDonating = realEstate.RealEstateID //重新更新房产ID
		if err := utils.WriteLedger(donating, stub, model.DonatingKey, []string{donating.Donor, objectOfDonating, grantee}); err != nil {
			return shim.Error(fmt.Sprintf("%s", err))
		}
		donatingGrantee.Donating = donating
		if err := utils.WriteLedger(donatingGrantee, stub, model.DonatingGranteeKey, []string{donatingGrantee.Grantee, donatingGrantee.CreateTime}); err != nil {
			return shim.Error(fmt.Sprintf("将本次捐赠交易写入账本失败%s", err))
		}
		data, err = json.Marshal(donatingGrantee)
		if err != nil {
			return shim.Error(fmt.Sprintf("序列化捐赠交易的信息出错: %s", err))
		}
		break
	case "cancelled":
		//重置房产信息担保状态
		realEstate.Encumbrance = false
		if err := utils.WriteLedger(realEstate, stub, model.RealEstateKey, []string{realEstate.Proprietor, realEstate.RealEstateID}); err != nil {
			return shim.Error(fmt.Sprintf("%s", err))
		}
		//更新捐赠状态
		donating.DonatingStatus = model.DonatingStatusConstant()["cancelled"]
		if err := utils.WriteLedger(donating, stub, model.DonatingKey, []string{donating.Donor, donating.ObjectOfDonating, donating.Grantee}); err != nil {
			return shim.Error(fmt.Sprintf("%s", err))
		}
		donatingGrantee.Donating = donating
		if err := utils.WriteLedger(donatingGrantee, stub, model.DonatingGranteeKey, []string{donatingGrantee.Grantee, donatingGrantee.CreateTime}); err != nil {
			return shim.Error(fmt.Sprintf("%s", err))
		}
		data, err = json.Marshal(donatingGrantee)
		if err != nil {
			return shim.Error(fmt.Sprintf("%s", err))
		}
		break
	default:
		return shim.Error(fmt.Sprintf("%s状态不支持", status))
	}
	return shim.Success(data)
}
