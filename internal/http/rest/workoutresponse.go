package handler

import (
	"gzclp-tracker/internal/types"
	"gzclp-tracker/internal/user"
)

// WorkoutResponse type contains all the derived fields from the database
type WorkoutResponse map[types.ExerciseID]exercise

// exercise has the weight and setxrep values
type exercise struct {
	Name   string  `firestore:"name"`
	Set    int     `firestore:"set"`
	Rep    int     `firestore:"rep"`
	Weight float64 `firestore:"weight"`
}

// Converts programInfo to the dto Workout
func ToResponse(programInfo map[types.ExerciseID]user.Info) WorkoutResponse {

	workout := make(WorkoutResponse, len(programInfo))

	for k, v := range programInfo {
		workout[k] = exercise{
			Name:   v.Name,
			Set:    v.Set,
			Rep:    v.Rep,
			Weight: v.Weight,
		}
	}

	return workout
}
