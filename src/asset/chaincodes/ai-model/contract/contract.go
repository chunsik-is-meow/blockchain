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
	Type        string `json:"type"`
	Name        string `json:"name"`
	Language    string `json:"language"`
	Price       int    `json:"price"`
	Owner       string `json:"owner"`
	Description string `json:"description"`
	Timestamp   string `json:"timestamp"`
}

// InitLedger ...
func (a *AIChaincode) InitLedger(ctx contractapi.TransactionContextInterface) error {
	aiModelInfos := []AIModelType{
		{Type: "aiModel", Name: "adult learning model", Language: "python", Price: 2400, Owner: "AAA", Description: "adult learning", Timestamp: "2021-08-27-09-11-49"},
		{Type: "aiModel", Name: "cancer learning model", Language: "go", Price: 3100, Owner: "BBB", Description: "cancer learning", Timestamp: "2021-08-27-09-11-49"},
	}

	isInitBytes, err := ctx.GetStub().GetState("isInit")
	if err != nil {
		return fmt.Errorf("failed GetState('isInit')")
	} else if isInitBytes == nil {
		for _, aiModel := range aiModelInfos {
			assetJSON, err := json.Marshal(aiModel)
			if err != nil {
				return fmt.Errorf("failed to json.Marshal(). %v", err)
			}

			err = ctx.GetStub().PutState(aiModel.Name, assetJSON)
			if err != nil {
				return fmt.Errorf("failed to put to world state. %v", err)
			}
		}

		return nil
	} else {
		return fmt.Errorf("already initialized")
	}
}

// AIModelInsert ...
func (a *AIChaincode) AIModelInsert(ctx contractapi.TransactionContextInterface, name string, language string, price int, owner string, description string, timestamp string) error {
	// TODO
	return nil
}

func (a *AIChaincode) PutAIModel(ctx contractapi.TransactionContextInterface, name string, language string, price int, owner string, description string, timestamp string) error {
	// a.AIModelInsert(ctx, name, language, price, owner, description, timestamp)
	exists, err := a.aiModelExists(ctx, name)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the aiModel %s already exists", name)
	}
	aiModel := AIModelType{
		Type:        "aiModel",
		Name:        name,
		Language:    language,
		Price:       price,
		Owner:       owner,
		Description: description,
		Timestamp:   timestamp,
	}
	aiModelAsBytes, err := json.Marshal(aiModel)
	if err != nil {
		return fmt.Errorf("failed to json.Marshal(). %v", err)
	}
	aiModelKey := makeAIModelKey(owner)
	ctx.GetStub().PutState(aiModelKey, aiModelAsBytes)
	if err != nil {
		return fmt.Errorf("failed to put to world state. %v", err)
	}
	return nil
}

func (a *AIChaincode) GetAllAIModelInfo(ctx contractapi.TransactionContextInterface) ([]*AIModelType, error) {
	// TODO
	var aiModelInfos []*AIModelType
	aiModelsAsBytes, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}

	for aiModelsAsBytes.HasNext() {
		queryResponse, err := aiModelsAsBytes.Next()
		if err != nil {
			return nil, err
		}

		var aiModel AIModelType
		err = json.Unmarshal(queryResponse.Value, &aiModel)
		if err != nil {
			return nil, err
		}
		aiModelInfos = append(aiModelInfos, &aiModel)
	}

	return aiModelInfos, nil
}

func (a *AIChaincode) GetAIModelInfo(ctx contractapi.TransactionContextInterface, name string) (*AIModelType, error) {
	// TODO
	aiModelInfo := &AIModelType{}
	aiModelAsBytes, err := ctx.GetStub().GetState(makeAIModelKey(name))
	if err != nil {
		return nil, err
	} else if aiModelAsBytes == nil {
		aiModelInfo.Type = "empty"
		aiModelInfo.Name = "empty"
		aiModelInfo.Language = "empty"
		aiModelInfo.Price = 0
		aiModelInfo.Owner = "empty"
		aiModelInfo.Description = "empty"
		aiModelInfo.Language = "empty"
		aiModelInfo.Timestamp = "empty"
	} else {
		err = json.Unmarshal(aiModelAsBytes, aiModelInfo)
		if err != nil {
			return nil, err
		}
	}

	return aiModelInfo, nil
}

func makeAIModelKey(classfication string) string {
	var sb strings.Builder

	sb.WriteString("D_")
	sb.WriteString(classfication)
	return sb.String()
}

func (a *AIChaincode) aiModelExists(ctx contractapi.TransactionContextInterface, name string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(name)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

func (a *AIChaincode) GetQueryAIModelHistory(ctx contractapi.TransactionContextInterface) ([]*AIModelType, error) {
	queryString := fmt.Sprintf(`{"selector":{"type":"aiModel"}}`)
	return getQueryResultForQueryString(ctx, queryString)
}

func getQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*AIModelType, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var transferHistorys []*AIModelType
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var transferHistory AIModelType
		err = json.Unmarshal(queryResult.Value, &transferHistory)
		if err != nil {
			return nil, err
		}
		transferHistorys = append(transferHistorys, &transferHistory)
	}

	return transferHistorys, nil
}
