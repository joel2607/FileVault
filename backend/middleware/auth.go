package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-joel2607/services"
)

type contextKey struct {
	name string
}

var userCtxKey = &contextKey{"user"}
var roleCtxKey = &contextKey{"role"}

// AuthMiddleware extracts the JWT token from the Authorization header,
// validates it, and adds the user information to the request context.
// It also attaches role information for RBAC in context.
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
			if err != nil {
				http.Error(w, "Invalid token", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), userCtxKey, user)
			ctx = context.WithValue(ctx, roleCtxKey, user.Role)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// REST routes can use this for RBAC. GraphQL resolvers can use the context directly and should not use this function.
func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role := r.Context().Value(roleCtxKey)
		if role == nil || role.(string) != "ADMIN" {
			http.Error(w, "Forbidden: Admins only", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}