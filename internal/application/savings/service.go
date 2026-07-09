package savings

import (
	"log"
	"time"

	"payme/internal/application/wallet"
	"payme/internal/config"

	"github.com/robfig/cron/v3"
)

// InitAutoSaveScheduler starts the global background cron job for auto-savings.
// This should be called once during application boot (e.g. in main.go).
func InitAutoSaveScheduler() {
	scheduler := cron.New()

	// Daily auto-saves at midnight: "0 0 * * *"
	_, err := scheduler.AddFunc("0 0 * * *", func() {
		log.Println("[Cron] Running daily auto-save check...")
		processAutoSaves("daily")
	})
	if err != nil {
		log.Println("[Cron] failed to add daily job:", err)
	}

	// Weekly auto-saves on Mondays at midnight: "0 0 * * 1"
	_, err = scheduler.AddFunc("0 0 * * 1", func() {
		log.Println("[Cron] Running weekly auto-save check...")
		processAutoSaves("weekly")
	})
	if err != nil {
		log.Println("[Cron] failed to add weekly job:", err)
	}

	// Monthly auto-saves on the 1st of the month at midnight: "0 0 1 * *"
	_, err = scheduler.AddFunc("0 0 1 * *", func() {
		log.Println("[Cron] Running monthly auto-save check...")
		processAutoSaves("monthly")
	})
	if err != nil {
		log.Println("[Cron] failed to add monthly job:", err)
	}

	scheduler.Start()
	log.Println("[Cron] Auto-save scheduler started ✓")
}

// processAutoSaves queries all active auto-saves for a given frequency and processes them.
func processAutoSaves(frequency string) {
    var users []PersonalSaving
    if err := config.DB.Where("auto_save = ? AND auto_save_frequency = ?", true, frequency).Find(&users).Error; err != nil {
        log.Println("[Cron] failed to fetch auto-save users:", err)
        return
    }

    for _, saving := range users {
        runAutoSave(saving)
    }
}

func runAutoSave(saving PersonalSaving) {
    if err := wallet.DeductWalletBalance(saving.UserID, int64(saving.AutoSaveAmount)); err != nil {
        log.Println("[Cron] deduct failed for user", saving.UserID, err)
        return
    }
    if err := wallet.UpdateSavingsWalletBalance(saving.UserID, int64(saving.AutoSaveAmount)); err != nil {
        log.Println("[Cron] savings wallet update failed for user", saving.UserID, err)
        return
    }

    saving.CurrentAmount +=saving.AutoSaveAmount
    saving.LastAutoSaveDate = time.Now()
    config.DB.Save(&saving)
  }