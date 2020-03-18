package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/togettoyou/blockchain-real-estate/chaincode/blockchain-real-estate/lib"
	"testing"
)

func initTest(t *testing.T) *shim.MockStub {
	scc := new(BlockChainRealEstate)
	stub := shim.NewMockStub("ex01", scc)
	checkInit(t, stub, [][]byte{[]byte("init")})
	return stub
}

func checkInit(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInit("1", args)
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Message))
		t.FailNow()
	}
}

func checkInvoke(t *testing.T, stub *shim.MockStub, args [][]byte) pb.Response {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		fmt.Println("Invoke", args, "failed", string(res.Message))
		t.FailNow()
	}
	return res
}

// 测试链码初始化
func TestBlockChainRealEstate_Init(t *testing.T) {
	initTest(t)
}

// 测试获取账户信息
func Test_QueryAccountList(t *testing.T) {
	stub := initTest(t)
	fmt.Println(fmt.Sprintf("1、测试获取所有数据\n%s",
		string(checkInvoke(t, stub, [][]byte{
			[]byte("queryAccountList"),
		}).Payload)))
	fmt.Println(fmt.Sprintf("2、测试获取多个数据\n%s",
		string(checkInvoke(t, stub, [][]byte{
			[]byte("queryAccountList"),
			[]byte("5feceb66ffc8"),
			[]byte("6b86b273ff34"),
		}).Payload)))
	fmt.Println(fmt.Sprintf("3、测试获取单个数据\n%s",
		string(checkInvoke(t, stub, [][]byte{
			[]byte("queryAccountList"),
			[]byte("4e07408562be"),
		}).Payload)))
	fmt.Println(fmt.Sprintf("4、测试获取无效数据\n%s",
		string(checkInvoke(t, stub, [][]byte{
			[]byte("queryAccountList"),
			[]byte("0"),
		}).Payload)))
}

// 测试创建房地产
func Test_CreateRealEstate(t *testing.T) {
	stub := initTest(t)
	//成功
	checkInvoke(t, stub, [][]byte{
		[]byte("createRealEstate"),
		[]byte("5feceb66ffc8"), //操作人
		[]byte("6b86b273ff34"), //所有者
		[]byte("50"),           //总面积
		[]byte("30"),           //生活空间
	})
	//操作人权限不足
	checkInvoke(t, stub, [][]byte{
		[]byte("createRealEstate"),
		[]byte("6b86b273ff34"), //操作人
		[]byte("4e07408562be"), //所有者
		[]byte("50"),           //总面积
		[]byte("30"),           //生活空间
	})
	//操作人应为管理员且与所有人不能相同
	checkInvoke(t, stub, [][]byte{
		[]byte("createRealEstate"),
		[]byte("5feceb66ffc8"), //操作人
		[]byte("5feceb66ffc8"), //所有者
		[]byte("50"),           //总面积
		[]byte("30"),           //生活空间
	})
	//业主proprietor信息验证失败
	checkInvoke(t, stub, [][]byte{
		[]byte("createRealEstate"),
		[]byte("5feceb66ffc8"),    //操作人
		[]byte("6b86b273ff34555"), //所有者
		[]byte("50"),              //总面积
		[]byte("30"),              //生活空间
	})
	//参数个数不满足
	checkInvoke(t, stub, [][]byte{
		[]byte("createRealEstate"),
		[]byte("5feceb66ffc8"), //操作人
		[]byte("6b86b273ff34"), //所有者
		[]byte("50"),           //总面积
	})
	//参数格式转换出错
	checkInvoke(t, stub, [][]byte{
		[]byte("createRealEstate"),
		[]byte("5feceb66ffc8"), //操作人
		[]byte("6b86b273ff34"), //所有者
		[]byte("50f"),          //总面积
		[]byte("30"),           //生活空间
	})
}

//手动创建一些房地产
func checkCreateRealEstate(stub *shim.MockStub, t *testing.T) []lib.RealEstate {
	var realEstateList []lib.RealEstate
	var realEstate lib.RealEstate
	//成功
	resp1 := checkInvoke(t, stub, [][]byte{
		[]byte("createRealEstate"),
		[]byte("5feceb66ffc8"), //操作人
		[]byte("6b86b273ff34"), //所有者
		[]byte("50"),           //总面积
		[]byte("30"),           //生活空间
	})
	resp2 := checkInvoke(t, stub, [][]byte{
		[]byte("createRealEstate"),
		[]byte("5feceb66ffc8"), //操作人
		[]byte("6b86b273ff34"), //所有者
		[]byte("80"),           //总面积
		[]byte("60.8"),         //生活空间
	})
	resp3 := checkInvoke(t, stub, [][]byte{
		[]byte("createRealEstate"),
		[]byte("5feceb66ffc8"), //操作人
		[]byte("4e07408562be"), //所有者
		[]byte("60"),           //总面积
		[]byte("40"),           //生活空间
	})
	resp4 := checkInvoke(t, stub, [][]byte{
		[]byte("createRealEstate"),
		[]byte("5feceb66ffc8"), //操作人
		[]byte("ef2d127de37b"), //所有者
		[]byte("80"),           //总面积
		[]byte("60"),           //生活空间
	})
	json.Unmarshal(bytes.NewBuffer(resp1.Payload).Bytes(), &realEstate)
	realEstateList = append(realEstateList, realEstate)
	json.Unmarshal(bytes.NewBuffer(resp2.Payload).Bytes(), &realEstate)
	realEstateList = append(realEstateList, realEstate)
	json.Unmarshal(bytes.NewBuffer(resp3.Payload).Bytes(), &realEstate)
	realEstateList = append(realEstateList, realEstate)
	json.Unmarshal(bytes.NewBuffer(resp4.Payload).Bytes(), &realEstate)
	realEstateList = append(realEstateList, realEstate)
	return realEstateList
}

// 测试获取房地产信息
func Test_QueryRealEstateList(t *testing.T) {
	stub := initTest(t)
	realEstateList := checkCreateRealEstate(stub, t)

	fmt.Println(fmt.Sprintf("1、测试获取所有数据\n%s",
		string(checkInvoke(t, stub, [][]byte{
			[]byte("queryRealEstateList"),
		}).Payload)))
	fmt.Println(fmt.Sprintf("2、测试获取指定数据\n%s",
		string(checkInvoke(t, stub, [][]byte{
			[]byte("queryRealEstateList"),
			[]byte(realEstateList[0].Proprietor),
			[]byte(realEstateList[0].RealEstateID),
		}).Payload)))
	fmt.Println(fmt.Sprintf("3、测试获取无效数据\n%s",
		string(checkInvoke(t, stub, [][]byte{
			[]byte("queryRealEstateList"),
			[]byte("0"),
		}).Payload)))
}

// 测试发起销售
func Test_CreateSelling(t *testing.T) {
	stub := initTest(t)
	realEstateList := checkCreateRealEstate(stub, t)
	//成功
	checkInvoke(t, stub, [][]byte{
		[]byte("createSelling"),
		[]byte(realEstateList[0].RealEstateID), //销售对象(正在出售的房地产RealEstateID)
		[]byte(realEstateList[0].Proprietor),   //卖家(卖家AccountId)
		[]byte("50"),                           //价格
		[]byte("30"),                           //智能合约的有效期(单位为天)
	})
	//验证销售对象objectOfSale属于卖家seller失败
	checkInvoke(t, stub, [][]byte{
		[]byte("createSelling"),
		[]byte(realEstateList[0].RealEstateID), //销售对象(正在出售的房地产RealEstateID)
		[]byte(realEstateList[2].Proprietor),   //卖家(卖家AccountId)
		[]byte("50"),                           //价格
		[]byte("30"),                           //智能合约的有效期(单位为天)
	})
	checkInvoke(t, stub, [][]byte{
		[]byte("createSelling"),
		[]byte("123"),                        //销售对象(正在出售的房地产RealEstateID)
		[]byte(realEstateList[0].Proprietor), //卖家(卖家AccountId)
		[]byte("50"),                         //价格
		[]byte("30"),                         //智能合约的有效期(单位为天)
	})
	//参数错误
	checkInvoke(t, stub, [][]byte{
		[]byte("createSelling"),
		[]byte(realEstateList[0].RealEstateID), //销售对象(正在出售的房地产RealEstateID)
		[]byte(realEstateList[0].Proprietor),   //卖家(卖家AccountId)
		[]byte("50"),                           //价格
	})
	checkInvoke(t, stub, [][]byte{
		[]byte("createSelling"),
		[]byte(""),                           //销售对象(正在出售的房地产RealEstateID)
		[]byte(realEstateList[0].Proprietor), //卖家(卖家AccountId)
		[]byte("50"),                         //价格
		[]byte("30"),                         //智能合约的有效期(单位为天)
	})
}

// 测试销售发起、购买等操作
func Test_QuerySellingList(t *testing.T) {
	stub := initTest(t)
	realEstateList := checkCreateRealEstate(stub, t)
	//先发起
	fmt.Println(fmt.Sprintf("发起\n%s", string(checkInvoke(t, stub, [][]byte{
		[]byte("createSelling"),
		[]byte(realEstateList[0].RealEstateID), //销售对象(正在出售的房地产RealEstateID)
		[]byte(realEstateList[0].Proprietor),   //卖家(卖家AccountId)
		[]byte("500000"),                       //价格
		[]byte("30"),                           //智能合约的有效期(单位为天)
	}).Payload)))
	fmt.Println(fmt.Sprintf("发起\n%s", string(checkInvoke(t, stub, [][]byte{
		[]byte("createSelling"),
		[]byte(realEstateList[2].RealEstateID), //销售对象(正在出售的房地产RealEstateID)
		[]byte(realEstateList[2].Proprietor),   //卖家(卖家AccountId)
		[]byte("600000"),                       //价格
		[]byte("40"),                           //智能合约的有效期(单位为天)
	}).Payload)))
	//查询成功
	fmt.Println(fmt.Sprintf("1、查询所有\n%s", string(checkInvoke(t, stub, [][]byte{
		[]byte("querySellingList"),
	}).Payload)))
	fmt.Println(fmt.Sprintf("2、查询指定%s\n%s", realEstateList[0].Proprietor, string(checkInvoke(t, stub, [][]byte{
		[]byte("querySellingList"),
		[]byte(realEstateList[0].Proprietor),
	}).Payload)))
	//购买
	fmt.Println(fmt.Sprintf("3、购买前先查询%s的账户余额\n%s", realEstateList[2].Proprietor, string(checkInvoke(t, stub, [][]byte{
		[]byte("queryAccountList"),
		[]byte(realEstateList[2].Proprietor),
	}).Payload)))
	fmt.Println(fmt.Sprintf("4、开始购买\n%s", string(checkInvoke(t, stub, [][]byte{
		[]byte("createSellingByBuy"),
		[]byte(realEstateList[0].RealEstateID), //销售对象(正在出售的房地产RealEstateID)
		[]byte(realEstateList[0].Proprietor),   //卖家(卖家AccountId)
		[]byte(realEstateList[2].Proprietor),   //买家(买家AccountId)
	}).Payload)))
	fmt.Println(fmt.Sprintf("》购买后再次查询%s的账户余额\n%s", realEstateList[2].Proprietor, string(checkInvoke(t, stub, [][]byte{
		[]byte("queryAccountList"),
		[]byte(realEstateList[2].Proprietor),
	}).Payload)))
	fmt.Println(fmt.Sprintf("》卖家查询购买成功信息\n%s", string(checkInvoke(t, stub, [][]byte{
		[]byte("querySellingList"),
		[]byte(realEstateList[0].Proprietor), //买家(买家AccountId)
	}).Payload)))
	fmt.Println(fmt.Sprintf("》买家查询购买成功信息\n%s", string(checkInvoke(t, stub, [][]byte{
		[]byte("querySellingListByBuyer"),
		[]byte(realEstateList[2].Proprietor), //买家(买家AccountId)
	}).Payload)))
	fmt.Println(fmt.Sprintf("》确认收款前卖家%s的账户余额\n%s", realEstateList[0].Proprietor, string(checkInvoke(t, stub, [][]byte{
		[]byte("queryAccountList"),
		[]byte(realEstateList[0].Proprietor),
	}).Payload)))
	fmt.Println(fmt.Sprintf("》确认收款前买家%s的账户余额\n%s", realEstateList[2].Proprietor, string(checkInvoke(t, stub, [][]byte{
		[]byte("queryAccountList"),
		[]byte(realEstateList[2].Proprietor),
	}).Payload)))
	fmt.Println(fmt.Sprintf("》确认收款前卖家%s的房产信息\n%s", realEstateList[0].Proprietor, string(checkInvoke(t, stub, [][]byte{
		[]byte("queryRealEstateList"),
		[]byte(realEstateList[0].Proprietor),
	}).Payload)))
	fmt.Println(fmt.Sprintf("》确认收款前买家%s的房产信息\n%s", realEstateList[2].Proprietor, string(checkInvoke(t, stub, [][]byte{
		[]byte("queryRealEstateList"),
		[]byte(realEstateList[2].Proprietor),
	}).Payload)))
	fmt.Println(fmt.Sprintf("》卖家确认收款\n%s", string(checkInvoke(t, stub, [][]byte{
		[]byte("updateSelling"),
		[]byte(realEstateList[0].RealEstateID), //销售对象(正在出售的房地产RealEstateID)
		[]byte(realEstateList[0].Proprietor),   //卖家(卖家AccountId)
		[]byte(realEstateList[2].Proprietor),   //买家(买家AccountId)
		[]byte("done"),                         //确认收款
	}).Payload)))
	//fmt.Println(fmt.Sprintf("》卖家取消收款\n%s", string(checkInvoke(t, stub, [][]byte{
	//	[]byte("updateSelling"),
	//	[]byte(realEstateList[0].RealEstateID), //销售对象(正在出售的房地产RealEstateID)
	//	[]byte(realEstateList[0].Proprietor),   //卖家(卖家AccountId)
	//	[]byte(realEstateList[2].Proprietor),   //买家(买家AccountId)
	//	[]byte("cancelled"),                    //取消收款
	//}).Payload)))
	fmt.Println(fmt.Sprintf("》确认收款后卖家%s的账户余额\n%s", realEstateList[0].Proprietor, string(checkInvoke(t, stub, [][]byte{
		[]byte("queryAccountList"),
		[]byte(realEstateList[0].Proprietor),
	}).Payload)))
	fmt.Println(fmt.Sprintf("》确认收款后买家%s的账户余额\n%s", realEstateList[2].Proprietor, string(checkInvoke(t, stub, [][]byte{
		[]byte("queryAccountList"),
		[]byte(realEstateList[2].Proprietor),
	}).Payload)))
	fmt.Println(fmt.Sprintf("》确认收款后卖家%s的房产信息\n%s", realEstateList[0].Proprietor, string(checkInvoke(t, stub, [][]byte{
		[]byte("queryRealEstateList"),
		[]byte(realEstateList[0].Proprietor),
	}).Payload)))
	fmt.Println(fmt.Sprintf("》确认收款后买家%s的房产信息\n%s", realEstateList[2].Proprietor, string(checkInvoke(t, stub, [][]byte{
		[]byte("queryRealEstateList"),
		[]byte(realEstateList[2].Proprietor),
	}).Payload)))
	fmt.Println(fmt.Sprintf("》卖家查询购买成功信息\n%s", string(checkInvoke(t, stub, [][]byte{
		[]byte("querySellingList"),
		[]byte(realEstateList[0].Proprietor), //买家(买家AccountId)
	}).Payload)))
	fmt.Println(fmt.Sprintf("》买家查询购买成功信息\n%s", string(checkInvoke(t, stub, [][]byte{
		[]byte("querySellingListByBuyer"),
		[]byte(realEstateList[2].Proprietor), //买家(买家AccountId)
	}).Payload)))
}

// 测试捐赠合约
func Test_Donating(t *testing.T) {
	stub := initTest(t)
	realEstateList := checkCreateRealEstate(stub, t)

	fmt.Println(fmt.Sprintf("获取房地产信息\n%s",
		string(checkInvoke(t, stub, [][]byte{
			[]byte("queryRealEstateList"),
		}).Payload)))
	//先发起
	fmt.Println(fmt.Sprintf("发起捐赠\n%s", string(checkInvoke(t, stub, [][]byte{
		[]byte("createDonating"),
		[]byte(realEstateList[0].RealEstateID),
		[]byte(realEstateList[0].Proprietor),
		[]byte(realEstateList[2].Proprietor),
	}).Payload)))

	fmt.Println(fmt.Sprintf("获取房地产信息\n%s",
		string(checkInvoke(t, stub, [][]byte{
			[]byte("queryRealEstateList"),
		}).Payload)))

	fmt.Println(fmt.Sprintf("1、查询所有\n%s", string(checkInvoke(t, stub, [][]byte{
		[]byte("queryDonatingList"),
	}).Payload)))
	fmt.Println(fmt.Sprintf("2、查询指定%s\n%s", realEstateList[0].Proprietor, string(checkInvoke(t, stub, [][]byte{
		[]byte("queryDonatingList"),
		[]byte(realEstateList[2].Proprietor),
	}).Payload)))
	fmt.Println(fmt.Sprintf("3、查询指定受赠%s\n%s", realEstateList[0].Proprietor, string(checkInvoke(t, stub, [][]byte{
		[]byte("queryDonatingListByGrantee"),
		[]byte(realEstateList[2].Proprietor),
	}).Payload)))

	//fmt.Println(fmt.Sprintf("4、接受受赠%s\n%s", realEstateList[0].Proprietor, string(checkInvoke(t, stub, [][]byte{
	//	[]byte("updateDonating"),
	//	[]byte(realEstateList[0].RealEstateID),
	//	[]byte(realEstateList[0].Proprietor),
	//	[]byte(realEstateList[2].Proprietor),
	//	[]byte("done"),
	//}).Payload)))
	fmt.Println(fmt.Sprintf("4、取消受赠%s\n%s", realEstateList[0].Proprietor, string(checkInvoke(t, stub, [][]byte{
		[]byte("updateDonating"),
		[]byte(realEstateList[0].RealEstateID),
		[]byte(realEstateList[0].Proprietor),
		[]byte(realEstateList[2].Proprietor),
		[]byte("cancelled"),
	}).Payload)))

	fmt.Println(fmt.Sprintf("获取房地产信息\n%s",
		string(checkInvoke(t, stub, [][]byte{
			[]byte("queryRealEstateList"),
		}).Payload)))
}
