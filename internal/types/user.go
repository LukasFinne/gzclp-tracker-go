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

// ExerciseTier represents the different progression schemes for an exercise
type ExerciseTier string

const (
	T1 ExerciseTier = "T1"
	T2 ExerciseTier = "T2"
	T3 ExerciseTier = "T3"
)

// User represents the workout information for a user
//
// # Fields
//
// 1. CurrentWorkout is the current workout program for the user and is also used as the key for the Program field to find the correspoding workout program.
//
// 2. WorkoutRotation is the set rotation of workout programs.
//
// 3. Program has the program for each workout
//
// 4. Exercises map key is the name of an exercise and has the repition for each exercise
type User struct {
	CurrentWorkout  Workout             `firestore:"currentWorkout" json:"currentWorkout"`
	WorkoutRotation []Workout           `firestore:"workoutRotation" json:"workoutRotation"`
	Program         map[Workout]Program `firestore:"program" json:"program"`
	Exercises       map[string]Repition `firestore:"exercises" json:"exercises"`
}

type Exercise struct {
	Name string `firestore:"exercise" json:"exercise"`
}

// Program defines the workout
//
// Each Program has 3 tiers with one corresponding exercise:
// T1 - Squat
// T2 - Bench
// T3 - Lat pulldown
type Program map[ExerciseTier]Exercise

// Repition has the repition scheme progression depending on the exercises tier (T1, T2, T3)
//
// *Progression - a pointer value is used as every tier is not present on each exercise (e.g., T3 exercises lack T1 and T2).
type Repition map[ExerciseTier]*Progression

// Progression defines the progression scheme of an exercise
//
// # StageIndex is used to get the current stage of an exercise from Stages
type Progression struct {
	StageIndex int     `firestore:"stageIndex" json:"stageIndex"`
	Stages     []Stage `firestore:"stages" json:"stages"`
}

// Stage defines the set, rep and weight for each stage of an progression
type Stage struct {
	Sets   int     `firestore:"sets" json:"sets"`
	Reps   int     `firestore:"reps" json:"reps"`
	Weight float64 `firestore:"weight" json:"weight"`
}
