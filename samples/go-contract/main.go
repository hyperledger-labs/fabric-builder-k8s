// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SampleContract struct {
	contractapi.Contract
}

// PutValue - Adds a key value pair to the world state
func (sc *SampleContract) PutValue(
	ctx contractapi.TransactionContextInterface,
	key string,
	value string,
) error {
	return ctx.GetStub().PutState(key, []byte(value))
}

// GetValue - Gets the value for a key from the world state
func (sc *SampleContract) GetValue(
	ctx contractapi.TransactionContextInterface,
	key string,
) (string, error) {
	bytes, err := ctx.GetStub().GetState(key)

	if err != nil {
		return "", nil
	}

	return string(bytes), nil
}

func main() {
	SampleContract := new(SampleContract)

	cc, err := contractapi.NewChaincode(SampleContract)

	if err != nil {
		panic(err.Error())
	}

	if err := cc.Start(); err != nil {
		panic(err.Error())
	}
}
