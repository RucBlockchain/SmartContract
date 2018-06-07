// main.go
package main


import (
  "fmt"
  "github.com/looplab/fsm"

  "github.com/hyperledger/fabric/core/chaincode/shim"
  pb "github.com/hyperledger/fabric/protos/peer"
)


type ContractChaincode struct {
}


// Policy结构体
type Policy struct {
  PolicyID              string       `json:"policyID"`
  OrderID               string       `json:"orderID"`
  ClientID              string       `json:"clientID"`
  InsurCompanyID        string       `json:"insurCompanyID"`
  ObjectType            string       `json:"objectType"`
  CreateTime            string       `json:"createTime"`
  Premium               float64      `json:"premium"`
  CashDeposit           float64      `json:"cashDeposit"`
}


// Client结构体
type Client struct {
  ClientID              string       `json:"clientID"`
  ObjectType            string       `json:"objectType"`
  PublicKeyX            string       `json:"publicKeyX"`
  PublicKeyY            string       `json:"publicKeyY"`
  Balance               float64      `json:"balance"`
}


// InsuranceCompany结构体
type InsuranceCompany struct {
  InsurCompanyID        string       `json:"insurCompanyID"`
  ObjectType            string       `json:"objectType"`
  PublicKeyX            string       `json:"publicKeyX"`
  PublicKeyY            string       `json:"publicKeyY"`
  Balance               float64      `json:"balance"`
}


// Order结构体
type Order struct {
  OrderID               string       `json:"orderID"`
  ClientID              string       `json:"clientID"`
  AirlineID             string       `json:"airlineID"`
  ObjectType            string       `json:"objectType"`
  CreateTime            string       `json:"createTime"`
  TicketStatus          string       `json:"ticketStatus"`
}


// =========================================
//       Init - initializes chaincode
// =========================================
func (c *ContractChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
  return shim.Success(nil)
}


// ============
//     Main
// ============
func main() {
  err := shim.Start(new(ContractChaincode))
  if err != nil {
    fmt.Printf("Error starting Contract chaincode: %s", err)
  }
}


// ======================================================
//       Invoke - Our entry point for Invocations
// ======================================================
func (c *ContractChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
  function, args := stub.GetFunctionAndParameters()
  fmt.Println("invoke is running " + function)

  if function == "check_buyTicket" {
    return c.check_buyTicket(stub)
  } else if function == "check_flightDelay" {
    return c.check_flightDelay(stub)
  } else if function == "clientRegist" {
    return c.clientRegist(stub, args)
  } else if function == "insurCompanyRegist" {
    return c.insurCompanyRegist(stub, args)
  } else if function == "buyTicket" {
    return c.buyTicket(stub, args)
  } else if function == "initPolicy" {
    return c.initPolicy(stub, args)
  } else if function == "queryByObjectType" {
    return c.queryByObjectType(stub, args)
  } else if function == "queryPolicy" {
    return c.queryPolicy(stub, args)
  } else if function == "queryOrder" {
    return c.queryOrder(stub, args)
  } else if function == "queryState" {
    return c.queryState(stub, args)
  } else if function == "clientDeposit" {
    return FsmEvent(stub, args, "clientDeposit")
  } else if function == "insurCompanyDeposit" {
    return FsmEvent(stub, args, "insurCompanyDeposit")
  } else if function == "clientRefund" {
    return FsmEvent(stub, args, "clientRefund")
  } else if function == "flightDelay" {
    return FsmEvent(stub, args, "flightDelay")
  } else if function == "insurCompanyRefund" {
    return FsmEvent(stub, args, "insurCompanyRefund")
  } else if function == "compensate" {
    return FsmEvent(stub, args, "compensate")
  } else if function == "timeOut" {
    return FsmEvent(stub, args, "timeOut")
  } else {
    return shim.Error("Function doesn't exits, make sure function is right!")
  }

  return shim.Success(nil)
}


func InitFSM(initStatus string) *fsm.FSM {
  f := fsm.NewFSM(
    initStatus,
    fsm.Events{
      {Name: "clientDeposit", Src: []string{"Bought"}, Dst: "Deposited"},
      {Name: "timeOut", Src: []string{"Bought"}, Dst: "Undeposited"},
      {Name: "insurCompanyDeposit", Src: []string{"Deposited"}, Dst: "Insuranced"},
      {Name: "timeOut", Src: []string{"Deposited"}, Dst: "Uninsuranced"},
      {Name: "clientRefund", Src: []string{"Uninsuranced"}, Dst: "ClientRefund"},
      {Name: "flightDelay", Src: []string{"Insuranced"}, Dst: "Delayed"},
      {Name: "compensate", Src: []string{"Delayed"}, Dst: "Success"},
      {Name: "insurCompanyRefund", Src: []string{"Undelayed"}, Dst: "InsurCompanyRefund"},
    },
    fsm.Callbacks{},
  )
  return f;
}


func FsmEvent(stub shim.ChaincodeStubInterface, args []string, event string) pb.Response{
  var resError pb.Response
  var str string

  policyID := args[0]
  bstatus, err := stub.GetState(policyID)         //从StateDB中读取对应表单的状态
  if err != nil{
    return shim.Error("Query policy status fail, policy ID: " + policyID)
  }

  status := string(bstatus)
  fmt.Println("Policy[" + policyID + "] status:" + status)
  //f := InitFSM(status)        //初始化状态机，并设置当前状态为表单的状态
  err = f.Event(event)          //触发状态机的事件
  if err != nil{
    return shim.Error("Current status is " + status + " does not support envent:" + event)
  } else {
    switch event {
      case "clientDeposit": resError, str = clientDeposit(stub, args)
      case "insurCompanyDeposit": resError, str = insurCompanyDeposit(stub, args)
      case "flightDelay": resError, str = changeOrder(stub, args)
      case "compensate": resError, str = compensate(stub, args)
      case "timeOut": str = "Done"
      case "clientRefund": resError, str = clientRefund(stub, args)
      case "insurCompanyRefund": resError, str = insurCompanyRefund(stub, args)
      default : resError = shim.Error("No event matches!!!")
    }

    if str != "Done" {
      return resError
    }
  }
  status = f.Current()
  fmt.Println("New status:" + status + "\n")
  stub.PutState(policyID, []byte(status))        //更新表单的状态
  return shim.Success([]byte(status));              //返回新状态
}
