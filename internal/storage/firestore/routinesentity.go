package firestore

import (
	"gzclp-tracker/internal/types"
	"gzclp-tracker/internal/user"
)

// RoutinesEntity represents the data that is stored in the database
type RoutinesEntity struct {
	Current types.Workout             `firestore:"current" json:"current"`
	Workout map[types.Workout]Workout `firestore:"workouts" json:"workouts"`
}

type Workout struct {
	Next types.Workout  `firestore:"nextworkout" json:"nextworkout"`
	T1   *ExerciseState `firestore:"T1" json:"T1"`
	T2   *ExerciseState `firestore:"T2" json:"T2"`
	T3   *ExerciseState `firestore:"T3" json:"T3"`
}

// ExerciseState contains the progression, current stage and weight configuration for an exercise.
type ExerciseState struct {
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

func ToRoutinesEntity(routines user.Routines) *RoutinesEntity {
	workoutMap := make(map[types.Workout]Workout, len(routines.Workouts))

	for workoutID, workout := range routines.Workouts {
		exercise := Workout{
			Next: workout.NextWorkout,
		}

		if d, ok := workout.Exercises[types.T1]; ok && d != nil {
			exercise.T1 = toExerciseStateEntity(d)
		}
		if d, ok := workout.Exercises[types.T2]; ok && d != nil {
			exercise.T2 = toExerciseStateEntity(d)
		}
		if d, ok := workout.Exercises[types.T3]; ok && d != nil {
			exercise.T3 = toExerciseStateEntity(d)
		}

		workoutMap[workoutID] = exercise
	}

	return &RoutinesEntity{
		Current: routines.Current,
		Workout: workoutMap,
	}
}

func toExerciseStateEntity(d *user.ExerciseState) *ExerciseState {
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

	return &ExerciseState{
		Name:                 d.Name,
		Repition:             d.Repition,
		RepitionsProgression: progressions,
		Weight:               d.Weight,
		Increase:             d.Increase,
		Reset:                d.Reset,
	}
}

func (p RoutinesEntity) ToRoutines() *user.Routines {
	workouts := make(map[types.Workout]user.Workout, len(p.Workout))

	for workoutID, storageWorkout := range p.Workout {
		exercises := make(map[types.ExerciseID]*user.ExerciseState, 3)

		if storageWorkout.T1 != nil {
			exercises[types.T1] = toUserExerciseState(storageWorkout.T1)
		}
		if storageWorkout.T2 != nil {
			exercises[types.T2] = toUserExerciseState(storageWorkout.T2)
		}
		if storageWorkout.T3 != nil {
			exercises[types.T3] = toUserExerciseState(storageWorkout.T3)
		}

		workouts[workoutID] = user.Workout{
			NextWorkout: storageWorkout.Next,
			Exercises:   exercises,
		}
	}

	return &user.Routines{
		Current:  p.Current,
		Workouts: workouts,
	}
}

func toUserExerciseState(d *ExerciseState) *user.ExerciseState {
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
