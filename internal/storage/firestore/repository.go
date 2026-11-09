package firestore

import (
	"chi-simple-server/internal/types"
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userRepo struct {
	client *firestore.Client
}

func NewUserRepo(client *firestore.Client) *userRepo {
	return &userRepo{
		client: client,
	}
}

func (up *userRepo) FindByID(ctx context.Context, userID string) (*types.User, error) {
	docSnap, err := up.client.Collection("users").Doc(userID).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, fmt.Errorf("workout not found for user: %s", userID)
		}
		return nil, fmt.Errorf("failed to fetch workout: %w", err)
	}

	var user types.User
	if err := docSnap.DataTo(&user); err != nil {
		return nil, fmt.Errorf("failed to decode document: %w", err)
	}
	return &user, nil
}

func (up *userRepo) Save(ctx context.Context, userID string, info *types.User) error {
	docRef := up.client.Collection("users").Doc(userID)
	return up.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		return tx.Set(docRef, info, firestore.MergeAll)
	})
}
