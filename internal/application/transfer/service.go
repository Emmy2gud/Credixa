package transfer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"payme/internal/application/accounts"
	"payme/internal/application/transaction"
	"payme/internal/application/transfer/dto"
	"payme/internal/application/wallet"

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
	InitializeFunding(ctx context.Context, req dto.CreateTransferRequest) (dto.TransferResponse, error)
	VerifyFunding(ctx context.Context) error
}

type transferService struct {
	db   *gorm.DB


	
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

func (s *transferService) InitializeFunding(ctx context.Context, req dto.CreateTransferRequest) (dto.TransferResponse, error) {
	var t transaction.Transaction
	var w wallet.Wallet
	var virtualAccount accounts.VirtualAccount

	// 1. Fetch sender wallet and idempotency

   if err := s.db.WithContext(ctx).Where("reference = ?", req.IdempotencyKey).First(&t).Error; err != nil {
		return dto.TransferResponse{}, fmt.Errorf("transaction not found: %w", err)
	}

   

	if err := s.db.WithContext(ctx).Where("user_id = ?", req.UserID).First(&w).Error; err != nil {
		return dto.TransferResponse{}, fmt.Errorf("wallet not found: %w", err)
	}

	// 2. Check balance
	if w.Balance < req.Amount {
		return dto.TransferResponse{}, fmt.Errorf("insufficient balance")
	}

	// 3. Determine transfer type
	isInternal := s.db.WithContext(ctx).Where("account_number = ?", req.AccountNumber).First(&virtualAccount).Error == nil
	if isInternal {
		return HandleInternalTransfer(ctx, req, w, virtualAccount,s.db)
	}

	return HandleExternalTransfer(ctx, req, w,s.db)
}


func (s *transferService) VerifyFunding(ctx context.Context) error {
	// VerifyFunding logic here
	return nil
}
