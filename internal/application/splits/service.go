package splits

import (
	"context"
	"fmt"
	"payme/internal/application/notifications"
	"payme/internal/application/wallet"

	"gorm.io/gorm"
)

type ParticipantInput struct {
	UserID     uint64 `json:"user_id"`
	Amount     int64 `json:"amount"`
	Percentage uint64 `json:"percentage"`
}

type CreateSplitBillInput struct {
	Title             string             `json:"title"`
	Description       string             `json:"description"`
	TotalAmount       int64             `json:"total_amount"`
	SplitType         string             `json:"split_type"`
	ParticipantsCount int                `json:"participants_count"`
	Participants      []ParticipantInput `json:"participants"`
}

type SplitService interface {
	CreateSplitBill(ctx context.Context, input CreateSplitBillInput, userID uint) (uint64, error)
	GetSplitBills(ctx context.Context, userID uint) ([]SplitBill, error)
	AcceptSplitBill(ctx context.Context, splitID uint64, userID uint) error
	DeclineSplitBill(ctx context.Context, splitID uint64, userID uint) error
}

type splitService struct {
	db *gorm.DB
}

func NewSplitService(db *gorm.DB) SplitService {
	return &splitService{db: db}
}

func (s *splitService) CreateSplitBill(ctx context.Context, input CreateSplitBillInput, userID uint) (uint64, error) {
	splitbill := SplitBill{
		CreatorID:         uint64(userID),
		Title:             input.Title,
		Description:       input.Description,
		TotalAmount:       input.TotalAmount,
		SplitType:         input.SplitType,
		ParticipantsCount: input.ParticipantsCount,
		Status:            "pending",
	}

	if err := s.db.WithContext(ctx).Create(&splitbill).Error; err != nil {
		return 0, fmt.Errorf("could not create split bill: %w", err)
	}

	var creatorAmount, creatorPercentage int64
	for _, p := range input.Participants {
		if p.UserID == uint64(userID) {
			creatorAmount = p.Amount
			creatorPercentage = int64(p.Percentage)
			break
		}
	}

	creatorParticipant := SplitBillParticipants{
		SplitBillID: splitbill.ID,
		UserID:      uint64(userID),
		Amount:      creatorAmount,
		Percentage:  uint64(creatorPercentage),
		Status:      "accepted",
	}
	if err := s.db.WithContext(ctx).Create(&creatorParticipant).Error; err != nil {
		return 0, fmt.Errorf("could not add creator as participant: %w", err)
	}

	for _, participant := range input.Participants {
		if participant.UserID == uint64(userID) {
			continue
		}

		invitedParticipant := SplitBillParticipants{
			SplitBillID: splitbill.ID,
			UserID:      participant.UserID,
			Amount:      participant.Amount,
			Percentage:  participant.Percentage,
			Status:      "pending",
		}
		if err := s.db.WithContext(ctx).Create(&invitedParticipant).Error; err != nil {
			fmt.Printf("Could not add participant %d: %v\n", participant.UserID, err)
			continue
		}

		inviteMsg := fmt.Sprintf(
			"You have been invited to a split bill: \"%s\". Please accept or decline.",
			splitbill.Title,
		)
		notifications.CreateNotification(uint(participant.UserID), "SplitBill", "Split Bill Invitation", inviteMsg)
	}

	creatorMsg := fmt.Sprintf(
		"Your split bill \"%s\" has been created successfully. Waiting for participants to respond.",
		splitbill.Title,
	)
	notifications.CreateNotification(userID, "SplitBill", "Split Bill Created", creatorMsg)

	return splitbill.ID, nil
}

func (s *splitService) GetSplitBills(ctx context.Context, userID uint) ([]SplitBill, error) {
	var splitBills []SplitBill
	if err := s.db.WithContext(ctx).Where("creator_id = ?", userID).Find(&splitBills).Error; err != nil {
		return nil, err
	}
	return splitBills, nil
}

func (s *splitService) AcceptSplitBill(ctx context.Context, splitID uint64, userID uint) error {
	var participant SplitBillParticipants
	if err := s.db.WithContext(ctx).
		Where("split_bill_id = ? AND user_id = ?", splitID, userID).
		First(&participant).Error; err != nil {
		return fmt.Errorf("participant record not found: %w", err)
	}

	participant.Status = "accepted"
	if err := s.db.WithContext(ctx).Save(&participant).Error; err != nil {
		return fmt.Errorf("could not update participant status: %w", err)
	}

	var splitBill SplitBill
	if err := s.db.WithContext(ctx).Where("id = ?", splitID).First(&splitBill).Error; err != nil {
		return fmt.Errorf("failed to load split bill: %w", err)
	}

	if err := wallet.DeductWalletBalance(uint(participant.UserID), participant.Amount); err != nil {
		return fmt.Errorf("failed to deduct participant wallet balance: %w", err)
	}

	if err := wallet.UpdateWalletBalance(uint(splitBill.CreatorID), participant.Amount); err != nil {
		return fmt.Errorf("failed to credit creator wallet: %w", err)
	}

	acceptMsg := fmt.Sprintf(
		"A participant has accepted your split bill: \"%s\".",
		splitBill.Title,
	)
	if _, err := notifications.CreateNotification(uint(splitBill.CreatorID), "SplitBill", "Split Bill Accepted", acceptMsg); err != nil {
		fmt.Printf("failed to send acceptance notification: %v\n", err)
	}

	return nil
}

func (s *splitService) DeclineSplitBill(ctx context.Context, splitID uint64, userID uint) error {
	var participant SplitBillParticipants
	if err := s.db.WithContext(ctx).
		Where("split_bill_id = ? AND user_id = ?", splitID, userID).
		First(&participant).Error; err != nil {
		return fmt.Errorf("participant record not found: %w", err)
	}

	participant.Status = "declined"
	if err := s.db.WithContext(ctx).Save(&participant).Error; err != nil {
		return fmt.Errorf("could not update participant status: %w", err)
	}

	var splitBill SplitBill
	if err := s.db.WithContext(ctx).Where("id = ?", splitID).First(&splitBill).Error; err != nil {
		return fmt.Errorf("failed to load split bill: %w", err)
	}

	declineMsg := fmt.Sprintf(
		"A participant has declined your split bill: \"%s\".",
		splitBill.Title,
	)
	if _, err := notifications.CreateNotification(uint(splitBill.CreatorID), "SplitBill", "Split Bill Declined", declineMsg); err != nil {
		fmt.Printf("failed to send decline notification: %v\n", err)
	}

	return nil
}
