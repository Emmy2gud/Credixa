package dto

type CreateVirtualAccountRequest struct {
	UserID    uint   `json:"user_id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Phone     string `json:"phonenumber"`
	Bvn       string `json:"bvn"`
}

type VirtualAccountResponse struct {
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
}
