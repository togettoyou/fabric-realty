package api

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// Hello 测试
func Hello(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success([]byte("hello world"))
}
