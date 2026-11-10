package user

import (
	"chi-simple-server/internal/types"
	"context"
)

// userRepo is used by services to make user-related queries to firestore
type userRepo interface {
	// FindByID fetches all data for a user
	//
	// # Parameters
	//
	//	ctx context.Context - handles request cancellations and timeouts
	//	userID string
	//
	// # Return values
	//
	//	*types.User - Contains the users workout information
	//	error - returns either a Notfound error or a Generic failed error
	FindByID(ctx context.Context, userID string) (*types.User, error)
	// Save updates the related fields in the data through an Transaction
	//
	// # Parameters
	//
	// ctx context.Context - handles request cancellations and timeouts
	// userID string
	// info *WorkoutInfo
	//
	// # Return values
	//
	//	error - returns an error if the updates has failed
	Save(ctx context.Context, userID string, info *types.User) error
}

// service is for handling user-related data from the database
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
		return Workout{}, err
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
