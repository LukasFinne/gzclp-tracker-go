package main

import (
	"encoding/json"
	"net/http"

	"cloud.google.com/go/firestore"
)

type Handler struct {
	db *firestore.Client
}

func NewHandler(db *firestore.Client) *Handler {
	return &Handler{db: db}
}

// FetchWorkout is the handler for /workout - It handles the users workout information
func (h *Handler) FetchWorkout(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		userID = "test-id"
	}

	userProvider := newWorkoutProvider(h.db)
	user := newUser(userID, userProvider)

	info, err := user.fetchWorkoutInformation(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

// Upload is the handler for /upload - it updates the exercise after a workout
func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	userProvider := newWorkoutProvider(h.db)
	user := newUser(userID, userProvider)

	err := user.upload(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Upload successs")
}
