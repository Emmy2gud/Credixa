package services

import (
	"context"

	"payme/internal/application/bill_payments/dto"
	"payme/pkg/utils"
)

func (s *billPaymentService) ProcessTV(ctx context.Context, userID uint, req dto.ChangeTvRequest) (*dto.BillPaymentResponse, error) {

	amount, err := utils.ParseAmount(req.Amount)
	if err != nil {
		return nil, err
	}

	return s.processPayment(ctx, userID, billPaymentParams{
		requestID:      req.RequestID,
		serviceID:      req.ServiceId,
		variationCode:  req.VariationCode,
		billType:       req.ServiceId,
		amount:         amount,
		idempotencyKey: req.IdempotencyKey,
		vtpassPayload:  req,
	})
}
