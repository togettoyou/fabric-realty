package main

import (
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
	fmt.Println("1、测试获取所有数据")
	response := checkInvoke(t, stub, [][]byte{[]byte("queryAccountList")})
	var allAccountList []lib.Account
	err := json.Unmarshal(response.Payload, &allAccountList)
	if err != nil {
		fmt.Printf("Unmarshal err: %s\n", err.Error())
		t.FailNow()
	}
	fmt.Println(allAccountList)

	fmt.Println("2、测试获取多个数据")
	response = checkInvoke(t, stub, [][]byte{[]byte("queryAccountList"), []byte("5feceb66ffc8"), []byte("6b86b273ff34")})
	var accounts []lib.Account
	err = json.Unmarshal(response.Payload, &accounts)
	if err != nil {
		fmt.Printf("Unmarshal err: %s\n", err.Error())
		t.FailNow()
	}
	fmt.Println(accounts)

	fmt.Println("3、测试获取单个数据")
	response = checkInvoke(t, stub, [][]byte{[]byte("queryAccountList"), []byte("4e07408562be")})
	var account []lib.Account
	err = json.Unmarshal(response.Payload, &account)
	if err != nil {
		fmt.Printf("Unmarshal err: %s\n", err.Error())
		t.FailNow()
	}
	fmt.Println(account)

	fmt.Println("4、测试获取无效数据")
	response = checkInvoke(t, stub, [][]byte{[]byte("queryAccountList"), []byte("0")})
	var noneAccount []lib.Account
	err = json.Unmarshal(response.Payload, &noneAccount)
	if err != nil {
		fmt.Printf("Unmarshal err: %s\n", err.Error())
		t.FailNow()
	}
	fmt.Println(noneAccount)
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

// 测试获取房地产信息
func Test_QueryRealEstateList(t *testing.T) {
	stub := initTest(t)
	//成功
	checkInvoke(t, stub, [][]byte{
		[]byte("createRealEstate"),
		[]byte("5feceb66ffc8"), //操作人
		[]byte("6b86b273ff34"), //所有者
		[]byte("50"),           //总面积
		[]byte("30"),           //生活空间
	})
	checkInvoke(t, stub, [][]byte{
		[]byte("createRealEstate"),
		[]byte("5feceb66ffc8"), //操作人
		[]byte("4e07408562be"), //所有者
		[]byte("60"),           //总面积
		[]byte("40"),           //生活空间
	})
	checkInvoke(t, stub, [][]byte{
		[]byte("createRealEstate"),
		[]byte("5feceb66ffc8"), //操作人
		[]byte("ef2d127de37b"), //所有者
		[]byte("80"),           //总面积
		[]byte("60"),           //生活空间
	})
	fmt.Println("1、测试获取所有数据")
	response := checkInvoke(t, stub, [][]byte{[]byte("queryRealEstateList")})
	var allRealEstateList []lib.RealEstate
	err := json.Unmarshal(response.Payload, &allRealEstateList)
	if err != nil {
		fmt.Printf("Unmarshal err: %s\n", err.Error())
		t.FailNow()
	}
	fmt.Println(allRealEstateList)

	fmt.Println("2、测试获取指定数据")
	response = checkInvoke(t, stub, [][]byte{[]byte("queryRealEstateList"), []byte("ef2d127de37b")})
	var realEstates []lib.RealEstate
	err = json.Unmarshal(response.Payload, &realEstates)
	if err != nil {
		fmt.Printf("Unmarshal err: %s\n", err.Error())
		t.FailNow()
	}
	fmt.Println(realEstates)

	fmt.Println("3、测试获取无效数据")
	response = checkInvoke(t, stub, [][]byte{[]byte("queryAccountList"), []byte("0")})
	var noneRealEstate []lib.RealEstate
	err = json.Unmarshal(response.Payload, &noneRealEstate)
	if err != nil {
		fmt.Printf("Unmarshal err: %s\n", err.Error())
		t.FailNow()
	}
	fmt.Println(noneRealEstate)
}
