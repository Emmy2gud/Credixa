package transfer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"gorm.io/gorm"
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

type TransferService interface {
	ResolveBankDetails(ctx context.Context, accountNumber, accountBank string) ([]byte, error)
	InitializeFunding(ctx context.Context, accountNumber, accountBank string, amount float64, narration string) ([]byte, error)
	VerifyFunding(ctx context.Context) error
}

type transferService struct {
	db *gorm.DB
}

func NewTransferService(db *gorm.DB) TransferService {
	return &transferService{
		db: db,
	}
}

func (s *transferService) ResolveBankDetails(ctx context.Context, accountNumber, accountBank string) ([]byte, error) {
	input := struct {
		AccountNumber string `json:"account_number"`
		AccountBank   string `json:"account_bank"`
	}{
		AccountNumber: accountNumber,
		AccountBank:   accountBank,
	}

	payload, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.flutterwave.com/v3/accounts/resolve", bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+os.Getenv("FLW_SECRET_KEY"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println("Flutterwave response:", string(body))
	return body, nil
}

func (s *transferService) InitializeFunding(ctx context.Context, accountNumber, accountBank string, amount float64, narration string) ([]byte, error) {
	transfer := InternalTransferPayload{
		AccountNumber: accountNumber,
		AccountBank:   accountBank,
		Amount:        amount,
		Narration:     narration,
		Currency:      "NGN",
		DebitCurrency: "NGN",
		Reference:     "TXN-" + time.Now().Format("20060102150405"),
	}

	payload, err := json.Marshal(transfer)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.flutterwave.com/v3/transfers", bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+os.Getenv("FLW_SECRET_KEY"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println("Flutterwave response:", string(body))
	return body, nil
}

func (s *transferService) VerifyFunding(ctx context.Context) error {
	// VerifyFunding logic here
	return nil
}
