package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/joho/godotenv"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson/mongo"
	"go.mongodb.org/mong-driver/mongo/options"
)

var *mongo.Collection

func init() {
	loadEnv()
	createDBInstance()
}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading the .env file")
	}
}

func createDBInstance() {
	connectionStr := os.Getenv("DB_URI")
	dbName := os.Getenv("DB_NAME")
	collectionName := os.Getenv("DB_COLLECTION_NAME")

	clientOptions := options.Client().ApplyURL(connectionStr)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nill {
		log.Fatal(err) {
		}
	}

	client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Huzzah! Connected to mongoDB!")

	collection = client.Database(dbName).Collection(collectionName)
	fmt.Println("Collection instance created")
}

func GetAllTasks(w http.ResponseWriter, r *http.Request) {}

func CreateTask() {}

func TaskComplete() {}

func DeleteTask() {}

func DeleteAllTasks() {}
