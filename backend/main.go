package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/emerconnelly/workout-tracker/router"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.Background()

func main() {
	log.Println("yo, we live")

	// Load the local .env file if not in a container
	if os.Getenv("IS_CONTAINER") != "true" {
		if err := godotenv.Load("../.env"); err != nil {
			log.Fatalln(err)
		}
		log.Print("loaded env vars from '.env' file")
	}

	// Get the MongoDB URI and port from the environment variables
	uri := os.Getenv("MONGODB_URI")
	port := os.Getenv("PORT")
	switch {
	case uri == "":
		log.Fatalln("MONGODB_URI not set")
	case port == "":
		port = "8080"
		log.Println("PORT not set, defaulting to 8080")
	}

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	} else {
		log.Println("connected to MongoDB")
	}
	defer client.Disconnect(ctx) // Close the MongoDB connection when the main function returns

	// Set up the routes
	mux := router.SetupRoutes(client)

	// Add CORS middleware
	handler := corsMiddleware(mux)

	// Start the server
	log.Fatalln(http.ListenAndServe("localhost:"+port, handler))
}
