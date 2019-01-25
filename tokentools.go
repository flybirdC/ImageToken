package ImageToken

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
)

//token transffer
func (token *Token)transfer(from *Account, to *Account, currency string,value float64) []byte  {

	var rev []byte

	if (token.Lock) {
		msg := &Msg{Status:false, Code:1,Message:"lock account,stop token transfer"}
		rev ,_ = json.Marshal(msg)
		return rev
	}

	if(from.Frozen) {
		msg := &Msg{Status:false,Code:1,Message:"From Account Frozen!"}
		rev, _ = json.Marshal(msg)
		return rev
	}

	if(to.Frozen) {
		msg := &Msg{Status:false,Code:1,Message:"To Account Frozen!"}
		rev, _ = json.Marshal(msg)
		return rev
	}

	if(!token.isCurrency(currency)) {
		msg := &Msg{Status:false, Code:1,Message:"Symbol none!"}
		rev, _ = json.Marshal(msg)
		return rev
	}

	if (from.BalanceOf[currency] >= value) {
		from.BalanceOf[currency] -= value
		to.BalanceOf[currency] += value

		msg := &Msg{Status:true,Code:0, Message:"transfer success!"}
		rev , _ = json.Marshal(msg)
		return rev
	} else {
		msg := &Msg{Status:false, Code:1, Message:" insufficient capital on the account balance!"}
		rev ,_ = json.Marshal(msg)
		return rev
	}

}

//init token supply first total
func (token *Token)initialSupply(name string,symbol string, supply float64,account *Account) []byte  {

	if _,ok := token.Currency[symbol];ok {
		msg := &Msg{Status:false,Code:1,Message:"token symbol has been setup!"}
		rev, _ := json.Marshal(msg)
		return rev
	} else {
		token.Currency[symbol] = Currency{TokenName:name,TokenSymbol:symbol,TotalSupply:supply}
		account.BalanceOf[symbol] = supply

		msg := &Msg{Status:true,Code:0,Message:"init token symbol success!"}
		rev, _ := json.Marshal(msg)
		return rev
	}



}

//add token
func (token *Token)mint(currency string, amount float64,account *Account ) []byte  {

	if(!token.isCurrency(currency)) {
		msg := &Msg{Status:false, Code:1,Message:"symbol is none "}
		rev ,_ := json.Marshal(msg)
		return rev
	}

	cur := token.Currency[currency]
	cur.TotalSupply += amount
	token.Currency[currency] = cur
	account.BalanceOf[currency] += amount

	msg := &Msg{Status:true,Code:0,Message:"token add success!"}
	rev, _ := json.Marshal(msg)
	return rev
}

//release token
func (token *Token)burn(currency string, amount float64,account *Account) []byte {

	if(!token.isCurrency(currency)) {
		msg := &Msg{Status:false,Code:1,Message:"symbol is none"}
		rev, _ := json.Marshal(msg)
		return rev
	}

	if(token.Currency[currency].TotalSupply >= amount) {

		cur := token.Currency[currency]
		cur.TotalSupply -= amount
		token.Currency[currency] = cur
		account.BalanceOf[currency] -= amount

		msg := &Msg{Status:true,Code:0,Message:"release token success!"}
		rev, _ := json.Marshal(msg)
		return rev
	} else {
		msg := &Msg{Status:false,Code:1,Message:"release token failed, amount is not enough"}
		rev, _ := json.Marshal(msg)
		return rev
	}
}




//Currency state
func (token *Token)isCurrency(currency string) bool  {

	if _,ok := token.Currency[currency]; ok {
		return true
	}else {
		return false
	}
}

//lock
func (token *Token)setLock(lock bool) bool  {

	token.Lock = lock
	return token.Lock
}


//account balance
func (account *Account)balance(currency string) map[string]float64  {

	bal := map[string]float64{currency:account.BalanceOf[currency]}
	return bal

}

//balance
func (account *Account) balanceAll() map[string]float64  {

	return account.BalanceOf
}


//user crypto - ethereum crypto
func (account *Account)makeCryptoEcc()  (address string,privatekey string) {

	key, err := crypto.GenerateKey()
	if err != nil {
		fmt.Println("Error: ,",err.Error())
	}

	//get Account address
	address = crypto.PubkeyToAddress(key.PublicKey).Hex()
	fmt.Printf("address[%d][%v]\n",len(address),address)

	//get the private key
	privatekey = hex.EncodeToString(key.D.Bytes())
	fmt.Printf("privatekey[%d][%v]\n", len(privatekey),privatekey)

	return address,privatekey

}