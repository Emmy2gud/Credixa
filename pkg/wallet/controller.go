package wallet

import (
	"encoding/json"
	"fmt"
	"net/http"

	"payme/pkg/middleware"
	"payme/pkg/utils"


	"github.com/gorilla/mux"
)

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
}

type ChargeWithTokenRequest struct {
	Token  string `json:"token"`
	Amount int    `json:"amount"`
	Email  string `json:"email"`
	TxRef  string `json:"tx_ref"`
}

func GetWalletBalance(w http.ResponseWriter, r *http.Request) {
	// GetBalance logic here
}

func InitiateWalletFunding(w http.ResponseWriter, r *http.Request) {
	var userCard ChargeRequest
	utils.ParseBody(r, &userCard)
	//importing service function
	txref, result, err := InitiateCardWalletFunding(userCard, r)
	response := map[string]interface{}{
		"tx_ref":      txref,
		"flutterwave": result,
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(result)
	body, err := json.Marshal(response)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)

}

func AuthorizeCardFunding(w http.ResponseWriter, r *http.Request) {
	// The user only needs to send their PIN and the tx_ref from step 1.
	// The full card payload is retrieved automatically from the database.
	var req struct {
		TxRef string `json:"tx_ref"`
		Pin   string `json:"pin"`
	}
	utils.ParseBody(r, &req)

	if req.TxRef == "" || req.Pin == "" {
		http.Error(w, "tx_ref and pin are required", http.StatusBadRequest)
		return
	}

	result, err := AuthorizeCardFundingService(req.TxRef, req.Pin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
func ValidateWalletFunding(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FlwRef string `json:"flw_ref"`
		Otp    string `json:"otp"`
	}

	utils.ParseBody(r, &req)

	if req.FlwRef == "" || req.Otp == "" {
		http.Error(w, "flw_ref and otp are required", http.StatusBadRequest)
		return
	}

	result, err := ValidateCardCharge(req.FlwRef, req.Otp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Just return result — DO NOT CREDIT WALLET HERE
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Payment processing, awaiting confirmation",
		"data":    result,
	})
}

func VerifyCardCharge(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		http.Error(w, "transaction id is required", http.StatusBadRequest)
		return
	}

	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "user not authenticated", http.StatusUnauthorized)
		return
	}

	result, err := VerifyCard(id, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}



// func ChargeWithToken(w http.ResponseWriter, r *http.Request) {
//     var tokendata ChargeWithTokenRequest
// 	   utils.ParseBody(r, &tokendata)

// 	payload := map[string]interface{}{
// 		"token":    tokendata.Token,
// 		"currency": "NGN",
// 		"amount":   tokendata.Amount,
// 		"email":    tokendata.Email,

// 	}
// 	body, _ := json.Marshal(payload)

// 	req, _ := http.NewRequest("POST", "https://api.flutterwave.com/v3/tokenized-charges", bytes.NewBuffer(body))
// 	req.Header.Set("Authorization", "Bearer "+os.Getenv("FLW_SECRET_KEY"))
// 	req.Header.Set("Content-Type", "application/json")

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	var result map[string]interface{}
// 	respBody, _ := io.ReadAll(resp.Body)
// 	json.Unmarshal(respBody, &result)

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(map[string]interface{}{
// 		"message": "Payment successful",
// 		"data":   result,
// 	})

// }


