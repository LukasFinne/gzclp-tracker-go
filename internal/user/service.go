package user

import (
	"context"
)

// userRepo is used by services to do user-related queries to the database
type userRepo interface {
	// FindByID fetches the user training program data
	//
	// # Parameters
	//
	//	ctx context.Context - handles request cancellations and timeouts
	//	userID string - is used to fetch the document from the database
	//
	// # Return values
	//
	//	*Routines - contains the user training program data
	//
	// Errors
	//   1. Returns ErrNotFound when the document does not exist
	// 	 2. Returns failed to fetch user error when fetching the doc fails
	// 	 3. Returns failed to decode document when decoding to the struct fails
	FindByID(ctx context.Context, userID string) (*Routines, error)
	// Save updates fields set in routines through an transaction
	//
	// # Parameters
	//
	// ctx context.Context - handles request cancellations and timeouts
	//
	// userID string
	//
	// routines *Routines - the updated routine data
	//
	// # Return values
	//
	//	Errors
	// 	 1. Returns failed to save routines error when the transaction fails
	Save(ctx context.Context, userID string, routines *Routines) error
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

func (s service) StartWorkout(ctx context.Context, userID string) (Exercises, error) {

	program, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return Exercises{}, err
	}

	workout := NewWorkout(program)

	return workout.Exercise.Start(), nil
}

// FinishWorkout save the new routine data by workout result
func (s service) FinishWorkout(ctx context.Context, userID string, result Result) error {

	program, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	workout := NewWorkout(program)

	updatedProgram, err := workout.Exercise.Finish(result)
	if err != nil {
		return err
	}

	err = s.userRepo.Save(ctx, userID, &updatedProgram)
	if err != nil {
		return err
	}

	return nil
}
