package ImageToken

import (
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"strconv"
	"time"
)



//get the address and key, private data
func (s *SmartContract)getUserCrypto(stub shim.ChaincodeStubInterface, args []string) peer.Response  {

	if len(args) != 2 {
		return shim.Error("enter password :first argument is password, the second is to remember password or not!")
	}
	password := args[0]
	rememberIs,_ := strconv.ParseBool(args[1])

	timestamp := strconv.FormatInt(time.Now().UnixNano(),64)

	userAccount := &Account{}
	address, privatekey := userAccount.makeCryptoEcc()


	userCrypto := &UserCrypto{

		TimeID:timestamp,
		Address:address,
		Privatekey:privatekey,
	}

	cryptoAsBytes,_ := json.Marshal(userCrypto)


	if rememberIs {

		stub.PutPrivateData("collectionUsersDetail",password,cryptoAsBytes)

	}


	return shim.Success(cryptoAsBytes)
}

//query crypto private: address, key

func (s *SmartContract)queryUserCrypto(stub shim.ChaincodeStubInterface, args []string) peer.Response  {

	if len(args) != 1 {
		return shim.Error("agument is password!")
	}

	password := args[0]

	cryptoAsBytes, err := stub.GetPrivateData("collectionUserDetail",password)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(cryptoAsBytes)
}