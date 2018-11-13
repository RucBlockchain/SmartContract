

package main


import (
  "fmt"
  "bytes"

  "github.com/hyperledger/fabric/core/chaincode/shim"
  pb "github.com/hyperledger/fabric/protos/peer"
)


// =============================================================================
//      queryPolicy - query policy by policyID from chaincode state
// =============================================================================
func (c *ContractChaincode) queryPolicy(stub shim.ChaincodeStubInterface, args []string) pb.Response {

  //      0
  // "PolicyID"
  if len(args) != 1 {
    return shim.Error("Incorrect number of arguments. Expecting 1")
  }

  if len(args[0]) <= 0 {
    return shim.Error("1st argument must be a non-empty string")
  }

  policyID := args[0]

  // ==== Create policy compositekey ====
  policyIndexName := "policy"
  policyIndexKey, err := stub.CreateCompositeKey(policyIndexName, []string{policyID})
  if err != nil {
    return shim.Error(err.Error())
  }

  // ==== Check if policy exists ====
  policyAsBytes, err := stub.GetState(policyIndexKey)
  if err != nil {
    return shim.Error("Failed to get policy: " + err.Error())
  } else if policyAsBytes == nil {
    return shim.Error("Policy doesn't exit, please check your policy ID!")
  }

  return shim.Success(policyAsBytes)
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
  // "PolicyID"
  if len(args) != 1 {
    return shim.Error("Incorrect number of arguments. Expecting 1")
  }

  if len(args[0]) <= 0 {
    return shim.Error("1st argument must be a non-empty string")
  }

  policyID := args[0]

  // ==== get current state of policy ====
  stateAsByte, err := stub.GetState(policyID)
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
