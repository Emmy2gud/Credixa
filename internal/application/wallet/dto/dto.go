package dto

type ChargeRequest struct {
	CardBrand   string `json:"card_brand"`
	Last4       string `json:"last4"`
	CardNumber  string `json:"card_number"`
	CVV         string `json:"cvv"`
	ExpiryMonth string `json:"expiry_month"`
	ExpiryYear  string `json:"expiry_year"`
	Currency    string `json:"currency"`
	Amount      string `json:"amount"`
	Email       string `json:"email"`
	Fullname    string `json:"fullname"`
	TxRef       string `json:"tx_ref"` // unique ID you generate
	Token       string `json:"token"`
	WalletID    uint   `json:"wallet_id"`
	City        string `json:"city"`
	Address     string `json:"address"`
	State       string `json:"state"`
	Zip         string `json:"zip"`
}

type ChargeWithTokenRequest struct {
	Token  string `json:"token"`
	Amount int    `json:"amount"`
	Email  string `json:"email"`
	TxRef  string `json:"tx_ref"`
}

type InitializeCardResponse struct {
	Mode   string   `json:"mode"`
	Fields []string `json:"fields"`
}

type AuthorizeCardResponse struct {
	Mode     string `json:"mode"`
	Endpoint string `json:"endpoint"`
}

type ValidateCardRequest struct {
	FlwRef string `json:"flw_ref"`
	Otp    string `json:"otp"`
	Type   string `json:"type"`
}
type VerifyCardRequest struct {
	TxRef string `json:"tx_ref"`
	Otp   string `json:"otp"`
}

type ValidateCardResponse struct {
	
		Status string `json:"status"`
		Amount int64  `json:"amount"`

		First6Digits string `json:"first_6_digits"`
		Last4Digits  string `json:"last_4_digits"`
		Issuer       string `json:"issuer"`
		Country      string `json:"country"`
		Type         string `json:"type"`
		Expiry       string `json:"expiry"`
	
}

type VerifyChargeResponse struct {

		Status string `json:"status"`
		Amount int64  `json:"amount"`
		Token  string `json:"token"`

		First6Digits string `json:"first_6_digits"`
		Last4Digits  string `json:"last_4_digits"`
		Issuer       string `json:"issuer"`
		Country      string `json:"country"`
		Type         string `json:"type"`
		Expiry       string `json:"expiry"`
	
}
