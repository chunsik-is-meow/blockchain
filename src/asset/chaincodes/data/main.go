package main

import (
	"log"

	"github.com/chunsik-is-meow/blockchain/src/asset/chaincodes/data/contract"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	dataChaincode, err := contractapi.NewChaincode(&contract.DataChaincode{})
	if err != nil {
		log.Panicf("Error creating dataChaincode: %v", err)
	}

	if err := dataChaincode.Start(); err != nil {
		log.Panicf("Error starting dataChaincode: %v", err)
	}
}
