package bill_payments

import (
	"context"
	"encoding/json"
	"fmt"

	"os"
	"payme/internal/application/bill_payments/dto"
	"payme/internal/application/bill_payments/sub-services"
	"payme/internal/application/transaction"
	"payme/internal/application/wallet"
	"payme/internal/config"
	"payme/pkg/httpx"
	"payme/pkg/utils"
	"strconv"

	"gorm.io/gorm"
)

type BillPaymentService interface {
	GetBillerCategories(ctx context.Context) (dto.BillerCategoriesResponse, error)
	GetBillerCategory(ctx context.Context, categoryID string) (dto.BillerCategoryResponse, error)
	GetBillCategory(ctx context.Context, categoryID string) (dto.BillCategoryResponseTwo, error)
	
	CreateBillPayment(ctx context.Context, userID uint, serviceID, variationCode string, airtimeInput dto.CreateBillPaymentAirtimeRequest, dataInput dto.CreateBillPaymentDataRequest, tvInput dto.ChangeTvRequest, electricInput dto.ElectricityRequest) (*dto.BillPaymentResponse, error)
	VerifySubscription(ctx context.Context, serviceID string, electricityInput dto.VerifyElectricityRequest, tvinput dto.VerifyTvSubscriptionRequest) (interface{}, error)
}

type billPaymentService struct {
	db *gorm.DB
}

func NewBillPaymentService(db *gorm.DB) BillPaymentService {
	return &billPaymentService{db: db}
}

func (s *billPaymentService) GetBillerCategories(ctx context.Context) (dto.BillerCategoriesResponse, error) {
	flwClient := httpx.New(
		"https://vtpass.com",
		map[string]string{
			"api-key":      os.Getenv("VTPASS_API_KEY"),
			"secret-key":   os.Getenv("VTPASS_SECRET_KEY"),
			"Content-Type": "application/json",
		},
	)

	respbody, err := flwClient.DoRequest(ctx, "GET", "/api/service-categories", nil, nil)
	if err != nil {
		return dto.BillerCategoriesResponse{}, fmt.Errorf("vtpass request failed: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respbody, &result); err != nil {
		return dto.BillerCategoriesResponse{}, fmt.Errorf("failed to parse API response: %v", err)
	}

	// content is an ARRAY, not a map
	contentRaw, ok := result["content"].([]interface{})
	if !ok {
		return dto.BillerCategoriesResponse{}, fmt.Errorf("unexpected content structure in API response")
	}

	var categories []dto.BillerCategory
	for _, item := range contentRaw {
		entry, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		categories = append(categories, dto.BillerCategory{
			Identifier: entry["identifier"].(string),
			Name:       entry["name"].(string),
		})
	}

	return dto.BillerCategoriesResponse{Categories: categories}, nil
}

func (s *billPaymentService) GetBillerCategory(ctx context.Context, categoryID string) (dto.BillerCategoryResponse, error) {

	flwClient := httpx.New(
		"https://vtpass.com",
		map[string]string{
			"api-key":      os.Getenv("VTPASS_API_KEY"),
			"secret-key":   os.Getenv("VTPASS_SECRET_KEY"),
			"Content-Type": "application/json",
		},
	)
	respbody, err := flwClient.DoRequest(ctx, "GET", "/api/services?identifier="+categoryID, nil, nil)
	if err != nil {
		return dto.BillerCategoryResponse{}, err
	}
	var vtresponse dto.BillerCategoryResponse
	if err := json.Unmarshal(respbody, &vtresponse); err != nil {
    return dto.BillerCategoryResponse{}, err
}
return dto.BillerCategoryResponse{
    Categories: vtresponse.Categories,
}, nil

}

func (s *billPaymentService) GetBillCategory(ctx context.Context, categoryID string) (dto.BillCategoryResponseTwo, error) {
	flwClient := httpx.New(
		"https://vtpass.com",
		map[string]string{
			"api-key":      os.Getenv("VTPASS_API_KEY"),
			"secret-key":   os.Getenv("VTPASS_SECRET_KEY"),
			"Content-Type": "application/json",
		},
	)
	respbody, err := flwClient.DoRequest(ctx, "GET", "/api/service-variations?serviceID="+categoryID, nil, nil)
	if err != nil {
		return dto.BillCategoryResponseTwo{}, err
	}
	categoryresp,_:= sub_services.BillServiceCategory(respbody,categoryID)
	// fmt.Println("categoryresp:",categoryresp)
	return  categoryresp,nil

}


func (s *billPaymentService) VerifySubscription(ctx context.Context, serviceID string, electricityInput dto.VerifyElectricityRequest, tvinput dto.VerifyTvSubscriptionRequest) (interface{}, error) {
	switch {
	case sub_services.IsElectricityService(serviceID):
		electricityInput.ServiceId = serviceID
		flwClient := httpx.New(
			"https://sandbox.vtpass.com",
			map[string]string{
				"api-key":      os.Getenv("VTPASS_API_KEY"),
				"secret-key":   os.Getenv("VTPASS_SECRET_KEY"),
				"Content-Type": "application/json",
			},
		)

		respbody, err := flwClient.DoRequest(ctx, "POST", "/api/merchant-verify", electricityInput, nil)

		if err != nil {
			return respbody, fmt.Errorf("vtpass request failed: %v", err)
		}

		fmt.Println("VtPass response:", string(respbody))

		var result map[string]interface{}
		if err := json.Unmarshal(respbody, &result); err != nil {
			return dto.VerifyElectricitySubscriptionResponse{}, fmt.Errorf("failed to parse API response: %v", err)
		}

		data, ok := result["content"].(map[string]interface{})
		if !ok {
			return dto.VerifyElectricitySubscriptionResponse{}, fmt.Errorf("invalid data structure in API response")
		}

		return dto.VerifyElectricitySubscriptionResponse{
			MeterNumber:         data["Meter_Number"].(string),
			CustomerName:        data["Customer_Name"].(string),
			Address:             data["Address"].(string),
			CustomerArrears:     data["Customer_Arrears"].(string),
			MinimumAmount:       data["Minimum_Amount"].(string),
			MinPurchaseAmount:   data["Min_Purchase_Amount"].(string),
			CanVend:             data["Can_Vend"].(string),
			BusinessUnit:        data["Business_Unit"].(string),
			CustomerAccountType: data["Customer_Account_Type"].(string),
			MeterType:           data["Meter_Type"].(string),
			WrongBillersCode:    data["WrongBillersCode"].(bool),
		}, nil

	case sub_services.IsTvService(serviceID):
		tvinput.ServiceId = serviceID
		flwClient := httpx.New(
			"https://sandbox.vtpass.com",
			map[string]string{
				"api-key":      os.Getenv("VTPASS_API_KEY"),
				"secret-key":   os.Getenv("VTPASS_SECRET_KEY"),
				"Content-Type": "application/json",
			},
		)

		respbody, err := flwClient.DoRequest(ctx, "POST", "/api/merchant-verify", tvinput, nil)
		if err != nil {
			return respbody, fmt.Errorf("flutterwave request failed: %v", err)
		}
		var result map[string]interface{}
		if err := json.Unmarshal(respbody, &result); err != nil {
			return dto.VerifyTvSubscriptionResponse{}, fmt.Errorf("failed to parse API response: %v", err)
		}

		data, ok := result["content"].(map[string]interface{})
		if !ok {
			return dto.VerifyTvSubscriptionResponse{
				CustomerName:   data["Customer_Name"].(string),
				Status:         data["Status"].(string),
				DueDate:        data["Due_Date"].(string),
				CustomerNumber: data["Customer_Number"].(string),
				CustomerType:   data["Customer_Type"].(string),
			}, fmt.Errorf("invalid data structure in API response")
		}

		return dto.VerifyTvSubscriptionResponse{
			CustomerName:   data["Customer_Name"].(string),
			Status:         data["Status"].(string),
			DueDate:        data["Due_Date"].(string),
			CustomerNumber: data["Customer_Number"].(string),
			CustomerType:   data["Customer_Type"].(string),
		}, nil

	default:
		return nil, fmt.Errorf("unsupported service type: %s", serviceID)
	}
}
func (s *billPaymentService) CreateBillPayment(ctx context.Context, userID uint, serviceID, variationCode string, airtimeInput dto.CreateBillPaymentAirtimeRequest, dataInput dto.CreateBillPaymentDataRequest, tvInput dto.ChangeTvRequest, electricInput dto.ElectricityRequest) (*dto.BillPaymentResponse, error) {

	// 1. Get wallet and idempotency
	var t transaction.Transaction	


   if err := s.db.WithContext(ctx).Where("reference = ?", req.IdempotencyKey).First(&t).Error; err != nil {
		return nil, fmt.Errorf("transaction not found: %w", err)
	}
	wallets, err := wallet.GetWallet(userID)
	if err != nil {
		return nil, fmt.Errorf("wallet not found")
	}

	// 2. Extract amount and phone per service type — these live in the DTOs
	var amountStr, phone string
	switch {
	case sub_services.IsElectricityService(serviceID):
		amountStr = electricInput.Amount
		phone = electricInput.Phone
		electricInput.ServiceId = serviceID
		electricInput.VariationCode = variationCode
	case sub_services.IsTvService(serviceID):
		amountStr = tvInput.Amount
		phone = tvInput.Phone
		tvInput.ServiceId = serviceID
		tvInput.VariationCode = variationCode
	case sub_services.IsMobileData(serviceID):
		amountStr = dataInput.Amount
		phone = dataInput.Phone
		dataInput.ServiceId = serviceID
		dataInput.VariationCode = variationCode
	case sub_services.IsMobileVtu(serviceID):
		amountStr = airtimeInput.Amount
		phone = airtimeInput.Phone
		airtimeInput.ServiceId = serviceID
	default:
		return nil, fmt.Errorf("unsupported service type: %s", serviceID)
	}

	// 3. Parse amount
	amountFloat, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amountFloat <= 0 {
		return nil, fmt.Errorf("invalid amount: %s", amountStr)
	}
	_ = phone // used in payload below

	// 4. Generate request ID
	requestID, err := utils.GenerateRequestID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate request ID: %v", err)
	}

	// 5. Deduct wallet
	if err := wallet.DeductWalletBalance(userID, int64(amountFloat)); err != nil {
		return nil, fmt.Errorf("insufficient balance or deduction failed: %v", err)
	}

	// 6. Create pending transaction
	trans := transaction.Transaction{
		Amount:    int64(amountFloat),
		Reference: requestID,
		Type:      "bill_payment",
		Status:    "pending",
		UserID:    userID,
		WalletID:  wallets.ID,
	}
	config.DB.Create(&trans)

	// 7. Create pending bill payment record
	billpayment := BillPayment{
		UserID:    userID,
		WalletID:  wallets.ID,
		BillType:  serviceID,
		Provider:  variationCode,
		Amount:    uint64(amountFloat),
		Reference: requestID,
		Status:    "pending",
	}
	config.DB.Create(&billpayment)

	// 8. Build the VTPass payload — inject the generated request_id
	var vtpassPayload interface{}
	switch {
	case sub_services.IsElectricityService(serviceID):
		electricInput.RequestID = requestID
		vtpassPayload = electricInput
	case sub_services.IsTvService(serviceID):
		tvInput.RequestID = requestID
		vtpassPayload = tvInput
	case sub_services.IsMobileData(serviceID):
		dataInput.RequestID = requestID
		vtpassPayload = dataInput
	case sub_services.IsMobileVtu(serviceID):
		airtimeInput.RequestID = requestID
		vtpassPayload = airtimeInput
	}

	// 9. Call VTPass
	vtClient := httpx.New(
		"https://sandbox.vtpass.com",
		map[string]string{
			"api-key":      os.Getenv("VTPASS_API_KEY"),
			"secret-key":   os.Getenv("VTPASS_SECRET_KEY"),
			"Content-Type": "application/json",
		},
	)

	respBody, err := vtClient.DoRequest(ctx, "POST", "/api/pay", vtpassPayload, nil)
	if err != nil {
		// Refund on network failure
		wallet.UpdateWalletBalance(userID, int64(amountFloat))
		trans.Status = "failed"
		config.DB.Save(&trans)
		billpayment.Status = "failed"
		config.DB.Save(&billpayment)
		return nil, fmt.Errorf("vtpass request failed: %v", err)
	}

	// 10. Parse and type the response
	var result dto.BillPaymentResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse vtpass response: %v", err)
	}

	fmt.Printf("VTPass [%s] response code=%s status=%s\n",
		serviceID, result.Code, result.Content.Transactions.Status)

	// 11. Update records based on outcome
	if result.Content.Transactions.Status == "failed" || result.Code != "000" {
		wallet.UpdateWalletBalance(userID, int64(amountFloat))
		trans.Status = "failed"
		billpayment.Status = "failed"
	} else {
		trans.Status = "success"
		billpayment.Status = "success"
	}
	s.db.WithContext(ctx).Save(&trans)
	s.db.WithContext(ctx).Save(&billpayment)

	return &result, nil
}

