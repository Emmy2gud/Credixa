package splits

import (
	"encoding/json"
	"fmt"
	"net/http"
	"payme/internal/application/notifications"
	"payme/internal/api/middleware"
	"payme/internal/config"
	"payme/pkg/utils"
	"payme/internal/application/wallet"
	"strconv"

	"github.com/gorilla/mux"
)

func CreateSplitBill(w http.ResponseWriter, r *http.Request) {
	type ParticipantInput struct {
		UserID     uint64 `json:"user_id"`
		Amount     uint64 `json:"amount"`
		Percentage uint64 `json:"percentage"`
	}
	var input struct {
		Title             string             `json:"title"`
		Description       string             `json:"description"`
		TotalAmount       uint64             `json:"total_amount"`
		SplitType         string             `json:"split_type"`
		ParticipantsCount int                `json:"participants_count"`
		Participants      []ParticipantInput `json:"participants"`
	}

	utils.ParseBody(r, &input)

	// Get creator's user ID from JWT context
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	//  Build and save the SplitBill record
	splitbill := SplitBill{
		CreatorID:         uint64(userID),
		Title:             input.Title,
		Description:       input.Description,
		TotalAmount:       input.TotalAmount,
		SplitType:         input.SplitType,
		ParticipantsCount: input.ParticipantsCount,
		Status:            "pending",
	}

	if err := config.DB.Create(&splitbill).Error; err != nil {
		http.Error(w, "Could not create split bill", http.StatusInternalServerError)
		return
	}

	// Find the creator's own amount/percentage from the participants list (if provided)
	var creatorAmount, creatorPercentage uint64
	for _, p := range input.Participants {
		if p.UserID == uint64(userID) {
			creatorAmount = p.Amount
			creatorPercentage = p.Percentage
			break
		}
	}

	//  Add the creator as the first participant
	creatorParticipant := SplitBillParticipants{
		SplitBillID: splitbill.ID,
		UserID:      uint64(userID),
		Amount:      creatorAmount,
		Percentage:  creatorPercentage,
		Status:      "accepted", // creator auto-accepts
	}
	if err := config.DB.Create(&creatorParticipant).Error; err != nil {
		http.Error(w, "Could not add creator as participant", http.StatusInternalServerError)
		return
	}

	//  Add every invited participant and notify them
	for _, participant := range input.Participants {
		if participant.UserID == uint64(userID) {
			continue // skip if creator accidentally included themselves
		}

		invitedParticipant := SplitBillParticipants{
			SplitBillID: splitbill.ID,
			UserID:      participant.UserID,
			Amount:      participant.Amount,
			Percentage:  participant.Percentage,
			Status:      "pending",
		}
		if err := config.DB.Create(&invitedParticipant).Error; err != nil {
			// log and continue — don't fail the whole request for one participant
			fmt.Printf("Could not add participant %d: %v\n", participant.UserID, err)
			continue
		}

		//  Notification 2: invite each participant to accept or decline
		inviteMsg := fmt.Sprintf(
			"You have been invited to a split bill: \"%s\". Please accept or decline.",
			splitbill.Title,
		)
		notifications.CreateNotification(
			uint(participant.UserID),
			"SplitBill",
			"Split Bill Invitation",
			inviteMsg,
		)
	}

	//   Notification 1: confirm to the creator that the split was created
	creatorMsg := fmt.Sprintf(
		"Your split bill \"%s\" has been created successfully. Waiting for participants to respond.",
		splitbill.Title,
	)
	notifications.CreateNotification(userID, "SplitBill", "Split Bill Created", creatorMsg)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Split bill created successfully",
		"split_id": splitbill.ID,
	})
}

func GetSplitBills(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var splitBills []SplitBill
	if err := config.DB.Where("creator_id = ?", userID).Find(&splitBills).Error; err != nil {
		http.Error(w, "Could not fetch split bills", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(splitBills)
}

// AcceptSplitBill lets a participant accept a split bill invitation.
// Route: PUT /splits/{id}/accept
func AcceptSplitBill(w http.ResponseWriter, r *http.Request) {
	
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	splitID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid split bill ID", http.StatusBadRequest)
		return
	}

	var participant SplitBillParticipants
	if err := config.DB.
		Where("split_bill_id = ? AND user_id = ?", splitID, userID).
		First(&participant).Error; err != nil {
		http.Error(w, "Participant record not found", http.StatusNotFound)
		return
	}

	participant.Status = "accepted"
	if err := config.DB.Save(&participant).Error; err != nil {
		http.Error(w, "Could not update participant status", http.StatusInternalServerError)
		return
	}

	//minus funds from participants wallet and then add the participants total amount to the creator wallet
	//fetching split content 
	var splitBill SplitBill
	if err := config.DB.Where("id = ?", splitID).First(&splitBill).Error; err != nil {
		http.Error(w, "failed to update wallet balance", http.StatusInternalServerError)
		return
	}
	//deducting the participant amount from his wallet
	if err := wallet.DeductWalletBalance(uint(participant.UserID), participant.Amount); err != nil {
		http.Error(w, "failed to update wallet balance", http.StatusInternalServerError)
		return
	}
	//add to the creator wallet

   fmt.Println("creator id",splitBill.CreatorID)


	if err := wallet.UpdateWalletBalance(uint(splitBill.CreatorID), participant.Amount); err != nil {
		http.Error(w, "failed to update wallet balance", http.StatusInternalServerError)
		return
	}
	// Fetch the split bill to notify the creator

	if err := config.DB.First(&splitBill, splitID).Error; err == nil {
		acceptMsg := fmt.Sprintf(
			"A participant has accepted your split bill: \"%s\".",
			splitBill.Title,
		)
		notifications.CreateNotification(uint(splitBill.CreatorID), "SplitBill", "Split Bill Accepted", acceptMsg)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Split bill accepted successfully",
	})
}

// DeclineSplitBill lets a participant decline a split bill invitation.
// Route: PUT /splits/{id}/decline
func DeclineSplitBill(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	splitID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid split bill ID", http.StatusBadRequest)
		return
	}

	var participant SplitBillParticipants
	if err := config.DB.
		Where("split_bill_id = ? AND user_id = ?", splitID, userID).
		First(&participant).Error; err != nil {
		http.Error(w, "Participant record not found", http.StatusNotFound)
		return
	}

	participant.Status = "declined"
	if err := config.DB.Save(&participant).Error; err != nil {
		http.Error(w, "Could not update participant status", http.StatusInternalServerError)
		return
	}

	// Fetch the split bill to notify the creator
	var splitBill SplitBill
	if err := config.DB.First(&splitBill, splitID).Error; err == nil {
		declineMsg := fmt.Sprintf(
			"A participant has declined your split bill: \"%s\".",
			splitBill.Title,
		)
		notifications.CreateNotification(uint(splitBill.CreatorID), "SplitBill", "Split Bill Declined", declineMsg)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Split bill declined successfully",
	})
}
