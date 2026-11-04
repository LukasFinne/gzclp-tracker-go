package firestore

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

// Init return an authanticated firestore client instance using the applications default credentials
//
// Parameters
//
//		ctx - handles request cancellations and timeouts
//	 serviceAccountFilePath - path to the serviceAccount.json file
func Init(ctx context.Context, serviceAccountFilePath string) (*firestore.Client, error) {
	option := option.WithAuthCredentialsFile(option.ServiceAccount, serviceAccountFilePath)
	app, err := firebase.NewApp(ctx, nil, option)
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase app: %w", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, fmt.Errorf("error initializing firestore client: %w", err)
	}

	return client, nil
}
