package transfer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
	// "payme/pkg/transaction"
)

type InternalTransferPayload struct {
	AccountNumber string  `json:"account_number"`
	AccountBank   string  `json:"account_bank"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	DebitCurrency string  `json:"debit_currency"`
	Narration     string  `json:"narration"`
	Reference     string  `json:"reference"`
}

func ResolveBankDetails(w http.ResponseWriter, r *http.Request) {
	type Input struct {
		AccountNumber string `json:"account_number"`
		AccountBank   string `json:"account_bank"`
	}
	var input Input
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	//convert to json
	payload, err := json.Marshal(input)
	req, err := http.NewRequest("POST", "https://api.flutterwave.com/v3/accounts/resolve", bytes.NewBuffer(payload))
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

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Flutterwave response:", string(body))

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}


func InitializeFunding(w http.ResponseWriter, r *http.Request) {
	// var transactions transaction.Transaction
	// userID, _ := middleware.GetUserID(r)
	type Input struct {
		AccountNumber string  `json:"account_number"`
		AccountBank   string  `json:"account_bank"`
		Amount        float64 `json:"amount"`
		Narration     string  `json:"narration"`
	}
	var input Input
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	transfer := InternalTransferPayload{
		AccountNumber: input.AccountNumber,
		AccountBank:   input.AccountBank,
		Amount:        input.Amount,
		Narration:     input.Narration,
		Currency:      "NGN",
		DebitCurrency: "NGN",
		Reference:     "TXN-" + time.Now().Format("20060102150405"),
	}
	//convert to json
	payload, err := json.Marshal(transfer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest("POST", "https://api.flutterwave.com/v3/transfers", bytes.NewBuffer(payload))
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

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Flutterwave response:", string(body))

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)

}

func VerifyFunding(w http.ResponseWriter, r *http.Request) {
	// VerifyFunding logic here
}
