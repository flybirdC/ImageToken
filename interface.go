package ImageToken

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)


type SmartContract struct {

}


func (s *SmartContract) Init(stub shim.ChaincodeStubInterface) peer.Response {

	token := &Token{Currency: map[string]Currency{}}

	tokenAsBytes, err := json.Marshal(token)
	err = stub.PutState(TokenKey,tokenAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	} else {
		fmt.Printf("Init Token %s \n",string(tokenAsBytes))
	}

	return shim.Success(nil)
}

func (s *SmartContract) Query(stub shim.ChaincodeStubInterface) peer.Response  {

	function, args := stub.GetFunctionAndParameters()
	switch function {
	case "balance" :
		return s.balance(stub,args)
	case "balanceAll":
		return s.balanceAll(stub,args)
	case "showAccount":
		return s.showAccount(stub,args)

	default:

	}

	return shim.Error("Invalid Smart Contract function name!")
}

func (s *SmartContract)Invoke(stub shim.ChaincodeStubInterface) peer.Response  {

	function, args := stub.GetFunctionAndParameters()
	switch function {
	case "initLedger":
		return s.initLedger(stub,args)
	case "createAccount":
		return s.createAccount(stub,args)
	case "initCurrency":
		return s.initCurrency(stub,args)
	case "setLock":
		return s.setLock(stub,args)
	case "transferToken":
		return s.transferToken(stub, args)
	case "frozenAccount":
		return s.frozenAccount(stub,args)
	case "mintToken":
		return s.mintToken(stub,args)
	case "balance":
		return s.balance(stub,args)
	case "balanceAll":
		return s.balanceAll(stub,args)
	case "showAccount":
		return s.showAccount(stub,args)
	case "showToken":
		return s.showToken(stub,args)


	default:

	}

	return shim.Error("Invalid Smart Contract function name.")
}

