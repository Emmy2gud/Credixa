package beneficiary

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type BeneficiaryService interface {
	AddBeneficiary(ctx context.Context, userID uint, name, accountNum, bankCode, bankName string) (Beneficiary, error)
	GetUserBeneficiaries(ctx context.Context, userID uint) ([]Beneficiary, error)
	DeleteBeneficiary(ctx context.Context, userID uint, beneficiaryID uint) error
}

type beneficiaryService struct {
	db *gorm.DB
}

func NewBeneficiaryService(db *gorm.DB) BeneficiaryService {
	return &beneficiaryService{db: db}
}

func (s *beneficiaryService) AddBeneficiary(ctx context.Context, userID uint, name, accountNum, bankCode, bankName string) (Beneficiary, error) {
	if accountNum == "" || bankCode == "" || name == "" {
		return Beneficiary{}, errors.New("account number, bank code, and account name are required")
	}

	var existing Beneficiary
	err := s.db.WithContext(ctx).
		Where("user_id = ? AND account_number = ? AND bank_code = ?", userID, accountNum, bankCode).
		First(&existing).Error

	if err == nil {
		return existing, nil // Already saved
	}

	b := Beneficiary{
		UserID:        userID,
		AccountName:   name,
		AccountNumber: accountNum,
		BankCode:      bankCode,
		BankName:      bankName,
	}

	if err := s.db.WithContext(ctx).Create(&b).Error; err != nil {
		return Beneficiary{}, fmt.Errorf("failed to save beneficiary: %w", err)
	}

	return b, nil
}

func (s *beneficiaryService) GetUserBeneficiaries(ctx context.Context, userID uint) ([]Beneficiary, error) {
	var list []Beneficiary
	if err := s.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, fmt.Errorf("failed to list beneficiaries: %w", err)
	}
	return list, nil
}

func (s *beneficiaryService) DeleteBeneficiary(ctx context.Context, userID uint, beneficiaryID uint) error {
	result := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", beneficiaryID, userID).Delete(&Beneficiary{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("beneficiary not found")
	}
	return nil
}
