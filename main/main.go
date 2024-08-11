package main

import (
	"Echo/mongo"
	"github.com/joho/godotenv"
	"log"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	mongo.ConnectDB()

}
