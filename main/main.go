package main

import (
	"Echo/api/controller"
	"Echo/api/route"
	"Echo/mongo"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"time"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	client, err := mongo.GetMongoClient()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Disconnect(ctx); err != nil {
		log.Fatal(err)
	}
	userRepo := controller.NewUserRepository(mongo.GetCollection("users"))

	router := gin.Default()
	route.UserRoutes(router, userRepo)
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
	log.Println("hehe")

}
