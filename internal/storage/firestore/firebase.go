package firestore

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

type Clients struct {
	Firestore *firestore.Client
	Auth      *auth.Client
}

// Init return an authanticated firestore client instance using the applications default credentials
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
func Init(ctx context.Context, serviceAccountFilePath string) (*firebase.App, error) {
	option := option.WithAuthCredentialsFile(option.ServiceAccount, serviceAccountFilePath)
	app, err := firebase.NewApp(ctx, nil, option)
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

func Auth(ctx context.Context, app *firebase.App) (firebaseAuth, error) {
	client, err := app.Auth(ctx)
	if err != nil {
		return firebaseAuth{}, fmt.Errorf("error getting Auth client: %w", err)
	}

	return firebaseAuth{
		client: client,
	}, nil
}

type firebaseAuth struct {
	client *auth.Client
}

func (fv firebaseAuth) VerifyToken(ctx context.Context, idToken string) (string, error) {

	token, err := fv.client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return "", fmt.Errorf("error verifying ID token: %v\n", err)
	}

	return token.UID, nil
}
