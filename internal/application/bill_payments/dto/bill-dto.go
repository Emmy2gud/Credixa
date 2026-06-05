package dto

//creating req and res dto for bill payment
type BillerCategory struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
}
type BillerCategoriesResponse struct {
	Categories []BillerCategory
}

type BillerCategoryResponse struct {
	Categories []BillerCategoryTwo `json:"content"`
}
type BillerCategoryTwo struct {
	ServiceID      string `json:"serviceID"`
	Name           string `json:"name"`
	MinimiumAmount string `json:"minimium_amount"`
	MaximumAmount  int `json:"maximum_amount"`
	ConvinienceFee string `json:"convinience_fee"`
	ProductType    string `json:"product_type"`
	Image          string `json:"image"`
}

// type BillCategoryResponse struct {
// 	CategoryType []BillCategoryDataResponse `json:"CategoryType"`
// }
type BillCategoryResponseTwo struct {
	Data map[string][]interface{} `json:"data"`
}

// type BillCategoryDataResponse struct {
// 	FixedPrice      string `json:"fixedPrice"`
// 	Name            string `json:"name"`
// 	VariationAmount string `json:"variation_amount"`
// 	VariationCode   string `json:"variation_code"`
// }

type StarttimesRequest struct {
	RequestID     string `json:"request_id"`
	ServiceId     string `json:"serviceID"`
	VariationCode string `json:"variation_code"`
	BillersCode   string `json:"billersCode"`
	Amount        string `json:"amount"`
	Phone         string `json:"phone"`
}
type ShowmaxRequest struct {
	RequestID     string `json:"request_id"`
	ServiceId     string `json:"serviceID"`
	VariationCode string `json:"variation_code"`
	BillersCode   string `json:"billersCode"`
	Amount        string `json:"amount"`
}
type VerifyElectricityRequest struct {
	ServiceId   string `json:"serviceID"`
	BillersCode string `json:"billersCode"`
	Type        string `json:"type"`
}
type VerifyTvSubscriptionRequest struct {
	ServiceId   string `json:"serviceID"`
	BillersCode string `json:"billersCode"`
}
type VerifyTvSubscriptionResponse struct {
	CustomerName   string `json:"Customer_Name"`
	Status         string `json:"Status"`
	DueDate        string `json:"Due_Date"`
	CustomerNumber string `json:"Customer_Number"`
	CustomerType   string `json:"Customer_Type"`
}

type VerifyElectricitySubscriptionResponse struct {
	CustomerName        string `json:"Customer_Name"`
	Address             string `json:"Address"`
	MeterNumber         string `json:"Meter_Number"`
	CustomerArrears     string `json:"Customer_Arrears"`
	MinimumAmount       string `json:"Minimum_Amount"`
	MinPurchaseAmount   string `json:"Min_Purchase_Amount"`
	CanVend             string `json:"Can_Vend"`
	BusinessUnit        string `json:"Business_Unit"`
	CustomerAccountType string `json:"Customer_Account_Type"`
	MeterType           string `json:"Meter_Type"`
	WrongBillersCode    bool   `json:"WrongBillersCode"`
}

type ElectricityRequest struct {
	RequestID     string `json:"request_id"`
	ServiceId     string `json:"serviceID"`
	VariationCode string `json:"variation_code"`
	BillersCode   string `json:"billersCode"`
	MeterNo       string `json:"meterNo"`
	Amount        string `json:"amount"`
	Phone         string `json:"phone"`
}

type ChangeTvRequest struct {
	RequestID        string `json:"request_id"`
	ServiceId        string `json:"serviceID"`
	VariationCode    string `json:"variation_code"`
	BillersCode      string `json:"billersCode"`
	Amount           string `json:"amount"`
	Phone            string `json:"phone"`
	SubscriptionType string `json:"subscription_type"`
	Quantity         string `json:"quantity"`
}

type CreateBillPaymentAirtimeRequest struct {
	RequestID string `json:"request_id"`
	ServiceId string `json:"serviceID"`
	Amount    string `json:"amount"`
	Phone     string `json:"phone"`
}
type CreateBillPaymentDataRequest struct {
	RequestID     string `json:"request_id"`
	ServiceId     string `json:"serviceID"`
	VariationCode string `json:"variation_code"`
	BillersCode   string `json:"billersCode"`
	Amount        string `json:"amount"`
	Phone         string `json:"phone"`
}

// ── SHARED RESPONSE DTO ───────────────────────────────────────────────────────
// All four service types return the same content.transactions shape.

type VTPassTransactions struct {
	Status         string  `json:"status"`
	ProductName    string  `json:"product_name"`
	UniqueElement  string  `json:"unique_element"`
	UnitPrice      string  `json:"unit_price"`
	Quantity       int     `json:"quantity"`
	Channel        string  `json:"channel"`
	Commission     float64 `json:"commission"`
	TotalAmount    float64 `json:"total_amount"`
	Type           string  `json:"type"`
	Email          string  `json:"email"`
	Phone          string  `json:"phone"`
	Amount         string  `json:"amount"`
	TransactionID  string  `json:"transactionId"`
	ConvenienceFee float64 `json:"convinience_fee"`
}

type VTPassContent struct {
	Transactions VTPassTransactions `json:"transactions"`
}

type BillPaymentResponse struct {
	Code                string        `json:"code"`
	Content             VTPassContent `json:"content"`
	ResponseDescription string        `json:"response_description"`
	RequestID           string        `json:"requestId"`
	Amount              float64       `json:"amount"`
	TransactionDate     string        `json:"transaction_date"`
	PurchasedCode       string        `json:"purchased_code"`
}
