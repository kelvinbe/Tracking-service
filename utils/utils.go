package utils

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	env := os.Getenv("GO_ENV")
	fmt.Printf("GO_ENV: %s", env)
	fmt.Println("Other env variables:::")
	fmt.Println(os.Getenv("DB_HOST"))
	fmt.Println(os.Getenv("DB_PORT"))
	fmt.Println(os.Getenv("DB_USERNAME"))
	fmt.Println(os.Getenv("DB_PASSWORD"))
	fmt.Println(os.Getenv("DB_NAME"))
	fmt.Println(os.Getenv("DB_SSLMODE"))
	
	if env == "" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Error loading .env file in development mode", err)
		}
	}
	// if production do nothin, env will be readily available via os.Getenv
}

// returns the coordinates from a string latitude then longitude
func ExtractCoordsFromText(text string) (string, string) {
	// TODO: unsure of the format of the text, will make necessary modifications once sure
	coords := strings.Split(text, " ");
	return coords[0], coords[1]
}
