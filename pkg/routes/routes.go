package routes

import (
	"payme/pkg/accounts"
	"payme/pkg/auth"
	"payme/pkg/bill_payments"
	"payme/pkg/webhooks"

	"payme/pkg/middleware"

	"payme/pkg/transactionpin"
	"payme/pkg/wallet"

	"github.com/gorilla/mux"
)

func SetupRoutes(router *mux.Router) {

	router.HandleFunc("/webhooks/flutterwave", webhooks.FlutterwaveWebhook).Methods("POST")
	UserWallet := router.PathPrefix("/wallet").Subrouter()
	UserWallet.Use(middleware.AuthMiddleware)
	UserWallet.HandleFunc("/balance", wallet.GetWalletBalance).Methods("GET")
	UserWallet.HandleFunc("/fund", wallet.InitiateWalletFunding).Methods("POST")
	UserWallet.HandleFunc("/fund/authorize", wallet.AuthorizeCardFunding).Methods("POST")
	UserWallet.HandleFunc("/fund/validate", wallet.ValidateWalletFunding).Methods("POST")
	UserWallet.HandleFunc("/fund/verify/{id}", wallet.VerifyCardCharge).Methods("GET")

	//transaction pin routes
	TransactionPin := router.PathPrefix("/transaction-pin").Subrouter()
	TransactionPin.Use(middleware.AuthMiddleware)
	TransactionPin.HandleFunc("/create/{userid}", transactionpin.CreateTransactionPin).Methods("POST")
	TransactionPin.HandleFunc("/verify/{userid}", transactionpin.VerifyTransactionPin).Methods("POST")
	TransactionPin.HandleFunc("/update/{userid}", transactionpin.UpdateTransactionPin).Methods("PUT")
	TransactionPin.HandleFunc("/delete/{userid}", transactionpin.DeleteTransactionPin).Methods("DELETE")

	//subscription routes for data,airtime,dstv,gotv,startimes,spectranet,smile,swift,electricity
	Subscription := router.PathPrefix("/subscription").Subrouter()
	Subscription.Use(middleware.AuthMiddleware)
	Subscription.HandleFunc("/biller-payments", bill_payments.BillerCategories).Methods("GET")
	Subscription.HandleFunc("/biller-payments/{category}", bill_payments.BillerCategory).Methods("GET")
	Subscription.HandleFunc("/bill-payments/{category}", bill_payments.BillCategory).Methods("GET")
	//collecting itemcode and number to validate
	// Subscription.HandleFunc("/bill-payments/validate/{itemcode}", bill_payments.ValidateBillerDetails).Methods("POST")
	Subscription.HandleFunc("/bill-payments/create/{serviceid}/payments/{variationcode}", bill_payments.CreateBillPayment).Methods("POST")
	//virtual account creation users
	VirtualAccount := router.PathPrefix("/virtual-account").Subrouter()
	VirtualAccount.Use(middleware.AuthMiddleware)
	VirtualAccount.HandleFunc("/create",accounts.CreateVirtualAccount).Methods("POST")

	router.HandleFunc("/register", auth.Register).Methods("POST")
	router.HandleFunc("/login", auth.Login).Methods("POST")
	router.HandleFunc("/logout", auth.Logout).Methods("POST")
	//reset password
	router.HandleFunc("/forgot-password", auth.ForgotPassword).Methods("POST")
	router.HandleFunc("/reset-password", auth.ResetPassword).Methods("POST")

}
