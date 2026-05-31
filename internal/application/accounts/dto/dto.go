package dto

	type CreateVirtualAccountInput struct {
		Email        string `json:"email"`
		Phone        string `json:"phonenumber"`
		Amount       int32  `json:"amount"`
		Firstname    string `json:"firstname"`
		Lastname     string `json:"lastname"`
		Bvn          string  `json:"bvn"`
		// TxRef        string `json:"tx_ref"`
		// Currency     string `json:"currency"`
		// Is_permanent bool   `json:"is_permanent"`
		// Narration    string `json:"narration"`
		// BankCode     string `json:"bank_code"`
	}


	type VirtualAccountResponse struct {
		ResponseCode    string  `json:"response_code"`
		ResponseMessage string  `json:"response_message"`
		FlwRef          string  `json:"flw_ref"`
		OrderRef        string  `json:"order_ref"`
		AccountNumber   string  `json:"account_number"`
		Frequency       string  `json:"frequency"`
		BankName        string  `json:"bank_name"`
		CreatedAt       string  `json:"created_at"`
		ExpiryDate      string  `json:"expiry_date"`
		Note            string  `json:"note"`
		Amount          string  `json:"amount"`

		
	}