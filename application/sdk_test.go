package main_test

import (
	"fmt"
	"github.com/togettoyou/blockchain-real-estate/application/blockchain"
	"testing"
)

func TestInvoke_QueryAccountList(t *testing.T) {
	blockchain.Init()
	response, e := blockchain.ChannelQuery("hello", [][]byte{})
	if e != nil {
		fmt.Println(e.Error())
		t.FailNow()
	}
	fmt.Println(string(response.Payload))
}
