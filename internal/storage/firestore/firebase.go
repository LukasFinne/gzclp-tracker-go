package firestore

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

// App intializses the firebase app with a service account file
//
// Parameters
//
//	ctx - handles request cancellations and timeouts
//	serviceAccountFilePath - path to the serviceAccount.json file
//
// # Return
//
//  1. A firebase app instance
//  2. An Error is returned when failing to initialize the instance
func App(ctx context.Context, serviceAccountFilePath string) (*firebase.App, error) {
	opt := option.WithAuthCredentialsFile(option.ServiceAccount, serviceAccountFilePath)

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase app: %w", err)
	}

	return app, nil
}

func DB(ctx context.Context, app *firebase.App) (*firestore.Client, error) {
	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting Firestore client: %w", err)
	}

	return client, nil
}

func Auth(ctx context.Context, app *firebase.App) (*FirebaseAuth, error) {
	client, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting Auth client: %w", err)
	}

	return &FirebaseAuth{
		client: client,
	}, nil
}

type FirebaseAuth struct {
	client *auth.Client
}

func (fv *FirebaseAuth) VerifyToken(ctx context.Context, idToken string) (string, error) {
	if fv == nil || fv.client == nil {
		return "", errors.New("firebase auth client is not initialized")
	}

	token, err := fv.client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return "", fmt.Errorf("error verifying ID token: %w", err)
	}

	return token.UID, nil
}
