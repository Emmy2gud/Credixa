package dto
type CreateTransferRequest struct{
	UserID uint 
	AccountNumber string  `json:"account_number"`
	AccountBank   string  `json:"account_bank"`
	Amount        int64 `json:"amount"`
	Narration     string  `json:"narration"`
	Firstname     string  `json:"firstname"`
	Lastname      string  `json:"lastname"`
	IdempotencyKey string `json:"idempotency_key"`
}
 

type TransferResponse struct {

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