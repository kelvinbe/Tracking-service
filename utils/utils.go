package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	env := os.Getenv("GO_ENV")
	if env == "" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Error loading .env file in development mode")
		}
	}
	// if production do nothin, env will be readily available via os.Getenv
}