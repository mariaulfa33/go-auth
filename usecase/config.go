package usecase

import (
	"fmt"

	"github.com/joho/godotenv"
)

func LoadEnv(envPath string) error {
	// load .env file
	err := godotenv.Load(envPath)

	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
