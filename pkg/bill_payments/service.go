package bill_payments

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"payme/pkg/config"
	"payme/pkg/transaction"
	"payme/pkg/utils"
	"payme/pkg/wallet"
	"strconv"
	"time"
)

func ProcessBillPayment(userID uint, serviceID, variationCode, phone, amountStr string) ([]byte, error) {
	var wallets wallet.Wallet
	if err := config.DB.Where("user_id = ?", userID).First(&wallets).Error; err != nil {
		return nil, fmt.Errorf("wallet not found")
	}

	// creating requestid using date formatting for strings
	const AlphaNumericBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	randomString, err := utils.RandomString(AlphaNumericBytes, 10)
	if err != nil {
		return nil, err
	}

	loc, _ := time.LoadLocation("Africa/Lagos")
	now := time.Now().In(loc)
	requestid := now.Format("200601021504") + randomString
	fmt.Println("Request ID:", requestid)

	fmt.Printf("Crediting wallet for user %d with amount %s\n", userID, amountStr)
	// convert amountStr to float64
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid amount")
	}

	if err := wallet.DeductWalletBalance(userID, amount); err != nil {
		return nil, fmt.Errorf("failed to update wallet balance: %v", err)
	}

	// Create pending transaction record
	var trans transaction.Transaction
	trans.Amount = amount
	trans.Reference = requestid
	trans.Type = "bill_payment"
	trans.Status = "pending"
	trans.UserID = userID
	trans.WalletID = wallets.ID
	config.DB.Create(&trans)

	// Create pending bill payment record
	var billpayment BillPayment
	billpayment.UserID = userID
	billpayment.WalletID = wallets.ID
	billpayment.BillType = serviceID
	billpayment.Provider = variationCode
	billpayment.Amount = amount
	billpayment.Reference = requestid
	billpayment.Status = "pending"
	config.DB.Create(&billpayment)

	formData := fmt.Sprintf("request_id=%s&serviceID=%s&billersCode=%s&variation_code=%s&amount=%s&phone=%s",
		requestid, serviceID, phone, variationCode, amountStr, phone)

	req, err := http.NewRequest("POST", "https://sandbox.vtpass.com/api/pay", bytes.NewBufferString(formData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("api-key", os.Getenv("VTPASS_API_KEY"))
	req.Header.Set("secret-key", os.Getenv("VTPASS_SECRET_KEY"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respbody, _ := io.ReadAll(resp.Body)
	
	// Handle response and update status/refund
	var result map[string]interface{}
	if err := json.Unmarshal(respbody, &result); err != nil {
		return respbody, fmt.Errorf("invalid response from provider")
	}

	content, ok := result["content"].(map[string]interface{})
	if ok && content["status"] == "failed" {
		wallet.UpdateWalletBalance(userID, amount)
		trans.Status = "failed"
		config.DB.Save(&trans)
		billpayment.Status = "failed"
		config.DB.Save(&billpayment)
	} else {
		trans.Status = "success"
		config.DB.Save(&trans)
		billpayment.Status = "success"
		config.DB.Save(&billpayment)
	}

	fmt.Println("Vtpass response:", string(respbody))
	return respbody, nil
}
