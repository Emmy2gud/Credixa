package accounts

import (
	"context"
	"log"

	"fmt"

	"payme/internal/adapters/flutterwave"
	"payme/internal/application/accounts/dto"
	"payme/internal/application/notifications"
	"payme/internal/application/wallet"

	"gorm.io/gorm"
)

type VirtualAccountService interface {
	CreateVirtualAccount(ctx context.Context, input dto.CreateVirtualAccountRequest) (dto.VirtualAccountResponse, error)
}

type virtualAccountService struct {
	db *gorm.DB
	notificationService notifications.NotificationService
}

func NewVirtualAccountService(db *gorm.DB) VirtualAccountService {
	return &virtualAccountService{
		db: db,
		notificationService: notifications.NewNotificationService(db),
	}
}

func (s *virtualAccountService) CreateVirtualAccount(ctx context.Context, req dto.CreateVirtualAccountRequest) (dto.VirtualAccountResponse, error) {

	var w wallet.Wallet

	if err := s.db.WithContext(ctx).Where("user_id = ?", req.UserID).First(&w).Error; err != nil {
		log.Printf("wallet not found: %v", err)
		return dto.VirtualAccountResponse{},
			fmt.Errorf("wallet not found: %w", err)
	}
	
	flwResp, err := adapters.NewClient().CreateVirtualAccount(ctx, req)

	if err != nil {
		return dto.VirtualAccountResponse{}, err
	}

	if flwResp.Data.ResponseCode != "02" {
		log.Printf("Flutterwave error: %s", flwResp.Data.ResponseMessage)
		return dto.VirtualAccountResponse{},
			fmt.Errorf(
				"flutterwave error: %s",
				flwResp.Data.ResponseMessage,
			)
	}

	virtualAcc := VirtualAccount{
		WalletID:      w.ID,
		AccountNumber: flwResp.Data.AccountNumber,
		AccountName:   req.Firstname + " " + req.Lastname,
		BankName:      flwResp.Data.BankName,
		Provider:      "flutterwave",
		Status:        flwResp.Status,
	}

	if err := s.db.WithContext(ctx).
		Create(&virtualAcc).Error; err != nil {

		return dto.VirtualAccountResponse{}, err
	}

	s.notificationService.CreateNotification(ctx, req.UserID, "Virtual Account Created", "wallet-fund", "Virtual Account Created", "")
	return dto.VirtualAccountResponse{
		ResponseCode:    flwResp.Data.ResponseCode,
		ResponseMessage: flwResp.Data.ResponseMessage,
		FlwRef:          flwResp.Data.FlwRef,
		OrderRef:        flwResp.Data.OrderRef,
		AccountNumber:   flwResp.Data.AccountNumber,
		Frequency:       flwResp.Data.Frequency,
		BankName:        flwResp.Data.BankName,
		CreatedAt:       flwResp.Data.CreatedAt,
		ExpiryDate:      flwResp.Data.ExpiryDate,
		Note:            flwResp.Data.Note,
		Amount:          flwResp.Data.Amount,
	}, nil
}
