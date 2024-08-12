package main

import (
	"Echo/api/controller"
	"Echo/api/route"
	"Echo/firebase"
	"Echo/mongo"
	"context"
	"github.com/madflojo/tasks"

	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	_, err = firebase.FirebaseInit(context.Background())

	if err != nil {
		fmt.Printf("error connecting to firebase: %v", err)
		os.Exit(1)
	}
	_, err = mongo.GetMongoClient()
	if err != nil {
		log.Fatal(err)
	}

	scheduler := tasks.New()

	schedulerRepo := controller.NewSchedulerRepository(scheduler, mongo.GetCollection("tasks"))

	userRepo := controller.NewUserRepository(mongo.GetCollection("users"))
	indexError := userRepo.CreateTtlIndex()
	if indexError != nil {
		fmt.Println(indexError)
	}
	router := gin.Default()
	route.SchedulerRoutes(router, schedulerRepo)
	route.UserRoutes(router, userRepo)
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
	log.Println("hehe")

}
