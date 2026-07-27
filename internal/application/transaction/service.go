package transaction

import (
	"context"
	"fmt"
	"payme/internal/application/transaction/dto"
	"payme/pkg/pagination"
	"gorm.io/gorm"
)

type TransactionService interface {
	GetTransactionHistory(ctx context.Context, userID uint, filter dto.TransactionFilterRequest) (pagination.Pagination, error)
	GetTransactionByID(ctx context.Context, id string, userID uint) (Transaction, error)
	GetWalletLogs(ctx context.Context, userID uint, filter dto.TransactionFilterRequest) (pagination.Pagination, error)
}

type transactionService struct {
	db *gorm.DB
}

func NewTransactionService(db *gorm.DB) TransactionService {
	return &transactionService{
		db: db,
	}
}

func (s *transactionService) GetTransactionHistory(ctx context.Context, userID uint, filter dto.TransactionFilterRequest) (pagination.Pagination, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 10
	}

	query := s.db.WithContext(ctx).Model(&Transaction{}).Where("user_id = ?", userID)

	if filter.Category != "" {
		query = query.Where("category = ?", filter.Category)
	}
	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Search != "" {
		searchTerm := "%" + filter.Search + "%"
		query = query.Where("reference LIKE ? OR description LIKE ?", searchTerm, searchTerm)
	}
	if filter.StartDate != "" {
		query = query.Where("created_at >= ?", filter.StartDate)
	}
	if filter.EndDate != "" {
		query = query.Where("created_at <= ?", filter.EndDate)
	}

	var totalRows int64
	if err := query.Count(&totalRows).Error; err != nil {
		return pagination.Pagination{}, fmt.Errorf("failed to count transactions: %w", err)
	}

	offset := (filter.Page - 1) * filter.Limit
	var transactions []Transaction
	if err := query.Order("created_at DESC").Offset(offset).Limit(filter.Limit).Find(&transactions).Error; err != nil {
		return pagination.Pagination{}, fmt.Errorf("failed to fetch transactions: %w", err)
	}

	return pagination.CreatePagination(filter.Page, filter.Limit, totalRows, transactions), nil
}

func (s *transactionService) GetTransactionByID(ctx context.Context, id string, userID uint) (Transaction, error) {
	var tx Transaction
	if err := s.db.WithContext(ctx).Where("(id = ? OR reference = ?) AND user_id = ?", id, id, userID).First(&tx).Error; err != nil {
		return Transaction{}, fmt.Errorf("transaction not found: %w", err)
	}
	return tx, nil
}

func (s *transactionService) GetWalletLogs(ctx context.Context, userID uint, filter dto.TransactionFilterRequest) (pagination.Pagination, error) {
	return s.GetTransactionHistory(ctx, userID, filter)
}

