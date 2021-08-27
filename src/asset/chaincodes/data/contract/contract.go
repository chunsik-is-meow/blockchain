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
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Downloaded  int    `json:"downloaded"`
	Owner       string `json:"owner"`
	Timestamp   string `json:"timestamp"`
}

// InitLedger ...
func (d *DataChaincode) InitLedger(ctx contractapi.TransactionContextInterface) error {
	dataInfos := []DataType{
		{Type: "data", Name: "adult", Description: "Census Income classfication", Downloaded: 0, Owner: "Ronny Kohavi and Barry Becker", Timestamp: "2021-08-27-09-11-49"},
		{Type: "data", Name: "breast-cancer-wisconsin", Description: "Cancer classfication", Downloaded: 0, Owner: "Olvi L. Mangasarian.", Timestamp: "2021-08-27-09-11-49"},
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

			dataKey := makeDataKey(data.Owner)
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
func (d *DataChaincode) DataInsert(ctx contractapi.TransactionContextInterface, name string, description string, owner string, timestamp string) error {
	// TODO
	// data file upload
	return nil
}

func (d *DataChaincode) PutCommonData(ctx contractapi.TransactionContextInterface, name string, description string, owner string, timestamp string) error {
	// d.DataInsert(ctx, name, description, owner, timestamp)
	exists, err := d.dataExists(ctx, name)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the data %s already exists", name)
	}

	down := 0
	data := DataType{
		Type:        "data",
		Name:        name,
		Description: description,
		Downloaded:  down,
		Owner:       owner,
		Timestamp:   timestamp,
	}
	dataAsBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to json.Marshal(). %v", err)
	}
	dataKey := makeDataKey(owner)
	ctx.GetStub().PutState(dataKey, dataAsBytes)
	if err != nil {
		return fmt.Errorf("failed to put to world state. %v", err)
	}

	return nil
}

func (d *DataChaincode) GetAllCommonDataInfo(ctx contractapi.TransactionContextInterface) ([]*DataType, error) {
	var dataInfos []*DataType
	datasAsBytes, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}

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

func (d *DataChaincode) GetCommonDataInfo(ctx contractapi.TransactionContextInterface, name string) (*DataType, error) {
	dataInfo := &DataType{}
	dataAsBytes, err := ctx.GetStub().GetState(makeDataKey(name))
	if err != nil {
		return nil, err
	} else if dataAsBytes == nil {
		dataInfo.Type = "empty"
		dataInfo.Name = "empty"
		dataInfo.Description = "empty"
		dataInfo.Downloaded = 0
		dataInfo.Owner = "empty"
		dataInfo.Timestamp = "empty"
	} else {
		err = json.Unmarshal(dataAsBytes, dataInfo)
		if err != nil {
			return nil, err
		}
	}

	return dataInfo, nil
}

func makeDataKey(classfication string) string {
	var sb strings.Builder

	sb.WriteString("D_")
	sb.WriteString(classfication)
	return sb.String()
}

func (d *DataChaincode) dataExists(ctx contractapi.TransactionContextInterface, name string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(name)
	if err != nil {
		return false, fmt.Errorf("data is exist...: %v", err)
	}

	return assetJSON != nil, nil
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
