package ImageToken

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"strconv"
)

//create account
func (s *SmartContract) createAccount(stub shim.ChaincodeStubInterface, args []string) peer.Response  {

	if len(args) != 1 {

		return shim.Error("Incorrect number of arguments.Parameter Expecting 1")

	}

	key := args[0]
	name :=args[0]
	existAsBytes ,err := stub.GetState(key)
	fmt.Printf("GetState(%s) \n",key, string(existAsBytes))
	if string(existAsBytes) != "" {
		fmt.Println("failed to create account, Duplicate key.")
		return shim.Error("failed to create account, Duplicate key")
	}

	account := Account{
		Name:name,
		Frozen:false,
		BalanceOf: map[string]float64{},
	}

	accountAsBytes, _ := json.Marshal(account)
	err = stub.PutState(key, accountAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("createAccount %s \n",string(accountAsBytes))

	return shim.Success(accountAsBytes)
}

//init ledger
func (s *SmartContract)initLedger(stub shim.ChaincodeStubInterface, args []string) peer.Response  {

	return shim.Success(nil)
}

//show token
func (s *SmartContract)showToken(stub shim.ChaincodeStubInterface, args []string) peer.Response  {

	tokenAsBytes, err := stub.GetState(TokenKey)
	if err != nil {
		return shim.Error(err.Error())
	}else {
		fmt.Printf("GetState(%s) %s \n",TokenKey,string(tokenAsBytes))
	}

	return shim.Success(tokenAsBytes)
}


//init currency
func (s *SmartContract)initCurrency(stub shim.ChaincodeStubInterface, args []string) peer.Response  {

	if len(args) != 4 {
		return shim.Error("Incorrect number of argument. Expecting 4!")
		}
	name := args[0]
	symbol := args[1]
	supply,_ := strconv.ParseFloat(args[2],64)
	account := args[3]

	coinbaseAsBytes, err := stub.GetState(account)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("Coinbase before %s \n", string(coinbaseAsBytes))

	coinbase := &Account{}

	json.Unmarshal(coinbaseAsBytes,&coinbase)

	token := Token{}

	existAsBytes, err := stub.GetState(TokenKey)

	if err != nil {
		return shim.Error(err.Error())
	}else {
		fmt.Printf("GetState(%s) %s \n", TokenKey,string(existAsBytes))
	}
	json.Unmarshal(existAsBytes, &token)

	result := token.initialSupply(name,symbol,supply,coinbase)

	tokenAsBytes, _ := json.Marshal(token)
	err = stub.PutState(TokenKey, tokenAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	} else {
		fmt.Printf("Init Token %s \n", string(tokenAsBytes))
	}

	coinbaseAsBytes, _ = json.Marshal(coinbase)

	err = stub.PutState(account,coinbaseAsBytes)

	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Printf("Coinbase after %s \n", string(coinbaseAsBytes))

	return shim.Success(result)
}

//transfertoken
func (s *SmartContract) transferToken(stub shim.ChaincodeStubInterface, args []string) peer.Response  {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments, Expecting 4")
	}

	from := args[0]
	to := args[1]
	currentcy := args[2]
	amount, _ := strconv.ParseFloat(args[3],32)

	if(amount <= 0) {
		return shim.Error("Incorrect number of amount!")
	}

	fromAsBytes,err := stub.GetState(from)

	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Printf("fromAccount %s \n",string(fromAsBytes))
	fromAccount := &Account{}
	json.Unmarshal(fromAsBytes,&fromAccount)

	toAsBytes, err := stub.GetState(to)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("toAccount %s \n", string(toAsBytes))
	toAccount := &Account{}
	json.Unmarshal(toAsBytes,&toAccount)

	tokenAsBytes, err := stub.GetState(TokenKey)

	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Printf("Token %s \n",string(toAsBytes))
	token := Token{Currency: map[string]Currency{}}
	json.Unmarshal(tokenAsBytes, &token)

	result := token.transfer(fromAccount,toAccount,currentcy,amount)
	fmt.Printf("Result %s \n", string(result))

	fromAsBytes, err = json.Marshal(fromAccount)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(from,fromAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	} else {
		fmt.Printf("fromAccount %s \n",string(fromAsBytes))

	}

	toAsBytes, err = json.Marshal(toAccount)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(to,toAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	} else {
		fmt.Printf("toAccount %s \n",string(toAsBytes))
	}

	return shim.Success(result)
}

//mint token, add token to poolaccount
func (s *SmartContract)mintToken(stub shim.ChaincodeStubInterface, args []string) peer.Response  {

	if len(args) != 3 {
		return shim.Error("incorrect number of arguments. Expecting 3")
	}

	currency := args[0]
	amount,_ := strconv.ParseFloat(args[1],32)
	account :=args[2]

	coinbaseAsbytes, err := stub.GetState(account)
	if err != nil {
		return shim.Error(err.Error())
	} else {
		fmt.Printf("Coinbase before %s \n",string(coinbaseAsbytes))
	}

	coinbase := &Account{}
	json.Unmarshal(coinbaseAsbytes, &coinbase)

	tokenAsBytes, err := stub.GetState(TokenKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("Token before %s \n",string(tokenAsBytes))

	token := Token{}

	json.Unmarshal(tokenAsBytes,&token)

	result := token.mint(currency,amount,coinbase)

	tokenAsBytes, err = json.Marshal(token)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("Token after %s \n",string(tokenAsBytes))

	coinbaseAsbytes, _ = json.Marshal(coinbase)
	err = stub.PutState(account,coinbaseAsbytes)
	if err != nil {
		return shim.Error(err.Error())
	}else {
		fmt.Printf("Coinbase after %s \n", string(tokenAsBytes))
	}

	fmt.Printf("mintToken %s \n",string(tokenAsBytes))

	return shim.Success(result)

}

//lock account
func (s *SmartContract) setLock(stub shim.ChaincodeStubInterface, args []string) peer.Response  {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	lock := args[0]

	tokenAsBytes, err := stub.GetState(TokenKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("setLock - begin %s \n",string(tokenAsBytes))

	token := Token{}

	json.Unmarshal(tokenAsBytes,&token)

	if (lock == "true") {
		token.setLock(true)
	}else {
		token.setLock(false)
	}

	tokenAsBytes, err = json.Marshal(token)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(TokenKey,tokenAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf("setLock - end %s \n",string(tokenAsBytes))

	return shim.Success(nil)
}

//frozen account
func (s *SmartContract)frozenAccount(stub shim.ChaincodeStubInterface, args []string) peer.Response  {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	_account := args[0]
	_status := args[1]

	accountAsBytes, err := stub.GetState(_account)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Printf("setLock - begin %s \n",string(accountAsBytes))

	account := Account{}

	json.Unmarshal(accountAsBytes,&account)

	var status bool
	if(_status == "true") {
		status = true
	}else {
		status = false
	}

	account.Frozen = status

	accountAsBytes, err = json.Marshal(account)

	if err != nil {
		return shim.Error(err.Error())
	}else {
		fmt.Printf("frozenAccount - end %s \n",string(accountAsBytes))
	}

	return shim.Success(nil)

	}

//show account
func (s *SmartContract) showAccount(stub shim.ChaincodeStubInterface, args []string) peer.Response  {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	_account := args[0]

	accountAsBytes, err := stub.GetState(_account)
	if err != nil {
		return shim.Error(err.Error())
	} else {
		fmt.Printf("Account balance %s \n",string(accountAsBytes))
	}
	return shim.Success(accountAsBytes)
}

//balance
func (s *SmartContract) balance(stub shim.ChaincodeStubInterface,args []string) peer.Response  {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments, Expecting 1")
	}

	_account := args[0]
	_currency := args[1]

	accountAsBytes, err := stub.GetState(_account)
	if err != nil {
		return shim.Error(err.Error())
	}else {
		fmt.Printf("Account balance %s \n",string(accountAsBytes))
	}

	account := Account{}
	json.Unmarshal(accountAsBytes,&account)
	result := account.balance(_currency)

	resultAsBytes, _ := json.Marshal(result)
	fmt.Printf("%s balance is %s \n",_account, string(resultAsBytes))

	return shim.Success(resultAsBytes)
}

//balanceAll
func (s *SmartContract) balanceAll(stub shim.ChaincodeStubInterface,args []string) peer.Response  {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	_account := args[0]

	accountAsBytes, err := stub.GetState(_account)
	if err != nil {
		return shim.Error(err.Error())
	}else {
		fmt.Printf("Account balance %s \n",string(accountAsBytes))
	}

	account :=Account{}
	json.Unmarshal(accountAsBytes, &account)
	result := account.balanceAll()
	resultAsBytes, _ := json.Marshal(result)
	fmt.Print("%s banlance is %s \n", _account, string(resultAsBytes))

	return shim.Success(resultAsBytes)
}




