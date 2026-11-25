package rest

import (
	"context"
	"net/http"
	"strings"
)

// Verfies a JWT token and extracts the userID
type TokenVerifier interface {
	VerifyToken(ctx context.Context, idToken string) (string, error) // returns userID
}

// AuthMiddlware verifies the IdToken then extract the userID
func AuthMiddlware(verifier TokenVerifier) func(http.Handler) http.Handler {
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
