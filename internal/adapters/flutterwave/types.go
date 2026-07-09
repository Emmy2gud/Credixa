package adapters

type CreateTransferRequest struct {
	AccountNumber string `json:"account_number"`
	AccountBank   string `json:"account_bank"`
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency" default:"NGN"`
	DebitCurrency string `json:"debit_currency" default:"NGN"`
	Narration     string `json:"narration"`
	Reference     string `json:"reference"`
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
		Frequency       int32  `json:"frequency"`
		BankName        string `json:"bank_name"`
		CreatedAt       string `json:"created_at"`
		ExpiryDate      string `json:"expiry_date"`
		Note            string `json:"note"`
		Amount          string `json:"amount"`
	} `json:"data"`
}

type CreateTransferResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		ID               int         `json:"id"`
		AccountNumber    string      `json:"account_number"`
		BankCode         string      `json:"bank_code"`
		FullName         string      `json:"full_name"`
		CreatedAt        string      `json:"created_at"`
		Currency         string      `json:"currency"`
		DebitCurrency    string      `json:"debit_currency"`
		Amount           int64       `json:"amount"`
		Fee              int64       `json:"fee"`
		Status           string      `json:"status"`
		Reference        string      `json:"reference"`
		Meta             interface{} `json:"meta"`
		Narration        string      `json:"narration"`
		CompleteMessage  string      `json:"complete_message"`
		RequiresApproval int         `json:"requires_approval"`
		IsApproved       int         `json:"is_approved"`
		BankName         string      `json:"bank_name"`
	}
}

// wallet funding dto

type InitializeCardResponse struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
	Meta    struct {
		Authorization struct {
			Mode   string   `json:"mode"`
			Fields []string `json:"fields"`
		}
	} `json:"meta"`
}

type AuthorizationCardResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Meta    struct {
		Authorization struct {
			Mode     string `json:"mode"`
			Redirect string `json:"redirect"`
			Endpoint string `json:"endpoint"`
		}
	}
}

type ValidateCardRequest struct {
	FlwRef string `json:"flw_ref"`
	Otp    string `json:"otp"`
	Type   string `json:"type"`
}
type VerifyCardRequest struct {
	TxRef string `json:"flw_ref"`
	Otp   string `json:"otp"`
}

type ValidateCardResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Status string `json:"status"`
		Amount int64  `json:"amount"`
		Card   struct {
			First6Digits string `json:"first_6_digits"`
			Last4Digits  string `json:"last_4_digits"`
			Issuer       string `json:"issuer"`
			Country      string `json:"country"`
			Type         string `json:"type"`
			Expiry       string `json:"expiry"`
		}
		Customer struct {
			ID          int    `json:"id"`
			PhoneNumber string `json:"phone_number"`
			Name        string `json:"name"`
			Email       string `json:"email"`
			CreatedAt   string `json:"created_at"`
		}
	}
}

type VerifyChargeResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Status string `json:"status"`
		Amount int64  `json:"amount"`
		Token  string `json:"token"`
		Card   struct {
			First6Digits string `json:"first_6_digits"`
			Last4Digits  string `json:"last_4_digits"`
			Issuer       string `json:"issuer"`
			Country      string `json:"country"`
			Type         string `json:"type"`
			Expiry       string `json:"expiry"`
		}
		Customer struct {
			ID          int    `json:"id"`
			PhoneNumber string `json:"phone_number"`
			Name        string `json:"name"`
			Email       string `json:"email"`
			CreatedAt   string `json:"created_at"`
		}
	}
}

// kyc verification
type InitiateKycTier2Request struct {
	Bvn         string `json:"bvn"`
	FirstName   string `json:"firstname"`
	LastName    string `json:"lastname"`
	RedirectUrl string `json:"redirect_url"`
}
type InitiateKycTier2Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Url string `json:"url"`
		Ref string `json:"ref"`
	}
}

type VerifyKycTier2Request struct {
	Ref string `json:"ref"`
}

type RetrieveBvnResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Status  string `json:"status"`
		BvnData struct {
			Bvn         string `json:"bvn"`
			DateOfBirth string `json:"date_of_birth"`
			PhoneNumber string `json:"phone_number"`
			FirstName   string `json:"firstname"`
			LastName    string `json:"lastname"`
			Gender      string `json:"gender"`
			Email       string `json:"email"`
			Address     string `json:"address"`
			Nin         string `json:"nin"`
		}
		Reference string `json:"reference"`
		CreatedAt string `json:"created_at"`
	}
}
