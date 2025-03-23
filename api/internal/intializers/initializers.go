package initializers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func EnvironmentVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		panic(err)
	}

	print(os.Getenv("JWT_SECRET_KEY"))
}
