package user

import (
	"fmt"
	"net/http"

	"payme/internal/application/user/dto"
	"payme/pkg/utils"
)

type UserHandler struct {
	service UserService
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateUserProfileRequest
	utils.ParseBody(r, &req)
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}
	// 2 Parse multipart form (image + fields)
	err := r.ParseMultipartForm(20 << 20) // 20MB
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}
	// Get uploaded file
	file, handler, err := r.FormFile("profile_picture")
	if err != nil {
		http.Error(w, "image is required", http.StatusBadRequest)
		return
	}
	defer file.Close()
	fmt.Println("this user is :", userID)
	resp, err := h.service.UpdateUserProfile(r.Context(), req, userID, file, handler)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.JSON(w, http.StatusOK, resp)
}
