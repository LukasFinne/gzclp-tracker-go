package internal

import (
	"chi-simple-server/internal/user"
	"context"
	"encoding/json"
	"net/http"
)

type userService interface {
	FetchWorkoutInformation(ctx context.Context, userID string) (user.Workout, error)
	Upload(ctx context.Context, userID string) error
}

type handler struct {
	userService userService
}

func NewHandler(userService userService) *handler {
	return &handler{userService: userService}
}

// fetchWorkout is the handler for /workout
func (h *handler) FetchWorkout(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		userID = "test-id"
	}

	info, err := h.userService.FetchWorkoutInformation(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
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
