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

func (up *userRepo) FindByID(ctx context.Context, userID string) (*user.Program, error) {
	doc, err := up.client.Collection("users").Doc(userID).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, types.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	var program ProgramEntity
	if err := doc.DataTo(&program); err != nil {
		return nil, fmt.Errorf("failed to find document: %w", err)
	}
	return program.ToProgram(), nil
}

func (up *userRepo) Save(ctx context.Context, userID string, info *user.Program) error {

	programEntiy := ToProgramEntity(*info)

	docRef := up.client.Collection("users").Doc(userID)
	return up.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		return tx.Set(docRef, programEntiy, firestore.MergeAll)
	})
}
