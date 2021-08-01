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

// MeowType ...
type MeowType struct {
	Type   string `json:"type"`
	Amount uint32 `json:"amount"`
}

// RewardType ...
type RewardType struct {
	Type   string `json:"type"`
	To     string `json:"to"`
	Amount uint32 `json:"amount"`
}

//BuyAIModelType ...
type BuyAIModelType struct {
	Timestamp string       `json:"timestamp"`
	History   []RewardType `json:"history"`
}

// ModelResult ...
type ModelResult struct {
	F1Score string `json:"f1_score"`
}

// Model ...
type Model struct {
	VerificationOrgs []string    `json:"verification_orgs"`
	Result           ModelResult `json:"model_result"`
	Price            uint32      `json:"price"`
}

// InitLedger ...
func (t *TradeChaincode) InitLedger(ctx contractapi.TransactionContextInterface) error {
	initMeow := MeowType{
		Type:   "GenesisMint",
		Amount: GENESIS_MINT_AMOUNT,
	}

	initMeowAsBytes, err := json.Marshal(initMeow)
	if err != nil {
		return fmt.Errorf("failed to json.Marshal(). %v", err)
	}

	ctx.GetStub().PutState(makeMeowKey("bank"), initMeowAsBytes)
	if err != nil {
		return fmt.Errorf("failed to put to world state. %v", err)
	}

	return nil
}

// GetCurrentMeow ...
func (t *TradeChaincode) GetCurrentMeow(ctx contractapi.TransactionContextInterface, uid string) (*MeowType, error) {
	currentMeow := &MeowType{}
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

// Transfer ...
func (t *TradeChaincode) Transfer(ctx contractapi.TransactionContextInterface, from string, to string, amount uint32, timestamp string, meowType string) error {
	// INSERT Transfer history
	rewardMeow := MeowType{
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

// BuyModel ...
func (t *TradeChaincode) BuyModel(ctx contractapi.TransactionContextInterface, uid string, modelKey string, price uint32, timestamp string) error {
	checkBuyAIModelAsBytes, err := ctx.GetStub().GetState(makeBuyAIModelKey(uid, modelKey))
	if err != nil {
		return fmt.Errorf("failed to get BuyAiModel. %v", err)
	} else if checkBuyAIModelAsBytes != nil {
		return fmt.Errorf("already buy model ...")
	}

	// TODO
	// GetModel from ai-model-channel
	// check isNotExist Model Error
	// modelKey -> ai-model channel query getModel(modelKey) ->
	// Model {
	// 	v orgs
	// 	result
	//  price
	// }

	model := Model{
		VerificationOrgs: []string{"verification-01"},
		Result: ModelResult{
			F1Score: "93.4%",
		},
		Price: 3000,
	}

	if price != model.Price {
		return fmt.Errorf("the price mismatch in blockchain ..")
	}

	currentMeow, err := t.GetCurrentMeow(ctx, uid)
	if err != nil {
		return fmt.Errorf("failed to get current meow. %v", err)
	}

	if currentMeow.Amount < price {
		return fmt.Errorf("meow is lacking.. %v", err)
	}

	// NOTE
	// modelKey -> AI_uid_modelName_version(unique)
	seller := strings.Split(modelKey, "_")[1]
	verificationOrgs := model.VerificationOrgs

	if price%10 != 0 {
		return fmt.Errorf("only available in units of 10 meow")
	}

	income := price * 8 / 10
	verifyReward := price * 1 / 10
	manageReward := price * 1 / 10

	t.Transfer(ctx, uid, seller, income, timestamp, "income")
	t.Transfer(ctx, uid, "admin", manageReward, timestamp, "reward")

	buyAIModel := BuyAIModelType{
		Timestamp: timestamp,
		History: []RewardType{
			{
				Type:   "income",
				To:     seller,
				Amount: income,
			},
			{
				Type:   "manageReward",
				To:     "admin",
				Amount: manageReward,
			},
		},
	}

	for _, org := range verificationOrgs {
		// TODO
		// divide verifyReward
		t.Transfer(ctx, uid, org, verifyReward, timestamp, "reward")
		verifyRewardHistory := RewardType{
			Type:   "verifyReward",
			To:     org,
			Amount: verifyReward,
		}
		buyAIModel.History = append(buyAIModel.History, verifyRewardHistory)
	}

	buyAIModelAsBytes, err := json.Marshal(buyAIModel)
	if err != nil {
		return fmt.Errorf("failed to json.Marshal(). %v", err)
	}
	ctx.GetStub().PutState(makeBuyAIModelKey(uid, modelKey), buyAIModelAsBytes)

	return nil
}

// TODO
// func getIsBuyModel
// func getAsset

func (t *TradeChaincode) getAllHistory(uid string) error {
	return nil
}

func (t *TradeChaincode) getQueryHistory(uid string) error {
	return nil
}

func makeBuyAIModelKey(uid string, model string) string {
	var sb strings.Builder

	sb.WriteString("B_")
	sb.WriteString(uid)
	sb.WriteString("_")
	sb.WriteString(model)

	return sb.String()
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
