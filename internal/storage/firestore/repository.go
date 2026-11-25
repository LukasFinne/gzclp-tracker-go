package firestore

import (
	"context"
	"fmt"
	"gzclp-tracker/internal/types"
	"gzclp-tracker/internal/user"

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

func (up *userRepo) FindByID(ctx context.Context, userID string) (*user.Routines, error) {
	doc, err := up.client.Collection("users").Doc(userID).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, types.ErrNotFound
		}
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	var routines RoutinesEntity
	if err := doc.DataTo(&routines); err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, types.ErrNotFound
		}
		return nil, fmt.Errorf("failed to decode document: %w", err)
	}
	return routines.ToRoutines(), nil
}

func (up *userRepo) Save(ctx context.Context, userID string, routines *user.Routines) error {
	routineEntity := ToRoutinesEntity(*routines)

	docRef := up.client.Collection("users").Doc(userID)

	err := up.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		return tx.Set(docRef, routineEntity, firestore.MergeAll)
	})
	if err != nil {
		return fmt.Errorf("failed to save routines for user: %w", err)
	}

	return nil
}
