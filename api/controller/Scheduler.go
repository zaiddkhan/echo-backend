package controller

import (
	"Echo/mongo/models"
	"context"
	"fmt"
	"github.com/madflojo/tasks"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type SchedulerRepository struct {
	scheduler  *tasks.Scheduler
	collection *mongo.Collection
}

func NewSchedulerRepository(scheduler *tasks.Scheduler, collection *mongo.Collection) *SchedulerRepository {
	return &SchedulerRepository{
		scheduler,
		collection,
	}
}

func (s *SchedulerRepository) AddTask(ctx context.Context, task *models.Task) {
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
	_, errInsert := s.collection.InsertOne(ctx, task)
	if errInsert != nil {
		fmt.Println(errInsert)
	}

	_, err := s.scheduler.Add(&tasks.Task{
		Interval: duration,
		TaskFunc: func() error {
			// Put your logic here
			fmt.Println("ddjdjd")
			return nil
		},
	})

	if err != nil {
		panic(err)
	}

}
