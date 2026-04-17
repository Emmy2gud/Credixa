package bill_payments

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"payme/pkg/middleware"

	"github.com/gorilla/mux"
)

// BillPaymentController handles HTTP requests for bill payments
type BillPaymentPayload struct {
	Phone         string `json:"phone"`
	Amount        string `json:"amount"`
	ServiceId     string `json:"serviceID"`
	RequestId     string `json:"request_id"`
	BillerCode    string `json:"billersCode"`
	VariationCode string `json:"variation_code"`
}

func BillerCategories(w http.ResponseWriter, r *http.Request) {
	// userID, _ := middleware.GetUserID(r)
	req, err := http.NewRequest("GET", "https://vtpass.com/api/service-categories", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Set("api-key", os.Getenv("VTPASS_API_KEY"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Flutterwave response:", string(body))

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func BillerCategory(w http.ResponseWriter, r *http.Request) {
	// Implementation will go here
	vars := mux.Vars(r)
	categoryId := vars["category"]

	req, err := http.NewRequest("GET", "https://vtpass.com/api/services?identifier="+categoryId, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Set("api-key", os.Getenv("VTPASS_API_KEY"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Flutterwave response:", string(body))

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func BillCategory(w http.ResponseWriter, r *http.Request) {
	// Implementation will go here
	vars := mux.Vars(r)
	categoryId := vars["category"]

	req, err := http.NewRequest("GET", "https://vtpass.com/api/service-variations?serviceID="+categoryId, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Set("api-key", os.Getenv("VTPASS_API_KEY"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Flutterwave response:", string(body))

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func CreateBillPayment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceid := vars["serviceid"]
	variationcode := vars["variationcode"]
	userID, _ := middleware.GetUserID(r)

	var input struct {
		Amount string `json:"amount"`
		Phone  string `json:"phone"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	respbody, err := ProcessBillPayment(userID, serviceid, variationcode, input.Phone, input.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respbody)
}
