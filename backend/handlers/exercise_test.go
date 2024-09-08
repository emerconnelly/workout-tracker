package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/emerconnelly/workout-tracker/models"
)

func TestListExercises(t *testing.T) {
	// Create and execute the HTTP request
	req, _ := http.NewRequest("GET", "/api/exercises", nil)
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
	if len(exercises) != len(testExercises) {
		t.Errorf("expected %d exercises, got %d", len(testExercises), len(exercises))
	}
	t.Logf("expected %d exercises, got %d", len(testExercises), len(exercises))

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

func TestGetExercise(t *testing.T) {
	// Create and execute the HTTP request
	req, _ := http.NewRequest("GET", "/api/exercise/{id}", nil)
	rr := httptest.NewRecorder()
	http.HandlerFunc(exerciseHandler.ListExercises).ServeHTTP(rr, req)
}
