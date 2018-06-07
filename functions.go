

package main


import (
  "fmt"
  "sync"
  "time"
  "io"
  "math/rand"
  "math/big"
  "encoding/json"
  "hash"
  "crypto/ecdsa"
  "crypto/md5"
  "crypto/elliptic"
  "github.com/looplab/fsm"

  "github.com/hyperledger/fabric/core/chaincode/shim"
  pb "github.com/hyperledger/fabric/protos/peer"
)


var f *fsm.FSM


func getRand(num int) int {
  rand.Seed(time.Now().UnixNano())
  var mu sync.Mutex
  mu.Lock()
  v := rand.Intn(num)
  mu.Unlock()
  return v
}


// =============================================================
//       check_buyTicket - check if client bought ticket
// =============================================================
func (c *ContractChaincode) check_buyTicket(stub shim.ChaincodeStubInterface) pb.Response {
  boolValue := getRand(2)
  if boolValue == 0 {
    return shim.Success([]byte("true"))
  } else {
    return shim.Success([]byte("false"))
  }
}


// ===================================================
//       check_delay - check if flight delayed
// ===================================================
func (c *ContractChaincode) check_flightDelay(stub shim.ChaincodeStubInterface) pb.Response {
  boolValue := getRand(2)
  if boolValue == 0 {
    return shim.Success([]byte("true"))
  } else {
    return shim.Success([]byte("false"))
  }
}


// =======================================================
//       clientRegist - save user's information to state
// =======================================================
func (c *ContractChaincode) clientRegist(stub shim.ChaincodeStubInterface, args []string) pb.Response {

  //       0       1     2
  // "ClientID", "$x", "$y"
  if len(args) != 3 {
    return shim.Error("Incorrect number of arguments. Expecting 3")
  }

  if len(args[0]) <= 0 {
    return shim.Error("1st argument must be a non-empty string")
  }
  if len(args[1]) <= 0 {
    return shim.Error("2nd argument must be a non-empty string")
  }
  if len(args[2]) <= 0 {
    return shim.Error("3rd argument must be a non-empty string")
  }

  clientID := args[0]
  publicKeyX := args[1]
  publicKeyY := args[2]
  objectType := "client"
  balance := float64(2000)

  // ==== Create client compositekey ====
  clientIndexName := "client"
  clientIndexKey, err := stub.CreateCompositeKey(clientIndexName, []string{clientID})
  if err != nil {
    return shim.Error(err.Error())
  }
  value := []byte{0x00}

  // ==== Check if client already exists ====
  clientAsBytes, err := stub.GetState(clientIndexKey)
  if err != nil {
    return shim.Error("Failed to get client: " + err.Error())
  } else if clientAsBytes != nil {
    fmt.Println("The client already exists! Please change your clientID.")
    return shim.Error("The client already exists! Please change your clientID.")
  }
  stub.PutState(clientIndexKey, value)

  clientJSON := &Client{clientID, objectType, publicKeyX, publicKeyY, balance}
  clientJSONAsBytes, err := json.Marshal(clientJSON)
  if err != nil {
    return shim.Error(err.Error())
  }

  // ==== Save client to state ====
  err = stub.PutState(clientIndexKey, clientJSONAsBytes)
  if err != nil {
    return shim.Error(err.Error())
  }

  fmt.Println("regist successfully", clientID)
  return shim.Success(nil)
}


// ===============================================================================
//       insurCompanyRegist - save insurance company's information to state
// ===============================================================================
func (c *ContractChaincode) insurCompanyRegist(stub shim.ChaincodeStubInterface, args []string) pb.Response {

  //         0           1     2
  // "InsurCompanyID", "$x", "$y"
  if len(args) != 3 {
    return shim.Error("Incorrect number of arguments. Expecting 3")
  }

  if len(args[0]) <= 0 {
    return shim.Error("1st argument must be a non-empty string")
  }
  if len(args[1]) <= 0 {
    return shim.Error("2nd argument must be a non-empty string")
  }
  if len(args[2]) <= 0 {
    return shim.Error("3rd argument must be a non-empty string")
  }

  insurCompanyID := args[0]
  objectType := "insurance company"
  publicKeyX := args[1]
  publicKeyY := args[2]
  balance := float64(200000)

  // ==== Create insurance company compositekey ====
  companyIndexName := "company"
  companyIndexKey, err := stub.CreateCompositeKey(companyIndexName, []string{insurCompanyID})
  if err != nil {
    return shim.Error(err.Error())
  }
  value := []byte{0x00}

  // ==== Check if insurance company already exists ====
  insurCompanyAsBytes, err := stub.GetState(companyIndexKey)
  if err != nil {
    return shim.Error("Failed to get insurance company: " + err.Error())
  } else if insurCompanyAsBytes != nil {
    fmt.Println("The insurance company already exists! Please change your insurCompanyID.")
    return shim.Error("The insurance company already exists! Please change your insurCompanyID.")
  }
  stub.PutState(companyIndexKey, value)

  companyJSON := &InsuranceCompany{companyIndexKey, objectType, publicKeyX, publicKeyY, balance}
  companyJSONAsBytes, err := json.Marshal(companyJSON)
  if err != nil {
    return shim.Error(err.Error())
  }

  // ==== Save insurance company to state ====
  err = stub.PutState(companyIndexKey, companyJSONAsBytes)
  if err != nil {
    return shim.Error(err.Error())
  }

  fmt.Println("regist successfully", insurCompanyID)
  return shim.Success(nil)
}


// ================================================
//       writeOrder - save orders into state
// ================================================
func (c *ContractChaincode) buyTicket(stub shim.ChaincodeStubInterface, args []string) pb.Response {
  var clientJSON Client

  //     0           1           2            3           4       5      6
  // "OrderID", "ClientID", "AirlineID", "CreateTime", "$act1", "$r1", "$s1"
  if len(args) != 7 {
    return shim.Error("Incorrect number of arguments. Expecting 7")
  }

  if len(args[0]) <= 0 {
    return shim.Error("1st argument must be a non-empty string")
  }
  if len(args[1]) <= 0 {
    return shim.Error("2nd argument must be a non-empty string")
  }
  if len(args[2]) <= 0 {
    return shim.Error("3rd argument must be a non-empty string")
  }
  if len(args[3]) <= 0 {
    return shim.Error("4th argument must be a non-empty string")
  }
  if len(args[4]) <= 0 {
    return shim.Error("5th argument must be a non-empty string")
  }
  if len(args[5]) <= 0 {
    return shim.Error("6th argument must be a non-empty string")
  }
  if len(args[6]) <= 0 {
    return shim.Error("7th argument must be a non-empty string")
  }

  orderID := args[0]
  clientID := args[1]
  airlineID := args[2]
  createTime := args[3]
  act := args[4]
  r := args[5]
  s := args[6]
  ticketStatus := "OPEN FORUSE"
  objectType := "order"

  // ==== Create client compositekey ====
  clientIndexName := "client"
  clientIndexKey, err := stub.CreateCompositeKey(clientIndexName, []string{clientID})
  if err != nil {
    return shim.Error(err.Error())
  }
  value := []byte{0x00}

  // ==== Check if client exists ====
  clientAsBytes, err := stub.GetState(clientIndexKey)
  if err != nil {
    return shim.Error("Failed to get client: " + err.Error())
  } else if clientAsBytes == nil {
    fmt.Println("The client not exists!")
    return shim.Error("The client not exists!")
  }

  err = json.Unmarshal([]byte(clientAsBytes), &clientJSON)
  if err != nil {
    jsonResp := "{\"Error\":\"Failed to decode JSON of: " + clientID + "\"}"
    return shim.Error(jsonResp)
  }

  x := clientJSON.PublicKeyX
  y := clientJSON.PublicKeyY

  verify := verifySignature(x, y, r, s, act)
  if verify == false {
    return shim.Error("Verify failed!!!")
  }

  // ==== Create order compositekey ====
  orderIndexName := "order"
  orderIndexKey, err := stub.CreateCompositeKey(orderIndexName, []string{orderID})
  if err != nil {
    return shim.Error(err.Error())
  }

  // ==== Check if order already exists ====
  orderAsBytes, err := stub.GetState(orderIndexKey)
  if err != nil {
    return shim.Error("Failed to get order: " + err.Error())
  } else if orderAsBytes != nil {
    fmt.Println("The order already exists!")
    return shim.Error("The order already exists!")
  }
  stub.PutState(orderIndexKey, value)

  orderJSON := &Order{orderID, clientID, airlineID, objectType, createTime, ticketStatus}
  orderJSONAsBytes, err := json.Marshal(orderJSON)
  if err != nil {
    return shim.Error(err.Error())
  }

  // ==== Save order to state ====
  err = stub.PutState(orderIndexKey, orderJSONAsBytes)
  if err != nil {
    return shim.Error(err.Error())
  }

  return shim.Success(nil)
}


// =============================================================
//        initPolicy - set policy init status Bought
// =============================================================
func (c *ContractChaincode) initPolicy(stub shim.ChaincodeStubInterface, args []string) pb.Response {

  //     0
  // "PolicyID"
  if len(args) != 1 {
    return shim.Error("Incorrect number of arguments. Expecting 1")
  }

  if len(args[0]) <= 0 {
    return shim.Error("1st argument must be a non-empty string")
  }

  policyID := args[0]

  status := string("Bought")
  fmt.Println("Policy[" + policyID + "] status:" + status)
  f = InitFSM(status)        //初始化状态机，并设置当前状态为表单的状态

  stub.PutState(policyID, []byte(status))        //更新表单的状态

  return shim.Success(nil)
}


func verifySignature(x string, y string, r string, s string, act string) bool {
  // 拼接x, y
  bigx := new(big.Int)
  nx, _ := bigx.SetString(x, 16)

  bigy := new(big.Int)
  ny, _ := bigy.SetString(y, 16)

  var pub ecdsa.PublicKey
  pub.Curve = elliptic.P256()
  pub.X = nx
  pub.Y = ny

  bigr := new(big.Int)
  nr, _ := bigr.SetString(r, 16)

  bigs := new(big.Int)
  ns, _ := bigs.SetString(s, 16)

  var h hash.Hash
  h = md5.New()
  io.WriteString(h, act)
  nsig := h.Sum(nil)

  verify := ecdsa.Verify(&pub, nsig, nr, ns)
  fmt.Println(verify) // should be true

  return verify
}
