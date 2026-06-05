package bill_payments

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"payme/internal/api/middleware"
	"payme/internal/application/bill_payments/dto"
	"payme/internal/application/bill_payments/sub-services"

	"payme/pkg/utils"

	"github.com/gorilla/mux"
)

type BillPaymentController struct {
	service BillPaymentService
}

func NewBillPaymentController(service BillPaymentService) *BillPaymentController {
	return &BillPaymentController{service: service}
}

func (h *BillPaymentController) BillerCategories(w http.ResponseWriter, r *http.Request) {
	body, err := h.service.GetBillerCategories(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.JSON(w, http.StatusOK, body)
}

func (h *BillPaymentController) BillerCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categoryID := vars["category"]

	body, err := h.service.GetBillerCategory(r.Context(), categoryID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.JSON(w, http.StatusOK, body)
}

func (h *BillPaymentController) BillCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categoryID := vars["category"]

	body, err := h.service.GetBillCategory(r.Context(), categoryID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// categorybody, _ := io.ReadAll(body)
    // respbody, err := h.service.BillServiceCategory(body, categoryID)
	utils.JSON(w, http.StatusOK, body)
}
// func (h *BillPaymentController)BillCategory(w http.ResponseWriter, r *http.Request) {
// 	// Implementation will go here
// 	vars := mux.Vars(r)
// 	categoryID := vars["category"]
// 	body:=utils.ParseBody(r)

// 	body, err := h.service.BillServiceCategory(body, categoryID)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	utils.JSON(w, http.StatusOK, body)
// }


func (h *BillPaymentController) VerifySubscription(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceID := vars["serviceid"]

	var ElecInput dto.VerifyElectricityRequest
	var TvInput dto.VerifyTvSubscriptionRequest
	switch {
	case sub_services.IsElectricityService(serviceID):
		utils.ParseBody(r, &ElecInput)
	case sub_services.IsTvService(serviceID):
		utils.ParseBody(r, &TvInput)
	}

	respbody, err := h.service.VerifySubscription(r.Context(), serviceID, ElecInput, TvInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(respbody)
	utils.JSON(w, http.StatusOK, respbody)
}

func (h *BillPaymentController) CreateBillPayment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceID := vars["serviceid"]
	variationCode := vars["variationcode"]

	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var (
		airtimeInput dto.CreateBillPaymentAirtimeRequest
		dataInput    dto.CreateBillPaymentDataRequest
		tvInput      dto.ChangeTvRequest
		electricInput dto.ElectricityRequest
	)

	// Decode directly into the right DTO based on serviceID
	switch {
	case sub_services.IsElectricityService(serviceID):
		if err := json.Unmarshal(body, &electricInput); err != nil {
			http.Error(w, "invalid electricity request body", http.StatusBadRequest)
			return
		}
	case sub_services.IsTvService(serviceID):
		if err := json.Unmarshal(body, &tvInput); err != nil {
			http.Error(w, "invalid TV request body", http.StatusBadRequest)
			return
		}
	case sub_services.IsMobileData(serviceID):
		if err := json.Unmarshal(body, &dataInput); err != nil {
			http.Error(w, "invalid data request body", http.StatusBadRequest)
			return
		}
	case sub_services.IsMobileVtu(serviceID):
		if err := json.Unmarshal(body, &airtimeInput); err != nil {
			http.Error(w, "invalid airtime request body", http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, "unsupported service type: "+serviceID, http.StatusBadRequest)
		return
	}
resp, err := h.service.CreateBillPayment(r.Context(), userID, serviceID, variationCode,airtimeInput, dataInput, tvInput, electricInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.JSON(w, http.StatusOK, resp)
}