package transactionpin

import (
	"context"
	"time"

	"errors"
	
	"payme/internal/application/transactionpin/dto"
	"payme/internal/application/user"
	"payme/pkg/utils"

	"gorm.io/gorm"
)

type TransactionPinService interface {
	SetTransactionPin(ctx context.Context, req dto.SetTransactionPinRequest, userID uint) (dto.SetTransactionPinResponse, error)
	UpdateTransactionPin(ctx context.Context, req dto.UpdateTransactionPin, userID uint) (dto.SetTransactionPinResponse, error)
	ForgotTransactionPin(ctx context.Context, req dto.SetTransactionPinRequest, userID uint) (dto.SetTransactionPinResponse, error)
	ResetTransactionPin(ctx context.Context, req dto.VerifyTransactionPinOTPRequest, userID uint) (dto.SetTransactionPinResponse, error)
}

type transactionPinService struct {
	db *gorm.DB
}

func NewTransactionPinService(db *gorm.DB) TransactionPinService {
	return &transactionPinService{
		db: db,
	}
}

func (s *transactionPinService) SetTransactionPin(ctx context.Context, req dto.SetTransactionPinRequest, userID uint) (dto.SetTransactionPinResponse, error) {
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
	}, nil
}

func (s *transactionPinService) UpdateTransactionPin(ctx context.Context, req dto.UpdateTransactionPin, userID uint) (dto.SetTransactionPinResponse, error) {
	// 1️⃣ Get user ID from context
	var tp TransactionPin

	//pin length is 4
	if len(req.OldPin) != 4 {
		return dto.SetTransactionPinResponse{}, errors.New("old pin must be 4 digits")
	}
	if len(req.NewPin) != 4 {
		return dto.SetTransactionPinResponse{}, errors.New("new pin must be 4 digits")
	}
	pin, err := utils.HashPassword(req.NewPin)
	if err != nil {
		return dto.SetTransactionPinResponse{}, errors.New("error hashing password")
	}
	// checking if the old pin is correct
	s.db.WithContext(ctx).Where("user_id = ?", userID).First(&tp)
	if tp.Pin != req.OldPin {
		return dto.SetTransactionPinResponse{}, errors.New("old pin is incorrect")
	}
	tp.UserID = userID

	tp.Pin = pin

	s.db.WithContext(ctx).Create(&tp)
	s.db.WithContext(ctx).Model(&tp).Where("user_id = ?", userID).Update("pin", pin)
	return dto.SetTransactionPinResponse{
		Message: "Transaction pin updated successfully",
	}, nil
}

func (s *transactionPinService) ForgotTransactionPin(ctx context.Context, req dto.SetTransactionPinRequest, userID uint) (dto.SetTransactionPinResponse, error) {
	//getting the user in other to get there email
	var u user.User
	s.db.WithContext(ctx).Where("id = ?", userID).First(&u)

	//validate and hash pin
	if len(req.Pin) != 4 {
		return dto.SetTransactionPinResponse{}, errors.New("pin must be 4 digits")
	}
	pin, err := utils.HashPassword(req.Pin)
	if err != nil {
		return dto.SetTransactionPinResponse{}, errors.New("error hashing password")
	}
	//check if otp verification already exist
	// --- RATE LIMIT CHECK ---
	var verification user.OtpVerification
	err = s.db.WithContext(ctx).Where("email = ? AND purpose = ?", u.Email, user.PurposeSignup).First(&verification).Error
	now := time.Now()
	//if there is an emailverification record for the user
	if err == nil {

		if now.Sub(verification.FirstOTPRequestAt) < time.Hour {
			// Still inside the 1-hour window
			if verification.OTPRequestCount >= 5 {
				return dto.SetTransactionPinResponse{}, errors.New("too many OTP requests, please try again later")
			}
			// under the limit → bump count
			verification.OTPRequestCount++
		} else {
			// window expired → reset
			verification.FirstOTPRequestAt = now
			verification.OTPRequestCount = 1
		}

		// update OTP + expiry + latest details
		verification.OTP = utils.GenerateOTP()
		verification.ExpiresAt = now.Add(5 * time.Minute)
		verification.FullName = u.FullName
		verification.Password = pin

		if err := s.db.WithContext(ctx).Save(&verification).Error; err != nil {
			return dto.SetTransactionPinResponse{}, errors.New("could not update verification")
		}
	} else {
		// No record yet — first request ever for this email
		verification = user.OtpVerification{
			Email:             u.Email,
			FullName:          u.FullName,
			Password:          pin,
			Purpose:           user.PurposePinReset,
			OTP:               utils.GenerateOTP(),
			ExpiresAt:         now.Add(5 * time.Minute),
			OTPRequestCount:   1,
			FirstOTPRequestAt: now,
		}
		if err := s.db.WithContext(ctx).Create(&verification).Error; err != nil {
			return dto.SetTransactionPinResponse{}, errors.New("could not create user")
		}
	}

	if err := utils.SendOTPEmail(u.Email, verification.OTP); err != nil {
		return dto.SetTransactionPinResponse{}, err
	}

	return dto.SetTransactionPinResponse{}, nil

}

func (s *transactionPinService) ResetTransactionPin(ctx context.Context, req dto.VerifyTransactionPinOTPRequest, userID uint) (dto.SetTransactionPinResponse, error) {
	var verification user.OtpVerification

	err := s.db.WithContext(ctx).Where("password = ? AND purpose = ?", req.Pin, user.PurposePinReset).First(&verification).Error

	if err != nil {
		return dto.SetTransactionPinResponse{},
			errors.New("otp not found")
	}

	if verification.OTP != req.OTP {
		return dto.SetTransactionPinResponse{},
			errors.New("invalid otp")
	}

	if time.Now().After(verification.ExpiresAt) {
		return dto.SetTransactionPinResponse{},
			errors.New("otp expired")
	}

	tp := TransactionPin{
		UserID: userID,
		Pin:    verification.Password,
	}
	if err := s.db.Create(&tp).Error; err != nil {
		return dto.SetTransactionPinResponse{}, err
	}

	s.db.Delete(&verification)

	return dto.SetTransactionPinResponse{
		Message: "Transaction Pin Change Successfully"}, nil

}