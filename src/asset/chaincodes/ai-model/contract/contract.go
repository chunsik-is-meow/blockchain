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
	Score       int    `json:"score"`
	Description string `json:"description"`
	Timestamp   string `json:"timestamp"`
}

// InitLedger ...
func (a *AIChaincode) InitLedger(ctx contractapi.TransactionContextInterface) error {
	aiModelInfos := []AIModelType{
		{Type: "ai-model", Name: "adult_learning_model", Language: "python", Price: 2400, Owner: "AAA", Score: 75, Description: "adult_learning", Timestamp: "2021-08-27-09-11-49"},
		{Type: "ai-model", Name: "cancer_learning_model", Language: "go", Price: 3100, Owner: "BBB", Score: 82, Description: "cancer_learning", Timestamp: "2021-08-27-09-11-49"},
	}

	isInitBytes, err := ctx.GetStub().GetState("isInit")
	if err != nil {
		return fmt.Errorf("failed GetState('isInit')")
	} else if isInitBytes == nil {
		for _, aiModel := range aiModelInfos {
			aiModelAsBytes, err := json.Marshal(aiModel)
			if err != nil {
				return fmt.Errorf("failed to json.Marshal(). %v", err)
			}

			aiModelKey := makeAIModelKey(aiModel.Name)
			ctx.GetStub().PutState(aiModelKey, aiModelAsBytes)
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
func (a *AIChaincode) AIModelInsert(ctx contractapi.TransactionContextInterface, name string, description string, owner string, timestamp string) error {
	// TODO
	// aiModel file upload
	return nil
}

func (a *AIChaincode) PutAIModel(ctx contractapi.TransactionContextInterface, name string, language string, price int, owner string, description string, timestamp string) error {
	// a.AIModelInsert(ctx, name, description, owner, timestamp)
	exists, err := a.aiModelExists(ctx, name)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the aiModel %s already exists", name)
	}
	file := ""
	score, err := evaluateScore(ctx, file)
	if err != nil {
		return err
	}

	aiModel := AIModelType{
		Type:        "ai-model",
		Name:        name,
		Language:    language,
		Price:       price,
		Owner:       owner,
		Score:       score,
		Description: description,
		Timestamp:   timestamp,
	}
	aiModelAsBytes, err := json.Marshal(aiModel)
	if err != nil {
		return fmt.Errorf("failed to json.Marshal(). %v", err)
	}
	aiModelKey := makeAIModelKey(name)
	ctx.GetStub().PutState(aiModelKey, aiModelAsBytes)
	if err != nil {
		return fmt.Errorf("failed to put to world state. %v", err)
	}

	return nil
}

func (a *AIChaincode) GetAllAIModelInfo(ctx contractapi.TransactionContextInterface) ([]*AIModelType, error) {
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
		aiModelInfo.Score = 0
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

func makeAIModelKey(key string) string {
	var sb strings.Builder

	sb.WriteString("A_")
	sb.WriteString(key)
	return sb.String()
}

func (a *AIChaincode) aiModelExists(ctx contractapi.TransactionContextInterface, name string) (bool, error) {
	aiModelAsBytes, err := ctx.GetStub().GetState(name)
	if err != nil {
		return false, fmt.Errorf("ai-model is exist...: %v", err)
	}

	return aiModelAsBytes != nil, nil
}

func (a *AIChaincode) GetQueryAIModelHistory(ctx contractapi.TransactionContextInterface) ([]*AIModelType, error) {
	queryString := fmt.Sprintf(`{"selector":{"type":"ai-model"}}`)
	return getQueryResultForQueryString(ctx, queryString)
}

func evaluateScore(ctx contractapi.TransactionContextInterface, aiModel string) (int, error) {
	// TODO
	score := 81
	return score, nil
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
