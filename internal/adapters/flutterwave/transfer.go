package adapters

import (
	"context"
	"encoding/json"
	"time"


	"payme/internal/application/transfer/dto"

)

// type Client struct {
// 	httpClient *httpx.Client
// }

// func NewClient() *Client {
// 	return &Client{
// 		httpClient: httpx.New(
// 			"https://api.flutterwave.com",
// 			map[string]string{
// 				"Authorization": "Bearer " + os.Getenv("FLW_SECRET_KEY"),
// 				"Content-Type": "application/json",
// 			},
// 		),
// 	}
// }
func (c *Client) CreateTransfers(ctx context.Context,req dto.CreateTransferRequest) (*CreateTransferResponse, error) {

	flwReq := CreateTransferRequest{
 	AccountNumber: req.AccountNumber,
    AccountBank:   req.AccountBank,
    Amount:        req.Amount,
    Currency:      "NGN",
    DebitCurrency: "NGN",
    Narration:     req.Narration,
    Reference:     "TXN-" + time.Now().Format("20060102150405"),
	}

	respBody, err := c.httpClient.DoRequest(
		ctx,
		"POST",
		"/v3/transfers",
		&flwReq,
		nil,
	)

	if err != nil {
		return nil, err
	}

	var resp CreateTransferResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}