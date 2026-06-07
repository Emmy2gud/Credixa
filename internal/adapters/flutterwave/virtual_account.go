package adapters

import (
	"context"

	"os"
	"payme/internal/application/accounts/dto"
	"payme/pkg/httpx"
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
func (c *Client) CreateVirtualAccount(
	ctx context.Context,
	input dto.CreateVirtualAccountInput,
) ( []byte,  error) {

	respBody, err := c.httpClient.DoRequest(
		ctx,
		"POST",
		"/v3/virtual-account-numbers",
		&input,
		nil,
	)


	return  respBody, err
}