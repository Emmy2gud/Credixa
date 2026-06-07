package accounts

import (
	"context"
	"encoding/json"
	"fmt"

	adapters "payme/internal/adapters/flutterwave"
	"payme/internal/application/accounts/dto"
	"payme/internal/application/wallet"


	"github.com/google/uuid"

	"gorm.io/gorm"
)

type VirtualAccountService interface {
	CreateVirtualAccount(ctx context.Context, input dto.CreateVirtualAccountInput, userID uint) (dto.VirtualAccountResponse, error)
}

type virtualAccountService struct {
	db *gorm.DB
	
}

func NewVirtualAccountService(db *gorm.DB) VirtualAccountService {
	return &virtualAccountService{
		db: db,
	}
}



func (s *virtualAccountService) CreateVirtualAccount(ctx context.Context, input dto.CreateVirtualAccountInput, userID uint) (dto.VirtualAccountResponse, error) {
	input.Currency = "NGN"
	input.BankCode = "090772"
	input.Narration = "Create a virtual account for this user"
	input.Is_permanent = true
	input.TxRef = "token_ch_" + uuid.New().String()

   	flwResp, err := adapters.NewClient().CreateVirtualAccount(ctx,input)
	if err != nil {
		return dto.VirtualAccountResponse{}, fmt.Errorf("flutterwave request failed: %v", err)
	}

	fmt.Println("FlutterWave response:", string(flwResp))

	var result map[string]interface{}
	if err := json.Unmarshal(flwResp, &result); err != nil {
		return dto.VirtualAccountResponse{}, fmt.Errorf("failed to parse API response: %v", err)
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return dto.VirtualAccountResponse{}, fmt.Errorf("invalid data structure in API response")
	}

	var w wallet.Wallet
	if err := s.db.WithContext(ctx).Where("user_id = ?", userID).First(&w).Error; err != nil {
		return dto.VirtualAccountResponse{}, fmt.Errorf("wallet not found: %v", err)
	}

	virtualAcc := VirtualAccount{
		WalletID:      uint(w.ID),
		AccountNumber: data["account_number"].(string),
		AccountName:   input.Firstname + " " + input.Lastname,
		BankName:      data["bank_name"].(string),
		Provider:      "flutterwave",
		Status:        data["account_status"].(string),
	}

	// Fallback if AccountName is missing from API
	if virtualAcc.AccountName == "" {
		virtualAcc.AccountName = fmt.Sprintf("%s %s", input.Firstname, input.Lastname)
	}

	if err := s.db.WithContext(ctx).Create(&virtualAcc).Error; err != nil {
		return dto.VirtualAccountResponse{}, fmt.Errorf("failed to create virtual account in database: %v", err)
	}

	return dto.VirtualAccountResponse{
    ResponseCode: data["response_code"].(string),
    ResponseMessage: data["response_message"].(string),
    FlwRef:          data["flw_ref"].(string),
    OrderRef:        data["order_ref"].(string),
    AccountNumber:   data["account_number"].(string),
    Frequency:       data["frequency"].(string),
    BankName:      data["bank_name"].(string),
    CreatedAt:       data["created_at"].(string),
    ExpiryDate:      data["expiry_date"].(string),
    Note:            data["note"].(string),
    Amount:          data["amount"].(string),
	}, nil
}
