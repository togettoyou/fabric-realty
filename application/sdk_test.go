package main_test

import (
	"fmt"
	"github.com/togettoyou/blockchain-real-estate/application/blockchain"
	"testing"
)

func TestInvoke_QueryAccountList(t *testing.T) {
	blockchain.Init()
	response, e := blockchain.ChannelExecute("createRealEstate", [][]byte{
		[]byte("5feceb66ffc8"),
		[]byte("6b86b273ff34"),
		[]byte("122.22"),
		[]byte("122.22"),
	})
	if e != nil {
		fmt.Println(e.Error())
		t.FailNow()
	}
	fmt.Println(string(response.Payload))
}
