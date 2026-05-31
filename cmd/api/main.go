package main

import (
	"log"
	"net/http"
	"payme/internal/api/routes"
	"payme/internal/application/accounts"
	"payme/internal/application/bill_payments"
	"payme/internal/config"

	"payme/internal/application/notifications"
	"payme/internal/application/pendingcard"

	"payme/internal/application/savings"
	"payme/internal/application/splits"
	"payme/internal/application/token"
	"payme/internal/application/transaction"
	"payme/internal/application/transactionpin"
	"payme/internal/application/transfer"
	"payme/internal/application/user"
	"payme/internal/application/wallet"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("../../.env")
	config.Connect()
	config.GetDB().AutoMigrate(&user.User{},&wallet.Wallet{},&savings.PersonalSaving{},&savings.GroupSaving{},&splits.SplitBill{},&splits.SplitBillParticipants{},&savings.PersonalSaving{},&savings.GroupSaving{},&savings.GroupMember{},&savings.GroupContribution{},&transactionpin.TransactionPin{},&wallet.SavingsWallet{},&accounts.VirtualAccount{},&notifications.Notification{},&transfer.Transfer{},&bill_payments.BillPayment{},&token.CardToken{},&pendingcard.PendingCard{},&transaction.Transaction{})
	r := mux.NewRouter()
	routes.SetupRoutes(r)
	http.Handle("/", r)
	log.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))

}
