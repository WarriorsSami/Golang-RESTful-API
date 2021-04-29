package main

import (
	m "RESTful-API-demo/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"time"
)

var client *mongo.Client

func getUsers(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var users []m.User

	dbName := os.Getenv("DB_NAME")
	clName := os.Getenv("COLLECTION_NAME")
	collection := client.Database(dbName).Collection(clName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		_, _ = response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		log.Fatal(err)
		return
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			_, _ = response.Write([]byte(`{"message": "` + err.Error() + `"}`))
			log.Fatal(err)
			return
		}
	}(cursor, ctx)

	for cursor.Next(ctx) {
		var user m.User
		err := cursor.Decode(&user)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			_, _ = response.Write([]byte(`{"message": "` + err.Error() + `"}`))
			log.Fatal(err)
			return
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		_, _ = response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		log.Fatal(err)
		return
	}
	err = json.NewEncoder(response).Encode(users)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		_, _ = response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		log.Fatal(err)
		return
	}
}

func getUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])

	var user m.User
	dbName := os.Getenv("DB_NAME")
	clName := os.Getenv("COLLECTION_NAME")
	collection := client.Database(dbName).Collection(clName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, m.User{ID: id}).Decode(&user)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		_, _ = response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		log.Fatal(err)
		return
	}
	err = json.NewEncoder(response).Encode(user)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		_, _ = response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		log.Fatal(err)
		return
	}
}

func createUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var user m.User
	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		_, _ = response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		log.Fatal(err)
		return
	}

	dbName := os.Getenv("DB_NAME")
	clName := os.Getenv("COLLECTION_NAME")
	collection := client.Database(dbName).Collection(clName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, _ := collection.InsertOne(ctx, user)
	err = json.NewEncoder(response).Encode(result)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		_, _ = response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		log.Fatal(err)
		return
	}
}

func updateUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var user m.User
	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		_, _ = response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		log.Fatal(err)
		return
	}

	dbName := os.Getenv("DB_NAME")
	clName := os.Getenv("COLLECTION_NAME")
	collection := client.Database(dbName).Collection(clName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": bson.M{"$eq": id}}
	update := bson.M{"$set": bson.M{
		"first_name": user.FirstName,
		"last_name": user.LastName,
		"password": user.Password,
		"email": user.Email,
	}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		_, _ = response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		log.Fatal(err)
		return
	}
	err = json.NewEncoder(response).Encode(result)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		_, _ = response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		log.Fatal(err)
		return
	}
}

func deleteUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])

	dbName := os.Getenv("DB_NAME")
	clName := os.Getenv("COLLECTION_NAME")
	collection := client.Database(dbName).Collection(clName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": bson.M{"$eq": id}}
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		_, _ = response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		log.Fatal(err)
		return
	}
	err = json.NewEncoder(response).Encode(result)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		_, _ = response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		log.Fatal(err)
		return
	}
}

func main() {
	fmt.Println("Starting server ...")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file!")
	}
	apiUri := os.Getenv("API_URI")
	mongoUri := os.Getenv("MONGO_URI")

	router := mux.NewRouter()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI(mongoUri))

	router.HandleFunc(fmt.Sprintf("%s", apiUri), getUsers).Methods("GET")
	router.HandleFunc(fmt.Sprintf("%s/{id}", apiUri), getUser).Methods("GET")
	router.HandleFunc(fmt.Sprintf("%s", apiUri), createUser).Methods("POST")
	router.HandleFunc(fmt.Sprintf("%s/{id}", apiUri), updateUser).Methods("PUT")
	router.HandleFunc(fmt.Sprintf("%s/{id}", apiUri), deleteUser).Methods("DELETE")

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
		return
	}
}