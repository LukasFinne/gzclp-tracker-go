package firestore_test

import (
	"encoding/json"
	"gzclp-tracker/internal/storage/firestore"
	"gzclp-tracker/internal/types"
	"reflect"
	"testing"
)

const testJSON = `{
  "current": "A1",
  "workouts": {
    "A1": {
      "nextworkout": "B1",
      "T1": {
        "name": "Squat",
        "repition": 0,
        "repitionsProgression": [
          { "set": 5, "rep": 3 },
          { "set": 6, "rep": 2 },
          { "set": 10, "rep": 1 }
        ],
        "weight": 100.0,
        "increase": 5.0,
        "reset": 0.85
      },
      "T2": {
        "name": "Bench Press",
        "repition": 0,
        "repitionsProgression": [
          { "set": 3, "rep": 10 },
          { "set": 3, "rep": 8 },
          { "set": 3, "rep": 6 }
        ],
        "weight": 60.0,
        "increase": 2.5,
        "reset": 60.0
      },
      "T3": {
        "name": "Lat Pulldown",
        "repition": 0,
        "repitionsProgression": [
          { "set": 3, "rep": 15 }
        ],
        "weight": 45.0,
        "increase": 2.5,
        "reset": 0.0
      }
    },
    "B1": {
      "nextworkout": "A2",
      "T1": {
        "name": "Overhead Press",
        "repition": 0,
        "repitionsProgression": [
          { "set": 5, "rep": 3 },
          { "set": 6, "rep": 2 },
          { "set": 10, "rep": 1 }
        ],
        "weight": 45.0,
        "increase": 2.5,
        "reset": 0.85
      },
      "T2": {
        "name": "Deadlift",
        "repition": 0,
        "repitionsProgression": [
          { "set": 3, "rep": 10 },
          { "set": 3, "rep": 8 },
          { "set": 3, "rep": 6 }
        ],
        "weight": 110.0,
        "increase": 5.0,
        "reset": 110.0
      },
      "T3": {
        "name": "Dumbbell Row",
        "repition": 0,
        "repitionsProgression": [
          { "set": 3, "rep": 15 }
        ],
        "weight": 25.0,
        "increase": 2.5,
        "reset": 0.0
      }
    },
    "A2": {
      "nextworkout": "B2",
      "T1": {
        "name": "Bench Press",
        "repition": 0,
        "repitionsProgression": [
          { "set": 5, "rep": 3 },
          { "set": 6, "rep": 2 },
          { "set": 10, "rep": 1 }
        ],
        "weight": 75.0,
        "increase": 2.5,
        "reset": 0.85
      },
      "T2": {
        "name": "Squat",
        "repition": 0,
        "repitionsProgression": [
          { "set": 3, "rep": 10 },
          { "set": 3, "rep": 8 },
          { "set": 3, "rep": 6 }
        ],
        "weight": 80.0,
        "increase": 5.0,
        "reset": 80.0
      },
      "T3": {
        "name": "Lat Pulldown",
        "repition": 0,
        "repitionsProgression": [
          { "set": 3, "rep": 15 }
        ],
        "weight": 45.0,
        "increase": 2.5,
        "reset": 0.0
      }
    },
    "B2": {
      "nextworkout": "A1",
      "T1": {
        "name": "Deadlift",
        "repition": 0,
        "repitionsProgression": [
          { "set": 5, "rep": 3 },
          { "set": 6, "rep": 2 },
          { "set": 10, "rep": 1 }
        ],
        "weight": 130.0,
        "increase": 5.0,
        "reset": 0.85
      },
      "T2": {
        "name": "Overhead Press",
        "repition": 0,
        "repitionsProgression": [
          { "set": 3, "rep": 10 },
          { "set": 3, "rep": 8 },
          { "set": 3, "rep": 6 }
        ],
        "weight": 35.0,
        "increase": 2.5,
        "reset": 35.0
      },
      "T3": {
        "name": "Dumbbell Row",
        "repition": 0,
        "repitionsProgression": [
          { "set": 3, "rep": 15 }
        ],
        "weight": 25.0,
        "increase": 2.5,
        "reset": 0.0
      }
    }
  }
}`

func TestProgramEntityUnmarshal(t *testing.T) {
	var entity firestore.ProgramEntity
	if err := json.Unmarshal([]byte(testJSON), &entity); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if entity.Current != types.A1 {
		t.Fatalf("expected Current to be A1, got %v", entity.Current)
	}

	if len(entity.Workout) != 4 {
		t.Fatalf("expected 4 workouts, got %d", len(entity.Workout))
	}

	a1, ok := entity.Workout[types.A1]
	if !ok {
		t.Fatalf("expected workout A1 to exist")
	}

	if a1.Next != types.B1 {
		t.Errorf("expected A1 next to be B1, got %v", a1.Next)
	}

	if a1.T1 == nil || a1.T1.Name != "Squat" || a1.T1.Weight != 100.0 {
		t.Errorf("unexpected A1 T1: %+v", a1.T1)
	}
	if len(a1.T1.RepitionsProgression) != 3 {
		t.Errorf("expected 3 progression steps for A1 T1, got %d", len(a1.T1.RepitionsProgression))
	}

	if a1.T2 == nil || a1.T2.Name != "Bench Press" || a1.T2.Weight != 60.0 {
		t.Errorf("unexpected A1 T2: %+v", a1.T2)
	}

	if a1.T3 == nil || a1.T3.Name != "Lat Pulldown" || a1.T3.Weight != 45.0 {
		t.Errorf("unexpected A1 T3: %+v", a1.T3)
	}

	// Test ToProgram
	prog := entity.ToProgram()
	if prog.Current != types.A1 {
		t.Errorf("expected prog.Current to be A1, got %v", prog.Current)
	}
	if prog.Workouts[types.A1].NextWorkout != types.B1 {
		t.Errorf("expected prog A1 next to be B1, got %v", prog.Workouts[types.A1].NextWorkout)
	}
	if prog.Workouts[types.A1].Exercises[types.T1].Name != "Squat" {
		t.Errorf("expected Squat, got %v", prog.Workouts[types.A1].Exercises[types.T1].Name)
	}

	// Test ToProgramEntity round trip
	roundTripEntity := firestore.ToProgramEntity(*prog)
	if !reflect.DeepEqual(&entity, roundTripEntity) {
		t.Errorf("round trip entity does not match original entity")
	}
}
