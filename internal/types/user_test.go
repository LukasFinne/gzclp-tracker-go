package types

import (
	"math"
	"testing"
)

func TestProgression_CalculateProgression(t *testing.T) {
	defaultStages := []Stage{
		{Sets: 5, Reps: 3},
		{Sets: 6, Reps: 2},
		{Sets: 10, Reps: 1},
	}

	tests := []struct {
		name               string
		initialProgression Progression
		exerciseResult     bool
		expectedStageIndex int
		expectedWeight     float64
	}{
		{
			name: "success at stage 0 increases weight and keeps stage index",
			initialProgression: Progression{
				StageIndex:     0,
				Stages:         defaultStages,
				Weight:         100.0,
				WeightIncrease: 5.0,
			},
			exerciseResult:     true,
			expectedStageIndex: 0,
			expectedWeight:     105.0,
		},
		{
			name: "success at intermediate stage increases weight and keeps stage index",
			initialProgression: Progression{
				StageIndex:     1,
				Stages:         defaultStages,
				Weight:         100.0,
				WeightIncrease: 2.5,
			},
			exerciseResult:     true,
			expectedStageIndex: 1,
			expectedWeight:     102.5,
		},
		{
			name: "failure at stage 0 advances to next stage without weight change",
			initialProgression: Progression{
				StageIndex:     0,
				Stages:         defaultStages,
				Weight:         100.0,
				WeightIncrease: 5.0,
			},
			exerciseResult:     false,
			expectedStageIndex: 1,
			expectedWeight:     100.0,
		},
		{
			name: "failure at stage 1 advances to stage 2 without weight change",
			initialProgression: Progression{
				StageIndex:     1,
				Stages:         defaultStages,
				Weight:         100.0,
				WeightIncrease: 5.0,
			},
			exerciseResult:     false,
			expectedStageIndex: 2,
			expectedWeight:     100.0,
		},
		{
			name: "failure at final stage resets stage to 0 and reduces weight by 15%",
			initialProgression: Progression{
				StageIndex:            2,
				Stages:                defaultStages,
				Weight:                100.0,
				WeightIncrease:        5.0,
				WeightResetPercentage: 0.85,
			},
			exerciseResult:     false,
			expectedStageIndex: 0,
			expectedWeight:     85.0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := tc.initialProgression
			p.CalculateProgression(tc.exerciseResult)

			if p.StageIndex != tc.expectedStageIndex {
				t.Errorf("got StageIndex = %d, want %d", p.StageIndex, tc.expectedStageIndex)
			}
			if math.Abs(p.Weight-tc.expectedWeight) > 1e-6 {
				t.Errorf("got Weight = %v, want %v", p.Weight, tc.expectedWeight)
			}
		})
	}
}

func TestProgression_CalculateProgression_Sequential(t *testing.T) {
	p := Progression{
		StageIndex: 0,
		Stages: []Stage{
			{Sets: 5, Reps: 3},
			{Sets: 6, Reps: 2},
			{Sets: 10, Reps: 1},
		},
		Weight:                100.0,
		WeightIncrease:        5.0,
		WeightResetPercentage: 0.85,
	}

	// 1. Success on Stage 0 -> weight increases by 5, stays on Stage 0
	p.CalculateProgression(true)
	if p.StageIndex != 0 || p.Weight != 105.0 {
		t.Fatalf("after success: got stage = %d, weight = %v; want stage = 0, weight = 105.0", p.StageIndex, p.Weight)
	}

	// 2. Fail on Stage 0 -> moves to Stage 1, weight unchanged
	p.CalculateProgression(false)
	if p.StageIndex != 1 || p.Weight != 105.0 {
		t.Fatalf("after fail on stage 0: got stage = %d, weight = %v; want stage = 1, weight = 105.0", p.StageIndex, p.Weight)
	}

	// 3. Fail on Stage 1 -> moves to Stage 2, weight unchanged
	p.CalculateProgression(false)
	if p.StageIndex != 2 || p.Weight != 105.0 {
		t.Fatalf("after fail on stage 1: got stage = %d, weight = %v; want stage = 2, weight = 105.0", p.StageIndex, p.Weight)
	}

	// 4. Fail on Stage 2 (final stage) -> resets to Stage 0, weight reduces by 15% (105 * 0.85 = 89.25)
	p.CalculateProgression(false)
	expectedWeight := 105.0 * 0.85
	if p.StageIndex != 0 || math.Abs(p.Weight-expectedWeight) > 1e-6 {
		t.Fatalf("after fail on stage 2 (reset): got stage = %d, weight = %v; want stage = 0, weight = %v", p.StageIndex, p.Weight, expectedWeight)
	}
}
