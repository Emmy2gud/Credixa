package transfer

import (
	"context"

	"fmt"

	"payme/internal/adapters/flutterwave"
	"payme/internal/application/accounts"
	"payme/internal/application/transaction"
	
	"payme/internal/application/transfer/dto"
	"payme/internal/application/wallet"

	"time"

	"gorm.io/gorm"
)



func HandleInternalTransfer(ctx context.Context,req dto.CreateTransferRequest,w wallet.Wallet,virtualAccount accounts.VirtualAccount,db *gorm.DB) (dto.TransferResponse, error) {

	ref := "Credix-" + time.Now().Format("20060102150405")

	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// Deduct sender
		if err := wallet.DeductWalletBalance(uint(req.UserID), req.Amount); err != nil {
			return fmt.Errorf("deduct sender: %w", err)
		}

		// Credit receiver
		if err := wallet.UpdateWalletBalance(uint(virtualAccount.UserID), req.Amount); err != nil {
			return fmt.Errorf("credit receiver: %w", err)
		}

		// Create transaction record
		t := transaction.Transaction{
			UserID:      req.UserID,
			WalletID:    w.ID,
			Type:        "debit",
			Category:    "transfer",
			Amount:      req.Amount,
			Fee:         50,
			Reference:   ref,
			Status:      "completed", 
			Description: req.Narration,
		}
		if err := tx.Create(&t).Error; err != nil {
			return fmt.Errorf("create transaction: %w", err)
		}

		// Create transfer record
		transfer :=Transfer{
			SenderID:      req.UserID,
			ReceiverID:    &virtualAccount.UserID,
			Amount:        req.Amount, 
			Reference:     ref,
			Status:        "completed",
			TransferType:  "internal",
			AccountNumber: req.AccountNumber,
			AccountName:   req.Firstname + " " + req.Lastname,
			Narration:     req.Narration,
		}
		if err := tx.Create(&transfer).Error; err != nil {
			return fmt.Errorf("create transfer: %w", err)
		}

		return nil
	})

	if err != nil {
		return dto.TransferResponse{}, err
	}

	// Internal transfers don't go through Flutterwave — return early
	return dto.TransferResponse{
		AccountNumber: req.AccountNumber,
		BankCode:      req.AccountBank,
		FullName:      req.Firstname + " " + req.Lastname,
		Amount:        req.Amount,
		Status:        "completed",
		Reference:     ref,
		Narration:     req.Narration,
	}, nil
}

func HandleExternalTransfer(ctx context.Context,req dto.CreateTransferRequest,w wallet.Wallet,db *gorm.DB) (dto.TransferResponse, error) {

	ref := "FLWX-" + time.Now().Format("20060102150405")

	// Deduct before sending to Flutterwave
	if err := wallet.DeductWalletBalance(uint(req.UserID), req.Amount); err != nil {
		return dto.TransferResponse{}, fmt.Errorf("deduct sender: %w", err)
	}

	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		t := transaction.Transaction{
			UserID:      req.UserID,
			WalletID:    w.ID,
			Type:        "debit",
			Category:    "transfer",
			Amount:      req.Amount,
			Fee:         50,
			Reference:   ref,
			Status:      "pending",
			Description: req.Narration,
		}
		if err := tx.Create(&t).Error; err != nil {
			return fmt.Errorf("create transaction: %w", err)
		}

		transfer := Transfer{
			SenderID:      req.UserID,
			Amount:        req.Amount,
			Reference:     ref,
			Status:        "pending",
			TransferType:  "external",
			AccountNumber: req.AccountNumber,
			BankCode:      &req.AccountBank,
			AccountName:   req.Firstname + " " + req.Lastname,
			Narration:     req.Narration,
		}
		if err := tx.Create(&transfer).Error; err != nil {
			return fmt.Errorf("create transfer: %w", err)
		}

		return nil
	})


	// Call Flutterwave
	flwResp, err := adapters.NewClient().CreateTransfers(ctx, req)
	if err != nil {
		// TODO: mark transfer as failed in DB, refund wallet
		return dto.TransferResponse{}, err
	}

	if flwResp.Status != "SUCCESS" {
		// TODO: mark transfer as failed in DB, refund wallet
		return dto.TransferResponse{}, fmt.Errorf("flutterwave error: %s", flwResp.Message)
	}

	return dto.TransferResponse{
		ID:               flwResp.Data.ID,
		AccountNumber:    flwResp.Data.AccountNumber,
		BankCode:         flwResp.Data.BankCode,
		FullName:         flwResp.Data.FullName,
		CreatedAt:        flwResp.Data.CreatedAt,
		Currency:         flwResp.Data.Currency,
		DebitCurrency:    flwResp.Data.DebitCurrency,
		Amount:           flwResp.Data.Amount,
		Fee:              flwResp.Data.Fee,
		Status:           flwResp.Data.Status,
		Reference:        flwResp.Data.Reference,
		Meta:             flwResp.Data.Meta,
		Narration:        flwResp.Data.Narration,
		CompleteMessage:  flwResp.Data.CompleteMessage,
		RequiresApproval: flwResp.Data.RequiresApproval,
		IsApproved:       flwResp.Data.IsApproved,
		BankName:         flwResp.Data.BankName,
	}, nil
}