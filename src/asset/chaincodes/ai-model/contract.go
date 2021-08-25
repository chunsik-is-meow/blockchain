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
	ObjectType         string `json:"docType"`
	Name       string `json:"name"`
	Score       int    `json:"score"`
	OwnelID       string `json:"ownelID"`
	SrcDataID       string `json:"srcDataID"`
	Description      string `json:"description"`
}


func (t *DataChaincode) PutAIModel(ctx contractapi.TransactionContextInterface, key string, model string, timestamp string) error {
	return nil
}


func (t *DataChaincode) GetAIModel(ctx contractapi.TransactionContextInterface, key string, timestamp string) error {
	return nil
}


func (t *DataChaincode) checkScore(ctx contractapi.TransactionContextInterface, model string) error {
	return nil
}


func (t *DataChaincode) isVaildData(ctx contractapi.TransactionContextInterface, srcDataID string) error {
	return nil
}




// Transfer ...
func (t *AIChaincode) Transfer(ctx contractapi.TransactionContextInterface, from string, to string, amount uint32, timestamp string) error {
	// INSERT Transfer history
	transferMeow := AIModelType{
		Type:   "transfer",
		From:   from,
		To:     to,
		Amount: amount,
	}

	transferMeowAsBytes, err := json.Marshal(transferMeow)
	if err != nil {
		return fmt.Errorf("failed to json.Marshal(). %v", err)
	}
	transferMeowKey := makeFromToMeowKey(from, to, timestamp)
	ctx.GetStub().PutState(transferMeowKey, transferMeowAsBytes)

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
	data, err := os.Open("../datas/%s.json", key)
	byteValue, _ := ioutil.ReadAll(data)
	var db_info Info

	json.Unmarshal(byteValue, &db_info)

	fmt.Println(db_info)

	return nil
}