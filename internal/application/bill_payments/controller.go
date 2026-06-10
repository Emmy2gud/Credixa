package bill_payments

import (
	"encoding/json"
	"fmt"

	"net/http"
	"payme/internal/api/middleware"

	"payme/internal/application/bill_payments/dto"
	"payme/internal/application/bill_payments/services"
	sub_services "payme/internal/application/bill_payments/sub-services"

	"payme/pkg/utils"

	"github.com/gorilla/mux"
)

type BillPaymentController struct {
	service services.BillPaymentService
}

func NewBillPaymentController(service services.BillPaymentService) *BillPaymentController {

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

	utils.JSON(w, http.StatusOK, body)
}

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
	requestID, err := utils.GenerateRequestID()
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	ctx := r.Context()

	switch {
	case sub_services.IsElectricityService(serviceID):
		var req dto.ElectricityRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid electricity request body", http.StatusBadRequest)
			return
		}
		req.ServiceId = serviceID
		req.VariationCode = variationCode
		req.RequestID = requestID

		resp, err := h.service.ProcessElectricity(ctx, userID, req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		utils.JSON(w, http.StatusOK, resp)

	case sub_services.IsTvService(serviceID):
		var req dto.ChangeTvRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid TV request body", http.StatusBadRequest)
			return
		}
		req.ServiceId = serviceID
		req.VariationCode = variationCode
		req.RequestID = requestID
		resp, err := h.service.ProcessTV(ctx, userID, req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		utils.JSON(w, http.StatusOK, resp)

	case sub_services.IsMobileData(serviceID):
		var req dto.CreateBillPaymentDataRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid data request body", http.StatusBadRequest)
			return
		}
		req.ServiceId = serviceID
		req.VariationCode = variationCode
		req.RequestID = requestID
		resp, err := h.service.ProcessData(ctx, userID, req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		utils.JSON(w, http.StatusOK, resp)

	case sub_services.IsMobileVtu(serviceID):
		var req dto.CreateBillPaymentAirtimeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid airtime request body", http.StatusBadRequest)
			return
		}
		req.ServiceId = serviceID
		req.RequestID = requestID
		resp, err := h.service.ProcessAirtime(ctx, userID, req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		utils.JSON(w, http.StatusOK, resp)

	default:
		http.Error(w, "unsupported service type: "+serviceID, http.StatusBadRequest)
	}
}
