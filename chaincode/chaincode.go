package main

import (
	"chaincode/api"
	"chaincode/model"
	"chaincode/pkg/utils"
	"fmt"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type BlockChainRealEstate struct {
}

// Init 链码初始化
func (t *BlockChainRealEstate) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("链码初始化")
	//初始化默认数据
	var accountIds = [6]string{
		"5feceb66ffc8",
		"6b86b273ff34",
		"d4735e3a265e",
		"4e07408562be",
		"4b227777d4dd",
		"ef2d127de37b",
	}
	var userNames = [6]string{"管理员", "①号业主", "②号业主", "③号业主", "④号业主", "⑤号业主"}
	var balances = [6]float64{0, 5000000, 5000000, 5000000, 5000000, 5000000}
	//初始化账号数据
	for i, val := range accountIds {
		account := &model.Account{
			AccountId: val,
			UserName:  userNames[i],
			Balance:   balances[i],
		}
		// 写入账本
		if err := utils.WriteLedger(account, stub, model.AccountKey, []string{val}); err != nil {
			return shim.Error(fmt.Sprintf("%s", err))
		}
	}
	return shim.Success(nil)
}

// Invoke 实现Invoke接口调用智能合约
func (t *BlockChainRealEstate) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	funcName, args := stub.GetFunctionAndParameters()
	switch funcName {
	case "hello":
		return api.Hello(stub, args)
	case "queryAccountList":
		return api.QueryAccountList(stub, args)
	case "createRealEstate":
		return api.CreateRealEstate(stub, args)
	case "queryRealEstateList":
		return api.QueryRealEstateList(stub, args)
	case "createSelling":
		return api.CreateSelling(stub, args)
	case "createSellingByBuy":
		return api.CreateSellingByBuy(stub, args)
	case "querySellingList":
		return api.QuerySellingList(stub, args)
	case "querySellingListByBuyer":
		return api.QuerySellingListByBuyer(stub, args)
	case "updateSelling":
		return api.UpdateSelling(stub, args)
	case "createDonating":
		return api.CreateDonating(stub, args)
	case "queryDonatingList":
		return api.QueryDonatingList(stub, args)
	case "queryDonatingListByGrantee":
		return api.QueryDonatingListByGrantee(stub, args)
	case "updateDonating":
		return api.UpdateDonating(stub, args)
	default:
		return shim.Error(fmt.Sprintf("没有该功能: %s", funcName))
	}
}

func main() {
	timeLocal, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	time.Local = timeLocal
	err = shim.Start(new(BlockChainRealEstate))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
