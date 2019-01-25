package ImageToken


type Msg struct {

	Status bool `json:"Status"`
	Code int `json:"Code"`
	Message string `json:"Message"`

}

type Currency struct {

	TokenName string `json:"TokenName"`
	TokenSymbol string `json:"TokenSybol"`
	TotalSupply float64 `json:"TotalSupply"`
	TokenAdd float64 `json:"TokenAdd"`
	TokenRelease float64 `json:"TokenRelease"`

}

type Token struct {
	Lock bool `json:"Lock"`
	Currency map[string]Currency `json:"Currency"`
}


type Account struct {

	Name string `json:"Name"`
	Frozen bool `json:"Frozen"`
	BalanceOf map[string]float64 `json:"BalanceOf"`
}

const TokenKey = "Token"



type UserCrypto struct {

	TimeID string `json:"TimeID"`
	Address string `json:"Address"`
	Privatekey string `json:"Privatekey"`

}


