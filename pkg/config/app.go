package config

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := "root:@tcp(localhost:3306)/payme?charset=utf8mb4&parseTime=True&loc=Local"
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Database connection error:", err)
		panic("Failed to connect to database!")
	}
	DB = database
}

func GetDB() *gorm.DB {
	return DB
}
