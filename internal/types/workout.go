package types

// Workout defines a workout program
//
// Each workout has 3 exercises in their respective program and each have this following this protocol:
//
// First exericse is a tier 1 exercise, Second exercise is a tier 2 exercise and the third exercise is a tier 3 exercise
type Workout string

const (
	A1 Workout = "A1"
	A2 Workout = "A2"
	B1 Workout = "B1"
	B2 Workout = "B2"
)

// ExerciseID represents the ID for an exercise in a program
type ExerciseID string

const (
	T1 ExerciseID = "T1" // The primary exercise
	T2 ExerciseID = "T2" // The secondary exercise
	T3 ExerciseID = "T3" // The tertiary exercise
)
