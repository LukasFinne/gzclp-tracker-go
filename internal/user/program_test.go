package user

import (
	"math"
	"testing"

	"gzclp-tracker/internal/types"
)

func TestProgram_Update(t *testing.T) {
	defaultStages := []Repitions{
		{Set: 5, Rep: 3},
		{Set: 6, Rep: 2},
		{Set: 10, Rep: 1},
	}

	tests := []struct {
		name             string
		exerciseID       types.ExerciseID
		initialDetail    Details
		exerciseResult   bool
		expectedRepition int
		expectedWeight   float64
		expectedReset    float64
	}{
		{
			name:       "T1 success at stage 0 increases weight and keeps repition stage",
			exerciseID: types.T1,
			initialDetail: Details{
				Name:                 "Squat",
				Repition:             0,
				RepitionsProgression: defaultStages,
				Weight:               100.0,
				Increase:             5.0,
				Reset:                0.85,
			},
			exerciseResult:   true,
			expectedRepition: 0,
			expectedWeight:   105.0,
			expectedReset:    0.85,
		},
		{
			name:       "T1 success at intermediate stage increases weight and keeps repition stage",
			exerciseID: types.T1,
			initialDetail: Details{
				Name:                 "Squat",
				Repition:             1,
				RepitionsProgression: defaultStages,
				Weight:               100.0,
				Increase:             2.5,
				Reset:                0.85,
			},
			exerciseResult:   true,
			expectedRepition: 1,
			expectedWeight:   102.5,
			expectedReset:    0.85,
		},
		{
			name:       "T1 failure at stage 0 advances to next stage without weight change",
			exerciseID: types.T1,
			initialDetail: Details{
				Name:                 "Squat",
				Repition:             0,
				RepitionsProgression: defaultStages,
				Weight:               100.0,
				Increase:             5.0,
				Reset:                0.85,
			},
			exerciseResult:   false,
			expectedRepition: 1,
			expectedWeight:   100.0,
			expectedReset:    0.85,
		},
		{
			name:       "T1 failure at stage 1 advances to stage 2 without weight change",
			exerciseID: types.T1,
			initialDetail: Details{
				Name:                 "Squat",
				Repition:             1,
				RepitionsProgression: defaultStages,
				Weight:               100.0,
				Increase:             5.0,
				Reset:                0.85,
			},
			exerciseResult:   false,
			expectedRepition: 2,
			expectedWeight:   100.0,
			expectedReset:    0.85,
		},
		{
			name:       "T1 failure at final stage resets stage to 0 and reduces weight by reset factor (15%)",
			exerciseID: types.T1,
			initialDetail: Details{
				Name:                 "Squat",
				Repition:             2,
				RepitionsProgression: defaultStages,
				Weight:               100.0,
				Increase:             5.0,
				Reset:                0.85,
			},
			exerciseResult:   false,
			expectedRepition: 0,
			expectedWeight:   85.0,
			expectedReset:    0.85,
		},
		{
			name:       "T2 failure at final stage resets stage to 0, updates reset and weight",
			exerciseID: types.T2,
			initialDetail: Details{
				Name:     "Bench",
				Repition: 2,
				RepitionsProgression: []Repitions{
					{Set: 3, Rep: 10},
					{Set: 3, Rep: 8},
					{Set: 3, Rep: 6},
				},
				Weight:   60.0,
				Increase: 2.5,
				Reset:    50.0,
			},
			exerciseResult:   false,
			expectedRepition: 0,
			expectedWeight:   52.5, // reset (50) + increase (2.5)
			expectedReset:    52.5,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			prog := &Program{
				Current: types.A1,
				Workouts: map[types.Workout]WorkoutSession{
					types.A1: {
						NextWorkout: types.B1,
						Exercises: map[types.ExerciseID]*Details{
							tc.exerciseID: &tc.initialDetail,
						},
					},
				},
			}

			prog.Update(Result{
				Workout: types.A1,
				Exercises: map[types.ExerciseID]bool{
					tc.exerciseID: tc.exerciseResult,
				},
			})

			got := prog.Workouts[types.A1].Exercises[tc.exerciseID]
			if got.Repition != tc.expectedRepition {
				t.Errorf("got Repition = %d, want %d", got.Repition, tc.expectedRepition)
			}
			if math.Abs(got.Weight-tc.expectedWeight) > 1e-6 {
				t.Errorf("got Weight = %v, want %v", got.Weight, tc.expectedWeight)
			}
			if math.Abs(got.Reset-tc.expectedReset) > 1e-6 {
				t.Errorf("got Reset = %v, want %v", got.Reset, tc.expectedReset)
			}
			if prog.Current != types.B1 {
				t.Errorf("got Current = %v, want %v", prog.Current, types.B1)
			}
		})
	}
}

func TestProgram_Update_Sequential(t *testing.T) {
	prog := &Program{
		Current: types.A1,
		Workouts: map[types.Workout]WorkoutSession{
			types.A1: {
				NextWorkout: types.B1,
				Exercises: map[types.ExerciseID]*Details{
					types.T1: {
						Name:     "Squat",
						Repition: 0,
						RepitionsProgression: []Repitions{
							{Set: 5, Rep: 3},
							{Set: 6, Rep: 2},
							{Set: 10, Rep: 1},
						},
						Weight:   100.0,
						Increase: 5.0,
						Reset:    0.85,
					},
				},
			},
		},
	}

	// 1. Success on Stage 0 -> weight increases by 5, stays on Stage 0
	prog.Update(Result{Workout: types.A1, Exercises: map[types.ExerciseID]bool{types.T1: true}})
	t1 := prog.Workouts[types.A1].Exercises[types.T1]
	if t1.Repition != 0 || t1.Weight != 105.0 {
		t.Fatalf("after success: got stage = %d, weight = %v; want stage = 0, weight = 105.0", t1.Repition, t1.Weight)
	}

	// 2. Fail on Stage 0 -> moves to Stage 1, weight unchanged
	prog.Update(Result{Workout: types.A1, Exercises: map[types.ExerciseID]bool{types.T1: false}})
	if t1.Repition != 1 || t1.Weight != 105.0 {
		t.Fatalf("after fail on stage 0: got stage = %d, weight = %v; want stage = 1, weight = 105.0", t1.Repition, t1.Weight)
	}

	// 3. Fail on Stage 1 -> moves to Stage 2, weight unchanged
	prog.Update(Result{Workout: types.A1, Exercises: map[types.ExerciseID]bool{types.T1: false}})
	if t1.Repition != 2 || t1.Weight != 105.0 {
		t.Fatalf("after fail on stage 1: got stage = %d, weight = %v; want stage = 2, weight = 105.0", t1.Repition, t1.Weight)
	}

	// 4. Fail on Stage 2 (final stage) -> resets to Stage 0, weight reduces by 15% (105 * 0.85 = 89.25)
	prog.Update(Result{Workout: types.A1, Exercises: map[types.ExerciseID]bool{types.T1: false}})
	expectedWeight := 105.0 * 0.85
	if t1.Repition != 0 || math.Abs(t1.Weight-expectedWeight) > 1e-6 {
		t.Fatalf("after fail on stage 2 (reset): got stage = %d, weight = %v; want stage = 0, weight = %v", t1.Repition, t1.Weight, expectedWeight)
	}
}
