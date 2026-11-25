package rest

import (
	"gzclp-tracker/internal/types"
	"gzclp-tracker/internal/user"
)

// ResultRequest is the result data that is sent by the client after a workout
//
// # Fields
//
//	Workout - name of the workout
//	Exercises - Contains all of the failed or succeeded exercises
type ResultRequest struct {
	Workout   types.Workout             `firestore:"workout" json:"workout"`
	Exercises map[types.ExerciseID]bool `firestore:"exercises" json:"exercises"`
}

func (r ResultRequest) ToResult() user.Result {
	return user.Result{
		WorkoutID: r.Workout,
		Exercises: r.Exercises,
	}
}
