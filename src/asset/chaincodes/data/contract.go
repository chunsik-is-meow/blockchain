package contract

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// DataChaincode ...
type DataChaincode struct {
	contractapi.Contract
}

// DataType ...
type DataType struct {
	Type           string `json:"type"`
	Name           string `json:"name"`
	Size           int    `json:"size"`
	Year           string `json:"year"`
	Classification string `json:"classification"`
	Description    string `json:"description"`
}

// INSERT Data history
func (d *DataChaincode) DataInsert(ctx contractapi.TransactionContextInterface, name string, size uint32, year string, classfication string, description string, timestamp string) error {
	dataMeow := DataType{
		Type:           "data",
		Name:           name,
		Size:           100,
		Year:           year,
		Classification: classfication,
		Description:    description,
	}

	dataMeowAsBytes, err := json.Marshal(dataMeow)
	if err != nil {
		return fmt.Errorf("failed to json.Marshal(). %v", err)
	}
	dataMeowKey := makedataMeowKey(classfication, timestamp)
	ctx.GetStub().PutState(dataMeowKey, dataMeowAsBytes)
	if err != nil {
		return fmt.Errorf("failed to put to world state. %v", err)
	}

	return nil
}

func makedataMeowKey(classfication string, timestamp string) string {
	var sb strings.Builder

	sb.WriteString("D_")
	sb.WriteString(classfication)
	sb.WriteString("_")
	sb.WriteString(timestamp)

	return sb.String()
}

func (d *DataChaincode) GetQueryDataHistory(ctx contractapi.TransactionContextInterface) ([]*DataType, error) {
	queryString := fmt.Sprintf(`{"selector":{"type":"data"}}`)
	return getQueryResultForQueryString(ctx, queryString)
}

func (d *DataChaincode) GetQueryDataClassficationHistory(ctx contractapi.TransactionContextInterface, classfication string) ([]*DataType, error) {
	queryString := fmt.Sprintf(`{"selector":{"type":"transfer","classfication":"%s"}}`, classfication)
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

func (d *DataChaincode) PutCommonData(ctx contractapi.TransactionContextInterface, key string, data string, timestamp string) error {
	//TODO

	return nil
}

type irisInfo struct {
	Min  string
	Max  string
	Mean string
	SD   string
}

func (d *DataChaincode) GetCommonData(ctx contractapi.TransactionContextInterface, key string, timestamp string) error {
	data, err := os.Open("../../datas/iris.json")
	if err != nil {
		return fmt.Errorf("failed to open. %v", err)
	}
	byteValue, _ := ioutil.ReadAll(data)
	var data_info irisInfo

	json.Unmarshal(byteValue, &data_info)

	fmt.Println(data_info)

	return nil
}
