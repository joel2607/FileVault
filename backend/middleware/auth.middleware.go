package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-joel2607/models"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-joel2607/services"
)

// ContextKey is a custom type to avoid key collisions in context.
type ContextKey struct {
	Name string
}

// UserCtxKey is the key for storing user information in the context.
var UserCtxKey = &ContextKey{"user"}

// RoleCtxKey is the key for storing user role in the context.
var RoleCtxKey = &ContextKey{"role"}

// AuthErrorCtxKey is the key for storing authentication errors in the context.
var AuthErrorCtxKey = &ContextKey{"auth-error"}

// AuthMiddleware extracts the JWT token, validates it, and then places
// either the user information or a specific authentication error into the request context.
// It no longer blocks the request, allowing downstream handlers (GraphQL/REST) to decide
// how to handle the authentication result.
func AuthMiddleware(authService *services.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				next.ServeHTTP(w, r)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			user, err := authService.GetUserFromToken(tokenString)

			var ctx context.Context
			if err != nil {
				// Put the specific error into the context
				ctx = context.WithValue(r.Context(), AuthErrorCtxKey, err)
			} else {
				// Put the user and role into the context on success
				ctx = context.WithValue(r.Context(), UserCtxKey, user)
				ctx = context.WithValue(ctx, RoleCtxKey, user.Role)
			}

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// AdminOnly is a middleware for REST routes that checks for admin privileges.
// It now also checks for authentication errors placed in the context by the upstream AuthMiddleware.
func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for authentication error first
		if err, ok := r.Context().Value(AuthErrorCtxKey).(error); ok && err != nil {
			http.Error(w, "Forbidden: "+err.Error(), http.StatusForbidden)
			return
		}

		// Check for admin role
		role := r.Context().Value(RoleCtxKey)
		if role == nil || role.(string) != "ADMIN" {
			http.Error(w, "Forbidden: Admins only", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func GetCurrentUser(ctx context.Context) (*models.User, error) {
	if err, ok := ctx.Value(AuthErrorCtxKey).(error); ok && err != nil {
		return nil, err
	}

	userVal := ctx.Value(UserCtxKey)
	if userVal == nil {
		return nil, fmt.Errorf("access denied: no token provided")
	}

	user, ok := userVal.(*models.User)
	if !ok {
		return nil, fmt.Errorf("internal server error: user context value is of the wrong type")
	}

	return user, nil
}
