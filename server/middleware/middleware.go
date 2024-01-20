package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sixam/go-ts-react-todo/models"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

func init() {
	loadEnv()
	createDBInstance()
}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading the .env file", err)
	}
}

func createDBInstance() {
	connectionStr := os.Getenv("DB_URI")
	dbName := os.Getenv("DB_NAME")
	collectionName := os.Getenv("DB_COLLECTION_NAME")

	clientOptions := options.Client().ApplyURI(connectionStr)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
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

	var task models.ToDoList
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

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Header", "Content-Type")

	params := mux.Vars(r)
	deleteOneTask(params["id"])
}

func DeleteAllTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	count := deleteAllTasksInDB()
	json.NewEncoder(w).Encode(count)
}

// Associated Helper Methods
func getAllTasksFromDB() []primitive.M {
	dbcursor, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	var results []primitive.M

	for dbcursor.Next(context.Background()) {
		var result bson.M
		dberror := dbcursor.Decode(&result)
		if dberror != nil {
			log.Fatal(dberror)
		}

		results = append(results, result)
	}

	if err := dbcursor.Err(); err != nil {
		log.Fatal(err)
	}

	dbcursor.Close(context.Background())
	return results
}

func updateTaskStatusInDB(task string) {
	id, _ := primitive.ObjectIDFromHex(task)
	retreivedTask := bson.M{"_id": id}
	updatedVal := bson.M{"$set": bson.M{"status": true}}
	result, err := collection.UpdateOne(context.Background(), retreivedTask, updatedVal)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Updated task counter: ", result.ModifiedCount)
}

func insertOneTask(task models.ToDoList) {
	inserted, err := collection.InsertOne(context.Background(), task)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Added one task!", inserted.InsertedID)
}

func deleteOneTask(task string) {
	fmt.Println(task)
	id, _ := primitive.ObjectIDFromHex(task)
	retreivedTask := bson.M{"_id": id}
	deleted, err := collection.DeleteOne(context.Background(), retreivedTask)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Deleted task", deleted.DeletedCount)
}

func deleteAllTasksInDB() int64 {
	deleted, err := collection.DeleteMany(context.Background(), bson.D{{}}, nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Deleted tasks", deleted.DeletedCount)
	return deleted.DeletedCount
}
