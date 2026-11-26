package rest

import (
	"context"
	"net/http"
	"strings"
)

// TokenVerifier, verfies a JWT token and extracts the userID
type TokenVerifier interface {
	VerifyToken(ctx context.Context, idToken string) (string, error) // returns userID
}

// AuthMiddleware checks if an user is authorized
//
// # Parameter
//
//	verifier - an implemenation that meets the TokenVerifter interface
//
// # Errors
//
//	Returns a 401 Unauthorized status code when header is missing or invalid. As well as when the token is invalid or expired
func AuthMiddleware(verifier TokenVerifier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")
			token, ok := strings.CutPrefix(authHeader, "Bearer ")
			if !ok || token == "" {
				http.Error(w, "missing or invalid authorization header", http.StatusUnauthorized)
				return
			}

			ctx := r.Context()

			userID, err := verifier.VerifyToken(ctx, token)
			if err != nil {
				http.Error(w, "invalid or expired token", http.StatusUnauthorized)
				return
			}

			ctx = ContextWithUser(ctx, userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
