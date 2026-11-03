package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Workout is used to handle a user's workout-related data from fetch to updating
type UserProvider interface {
	// findByID fetches all the information for a user by ID
	findByID(ctx context.Context, userID string) (*WorkoutInfo, error)
}

// WorkoutProvider implements the Workout API
type WorkoutProvider struct {
	firestoreClient *firestore.Client // firestoreClient is the instance of the firestore database
}

// NewWorkoutProviders creates new instance of WorkoutProvider
func newWorkoutProvider(firestoreClient *firestore.Client) *WorkoutProvider {
	return &WorkoutProvider{
		firestoreClient: firestoreClient,
	}
}

// FetchInformation gets all workout information for a user by userId from the database
//
// # Parameters
//
//	ctx context.Context - handles request cancellations and timeouts
//
// # Return values
//
//	WorkoutInfo - Contains the users workout information
//	error - returns either a Notfound error, Failed to decode error or General failed to fetch error
func (wp WorkoutProvider) findByID(ctx context.Context, userID string) (*WorkoutInfo, error) {
	docSnap, err := wp.firestoreClient.Collection("users").Doc(userID).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, fmt.Errorf("workout not found for user: %s", userID)
		}
		return nil, fmt.Errorf("failed to fetch workout: %w", err)
	}

	var workoutInfo WorkoutInfo
	if err := docSnap.DataTo(&workoutInfo); err != nil {
		return nil, fmt.Errorf("failed to decode document: %w", err)
	}
	return &workoutInfo, nil
}

// Upload updates exerices values depending on their success or failure
//
// # Parameters
//
// ctx context.Context - handles request cancellations and timeouts
//
// # Return values
//
//	error - returns an error if the updates has failed
func (wp WorkoutProvider) Upload(ctx context.Context, userID string) error {
	docRef := wp.firestoreClient.Collection("users").Doc(userID)
	err := wp.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		var data WorkoutInfo
		return tx.Set(docRef, data, firestore.MergeAll)
	})
	if err != nil {
		return fmt.Errorf("Failed to upload new workout data: %w", err)
	}
	return err
}
