package firestore

import (
	"gzclp-tracker/internal/types"
	"gzclp-tracker/internal/user"
)

// ProgramEntity represents the data that is stored in the database
type ProgramEntity struct {
	Current types.Workout                    `firestore:"current" json:"current"`
	Workout map[types.Workout]WorkoutSession `firestore:"workouts" json:"workouts"`
}

type WorkoutSession struct {
	Next types.Workout `firestore:"nextworkout" json:"nextworkout"`
	T1   *Details      `firestore:"T1" json:"T1"`
	T2   *Details      `firestore:"T2" json:"T2"`
	T3   *Details      `firestore:"T3" json:"T3"`
}

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
	workoutMap := make(map[types.Workout]WorkoutSession, len(program.Workouts))

	for workoutID, programWorkout := range program.Workouts {
		session := WorkoutSession{
			Next: programWorkout.NextWorkout,
		}

		if d, ok := programWorkout.Exercises[types.T1]; ok && d != nil {
			session.T1 = toEntityDetails(d)
		}
		if d, ok := programWorkout.Exercises[types.T2]; ok && d != nil {
			session.T2 = toEntityDetails(d)
		}
		if d, ok := programWorkout.Exercises[types.T3]; ok && d != nil {
			session.T3 = toEntityDetails(d)
		}

		workoutMap[workoutID] = session
	}

	return &ProgramEntity{
		Current: program.Current,
		Workout: workoutMap,
	}
}

func toEntityDetails(d *user.ExerciseState) *Details {
	if d == nil {
		return nil
	}
	progressions := make([]RepitionProgression, len(d.RepitionsProgression))
	for i, p := range d.RepitionsProgression {
		progressions[i] = RepitionProgression{
			Set: p.Set,
			Rep: p.Rep,
		}
	}

	return &Details{
		Name:                 d.Name,
		Repition:             d.Repition,
		RepitionsProgression: progressions,
		Weight:               d.Weight,
		Increase:             d.Increase,
		Reset:                d.Reset,
	}
}

func (p ProgramEntity) ToProgram() *user.Program {
	workouts := make(map[types.Workout]user.WorkoutSession, len(p.Workout))

	for workoutID, storageWorkout := range p.Workout {
		exercises := make(map[types.ExerciseID]*user.ExerciseState, 3)

		if storageWorkout.T1 != nil {
			exercises[types.T1] = toUserDetails(storageWorkout.T1)
		}
		if storageWorkout.T2 != nil {
			exercises[types.T2] = toUserDetails(storageWorkout.T2)
		}
		if storageWorkout.T3 != nil {
			exercises[types.T3] = toUserDetails(storageWorkout.T3)
		}

		workouts[workoutID] = user.WorkoutSession{
			NextWorkout: storageWorkout.Next,
			Exercises:   exercises,
		}
	}

	return &user.Program{
		Current:  p.Current,
		Workouts: workouts,
	}
}

func toUserDetails(d *Details) *user.ExerciseState {
	if d == nil {
		return nil
	}
	progressions := make([]user.Repitions, len(d.RepitionsProgression))
	for i, p := range d.RepitionsProgression {
		progressions[i] = user.Repitions{
			Set: p.Set,
			Rep: p.Rep,
		}
	}

	return &user.ExerciseState{
		Name:                 d.Name,
		Repition:             d.Repition,
		RepitionsProgression: progressions,
		Weight:               d.Weight,
		Increase:             d.Increase,
		Reset:                d.Reset,
	}
}
