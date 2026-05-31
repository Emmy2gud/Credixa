package bill_payments

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"payme/internal/application/transaction"
	"payme/internal/application/wallet"
	"payme/internal/config"
	"payme/pkg/utils"
	"strconv"
	"strings"
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

	if err := wallet.DeductWalletBalance(userID, uint64(amount)); err != nil {
		return nil, fmt.Errorf("failed to update wallet balance: %v", err)
	}

	// Create pending transaction record
	var trans transaction.Transaction
	trans.Amount = uint64(amount)
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
		wallet.UpdateWalletBalance(userID, uint64(amount))
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

func BillServiceCategory(body []byte, categoryId string) map[string]interface{} {
	var responseBody map[string]interface{}
	// Convert JSON response into a Go map
	json.Unmarshal(body, &responseBody)

	content := responseBody["content"].(map[string]interface{})
	variations := content["variations"].([]interface{})
	delete(content, "ServiceName")
	delete(content, "serviceID")
	delete(content, "convinience_fee")
	delete(content, "varations")

	categories := make(map[string][]interface{})
	// Loop through each data plan
	for _, item := range variations {

		// Convert current item to a map
		plan := item.(map[string]interface{})

		planCode := plan["variation_code"].(string)
		category := ""
		// Decide category based on network type
		switch {
		case categoryId == "mtn-data":
			category = GetMTNDataCategory(planCode)
		case categoryId == "glo-data":
			category = GetGloDataCategory(planCode)
		case categoryId == "airtel-data":
			category = GetAirtelDataCategory(planCode)
		case categoryId == "9mobile-sme-data":
			category = Get9mobileDataCategory(planCode)
		default:
			category = "Monthly"
		}

		categories[category] = append(categories[category], plan)
	}

	content["variations"] = categories
return responseBody
	
}



func GetGloDataCategory(variationCode string) string {
	switch {
	case strings.HasPrefix(variationCode, "glo-social-") ||
		strings.HasPrefix(variationCode, "glo-telegram-") ||
		strings.HasPrefix(variationCode, "glo-insta-") ||
		strings.HasPrefix(variationCode, "glo-tiktok-") ||
		strings.HasPrefix(variationCode, "glo-opera-") ||
		strings.HasPrefix(variationCode, "glo-youtube-"):
		return "Social Media"
	case strings.HasPrefix(variationCode, "glo-daily-"):
		return "Daily"
	case strings.HasPrefix(variationCode, "glo-2days-"),
		strings.HasPrefix(variationCode, "glo-2weeks-"):
		return "Short-Term"
	case strings.HasPrefix(variationCode, "glo-monthly-"):
		return "Monthly (Day + Night)"
	case strings.HasPrefix(variationCode, "glo-weekend-")||strings.HasPrefix(variationCode, "glo-sunday"):
		return "Weekend"
	case strings.HasPrefix(variationCode, "glo-mega-"):
		return "Mega"
	case strings.HasPrefix(variationCode, "glo-tv-"):
		return "Glo TV"
	case strings.HasPrefix(variationCode, "glo-wtf-"):
		return "WTF"
	case strings.HasPrefix(variationCode, "glo-dg-"):
		return "SME"
	case strings.HasPrefix(variationCode, "glo-special-"):
		return "Special"
	default:
		return "Standard"
	}
}

func GetMTNDataCategory(variationCode string) string {
    switch {
    case variationCode == "mtn-10mb-100" || variationCode == "mtn-50mb-200" ||
        variationCode == "mtn-2-5gb-600" || variationCode == "mtn-3gb-800"||variationCode=="mtn-230mb-200":
        return "Daily / 2-Day"
    case variationCode == "mtn-20hrs-1500" || variationCode == "mtn-7gb-2000"||variationCode=="mtn-1500mb-1000":
        return "Weekly"
    case variationCode == "mtn-100gb-20000" || variationCode == "mtn-160gb-30000" ||
        variationCode == "mtn-400gb-50000" || variationCode == "mtn-600gb-75000" ||
        variationCode == "mtn-120gb-22000"||variationCode=="monthly":
        return "Multi-Month"
    case variationCode == "mtn-4-5tb-450000" || variationCode == "mtn-1tb-110000":
        return "Yearly"
    default:
        return "Monthly"
    }
}
func GetAirtelDataCategory(variationCode string) string {
	switch {

	case strings.HasPrefix(variationCode, "social"):
		return "social"
		case strings.HasPrefix(variationCode, "mifi"):
		return "mifi"
	case strings.HasSuffix(variationCode, "-1"),
		strings.HasSuffix(variationCode, "-2"),
		strings.HasSuffix(variationCode, "-7"):
		return "Short-Term"

	case strings.Contains(variationCode, "6000-30"),
		strings.Contains(variationCode, "10000"),
		strings.Contains(variationCode, "15000"),
		strings.Contains(variationCode, "20000"):
		return "Mega"

	case strings.Contains(variationCode, "1500-2"):
		return "Binge"

	case strings.Contains(variationCode, "50"),
		strings.Contains(variationCode, "100"),
		strings.Contains(variationCode, "200"),
		strings.Contains(variationCode, "300"):

		return "Daily"

	default:
		return "Monthly"
	}
}
func Get9mobileDataCategory(variationCode string) string {
	switch {
	case strings.Contains(variationCode, "150-1") ||
		strings.Contains(variationCode, "100") ||
		strings.Contains(variationCode, "200") ||
		strings.Contains(variationCode, "300"):
		return "Daily"

	case strings.Contains(variationCode, "1500-7"):
		return "Weekly"

	case strings.Contains(variationCode, "2500"):
		return "Night + Weekend"

	case strings.Contains(variationCode, "27500") ||
		strings.Contains(variationCode, "55000") ||
		strings.Contains(variationCode, "110000"):
		return "Long-Term"

	default:
		return "Monthly"
	}
}