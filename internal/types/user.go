package types

// WorkoutDay represents what the users current day is in the program
type workoutDay string

const (
	A1 workoutDay = "A1"
	A2 workoutDay = "A2"
	B1 workoutDay = "B1"
	B2 workoutDay = "B2"
)

// User represents the workout information for a user
type User struct {
	CurrentWorkoutDay workoutDay                  `firestore:"currentWorkoutDay" json:"currentWorkoutDay"` // CurrentWorkoutDay holds the users current day (A1, A2, B1, B2)
	DayOrder          []workoutDay                `firestore:"dayOrder" json:"dayOrder"`                   // DayOrder is the current set rotation of WorkoutDays
	Program           map[workoutDay]programTiers `firestore:"days" json:"days"`                           // Program map has the exercise program for each day
	Exercises         map[string]exerciseTiers    `firestore:"exercises" json:"exercises"`                 // Exercises map has the weight, set&rep values for each exercise
}

// ProgramTiers defines the exercise assigned to each tier for a specific day.
type programTiers struct {
	T1 string `firestore:"T1,omitempty" json:"T1,omitempty"`
	T2 string `firestore:"T2,omitempty" json:"T2,omitempty"`
	T3 string `firestore:"T3,omitempty" json:"T3,omitempty"`
}

// ExerciseTiers maps tier progression configurations for an exercise.
// Pointers with omitempty handle tiers not present on all exercises (e.g., T3 exercises lack T1 and T2).
type exerciseTiers struct {
	T1 *tierConfig `firestore:"T1,omitempty" json:"T1,omitempty"`
	T2 *tierConfig `firestore:"T2,omitempty" json:"T2,omitempty"`
	T3 *tierConfig `firestore:"T3,omitempty" json:"T3,omitempty"`
}

// tierConfig defines the weight, current stage index, and progression scheme for a tier.
// StageIndex is used to get the current stage from []Stages
type tierConfig struct {
	Weight     float64 `firestore:"weight" json:"weight"`
	StageIndex int     `firestore:"stageIndex" json:"stageIndex"` // Represent the current stage of the exercise
	Stages     []stage `firestore:"stages" json:"stages"`         // Array of each set & rep for each stage(1,2,3)
}

// stage defines set and repetition counts.
type stage struct {
	Sets int `firestore:"sets" json:"sets"`
	Reps int `firestore:"reps" json:"reps"`
}

// Flatten exerice dervied from WorkoutInfo
type flattenExercise struct {
	Name   string  `firestore:"name"`
	Tier   string  `firestore:"tier"`
	Set    int     `firestore:"set"`
	Rep    int     `firestore:"rep"`
	Weight float64 `firestore:"weight"`
}
