package accounts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

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

	respbody, _ := io.ReadAll(resp.Body)
	fmt.Println("FlutterWave response:", string(respbody))

	w.Header().Set("Content-Type", "application/json")
	w.Write(respbody)
}
