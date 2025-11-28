package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/todoist/backend/pkg/jwt"
)

type contextKey string

const UserIDKey contextKey = "user_id"
const EmailKey contextKey = "email"

// Auth middleware validates JWT tokens
func Auth(jwtSecret string) func(http.Handler) http.Handler {
	jwtService := jwt.NewService(jwtSecret, 0, 0)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error":{"code":"UNAUTHORIZED","message":"missing authorization header"}}`, http.StatusUnauthorized)
				return
			}

			// Extract token
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, `{"error":{"code":"UNAUTHORIZED","message":"invalid authorization header format"}}`, http.StatusUnauthorized)
				return
			}

			token := parts[1]

			// Validate token
			claims, err := jwtService.ValidateToken(token)
			if err != nil {
				if err == jwt.ErrExpiredToken {
					http.Error(w, `{"error":{"code":"UNAUTHORIZED","message":"token expired"}}`, http.StatusUnauthorized)
					return
				}
				http.Error(w, `{"error":{"code":"UNAUTHORIZED","message":"invalid token"}}`, http.StatusUnauthorized)
				return
			}

			// Add user info to context
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, EmailKey, claims.Email)

			// Add headers for downstream services
			r.Header.Set("X-User-ID", claims.UserID.String())
			r.Header.Set("X-User-Email", claims.Email)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
