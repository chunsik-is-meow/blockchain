package main

import (
	"log"

	"github.com/chunsik-is-meow/blockchain/src/asset/chaincodes/ai/contract"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	dataChaincode, err := contractapi.NewChaincode(&contract.AIChaincode{})
	if err != nil {
		log.Panicf("Error creating aiChaincode: %v", err)
	}

	if err := dataChaincode.Start(); err != nil {
		log.Panicf("Error starting aiChaincode: %v", err)
	}
}
