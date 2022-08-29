package main

import (
	"fmt"
	"context"
	"net/http"
	"time"
	"encoding/json"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var collection *mongo.Collection

// define upfront the data struct 'blueprint'
// map the fields to match with api request and mongodb
type Person struct {
	ID			primitive.ObjectID	`json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname	string				`json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname	string				`json:"lastname,omitempty" bson:"lastname,omitempty"`
}

func CreatePersonEndPoint(response http.ResponseWriter, request *http.Request) {
	// define the response content-type as json
	response.Header().Set("content-type", "application/json")

	// define a variable that will receive the request data
	var person Person

	// decode request data into person variable
	json.NewDecoder(request.Body).Decode(&person)

	// define a context that carries the time that will be used as limit to db operation attempt
	// skip the callback function since errors will be handled if find does not match
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// try to insert data into db
	// if fails, messege error is returned
	result, err := collection.InsertOne(ctx, person)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}

	// encode the insertion message return into json formart and append it to the response
	json.NewEncoder(response).Encode(result)
}

func GetPeopleEndpoint(response http.ResponseWriter, request *http.Request) { 
	// define the response content-type as json
	response.Header().Set("content-type", "application/json")

	// define a 'dynamically-sized array' that will receive queries from db
	var people []Person

	// define a context that carries the time that will be used as limit to db operation attempt
	// skip the callback function since errors will be handled if find does not match
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// get a cursor 'pointer' to the entries in db with an empty filter 'bson.M{}'
	// if fails, messege error is returned
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}

	//if not fail will close cursor at the end of fucntion's scope
	defer cursor.Close(ctx)

	// iteraates through cursor and appen to the people slice
	// if fails, messege error is returned
	if err := cursor.All(ctx, &people); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}

	// encode people slice into json format and append it to the response
	json.NewEncoder(response).Encode(people)
}

func main() {
	fmt.Println("application has started...")

	// define a context that carries the time that will be used as limit to db connection attempt
	// define a function callback in case timeout is reached
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	
	// will release the resources if application hits time out
	// this is going to be executed when the function reaches the end of its scope
	defer cancel()

	// try to establish connection with the mongodb
	// if fails, program will finish execution
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:example@mongodb:27017"))
	if err != nil {
		panic(err)
	}

	// will disconnect from database at the end of program execution
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// define which database and which collection 'table' is going to be accessed
	collection = client.Database("mymongdb").Collection("people")

	// define a router to handle requests
	router := mux.NewRouter()

	// define the specific handlers to each request and endpoint
	router.HandleFunc("/person", CreatePersonEndPoint).Methods("POST")
	router.HandleFunc("/people", GetPeopleEndpoint).Methods("GET")

	// loop to listen and respond for requests
	http.ListenAndServe(":8000", router)
}
