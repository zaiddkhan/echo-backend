package models

import (
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name          string             `json:"name" bson:"name" validate:"required,min=2,max=32"`
	Email         string             `json:"email" bson:"email" validate:"required,email"`
	FirebaseToken string             `json:"firebase_token,omitempty" bson:"firebase_token,omitempty"`
}

func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
