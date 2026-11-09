package user

import (
	"chi-simple-server/internal/types"
	"context"
	"fmt"
)

// UserStore is used by services to make user-related queries to firestore
type userRepo interface {
	FindByID(ctx context.Context, userID string) (*types.User, error)
	Save(ctx context.Context, userID string, info *types.User) error
}

// Userservice is for handling user-related data from the database
// And refacroting the data structure to a resonable format for the frontend to use.
type service struct {
	userRepo userRepo
}

func NewService(userRepo userRepo) *service {
	return &service{
		userRepo: userRepo,
	}
}

func (s service) FetchWorkoutInformation(ctx context.Context, userID string) (Workout, error) {
	user, err := s.userRepo.FindByID(ctx, userID)

	if err != nil {
		return Workout{}, fmt.Errorf("Could not find user: %w", err)
	}

	program := user.Program[user.CurrentWorkout]

	m := make(map[types.ExerciseTier]exercise, len(program))

	for k, v := range program {
		repetition := user.Exercises[v.Name]
		progression := repetition[k]
		m[k] = exercise{
			Name:   v.Name,
			Set:    progression.Stages[progression.StageIndex].Sets,
			Rep:    progression.Stages[progression.StageIndex].Reps,
			Weight: progression.Stages[progression.StageIndex].Weight,
		}
	}
	return m, nil
}

func (s service) Upload(ctx context.Context, userID string) error {
	err := s.userRepo.Save(ctx, userID, &types.User{})
	return err
}
