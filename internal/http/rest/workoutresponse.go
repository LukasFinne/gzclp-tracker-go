package handler

import (
	"gzclp-tracker/internal/types"
	"gzclp-tracker/internal/user"
)

// WorkoutResponse type contains all the derived fields from the training program
type WorkoutResponse map[types.ExerciseID]Exercise

type Exercise struct {
	Name   string  `json:"name"`
	Set    int     `json:"set"`
	Rep    int     `json:"rep"`
	Weight float64 `json:"weight"`
}

// Converts programInfo to the dto Workout
func ToResponse(programInfo map[types.ExerciseID]user.Info) WorkoutResponse {
	workout := make(WorkoutResponse, len(programInfo))

	for k, v := range programInfo {
		workout[k] = Exercise{
			Name:   v.Name,
			Set:    v.Set,
			Rep:    v.Rep,
			Weight: v.Weight,
		}
	}

	return workout
}
