package adapters

type CreateTransferRequest struct {
	AccountNumber string  `json:"account_number"`
	AccountBank   string  `json:"account_bank"`
	Amount        int64 `json:"amount"`
	Currency      string  `json:"currency" default:"NGN"`
	DebitCurrency string  `json:"debit_currency" default:"NGN"`
	Narration     string  `json:"narration"`
	Reference     string  `json:"reference"`
}


type CreateVirtualAccountRequest struct {
	Email       string `json:"email"`
	Phone       string `json:"phonenumber"`
	Amount      int64  `json:"amount"`
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


type CreateTransferResponse struct {
	Status string `json:"status"`
	Message string `json:"message"`
	Data struct {
		ID int `json:"id"`
		AccountNumber string `json:"account_number"`
		BankCode string `json:"bank_code"`
		FullName string `json:"full_name"`
		CreatedAt string `json:"created_at"`
		Currency string `json:"currency"`
		DebitCurrency string `json:"debit_currency"`
		Amount int64 `json:"amount"`
		Fee int64 `json:"fee"`
		Status string `json:"status"`
		Reference string `json:"reference"`
		Meta interface{} `json:"meta"`
		Narration string `json:"narration"`
		CompleteMessage string `json:"complete_message"`
		RequiresApproval int `json:"requires_approval"`
		IsApproved int `json:"is_approved"`
		BankName string `json:"bank_name"`
	}
}

