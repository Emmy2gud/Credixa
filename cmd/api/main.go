package main

import (
	"log"
	"net/http"
	"payme/pkg/accounts"
	"payme/pkg/bill_payments"
	"payme/pkg/config"
	"payme/pkg/notifications"
	"payme/pkg/pendingcard"
	"payme/pkg/routes"
	"payme/pkg/savings"
	"payme/pkg/token"
	"payme/pkg/transaction"
	"payme/pkg/transactionpin"
	"payme/pkg/transfer"
	"payme/pkg/user"
	"payme/pkg/wallet"
	"payme/pkg/withdrawal"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("../../.env")
	config.Connect()
	config.GetDB().AutoMigrate(&user.User{},&wallet.Wallet{},&savings.PersonalSaving{},&savings.GroupSaving{},&savings.PersonalSaving{},&savings.GroupSaving{},&savings.GroupMember{},&savings.GroupContribution{},&transactionpin.TransactionPin{},&wallet.SavingsWallet{},&accounts.VirtualAccount{},&notifications.Notification{},&transfer.Transfer{},&withdrawal.Withdrawal{},bill_payments.BillPayment{},&token.CardToken{},&pendingcard.PendingCard{},&transaction.Transaction{})
	r := mux.NewRouter()
	routes.SetupRoutes(r)
	http.Handle("/", r)
	log.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
