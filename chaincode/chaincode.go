package main

import (
	"fmt"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

const (
	InitLedgerKey = "InitLedgerKey"
)

type Data struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	log.Println("InitLedger", ctx.GetStub().GetTxID())

	_ = ctx.GetStub().DelState(InitLedgerKey)
	return ctx.GetStub().PutState(InitLedgerKey, []byte("hello"))
}

func (s *SmartContract) Hello(ctx contractapi.TransactionContextInterface) (string, error) {
	log.Println("Hello", ctx.GetStub().GetTxID())

	str, err := ctx.GetStub().GetState(InitLedgerKey)
	if err != nil {
		return "", fmt.Errorf("failed to read from world state: %v", err)
	}
	if str == nil {
		return "", fmt.Errorf("the InitLedgerKey does not exist")
	}

	return string(str), nil
}

func main() {
	cc, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("Error creating chaincode: %v", err)
	}

	if err := cc.Start(); err != nil {
		log.Panicf("Error starting chaincode: %v", err)
	}
}
