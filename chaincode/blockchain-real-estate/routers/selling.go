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
	resultsRealEstate, err := utils.GetStateByPartialCompositeKeys2(stub, lib.RealEstateKey, []string{seller, objectOfSale})
	if err != nil || len(resultsRealEstate) != 1 {
		return shim.Error(fmt.Sprintf("验证%s属于%s失败: %s", objectOfSale, seller, err))
	}
	var realEstate lib.RealEstate
	if err = json.Unmarshal(resultsRealEstate[0], &realEstate); err != nil {
		return shim.Error(fmt.Sprintf("CreateSelling-反序列化出错: %s", err))
	}
	//判断记录是否已存在，不能重复发起销售
	//若Encumbrance为true即说明此房产已经正在担保状态
	if realEstate.Encumbrance {
		return shim.Error("此房地产已经作为担保状态，不能重复发起销售")
	}
	selling := &lib.Selling{
		ObjectOfSale:  objectOfSale,
		Seller:        seller,
		Buyer:         "",
		Price:         formattedPrice,
		CreateTime:    time.Now().Local().Format("2006-01-02 15:04:05"),
		SalePeriod:    formattedSalePeriod,
		SellingStatus: lib.SellingStatusConstant()["saleStart"],
	}
	// 写入账本
	if err := utils.WriteLedger(selling, stub, lib.SellingKey, []string{selling.Seller, selling.ObjectOfSale}); err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	//将房子状态设置为正在担保状态
	realEstate.Encumbrance = true
	if err := utils.WriteLedger(realEstate, stub, lib.RealEstateKey, []string{realEstate.Proprietor, realEstate.RealEstateID}); err != nil {
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
func CreateSellingByBuy(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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
		return shim.Error("此交易不属于销售中状态，已经无法购买")
	}
	//根据buyer获取买家信息
	resultsAccount, err := utils.GetStateByPartialCompositeKeys(stub, lib.AccountKey, []string{buyer})
	if err != nil || len(resultsAccount) != 1 {
		return shim.Error(fmt.Sprintf("buyer买家信息验证失败%s", err))
	}
	var buyerAccount lib.Account
	if err = json.Unmarshal(resultsAccount[0], &buyerAccount); err != nil {
		return shim.Error(fmt.Sprintf("查询buyer买家信息-反序列化出错: %s", err))
	}
	if buyerAccount.UserName == "管理员" {
		return shim.Error(fmt.Sprintf("管理员不能购买%s", err))
	}
	//判断余额是否充足
	if buyerAccount.Balance < selling.Price {
		return shim.Error(fmt.Sprintf("房产售价为%f,您的当前余额为%f,购买失败", selling.Price, buyerAccount.Balance))
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
		CreateTime: time.Now().Local().Format("2006-01-02 15:04:05"),
		Selling:    selling,
	}
	local, _ := time.LoadLocation("Local")
	createTimeUnixNano, _ := time.ParseInLocation("2006-01-02 15:04:05", sellingBuy.CreateTime, local)
	if err := utils.WriteLedger(sellingBuy, stub, lib.SellingBuyKey, []string{sellingBuy.Buyer, fmt.Sprintf("%d", createTimeUnixNano.UnixNano())}); err != nil {
		return shim.Error(fmt.Sprintf("将本次购买交易写入账本失败%s", err))
	}
	sellingBuyByte, err := json.Marshal(sellingBuy)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化成功创建的信息出错: %s", err))
	}
	//购买成功，扣取余额，更新账本余额，注意，此时需要卖家确认收款，款项才会转入卖家账户，此处先扣除买家的余额
	buyerAccount.Balance -= selling.Price
	if err := utils.WriteLedger(buyerAccount, stub, lib.AccountKey, []string{buyerAccount.AccountId}); err != nil {
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

// 更新销售状态（买家确认、买卖家取消）
func UpdateSelling(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 验证参数
	if len(args) != 4 {
		return shim.Error("参数个数不满足")
	}
	objectOfSale := args[0]
	seller := args[1]
	buyer := args[2]
	status := args[3]
	if objectOfSale == "" || seller == "" || status == "" {
		return shim.Error("参数存在空值")
	}
	if buyer == seller {
		return shim.Error("买家和卖家不能同一人")
	}
	//根据objectOfSale和seller获取想要购买的房产信息，确认存在该房产
	resultsRealEstate, err := utils.GetStateByPartialCompositeKeys2(stub, lib.RealEstateKey, []string{seller, objectOfSale})
	if err != nil || len(resultsRealEstate) != 1 {
		return shim.Error(fmt.Sprintf("根据%s和%s获取想要购买的房产信息失败: %s", objectOfSale, seller, err))
	}
	var realEstate lib.RealEstate
	if err = json.Unmarshal(resultsRealEstate[0], &realEstate); err != nil {
		return shim.Error(fmt.Sprintf("UpdateSellingBySeller-反序列化出错: %s", err))
	}
	//根据objectOfSale和seller获取销售信息
	resultsSelling, err := utils.GetStateByPartialCompositeKeys2(stub, lib.SellingKey, []string{seller, objectOfSale})
	if err != nil || len(resultsSelling) != 1 {
		return shim.Error(fmt.Sprintf("根据%s和%s获取销售信息失败: %s", objectOfSale, seller, err))
	}
	var selling lib.Selling
	if err = json.Unmarshal(resultsSelling[0], &selling); err != nil {
		return shim.Error(fmt.Sprintf("UpdateSellingBySeller-反序列化出错: %s", err))
	}
	//根据buyer获取买家购买信息sellingBuy
	var sellingBuy lib.SellingBuy
	//如果当前状态是saleStart销售中，是不存在买家的
	if selling.SellingStatus != lib.SellingStatusConstant()["saleStart"] {
		resultsSellingByBuyer, err := utils.GetStateByPartialCompositeKeys2(stub, lib.SellingBuyKey, []string{buyer})
		if err != nil || len(resultsSellingByBuyer) == 0 {
			return shim.Error(fmt.Sprintf("根据%s获取买家购买信息失败: %s", buyer, err))
		}
		for _, v := range resultsSellingByBuyer {
			if v != nil {
				var s lib.SellingBuy
				err := json.Unmarshal(v, &s)
				if err != nil {
					return shim.Error(fmt.Sprintf("UpdateSellingBySeller-反序列化出错: %s", err))
				}
				if s.Selling.ObjectOfSale == objectOfSale && s.Selling.Seller == seller && s.Buyer == buyer {
					//还必须判断状态必须为交付中,防止房子已经交易过，只是被取消了
					if s.Selling.SellingStatus == lib.SellingStatusConstant()["delivery"] {
						sellingBuy = s
						break
					}
				}
			}
		}
	}
	var data []byte
	//判断销售状态
	switch status {
	case "done":
		//如果是买家确认收款操作,必须确保销售处于交付状态
		if selling.SellingStatus != lib.SellingStatusConstant()["delivery"] {
			return shim.Error("此交易并不处于交付中，确认收款失败")
		}
		//根据seller获取卖家信息
		resultsSellerAccount, err := utils.GetStateByPartialCompositeKeys(stub, lib.AccountKey, []string{seller})
		if err != nil || len(resultsSellerAccount) != 1 {
			return shim.Error(fmt.Sprintf("seller卖家信息验证失败%s", err))
		}
		var accountSeller lib.Account
		if err = json.Unmarshal(resultsSellerAccount[0], &accountSeller); err != nil {
			return shim.Error(fmt.Sprintf("查询seller卖家信息-反序列化出错: %s", err))
		}
		//确认收款,将款项加入到卖家账户
		accountSeller.Balance += selling.Price
		if err := utils.WriteLedger(accountSeller, stub, lib.AccountKey, []string{accountSeller.AccountId}); err != nil {
			return shim.Error(fmt.Sprintf("卖家确认接收资金失败%s", err))
		}
		//将房产信息转入买家，并重置担保状态
		realEstate.Proprietor = buyer
		realEstate.Encumbrance = false
		realEstate.RealEstateID = fmt.Sprintf("%d", time.Now().Local().UnixNano()) //重新更新房产ID
		if err := utils.WriteLedger(realEstate, stub, lib.RealEstateKey, []string{realEstate.Proprietor, realEstate.RealEstateID}); err != nil {
			return shim.Error(fmt.Sprintf("%s", err))
		}
		//清除原来的房产信息
		if err := utils.DelLedger(stub, lib.RealEstateKey, []string{seller, objectOfSale}); err != nil {
			return shim.Error(fmt.Sprintf("%s", err))
		}
		//订单状态设置为完成，写入账本
		selling.SellingStatus = lib.SellingStatusConstant()["done"]
		selling.ObjectOfSale = realEstate.RealEstateID //重新更新房产ID
		if err := utils.WriteLedger(selling, stub, lib.SellingKey, []string{selling.Seller, objectOfSale}); err != nil {
			return shim.Error(fmt.Sprintf("%s", err))
		}
		sellingBuy.Selling = selling
		local, _ := time.LoadLocation("Local")
		sellingBuyCreateTimeUnixNano, _ := time.ParseInLocation("2006-01-02 15:04:05", sellingBuy.CreateTime, local)
		if err := utils.WriteLedger(sellingBuy, stub, lib.SellingBuyKey, []string{sellingBuy.Buyer, fmt.Sprintf("%d", sellingBuyCreateTimeUnixNano.UnixNano())}); err != nil {
			return shim.Error(fmt.Sprintf("将本次购买交易写入账本失败%s", err))
		}
		data, err = json.Marshal(sellingBuy)
		if err != nil {
			return shim.Error(fmt.Sprintf("序列化购买交易的信息出错: %s", err))
		}
		break
	case "cancelled":
		data, err = closeSelling("cancelled", selling, realEstate, sellingBuy, buyer, stub)
		if err != nil {
			return shim.Error(fmt.Sprintf("%s", err))
		}
		break
	case "expired":
		data, err = closeSelling("expired", selling, realEstate, sellingBuy, buyer, stub)
		if err != nil {
			return shim.Error(fmt.Sprintf("%s", err))
		}
		break
	default:
		return shim.Error(fmt.Sprintf("%s状态不支持", status))
	}
	return shim.Success(data)
}

//不管是取消还是过期，都分两种情况
//1、当前处于saleStart销售状态
//2、当前处于delivery交付中状态
func closeSelling(closeStart string, selling lib.Selling, realEstate lib.RealEstate, sellingBuy lib.SellingBuy, buyer string, stub shim.ChaincodeStubInterface) ([]byte, error) {
	switch selling.SellingStatus {
	case lib.SellingStatusConstant()["saleStart"]:
		selling.SellingStatus = lib.SellingStatusConstant()[closeStart]
		//重置房产信息担保状态
		realEstate.Encumbrance = false
		if err := utils.WriteLedger(realEstate, stub, lib.RealEstateKey, []string{realEstate.Proprietor, realEstate.RealEstateID}); err != nil {
			return nil, err
		}
		if err := utils.WriteLedger(selling, stub, lib.SellingKey, []string{selling.Seller, selling.ObjectOfSale}); err != nil {
			return nil, err
		}
		data, err := json.Marshal(selling)
		if err != nil {
			return nil, err
		}
		return data, nil
	case lib.SellingStatusConstant()["delivery"]:
		//根据buyer获取卖家信息
		resultsBuyerAccount, err := utils.GetStateByPartialCompositeKeys(stub, lib.AccountKey, []string{buyer})
		if err != nil || len(resultsBuyerAccount) != 1 {
			return nil, err
		}
		var accountBuyer lib.Account
		if err = json.Unmarshal(resultsBuyerAccount[0], &accountBuyer); err != nil {
			return nil, err
		}
		//此时取消操作，需要将资金退还给买家
		accountBuyer.Balance += selling.Price
		if err := utils.WriteLedger(accountBuyer, stub, lib.AccountKey, []string{accountBuyer.AccountId}); err != nil {
			return nil, err
		}
		//重置房产信息担保状态
		realEstate.Encumbrance = false
		if err := utils.WriteLedger(realEstate, stub, lib.RealEstateKey, []string{realEstate.Proprietor, realEstate.RealEstateID}); err != nil {
			return nil, err
		}
		//更新销售状态
		selling.SellingStatus = lib.SellingStatusConstant()[closeStart]
		if err := utils.WriteLedger(selling, stub, lib.SellingKey, []string{selling.Seller, selling.ObjectOfSale}); err != nil {
			return nil, err
		}
		sellingBuy.Selling = selling
		local, _ := time.LoadLocation("Local")
		sellingBuyCreateTimeUnixNano, _ := time.ParseInLocation("2006-01-02 15:04:05", sellingBuy.CreateTime, local)
		if err := utils.WriteLedger(sellingBuy, stub, lib.SellingBuyKey, []string{sellingBuy.Buyer, fmt.Sprintf("%d", sellingBuyCreateTimeUnixNano.UnixNano())}); err != nil {
			return nil, err
		}
		data, err := json.Marshal(sellingBuy)
		if err != nil {
			return nil, err
		}
		return data, nil
	default:
		return nil, nil
	}
}
