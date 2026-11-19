package testdata

import (
	_ "embed"
	"encoding/json"
	"testing"

	"gzclp-tracker/internal/types"
)

//go:embed defaultuser.json
var defaultUserJSON []byte

// TestProgramEntity is a local JSON representation of the test data
// to avoid import cycles between user and storage/firestore.
type TestProgramEntity struct {
	Current types.Workout               `json:"current"`
	Workout map[types.Workout]Exercises `json:"workout"`
}

type Exercises map[types.ExerciseID]TestDetails

type TestDetails struct {
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

// LoadDefaultProgramEntity unmarshals defaultuser.json into a TestProgramEntity.
func LoadDefaultProgramEntity(t *testing.T) *TestProgramEntity {
	t.Helper()
	var entity TestProgramEntity
	if err := json.Unmarshal(defaultUserJSON, &entity); err != nil {
		t.Fatalf("failed to unmarshal defaultuser.json: %v", err)
	}
	return &entity
}
