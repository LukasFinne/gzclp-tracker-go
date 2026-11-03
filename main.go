package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// main intializes everything required for the applicaiton server to start
func main() {
	ctx := context.Background()
	serviceAccountFilePath := os.Getenv("SERVICE_ACCOUNT_PATH")

	db, err := InitFirestore(ctx, serviceAccountFilePath)
	if err != nil {
		log.Fatalf("Failed to initialize Firestore: %v", err)
	}
	defer db.Close()

	h := NewHandler(db)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/workout", h.FetchWorkout)
	r.Post("/upload", h.Upload)

	log.Println("Server running on :3000")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
