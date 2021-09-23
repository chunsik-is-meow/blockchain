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
	Type             string   `json:"type"`
	Uploader         string   `json:"uploader"`
	Name             string   `json:"name"`
	Version          string   `json:"version"`
	Language         string   `json:"language"`
	Price            uint32   `json:"price"`
	Owner            string   `json:"owner"`
	Score            uint64   `json:"score"`
	Downloaded       uint32   `json:"downloaded"`
	Description      string   `json:"description"`
	VerificationOrgs []string `json:"verification_orgs"`
	Contents         string   `json:"contents`
	Timestamp        string   `json:"timestamp"`
}

type AIModelCount struct {
	Type  string `json:"type"`
	List  string `json:"list"`
	Count int    `json:"count"`
}

// InitLedger ...
func (a *AIChaincode) InitLedger(ctx contractapi.TransactionContextInterface) error {
	isInitBytes, err := ctx.GetStub().GetState("isInit")
	if err != nil {
		return fmt.Errorf("failed GetState('isInit')")
	} else if isInitBytes == nil {
		initCount := AIModelCount{
			Type:  "AIModelCount",
			List:  "",
			Count: 0,
		}
		initAIModelCountAsBytes, err := json.Marshal(initCount)
		ctx.GetStub().PutState(makeAIModelCountKey("AC"), initAIModelCountAsBytes)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
		if err != nil {
			return fmt.Errorf("failed to json.Marshal(). %v", err)
		}
		return nil
	} else {
		return fmt.Errorf("already initialized")
	}
}

func (a *AIChaincode) PutAIModel(ctx contractapi.TransactionContextInterface, uploader string, name string, version string, language string, price uint32, owner string, description string, contents string, timestamp string, score uint64) error {
	exists, err := a.aiModelExists(ctx, uploader, name, version)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the aiModel %s already exists", name)
	}
	var download uint32
	download = 0
	verificationOrgs := []string{"verification-01"}
	aiModelInfo := AIModelType{
		Type:             "AI-Model",
		Uploader:         uploader,
		Name:             name,
		Version:          version,
		Language:         language,
		Price:            price,
		Owner:            owner,
		Score:            score,
		Downloaded:       download,
		Description:      description,
		VerificationOrgs: verificationOrgs,
		Contents:         contents,
		Timestamp:        timestamp,
	}
	aiModelAsBytes, err := json.Marshal(aiModelInfo)
	if err != nil {
		return fmt.Errorf("failed to json.Marshal(). %v", err)
	}
	aiModelKey := makeAIModelKey(uploader, name, version)
	ctx.GetStub().PutState(aiModelKey, aiModelAsBytes)
	if err != nil {
		return fmt.Errorf("failed to put to world state. %v", err)
	}

	currentAIModelCount, err := a.GetAIModelCount(ctx, "AC")
	if err != nil {
		return fmt.Errorf("failed to get count. %v", err)
	}
	currentAIModelCount.Count++
	currentAIModelCount.List += "/" + name

	currentAIModelCountAsBytes, err := json.Marshal(currentAIModelCount)
	if err != nil {
		return fmt.Errorf("failed to json.Marshal(). %v", err)
	}
	ctx.GetStub().PutState(makeAIModelCountKey("AC"), currentAIModelCountAsBytes)
	return nil
}

func makeAIModelCountKey(key string) string {
	var sb strings.Builder

	sb.WriteString("Count_D_")
	sb.WriteString(key)
	return sb.String()
}

func makeAIModelKey(uploader string, name string, version string) string {
	var sb strings.Builder

	sb.WriteString("A_")
	sb.WriteString(uploader)
	sb.WriteString("_")
	sb.WriteString(name)
	sb.WriteString("_")
	sb.WriteString(version)
	return sb.String()
}

func (a *AIChaincode) aiModelExists(ctx contractapi.TransactionContextInterface, uploader string, name string, version string) (bool, error) {
	aiModelAsBytes, err := ctx.GetStub().GetState(makeAIModelKey(uploader, name, version))
	if err != nil {
		return false, fmt.Errorf("ai-model is exist...: %v", err)
	}

	return aiModelAsBytes != nil, nil
}

func (a *AIChaincode) GetAllAIModelInfo(ctx contractapi.TransactionContextInterface) ([]*AIModelType, error) {
	queryString := fmt.Sprintf(`{"selector":{"type":"AI-Model"}}`)
	return getQueryResultForQueryString(ctx, queryString)
}

func (a *AIChaincode) GetAIModelInfo(ctx contractapi.TransactionContextInterface, uploader string, name string, version string) (*AIModelType, error) {
	aiModelInfo := &AIModelType{}
	aiModelAsBytes, err := ctx.GetStub().GetState(makeAIModelKey(uploader, name, version))
	if err != nil {
		return nil, err
	} else if aiModelAsBytes == nil {
		aiModelInfo.Type = "empty"
		aiModelInfo.Uploader = "empty"
		aiModelInfo.Name = "empty"
		aiModelInfo.Version = "empty"
		aiModelInfo.Language = "empty"
		aiModelInfo.Price = 0
		aiModelInfo.Owner = "empty"
		aiModelInfo.Score = 0
		aiModelInfo.Downloaded = 0
		aiModelInfo.Description = "empty"
		aiModelInfo.VerificationOrgs = []string{""}
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

func (a *AIChaincode) GetAIModelInfoWithKey(ctx contractapi.TransactionContextInterface, aiModelKey string) (*AIModelType, error) {
	aiModelInfo := &AIModelType{}
	aiModelAsBytes, err := ctx.GetStub().GetState(aiModelKey)
	if err != nil {
		return nil, err
	} else if aiModelAsBytes == nil {
		aiModelInfo.Type = "empty"
		aiModelInfo.Uploader = "empty"
		aiModelInfo.Name = "empty"
		aiModelInfo.Version = "empty"
		aiModelInfo.Language = "empty"
		aiModelInfo.Price = 0
		aiModelInfo.Owner = "empty"
		aiModelInfo.Score = 0
		aiModelInfo.Downloaded = 0
		aiModelInfo.Description = "empty"
		aiModelInfo.VerificationOrgs = []string{""}
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

func (a *AIChaincode) GetAIModelContents(ctx contractapi.TransactionContextInterface, uploader string, name string, version string, downloader string) (string, error) {
	aiModelInfo, err := a.GetAIModelInfo(ctx, uploader, name, version)
	if err != nil {
		return "failed to get Info", err
	}
	aiModelInfo.Downloaded++
	aiModelAsBytes, err := json.Marshal(aiModelInfo)
	if err != nil {
		return "failed to json.Marshal().", err
	}
	aiModelKey := makeAIModelKey(uploader, name, version)
	ctx.GetStub().PutState(aiModelKey, aiModelAsBytes)
	if err != nil {
		return "failed to put to world state.", err
	}

	currentAIModelCount, err := a.GetAIModelCount(ctx, downloader)
	if err != nil {
		return "failed to get count", err
	}
	currentAIModelCount.Count++
	currentAIModelCount.List += "/" + name

	currentAIModelCountAsBytes, err := json.Marshal(currentAIModelCount)
	if err != nil {
		return "failed to json.Marshal()", err
	}
	ctx.GetStub().PutState(makeAIModelCountKey(downloader), currentAIModelCountAsBytes)

	return aiModelInfo.Contents, nil
}

func (a *AIChaincode) GetAIModelCount(ctx contractapi.TransactionContextInterface, key string) (*AIModelCount, error) {
	currentAIModelCount := &AIModelCount{}
	currentAIModelCountAsBytes, err := ctx.GetStub().GetState(makeAIModelCountKey(key))
	if err != nil {
		return nil, err
	} else if currentAIModelCountAsBytes == nil {
		currentAIModelCount.Type = "AIModelCount"
		currentAIModelCount.List = ""
		currentAIModelCount.Count = 0
	} else {
		err = json.Unmarshal(currentAIModelCountAsBytes, currentAIModelCount)
		if err != nil {
			return nil, err
		}
	}

	return currentAIModelCount, nil
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
