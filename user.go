package main

import (
	"context"
)

type User struct {
	userID       string
	UserProvider UserProvider
}

func newUser(userID string, UserProvider UserProvider) User {
	return User{
		userID:       userID,
		UserProvider: UserProvider,
	}
}

func (u User) fetchWorkoutInformation(ctx context.Context) (*WorkoutInfo, error) {
	info, err := u.UserProvider.findByID(ctx, u.userID)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (u User) upload(ctx context.Context) error {

	return nil
}
