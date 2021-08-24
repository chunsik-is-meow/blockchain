package contract

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// AIChaincode ...
type AIChaincode struct {
	contractapi.Contract
}

// AIModelType ...
type AIModelType struct {
	ObjectType         string `json:"docType"`
	Name       string `json:"name"`
	Score       int    `json:"score"`
	OwnelID       string `json:"ownelID"`
	SrcDataID       string `json:"srcDataID"`
	Description      string `json:"description"`
}


func (t *DataChaincode) PutAIModel(ctx contractapi.TransactionContextInterface, key string, model string, timestamp string) error {
	return nil
}


func (t *DataChaincode) GetAIModel(ctx contractapi.TransactionContextInterface, key string, timestamp string) error {
	return nil
}


func (t *DataChaincode) checkScore(ctx contractapi.TransactionContextInterface, model string) error {
	return nil
}


func (t *DataChaincode) isVaildData(ctx contractapi.TransactionContextInterface, srcDataID string) error {
	return nil
}




// Transfer ...
func (t *AIChaincode) Transfer(ctx contractapi.TransactionContextInterface, from string, to string, amount uint32, timestamp string) error {
	// INSERT Transfer history
	transferMeow := AIModelType{
		Type:   "transfer",
		From:   from,
		To:     to,
		Amount: amount,
	}

	transferMeowAsBytes, err := json.Marshal(transferMeow)
	if err != nil {
		return fmt.Errorf("failed to json.Marshal(). %v", err)
	}
	transferMeowKey := makeFromToMeowKey(from, to, timestamp)
	ctx.GetStub().PutState(transferMeowKey, transferMeowAsBytes)

	if err != nil {
		return fmt.Errorf("failed to put to world state. %v", err)
	}

	// UPDATE Current From Meow
	currentFromMeow, err := t.GetCurrentMeow(ctx, from)
	if err != nil {
		return fmt.Errorf("failed to get current meow. %v", err)
	}

	if currentFromMeow.Amount < amount {
		return fmt.Errorf("meow is lacking.. %v", err)
	}

	currentFromMeow.Amount -= amount

	currentFromMeowAsBytes, err := json.Marshal(currentFromMeow)
	if err != nil {
		return fmt.Errorf("failed to json.Marshal(). %v", err)
	}
	ctx.GetStub().PutState(makeMeowKey(from), currentFromMeowAsBytes)

	// UPDATE Current To Meow
	currentToMeow, err := t.GetCurrentMeow(ctx, to)
	if err != nil {
		return fmt.Errorf("failed to get current meow. %v", err)
	}

	currentToMeow.Amount += amount

	currentToMeowAsBytes, err := json.Marshal(currentToMeow)
	if err != nil {
		return fmt.Errorf("failed to json.Marshal(). %v", err)
	}
	ctx.GetStub().PutState(makeMeowKey(to), currentToMeowAsBytes)

	// TODO
	// Transfer
	// Before amount (from, to)
	// After amount (from, to)
	return nil
}

