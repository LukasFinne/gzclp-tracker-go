package types

// Workout defines each workout name
type Workout string

const (
	A1 Workout = "A1"
	A2 Workout = "A2"
	B1 Workout = "B1"
	B2 Workout = "B2"
)

// ExerciseID represents the ID for an exercise in a workout
type ExerciseID string

const (
	T1 ExerciseID = "T1" // The primary exercise
	T2 ExerciseID = "T2" // The secondary exercise
	T3 ExerciseID = "T3" // The tertiary exercise
)
