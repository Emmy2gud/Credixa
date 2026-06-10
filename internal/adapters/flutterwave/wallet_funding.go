package adapters

import (
	"context"
	"encoding/json"


	"payme/internal/application/wallet/dto"
)


func (c *Client) InitiateCardWalletFunding(ctx context.Context,req []byte) (*InitializeCardResponse, error) {



	respBody, err := c.httpClient.DoRequest(
		ctx,
		"POST",
		"/v3/charges?type=card",
		&req,
		nil,
	)

	if err != nil {
		return nil, err
	}

	var resp InitializeCardResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
func (c *Client) AuthorizeCardFunding(ctx context.Context, req []byte) (*AuthorizationCardResponse, error) {

	respBody, err := c.httpClient.DoRequest(
		ctx,
		"POST",
		"/v3/charges?type=card",
		&req,
		nil,
	)

	if err != nil {
		return nil, err
	}

	var resp AuthorizationCardResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
func (c *Client) ValidateCardWalletFunding(ctx context.Context,req dto.ValidateCardRequest) (*ValidateCardResponse, error) {
   payload:=VerifyCardRequest{
	Otp:req.Otp,
	TxRef:req.FlwRef,
   }


	respBody, err := c.httpClient.DoRequest(
		ctx,
		"POST",
		"/v3/validate-charge",
		&payload,
		nil,
	)

	if err != nil {
		return nil, err
	}

	var resp ValidateCardResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
func (c *Client) VerifyCardWalletFunding(ctx context.Context,id string) (*VerifyChargeResponse, error) {
   

	respBody, err := c.httpClient.DoRequest(
		ctx,
		"GET",
		"/v3/transactions/"+id+"/verify",
		nil,
		nil,
	)

	if err != nil {
		return nil, err
	}

	var resp VerifyChargeResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}