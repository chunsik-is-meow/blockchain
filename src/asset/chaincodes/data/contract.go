package contract

import (
	"encoding/json"
	"fmt"
	"strings"
	"io/ioutil"
    "os"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Info struct{
    Username string
    Password string
    Hostname string
    Port string
}

// DataChaincode ...
type DataChaincode struct {
	contractapi.Contract
}

// DataType ...
type DataType struct {
	Type         string `json:"type"`
	Name       string `json:"name"`
	Size       int    `json:"size"`
	Year       string `json:"year"`
	Classification   string `json:"classification"`
	Description      string `json:"description"`
}

// INSERT Data history
func (t *DataChaincode) DataInsert(ctx contractapi.TransactionContextInterface, name string, size uint32, year string, classfication string, description string, timestamp string) error {
	dataMeow := DataType{
		Type:	"data",
		Name:	name,
		Size:	size,
		Year:	year,
		Classification:	classfication,
		Description:	description,
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

func (t *TradeChaincode) GetQueryDataHistory(ctx contractapi.TransactionContextInterface) ([]*TransferType, error) {
	queryString := fmt.Sprintf(`{"selector":{"type":"data"}}`)
	return getQueryResultForQueryString(ctx, queryString)
}

func (t *TradeChaincode) GetQueryDataClassficationHistory(ctx contractapi.TransactionContextInterface, classfication string) ([]*TransferType, error) {
	queryString := fmt.Sprintf(`{"selector":{"type":"transfer","classfication":"%s"}}`, classfication)
	return getQueryResultForQueryString(ctx, queryString)
}

func (t *DataChaincode) PutCommonData(ctx contractapi.TransactionContextInterface, key string, data string, timestamp string) error {
	//TODO

	return nil
}


func (t *DataChaincode) GetCommonData(ctx contractapi.TransactionContextInterface, key string, timestamp string) error {
	data, err := os.Open("../datas/iris.json")
	byteValue, _ := ioutil.ReadAll(data)
	var data_info Info

	json.Unmarshal(byteValue, &db_info)

	fmt.Println(db_info)

	return nil
}