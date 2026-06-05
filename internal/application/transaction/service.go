package transaction

import (
	"context"
	"gorm.io/gorm"
)

type TransactionService interface {
	GetTransactionHistory(ctx context.Context, userID uint) ([]Transaction, error)
	GetTransactionByID(ctx context.Context, id string, userID uint) (Transaction, error)
	GetWalletLogs(ctx context.Context, userID uint) ([]Transaction, error)
}

type transactionService struct {
	db *gorm.DB
}

func NewTransactionService(db *gorm.DB) TransactionService {
	return &transactionService{
		db: db,
	}
}

func (s *transactionService) GetTransactionHistory(ctx context.Context, userID uint) ([]Transaction, error) {
	return nil, nil
}

func (s *transactionService) GetTransactionByID(ctx context.Context, id string, userID uint) (Transaction, error) {
	return Transaction{}, nil
}

func (s *transactionService) GetWalletLogs(ctx context.Context, userID uint) ([]Transaction, error) {
	return nil, nil
}
