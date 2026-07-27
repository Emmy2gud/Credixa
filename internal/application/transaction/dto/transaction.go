package dto

type TransactionFilterRequest struct {
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`
	Category  string `json:"category"`
	Type      string `json:"type"`
	Status    string `json:"status"`
	Search    string `json:"search"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}
