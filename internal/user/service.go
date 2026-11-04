package user

import (
	"chi-simple-server/internal/types"
	"context"
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

func (s service) FetchWorkoutInformation(ctx context.Context, userID string) (types.User, error) {

	return types.User{}, nil
}
