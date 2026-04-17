package auth

import (
	"encoding/json"
	"net/http"
	"payme/pkg/config"
	"payme/pkg/user"
	"payme/pkg/utils"
	"payme/pkg/wallet"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var u user.User
	utils.ParseBody(r, &u)

	// Check if user already exists
	var existingUser user.User
	if err := config.DB.Where("email = ?", u.Email).First(&existingUser).Error; err == nil {
		http.Error(w, "Email already exists", http.StatusBadRequest)
		return
	}

	if err := utils.ValidateRegister(u.FullName, u.Email, u.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	u.Password = hashedPassword
	u.Role = "users"

	if err := config.DB.Create(&u).Error; err != nil {
		http.Error(w, "Could not create user", http.StatusInternalServerError)
		return
	}

	wlt := wallet.Wallet{
		UserID:   u.ID,
		Balance:  0,
		Currency: "NGN",
		Status:   "active",
	}

	if err := config.DB.Create(&wlt).Error; err != nil {
		http.Error(w, "Could not create wallet", http.StatusInternalServerError)
		return
	}

	// ✅ Generate JWT token
	token, err := utils.GenerateToken(u.ID, u.Role)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message":      "User registered successfully",
		"access_token": token,
	})

}

func Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	utils.ParseBody(r, &input)
	var user user.User

	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	err := utils.CheckPassword(user.Password, input.Password)
	if err != nil {
		http.Error(w, "Error hashing password ", http.StatusInternalServerError)
		return
	}

	// ✅ Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":      "User Logged successfully",
		"access_token": token,
	})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Logged out successfully. Please clear your auth token.",
	})
}
func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	// ForgotPassword logic here
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	// ResetPassword logic here
}

func VerifyEmail(w http.ResponseWriter, r *http.Request) {
	// VerifyEmail logic here
}
