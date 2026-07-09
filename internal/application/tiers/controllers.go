package tiers

import (
	"net/http"
	"payme/internal/api/middleware"
	"payme/internal/application/tiers/dto"
	"payme/internal/config"
	"payme/pkg/utils"
)

type TierHandler struct {
	service TierService
}

func NewTierHandler(service TierService) *TierHandler {
	return &TierHandler{service: service}
}

func (h *TierHandler) Tier2KycUpload(w http.ResponseWriter, r *http.Request) {
	var input struct {
		BVNNumber string `json:"bvnNumber"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	}

	utils.ParseBody(r, &input)

	userID, ok := middleware.GetUserID(r)

	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	req := dto.Tier2Request{
		BVN:       input.BVNNumber,
		FirstName: input.FirstName,
		LastName:  input.LastName,
	}

	tier2Response, err := h.service.Tier2KycUpload(r.Context(), req, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.JSON(w, http.StatusOK, tier2Response)

}

func (h *TierHandler) BVNCallback(w http.ResponseWriter, r *http.Request) {
	reference := r.URL.Query().Get("reference")
	if reference == "" {
		http.Error(w, "missing reference", http.StatusBadRequest)
		return
	}

	_, err := h.service.RetrieveBvnKycUpload(r.Context(), reference)
	if err != nil {
		http.Redirect(w, r, config.FrontendFailedURL, http.StatusFound)
		return
	}

	http.Redirect(w, r, config.FrontendSuccessURL, http.StatusFound)
}

// func (h *TierHandler) Tier3KycUpload(w http.ResponseWriter, r *http.Request) {
// 	var input struct {
// 		NIN        string `json:"nin"`
// 		CustomerID string `json:"customer_id"`
// 		Street     string `json:"street"`
// 		LgaName    string `json:"lgaName"`
// 		StateName  string `json:"stateName"`
// 		City       string `json:"city"`
// 		Landmark   string `json:"landmark"`
// 		FirstName  string `json:"firstName"`
// 		LastName   string `json:"lastName"`
// 		Phone      string `json:"phone"`
// 		DOB        string `json:"dob"`
// 	}

// 	utils.ParseBody(r, &input)

// 	userID, ok := middleware.GetUserID(r)

// 	if !ok {
// 		http.Error(w, "unauthorized", http.StatusUnauthorized)
// 		return
// 	}
// 	switch {
// 	case input.NIN != "":
// 		{
// 			req := dto.Tier3NinRequest{
// 				NIN:       input.NIN,
// 				FirstName: input.FirstName,
// 				LastName:  input.LastName,
// 				DOB:       input.DOB,
// 			}
// 			response, err := h.service.Tier3KycNinUpload(r.Context(), req, userID)
// 			if err != nil {
// 				http.Error(w, err.Error(), http.StatusBadRequest)
// 				return
// 			}
// 			utils.JSON(w, http.StatusOK, response)
// 		}
// 	case input.Street != "":
// 		{

// 			req := dto.Tier3AddressRequest{

// 				Street:             input.Street,
// 				LgaName:            input.LgaName,
// 				StateName:          input.StateName,
// 				City:               input.City,
// 				Landmark:           input.Landmark,
// 				ApplicantFirstName: input.FirstName,
// 				ApplicantLastName:  input.LastName,
// 				ApplicantPhone:     input.Phone,
// 				ApplicantDOB:       input.DOB,
// 			}
// 			response, err := h.service.Tier3KycAddressUpload(r.Context(), req, userID)
// 			if err != nil {
// 				http.Error(w, err.Error(), http.StatusBadRequest)
// 				return
// 			}
// 			utils.JSON(w, http.StatusOK, response)
// 		}
// 	default:
// 		http.Error(w, "either nin or street is required", http.StatusBadRequest)
// 		return
// 	}

// }


func (h *TierHandler) Tier3KycUpload(w http.ResponseWriter, r *http.Request) {
	var input dto.Tier3Request
	utils.ParseBody(r, &input)

	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	response, err := h.service.Tier3Verification(r.Context(), input, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.JSON(w, http.StatusOK, response)
}