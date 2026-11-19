package user

import (
	"gzclp-tracker/internal/types"
)

// Program has information of each workout and the workouts each exercise
type Program struct {
	Current types.Workout              // The current workout
	Workout map[types.Workout]Exercise // List of all workouts
}

type Exercise map[types.ExerciseID]Details

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

func (p *Program) Info() map[types.ExerciseID]types.Info {
	exercises := p.Workout[p.Current]

	var allExercises = make(map[types.ExerciseID]types.Info, len(exercises))

	for k, v := range exercises {
		repition := v.RepitionsProgression[v.Repition]
		allExercises[k] = types.Info{
			Name:   v.Name,
			Set:    repition.Set,
			Rep:    repition.Rep,
			Weight: v.Weight,
		}
	}

	return allExercises
}

func (p *Program) Update(result Result) {
	exercises := p.Workout[result.Workout]

	for ID, result := range result.Exercises {
		details := exercises[ID]
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
}
