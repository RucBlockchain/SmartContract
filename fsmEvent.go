

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
  // "GuaranteeID", "OrderID", "ClientID", "InsurCompanyID", "CreateTime", "$act2", "$r2", "$s2"
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

  guaranteeID := args[0]
  orderID := args[1]
  clientID := args[2]
  insurCompanyID := args[3]
  createTime := args[4]
  act := args[5]
  r := args[6]
  s := args[7]
  premium := float64(100)
  cashDeposit := float64(0)
  objectType := "guarantee"

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
  newAccount := clientJSON.Account - premium

  newClientJSON := &Client{clientID, clientJSON.ObjectType, x, y, newAccount}
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

  // ==== Create guarantee compositekey ====
  guaranteeIndexName := "guarantee"
  guaranteeIndexKey, err := stub.CreateCompositeKey(guaranteeIndexName, []string{guaranteeID})
  if err != nil {
    return shim.Error(err.Error()), ""
  }

  // ==== Check whether guarantee alread exists ====
  guaranteeAsBytes, err := stub.GetState(guaranteeIndexKey)
  if err != nil {
    return shim.Error("Failed to get guarantee: " + err.Error()), ""
  } else if guaranteeAsBytes != nil {
    fmt.Println("Guarantee alread exits, you cannot deposit again!")
    return shim.Error("Guarantee already exits, you cannot deposit again!"), ""
  }
  stub.PutState(guaranteeIndexKey, value)

  guaranteeJSON := &Guarantee{guaranteeID, orderID, clientID, insurCompanyID, objectType, createTime, premium, cashDeposit}
  guaranteeJSONAsBytes, err := json.Marshal(guaranteeJSON)
  if err != nil {
    return shim.Error(err.Error()), ""
  }

  // ==== Save guarantee to state ====
  err = stub.PutState(guaranteeIndexKey, guaranteeJSONAsBytes)
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
  var guaranteeJSON Guarantee

  //       0            1           2             3                4           5      6      7
  // "GuaranteeID", "OrderID", "ClientID", "InsurCompanyID", "CreateTime", "$act1", "$r1", "$s1"
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

  guaranteeID := args[0]
  orderID := args[1]
  clientID := args[2]
  insurCompanyID := args[3]
  createTime := args[4]
  act := args[5]
  r := args[6]
  s := args[7]
  cashDeposit := float64(10000)
  objectType := "guarantee"

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
  newAccount := insurCompanyJSON.Account - cashDeposit

  newCompanyJSON := &InsuranceCompany{insurCompanyID, insurCompanyJSON.ObjectType, x, y, newAccount}
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

  // ==== Create guarantee compositekey ====
  guaranteeIndexName := "guarantee"
  guaranteeIndexKey, err := stub.CreateCompositeKey(guaranteeIndexName, []string{guaranteeID})
  if err != nil {
    return shim.Error(err.Error()), ""
  }

  // ==== Check if guarantee exists ====
  guaranteeAsBytes, err := stub.GetState(guaranteeIndexKey)
  if err != nil {
    return shim.Error("Failed to get guarantee: " + err.Error()), ""
  } else if guaranteeAsBytes == nil {
    fmt.Println("Guarantee not exits, please check guarantee ID!")
    return shim.Error("Guarantee not exits, please check guarantee ID!"), ""
  }
  stub.PutState(guaranteeIndexKey, value)

  err = json.Unmarshal([]byte(guaranteeAsBytes), &guaranteeJSON)
  if err != nil {
    jsonResp := "{\"Error\":\"Failed to decode JSON of: " + guaranteeID + "\"}"
    return shim.Error(jsonResp), ""
  }

  premium := guaranteeJSON.Premium

  newGuaranteeJSON := &Guarantee{guaranteeID, orderID, clientID, insurCompanyID, objectType, createTime, premium, cashDeposit}
  newGuaranteeJSONAsBytes, err := json.Marshal(newGuaranteeJSON)
  if err != nil {
    return shim.Error(err.Error()), ""
  }

  // ==== Save guarantee to state ====
  err = stub.PutState(guaranteeIndexKey, newGuaranteeJSONAsBytes)
  if err != nil {
    return shim.Error(err.Error()), ""
  }
  return shim.Success(nil), "Done"
}


// ==========================================================
//        clientWithdraw - refund premium to client
// ==========================================================
func clientWithdraw(stub shim.ChaincodeStubInterface, args []string) (pb.Response, string) {
  var clientJSON Client

  //       0             1         2       3      4
  // "GuaranteeID", "ClientID", "$act3", "$r3", "$s3"
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
  account := clientJSON.Account + 100

  verify := verifySignature(x, y, r, s, act)
  if verify == false {
    return shim.Error("Verify failed!!!"), ""
  }

  newClientJSON := &Client{clientID, objectType, x, y, account}
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
  // "GuaranteeID", "OrderID", "TicketStatus"
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

  // ==== Save guarantee to state ====
  err = stub.PutState(orderIndexKey, orderJSONAsBytes)
  if err != nil {
    return shim.Error(err.Error()), ""
  }

  return shim.Success(nil), "Done"
}


// ===============================================================================
//        insurCompanyWithdraw - refund cash deposit to insurance company
// ===============================================================================
func insurCompanyWithdraw(stub shim.ChaincodeStubInterface, args []string) (pb.Response, string) {
  var insuranceCompanyJSON InsuranceCompany

  //       0                1            2       3      4
  // "GuaranteeID", "InsurCompanyID", '$act2', "$r2", "$s2"
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
  account := insuranceCompanyJSON.Account + 10100

  verify := verifySignature(x, y, r, s, act)
  if verify == false {
    return shim.Error("Verify failed!!!"), ""
  }

  newCompanyJSON := &InsuranceCompany{insurCompanyID, objectType, x, y, account}
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
func indemnify(stub shim.ChaincodeStubInterface, args []string) (pb.Response, string) {
  var clientJSON Client

  //       0             1         2       3      4
  // "GuaranteeID", "ClientID", "$act4", "$r4", "$s4"
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
  account := clientJSON.Account + 10000

  verify := verifySignature(x, y, r, s, act)
  if verify == false {
    return shim.Error("Verify failed!!!"), ""
  }

  newClientJSON := &Client{clientID, objectType, x, y, account}
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
