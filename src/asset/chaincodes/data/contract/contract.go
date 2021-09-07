package contract

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// DataChaincode ...
type DataChaincode struct {
	contractapi.Contract
}

// DataType ...
type DataType struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Downloaded  int    `json:"downloaded"`
	Owner       string `json:"owner"`
	Contents    string `json:"contents`
	Timestamp   string `json:"timestamp"`
}

type DataCount struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
}

// InitLedger ...
func (d *DataChaincode) InitLedger(ctx contractapi.TransactionContextInterface) error {
	dataInfos := []DataType{
		{Type: "data", Name: "test", Description: "test", Downloaded: 0, Owner: "test", Contents: "test", Timestamp: "2021-08-27-09-11-49"},
	}

	isInitBytes, err := ctx.GetStub().GetState("isInit")
	if err != nil {
		return fmt.Errorf("failed GetState('isInit')")
	} else if isInitBytes == nil {
		for _, data := range dataInfos {
			dataAsBytes, err := json.Marshal(data)
			if err != nil {
				return fmt.Errorf("failed to json.Marshal(). %v", err)
			}

			dataKey := makeDataKey("admin", data.Name, "0.0")
			ctx.GetStub().PutState(dataKey, dataAsBytes)
			if err != nil {
				return fmt.Errorf("failed to put to world state. %v", err)
			}
		}

		return nil
	} else {
		return fmt.Errorf("already initialized")
	}
}

// DataInsert ...
func (d *DataChaincode) PutCommonData(ctx contractapi.TransactionContextInterface, username string, name string, version string, description string, owner string, contents string, timestamp string) error {
	exists, err := d.dataExists(ctx, username, name, version)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the data %s already exists", name)
	}

	download := 0
	data := DataType{
		Name:        name,
		Type:        "data",
		Description: description,
		Downloaded:  download,
		Owner:       owner,
		Contents:    contents,
		Timestamp:   timestamp,
	}
	dataAsBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to json.Marshal(). %v", err)
	}
	dataKey := makeDataKey(username, name, version)
	ctx.GetStub().PutState(dataKey, dataAsBytes)
	if err != nil {
		return fmt.Errorf("failed to put to world state. %v", err)
	}

	return nil
}

func (d *DataChaincode) GetAllCommonDataInfo(ctx contractapi.TransactionContextInterface) ([]*DataType, error) {
	datasAsBytes, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}

	var dataInfos []*DataType
	for datasAsBytes.HasNext() {
		queryResponse, err := datasAsBytes.Next()
		if err != nil {
			return nil, err
		}

		var data DataType
		err = json.Unmarshal(queryResponse.Value, &data)
		if err != nil {
			return nil, err
		}
		dataInfos = append(dataInfos, &data)
	}

	return dataInfos, nil
}

func (d *DataChaincode) GetCommonDataInfo(ctx contractapi.TransactionContextInterface, username string, name string, version string) (*DataType, error) {
	var dataInfo DataType
	dataAsBytes, err := ctx.GetStub().GetState(makeDataKey(username, name, version))
	if err != nil {
		return nil, err
	} else if dataAsBytes == nil {
		dataInfo.Name = "empty"
		dataInfo.Type = "empty"
		dataInfo.Description = "empty"
		dataInfo.Downloaded = 0
		dataInfo.Owner = "empty"
		dataInfo.Contents = "empty"
		dataInfo.Timestamp = "empty"
	} else {
		err = json.Unmarshal(dataAsBytes, &dataInfo)
		if err != nil {
			return nil, err
		}
	}
	return &dataInfo, nil
}

func (d *DataChaincode) GetCommonDataContents(ctx contractapi.TransactionContextInterface, username string, name string, version string) (string, error) {
	var dataInfo DataType
	dataAsBytes, err := ctx.GetStub().GetState(makeDataKey(username, name, version))
	if err != nil {
		return "not existed...", err
	} else if dataAsBytes == nil {
		dataInfo.Name = "empty"
		dataInfo.Type = "empty"
		dataInfo.Description = "empty"
		dataInfo.Downloaded = 0
		dataInfo.Owner = "empty"
		dataInfo.Contents = "empty"
		dataInfo.Timestamp = "empty"
	} else {
		err = json.Unmarshal(dataAsBytes, &dataInfo)
		if err != nil {
			return "failed...", err
		}
	}
	return dataInfo.Contents, nil
}

func GetAllDataCount(ctx contractapi.TransactionContextInterface) string {
	return "aa"
}

func makeDataKey(username string, name string, version string) string {
	var sb strings.Builder

	sb.WriteString("D_")
	sb.WriteString(username)
	sb.WriteString("_")
	sb.WriteString(name)
	sb.WriteString("_")
	sb.WriteString(version)
	return sb.String()
}

func (d *DataChaincode) dataExists(ctx contractapi.TransactionContextInterface, username string, name string, version string) (bool, error) {
	dataAsBytes, err := ctx.GetStub().GetState(makeDataKey(username, name, version))
	if err != nil {
		return false, fmt.Errorf("data is exist...: %v", err)
	}

	return dataAsBytes != nil, nil
}

func (d *DataChaincode) GetQueryDataHistory(ctx contractapi.TransactionContextInterface) ([]*DataType, error) {
	queryString := fmt.Sprintf(`{"selector":{"type":"data"}}`)
	return getQueryResultForQueryString(ctx, queryString)
}

func getQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*DataType, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var transferHistorys []*DataType
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var transferHistory DataType
		err = json.Unmarshal(queryResult.Value, &transferHistory)
		if err != nil {
			return nil, err
		}
		transferHistorys = append(transferHistorys, &transferHistory)
	}

	return transferHistorys, nil
}
