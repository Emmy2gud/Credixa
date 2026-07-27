package services

import (
	"context"
	"encoding/json"
	"fmt"

	"os"
	adapters "payme/internal/adapters/vtpass"
	"payme/internal/application/bill_payments/dto"
	"payme/internal/application/bill_payments/models"
	sub_services "payme/internal/application/bill_payments/sub-services"
	"payme/internal/application/transaction"
	"payme/internal/application/wallet"

	"payme/pkg/httpx"
	"payme/pkg/utils"

	"gorm.io/gorm"
)

type billPaymentParams struct {
	requestID      string
	serviceID      string
	variationCode  string
	billType       string
	amount         int64
	idempotencyKey string
	vtpassPayload  interface{}
}

type BillPaymentService interface {
	GetBillerCategories(ctx context.Context) (dto.BillerCategoriesResponse, error)
	GetBillerCategory(ctx context.Context, categoryID string) (dto.BillerCategoryResponse, error)
	GetBillCategory(ctx context.Context, categoryID string) (dto.BillCategoryResponseTwo, error)
	ProcessAirtime(ctx context.Context, userID uint, req dto.CreateBillPaymentAirtimeRequest) (*dto.BillPaymentResponse, error)
	ProcessData(ctx context.Context, userID uint, req dto.CreateBillPaymentDataRequest) (*dto.BillPaymentResponse, error)
	ProcessTV(ctx context.Context, userID uint, req dto.ChangeTvRequest) (*dto.BillPaymentResponse, error)
	ProcessElectricity(ctx context.Context, userID uint, req dto.ElectricityRequest) (*dto.BillPaymentResponse, error)
	VerifySubscription(ctx context.Context, serviceID string, electricityInput dto.VerifyElectricityRequest, tvinput dto.VerifyTvSubscriptionRequest) (interface{}, error)
}

type billPaymentService struct {
	db *gorm.DB
}

func NewBillPaymentService(db *gorm.DB) BillPaymentService {
	return &billPaymentService{db: db}
}

func (s *billPaymentService) idempotencyCheck(ctx context.Context, key string) (*transaction.Transaction, error) {
	if key == "" {
		return nil, nil
	}
	var t transaction.Transaction
	err := s.db.WithContext(ctx).Where("reference = ? AND status IN ('success','failed')", key).First(&t).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("idempotency check failed: %w", err)
	}
	return &t, nil
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
	categoryresp, _ := sub_services.BillServiceCategory(respbody, categoryID)
	// fmt.Println("categoryresp:",categoryresp)
	return categoryresp, nil

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

func (s *billPaymentService) processPayment(ctx context.Context, userID uint, p billPaymentParams) (*dto.BillPaymentResponse, error) {
	var billPay models.BillPayment
	// 1. Idempotency guard
	if existing, err := s.idempotencyCheck(ctx, p.idempotencyKey); err != nil {
		return nil, err
	} else if existing != nil {

		return &dto.BillPaymentResponse{
			Code: "000",
		}, nil
	}

	// 2. Get wallet
	w, err := wallet.GetWallet(userID)
	if err != nil {
		return nil, fmt.Errorf("wallet not found: %w", err)
	}

	// 3. Generate request ID (this also becomes the transaction reference)
	requestID, err := utils.GenerateRequestID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate request ID: %w", err)
	}
	s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 4. Deduct wallet (optimistic — refund on failure)
		if err := wallet.DeductWalletBalance(tx,w.ID, 0, p.amount, requestID, "bill payment", "bill_payment", "pending", "bill_payment"); err != nil {
			return fmt.Errorf("insufficient balance or deduction failed: %w", err)
		}

		// 5. Persist pending bill payment
		billPay := models.BillPayment{
			UserID:    userID,
			WalletID:  w.ID,
			BillType:  p.billType,
			Provider:  p.variationCode,
			Amount:    uint64(p.amount),
			Reference: requestID,
			Status:    "pending",
		}
		if err := s.db.WithContext(ctx).Create(&billPay).Error; err != nil {
			wallet.UpdateWalletBalance(tx,w.ID, 0, p.amount, requestID, "bill payment", "bill_payment", "failed", "bill_payment")
			return fmt.Errorf("failed to create bill payment record: %w", err)
		}
		return nil
	})

	// 6. Call VTPass
	respBody, err := adapters.NewClient().CreateBillPayment(ctx, p.vtpassPayload)
	if err != nil {
		defer s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			wallet.UpdateWalletBalance(tx,w.ID, 0, p.amount, requestID, "bill payment", "bill_payment", "failed", "bill_payment")
			billPay.Status = "failed"
			s.db.WithContext(ctx).Save(&billPay)
			return nil
		})
		return nil, fmt.Errorf("vtpass request failed: %w", err)
	}

	// 7. Parse response
	var result dto.BillPaymentResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse vtpass response: %w", err)
	}

	// 8. Update records based on VTPass outcome
	if result.Content.Transactions.Status == "failed" || result.Code != "000" {
		defer s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			wallet.UpdateWalletBalance(tx,w.ID, 0, p.amount, requestID, "bill payment", "bill_payment", "failed", "bill_payment")
			billPay.Status = "failed"
			s.db.WithContext(ctx).Save(&billPay)
			return nil
		})
	} else {
		billPay.Status = "success"
		s.db.WithContext(ctx).Save(&billPay)
	}
	s.db.WithContext(ctx).Save(&billPay)

	return &result, nil
}
