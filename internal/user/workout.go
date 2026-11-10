package user

import "gzclp-tracker/internal/types"

// Workout type contains all the derived fields from the database
type Workout map[types.ExerciseTier]exercise

// exercise has the weight and setxrep values
type exercise struct {
	Name   string  `firestore:"name"`
	Set    int     `firestore:"set"`
	Rep    int     `firestore:"rep"`
	Weight float64 `firestore:"weight"`
}
