package transactionpin

import (
	"context"

	"payme/internal/application/transactionpin/dto"
	"payme/pkg/utils"
	"errors"

	"gorm.io/gorm"
)

type TransactionPinService interface {
	SetTransactionPin(ctx context.Context,req dto.SetTransactionPinRequest,userID uint) (dto.SetTransactionPinResponse,error)
}

type transactionPinService struct {
	db *gorm.DB	
}

func NewTransactionPinService(db *gorm.DB) TransactionPinService {
	return &transactionPinService{
		db: db,
	}
}

func (s *transactionPinService) SetTransactionPin(ctx context.Context,req dto.SetTransactionPinRequest,userID uint) (dto.SetTransactionPinResponse,error) {
		// 1️⃣ Get user ID from context
		var tp TransactionPin


	//pin length is 4 
	if len(req.Pin) != 4 {
		return dto.SetTransactionPinResponse{}, errors.New("pin must be 4 digits")
	}
	pin, err := utils.HashPassword(req.Pin)
	if err != nil {
		return dto.SetTransactionPinResponse{}, errors.New("error hashing password")
	}

	tp.UserID = userID

	tp.Pin = pin

	s.db.WithContext(ctx).Create(&tp)

	return dto.SetTransactionPinResponse{
		Message: "Transaction pin created successfully",
	},nil
}