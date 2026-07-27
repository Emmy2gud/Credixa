package transfer

import (
	"context"

	"fmt"

	"payme/internal/adapters/flutterwave"


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
	ResolveBankDetails(ctx context.Context, input dto.ResolveBankDetailsRequest) (dto.ResolveBankDetailsResponse, error)
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

func (s *transferService) ResolveBankDetails(ctx context.Context, input dto.ResolveBankDetailsRequest) (dto.ResolveBankDetailsResponse, error) {


	flwResp, err := adapters.NewClient().ResolveAccountDetails(ctx, input)
	if err != nil {
		return dto.ResolveBankDetailsResponse{}, err
	}
    return dto.ResolveBankDetailsResponse{
        AccountName: flwResp.Data.AccountName,
        AccountNumber: flwResp.Data.AccountNumber,
        Status: flwResp.Status,
        Message: flwResp.Message,
	},nil
}

func (s *transferService) InitializeFunding(ctx context.Context, req dto.CreateTransferRequest) (dto.TransferResponse, error) {
	var t transaction.Transaction
	var w wallet.Wallet
	var virtualAccount accounts.VirtualAccount

	// 1. Fetch sender wallet and idempotency

   if err := s.db.WithContext(ctx).Where("idempotency_key = ?", req.IdempotencyKey).First(&t).Error; err != nil {
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
