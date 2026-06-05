package auth

import (
	"context"
	"errors"
	"fmt"


	"payme/internal/application/auth/dto"
	"payme/internal/application/user"
	"payme/internal/application/wallet"
	"payme/internal/config"
	"payme/pkg/utils"

	"github.com/golang-jwt/jwt/v5"

	"gorm.io/gorm"
)

type AuthService interface {
	SignUp(ctx context.Context,req dto.SignUpRequest) (dto.SignUpResponse,error)
	Login(ctx context.Context,req dto.LoginRequest)(dto.LoginResponse,error)
	
	ForgotPassword(ctx context.Context,req dto.ForgotPasswordRequest) (dto.ForgotPasswordResponse,error)
	ResetPassword(ctx context.Context,req dto.ResetPasswordRequest,tokenString string) (dto.ResetPasswordResponse,error)
}

type authService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) AuthService {
	return &authService{
		db: db,
	}
}

func (s *authService) SignUp(ctx context.Context,req dto.SignUpRequest) (dto.SignUpResponse,error) {
	// Check if user already exists
	var existingUser user.User
	if err := s.db.WithContext(ctx).Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return dto.SignUpResponse{}, errors.New("email already exists")
	
	}

	if err := utils.ValidateRegister(req.FullName, req.Email, req.Password); err != nil {
		return dto.SignUpResponse{}, err
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return dto.SignUpResponse{}, errors.New("error hashing password")
	}
	user := user.User{
		Email:    req.Email,
		Password: hashedPassword,
		FullName: req.FullName,
		Role:     "users",
	}

	if err := s.db.WithContext(ctx).Create(&user).Error; err != nil {
		return dto.SignUpResponse{}, errors.New("could not create user")
	}

	wlt := wallet.Wallet{
		UserID:   user.ID,
		Balance:  0,
		Currency: "NGN",
		Status:   "active",
	}

	if err := s.db.WithContext(ctx).Create(&wlt).Error; err != nil {
		return dto.SignUpResponse{}, errors.New("could not create wallet")
		
	}

	// ✅ Generate JWT token
	token, err := utils.GenerateToken(user.ID, req.Role)
	if err != nil {
		return dto.SignUpResponse{}, errors.New("could not generate token")
	}

	return dto.SignUpResponse{
		Message: "User registered successfully",
		UserID:  user.ID,
		Token:   token,
	}, nil

}


func (s *authService) Login(ctx context.Context,req dto.LoginRequest) (dto.LoginResponse,error) {
	var user user.User

	if err := s.db.WithContext(ctx).Where("email = ?", req.Email).First(&user).Error; err != nil {
		return dto.LoginResponse{}, errors.New("email not found")
	}

	err := utils.CheckPassword(user.Password, req.Password)
	if err != nil {
		return dto.LoginResponse{}, errors.New("error hashing password")
		
	}

	// ✅ Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return dto.LoginResponse{}, errors.New("could not generate token")
	}
	
	return dto.LoginResponse{
		Message: "User logged in successfully",
		Token:   token,
	}, nil

}

func (s *authService) ForgotPassword(ctx context.Context,req dto.ForgotPasswordRequest) (dto.ForgotPasswordResponse,error) {
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
func (s *authService) ResetPassword(ctx context.Context,req dto.ResetPasswordRequest,tokenString string) (dto.ResetPasswordResponse,error) {


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
