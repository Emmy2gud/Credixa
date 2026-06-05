package transaction

import (
	"net/http"
)

type TransactionController struct {
	service TransactionService
}

func NewTransactionController(service TransactionService) *TransactionController {
	return &TransactionController{
		service: service,
	}
}

func (h *TransactionController) GetTransactionHistory(w http.ResponseWriter, r *http.Request) {
	// GetTransactionHistory logic here
}

func (h *TransactionController) GetTransactionByID(w http.ResponseWriter, r *http.Request) {
	// GetTransactionByID logic here
}

func (h *TransactionController) GetWalletLogs(w http.ResponseWriter, r *http.Request) {
	// GetWalletLogs logic here
}
