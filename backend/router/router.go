package router

import (
	"net/http"

	"github.com/emerconnelly/workout-tracker/handlers"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRoutes(client *mongo.Client) *http.ServeMux {
	// Exercise handler
	exerciseCollection := client.Database("workout_tracker").Collection("exercises")
	exerciseHandler := &handlers.ExerciseHandler{Collection: exerciseCollection}

	// Workout handler
	workoutCollection := client.Database("workout_tracker").Collection("workouts")
	workoutHandler := &handlers.WorkoutHandler{Collection: workoutCollection}

	mux := http.NewServeMux()

	// Exercise routes
	mux.HandleFunc("GET /api/exercises/", exerciseHandler.ListExercises)
	mux.HandleFunc("GET /api/exercise/{id}/", exerciseHandler.GetExercise)
	mux.HandleFunc("POST /api/exercise/", exerciseHandler.CreateExercise)
	mux.HandleFunc("PATCH /api/exercise/{id}/", exerciseHandler.UpdateExercise)
	mux.HandleFunc("DELETE /api/exercise/{id}/", exerciseHandler.DeleteExercise)
	mux.HandleFunc("DELETE /api/exercises/", exerciseHandler.DeleteExercises)

	// Workout routes
	mux.HandleFunc("GET /api/workouts/", workoutHandler.ListWorkouts)
	mux.HandleFunc("POST /api/workout/", workoutHandler.CreateWorkout)

	return mux
}
