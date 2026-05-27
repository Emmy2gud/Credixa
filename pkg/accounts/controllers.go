package accounts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"payme/pkg/config"
	"payme/pkg/middleware"
	"payme/pkg/wallet"

	"github.com/google/uuid"
)

// BillPaymentController handles HTTP requests for bill payments
type VirtualAccountPayload struct {
	Email        string `json:"email"`
	Phone        string  `json:"phonenumber"`
	Amount       int32 `json:"amount"`
	Currency     string `json:"currency"`
	Firstname    string `json:"firstname"`
	Lastname     string `json:"lastname"`
	Is_permanent bool   `json:"is_permanent"`
	TxRef        string `json:"tx_ref"`
	Narration    string `json:"narration"`
	BankCode     string `json:"bank_code"`
	Bvn          string  `json:"bvn"`
}

func CreateVirtualAccount(w http.ResponseWriter, r *http.Request) {
	var payload VirtualAccountPayload


	var input struct {
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

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}
    payload.Amount=input.Amount
	payload.Phone=input.Phone
	payload.Firstname=input.Firstname
	payload.Lastname=input.Lastname
	payload.Email=input.Email
	payload.Currency="NGN"
	payload.BankCode="090772"
	payload.Bvn = input.Bvn
	payload.Narration="Create a virtual account for this user"
	payload.Is_permanent=true
	payload.TxRef = "token_ch_" + uuid.New().String()
    flwPayload, _ := json.Marshal(&payload)
	req, err := http.NewRequest("POST", "https://api.flutterwave.com/v3/virtual-account-numbers", bytes.NewBuffer(flwPayload))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Authorization", "Bearer "+os.Getenv("FLW_SECRET_KEY"))
	req.Header.Set("Content-Type", "application/json")
	

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	respbody, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}
	fmt.Println("FlutterWave response:", string(respbody))

	var result map[string]interface{}
	if err := json.Unmarshal(respbody, &result); err != nil {
		http.Error(w, "Failed to parse API response", http.StatusInternalServerError)
		return
	}



	data, ok := result["data"].(map[string]interface{})
	if !ok {
		http.Error(w, "Invalid data structure in API response", http.StatusInternalServerError)
		return
	}

	var wallet wallet.Wallet
	userID, _ := middleware.GetUserID(r)
	if err := config.DB.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		http.Error(w, "wallet not found", http.StatusBadRequest)
		return
	}

	virtualAcc := VirtualAccount{
		WalletID:      uint(wallet.ID),
		AccountNumber: data["account_number"].(string),
		AccountName:   input.Firstname + " " + input.Lastname,
		BankName:      data["bank_name"].(string),
		Provider:      "flutterwave",
		Status:        data["account_status"].(string),
	}

	// Fallback if AccountName is missing from API
	if virtualAcc.AccountName == "" {
		virtualAcc.AccountName = fmt.Sprintf("%s %s", input.Firstname, input.Lastname)
	}

	if err := config.DB.Create(&virtualAcc).Error; err != nil {
		http.Error(w, "Failed to create virtual account in database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respbody)
}
