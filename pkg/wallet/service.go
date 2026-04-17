package wallet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"payme/pkg/config"
	"payme/pkg/middleware"
	"payme/pkg/pendingcard"
	"payme/pkg/token"
	"payme/pkg/transaction"
	"payme/pkg/utils"
	"strconv"

	"github.com/google/uuid"
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

func InitiateCardWalletFunding(userCard ChargeRequest, r *http.Request) (string, map[string]interface{}, error) {
	userID, _ := middleware.GetUserID(r)
	fmt.Println("user id is",userID)
	// Generate a unique transaction reference server-side.
	// This is returned to the client so they can reference it in step 2 (PIN auth).
	userCard.TxRef = "token_ch_" + uuid.New().String()
	// Marshal the raw card data FIRST — this is what we store in the DB.
	// We store it unencrypted so we can unmarshal it back in step 2.
	body, err := json.Marshal(userCard)
	if err != nil {
		return "", nil, fmt.Errorf("failed to marshal card request: %v", err)
	}

	// Save the raw payload to DB keyed by tx_ref.
	// Step 2 will look this up using tx_ref.
	pendingEntry := pendingcard.PendingCard{
		UserID:  userID,
		Payload: body,
		TxRef:   userCard.TxRef,
	}
	config.DB.Create(&pendingEntry)
//get the wallet id for a particular user in the transaction record
var w Wallet
if err := config.DB.Where("user_id = ?", userID).First(&w).Error; err != nil {
    return "", nil, fmt.Errorf("wallet not found: %v", err)
}

	// Create a FundingSession to track the transaction synchronously and asynchronously
	amount := userCard.Amount
	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse amount: %v", err)
	}
	fundingSession := transaction.Transaction{
		UserID:    userID,
		Amount:    amountFloat,
		Status:    "pending",
		Reference: userCard.TxRef,
		WalletID:  w.ID,
		Type:      "credit",
		Category:  "funding",
	}
	config.DB.Create(&fundingSession)

	// Encrypt ONLY for sending to Flutterwave — not for storage.
	encryptedBody, err := utils.Encryption3des(string(body))
	if err != nil {
		return "", nil, fmt.Errorf("encryption failed: %v", err)
	}

	flwPayload, _ := json.Marshal(map[string]string{"client": encryptedBody})

	req, err := http.NewRequest("POST", "https://api.flutterwave.com/v3/charges?type=card", bytes.NewBuffer(flwPayload))
	if err != nil {
		return "", nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+os.Getenv("FLW_SECRET_KEY"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &result)

	// Return tx_ref alongside Flutterwave response.
	// Controller will include tx_ref in the API response so the client can use it in step 2.
	return userCard.TxRef, result, nil
}

// AuthorizeCardCharge re-sends the full card payload WITH the PIN added inside,
// then re-encrypts and posts to the same /v3/charges?type=card endpoint.
// This is exactly how Flutterwave PIN mode works.
func AuthorizeCardCharge(pin string, chargeRequest ChargeRequest) (map[string]interface{}, error) {
	type AuthorizeInfo struct {
		ChargeRequest
		Authorization map[string]string `json:"authorization"`
	}
	// Add the PIN authorization to the card payload before encrypting
	// Initialize embedded fields by putting the struct itself
	info := AuthorizeInfo{
		ChargeRequest: chargeRequest,
		Authorization: map[string]string{
			"mode": "pin",
			"pin":  pin,
		},
	}

	body, err := json.Marshal(info)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal card request: %v", err)
	}

	// Encrypt the full payload (same as step 1, but now includes authorization)
	encryptedData, err := utils.Encryption3des(string(body))
	if err != nil {
		return nil, fmt.Errorf("encryption failed: %v", err)
	}

	flwPayload, _ := json.Marshal(map[string]string{"client": encryptedData})
	fmt.Println(flwPayload)

	req, err := http.NewRequest("POST", "https://api.flutterwave.com/v3/charges?type=card", bytes.NewBuffer(flwPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+os.Getenv("FLW_SECRET_KEY"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &result)

	fmt.Println(result)
	return result, nil
}

// AuthorizeCardFundingService looks up the saved ChargeRequest by tx_ref,
// calls AuthorizeCardCharge with the PIN, then cleans up the pending record.
func AuthorizeCardFundingService(txRef string, pin string) (map[string]interface{}, error) {
	// 1. Fetch the saved pending card from DB
	var pending pendingcard.PendingCard
	if err := config.DB.Where("tx_ref = ?", txRef).First(&pending).Error; err != nil {
		return nil, fmt.Errorf("no pending charge found for tx_ref %s: %v", txRef, err)
	}

	// 2. Deserialize the stored payload back into a ChargeRequest
	var chargeRequest ChargeRequest
	if err := json.Unmarshal(pending.Payload, &chargeRequest); err != nil {
		return nil, fmt.Errorf("failed to parse saved card payload: %v", err)
	}

	// 3. Run the authorization with the PIN
	result, err := AuthorizeCardCharge(pin, chargeRequest)
	if err != nil {
		return nil, err
	}

	// 4. Clean up — delete the pending record now that it's been used
	config.DB.Delete(&pending)

	return result, nil
}

func ValidateCardCharge(ref, otp string) (map[string]interface{}, error) {
	validatePayload := map[string]string{
		"flw_ref": ref,
		"otp":     otp,
	}
	body, err := json.Marshal(validatePayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal validation request: %v", err)
	}

	req, err := http.NewRequest("POST", "https://api.flutterwave.com/v3/validate-charge", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+os.Getenv("FLW_SECRET_KEY"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &result)

	fmt.Println(result)
	return result, nil
}

func UpdateWalletBalance(userID uint, amount float64) error {
	var wallet Wallet

	if err := config.DB.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		return fmt.Errorf("wallet not found for user %d: %v", userID, err)
	}
	newBalance := wallet.Balance + amount
	if err := config.DB.Model(&wallet).Update("balance", newBalance).Error; err != nil {
		return fmt.Errorf("failed to update wallet balance: %v", err)
	}

	return nil
}

func DeductWalletBalance(userID uint, amount float64) error {
	var wallet Wallet

	if err := config.DB.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		return fmt.Errorf("wallet not found for user %d: %v", userID, err)
	}
	newBalance := wallet.Balance - amount
	if err := config.DB.Model(&wallet).Update("balance", newBalance).Error; err != nil {
		return fmt.Errorf("failed to update wallet balance: %v", err)
	}

	return nil
}

func VerifyCard(id string, userID uint) (map[string]interface{}, error) {

	verifyReq, err := http.NewRequest("GET", "https://api.flutterwave.com/v3/transactions/"+id+"/verify", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create verify request: %v", err)
	}

	verifyReq.Header.Set("Authorization", "Bearer "+os.Getenv("FLW_SECRET_KEY"))
	verifyReq.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	verifyResp, err := client.Do(verifyReq)
	if err != nil {
		return nil, err
	}
	defer verifyResp.Body.Close()

	var verifyResult FlwVerifyResponse
	verifyBody, _ := io.ReadAll(verifyResp.Body)
	json.Unmarshal(verifyBody, &verifyResult)
	fmt.Println("Verify result:", verifyResult)

	if verifyResult.Status == "success" && verifyResult.Data.Status == "successful" {
		// Store the card in the database
		cardToken := token.CardToken{
			UserID:    userID,
			Token:     verifyResult.Data.Card.Token,
			CardBrand: verifyResult.Data.Card.Type,
			Last4:     verifyResult.Data.Card.Last4,
			Expiry:    verifyResult.Data.Card.Expiry,
			First6:    verifyResult.Data.Card.First6,
			Issuer:    verifyResult.Data.Card.Issuer,
			Country:   verifyResult.Data.Card.Country,
			Type:      verifyResult.Data.Card.Type,
		}

		// Use FirstOrCreate or similar to avoid duplicates if needed
		// For now, let's just create a new record
		if err := config.DB.Create(&cardToken).Error; err != nil {
			fmt.Printf("Failed to save card token: %v\n", err)
		}
	}

	// Return the result
	combined := map[string]interface{}{
		"verification": verifyResult,
	}
	return combined, nil
}
