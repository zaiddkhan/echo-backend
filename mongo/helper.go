package mongo

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	clientInstance *mongo.Client
	clientError    error
	mongoOnce      sync.Once
)

func GetMongoClient() (*mongo.Client, error) {

	mongoOnce.Do(func() {

		clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI"))
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Use the context created above
		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			clientError = err
			return
		}
		if err := client.Ping(ctx, readpref.Primary()); err != nil {
			clientError = err
			return
		}
		clientInstance = client
	})

	return clientInstance, clientError

}

func GetCollection(collectionName string) *mongo.Collection {
	client, err := GetMongoClient()
	if err != nil {
		return nil
	}
	return client.Database("echo").Collection(collectionName)
}

type ErrorResponse struct {
	StatusCode   int    `json:"status_code"`
	ErrorMessage string `json:"message"`
}

func GetError(err error, w http.ResponseWriter) {
	var response = ErrorResponse{
		ErrorMessage: err.Error(),
		StatusCode:   http.StatusInternalServerError,
	}
	message, _ := json.Marshal(response)

	w.WriteHeader(response.StatusCode)
	w.Write(message)
}
