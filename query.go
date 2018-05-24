

package main


import (
	"fmt"
  "bytes"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)


// =============================================================================
//      queryGuarantee - query guarantee by guaranteeID from chaincode state
// =============================================================================
func (c *ContractChaincode) queryGuarantee(stub shim.ChaincodeStubInterface, args []string) pb.Response {

  //      0
	// "GuaranteeID"
  if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

  if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}

  guaranteeID := args[0]

	// ==== Create guarantee compositekey ====
	guaranteeIndexName := "guarantee"
	guaranteeIndexKey, err := stub.CreateCompositeKey(guaranteeIndexName, []string{guaranteeID})
	if err != nil {
		return shim.Error(err.Error())
	}

	// ==== Check if guarantee exists ====
	guaranteeAsBytes, err := stub.GetState(guaranteeIndexKey)
	if err != nil {
		return shim.Error("Failed to get guarantee: " + err.Error())
	} else if guaranteeAsBytes == nil {
		return shim.Error("Guarantee doesn't exit, please check your guarantee ID!")
	}

	return shim.Success(guaranteeAsBytes)
}


// ======================================================================
//      queryOrder - query order by order ID from chaincode state
// ======================================================================
func (c *ContractChaincode) queryOrder(stub shim.ChaincodeStubInterface, args []string) pb.Response {

  //     0
	// "OrderID"
  if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

  if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}

  orderID := args[0]

	// ==== Create order compositekey ====
	orderIndexName := "order"
	orderIndexKey, err := stub.CreateCompositeKey(orderIndexName, []string{orderID})
	if err != nil {
		return shim.Error(err.Error())
	}

	// ==== Check if order exists ====
	orderAsBytes, err := stub.GetState(orderIndexKey)
	if err != nil {
		return shim.Error("Failed to get order: " + err.Error())
	} else if orderAsBytes == nil {
		return shim.Error("Order doesn't exit, please check your order ID!")
	}

	return shim.Success(orderAsBytes)
}


func (c *ContractChaincode) queryState(stub shim.ChaincodeStubInterface, args []string) pb.Response {

  //       0
	// "GuaranteeID"
  if len(args) != 1 {
    return shim.Error("Incorrect number of arguments. Expecting 1")
  }

	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}

  guaranteeID := args[0]

  // ==== get current state of guarantee ====
	stateAsByte, err := stub.GetState(guaranteeID)
	if err != nil {
		fmt.Println("Failed to read state!")
	}
	state := string(stateAsByte)

  fmt.Println("current state is :", state)
  return shim.Success(stateAsByte)
}


// ==========================================================================================================
//      queryByObjectType - query objects by objectType from chaincode state
// ==========================================================================================================
func (c *ContractChaincode) queryByObjectType(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//       0
	// "ObjectType"
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	objectType := args[0]

	queryString := fmt.Sprintf("{\"selector\":{\"objectType\":\"%s\"}}", objectType)
	queryResults, err := getResultForQueryString(stub, queryString)

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(queryResults)
}


// =========================================================================================
// getQueryResultForQueryString executes the passed in query string.
// Result set is built and returned as a byte array containing the JSON results.
// =========================================================================================
func getResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// ==== buffer is a JSON array containing QueryRecords ====
	var buffer bytes.Buffer
	buffer.WriteString("[")

	providerAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// ==== Add a comma before array members, suppress it for the first array member =====
		if providerAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		//buffer.WriteString("}")
		providerAlreadyWritten = true
	}
	buffer.WriteString("]\n")

	fmt.Printf("- getResultForQueryString queryResult:\n%s\n", buffer.String())
	return buffer.Bytes(), nil
}
