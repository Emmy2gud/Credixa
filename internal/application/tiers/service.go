package tiers

import (
	"context"
	"errors"
	"fmt"
	"time"

	qoerid "payme/internal/adapters/QoerId"
	adapters "payme/internal/adapters/flutterwave"
	"payme/internal/application/tiers/dto"
	"payme/internal/application/user"

	"gorm.io/gorm"
)

type TierService interface {
	Tier2KycUpload(ctx context.Context, input dto.Tier2Request, userID uint) (dto.InitiateTier2Response, error)
	RetrieveBvnKycUpload(ctx context.Context, ref string) (dto.BvnRetrievalResponse, error)
	Tier3Verification(ctx context.Context, input dto.Tier3Request, userID uint) (dto.Tier3Response, error)
}

type tierService struct {
	db *gorm.DB
}

func NewTierService(db *gorm.DB) TierService {
	return &tierService{
		db: db,
	}
}

func (s *tierService) Tier2KycUpload(ctx context.Context, input dto.Tier2Request, userID uint) (dto.InitiateTier2Response, error) {

	flwResp, err := adapters.NewClient().InitiateKYCTIER2(ctx, input)
	if err != nil {
		return dto.InitiateTier2Response{}, fmt.Errorf("failed to initiate tier 2 kyc: %w", err)
	}

	// 2. Only now save — and only the reference, NOT the raw BVN
	kycRecord := KYC{
		UserID:    userID,
		Tier:      2,
		Reference: flwResp.Data.Ref,
		Status:    "pending",
	}
	//3.update user tier record
	var user user.User
	s.db.WithContext(ctx).Where("id = ?", userID).First(&user)

	user.Tier = 2
	user.KYCStatus = "pending"

	if err := s.db.WithContext(ctx).Save(&user).Error; err != nil {
		return dto.InitiateTier2Response{}, fmt.Errorf("failed to update user tier: %w", err)
	}

	if err := s.db.WithContext(ctx).Create(&kycRecord).Error; err != nil {
		return dto.InitiateTier2Response{}, fmt.Errorf("failed to save kyc: %w", err)
	}

	return dto.InitiateTier2Response{
		Status:  flwResp.Status,
		Message: flwResp.Message,
		Url:     flwResp.Data.Url,
		Ref:     flwResp.Data.Ref,
	}, nil
}

func (s *tierService) RetrieveBvnKycUpload(ctx context.Context, ref string) (dto.BvnRetrievalResponse, error) {

	var kycRecord KYC
	if err := s.db.WithContext(ctx).
		Where("reference = ?", ref).
		First(&kycRecord).Error; err != nil {
		return dto.BvnRetrievalResponse{}, errors.New("kyc record not found for this reference")
	}

	if kycRecord.Status == "approved" {
		return dto.BvnRetrievalResponse{Status: "approved", Message: "already verified"}, nil
	}

	// 2. Call Flutterwave's retrieve endpoint
	flwResp, err := adapters.NewClient().RetrieveBVN(ctx, ref)
	if err != nil {
		return dto.BvnRetrievalResponse{}, fmt.Errorf("failed to retrieve bvn: %w", err)
	}

	// 3. Update based on result — save ONLY the verdict, not raw personal data
	if flwResp.Data.Status == "COMPLETED" {
		now := time.Now()
		kycRecord.Status = "approved"
		kycRecord.VerifiedAt = &now
		s.db.WithContext(ctx).Save(&kycRecord)

		// 4. Bump the user's tier
		s.db.WithContext(ctx).Model(&user.User{}).
			Where("id = ?", kycRecord.UserID).
			Updates(user.User{Tier: 2, KYCStatus: "verified"})
	} else {
		kycRecord.Status = "rejected"
		s.db.WithContext(ctx).Save(&kycRecord)
	}

	return dto.BvnRetrievalResponse{
		Status:  kycRecord.Status,
		Message: flwResp.Message,
	}, nil
}
func (s *tierService) Tier3Verification(ctx context.Context, input dto.Tier3Request, userID uint) (dto.Tier3Response, error) {
	// 1. Validate inputs
	if input.NIN == "" {
		return dto.Tier3Response{}, errors.New("NIN is required for Tier 3 verification")
	}
	if input.Street == "" || input.LgaName == "" || input.StateName == "" || input.City == "" {
		return dto.Tier3Response{}, errors.New("complete address details (street, lga, state, city) are required for Tier 3 verification")
	}
	if input.FirstName == "" || input.LastName == "" || input.DOB == "" {
		return dto.Tier3Response{}, errors.New("first name, last name, and date of birth are required for verification")
	}

	// 2. Initialize the KYC record to keep track of the status/results
	kycRecord := KYC{
		UserID:          userID,
		Tier:            3,
		Status:          "pending",
		NINVerified:     false,
		AddressVerified: false,
	}

	// 3. Step 1: NIN Verification via QoreID
	ninReq := dto.Tier3NinRequest{
		NIN:       input.NIN,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		DOB:       input.DOB,
	}
	qoreIdNinResp, err := qoerid.QoerIdNewClient().KYC_Nin(ctx, ninReq)
	if err != nil  {
		kycRecord.Status = "rejected"
		kycRecord.RejectionReason = "NIN verification API error: " + err.Error()
		if saveErr := s.db.WithContext(ctx).Create(&kycRecord).Error; saveErr != nil {
			return dto.Tier3Response{Status: "failed", Message: err.Error()}, fmt.Errorf("failed to save KYC record after NIN error: %w (original error: %v)", saveErr, err)
		}
		return dto.Tier3Response{Status: "rejected", Message: "NIN verification failed due to system error"}, fmt.Errorf("failed to verify NIN: %w", err)
	}

	if qoreIdNinResp.Status.Status != "verified" {
		kycRecord.Status = "rejected"
		kycRecord.RejectionReason = "NIN not verified: " + qoreIdNinResp.Status.State
		if saveErr := s.db.WithContext(ctx).Create(&kycRecord).Error; saveErr != nil {
			return dto.Tier3Response{Status: "rejected", Message: qoreIdNinResp.Status.State}, fmt.Errorf("failed to save KYC record after NIN rejection: %w", saveErr)
		}
		return dto.Tier3Response{Status: "rejected", Message: "NIN verification status not verified: " + qoreIdNinResp.Status.State}, nil
	}

	// NIN verified successfully
	kycRecord.NINVerified = true

	// 4. Step 2: Address Verification via QoreID
	addressReq := dto.Tier3AddressRequest{
		Street:             input.Street,
		LgaName:            input.LgaName,
		StateName:          input.StateName,
		City:               input.City,
		Landmark:           input.Landmark,
		ApplicantFirstName: input.FirstName,
		ApplicantLastName:  input.LastName,
		ApplicantPhone:     input.Phone,
		ApplicantDOB:       input.DOB,
	}
	qoreIdAddrResp, err := qoerid.QoerIdNewClient().KYC_Address(ctx, addressReq)
	if err != nil {
		kycRecord.Status = "rejected"
		kycRecord.RejectionReason = "Address verification API error: " + err.Error()
		if saveErr := s.db.WithContext(ctx).Create(&kycRecord).Error; saveErr != nil {
			return dto.Tier3Response{Status: "failed", Message: err.Error()}, fmt.Errorf("failed to save KYC record after Address error: %w (original error: %v)", saveErr, err)
		}
		return dto.Tier3Response{Status: "rejected", Message: "Address verification failed due to system error"}, fmt.Errorf("failed to verify address: %w", err)
	}

	if qoreIdAddrResp.Status.Status != "verified" {
		kycRecord.Status = "rejected"
		kycRecord.RejectionReason = "Address not verified: " + qoreIdAddrResp.Status.State
		if saveErr := s.db.WithContext(ctx).Create(&kycRecord).Error; saveErr != nil {
			return dto.Tier3Response{Status: "rejected", Message: qoreIdAddrResp.Status.State}, fmt.Errorf("failed to save KYC record after Address rejection: %w", saveErr)
		}
		return dto.Tier3Response{Status: "rejected", Message: "Address verification status not verified: " + qoreIdAddrResp.Status.State}, nil
	}

	// Address verified successfully
	kycRecord.AddressVerified = true

	// 5. Both verifications succeeded! Update DB inside a transaction
	now := time.Now()
	kycRecord.Status = "approved"
	kycRecord.VerifiedAt = &now

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Save the successful KYC record
		if err := tx.Create(&kycRecord).Error; err != nil {
			return fmt.Errorf("failed to save KYC record: %w", err)
		}

		// Update the user's tier to 3 and KYC status to verified
		// Also update user's address field with the verified address
		fullAddress := fmt.Sprintf("%s, %s, %s, %s", input.Street, input.City, input.LgaName, input.StateName)
		if err := tx.Model(&user.User{}).
			Where("id = ?", userID).
			Updates(user.User{
				Tier:      3,
				KYCStatus: "verified",
				Address:   fullAddress,
			}).Error; err != nil {
			return fmt.Errorf("failed to update user tier details: %w", err)
		}

		return nil
	})

	if err != nil {
		return dto.Tier3Response{Status: "failed", Message: "failed to complete verification process"}, err
	}

	return dto.Tier3Response{
		Status:  "approved",
		Message: "Tier 3 verification successful",
	}, nil
}
