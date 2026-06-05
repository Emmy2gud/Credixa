package accounts

import (

	"net/http"

	"payme/internal/api/middleware"
	"payme/internal/application/accounts/dto"

	"payme/pkg/utils"
)

type VirtualAccountController struct {
	service VirtualAccountService
}

func NewVirtualAccountController(service VirtualAccountService) *VirtualAccountController {
	return &VirtualAccountController{service: service}
}

func (h *VirtualAccountController) CreateVirtualAccount(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateVirtualAccountInput
	
	utils.ParseBody(r, &input)
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	respbody, err := h.service.CreateVirtualAccount(r.Context(), input, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.JSON(w,http.StatusCreated,respbody)

}