/**
 * @Author: 夜央 Oh oh oh oh oh oh (https://github.com/togettoyou)
 * @Email: zoujh99@qq.com
 * @Date: 2020/3/10 6:40 下午
 * @Description: 销售相关合约路由
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

//发起销售
func CreateSelling(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 验证参数
	if len(args) != 4 {
		return shim.Error("参数个数不满足")
	}
	objectOfSale := args[0]
	seller := args[1]
	price := args[2]
	salePeriod := args[3]
	if objectOfSale == "" || seller == "" || price == "" || salePeriod == "" {
		return shim.Error("参数存在空值")
	}
	// 参数数据格式转换
	var formattedPrice float64
	if val, err := strconv.ParseFloat(price, 64); err != nil {
		return shim.Error(fmt.Sprintf("price参数格式转换出错: %s", err))
	} else {
		formattedPrice = val
	}
	var formattedSalePeriod int
	if val, err := strconv.Atoi(salePeriod); err != nil {
		return shim.Error(fmt.Sprintf("salePeriod参数格式转换出错: %s", err))
	} else {
		formattedSalePeriod = val
	}
	//判断objectOfSale是否属于seller
	resultsSeller, err := utils.GetStateByPartialCompositeKeys2(stub, lib.RealEstateKey, []string{seller, objectOfSale})
	if err != nil || len(resultsSeller) != 1 {
		return shim.Error(fmt.Sprintf("验证%s属于%s失败: %s", objectOfSale, seller, err))
	}
	selling := &lib.Selling{
		ObjectOfSale:  objectOfSale,
		Seller:        seller,
		Buyer:         "",
		Price:         formattedPrice,
		CreateTime:    time.Now(),
		SalePeriod:    formattedSalePeriod,
		SellingStatus: lib.SellingStatusConstant()["saleStart"],
	}
	// 写入账本
	if err := utils.WriteLedger(selling, stub, lib.SellingKey, []string{selling.Seller, selling.ObjectOfSale}); err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	//将成功创建的信息返回
	sellingByte, err := json.Marshal(selling)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化成功创建的信息出错: %s", err))
	}
	// 成功返回
	return shim.Success(sellingByte)
}

//参与销售(买家购买)
func CreateSellingBuy(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 验证参数
	if len(args) != 3 {
		return shim.Error("参数个数不满足")
	}
	objectOfSale := args[0]
	seller := args[1]
	buyer := args[2]
	if objectOfSale == "" || seller == "" || buyer == "" {
		return shim.Error("参数存在空值")
	}
	if seller == buyer {
		return shim.Error("买家和卖家不能同一人")
	}
	//根据objectOfSale和seller获取想要购买的房产信息，确认存在该房产
	resultsRealEstate, err := utils.GetStateByPartialCompositeKeys2(stub, lib.RealEstateKey, []string{seller, objectOfSale})
	if err != nil || len(resultsRealEstate) != 1 {
		return shim.Error(fmt.Sprintf("根据%s和%s获取想要购买的房产信息失败: %s", objectOfSale, seller, err))
	}
	//根据objectOfSale和seller获取销售信息
	resultsSelling, err := utils.GetStateByPartialCompositeKeys2(stub, lib.SellingKey, []string{seller, objectOfSale})
	if err != nil || len(resultsSelling) != 1 {
		return shim.Error(fmt.Sprintf("根据%s和%s获取销售信息失败: %s", objectOfSale, seller, err))
	}
	var selling lib.Selling
	if err = json.Unmarshal(resultsSelling[0], &selling); err != nil {
		return shim.Error(fmt.Sprintf("CreateSellingBuy-反序列化出错: %s", err))
	}
	//判断selling的状态是否为销售中
	if selling.SellingStatus != lib.SellingStatusConstant()["saleStart"] {
		return shim.Error("此交易已经正在进行中，购买失败")
	}
	//根据buyer获取买家信息
	resultsAccount, err := utils.GetStateByPartialCompositeKeys(stub, lib.AccountKey, []string{buyer})
	if err != nil || len(resultsAccount) != 1 {
		return shim.Error(fmt.Sprintf("buyer买家信息验证失败%s", err))
	}
	var account lib.Account
	if err = json.Unmarshal(resultsAccount[0], &account); err != nil {
		return shim.Error(fmt.Sprintf("查询buyer买家信息-反序列化出错: %s", err))
	}
	//判断余额是否充足
	if account.Balance < selling.Price {
		return shim.Error(fmt.Sprintf("房产售价为%f,您的当前余额为%f,购买失败", selling.Price, account.Balance))
	}
	//将buyer写入交易selling,修改交易状态
	selling.Buyer = buyer
	selling.SellingStatus = lib.SellingStatusConstant()["delivery"]
	if err := utils.WriteLedger(selling, stub, lib.SellingKey, []string{selling.Seller, selling.ObjectOfSale}); err != nil {
		return shim.Error(fmt.Sprintf("将buyer写入交易selling,修改交易状态 失败%s", err))
	}
	//将本次购买交易写入账本,可供买家查询
	sellingBuy := &lib.SellingBuy{
		Buyer:      buyer,
		CreateTime: time.Now(),
		Selling:    selling,
	}
	if err := utils.WriteLedger(sellingBuy, stub, lib.SellingBuyKey, []string{sellingBuy.Buyer, sellingBuy.CreateTime.String()}); err != nil {
		return shim.Error(fmt.Sprintf("将本次购买交易写入账本失败%s", err))
	}
	sellingBuyByte, err := json.Marshal(sellingBuy)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化成功创建的信息出错: %s", err))
	}
	//购买成功，扣取余额，更新账本余额，注意，此时需要卖家确认收款，款项才会转入卖家账户，此处先扣除买家的余额
	account.Balance -= selling.Price
	if err := utils.WriteLedger(account, stub, lib.AccountKey, []string{account.AccountId}); err != nil {
		return shim.Error(fmt.Sprintf("扣取买家余额失败%s", err))
	}
	// 成功返回
	return shim.Success(sellingBuyByte)
}

//查询销售(可查询所有，也可根据发起销售人查询)(发起的)(供卖家查询)
func QuerySellingList(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var sellingList []lib.Selling
	results, err := utils.GetStateByPartialCompositeKeys2(stub, lib.SellingKey, args)
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	for _, v := range results {
		if v != nil {
			var selling lib.Selling
			err := json.Unmarshal(v, &selling)
			if err != nil {
				return shim.Error(fmt.Sprintf("QuerySellingList-反序列化出错: %s", err))
			}
			sellingList = append(sellingList, selling)
		}
	}
	sellingListByte, err := json.Marshal(sellingList)
	if err != nil {
		return shim.Error(fmt.Sprintf("QuerySellingList-序列化出错: %s", err))
	}
	return shim.Success(sellingListByte)
}

//根据参与销售人、买家(买家AccountId)查询销售(参与的)(供买家查询)
func QuerySellingListByBuyer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error(fmt.Sprintf("必须指定买家AccountId查询"))
	}
	var sellingBuyList []lib.SellingBuy
	results, err := utils.GetStateByPartialCompositeKeys2(stub, lib.SellingBuyKey, args)
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	for _, v := range results {
		if v != nil {
			var sellingBuy lib.SellingBuy
			err := json.Unmarshal(v, &sellingBuy)
			if err != nil {
				return shim.Error(fmt.Sprintf("QuerySellingListByBuyer-反序列化出错: %s", err))
			}
			sellingBuyList = append(sellingBuyList, sellingBuy)
		}
	}
	sellingBuyListByte, err := json.Marshal(sellingBuyList)
	if err != nil {
		return shim.Error(fmt.Sprintf("QuerySellingListByBuyer-序列化出错: %s", err))
	}
	return shim.Success(sellingBuyListByte)
}
