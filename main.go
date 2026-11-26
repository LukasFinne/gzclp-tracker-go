package main

import (
	"context"
	handler "gzclp-tracker/internal/http/rest"
	"gzclp-tracker/internal/storage/firestore"
	"gzclp-tracker/internal/user"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// main intializes everything required for the applicaiton server to start
func main() {
	ctx := context.Background()

	serviceAccountFilePath, ok := os.LookupEnv("SERVICE_ACCOUNT_PATH")
	if !ok {
		log.Fatalf("error SERVICE_ACCOUNT_PATH env is not set")
	}

	app, err := firestore.App(ctx, serviceAccountFilePath)
	if err != nil {
		log.Fatalf("error initializing firebase app: %v", err)
	}

	auth, err := firestore.Auth(ctx, app)
	if err != nil {
		log.Fatalf("error initializing firebase auth: %v", err)
	}

	db, err := firestore.DB(ctx, app)
	if err != nil {
		log.Fatalf("error initializing firestore: %v", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("error closing firestore client: %v", err)
		}
	}()

	userRepo := firestore.NewUserRepo(db)
	userService := user.NewService(userRepo)
	h := handler.New(userService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(handler.AuthMiddleware(auth))

	r.Get("/workout", h.StartWorkout)
	r.Post("/upload", h.FinishWorkout)

	log.Println("Server running on :3000")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
