package main

import (
	"db-less-store/internal/auth"
	"db-less-store/internal/product"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(
		&product.Product{},
		&auth.Session{},
		&auth.User{})

	if err != nil {
		panic(err)
	}
}
