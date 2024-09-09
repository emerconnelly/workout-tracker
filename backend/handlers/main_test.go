package handlers

import (
	"log"
	"os"
	"testing"

	"github.com/emerconnelly/workout-tracker/models"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	testExercises = []interface{}{
		&models.Exercise{Name: "Shoulder Press", MuscleGroup: "Shoulders"},
		&models.Exercise{Name: "Bicep Curl", MuscleGroup: "Arms"},
		&models.Exercise{Name: "Push-up", MuscleGroup: "Chest"},
		&models.Exercise{Name: "Pull-up", MuscleGroup: "Back"},
		&models.Exercise{Name: "Squat", MuscleGroup: "Legs"},
	}

	client     *mongo.Client
	collection *mongo.Collection

	exerciseHandler = &ExerciseHandler{}
)

func TestMain(m *testing.M) {
	log.Println("TestMain: called")

	// Start a new MongoDB container
	container, err := mongodb.Run(ctx, "mongo:7")
	if err != nil {
		log.Fatal("TestMain: failed to start MongoDB container:", err)
	} else {
		log.Println("TestMain: started MongoDB container")
	}

	// Ensure the container is terminated when the test finishes
	defer func() {
		if err := container.Terminate(ctx); err != nil {
			log.Fatal("TestMain: failed to terminate MongoDB container:", err)
		} else {
			log.Println("TestMain: terminated MongoDB container")
		}
	}()

	// Connect to the MongoDB container
	endpoint, err := container.ConnectionString(ctx)
	if err != nil {
		log.Fatal("TestMain: failed to get MonogoDB connection string:", err)
	} else {
		log.Println("TestMain: got MongoDB connection string")
	}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(endpoint))
	if err != nil {
		log.Fatal("TestMain: failed to connect to MongoDB:", err)
	} else {
		log.Println("TestMain: connected to MongoDB")
	}

	// Create a new MongoDB collection
	collection = client.Database("test").Collection("exercises")

	// Create a new ExerciseHandler
	exerciseHandler = &ExerciseHandler{
		Collection: collection,
	}

	// Insert test data
	result, err := collection.InsertMany(ctx, testExercises)
	if err != nil {
		log.Fatal("TestMain: failed to insert test data in MongoDB collection:", err)
	} else {
		log.Println("TestMain: inserted test data into MongoDB collection")
	}

	// Update testExercises with the generated IDs
	for i, id := range result.InsertedIDs {
		testExercises[i].(*models.Exercise).ID = id.(primitive.ObjectID)
	}

	// Run the tests
	log.Println("TestMain: running tests")
	code := m.Run()

	// Disconnect from the MongoDB container
	if err := client.Disconnect(ctx); err != nil {
		log.Fatal("TestMain: failed to disconnect from MongoDB:", err)
	} else {
		log.Println("TestMain: disconnected from MongoDB")
	}

	os.Exit(code)
}
