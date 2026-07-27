package adapters

import (
	"context"
	"encoding/json"
	"time"


	"payme/internal/application/transfer/dto"

)


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

func(c *Client) ResolveAccountDetails(ctx context.Context, req dto.ResolveBankDetailsRequest) (*ResolveAccountDetailsResponse, error) {
	
	flwReq := ResolveAccountDetailsRequest{
		AccountNumber: req.AccountNumber,
		AccountBank:   req.AccountBank,
	}

	respBody, err := c.httpClient.DoRequest(
		ctx,
		"POST",
		"/v3/accounts/resolve",
		&flwReq,
		nil,
	)

	if err != nil {
		return nil, err
	}

	var resp ResolveAccountDetailsResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}