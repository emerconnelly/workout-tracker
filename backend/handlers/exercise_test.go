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
	req, _ := http.NewRequest("GET", "/api/exercises/", nil)
	rr := httptest.NewRecorder()
	mux := http.NewServeMux()
	mux.HandleFunc("/api/exercises/", exerciseHandler.ListExercises)
	mux.ServeHTTP(rr, req)

	// Check the response code
	code := rr.Code
	if code != http.StatusOK {
		t.Errorf("response code returned %v, want %v", code, http.StatusOK)
		return
	} else {
		t.Logf("response code good: %v", code)
	}

	// Read the response body
	body, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Fatal("failed to read response body:", err)
	} else {
		t.Log("response body read")
	}

	// Parse the response body
	var exercises []models.Exercise
	if err := json.Unmarshal(body, &exercises); err != nil {
		t.Fatal("failed to unmarshal response body:", err)
	} else {
		t.Log("response body unmarshalled")
	}

	// Check the number of exercises returned
	if len(exercises) != len(testExercises) {
		t.Errorf("expected %d exercises, got %d", len(testExercises), len(exercises))
		return
	} else {
		t.Logf("expected %d exercises, got %d", len(testExercises), len(exercises))
	}

	// Check the content of the exercises
	for i, exercise := range exercises {
		expected := testExercises[i].(*models.Exercise)
		t.Logf("checking exercise '%v'", exercise.Name)
		switch {
		case exercise.Name != expected.Name:
			t.Errorf("'Name' bad match: got %v, want: %v", exercise.Name, expected.Name)
		case exercise.MuscleGroup != expected.MuscleGroup:
			t.Errorf("'MuscleGroup' bad match: got %v, want: %v", exercise.MuscleGroup, expected.MuscleGroup)
		default:
			t.Logf("exercise matches test data")
		}
	}
}

func TestGetExercise(t *testing.T) {
	for _, testExerciseInterface := range testExercises {
		testExercise := testExerciseInterface.(*models.Exercise)
		t.Logf("testing exercise '%v' with %v", testExercise.Name, testExercise.ID)

		// Create and execute the HTTP request
		req, _ := http.NewRequest("GET", "/api/exercise/"+testExercise.ID.Hex()+"/", nil)
		rr := httptest.NewRecorder()
		mux := http.NewServeMux()
		mux.HandleFunc("/api/exercise/{id}/", exerciseHandler.GetExercise)
		mux.ServeHTTP(rr, req)

		// Check the response code
		code := rr.Code
		if code != http.StatusOK {
			t.Errorf("response code returned %v, want %v", code, http.StatusOK)
		} else {
			t.Logf("response code good: %v", code)
		}

		// Read the response body
		// t.Logf("response body: %v", rr.Body)
		body, err := io.ReadAll(rr.Body)
		if err != nil {
			t.Errorf("failed to read response body: %e", err)
		} else {
			t.Logf("response body read")
		}

		// Parse the response body
		var responseExercise models.Exercise
		if err := json.Unmarshal(body, &responseExercise); err != nil {
			t.Errorf("failed to unmarshal response body: %e", err)
		} else {
			t.Logf("response body unmarshalled")
		}

		// Check the content of the exercise
		switch {
		case responseExercise.Name != testExercise.Name:
			t.Errorf("'Name' bad match: got %v, want %v", responseExercise.Name, testExercise.Name)
		case responseExercise.MuscleGroup != testExercise.MuscleGroup:
			t.Errorf("'MuslceGroup' bad match: got %v, want %v", responseExercise.MuscleGroup, testExercise.MuscleGroup)
		default:
			t.Logf("exercise matches test data")
		}
	}
}
