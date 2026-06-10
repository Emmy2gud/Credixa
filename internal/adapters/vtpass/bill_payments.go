package adapters

import (
	"context"
	"encoding/json"
	"fmt"

	"os"

	"payme/pkg/httpx"
)


type Client struct {
	httpClient *httpx.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: httpx.New(
			"https://sandbox.vtpass.com",
			map[string]string{
				"api-key":      os.Getenv("VTPASS_API_KEY"),
				"secret-key":   os.Getenv("VTPASS_SECRET_KEY"),
				"Content-Type": "application/json",
			},
		),
	}
}
	// 7. Call VTPass
// var vtpastPayload = dto.CreateBillPaymentRequest{
// 		ServiceId:      p.serviceID,
// 		VariationCode:  p.variationCode,
// 		BillersCode:    p.billersCode,
// 		Amount:         p.amount,
// 		Phone:          p.phone,
// 		MeterNo:        p.meterNo,
// 		SubscriptionType: p.SubscriptionType,
// 		Quantity:         p.Quantity,
// 	}

func (c *Client) CreateBillPayment(ctx context.Context,payload interface{}) ([]byte, error) {

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal vtpass payload: %w", err)
	}
	_ = body // httpx.DoRequest accepts the struct directly in your current setup
 

	respBody, err := c.httpClient.DoRequest(ctx,"POST","/api/pay",&payload,nil)

	if err != nil {
		return nil, err
	}

	return respBody, nil
}