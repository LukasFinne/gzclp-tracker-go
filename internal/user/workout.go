package user

import (
	"gzclp-tracker/internal/types"
)

// exerciser describes how to find and update all exercises
type exerciser interface {
	// Info returns the derived data from a workout and includes an exercises: Name, Weight and SetxRep
	Info() map[types.ExerciseID]types.Info

	// Update updates weight, setxrep and changes the current workout based on the result for each exercise in a workout
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
