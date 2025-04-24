//go:build gofuzz
// +build gofuzz

//sasass

package fuzz

import (
	"bytes"
	"sync"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"

	cc "github.com/hyperledger/fabric-samples/asset-transfer-basic/chaincode-go/chaincode"
)

var (
	stub *shimtest.MockStub
	once sync.Once
	mu   sync.Mutex
)

func initStub() {
	
	cc, err := contractapi.NewChaincode(new(cc.SmartContract))
	if err != nil {
		panic(err)
	}
	stub = shimtest.NewMockStub("fuzzstub", cc)
	res := stub.MockInit("tx0", nil)
	if res.Status != shim.OK {
		panic(res.Message)
	}
}

func Fuzz(data []byte) int {
	once.Do(initStub)

	args := bytes.Split(data, []byte{0})
	if len(args) == 0 || len(args[0]) == 0 {
		return 0
	}

	mu.Lock()
	_ = stub.MockInvoke("tx", args)
	mu.Unlock()

	return 1
}
