package main

import (
	"log"

	"github.com/chunsik-is-meow/blockchain/src/asset/chaincodes/ai-model/contract"
	"github.com/chunsik-is-meow/blockchain/src/asset/chaincodes/data/contract"
	"github.com/chunsik-is-meow/blockchain/src/asset/chaincodes/trade/contract"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	tradeChaincode, err := contractapi.NewChaincode(&contract.TradeChaincode{})
	if err != nil {
		log.Panicf("Error creating tradeChaincode: %v", err)
	}

	if err := tradeChaincode.Start(); err != nil {
		log.Panicf("Error starting tradeChaincode: %v", err)
	}
	// dataChaincode, err := contractapi.NewChaincode(&contract.DataChaincode{})
	// if err != nil {
	// 	log.Panicf("Error creating dataChaincode: %v", err)
	// }

	// if err := dataChaincode.Start(); err != nil {
	// 	log.Panicf("Error starting dataChaincode: %v", err)
	// }
	// aiChaincode, err := contractapi.NewChaincode(&contract.AIChaincode{})
	// if err != nil {
	// 	log.Panicf("Error creating aiChaincode: %v", err)
	// }

	// if err := aiChaincode.Start(); err != nil {
	// 	log.Panicf("Error starting aiChaincode: %v", err)
	// }
}
