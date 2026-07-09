package config

import "os"

var (
	BackendCallbackURL = os.Getenv("KYC_BACKEND_CALLBACK_URL")
	FrontendSuccessURL = os.Getenv("KYC_FRONTEND_SUCCESS_URL")
	FrontendFailedURL  = os.Getenv("KYC_FRONTEND_FAILED_URL")
)