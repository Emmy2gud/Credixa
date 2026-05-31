package auth

import (
	"fmt"
	"net/http"
	"payme/internal/config"
	"payme/internal/application/user"

	"payme/internal/application/wallet"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

func SeedFakeUsers(w http.ResponseWriter, r *http.Request) {
	// Use a time-based seed so every run generates unique emails
	gofakeit.Seed(time.Now().UnixNano())

	created := 0
	for i := 0; i < 20; i++ {

		u := user.User{
			FullName: gofakeit.Name(),
			Email:    gofakeit.Email(),
			Password: gofakeit.Password(true, true, true, true, false, 12),
			Role:     "user",
		}

		// hashed, err := utils.HashPassword(u.Password)
		// if err != nil {
		// 	fmt.Printf("Seed: could not hash password for user %d: %v\n", i, err)
		// 	continue
		// }
		u.Password =u.Password

		if err := config.DB.Create(&u).Error; err != nil {
			fmt.Printf("Seed: could not create user %d (%s): %v\n", i, u.Email, err)
			continue
		}

		walletModel := wallet.Wallet{
			UserID:   u.ID,
			Balance:  int64(gofakeit.Number(1000, 50000)),
			Currency: "NGN",
			Status:   "active",
		}

		if err := config.DB.Create(&walletModel).Error; err != nil {
			fmt.Printf("Seed: could not create wallet for user %d: %v\n", u.ID, err)
		}

		created++
	}
    w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"message":"Seeded %d users successfully"}`, created)
}
