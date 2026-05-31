package savings

import (

    "payme/internal/config"

	"payme/internal/application/wallet"

)


func runAutoSaveBatch(res Response, frequency string,userID uint) {
		var savingswallet wallet.SavingsWallet
  //check if frequenc is daily,weekly or monthly
  
  switch frequency {
  case "daily":
    runDailyAutoSave(res,userID,savingswallet)
  case "weekly":
    runWeeklyAutoSave(res,userID,savingswallet)
  case "monthly":
    runMonthlyAutoSave(res,userID,savingswallet)
  }

}


func runDailyAutoSave(res Response,userID uint,savingswallet wallet.SavingsWallet) {
   savings := PersonalSaving{
	UserID: userID,
	TargetAmount: res.TargetAmount,
    CurrentAmount: res.CurrentAmount,
    Purpose: res.Purpose,
    Status: res.Status,
    AutoSave: res.AutoSave,
    AutoSaveFrequency: res.AutoSaveFrequency,
    AutoSaveAmount: res.AutoSaveAmount,
    LastAutoSaveDate: res.LastAutoSaveDate,

}
	config.DB.Create(&savings)

	//deducting the participant amount from his wallet
	if err := wallet.UpdateSavingsWalletBalance(uint(userID), res.AutoSaveAmount); err != nil {
		return 
	}
	if err := wallet.DeductWalletBalance(uint(userID), res.AutoSaveAmount); err != nil {
		return 
	}
	config.DB.Save(&savingswallet)


}

	

func runWeeklyAutoSave(res Response, userID uint,savingswallet wallet.SavingsWallet) {

}
func runMonthlyAutoSave(res Response, userID uint,savingswallet wallet.SavingsWallet) {

}