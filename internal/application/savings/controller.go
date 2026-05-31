package savings

import (
	"log"
	"net/http"
	"payme/internal/config"
	"payme/internal/api/middleware"
	"time"

	"github.com/robfig/cron/v3"
)
	type Response struct {

    TargetAmount      float64   `json:"target_amount"`
    CurrentAmount     float64   `json:"current_amount"`
    Purpose           string    `json:"purpose"`
    Status            string    `json:"status"`
    AutoSave          bool      `json:"auto_save"`
    AutoSaveFrequency string    `json:"auto_save_frequency"`
    AutoSaveAmount    uint64   `json:"auto_save_amount"`
    LastAutoSaveDate  time.Time `json:"last_auto_save_date"`
	}
func StartAutoSaveScheduler(w http.ResponseWriter, r *http.Request) {

	userID, _ :=middleware.GetUserID(r)
	// unauthorized user
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}	

	var savings PersonalSaving
	if err := config.DB.Where("user_id = ?", userID).First(&savings).Error; err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	var res Response
	// Create a new cron scheduler
	// cron.New() with WithSeconds() lets you use second-level precision if needed.
	// Without it, the smallest unit is minutes.
	scheduler := cron.New()
 
	// ---------------------------------------------------------------
	// JOB 1: Run every day at midnight (00:00)
	// Cron format: "minute hour day month weekday"
	// "0 0 * * *" = at 00:00 every day
	// ---------------------------------------------------------------
	scheduler.AddFunc("0 0 * * *", func() {
		log.Println("[Cron] Running daily auto-save check...")
		runAutoSaveBatch(res, "daily",userID)
	})
 
	// ---------------------------------------------------------------
	// JOB 2: Run every Monday at midnight for weekly savers
	// "0 0 * * 1" = at 00:00 every Monday
	// ---------------------------------------------------------------
	scheduler.AddFunc("0 0 * * 1", func() {
		log.Println("[Cron] Running weekly auto-save check...")
		runAutoSaveBatch(res, "weekly",userID)
	})
 
	// ---------------------------------------------------------------
	// JOB 3: Run on the 1st of every month at midnight for monthly savers
	// "0 0 1 * *" = at 00:00 on day 1 of every month
	// ---------------------------------------------------------------
	scheduler.AddFunc("0 0 1 * *", func() {
		log.Println("[Cron] Running monthly auto-save check...")
		runAutoSaveBatch(res, "monthly",userID)
	})
 
	// Start the scheduler — it now runs in the background forever
	// as long as your server is running
	scheduler.Start()
 
	log.Println("[Cron] Auto-save scheduler started ✓")
 
	// Optional: keep a reference if you ever want to stop it gracefully
	// e.g. on server shutdown: scheduler.Stop()
}
