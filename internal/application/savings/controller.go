package savings

import (
	"net/http"
	"payme/internal/api/middleware"
	"payme/internal/config"
	"payme/pkg/utils"
	"time"
)

type Response struct {
	TargetAmount      uint64    `json:"target_amount"`
	CurrentAmount     uint64    `json:"current_amount"`
	Purpose           string    `json:"purpose"`
	Status            string    `json:"status"`
	AutoSave          bool      `json:"auto_save"`
	AutoSaveFrequency string    `json:"auto_save_frequency"`
	AutoSaveAmount    uint64    `json:"auto_save_amount"`
	LastAutoSaveDate  time.Time `json:"last_auto_save_date"`
}

func CreateSavingGoal(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req Response
	utils.ParseBody(r, &req)

	var s PersonalSaving
	s.UserID = userID
	s.Purpose = req.Purpose
	s.TargetAmount = req.TargetAmount
	s.AutoSave = true
	s.AutoSaveFrequency = req.AutoSaveFrequency
	s.AutoSaveAmount = req.AutoSaveAmount
	config.DB.Create(&s)

	w.WriteHeader(http.StatusOK)
}

