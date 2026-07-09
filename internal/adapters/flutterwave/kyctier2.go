package adapters

import (
	"context"
	"encoding/json"
	"payme/internal/application/tiers/dto"
)



func (c *Client) InitiateKYCTIER2(ctx context.Context,req dto.Tier2Request) (*InitiateKycTier2Response, error) {

		qoeridReq :=InitiateKycTier2Request{
			FirstName: req.FirstName,
			LastName: req.LastName,
			RedirectUrl: "https://yourapp.com/kyc-success",
		}

	respBody, err := c.httpClient.DoRequest(ctx,"POST","/bvn/verifications/"+req.BVN,&qoeridReq,nil,)

	if err != nil {
		return nil, err
	}

	var resp InitiateKycTier2Response
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) RetrieveBVN(ctx context.Context, ref string) (*RetrieveBvnResponse, error) {
	respBody, err := c.httpClient.DoRequest(ctx,"GET","/bvn/verifications/"+ref,nil,nil)
	if err != nil {
		return nil, err
	}
	var resp RetrieveBvnResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}