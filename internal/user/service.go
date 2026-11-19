package user

import (
	"context"
	"gzclp-tracker/internal/types"
)

// userRepo is used by services to make user-related queries to firestore
type userRepo interface {
	// FindByID fetches the users training program
	//
	// # Parameters
	//
	//	ctx context.Context - handles request cancellations and timeouts
	//	userID string
	//
	// # Return values
	//
	//	*Program is the users training program and contains the workout id and the corresponding exercises
	//	error - returns either a Notfound error or a Generic failed error
	FindByID(ctx context.Context, userID string) (*Program, error)
	// Save updates the related fields in the data through an Transaction
	//
	// # Parameters
	//
	// ctx context.Context - handles request cancellations and timeouts
	//
	// userID string
	//
	// data *types.User - updated data from the client
	//
	// # Return values
	//
	//	error - returns a update failed error
	Save(ctx context.Context, userID string, info *Program) error
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

func (s service) FetchWorkoutInformation(ctx context.Context, userID string) (map[types.ExerciseID]types.Info, error) {
	program, err := s.userRepo.FindByID(ctx, userID)

	if err != nil {
		return map[types.ExerciseID]types.Info{}, err
	}

	workout := NewWorkout(program)

	return workout.Exercise.Info(), nil
}

// Upload calculates the new workout weight, set and rep progession as well as deciding what workout is next.
func (s service) Upload(ctx context.Context, userID string, result Result) error {
	program, err := s.userRepo.FindByID(ctx, userID)

	if err != nil {
		return err
	}

	workout := NewWorkout(program)

	workout.Exercise.Update(result)

	err = s.userRepo.Save(ctx, userID, program)
	if err != nil {
		return err
	}

	return nil
}
