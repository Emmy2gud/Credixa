package adapters

type CreateBillPaymentRequest struct {
	RequestID        string `json:"request_id"`
	ServiceId        string `json:"serviceID"`
	VariationCode    string `json:"variation_code"`
	BillersCode      string `json:"billersCode"`
	Amount           int64  `json:"amount"`
	Phone            string `json:"phone"`
	MeterNo        string `json:"meterNo"`
	SubscriptionType string `json:"subscription_type"`
	Quantity         string `json:"quantity"`
	
}

// ── Response ──────────────────────────────────────────────────────────────────
 
type BillPaymentResponse struct {
	Code    string `json:"code"`
	Content struct {
		Transactions struct {
			Status string `json:"status"`
		} `json:"transactions"`
	} `json:"content"`
}