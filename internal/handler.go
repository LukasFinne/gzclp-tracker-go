package internal

import (
	"chi-simple-server/internal/types"
	"context"
	"encoding/json"
	"net/http"
)

type userService interface {
	FetchWorkoutInformation(ctx context.Context, userID string) (types.User, error)
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
		return
	}

	var data types.User
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	//	if err := h.userService.fetchWorkoutInformation().Save(r.Context(), userID, &data); err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	// }

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Upload success"})
}
