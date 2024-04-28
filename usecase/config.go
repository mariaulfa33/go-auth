package usecase

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}
