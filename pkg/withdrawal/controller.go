package withdrawal

import (
	"net/http"
)

type WithdrawalController struct {
	// dependencies
}

func NewWithdrawalController() *WithdrawalController {
	return &WithdrawalController{}
}

func (c *WithdrawalController) RequestWithdrawal(w http.ResponseWriter, r *http.Request) {
	// logic for requesting a withdrawal
}

func (c *WithdrawalController) GetWithdrawalStatus(w http.ResponseWriter, r *http.Request) {
	// logic for checking status
}
