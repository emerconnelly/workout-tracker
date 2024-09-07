package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/emerconnelly/workout-tracker/models"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestListExercises(t *testing.T) {
	ctx := context.Background()

	// Start a new MongoDB container
	mongodbContainer, err := mongodb.Run(ctx, "mongo:7")
	if err != nil {
		t.Fatal("failed to start MongoDB container:", err)
	}

	// Ensure the container is terminated when the test finishes
	defer func() {
		if err := mongodbContainer.Terminate(ctx); err != nil {
			t.Fatal("failed to terminate MongoDB container:", err)
		}
	}()

	// Connect to the MongoDB container
	mongoEndpoint, err := mongodbContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatal("failed to get MonogoDB connection string:", err)
	}
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoEndpoint))
	if err != nil {
		t.Fatal("failed to connect to MongoDB:", err)
	}

	// Create a new MongoDB collection
	collection := mongoClient.Database("test").Collection("exercises")

	// Insert test data
	testExercises := []interface{}{
		models.Exercise{Name: "Push-ups", MuscleGroup: "Chest"},
		models.Exercise{Name: "Squats", MuscleGroup: "Legs"},
	}
	_, err = collection.InsertMany(ctx, testExercises)
	if err != nil {
		t.Fatal("failed to insert test data in MongoDB collection:", err)
	}
	t.Log("inserted test data into MongoDB collection")

	/// Create an instance of ExerciseHandler
	exerciseHandler := &ExerciseHandler{
		Collection: collection,
	}

	// Create and execute the HTTP request
	req, _ := http.NewRequest("GET", "/exercises", nil)
	rr := httptest.NewRecorder()
	http.HandlerFunc(exerciseHandler.ListExercises).ServeHTTP(rr, req)

	// Check the response
	status := rr.Code
	if status != http.StatusOK {
		t.Errorf("handler status code returned %v, want %v", status, http.StatusOK)
	}
	t.Logf("wanted response code %v, got %v", http.StatusOK, status)

	// Read the response body
	body, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Fatal("failed to read response body:", err)
	}
	t.Log("response body read successfully")

	// Parse the response body
	var exercises []models.Exercise
	if err := json.Unmarshal(body, &exercises); err != nil {
		t.Fatal("failed to unmarshal response body:", err)
	}

	// Check the number of exercises returned
	if len(exercises) != 2 {
		t.Errorf("expected %d exercises, got %d", len(testExercises), len(exercises))
	}
	t.Log("expected 2 exercises, got", len(exercises))

	// Check the content of the exercises
	for i, exercise := range exercises {
		expected := testExercises[i].(models.Exercise)
		switch {
		case exercise.Name != expected.Name:
			t.Errorf("exercise %d does not match: got %v, want: %v", i, exercise.Name, expected.Name)
		case exercise.MuscleGroup != expected.MuscleGroup:
			t.Errorf("exercise %d does not match: got %v, want: %v", i, exercise.MuscleGroup, expected.MuscleGroup)
		}
	}
	t.Log("exercises match the test data")
}

// test comment
