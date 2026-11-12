package testdata

import (
	_ "embed"
	"encoding/json"
	"testing"

	"gzclp-tracker/internal/types"
)

//go:embed defaultuser.json
var defaultUserJSON []byte

//go:embed workoutresult.json
var workoutResultJSON []byte

// LoadDefaultUser unmarshals a clean copy of defaultuser.json for testing.
func LoadDefaultUser(t *testing.T) *types.User {
	t.Helper()
	var u types.User
	if err := json.Unmarshal(defaultUserJSON, &u); err != nil {
		t.Fatalf("failed to unmarshal defaultuser.json: %v", err)
	}
	return &u
}

// LoadWorkoutResult unmarshals a clean copy of workoutresult.json for testing.
func LoadWorkoutResult(t *testing.T) *types.User {
	t.Helper()
	var u types.User
	if err := json.Unmarshal(workoutResultJSON, &u); err != nil {
		t.Fatalf("failed to unmarshal workoutresult.json: %v", err)
	}
	return &u
}
