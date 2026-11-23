package user

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"gzclp-tracker/internal/types"
	"gzclp-tracker/internal/user/testdata"
)

// mockUserRepo is a test double implementing the userRepo interface.
type mockUserRepo struct {
	findByIDFunc func(ctx context.Context, userID string) (*Program, error)
	saveFunc     func(ctx context.Context, userID string, info *Program) error
	savedUser    *Program
}

func (m *mockUserRepo) FindByID(ctx context.Context, userID string) (*Program, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, userID)
	}
	return nil, nil
}

func (m *mockUserRepo) Save(ctx context.Context, userID string, info *Program) error {
	m.savedUser = info
	if m.saveFunc != nil {
		return m.saveFunc(ctx, userID, info)
	}
	return nil
}

func loadTestProgram(t *testing.T) *Program {
	t.Helper()
	entity := testdata.LoadDefaultProgramEntity(t)

	workouts := make(map[types.Workout]WorkoutSession, len(entity.Workouts))
	for wID, session := range entity.Workouts {
		exMap := make(map[types.ExerciseID]*ExerciseState, len(session.Exercises))
		for eID, d := range session.Exercises {
			progressions := make([]Repitions, len(d.RepitionsProgression))
			for i, p := range d.RepitionsProgression {
				progressions[i] = Repitions{
					Set: p.Set,
					Rep: p.Rep,
				}
			}
			exMap[eID] = &ExerciseState{
				Name:                 d.Name,
				Repition:             d.Repition,
				RepitionsProgression: progressions,
				Weight:               d.Weight,
				Increase:             d.Increase,
				Reset:                d.Reset,
			}
		}
		workouts[wID] = WorkoutSession{
			NextWorkout: session.NextWorkout,
			Exercises:   exMap,
		}
	}

	return &Program{
		Current:  entity.Current,
		Workouts: workouts,
	}
}

func TestFetchWorkoutInformation(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		repo    userRepo
		want    Exercises
		wantErr bool
	}{
		{
			name:   "success - returns workout information",
			userID: "test-user-id",
			repo: &mockUserRepo{
				findByIDFunc: func(ctx context.Context, userID string) (*Program, error) {
					return loadTestProgram(t), nil
				},
			},
			want: Exercises{
				WorkoutID: types.B1,
				Exercises: map[types.ExerciseID]Exercise{
					types.T1: {
						Name:   "OHP",
						Set:    5,
						Rep:    3,
						Weight: 10,
					},
					types.T2: {
						Name:   "Deadlift",
						Set:    3,
						Rep:    10,
						Weight: 10,
					},
					types.T3: {
						Name:   "Dumbbell row",
						Set:    3,
						Rep:    15,
						Weight: 10,
					},
				},
			},
			wantErr: false,
		},
		{
			name:   "error - repository fails to find user",
			userID: "non-existent-user",
			repo: &mockUserRepo{
				findByIDFunc: func(ctx context.Context, userID string) (*Program, error) {
					return nil, errors.New("database connection failed")
				},
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := NewService(tc.repo)

			got, err := svc.FetchWorkoutInformation(context.Background(), tc.userID)
			if (err != nil) != tc.wantErr {
				t.Fatalf("FetchWorkoutInformation() error = %v, wantErr %v", err, tc.wantErr)
			}
			if tc.wantErr {
				return
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("FetchWorkoutInformation() got = %+v, want %+v", got, tc.want)
			}
		})
	}
}

func TestUpload(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		result  Result
		repo    userRepo
		want    *Program
		wantErr bool
	}{
		{
			name:   "success - updates the data correctly",
			userID: "test-user-id",
			result: Result{
				WorkoutID: types.A1,
				Exercises: map[types.ExerciseID]bool{
					types.T1: false, // Squat T1 fails -> advances stageIndex (Repition) to 1
					types.T2: true,  // Bench T2 succeeds -> weight increased by 10 (10 + 10 = 20)
					types.T3: true,  // Lat pulldown T3 succeeds -> weight increased by 2.5 (10 + 2.5 = 12.5)
				},
			},
			repo: &mockUserRepo{
				saveFunc: func(ctx context.Context, userID string, info *Program) error {
					return nil
				},
				findByIDFunc: func(ctx context.Context, userID string) (*Program, error) {
					return loadTestProgram(t), nil
				},
			},
			want: func() *Program {
				p := loadTestProgram(t)
				// Squat T1 failed -> stage 1
				p.Workouts[types.A1].Exercises[types.T1].Repition = 1
				// Bench T2 succeeded -> weight 20
				p.Workouts[types.A1].Exercises[types.T2].Weight = 20
				// Lat pulldown T3 succeeded -> weight 12.5
				p.Workouts[types.A1].Exercises[types.T3].Weight = 12.5
				// Current advances to next workout of A1 ("B1")
				p.Current = p.Workouts[types.A1].NextWorkout
				return p
			}(),
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := NewService(tc.repo)

			err := svc.Upload(context.Background(), tc.userID, tc.result)
			if (err != nil) != tc.wantErr {
				t.Fatalf("Upload() error = %v, wantErr %v", err, tc.wantErr)
			}
			if tc.wantErr {
				return
			}

			mockRepo, ok := tc.repo.(*mockUserRepo)
			if !ok || mockRepo.savedUser == nil {
				t.Fatal("Upload() did not save user")
			}

			if !reflect.DeepEqual(mockRepo.savedUser, tc.want) {
				t.Errorf("Upload() got = %+v, want %+v", mockRepo.savedUser, tc.want)
			}
		})
	}
}
