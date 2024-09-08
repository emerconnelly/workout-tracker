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

type ExerciseHandler struct {
	Collection *mongo.Collection
}

func (h *ExerciseHandler) ListExercises(w http.ResponseWriter, r *http.Request) {
	log.Println("ListExercises: called")

	// Find all documents in the MongoDB collection
	var exercises []models.Exercise
	cursor, err := h.Collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx) // Close the cursor when the function returns

	// Decode each document into an Exercise struct and append it to the Exercises struct slice
	for cursor.Next(ctx) {
		var exercise models.Exercise
		if err := cursor.Decode(&exercise); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		exercises = append(exercises, exercise)
	}
	log.Printf("ListExercises: found %d exercises\n", len(exercises))

	// Return the Exercises struct slice as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(exercises); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("ListExercises: response sent successfully")
}

func (h *ExerciseHandler) GetExercise(w http.ResponseWriter, r *http.Request) {
	log.Println("GetExercise: called")

	// Get the id from the URLs
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	// Convert the id to a MongoDB ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Find the exercise by the ObjectID
	var exercise models.Exercise
	if err := h.Collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&exercise); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the exercise as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(exercise); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *ExerciseHandler) CreateExercise(w http.ResponseWriter, r *http.Request) {
	log.Println("CreateExercise: called")

	// Decode the JSON body into an Exercise struct
	exercise := new(models.Exercise)
	if err := json.NewDecoder(r.Body).Decode(&exercise); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Validate the Exercise struct
	switch {
	case exercise.Name == "":
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	case exercise.MuscleGroup == "":
		http.Error(w, "MuscleGroup is required", http.StatusBadRequest)
		return
	}

	// Insert the Exercise struct into the MongoDB collection
	result, err := h.Collection.InsertOne(ctx, exercise)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the Exercise struct ID to the MongoDB document ObjectID
	exercise.ID = result.InsertedID.(primitive.ObjectID)

	// Return the Exercise struct as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(exercise); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *ExerciseHandler) UpdateExercise(w http.ResponseWriter, r *http.Request) {
	log.Println("UpdateExercise: called")

	// Get the id from the URLs
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	// Convert the id to a MongoDB ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Decode the JSON body into an Exercise struct
	exercise := new(models.Exercise)
	if err := json.NewDecoder(r.Body).Decode(&exercise); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Validate the Exercise struct
	switch {
	case exercise.Name == "":
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	case exercise.MuscleGroup == "":
		http.Error(w, "MuscleGroup is required", http.StatusBadRequest)
		return
	}

	// Update the exercise by the ObjectID
	result, err := h.Collection.ReplaceOne(ctx, bson.M{"_id": objectID}, exercise)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the Exercise struct as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *ExerciseHandler) DeleteExercise(w http.ResponseWriter, r *http.Request) {
	log.Println("DeleteExercise: called")

	// Get the id from the URLs
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	// Convert the id to a MongoDB ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Delete the exercise by the ObjectID
	result, err := h.Collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return a success message
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *ExerciseHandler) DeleteExercises(w http.ResponseWriter, r *http.Request) {
	log.Println("DeleteExercises: called")

	// Delete all documents in the MongoDB collection
	result, err := h.Collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return a success message
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
