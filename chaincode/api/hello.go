package api

import (
	"chaincode/pkg/utils"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// Hello 测试
func Hello(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	err := utils.WriteLedger(map[string]interface{}{"msg": "hello"}, stub, "hello", []string{})
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte("hello world"))
}
