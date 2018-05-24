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


// Guarantee结构体
type Guarantee struct {
  GuaranteeID           string       `json:"guaranteeID"`
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
  Account               float64      `json:"account"`
}


// InsuranceCompany结构体
type InsuranceCompany struct {
  InsurCompanyID        string       `json:"insurCompanyID"`
  ObjectType            string       `json:"objectType"`
  PublicKeyX            string       `json:"publicKeyX"`
  PublicKeyY            string       `json:"publicKeyY"`
  Account               float64      `json:"account"`
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

  if function == "getRandBool" {
    return c.getRandBool(stub)
  } else if function == "clientRegist" {
    return c.clientRegist(stub, args)
  } else if function == "insurCompanyRegist" {
    return c.insurCompanyRegist(stub, args)
  } else if function == "writeOrder" {
    return c.writeOrder(stub, args)
  } else if function == "initGuarantee" {
    return c.initGuarantee(stub, args)
  } else if function == "queryByObjectType" {
    return c.queryByObjectType(stub, args)
  } else if function == "queryGuarantee" {
    return c.queryGuarantee(stub, args)
  } else if function == "queryOrder" {
    return c.queryOrder(stub, args)
  } else if function == "queryState" {
    return c.queryState(stub, args)
  } else if function == "clientDeposit" {
    return FsmEvent(stub, args, "clientDeposit")
  } else if function == "insurCompanyDeposit" {
    return FsmEvent(stub, args, "insurCompanyDeposit")
  } else if function == "clientWithdraw" {
    return FsmEvent(stub, args, "clientWithdraw")
  } else if function == "flightDelay" {
    return FsmEvent(stub, args, "flightDelay")
  } else if function == "insurCompanyWithdraw" {
    return FsmEvent(stub, args, "insurCompanyWithdraw")
  } else if function == "indemnify" {
    return FsmEvent(stub, args, "indemnify")
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
      {Name: "clientWithdraw", Src: []string{"Uninsuranced"}, Dst: "Client-Withdrawn"},
      {Name: "flightDelay", Src: []string{"Insuranced"}, Dst: "Delayed"},
      {Name: "indemnify", Src: []string{"Delayed"}, Dst: "Success"},
      {Name: "insurCompanyWithdraw", Src: []string{"Undelayed"}, Dst: "InsurCompany-Withdrawn"},
    },
    fsm.Callbacks{},
  )
  return f;
}


func FsmEvent(stub shim.ChaincodeStubInterface, args []string, event string) pb.Response{
  var resError pb.Response
  var str string

  guaranteeID := args[0]
  bstatus, err := stub.GetState(guaranteeID)         //从StateDB中读取对应表单的状态
  if err != nil{
    return shim.Error("Query guarantee status fail, guarantee ID: " + guaranteeID)
  }

  status := string(bstatus)
  fmt.Println("Guarantee[" + guaranteeID + "] status:" + status)
  //f := InitFSM(status)        //初始化状态机，并设置当前状态为表单的状态
  err = f.Event(event)          //触发状态机的事件
  if err != nil{
    return shim.Error("Current status is " + status + " does not support envent:" + event)
  } else {
    switch event {
      case "clientDeposit": resError, str = clientDeposit(stub, args)
      case "insurCompanyDeposit": resError, str = insurCompanyDeposit(stub, args)
      case "flightDelay": resError, str = changeOrder(stub, args)
      case "indemnify": resError, str = indemnify(stub, args)
      case "timeOut": str = "Done"
      case "clientWithdraw": resError, str = clientWithdraw(stub, args)
      case "insurCompanyWithdraw": resError, str = insurCompanyWithdraw(stub, args)
      default : resError = shim.Error("No event matches!!!")
    }

    if str != "Done" {
      return resError
    }
  }
  status = f.Current()
  fmt.Println("New status:" + status + "\n")
  stub.PutState(guaranteeID, []byte(status))        //更新表单的状态
  return shim.Success([]byte(status));              //返回新状态
}
