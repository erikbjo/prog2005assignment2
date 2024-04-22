package config

import (
	"assignment-2/internal/constants"
	"github.com/joho/godotenv"
	"log"
)

func InitConfig() error {
	// Load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Println(constants.ErrLoadingEnvFile + err.Error())
		return err
	}

	return nil
}
