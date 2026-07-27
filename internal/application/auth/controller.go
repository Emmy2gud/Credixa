package auth

import (
	"encoding/json"

	"net/http"
	"strings"

	"payme/internal/application/auth/dto"

	"payme/pkg/utils"
)

type AuthHandler struct {
	service AuthService
}

func NewAuthHandler(service AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {

	var u dto.SignUpRequest
	utils.ParseBody(r, &u)

	resp, err := h.service.SignUp(r.Context(), u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)

}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input dto.LoginRequest
	utils.ParseBody(r, &input)
	resp, err := h.service.Login(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.JSON(w,http.StatusOK,resp)
}

func Logout(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Logged out successfully. Please clear your auth token.",
	})
}

func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var fp dto.ForgotPasswordRequest
	utils.ParseBody(r, &fp)
    resp,err := h.service.ForgotPassword(r.Context(),fp)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
	utils.JSON(w,http.StatusCreated,resp)
}

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
    var req dto.ResetPasswordRequest

	utils.ParseBody(r, &req)
	// 1️⃣ Get token from Authorization header

	auth := r.Header.Get("Authorization")
	if auth == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
		

	tokenString := strings.Replace(auth, "Bearer ", "", 1)


	resp,err:=h.service.ResetPassword(r.Context(),req,tokenString)
	if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }




	utils.JSON(w,http.StatusOK,resp)
}

func  (h *AuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	var u dto.VerifyOTPRequest
	utils.ParseBody(r, &u)

	resp, err := h.service.VerifyOTP(r.Context(), u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
