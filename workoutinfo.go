package main

// WorkoutDay represents what the users current day is in the program
type WorkoutDay string

const (
	A1 WorkoutDay = "A1"
	A2 WorkoutDay = "A2"
	B1 WorkoutDay = "B1"
	B2 WorkoutDay = "B2"
)

// WorkoutInfo represents the workout information for a user
type WorkoutInfo struct {
	CurrentWorkoutDay WorkoutDay                  `firestore:"currentWorkoutDay" json:"currentWorkoutDay"` // CurrentWorkoutDay holds the users current day (A1, A2, B1, B2)
	DayOrder          []WorkoutDay                `firestore:"dayOrder" json:"dayOrder"`                   // DayOrder is the current set rotation of WorkoutDays
	Program           map[WorkoutDay]ProgramTiers `firestore:"days" json:"days"`                           // Program map has the exercise program for each day
	Exercises         map[string]ExerciseTiers    `firestore:"exercises" json:"exercises"`                 // Exercises map has the weight, set&rep values for each exercise
}

// ProgramTiers defines the exercise assigned to each tier for a specific day.
type ProgramTiers struct {
	T1 string `firestore:"T1,omitempty" json:"T1,omitempty"`
	T2 string `firestore:"T2,omitempty" json:"T2,omitempty"`
	T3 string `firestore:"T3,omitempty" json:"T3,omitempty"`
}

// ExerciseTiers maps tier progression configurations for an exercise.
// Pointers with omitempty handle tiers not present on all exercises (e.g., T3 exercises lack T1 and T2).
type ExerciseTiers struct {
	T1 *TierConfig `firestore:"T1,omitempty" json:"T1,omitempty"`
	T2 *TierConfig `firestore:"T2,omitempty" json:"T2,omitempty"`
	T3 *TierConfig `firestore:"T3,omitempty" json:"T3,omitempty"`
}

// TierConfig defines the weight, current stage index, and progression scheme for a tier.
// StageIndex is used to get the current stage from []Stages
type TierConfig struct {
	Weight     float64 `firestore:"weight" json:"weight"`
	StageIndex int     `firestore:"stageIndex" json:"stageIndex"` // Represent the current stage of the exercise
	Stages     []Stage `firestore:"stages" json:"stages"`         // Array of each set & rep for each stage(1,2,3)
}

// Stage defines set and repetition counts.
type Stage struct {
	Sets int `firestore:"sets" json:"sets"`
	Reps int `firestore:"reps" json:"reps"`
}

// Flatten exerice dervied from WorkoutInfo
type FlattenExercise struct {
	Name   string  `firestore:"name"`
	Tier   string  `firestore:"tier"`
	Set    int     `firestore:"set"`
	Rep    int     `firestore:"rep"`
	Weight float64 `firestore:"weight"`
}
