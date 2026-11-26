package rest

import "context"

type userKey struct{}

// userFromContext fetches the userID
func UserFromContext(ctx context.Context) (string, bool) {
	user, ok := ctx.Value(userKey{}).(string)
	return user, ok
}

// ContextWithUser sets the userID in the Context
func ContextWithUser(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userKey{}, userID)
}
