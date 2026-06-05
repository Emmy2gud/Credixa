package routes

import (
	"payme/internal/api/middleware"
	"payme/internal/application/accounts"
	"payme/internal/application/auth"
	"payme/internal/application/bill_payments"
	"payme/internal/application/notifications"
	"payme/internal/application/splits"
	"payme/internal/application/transaction"

	"payme/internal/application/transfer"
	"payme/internal/application/webhooks"

	"payme/internal/application/transactionpin"
	"payme/internal/application/wallet"
	"payme/internal/config"

	"github.com/gorilla/mux"
)

func SetupRoutes(router *mux.Router) {

	webhookSvc := webhooks.NewWebhookService()
	webhookController := webhooks.NewWebhookController(webhookSvc)

	walletSvc := wallet.NewWalletService(config.DB)
	walletController := wallet.NewWalletController(walletSvc)

	accountsSvc := accounts.NewVirtualAccountService(config.DB)
	accountsController := accounts.NewVirtualAccountController(accountsSvc)

	transferSvc := transfer.NewTransferService(config.DB)
	transferController := transfer.NewTransferController(transferSvc)

	transactionSvc := transaction.NewTransactionService(config.DB)
	transactionController := transaction.NewTransactionController(transactionSvc)

	notificationSvc := notifications.NewNotificationService(config.DB)
	notificationController := notifications.NewNotificationController(notificationSvc)

	billPaymentsSvc := bill_payments.NewBillPaymentService(config.DB)
	billPaymentsController := bill_payments.NewBillPaymentController(billPaymentsSvc)

	splitSvc := splits.NewSplitService(config.DB)
	splitController := splits.NewSplitController(splitSvc)

	router.HandleFunc("/webhooks/flutterwave", webhookController.FlutterwaveWebhook).Methods("POST")
	UserWallet := router.PathPrefix("/wallet").Subrouter()
	UserWallet.Use(middleware.AuthMiddleware)
	UserWallet.HandleFunc("/balance", walletController.GetWalletBalance).Methods("GET")
	UserWallet.HandleFunc("/fund", walletController.InitiateWalletFunding).Methods("POST")
	UserWallet.HandleFunc("/fund/authorize", walletController.AuthorizeCardFunding).Methods("POST")
	UserWallet.HandleFunc("/fund/validate", walletController.ValidateWalletFunding).Methods("POST")
	UserWallet.HandleFunc("/fund/verify/{id}", walletController.VerifyCardCharge).Methods("GET")

	//transaction pin routes
	pinSvc := transactionpin.NewTransactionPinService(config.DB)
	pinHandler := transactionpin.NewTransactionPinController(pinSvc)
	TransactionPin := router.PathPrefix("/transaction-pin").Subrouter()
	TransactionPin.Use(middleware.AuthMiddleware)
	TransactionPin.HandleFunc("/create/{userid}", pinHandler.CreateTransactionPin).Methods("POST")
	// TransactionPin.HandleFunc("/verify/{userid}", pinHandler.VerifyTransactionPin).Methods("POST")
	// TransactionPin.HandleFunc("/update/{userid}", transactionpin.UpdateTransactionPin).Methods("PUT")
	// TransactionPin.HandleFunc("/delete/{userid}", transactionpin.DeleteTransactionPin).Methods("DELETE")

	//subscription routes for data,airtime,dstv,gotv,startimes,spectranet,smile,swift,electricity
	Subscription := router.PathPrefix("/subscription").Subrouter()
	Subscription.Use(middleware.AuthMiddleware)
	Subscription.HandleFunc("/biller-payments", billPaymentsController.BillerCategories).Methods("GET")
	Subscription.HandleFunc("/biller-payments/{category}", billPaymentsController.BillerCategory).Methods("GET")
	Subscription.HandleFunc("/bill-payments/{category}", billPaymentsController.BillCategory).Methods("GET")
	//verify tv and electricity details
	Subscription.HandleFunc("/bill-payments/verify/{serviceid}/payments", billPaymentsController.VerifySubscription).Methods("POST")
	//collecting itemcode and number to validate
	Subscription.HandleFunc("/bill-payments/create/{serviceid}/payments/{variationcode}", billPaymentsController.CreateBillPayment).Methods("POST")

	//virtual account creation users
	VirtualAccount := router.PathPrefix("/virtual-account").Subrouter()
	VirtualAccount.Use(middleware.AuthMiddleware)
	VirtualAccount.HandleFunc("/create", accountsController.CreateVirtualAccount).Methods("POST")

	//transaction history
	Transaction := router.PathPrefix("/transaction").Subrouter()
	Transaction.Use(middleware.AuthMiddleware)
	Transaction.HandleFunc("/history", transactionController.GetTransactionHistory).Methods("GET")
	Transaction.HandleFunc("/{id}", transactionController.GetTransactionByID).Methods("GET")
	Transaction.HandleFunc("/wallet", transactionController.GetWalletLogs).Methods("GET")

	//transfer routes
	Transfer := router.PathPrefix("/transfer").Subrouter()
	Transfer.Use(middleware.AuthMiddleware)
	Transfer.HandleFunc("/resolve-bank-details", transferController.ResolveBankDetails).Methods("POST")
	Transfer.HandleFunc("/initialize", transferController.InitializeFunding).Methods("POST")
	Transfer.HandleFunc("/verify", transferController.VerifyFunding).Methods("POST")

	// ── Split Bill routes ──────────────────────────────────────────────────────
	Splits := router.PathPrefix("/splits").Subrouter()
	Splits.Use(middleware.AuthMiddleware)
	Splits.HandleFunc("", splitController.CreateSplitBill).Methods("POST")              // POST   /splits
	Splits.HandleFunc("", splitController.GetSplitBills).Methods("GET")                 // GET    /splits
	Splits.HandleFunc("/{id}/accept", splitController.AcceptSplitBill).Methods("PUT")   // PUT    /splits/{id}/accept
	Splits.HandleFunc("/{id}/decline", splitController.DeclineSplitBill).Methods("PUT") // PUT    /splits/{id}/decline

	// ── Notification routes ────────────────────────────────────────────────────
	Notifications := router.PathPrefix("/notifications").Subrouter()
	Notifications.Use(middleware.AuthMiddleware)
	Notifications.HandleFunc("", notificationController.GetNotifications).Methods("GET")            // GET  /notifications
	Notifications.HandleFunc("/read", notificationController.MarkNotificationRead).Methods("PATCH") // PATCH /notifications/read

	authSvc := auth.NewAuthService(config.DB)
	authHandler := auth.NewAuthHandler(authSvc)

	router.HandleFunc("/dev/seed/users", auth.SeedFakeUsers).Methods("POST")
	router.HandleFunc("/login", authHandler.Login).Methods("POST")
	router.HandleFunc("/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/logout", auth.Logout).Methods("POST")
	//reset password
	router.HandleFunc("/forgot-password", authHandler.ForgotPassword).Methods("POST")
	router.HandleFunc("/reset-password", authHandler.ResetPassword).Methods("POST")

}
