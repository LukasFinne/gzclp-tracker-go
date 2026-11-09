package testdata

import (
	_ "embed"
	"encoding/json"
	"testing"

	"chi-simple-server/internal/types"
)

//go:embed defaultuser.json
var defaultUserJSON []byte

// LoadDefaultUser unmarshals a clean copy of defaultuser.json for testing.
func LoadDefaultUser(t *testing.T) *types.User {
	t.Helper()
	var u types.User
	if err := json.Unmarshal(defaultUserJSON, &u); err != nil {
		t.Fatalf("failed to unmarshal defaultuser.json: %v", err)
	}
	return &u
}
