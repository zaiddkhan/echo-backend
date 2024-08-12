package controller

import (
	"Echo/firebase"
	"Echo/mongo/models"
	"context"
	"firebase.google.com/go/v4/messaging"
	"fmt"
	"github.com/madflojo/tasks"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type SchedulerRepository struct {
	scheduler       *tasks.Scheduler
	collection      *mongo.Collection
	messagingClient *messaging.Client
}

func NewSchedulerRepository(scheduler *tasks.Scheduler, collection *mongo.Collection) *SchedulerRepository {
	client, err := firebase.FirebaseInit(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	return &SchedulerRepository{
		scheduler,
		collection,
		client,
	}
}

func (s *SchedulerRepository) SendPushNotification(context context.Context, title, description, token string) {

	response, err := s.messagingClient.Send(context, &messaging.Message{
		Notification: &messaging.Notification{
			Title: title,
			Body:  description,
		},
		Token: token,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response)

}
func (s *SchedulerRepository) AddTask(ctx context.Context, task *models.Task, userRepo *UserRepository) {
	var duration time.Duration
	switch task.Unit {
	case "second":
		duration = time.Duration(task.Interval) * time.Second
	case "minute":
		duration = time.Duration(task.Interval) * time.Minute
	case "hour":
		duration = time.Duration(task.Interval) * time.Hour
	default:
		panic("invalid unit")

	}

	user, errUser := userRepo.FindUserByID(
		ctx, task.UserId,
	)
	if errUser != nil {
		log.Fatal(errUser)
	}

	_, errInsert := s.collection.InsertOne(ctx, task)
	if errInsert != nil {
		fmt.Println(errInsert)
	}

	_, err := s.scheduler.Add(&tasks.Task{
		Interval: duration,
		TaskFunc: func() error {
			// Put your logic here
			s.SendPushNotification(
				context.Background(),
				task.Title,
				task.Description,
				user.FirebaseToken,
			)
			return nil
		},
	})

	if err != nil {
		panic(err)
	}

}
