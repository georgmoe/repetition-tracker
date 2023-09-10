package configs

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func EnvMongoURI() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	mongoUser := os.Getenv("MONGOUSER")
	mongoPass := os.Getenv("MONGOPASS")
	mongoPort := os.Getenv("MONGOPORT")

	return fmt.Sprintf("mongodb://%v:%v@localhost:%v", mongoUser, mongoPass, mongoPort)
}
