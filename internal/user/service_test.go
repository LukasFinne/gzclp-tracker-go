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
}

func (m *mockUserRepo) FindByID(ctx context.Context, userID string) (*types.User, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, userID)
	}
	return nil, nil
}

func (m *mockUserRepo) Save(ctx context.Context, userID string, info *types.User) error {
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
					Weight: 14.5,
				},
				"T2": exercise{
					Name:   "Deadlift",
					Set:    3,
					Rep:    10,
					Weight: 25,
				},
				"T3": exercise{
					Name:   "Dumbbell row",
					Set:    3,
					Rep:    15,
					Weight: 8,
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
