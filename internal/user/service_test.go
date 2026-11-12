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
	findByIDFunc func(ctx context.Context, userID string) (*types.User, error)
	saveFunc     func(ctx context.Context, userID string, info *types.User) error
	savedUser    *types.User
}

func (m *mockUserRepo) FindByID(ctx context.Context, userID string) (*types.User, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, userID)
	}
	return nil, nil
}

func (m *mockUserRepo) Save(ctx context.Context, userID string, info *types.User) error {
	m.savedUser = info
	if m.saveFunc != nil {
		return m.saveFunc(ctx, userID, info)
	}
	return nil
}

func TestFetchWorkoutInformation(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		repo    userRepo
		want    Workout
		wantErr bool
	}{
		{
			name:   "success - returns workout information",
			userID: "test-user-id",
			repo: &mockUserRepo{
				findByIDFunc: func(ctx context.Context, userID string) (*types.User, error) {
					return testdata.LoadDefaultUser(t), nil
				},
			},
			want: Workout{
				"T1": exercise{
					Name:   "OHP",
					Set:    5,
					Rep:    3,
					Weight: 10,
				},
				"T2": exercise{
					Name:   "Deadlift",
					Set:    3,
					Rep:    10,
					Weight: 10,
				},
				"T3": exercise{
					Name:   "Dumbbell row",
					Set:    3,
					Rep:    15,
					Weight: 10,
				},
			},
			wantErr: false,
		},
		{
			name:   "error - repository fails to find user",
			userID: "non-existent-user",
			repo: &mockUserRepo{
				findByIDFunc: func(ctx context.Context, userID string) (*types.User, error) {
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
		result  WorkoutResult
		repo    userRepo
		want    types.User
		wantErr bool
	}{
		{
			name:   "success - updates the data correctly",
			userID: "test-user-id",
			result: WorkoutResult{
				Workout: types.A1,
				Exercises: map[types.ExerciseTier]ExerciseResult{
					types.T1: {
						Name:      "Squat",
						IsSuccess: false,
					},
					types.T2: {
						Name:      "Bench",
						IsSuccess: true,
					},
					types.T3: {
						Name:      "Lat pulldown",
						IsSuccess: true,
					},
				},
			},
			repo: &mockUserRepo{
				saveFunc: func(ctx context.Context, userID string, info *types.User) error {
					return nil
				},
				findByIDFunc: func(ctx context.Context, userID string) (*types.User, error) {
					return testdata.LoadDefaultUser(t), nil
				},
			},
			want: types.User{
				CurrentWorkout: types.B1,
				Exercises: func() map[string]types.Repition {
					u := testdata.LoadDefaultUser(t)
					// Squat T1 failed -> moves to stage 1
					u.Exercises["Squat"][types.T1].StageIndex = 1
					u.Exercises["Bench"][types.T2].Weight = 20
					return u.Exercises
				}(),
			},
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

			if !reflect.DeepEqual(*mockRepo.savedUser, tc.want) {
				t.Errorf("Upload() got = %+v, want %+v", *mockRepo.savedUser, tc.want)
			}
		})
	}
}
