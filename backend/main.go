package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	log.Println("Hello, World!")

	// Load the .env file if not in production
	if os.Getenv("ENV") != "production" {
		if err := godotenv.Load("../.env"); err != nil {
			log.Fatal(err)
		}
	}

	// Get the MongoDB URI and port from the environment variables
	uri := os.Getenv("MONGODB_URI")
	port := os.Getenv("PORT")
	switch {
	case uri == "":
		log.Fatal("MONGODB_URI is empty")
	case port == "":
		log.Fatal("PORT is empty")
	}

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	} else {
		log.Println("Connected to MongoDB")
	}
	defer client.Disconnect(context.TODO()) // Close the MongoDB connection when the main function returns

	// Set up the routes and start the server
	mux := setupRoutes(client)

	// Add CORS middleware
	handler := corsMiddleware(mux)

	log.Fatal(http.ListenAndServe("localhost:"+port, handler))
}
