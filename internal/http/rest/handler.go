package rest

import (
	"context"
	"encoding/json"
	"errors"
	"gzclp-tracker/internal/types"
	"gzclp-tracker/internal/user"
	"net/http"
)

type userService interface {
	// StartWorkout fetches the current workout for the user
	//
	// # Return
	//
	//	Workout - a map object that contains data for the workout program(name, weight, setxrep, etc...)
	//
	//	Error
	// 	 1. Returns a ErrNotFound when a document cannot be found
	//   2. Returns a failed error when failing to fetch a document
	StartWorkout(ctx context.Context, userID string) (user.Exercises, error)

	// FinishWorkout updates the database with the new workout result from the client
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
	FinishWorkout(ctx context.Context, userID string, result user.Result) error
}

type handler struct {
	userService userService
}

func New(userService userService) *handler {
	return &handler{userService: userService}
}

// StartWorkout is the handler for /workout
func (h *handler) StartWorkout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := UserFromContext(ctx)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	workout, err := h.userService.StartWorkout(ctx, userID)
	if err != nil {
		if errors.Is(err, types.ErrNotFound) {
			http.Error(w, "start workout not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to start workout for user", http.StatusInternalServerError)
		return

	}

	response := ToResponse(workout)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "failed to encode start workout response", http.StatusInternalServerError)
		return
	}
}

// FinishWorkout is the handler for /upload - it updates the exercise after a workout
func (h *handler) FinishWorkout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := UserFromContext(ctx)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var request ResultRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	result := request.ToResult()

	err := h.userService.FinishWorkout(ctx, userID, result)
	if err != nil {
		if errors.Is(err, types.ErrNotFound) {
			http.Error(w, "finish workout not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to finish workout", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]string{"message": "Upload success"})
	if err != nil {
		http.Error(w, "failed to encode finish workout response", http.StatusInternalServerError)
		return
	}
}
