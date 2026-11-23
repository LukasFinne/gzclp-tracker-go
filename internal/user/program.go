package user

import (
	"fmt"
	"gzclp-tracker/internal/types"
)

// Program has information of each workout and the workouts each exercise
type Program struct {
	Current  types.Workout                    // The current workout
	Workouts map[types.Workout]WorkoutSession // List of all workouts
}

type WorkoutSession struct {
	NextWorkout types.Workout
	Exercises   map[types.ExerciseID]*ExerciseState
}

type ExerciseState struct {
	Name                 string
	Repition             int
	RepitionsProgression []Repitions
	Weight               float64
	Increase             float64 // How much the weight should be increased
	Reset                float64 // The value that decides how much the weight should be resetted to when a new progession begins
}

type Repitions struct {
	Set int
	Rep int
}

func (p *Program) Start() Exercises {
	workout := p.Workouts[p.Current]

	var allExercises = make(map[types.ExerciseID]Exercise, len(workout.Exercises))

	for k, v := range workout.Exercises {
		repition := v.RepitionsProgression[v.Repition]
		allExercises[k] = Exercise{
			Name:   v.Name,
			Set:    repition.Set,
			Rep:    repition.Rep,
			Weight: v.Weight,
		}
	}

	return Exercises{
		WorkoutID: p.Current,
		Exercises: allExercises,
	}
}

func (p *Program) Finish(result Result) (Program, error) {
	workout, ok := p.Workouts[result.WorkoutID]
	if !ok {
		return Program{}, fmt.Errorf("workoutID %s not found", result.WorkoutID)
	}

	for ID, success := range result.Exercises {
		exercise, ok := workout.Exercises[ID]
		if !ok || exercise == nil {
			return Program{}, fmt.Errorf("exercise %s not found in workout %s", ID, result.WorkoutID)
		}

		switch {

		// Restart protocol and weight progression
		case exercise.Repition == len(exercise.RepitionsProgression)-1 && !success:

			if ID == types.T1 {
				exercise.Weight = exercise.Weight * exercise.Reset
			}

			if ID == types.T2 {
				exercise.Weight = exercise.Reset + exercise.Increase
				exercise.Reset = exercise.Weight
			}

			exercise.Repition = 0 // resets the progression to the beginning of the protocol

			// Updates to next protocol on failure
		case !success:
			newRepition := (exercise.Repition + 1) % len(exercise.RepitionsProgression)
			exercise.Repition = newRepition

			// Increases weight on success
		case success:
			exercise.Weight = exercise.Weight + exercise.Increase

		default:
			return Program{}, fmt.Errorf("exercise %s is in an invalid state: repition %d", ID, exercise.Repition)
		}
	}

	// Rotates the workout program
	p.Current = workout.NextWorkout

	return *p, nil
}
