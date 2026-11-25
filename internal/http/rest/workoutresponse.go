package handler

import (
	"gzclp-tracker/internal/types"
	"gzclp-tracker/internal/user"
)

// WorkoutResponse contains the workout ID and exercises
type WorkoutResponse struct {
	Workout   types.Workout                 `json:"workout"`
	Exercises map[types.ExerciseID]Exercise `json:"exercises"`
}

type Exercise struct {
	Name   string  `json:"name"`
	Set    int     `json:"set"`
	Rep    int     `json:"rep"`
	Weight float64 `json:"weight"`
}

func ToResponse(exercises user.Exercises) WorkoutResponse {
	workout := make(map[types.ExerciseID]Exercise, len(exercises.Exercises))

	for k, v := range exercises.Exercises {
		workout[k] = Exercise{
			Name:   v.Name,
			Set:    v.Set,
			Rep:    v.Rep,
			Weight: v.Weight,
		}
	}

	return WorkoutResponse{
		Workout:   exercises.WorkoutID,
		Exercises: workout,
	}
}
