package adapters

import (
	"context"
	"encoding/json"
	"os"

	"payme/internal/application/tiers/dto"
	"payme/pkg/httpx"

	"github.com/google/uuid"
)

type Client struct {
	httpClient *httpx.Client
}

func QoerIdNewClient() *Client {
	return &Client{
		httpClient: httpx.New(
			"https://api.qoreid.com/v1/ng",
			map[string]string{
				"Authorization": "Bearer " + os.Getenv("QOEID_SECRET_KEY"),
				"Content-Type":  "application/json",
			},
		),
	}
}

func (c *Client) KYC_Nin(ctx context.Context, req dto.Tier3NinRequest) (*KycTier3NinResponse, error) {


	qoreIdReq := KycTier3NinRequest{
		NIN:       req.NIN,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		DOB:       req.DOB,
	}
	respBody, err := c.httpClient.DoRequest(ctx, "POST", "/identities/virtual-nin/"+req.NIN, &qoreIdReq, nil)
	if err != nil {
		return nil, err
	}

	var resp KycTier3NinResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) KYC_Address(ctx context.Context, req dto.Tier3AddressRequest) (*KycTier3AddressResponse, error) {
     qoreidAddressRequest := KycTier3AddressRequest{
		CustomerReference:"Ref_" + uuid.New().String(),
		Street: req.Street,
		LgaName: req.LgaName,
		StateName: req.StateName,
		City: req.City,
		Landmark: req.Landmark,
		ApplicantFirstName: req.ApplicantFirstName,
		ApplicantLastName: req.ApplicantLastName,
		ApplicantPhone: req.ApplicantPhone,
		ApplicantDOB: req.ApplicantDOB,
	}

	respBody, err := c.httpClient.DoRequest(ctx, "POST", "/address/verifications", &qoreidAddressRequest, nil)
	if err != nil {
		return nil, err
	}

	var resp KycTier3AddressResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
