package token

import "gorm.io/gorm"

type CardToken struct {
	gorm.Model	
	ID        uint
	UserID    uint
	Token     string `json:"token"`
	CardBrand string `json:"card_brand"`
	Last4     string `json:"last_4digits"`
	Expiry    string `json:"expiry"`
	First6    string `json:"first_6digits"`
	Issuer    string `json:"issuer"`
	Country   string `json:"country"`
	Type      string `json:"type"`
}