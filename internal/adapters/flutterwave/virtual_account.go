package adapters

import (
	"context"
	"encoding/json"

	"os"

	"payme/internal/application/accounts/dto"
	"payme/pkg/httpx"

	"github.com/google/uuid"
)


type Client struct {
	httpClient *httpx.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: httpx.New(
			"https://api.flutterwave.com",
			map[string]string{
				"Authorization": "Bearer " + os.Getenv("FLW_SECRET_KEY"),
				"Content-Type": "application/json",
			},
		),
	}
}
func (c *Client) CreateVirtualAccount(ctx context.Context,req dto.CreateVirtualAccountRequest) (*CreateVirtualAccountResponse, error) {

	flwReq := CreateVirtualAccountRequest{
		Email:       req.Email,
		Phone:       req.Phone,
		Firstname:   req.Firstname,
		Lastname:    req.Lastname,
		Bvn:         req.Bvn,
		TxRef:       "token_ch_" + uuid.New().String(),
		Narration:   "Virtual account for user " + req.Firstname + " " + req.Lastname,
		BankCode:    "035",
		Currency:    "NGN",
		IsPermanent: true,
		Amount: 	 1020,
	}

	respBody, err := c.httpClient.DoRequest(
		ctx,
		"POST",
		"/v3/virtual-account-numbers",
		&flwReq,
		nil,
	)

	if err != nil {
		return nil, err
	}

	var resp CreateVirtualAccountResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}