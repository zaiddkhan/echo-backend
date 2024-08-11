package mongo

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
)

func ConnectDB() {
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI"))
	_, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

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
