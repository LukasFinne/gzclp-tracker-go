package user

import (
	"gzclp-tracker/internal/types"
)

// exerciser finishes the workout and updates
type exerciser interface {
	// Start finds the current workout for a user
	//
	// # Return
	//  Exercises has the current workout ID and the workout corresponding exercises
	Start() Exercises

	// Finish completes a workout and updates each exercises value
	//
	// # Parameters
	// 	Results - The results of each exercise in an workout also contains the corresponding workoutID
	// # Return
	//
	//  Program - the updated data according to the results
	//  Error
	// 		1. Returns a not found error when workoutID or ExerciseID is missing
	// 		2. Returns a invaild state error when the progression of an exercise cannot be done
	//
	Finish(results Result) (Routines, error)
}

func NewWorkout(exerciser exerciser) workout {
	return workout{
		Exercise: exerciser,
	}
}

type workout struct {
	Exercise exerciser
}

// Contains the exercises for the current workout and the workout ID
type Exercises struct {
	WorkoutID types.Workout
	Exercises map[types.ExerciseID]Exercise
}

type Exercise struct {
	Name   string
	Set    int
	Rep    int
	Weight float64
}

// Result represents a finished workout
type Result struct {
	WorkoutID types.Workout             // id of the workout
	Exercises map[types.ExerciseID]bool // exercises has the data about which exercise is failed or successed
}
