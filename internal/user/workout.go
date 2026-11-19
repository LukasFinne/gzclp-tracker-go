package user

import (
	"gzclp-tracker/internal/types"
)

// exerciser describes how to find and update workout programs
type exerciser interface {
	// Info returns the derived data from a workout and includes an exercises: Name, Weight and SetxRep
	//
	// # Return
	//  map[types.ExerciseID]Info is a concise object that includes all exercises from a specific workout
	Info() map[types.ExerciseID]Info

	// Update updates the workout program exercises(Weight, Set & Rep) based on if it was a success or failure.
	// As well as, changing to the next workout program in the rotation.
	//
	// # Fields
	//  results Result - The result data of a workout sent from the client
	//
	Update(results Result)
}

func NewWorkout(exerciser exerciser) workout {
	return workout{
		Exercise: exerciser,
	}
}

type workout struct {
	Exercise exerciser
}

type Info struct {
	Name   string
	Set    int
	Rep    int
	Weight float64
}
