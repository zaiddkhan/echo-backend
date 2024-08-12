package controller

import (
	mongo2 "Echo/mongo"
	"Echo/mongo/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) *UserRepository {
	return &UserRepository{collection}
}

func (r *UserRepository) GetAllUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {

		}
	}(cursor, ctx)
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepository) FindUserByID(ctx context.Context, id primitive.ObjectID) (models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err

}

func (r *UserRepository) UpsertUser(ctx context.Context, user models.User) error {
	opts := options.Replace().SetUpsert(true)
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": user.ID}, user, opts)
	return err
}

func (r *UserRepository) CreateTtlIndex() error {
	collection := mongo2.GetCollection("users")
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"timestamp": 1},
		Options: options.Index().SetExpireAfterSeconds(int32(time.Hour * 24 * 21 / time.Second)),
	}
	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	return err
}
