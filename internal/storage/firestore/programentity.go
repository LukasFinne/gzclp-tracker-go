package firestore

import (
	"gzclp-tracker/internal/types"
	"gzclp-tracker/internal/user"
)

// ProgramEntity represents the data that is stored in the database
type ProgramEntity struct {
	Current types.Workout               `firestore:"current" json:"current"`
	Workout map[types.Workout]Exercises `firestore:"workout" json:"workout"`
}

type Exercises map[types.ExerciseID]Details

// Details contains the progression, current stage and weight configuration for an exercise.
type Details struct {
	Name                 string                `firestore:"name" json:"name"`
	Repition             int                   `firestore:"repition" json:"repition"`
	RepitionsProgression []RepitionProgression `firestore:"repitionsProgression" json:"repitionsProgression"`
	Weight               float64               `firestore:"weight" json:"weight"`
	Increase             float64               `firestore:"increase" json:"increase"`
	Reset                float64               `firestore:"reset" json:"reset"`
}

type RepitionProgression struct {
	Set int `firestore:"set" json:"set"`
	Rep int `firestore:"rep" json:"rep"`
}

func ToProgramEntity(program user.Program) *ProgramEntity {
	workout := make(map[types.Workout]Exercises, len(program.Workout))

	for workoutID, programExercises := range program.Workout {
		exercises := make(Exercises, len(programExercises))

		for exID, d := range programExercises {
			progressions := make([]RepitionProgression, len(d.RepitionsProgression))
			for i, p := range d.RepitionsProgression {
				progressions[i] = RepitionProgression{
					Set: p.Set,
					Rep: p.Rep,
				}
			}

			exercises[exID] = Details{
				Name:                 d.Name,
				Repition:             d.Repition,
				RepitionsProgression: progressions,
				Weight:               d.Weight,
				Increase:             d.Increase,
				Reset:                d.Reset,
			}
		}

		workout[workoutID] = exercises
	}

	return &ProgramEntity{
		Current: program.Current,
		Workout: workout,
	}
}

func (p ProgramEntity) ToProgram() *user.Program {
	workout := make(map[types.Workout]user.Exercise, len(p.Workout))

	// 2. Iterate through each workout  ID (e.g. A1, B1, A2, B2)
	for workoutID, storageExercises := range p.Workout {
		exercises := make(user.Exercise, len(storageExercises))

		// 3. Iterate through each exercise ID in this workout (e.g. T1, T2, T3)
		for exID, d := range storageExercises {

			progressions := make([]user.Repitions, len(d.RepitionsProgression))
			for i, p := range d.RepitionsProgression {
				progressions[i] = user.Repitions{
					Set: p.Set,
					Rep: p.Rep,
				}
			}

			exercises[exID] = user.Details{
				Name:                 d.Name,
				Repition:             d.Repition,
				RepitionsProgression: progressions,
				Weight:               d.Weight,
				Increase:             d.Increase,
				Reset:                d.Reset,
			}
		}

		workout[workoutID] = exercises
	}

	return &user.Program{
		Current: p.Current,
		Workout: workout,
	}
}
