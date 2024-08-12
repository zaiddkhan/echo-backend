package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Task struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	UserId      primitive.ObjectID `json:"user_id" bson:"user_id"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	Interval    int                `json:"interval" bson:"interval"`
	Unit        string             `json:"unit" bson:"unit"`
	Count       int                `json:"count,omitempty" bson:"count,omitempty"`
}
