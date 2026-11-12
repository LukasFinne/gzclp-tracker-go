package user

import (
	"context"
	"gzclp-tracker/internal/types"
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
	//
	// userID string
	//
	// data *types.User - updated data from the client
	//
	// # Return values
	//
	//	error - returns a update failed error
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
			Weight: progression.Weight,
		}
	}
	return m, nil
}

// Upload calculates the new workout weight, set and rep progession as well as deciding what workout is next.
func (s service) Upload(ctx context.Context, userID string, workoutResult WorkoutResult) error {

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	rotate := user.WorkoutRotation[workoutResult.Workout]

	// Loop iteration:
	// 1. Finds exercises weight and setxrep for each tier (T1,T2,T3)
	// 2. Selects the revelant tier (T1,T2,T3)
	// 3. Updates the weight and setxrep
	for tier, exercise := range workoutResult.Exercises {
		reptitions := user.Exercises[exercise.Name]
		progression := reptitions[tier]
		progression.CalculateProgression(exercise.IsSuccess)
	}

	err = s.userRepo.Save(ctx, userID, &types.User{
		CurrentWorkout: rotate.Next,
		Exercises:      user.Exercises,
	})
	if err != nil {
		return err
	}

	return nil
}
