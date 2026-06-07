package adapters

type CreateVirtualAccountRequest struct {
	Email       string `json:"email"`
	Phone       string `json:"phonenumber"`
	Amount      int32  `json:"amount"`
	Currency    string `json:"currency,omitempty"`
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	IsPermanent bool   `json:"is_permanent,omitempty"`
	TxRef       string `json:"tx_ref,omitempty"`
	Narration   string `json:"narration,omitempty"`
	BankCode    string `json:"bank_code,omitempty"`
	Bvn         string `json:"bvn"`
}

type CreateVirtualAccountResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`

	Data struct {
		ResponseCode    string `json:"response_code"`
		ResponseMessage string `json:"response_message"`
		FlwRef          string `json:"flw_ref"`
		OrderRef        string `json:"order_ref"`
		AccountNumber   string `json:"account_number"`
		Frequency       int32 `json:"frequency"`
		BankName        string `json:"bank_name"`
		CreatedAt       string `json:"created_at"`
		ExpiryDate      string `json:"expiry_date"`
		Note            string `json:"note"`
		Amount          string `json:"amount"`
	} `json:"data"`
}
