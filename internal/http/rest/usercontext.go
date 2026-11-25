package rest

import "context"

type userKey struct{}

func UserFromContext(ctx context.Context) (string, bool) {
	user, ok := ctx.Value(userKey{}).(string)
	return user, ok
}

func ContextWithUser(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userKey{}, userID)
}
