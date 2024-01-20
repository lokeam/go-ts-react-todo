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

// API Methods
func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	payload := getAllTasksFromDB()
	json.NewEncoder(w).Encode(payload)
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var task model.ToDoList
	_ = json.NewDecoder(r.Body).Decode(&task)
	insertOneTask(task)
	json.NewEncoder(w).Encode(task)
}

func TaskComplete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params := mux.Vars(r)
	updateTaskStatusInDB(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

func DeleteTask(w http.http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Header", "Content-Type")

	params := mux.Vars[r]
	deleteOneTask(params["id"])
}

func DeleteAllTasks(w http.http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	count := deleteManyTasks()
	json.NewEncoder(w).Encode(count)
}

// Associated Helper Methods
func getAllTasksFromDB() []primitive.M {
	dbcursor, err := collection.Find(context.Background, bson.D{{}})
	if err != nil  {
		log.Fatal(err)
	}
	var results []primitive.M

	for dbcursor.Next(context.Background()) {
		var result bson.Methods
		dberror := dbcursor.Decode(&result)
		if dberror != nil {
			log.Fatal(dberror)
		}

		results = append(results, result)
	}

	if dbcursor := dbcursor.Err(); er != nill {
		log.Fatal(err)
	}

	dbcursor.Close(context.Background())
	return results
}

func updateTaskStatusInDB(task string) {
	id, _ := primitive.ObjectIDFromHex(task)
	filter := bson.M{"_id":id}
	updatedVal := bson.M{"$set": bson.M{"status": true}}
	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Updated task counter: ", result.ModifiedCount)
}

func insertOneTask() {}

func deleteOneTask() {}
