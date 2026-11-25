package testdata

import (
	_ "embed"
	"encoding/json"
	"testing"

	"gzclp-tracker/internal/types"
)

//go:embed defaultuser.json
var defaultUserJSON []byte

// TestRoutinesEntity is a local JSON representation of the test data
// to avoid import cycles between user and storage/firestore.
type TestRoutinesEntity struct {
	Current  types.Workout                        `json:"current"`
	Workouts map[types.Workout]TestWorkoutSession `json:"workouts"`
}

type TestWorkoutSession struct {
	NextWorkout types.Workout                     `json:"nextworkout"`
	Exercises   map[types.ExerciseID]TestExercise `json:"exercises"`
}

type TestExercise struct {
	Name                 string                    `json:"name"`
	Repition             int                       `json:"repition"`
	RepitionsProgression []TestRepitionProgression `json:"repitionsProgression"`
	Weight               float64                   `json:"weight"`
	Increase             float64                   `json:"increase"`
	Reset                float64                   `json:"reset"`
}

type TestRepitionProgression struct {
	Set int `json:"set"`
	Rep int `json:"rep"`
}

// LoadDefaultRoutinesEntity unmarshals defaultuser.json into a TestRoutinesEntity.
func LoadDefaultRoutinesEntity(t *testing.T) *TestRoutinesEntity {
	t.Helper()
	var entity TestRoutinesEntity
	if err := json.Unmarshal(defaultUserJSON, &entity); err != nil {
		t.Fatalf("failed to unmarshal defaultuser.json: %v", err)
	}
	return &entity
}
