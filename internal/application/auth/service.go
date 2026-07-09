package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"payme/internal/application/auth/dto"
	"payme/internal/application/user"
	"payme/internal/application/wallet"
	"payme/internal/config"
	"payme/pkg/utils"

	"github.com/golang-jwt/jwt/v5"

	"gorm.io/gorm"
)

type AuthService interface {
	SignUp(ctx context.Context, req dto.SignUpRequest) (dto.SignUpResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error)
	VerifyOTP(ctx context.Context, req dto.VerifyOTPRequest) (dto.SignUpResponse, error)
	ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest) (dto.ForgotPasswordResponse, error)
	ResetPassword(ctx context.Context, req dto.ResetPasswordRequest, tokenString string) (dto.ResetPasswordResponse, error)
}

type authService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) AuthService {
	return &authService{
		db: db,
	}
}

func (s *authService) SignUp(ctx context.Context, req dto.SignUpRequest) (dto.SignUpResponse, error) {
	// Check if user already exists
	var existingUser user.User
	if err := s.db.WithContext(ctx).Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return dto.SignUpResponse{}, errors.New("email already exists")

	}

	if err := utils.ValidateRegister(req.FullName, req.Email, req.PhoneNumber, req.Password); err != nil {
		return dto.SignUpResponse{}, err
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return dto.SignUpResponse{}, errors.New("error hashing password")
	}
	now := time.Now()

	// --- RATE LIMIT CHECK ---
	var verification user.EmailVerification
	err = s.db.WithContext(ctx).Where("email = ?", req.Email).First(&verification).Error

	if err == nil {

		if now.Sub(verification.FirstOTPRequestAt) < time.Hour {
			// Still inside the 1-hour window
			if verification.OTPRequestCount >= 5 {
				return dto.SignUpResponse{}, errors.New("too many OTP requests, please try again later")
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
		verification.FullName = req.FullName
		verification.Password = hashedPassword

		if err := s.db.WithContext(ctx).Save(&verification).Error; err != nil {
			return dto.SignUpResponse{}, errors.New("could not update verification")
		}
	} else {
		// No record yet — first request ever for this email
		verification = user.EmailVerification{
			Email:             req.Email,
			FullName:          req.FullName,
			Password:          hashedPassword,
			OTP:               utils.GenerateOTP(),
			ExpiresAt:         now.Add(5 * time.Minute),
			OTPRequestCount:   1,
			FirstOTPRequestAt: now,
		}
		if err := s.db.WithContext(ctx).Create(&verification).Error; err != nil {
			return dto.SignUpResponse{}, errors.New("could not create user")
		}
	}

	if err := utils.SendOTPEmail(req.Email, verification.OTP); err != nil {
		return dto.SignUpResponse{}, err
	}

	return dto.SignUpResponse{}, nil

}
func (s *authService) VerifyOTP(ctx context.Context, req dto.VerifyOTPRequest) (dto.SignUpResponse, error) {

	var verification user.EmailVerification

	err := s.db.WithContext(ctx).Where("email = ?", req.Email).First(&verification).Error

	if err != nil {
		return dto.SignUpResponse{},
			errors.New("otp not found")
	}

	if verification.OTP != req.OTP {
		return dto.SignUpResponse{},
			errors.New("invalid otp")
	}

	if time.Now().After(verification.ExpiresAt) {
		return dto.SignUpResponse{},
			errors.New("otp expired")
	}

	user := user.User{
		Email:    verification.Email,
		Password: verification.Password,
		FullName: verification.FullName,
		Role:     "users",
		KYCStatus: "pending",
		Tier:      1,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return dto.SignUpResponse{}, err
	}

	wallet := wallet.Wallet{
		UserID:   user.ID,
		Balance:  0,
		Currency: "NGN",
		Status:   "active",
	}

	if err := s.db.Create(&wallet).Error; err != nil {
		return dto.SignUpResponse{}, err
	}

	// token, err := utils.GenerateToken(
	// 	user.ID,
	// 	user.Role,
	// )

	if err != nil {
		return dto.SignUpResponse{}, err
	}

	s.db.Delete(&verification)

	return dto.SignUpResponse{
		Message: "Account verified",
		UserID:  user.ID,
		// Token:   token,
		FullName: user.FullName,
		Email:  user.Email,
	}, nil
}
func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error) {
	var u user.User

	if err := s.db.WithContext(ctx).Where("email = ?", req.Email).First(&u).Error; err != nil {
		return dto.LoginResponse{}, errors.New("email not found")
	}

	// --- CHECK IF ACCOUNT IS CURRENTLY LOCKED ---
	if u.LockedUntil != nil && time.Now().Before(*u.LockedUntil) {
		remaining := time.Until(*u.LockedUntil).Round(time.Second)
		return dto.LoginResponse{}, fmt.Errorf("account locked, try again in %s", remaining)
	}

	// --- CHECK PASSWORD ---
	err := utils.CheckPassword(u.Password, req.Password)
	if err != nil {
		u.FailedAttempts++

		if u.FailedAttempts >= 3 {
			lockUntil := time.Now().Add(15 * time.Minute)
			u.LockedUntil = &lockUntil
			u.FailedAttempts = 0
			s.db.WithContext(ctx).Save(&u)
			return dto.LoginResponse{}, errors.New("too many failed attempts, account locked for 15 minutes")
		}

		s.db.WithContext(ctx).Save(&u)
		return dto.LoginResponse{}, errors.New("invalid password")
	}

	// --- SUCCESSFUL LOGIN: reset attempts ---
	u.FailedAttempts = 0
	u.LockedUntil = nil
	s.db.WithContext(ctx).Save(&u)

	token, err := utils.GenerateToken(u.ID, u.Role)
	if err != nil {
		return dto.LoginResponse{}, errors.New("could not generate token")
	}

	return dto.LoginResponse{
		Message:  "User logged in successfully",
		Token:    token,
		FullName: u.FullName,
		Email:    u.Email,
	}, nil
}

func (s *authService) ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest) (dto.ForgotPasswordResponse, error) {
	//check if email exist
	var user user.User
	result := s.db.WithContext(ctx).Where("email = ?", req.Email).First(&user)

	if result.RowsAffected == 0 {
		return dto.ForgotPasswordResponse{}, errors.New("email does not exist")
	}

	// ✅ Generate JWT token
	token, err := utils.GeneratePasswordToken(req.Email)
	if err != nil {
		return dto.ForgotPasswordResponse{}, errors.New("could not generate token")
	}
	resetLink := req.FrontendUrl + "/reset-password?token=" + token

	err = utils.SendResetEmail(req.Email, resetLink)
	if err != nil {
		return dto.ForgotPasswordResponse{}, errors.New("could not send email")
	}
	return dto.ForgotPasswordResponse{
		Message: "Check your email",
		Token:   token,
	}, nil
}
func (s *authService) ResetPassword(ctx context.Context, req dto.ResetPasswordRequest, tokenString string) (dto.ResetPasswordResponse, error) {

	// 2️⃣ Validate token (this already checks exp)
	token, err := utils.ValidateToken(tokenString)
	if err != nil || !token.Valid {

		return dto.ResetPasswordResponse{}, errors.New("Invalid or expired token")
	}

	// 3️⃣ Extract claims
	claims := token.Claims.(jwt.MapClaims)
	fmt.Println(claims)

	email, ok := claims["email"].(string)
	if !ok {
		return dto.ResetPasswordResponse{}, errors.New("Invalid token payload")
	}

	// 4️⃣ Find user
	var user user.User
	result := config.DB.Where("email = ?", email).First(&user)
	if result.RowsAffected == 0 {
		return dto.ResetPasswordResponse{}, errors.New("User not found")
	}

	// 5️⃣ Hash new password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return dto.ResetPasswordResponse{}, errors.New("Could not hash password")
	}

	user.Password = hashedPassword
	config.DB.Save(&user)

	return dto.ResetPasswordResponse{
		Message: "Password reset successful",
	}, nil
}
