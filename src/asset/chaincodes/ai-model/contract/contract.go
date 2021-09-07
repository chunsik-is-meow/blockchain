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
	Contents    string `json:"contents`
	Timestamp   string `json:"timestamp"`
}

type DataCount struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
}

// InitLedger ...
func (a *AIChaincode) InitLedger(ctx contractapi.TransactionContextInterface) error {
	aiModelInfos := []AIModelType{
		{Type: "ai-model", Name: "test", Language: "test", Price: 0, Owner: "test", Score: 100, Description: "test", Contents: "test", Timestamp: "2021-08-27-09-11-49"},
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

			aiModelKey := makeAIModelKey("admin", aiModel.Name, "0.0")
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

func (a *AIChaincode) PutAIModel(ctx contractapi.TransactionContextInterface, username string, name string, version string, language string, price int, owner string, description string, contents string, timestamp string) error {
	// a.AIModelInsert(ctx, name, description, owner, timestamp)
	exists, err := a.aiModelExists(ctx, username, name, version)
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
		Contents:    contents,
		Timestamp:   timestamp,
	}
	aiModelAsBytes, err := json.Marshal(aiModel)
	if err != nil {
		return fmt.Errorf("failed to json.Marshal(). %v", err)
	}
	aiModelKey := makeAIModelKey(username, name, version)
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

func (a *AIChaincode) GetAIModelInfo(ctx contractapi.TransactionContextInterface, username string, name string, version string) (*AIModelType, error) {
	aiModelInfo := &AIModelType{}
	aiModelAsBytes, err := ctx.GetStub().GetState(makeAIModelKey(username, name, version))
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
		aiModelInfo.Contents = "empty"
		aiModelInfo.Timestamp = "empty"
	} else {
		err = json.Unmarshal(aiModelAsBytes, &aiModelInfo)
		if err != nil {
			return nil, err
		}
	}
	return aiModelInfo, nil
}

func (a *AIChaincode) GetAIModelContents(ctx contractapi.TransactionContextInterface, username string, name string, version string) (string, error) {
	var aiModelInfo AIModelType
	aiModelAsBytes, err := ctx.GetStub().GetState(makeAIModelKey(username, name, version))
	if err != nil {
		return "not existed...", err
	} else if aiModelAsBytes == nil {
		aiModelInfo.Type = "ai-model"
		aiModelInfo.Name = "empty"
		aiModelInfo.Language = "empty"
		aiModelInfo.Price = 0
		aiModelInfo.Owner = "empty"
		aiModelInfo.Score = 0
		aiModelInfo.Description = "empty"
		aiModelInfo.Contents = "empty"
		aiModelInfo.Timestamp = "empty"
	} else {
		err = json.Unmarshal(aiModelAsBytes, &aiModelInfo)
		if err != nil {
			return "failed...", err
		}
	}
	return aiModelInfo.Contents, nil
}

func GetAllAIModelCount(ctx contractapi.TransactionContextInterface) string {
	return "aa"
}

func makeAIModelKey(username string, name string, version string) string {
	var sb strings.Builder

	sb.WriteString("A_")
	sb.WriteString(username)
	sb.WriteString("_")
	sb.WriteString(name)
	sb.WriteString("_")
	sb.WriteString(version)
	return sb.String()
}

func (a *AIChaincode) aiModelExists(ctx contractapi.TransactionContextInterface, username string, name string, version string) (bool, error) {
	aiModelAsBytes, err := ctx.GetStub().GetState(makeAIModelKey(username, name, version))
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
