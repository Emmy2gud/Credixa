package wallet

import (
	"context"
	"encoding/json"
	"fmt"

	adapters "payme/internal/adapters/flutterwave"
	"payme/internal/application/pendingcard"
	"payme/internal/application/token"
	"payme/internal/application/transaction"
	"payme/internal/application/wallet/dto"
	"payme/internal/config"
	"payme/pkg/utils"
	"strconv"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FlwVerifyResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		ID            int     `json:"id"`
		TxRef         string  `json:"tx_ref"`
		FlwRef        string  `json:"flw_ref"`
		Amount        float64 `json:"amount"`
		Currency      string  `json:"currency"`
		ProcessorResp string  `json:"processor_response"`
		Status        string  `json:"status"`
		Card          struct {
			First6  string `json:"first_6digits"`
			Last4   string `json:"last_4digits"`
			Issuer  string `json:"issuer"`
			Country string `json:"country"`
			Type    string `json:"type"`
			Expiry  string `json:"expiry"`
			Token   string `json:"token"`
		} `json:"card"`
	} `json:"data"`
}

type WalletService interface {
	InitiateCardWalletFunding(ctx context.Context, req dto.ChargeRequest, userID uint) (string, dto.InitializeCardResponse, error)
	AuthorizeCardFundingService(ctx context.Context, txRef string, pin string) (dto.AuthorizeCardResponse, error)
	ValidateCardCharge(ctx context.Context, ref, otp string) (dto.ValidateCardResponse, error)
	VerifyCard(ctx context.Context, id string, userID uint) (dto.VerifyChargeResponse, error)
}

type walletService struct {
	db *gorm.DB
}

func NewWalletService(db *gorm.DB) WalletService {
	return &walletService{
		db: db,
	}
}

func (s *walletService) InitiateCardWalletFunding(ctx context.Context, req dto.ChargeRequest, userID uint) (string, dto.InitializeCardResponse, error) {

	// Generate a unique transaction reference server-side.
	// This is returned to the client so they can reference it in step 2 (PIN auth).
	req.TxRef = "token_ch_" + uuid.New().String()
	// Marshal the raw card data FIRST — this is what we store in the DB.
	// We store it unencrypted so we can unmarshal it back in step 2.
	body, err := json.Marshal(req)
	if err != nil {
		return "", dto.InitializeCardResponse{}, fmt.Errorf("failed to marshal card request: %v", err)
	}

	// Save the raw payload to DB keyed by tx_ref.
	// Step 2 will look this up using tx_ref.
	pendingEntry := pendingcard.PendingCard{
		UserID:  userID,
		Payload: body,
		TxRef:   req.TxRef,
	}
	if err := s.db.WithContext(ctx).Create(&pendingEntry).Error; err != nil {
		return "", dto.InitializeCardResponse{}, fmt.Errorf("failed to create pending charge entry: %v", err)
	}

	//get the wallet id for a particular user in the transaction record
	var w Wallet
	if err := s.db.WithContext(ctx).Where("user_id = ?", userID).First(&w).Error; err != nil {
		return "", dto.InitializeCardResponse{}, fmt.Errorf("wallet not found: %v", err)
	}

	// Create a FundingSession to track the transaction synchronously and asynchronously
	amount := req.Amount
	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return "", dto.InitializeCardResponse{}, fmt.Errorf("failed to parse amount: %v", err)
	}
	fundingSession := transaction.Transaction{
		UserID:    userID,
		Amount:    int64(amountFloat),
		Status:    "pending",
		Reference: req.TxRef,
		WalletID:  w.ID,
		Type:      "credit",
		Category:  "funding",
	}
	if err := s.db.WithContext(ctx).Create(&fundingSession).Error; err != nil {
		return "", dto.InitializeCardResponse{}, fmt.Errorf("failed to create funding session: %v", err)
	}

	// Encrypt ONLY for sending to Flutterwave — not for storage.
	encryptedBody, err := utils.Encryption3des(string(body))
	if err != nil {
		return "", dto.InitializeCardResponse{}, fmt.Errorf("encryption failed: %v", err)
	}

	flwPayload, _ := json.Marshal(map[string]string{"client": encryptedBody})
	flwResp, err := adapters.NewClient().InitiateCardWalletFunding(ctx, flwPayload)
	if err != nil {
		return "", dto.InitializeCardResponse{}, fmt.Errorf("failed to initiate card wallet funding: %v", err)
	}

	// Return tx_ref alongside Flutterwave response.
	// Controller will include tx_ref in the API response so the client can use it in step 2.
	return req.TxRef, dto.InitializeCardResponse{

		Mode:   flwResp.Meta.Authorization.Mode,
		Fields: flwResp.Meta.Authorization.Fields,
	}, nil
}

// AuthorizeCardCharge re-sends the full card payload WITH the PIN added inside,
// then re-encrypts and posts to the same /v3/charges?type=card endpoint.
// This is exactly how Flutterwave PIN mode works.
func (s *walletService) authorizeCardCharge(ctx context.Context, pin string, chargeRequest dto.ChargeRequest) (dto.AuthorizeCardResponse, error) {
	type AuthorizeInfo struct {
		dto.ChargeRequest
		Authorization map[string]string `json:"authorization"`
	}
	// Add the PIN authorization to the card payload before encrypting
	// Initialize embedded fields by putting the struct itself
	info := AuthorizeInfo{
		ChargeRequest: chargeRequest,
		Authorization: map[string]string{
			"mode":    "pin",
			"pin":     pin,
			"city":    chargeRequest.City,
			"address": chargeRequest.Address,
			"state":   chargeRequest.State,
			"zip":     chargeRequest.Zip,
		},
	}

	body, err := json.Marshal(info)
	if err != nil {
		return dto.AuthorizeCardResponse{}, fmt.Errorf("failed to marshal card request: %v", err)
	}

	// Encrypt the full payload (same as step 1, but now includes authorization)
	encryptedData, err := utils.Encryption3des(string(body))
	if err != nil {
		return dto.AuthorizeCardResponse{}, fmt.Errorf("encryption failed: %v", err)
	}

	flwPayload, _ := json.Marshal(map[string]string{"client": encryptedData})

	flwResp, err := adapters.NewClient().AuthorizeCardFunding(ctx, flwPayload)
	if err != nil {
		return dto.AuthorizeCardResponse{}, fmt.Errorf("failed to initiate card wallet funding: %v", err)
	}

	fmt.Println(flwResp)
	return dto.AuthorizeCardResponse{

		Mode:     flwResp.Meta.Authorization.Mode,
		Endpoint: flwResp.Meta.Authorization.Endpoint,
	}, nil
}

// AuthorizeCardFundingService looks up the saved ChargeRequest by tx_ref,
// calls authorizeCardCharge with the PIN, then cleans up the pending record.
func (s *walletService) AuthorizeCardFundingService(ctx context.Context, txRef string, pin string) (dto.AuthorizeCardResponse, error) {
	// 1. Fetch the saved pending card from DB
	var pending pendingcard.PendingCard
	if err := s.db.WithContext(ctx).Where("tx_ref = ?", txRef).First(&pending).Error; err != nil {
		return dto.AuthorizeCardResponse{}, fmt.Errorf("no pending charge found for tx_ref %s: %v", txRef, err)
	}

	// 2. Deserialize the stored payload back into a ChargeRequest
	var chargeRequest dto.ChargeRequest
	if err := json.Unmarshal(pending.Payload, &chargeRequest); err != nil {
		return dto.AuthorizeCardResponse{}, fmt.Errorf("failed to parse saved card payload: %v", err)
	}

	// 3. Run the authorization with the PIN
	result, err := s.authorizeCardCharge(ctx, pin, chargeRequest)
	if err != nil {
		return dto.AuthorizeCardResponse{}, err
	}

	// 4. Clean up — delete the pending record now that it's been used
	s.db.WithContext(ctx).Delete(&pending)

	return result, nil
}

func (s *walletService) ValidateCardCharge(ctx context.Context, ref, otp string) (dto.ValidateCardResponse, error) {

	validatePayload := dto.ValidateCardRequest{
		Otp:    otp,
		FlwRef: ref,
	}

	flwResp, err := adapters.NewClient().ValidateCardWalletFunding(ctx, validatePayload)
	if err != nil {
		return dto.ValidateCardResponse{}, fmt.Errorf("failed to initiate card wallet funding: %v", err)
	}

	return dto.ValidateCardResponse{
		Status:       flwResp.Status,
		Amount:       flwResp.Data.Amount,
		First6Digits: flwResp.Data.Card.First6Digits,
		Last4Digits:  flwResp.Data.Card.Last4Digits,
		Issuer:       flwResp.Data.Card.Issuer,
		Country:      flwResp.Data.Card.Country,
		Type:         flwResp.Data.Card.Type,
		Expiry:       flwResp.Data.Card.Expiry,
	}, nil
}

func (s *walletService) VerifyCard(ctx context.Context, id string, userID uint) (dto.VerifyChargeResponse, error) {

	flwResp, err := adapters.NewClient().VerifyCardWalletFunding(ctx, id)
	if err != nil {
		return dto.VerifyChargeResponse{}, fmt.Errorf("failed to initiate card wallet funding: %v", err)
	}

	if flwResp.Status == "success" && flwResp.Data.Status == "successful" {
		// Store the card in the database
		cardToken := token.CardToken{
			UserID:    userID,
			Token:     flwResp.Data.Token,
			CardBrand: flwResp.Data.Card.Type,
			Last4:     flwResp.Data.Card.Last4Digits,
			Expiry:    flwResp.Data.Card.Expiry,
			First6:    flwResp.Data.Card.First6Digits,
			Issuer:    flwResp.Data.Card.Issuer,
			Country:   flwResp.Data.Card.Country,
			Type:      flwResp.Data.Card.Type,
		}

		if err := s.db.WithContext(ctx).Create(&cardToken).Error; err != nil {
			fmt.Printf("Failed to save card token: %v\n", err)
		}
	}

	return dto.VerifyChargeResponse{
		Status:       flwResp.Status,
		Amount:       flwResp.Data.Amount,
		First6Digits: flwResp.Data.Card.First6Digits,
		Last4Digits:  flwResp.Data.Card.Last4Digits,
		Issuer:       flwResp.Data.Card.Issuer,
		Country:      flwResp.Data.Card.Country,
		Type:         flwResp.Data.Card.Type,
		Expiry:       flwResp.Data.Card.Expiry,
	}, nil
}

// Keep these package-level functions intact for external imports by other packages.
func GetWallet(userID uint) (Wallet, error) {
	var w Wallet
	if err := config.DB.Where("user_id = ?", userID).First(&w).Error; err != nil {
		return Wallet{}, fmt.Errorf("wallet not found for user %d: %v", userID, err)
	}
	return w, nil
}

func UpdateWalletBalance(db *gorm.DB,walletID uint,fee, amount int64,ref,description,entrytype,status,category string) error {

	var wallet Wallet

	if err := config.DB.Where("id = ?", walletID).First(&wallet).Error; err != nil {
		return fmt.Errorf("wallet not found for user %d: %v", walletID, err)
	}
	newBalance := wallet.Balance + int64(amount)
	if err := config.DB.Model(&wallet).Update("balance", newBalance).Error; err != nil {
		return fmt.Errorf("failed to update wallet balance: %v", err)
	}
	//wallet ledger entry
	walletLedger := WalletLedger{
		WalletID:      wallet.ID,
		UserID:        wallet.UserID,
		TransactionID: wallet.ID,
		Amount:        amount,
		BalanceBefore: wallet.Balance,
		BalanceAfter:  newBalance,
		Description:   description,
		Status:        status,
		EntryType:    entrytype,
	}
	if err := config.DB.Create(&walletLedger).Error; err != nil {
		return fmt.Errorf("failed to create wallet ledger entry: %v", err)
	}
    //creating transaction entry
	tx:=transaction.Transaction{
		UserID: wallet.UserID,
		WalletID: wallet.ID,
		Type: entrytype,
		Category: category,
		Amount: amount,
		Fee: fee,
		Reference: ref,
		Status: status,
		Description: description,
	}
	if err := config.DB.Create(&tx).Error; err != nil {
		return fmt.Errorf("failed to create transaction entry: %v", err)
	}
	
	return nil
}

func DeductWalletBalance(db *gorm.DB,walletID uint,fee, amount int64,ref,description,entrytype,status,category string) error {
	var wallet Wallet

	if err := config.DB.Where("id = ?", walletID).First(&wallet).Error; err != nil {
		return fmt.Errorf("wallet not found for user %d: %v", walletID, err)
	}
	newBalance := wallet.Balance - int64(amount)
	if err := config.DB.Model(&wallet).Update("balance", newBalance).Error; err != nil {
		return fmt.Errorf("failed to update wallet balance: %v", err)
	}
	//wallet ledger entry
	walletLedger := WalletLedger{
		WalletID:      wallet.ID,
		UserID:        wallet.UserID,
		TransactionID: wallet.ID,
		Amount:        amount,
		BalanceBefore: wallet.Balance,
		BalanceAfter:  newBalance,
		Description:   description,
		Status:        status,
		EntryType:    entrytype,
	}
	if err := config.DB.Create(&walletLedger).Error; err != nil {
		return fmt.Errorf("failed to create wallet ledger entry: %v", err)
	}
    //creating transaction entry
	tx:=transaction.Transaction{
		UserID: wallet.UserID,
		WalletID: wallet.ID,
		Type: entrytype,
		Category: category,
		Amount: amount,
		Fee: fee,
		Reference: ref,
		Status: status,
		Description: description,
	}
	if err := config.DB.Create(&tx).Error; err != nil {
		return fmt.Errorf("failed to create transaction entry: %v", err)
	}
	
	return nil
}

func UpdateSavingsWalletBalance(db *gorm.DB,walletID uint,fee, amount int64,ref,description,entrytype,status,category string) error {
	var wallet SavingsWallet

	if err := config.DB.Where("id = ?", walletID).First(&wallet).Error; err != nil {
		return fmt.Errorf("wallet not found for user %d: %v", walletID, err)
	}
	newBalance := wallet.Balance + int64(amount)
	if err := config.DB.Model(&wallet).Update("balance", newBalance).Error; err != nil {
		return fmt.Errorf("failed to update wallet balance: %v", err)
	}
	//wallet ledger entry
	walletLedger := WalletLedger{
		WalletID:      wallet.ID,
		UserID:        wallet.UserID,
		TransactionID: wallet.ID,
		Amount:        amount,
		BalanceBefore: wallet.Balance,
		BalanceAfter:  newBalance,
		Description:   description,
		Status:        status,
		EntryType:    entrytype,
	}
	if err := config.DB.Create(&walletLedger).Error; err != nil {
		return fmt.Errorf("failed to create wallet ledger entry: %v", err)
	}
    //creating transaction entry
	tx:=transaction.Transaction{
		UserID: wallet.UserID,
		WalletID: wallet.ID,
		Type: entrytype,
		Category: category,
		Amount: amount,
		Fee: fee,
		Reference: ref,
		Status: status,
		Description: description,
	}
	if err := config.DB.Create(&tx).Error; err != nil {
		return fmt.Errorf("failed to create transaction entry: %v", err)
	}

	return nil
}

func DeductSavingsWalletBalance(db *gorm.DB,walletID uint,fee,amount int64,ref,description,entrytype,status,category string) error {
	var wallet SavingsWallet

	if err := config.DB.Where("id = ?", walletID).First(&wallet).Error; err != nil {
		return fmt.Errorf("wallet not found for user %d: %v", walletID, err)
	}
	newBalance := wallet.Balance - int64(amount)
	if err := config.DB.Model(&wallet).Update("balance", newBalance).Error; err != nil {
		return fmt.Errorf("failed to update wallet balance: %v", err)
	}
	//wallet ledger entry
	walletLedger := WalletLedger{
		WalletID:      wallet.ID,
		UserID:        wallet.UserID,
		TransactionID: wallet.ID,
		Amount:        amount,
		BalanceBefore: wallet.Balance,
		BalanceAfter:  newBalance,
		Description:   description,
		Status:        status,
		EntryType:    entrytype,
	}
	if err := config.DB.Create(&walletLedger).Error; err != nil {
		return fmt.Errorf("failed to create wallet ledger entry: %v", err)
	}
    //creating transaction entry
	tx:=transaction.Transaction{
		UserID: wallet.UserID,
		WalletID: wallet.ID,
		Type: entrytype,
		Category: category,
		Amount: amount,
		Fee: fee,
		Reference: ref,
		Status: status,
		Description: description,
	}
	if err := config.DB.Create(&tx).Error; err != nil {
		return fmt.Errorf("failed to create transaction entry: %v", err)
	}

	return nil
}
