package handler

import (
	"context"
	"encoding/json"
	"errors"
	"gzclp-tracker/internal/types"
	"gzclp-tracker/internal/user"
	"net/http"
)

type userService interface {
	// FetchWorkoutInformation fetches the workout program for a user
	//
	// # Return
	//
	//	Workout - a map object that contains data for the workout program(name, weight, setxrep, etc...)
	//
	//	Error - Returns either ErrUserNotFound or Generic failed error
	FetchWorkoutInformation(ctx context.Context, userID string) (map[types.ExerciseID]types.Info, error)

	// Upload will updated the database with the new workout result from the client
	//
	// # Parameters
	//
	// ctx context.Context - handles metadata, request cancellations and timeouts
	//
	// userID string
	//
	// workoutResult - the result of an workout sent by the client
	//
	// # Return
	//
	// error - returns a updated failed error
	Upload(ctx context.Context, userID string, result user.Result) error
}

type handler struct {
	userService userService
}

func New(userService userService) *handler {
	return &handler{userService: userService}
}

// FetchWorkout is the handler for /workout
func (h *handler) FetchWorkout(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")

	workoutInfo, err := h.userService.FetchWorkoutInformation(r.Context(), userID)
	if err != nil {
		if errors.Is(err, types.ErrUserNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "failed to fetch workout for user", http.StatusInternalServerError)
		}
	}

	response := ToResponse(workoutInfo)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Upload is the handler for /upload - it updates the exercise after a workout
func (h *handler) Upload(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "missing user_id", http.StatusBadRequest)
	}

	var request ResultRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
	}

	result := request.ToResult()

	err := h.userService.Upload(r.Context(), userID, result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Upload success"})
}
