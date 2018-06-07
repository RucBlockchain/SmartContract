

package main


import (
  "fmt"
  "encoding/json"

  "github.com/hyperledger/fabric/core/chaincode/shim"
  pb "github.com/hyperledger/fabric/protos/peer"
)


// =================================================
//        clientDeposit - client pay deposit
// =================================================
func clientDeposit(stub shim.ChaincodeStubInterface, args []string) (pb.Response, string) {
  var clientJSON Client

  //       0            1           2             3                4           5      6      7
  // "PolicyID", "OrderID", "ClientID", "InsurCompanyID", "CreateTime", "$act2", "$r2", "$s2"
  if len(args) != 8 {
    return shim.Error("Incorrect number of arguments. Expecting 8"), ""
  }

  if len(args[0]) <= 0 {
    return shim.Error("1st argument must be a non-empty string"), ""
  }
  if len(args[1]) <= 0 {
    return shim.Error("2nd argument must be a non-empty string"), ""
  }
  if len(args[2]) <= 0 {
    return shim.Error("3rd argument must be a non-empty string"), ""
  }
  if len(args[3]) <= 0 {
    return shim.Error("4th argument must be a non-empty string"), ""
  }
  if len(args[4]) <= 0 {
    return shim.Error("5th argument must be a non-empty string"), ""
  }
  if len(args[5]) <= 0 {
    return shim.Error("6th argument must be a non-empty string"), ""
  }
  if len(args[6]) <= 0 {
    return shim.Error("7th argument must be a non-empty string"), ""
  }
  if len(args[7]) <= 0 {
    return shim.Error("8th argument must be a non-empty string"), ""
  }

  policyID := args[0]
  orderID := args[1]
  clientID := args[2]
  insurCompanyID := args[3]
  createTime := args[4]
  act := args[5]
  r := args[6]
  s := args[7]
  premium := float64(100)
  cashDeposit := float64(0)
  objectType := "policy"

  // ==== Create client compositekey ====
  clientIndexName := "client"
  clientIndexKey, err := stub.CreateCompositeKey(clientIndexName, []string{clientID})
  if err != nil {
    return shim.Error(err.Error()), ""
  }  
  value := []byte{0x00}

// ==== Check if client exists ====
clientAsBytes, err := stub.GetState(clientIndexKey)
  if err != nil {
    return shim.Error("Failed to get client: " + err.Error()), ""
  } else if clientAsBytes == nil {
    fmt.Println("Client not exits, please check client ID!")
    return shim.Error("Client not exits, please check client ID!"), ""
  }

  err = json.Unmarshal([]byte(clientAsBytes), &clientJSON)
  if err != nil {
    jsonResp := "{\"Error\":\"Failed to decode JSON of: " + clientID + "\"}"
    return shim.Error(jsonResp), ""
  }

  x := clientJSON.PublicKeyX
  y := clientJSON.PublicKeyY
  newBalance := clientJSON.Balance - premium

  newClientJSON := &Client{clientID, clientJSON.ObjectType, x, y, newBalance}
  newClientJSONAsBytes, err := json.Marshal(newClientJSON)
  if err != nil {
    return shim.Error(err.Error()), ""
  }

  // ==== Save new client to state ====
  err = stub.PutState(clientIndexKey, newClientJSONAsBytes)
  if err != nil {
    return shim.Error(err.Error()), ""
  }

  verify := verifySignature(x, y, r, s, act)
  if verify == false {
    return shim.Error("Verify failed!!!"), ""
  }

  // ==== Create policy compositekey ====
  policyIndexName := "policy"
  policyIndexKey, err := stub.CreateCompositeKey(policyIndexName, []string{policyID})
  if err != nil {
    return shim.Error(err.Error()), ""
  }

  // ==== Check whether policy alread exists ====
  policyAsBytes, err := stub.GetState(policyIndexKey)
  if err != nil {
    return shim.Error("Failed to get policy: " + err.Error()), ""
  } else if policyAsBytes != nil {
    fmt.Println("Policy alread exits, you cannot deposit again!")
    return shim.Error("Policy already exits, you cannot deposit again!"), ""
  }
  stub.PutState(policyIndexKey, value)

  policyJSON := &Policy{policyID, orderID, clientID, insurCompanyID, objectType, createTime, premium, cashDeposit}
  policyJSONAsBytes, err := json.Marshal(policyJSON)
  if err != nil {
    return shim.Error(err.Error()), ""
  }

  // ==== Save policy to state ====
  err = stub.PutState(policyIndexKey, policyJSONAsBytes)
  if err != nil {
    return shim.Error(err.Error()), ""
  }
  return shim.Success(nil), "Done"
}


// ==========================================================================
//        insurCompanyDeposit - insurance company pay for cash deposit
// ==========================================================================
func insurCompanyDeposit(stub shim.ChaincodeStubInterface, args []string) (pb.Response, string) {
  var insurCompanyJSON InsuranceCompany
  var policyJSON Policy

  //       0            1           2             3                4           5      6      7
  // "PolicyID", "OrderID", "ClientID", "InsurCompanyID", "CreateTime", "$act1", "$r1", "$s1"
  if len(args) != 8 {
    return shim.Error("Incorrect number of arguments. Expecting 8"), ""
  }

  if len(args[0]) <= 0 {
    return shim.Error("1st argument must be a non-empty string"), ""
  }
  if len(args[1]) <= 0 {
    return shim.Error("2nd argument must be a non-empty string"), ""
  }
  if len(args[2]) <= 0 {
    return shim.Error("3rd argument must be a non-empty string"), ""
  }
  if len(args[3]) <= 0 {
    return shim.Error("4th argument must be a non-empty string"), ""
  }
  if len(args[4]) <= 0 {
    return shim.Error("5th argument must be a non-empty string"), ""
  }
  if len(args[5]) <= 0 {
    return shim.Error("6th argument must be a non-empty string"), ""
  }
  if len(args[6]) <= 0 {
    return shim.Error("7th argument must be a non-empty string"), ""
  }
  if len(args[7]) <= 0 {
    return shim.Error("8th argument must be a non-empty string"), ""
  }

  policyID := args[0]
  orderID := args[1]
  clientID := args[2]
  insurCompanyID := args[3]
  createTime := args[4]
  act := args[5]
  r := args[6]
  s := args[7]
  cashDeposit := float64(10000)
  objectType := "policy"

  // ==== Create insurance company compositekey ====
  companyIndexName := "company"
  companyIndexKey, err := stub.CreateCompositeKey(companyIndexName, []string{insurCompanyID})
  if err != nil {
    return shim.Error(err.Error()), ""
  }
  value := []byte{0x00}

  // ==== Check if insurance company exists ====
  companyAsBytes, err := stub.GetState(companyIndexKey)
  if err != nil {
    return shim.Error("Failed to get insurance company: " + err.Error()), ""
  } else if companyAsBytes == nil {
    fmt.Println("Insurance company not exits, please check insurance company ID!")
    return shim.Error("Insurance company not exits, please check insurance company ID!"), ""
  }

  err = json.Unmarshal([]byte(companyAsBytes), &insurCompanyJSON)
  if err != nil {
    jsonResp := "{\"Error\":\"Failed to decode JSON of: " + insurCompanyID + "\"}"
    return shim.Error(jsonResp), ""
  }

  x := insurCompanyJSON.PublicKeyX
  y := insurCompanyJSON.PublicKeyY
  newBalance := insurCompanyJSON.Balance - cashDeposit

  newCompanyJSON := &InsuranceCompany{insurCompanyID, insurCompanyJSON.ObjectType, x, y, newBalance}
  newCompanyJSONAsBytes, err := json.Marshal(newCompanyJSON)
  if err != nil {
    return shim.Error(err.Error()), ""
  }

  // ==== Save new insurance company to state ====
  err = stub.PutState(companyIndexKey, newCompanyJSONAsBytes)
  if err != nil {
    return shim.Error(err.Error()), ""
  }

  verify := verifySignature(x, y, r, s, act)
  if verify == false {
    return shim.Error("Verify failed!!!"), ""
  }

  // ==== Create policy compositekey ====
  policyIndexName := "policy"
  policyIndexKey, err := stub.CreateCompositeKey(policyIndexName, []string{policyID})
  if err != nil {
    return shim.Error(err.Error()), ""
  }

  // ==== Check if policy exists ====
  policyAsBytes, err := stub.GetState(policyIndexKey)
  if err != nil {
    return shim.Error("Failed to get policy: " + err.Error()), ""
  } else if policyAsBytes == nil {
    fmt.Println("Policy not exits, please check policy ID!")
    return shim.Error("Policy not exits, please check policy ID!"), ""
  }
  stub.PutState(policyIndexKey, value)

  err = json.Unmarshal([]byte(policyAsBytes), &policyJSON)
  if err != nil {
    jsonResp := "{\"Error\":\"Failed to decode JSON of: " + policyID + "\"}"
    return shim.Error(jsonResp), ""
  }

  premium := policyJSON.Premium

  newPolicyJSON := &Policy{policyID, orderID, clientID, insurCompanyID, objectType, createTime, premium, cashDeposit}
  newPolicyJSONAsBytes, err := json.Marshal(newPolicyJSON)
  if err != nil {
    return shim.Error(err.Error()), ""
  }

  // ==== Save policy to state ====
  err = stub.PutState(policyIndexKey, newPolicyJSONAsBytes)
  if err != nil {
    return shim.Error(err.Error()), ""
  }
  return shim.Success(nil), "Done"
}


// ==========================================================
//        clientWithdraw - refund premium to client
// ==========================================================
func clientRefund(stub shim.ChaincodeStubInterface, args []string) (pb.Response, string) {
  var clientJSON Client

  //       0             1         2       3      4
  // "PolicyID", "ClientID", "$act3", "$r3", "$s3"
  if len(args) != 5 {
    return shim.Error("Incorrect number of arguments. Expecting 5"), ""
  }

  if len(args[0]) <= 0 {
    return shim.Error("1st argument must be a non-empty string"), ""
  }
  if len(args[1]) <= 0 {
    return shim.Error("2nd argument must be a non-empty string"), ""
  }
  if len(args[2]) <= 0 {
    return shim.Error("3rd argument must be a non-empty string"), ""
  }
  if len(args[3]) <= 0 {
    return shim.Error("4th argument must be a non-empty string"), ""
  }
  if len(args[4]) <= 0 {
    return shim.Error("5th argument must be a non-empty string"), ""
  }

  clientID := args[1]
  act := args[2]
  r := args[3]
  s := args[4]

  // ==== Create client compositekey ====
  clientIndexName := "client"
  clientIndexKey, err := stub.CreateCompositeKey(clientIndexName, []string{clientID})
  if err != nil {
    return shim.Error(err.Error()), ""
  }
  value := []byte{0x00}

  // ==== Check if client exists ====
  clientAsBytes, err := stub.GetState(clientIndexKey)
  if err != nil {
    return shim.Error("Failed to get client: " + err.Error()), ""
  } else if clientAsBytes == nil {
    return shim.Error("Client doesn't exit, please check your client ID!"), ""
  }

  err = json.Unmarshal([]byte(clientAsBytes), &clientJSON)
  if err != nil {
    jsonResp := "{\"Error\":\"Failed to decode JSON of: " + clientID + "\"}"
    return shim.Error(jsonResp), ""
  }

  objectType := clientJSON.ObjectType
  x := clientJSON.PublicKeyX
  y := clientJSON.PublicKeyY
  balance := clientJSON.Balance + 100

  verify := verifySignature(x, y, r, s, act)
  if verify == false {
    return shim.Error("Verify failed!!!"), ""
  }

  newClientJSON := &Client{clientID, objectType, x, y, balance}
  clientJSONAsBytes, err := json.Marshal(newClientJSON)
  if err != nil {
    return shim.Error(err.Error()), ""
  }
  stub.PutState(clientIndexKey, value)

  // ==== Save client to state ====
  err = stub.PutState(clientIndexKey, clientJSONAsBytes)
  if err != nil {
    return shim.Error(err.Error()), ""
  }

  return shim.Success(nil), "Done"
}


// ==========================================================
//        changeOrder - save ticketStatus into state
// ==========================================================
func changeOrder(stub shim.ChaincodeStubInterface, args []string) (pb.Response, string) {
  var orderJSON Order

  //       0            1             2
  // "PolicyID", "OrderID", "TicketStatus"
  if len(args) != 3 {
    return shim.Error("Incorrect number of arguments. Expecting 3"), ""
  }

  if len(args[0]) <= 0 {
    return shim.Error("1st argument must be a non-empty string"), ""
  }
  if len(args[1]) <= 0 {
    return shim.Error("2nd argument must be a non-empty string"), ""
  }
  if len(args[2]) <= 0 {
    return shim.Error("3rd argument must be a non-empty string"), ""
  }

  orderID := args[1]
  ticketStatus := args[2]

  // ==== Create order compositekey ====
  orderIndexName := "order"
  orderIndexKey, err := stub.CreateCompositeKey(orderIndexName, []string{orderID})
  if err != nil {
    return shim.Error(err.Error()), ""
  }
  value := []byte{0x00}

  // ==== Check if order exists ====
  orderAsBytes, err := stub.GetState(orderIndexKey)
  if err != nil {
    return shim.Error("Failed to get order: " + err.Error()), ""
  } else if orderAsBytes == nil {
    return shim.Error("Order doesn't exit, please check your order ID!"), ""
  }

  err = json.Unmarshal([]byte(orderAsBytes), &orderJSON)
  if err != nil {
    jsonResp := "{\"Error\":\"Failed to decode JSON of: " + orderID + "\"}"
    return shim.Error(jsonResp), ""
  }

  clientID := orderJSON.ClientID
  airlineID := orderJSON.AirlineID
  objectType := orderJSON.ObjectType
  createTime := orderJSON.CreateTime

  newOrderJSON := &Order{orderID, clientID, airlineID, objectType, createTime, ticketStatus}
  orderJSONAsBytes, err := json.Marshal(newOrderJSON)
  if err != nil {
    return shim.Error(err.Error()), ""
  }
  stub.PutState(orderIndexKey, value)

  // ==== Save policy to state ====
  err = stub.PutState(orderIndexKey, orderJSONAsBytes)
  if err != nil {
    return shim.Error(err.Error()), ""
  }

  return shim.Success(nil), "Done"
}


// ===============================================================================
//        insurCompanyWithdraw - refund cash deposit to insurance company
// ===============================================================================
func insurCompanyRefund(stub shim.ChaincodeStubInterface, args []string) (pb.Response, string) {
  var insuranceCompanyJSON InsuranceCompany

  //       0                1            2       3      4
  // "PolicyID", "InsurCompanyID", '$act2', "$r2", "$s2"
  if len(args) != 5 {
    return shim.Error("Incorrect number of arguments. Expecting 5"), ""
  }

  if len(args[0]) <= 0 {
    return shim.Error("1st argument must be a non-empty string"), ""
  }
  if len(args[1]) <= 0 {
    return shim.Error("2nd argument must be a non-empty string"), ""
  }
  if len(args[2]) <= 0 {
    return shim.Error("3rd argument must be a non-empty string"), ""
  }
  if len(args[3]) <= 0 {
    return shim.Error("4th argument must be a non-empty string"), ""
  }
  if len(args[4]) <= 0 {
    return shim.Error("5th argument must be a non-empty string"), ""
  }

  insurCompanyID := args[1]
  act := args[2]
  r := args[3]
  s := args[4]

  // ==== Create insurance company compositekey ====
  companyIndexName := "company"
  companyIndexKey, err := stub.CreateCompositeKey(companyIndexName, []string{insurCompanyID})
  if err != nil {
    return shim.Error(err.Error()), ""
  }
  value := []byte{0x00}

  // ==== Check if insurance company exists ====
  companyAsBytes, err := stub.GetState(companyIndexKey)
  if err != nil {
    return shim.Error("Failed to get insurance company: " + err.Error()), ""
  } else if companyAsBytes == nil {
    return shim.Error("Insurance company doesn't exit, please check your insurance company ID!"), ""
  }

  err = json.Unmarshal([]byte(companyAsBytes), &insuranceCompanyJSON)
  if err != nil {
    jsonResp := "{\"Error\":\"Failed to decode JSON of: " + insurCompanyID + "\"}"
    return shim.Error(jsonResp), ""
  }

  objectType := insuranceCompanyJSON.ObjectType
  x := insuranceCompanyJSON.PublicKeyX
  y := insuranceCompanyJSON.PublicKeyY
  balance := insuranceCompanyJSON.Balance + 10100

  verify := verifySignature(x, y, r, s, act)
  if verify == false {
    return shim.Error("Verify failed!!!"), ""
  }

  newCompanyJSON := &InsuranceCompany{insurCompanyID, objectType, x, y, balance}
  companyJSONAsBytes, err := json.Marshal(newCompanyJSON)
  if err != nil {
    return shim.Error(err.Error()), ""
  }
  stub.PutState(companyIndexKey, value)

  // ==== Save insurance company to state ====
  err = stub.PutState(companyIndexKey, companyJSONAsBytes)
  if err != nil {
    return shim.Error(err.Error()), ""
  }

  return shim.Success(nil), "Done"
}


// ==========================================================
//        indemnify - refund cash deposit to client
// ==========================================================
func compensate(stub shim.ChaincodeStubInterface, args []string) (pb.Response, string) {
  var clientJSON Client

  //       0             1         2       3      4
  // "PolicyID", "ClientID", "$act4", "$r4", "$s4"
  if len(args) != 5 {
    return shim.Error("Incorrect number of arguments. Expecting 5"), ""
  }

  if len(args[0]) <= 0 {
    return shim.Error("1st argument must be a non-empty string"), ""
  }
  if len(args[1]) <= 0 {
    return shim.Error("2nd argument must be a non-empty string"), ""
  }
  if len(args[2]) <= 0 {
    return shim.Error("3rd argument must be a non-empty string"), ""
  }
  if len(args[3]) <= 0 {
    return shim.Error("4th argument must be a non-empty string"), ""
  }
  if len(args[4]) <= 0 {
    return shim.Error("5th argument must be a non-empty string"), ""
  }

  clientID := args[1]
  act := args[2]
  r := args[3]
  s := args[4]

  // ==== Create client compositekey ====
  clientIndexName := "client"
  clientIndexKey, err := stub.CreateCompositeKey(clientIndexName, []string{clientID})
  if err != nil {
    return shim.Error(err.Error()), ""
  }
  value := []byte{0x00}

  // ==== Check if client exists ====
  clientAsBytes, err := stub.GetState(clientIndexKey)
  if err != nil {
    return shim.Error("Failed to get client: " + err.Error()), ""
  } else if clientAsBytes == nil {
    return shim.Error("Client doesn't exit, please check your client ID!"), ""
  }

  err = json.Unmarshal([]byte(clientAsBytes), &clientJSON)
  if err != nil {
    jsonResp := "{\"Error\":\"Failed to decode JSON of: " + clientID + "\"}"
    return shim.Error(jsonResp), ""
  }

  objectType := clientJSON.ObjectType
  x := clientJSON.PublicKeyX
  y := clientJSON.PublicKeyY
  balance := clientJSON.Balance + 10000

  verify := verifySignature(x, y, r, s, act)
  if verify == false {
    return shim.Error("Verify failed!!!"), ""
  }

  newClientJSON := &Client{clientID, objectType, x, y, balance}
  clientJSONAsBytes, err := json.Marshal(newClientJSON)
  if err != nil {
    return shim.Error(err.Error()), ""
  }
  stub.PutState(clientIndexKey, value)

  // ==== Save client to state ====
  err = stub.PutState(clientIndexKey, clientJSONAsBytes)
  if err != nil {
    return shim.Error(err.Error()), ""
  }

  return shim.Success(nil), "Done"
}
