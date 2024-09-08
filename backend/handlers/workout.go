package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/emerconnelly/workout-tracker/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type WorkoutHandler struct {
	Collection *mongo.Collection
}

func (h *WorkoutHandler) ListWorkouts(w http.ResponseWriter, r *http.Request) {
	log.Println("ListWorkouts: called")

	// Find all documents in the MongoDB collection
	var workouts []models.Workout
	cursor, err := h.Collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx) // Close the cursor when the function returns

	// Decode each document into a Workout struct and append it to the Workouts struct slice
	for cursor.Next(ctx) {
		var workout models.Workout
		if err := cursor.Decode(&workout); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		workouts = append(workouts, workout)
	}

	// Return the Workouts struct slice as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(workouts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *WorkoutHandler) CreateWorkout(w http.ResponseWriter, r *http.Request) {
	log.Println("CreateWorkout: called")

	// Decode the request body into a Workout struct
	var workout models.Workout
	if err := json.NewDecoder(r.Body).Decode(&workout); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the Workout struct
	switch {
	case workout.Name == "":
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	// Insert the Workout struct into the MongoDB collection
	result, err := h.Collection.InsertOne(ctx, workout)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set Workout struct ID to the MongoDB document ObjectID
	workout.ID = result.InsertedID.(primitive.ObjectID)

	// Return the Workout struct as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(workout); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
