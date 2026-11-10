package internal

import (
	"context"
	"encoding/json"
	"errors"
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
	FetchWorkoutInformation(ctx context.Context, userID string) (user.Workout, error)
	Upload(ctx context.Context, userID string) error
}

type handler struct {
	userService userService
}

func NewHandler(userService userService) *handler {
	return &handler{userService: userService}
}

// FetchWorkout is the handler for /workout
func (h *handler) FetchWorkout(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")

	workout, err := h.userService.FetchWorkoutInformation(r.Context(), userID)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "failed to fetch workout for user", http.StatusInternalServerError)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workout)
}

// Upload is the handler for /upload - it updates the exercise after a workout
func (h *handler) Upload(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "missing user_id", http.StatusBadRequest)
	}

	err := h.userService.Upload(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Upload success"})
}
