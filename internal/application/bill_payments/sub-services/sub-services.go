package sub_services

import (
	"encoding/json"
	"payme/internal/application/bill_payments/dto"

	"strings"
)

func IsTvService(serviceID string) bool {
switch strings.ToLower(serviceID) {
	case "dstv", "gotv", "startimes", "showmax":
		return true
	default:
		return false
	}
}

func IsElectricityService(serviceID string) bool {
	switch strings.ToLower(serviceID) {
	case "ikeja-electric", "eko-electric", "abuja-electric", "kano-electric","portharcourt-electric","kaduna-electric","enugu-electric","ibadan-electric","benin-electric","aba-electric","yola-electric","yedc":
	return true
	default:
		return false
	}
}

func IsMobileVtu(serviceID string) bool {
	switch strings.ToLower(serviceID) {
	case "mtn", "airtel", "glo", "9mobile":
		return true
	default:
		return false
	}
}

func IsMobileData(serviceID string) bool {
	switch strings.ToLower(serviceID) {
	case "mtn-data", "airtel-data", "glo-data", "9mobile-data":
		return true
	default:
		return false
	}
}

func BillServiceCategory(body []byte, categoryId string) (dto.BillCategoryResponseTwo,error) {

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

	return dto.BillCategoryResponseTwo{
		Data: categories,
	},nil

}
