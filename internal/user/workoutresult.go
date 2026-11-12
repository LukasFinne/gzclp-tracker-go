package user

import "gzclp-tracker/internal/types"

// WorkoutResult is the result data that is sent by the client after a workout
//
// # Fields
//
//	Workout - name of the workout
//	Exercises - Contains all of the failed or succeeded exercises by tier
type WorkoutResult struct {
	Workout   types.Workout                         `firestore:"workout" json:"workout"`
	Exercises map[types.ExerciseTier]ExerciseResult `firestore:"exercises" json:"exercises"`
}

// Exercise contains the name and result of the exercise
type ExerciseResult struct {
	Name      string `firestore:"name" json:"name"`
	IsSuccess bool   `firestore:"isSuccess" json:"isSuccess"`
}
