package user

import (
	"gzclp-tracker/internal/types"
)

// Program has information of each workout and the workouts each exercise
type Program struct {
	Current  types.Workout                    // The current workout
	Workouts map[types.Workout]WorkoutSession // List of all workouts
}

type WorkoutSession struct {
	NextWorkout types.Workout
	Exercises   map[types.ExerciseID]*Details
}

type Details struct {
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

// Result represents a finished workout
type Result struct {
	Workout   types.Workout             // id of the workout
	Exercises map[types.ExerciseID]bool // exercises has the data about which exercise is failed or successed
}

func (p *Program) Info() map[types.ExerciseID]Info {
	workout := p.Workouts[p.Current]

	var allExercises = make(map[types.ExerciseID]Info, len(workout.Exercises))

	for k, v := range workout.Exercises {
		repition := v.RepitionsProgression[v.Repition]
		allExercises[k] = Info{
			Name:   v.Name,
			Set:    repition.Set,
			Rep:    repition.Rep,
			Weight: v.Weight,
		}
	}

	return allExercises
}

func (p *Program) Update(result Result) {
	workout := p.Workouts[result.Workout]

	for ID, result := range result.Exercises {
		details := workout.Exercises[ID]
		if details == nil {
			continue
		}
		if result {
			details.Weight = details.Weight + details.Increase
		} else {

			repition := details.Repition
			newRepition := (details.Repition + 1) % len(details.RepitionsProgression)

			// Checks if the repition has been restarted
			if repition == len(details.RepitionsProgression)-1 && repition > newRepition {

				if ID == types.T1 {
					details.Weight = details.Weight * details.Reset // reduces new starting weight
				}

				if ID == types.T2 {
					details.Weight = details.Reset + details.Increase // starting weight is slightly higher than previous time on the first step.
					details.Reset = details.Weight
				}
			}

			details.Repition = newRepition
		}
	}

	// Rotates the workout program
	p.Current = workout.NextWorkout
}
