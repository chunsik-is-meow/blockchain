package contract

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const GENESIS_MINT_AMOUNT = 100000000

// TradeChaincode ...
type TradeChaincode struct {
	contractapi.Contract
}

// Meow ...
type Meow struct {
	Type   string `json:"type"`
	Amount uint32 `json:"amount"`
}

// InitLedger ...
func (t *TradeChaincode) InitLedger(ctx contractapi.TransactionContextInterface) error {
	initMeow := Meow{
		Type:   "GenesisMint",
		Amount: GENESIS_MINT_AMOUNT,
	}

	initMeowAsBytes, err := json.Marshal(initMeow)
	if err != nil {
		return fmt.Errorf("failed to json.Marshal(). %v", err)
	}

	ctx.GetStub().PutState(makeMeowKey("admin"), initMeowAsBytes)
	if err != nil {
		return fmt.Errorf("failed to put to world state. %v", err)
	}

	return nil
}

// GetCurrentMeow ...
func (t *TradeChaincode) GetCurrentMeow(ctx contractapi.TransactionContextInterface, uid string) (*Meow, error) {
	currentMeow := &Meow{}
	currentMeowAsBytes, err := ctx.GetStub().GetState(makeMeowKey(uid))
	if err != nil {
		return nil, err
	} else if currentMeowAsBytes == nil {
		currentMeow.Type = "CurrentMeowAmount"
		currentMeow.Amount = 0
	} else {
		err = json.Unmarshal(currentMeowAsBytes, currentMeow)
		if err != nil {
			return nil, err
		}
	}

	return currentMeow, nil
}

// Reward ...
func (t *TradeChaincode) Reward(ctx contractapi.TransactionContextInterface, uid string, amount uint32, timestamp string) error {
	return t.exec(ctx, "admin", uid, amount, timestamp, "Reward")
}

// Transfer ...
func (t *TradeChaincode) Transfer(ctx contractapi.TransactionContextInterface, from string, to string, amount uint32, timestamp string) error {
	return t.exec(ctx, from, to, amount, timestamp, "Transfer")
}

func (t *TradeChaincode) exec(ctx contractapi.TransactionContextInterface, from string, to string, amount uint32, timestamp string, meowType string) error {
	// INSERT Transfer history
	rewardMeow := Meow{
		Type:   meowType,
		Amount: amount,
	}

	rewardMeowAsBytes, err := json.Marshal(rewardMeow)
	if err != nil {
		return fmt.Errorf("failed to json.Marshal(). %v", err)
	}
	rewardMeowKey := makeFromToMeowKey(from, to, timestamp)
	ctx.GetStub().PutState(rewardMeowKey, rewardMeowAsBytes)

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

func makeMeowKey(uid string) string {
	var sb strings.Builder

	sb.WriteString("M_")
	sb.WriteString(uid)

	return sb.String()
}

func makeFromToMeowKey(from string, to string, timestamp string) string {
	var sb strings.Builder

	sb.WriteString("F_")
	sb.WriteString(from)
	sb.WriteString("_T_")
	sb.WriteString(to)
	sb.WriteString("_")
	sb.WriteString(timestamp)

	return sb.String()
}
