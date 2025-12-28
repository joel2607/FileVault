package middleware

import (
	"context"
	"log"
	"strings"

	"github.com/joel2607/FileVault/services"
	"github.com/99designs/gqlgen/graphql/handler/transport"
)

// WebSocketInitFunc is used to authenticate WebSocket connections.
// It extracts the JWT from the connection parameters, validates it,
// and adds the user to the context.
func WebSocketInitFunc(authService *services.AuthService) transport.WebsocketInitFunc {
	return func(ctx context.Context, initPayload transport.InitPayload) (context.Context, *transport.InitPayload, error) {
		// Attempt to retrieve the token from the authorization header.
		authHeader, ok := initPayload["Authorization"].(string)
		if !ok {
			// This is not an error, as it allows tools like the GraphQL Playground
			// to connect without authentication.
			log.Println("WebSocket connection attempt without Authorization header.")
			return ctx, &initPayload, nil
		}

		// The token is expected to be in the format "Bearer <token>".
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == "" {
			log.Println("WebSocket connection attempt with an empty token.")
			return ctx, &initPayload, nil
		}

		// Validate the token and retrieve the user.
		user, err := authService.GetUserFromToken(tokenStr)
		if err != nil {
			// If the token is invalid or expired, we log the error but do not
			// fail the connection. This maintains consistency with the HTTP middleware.
			log.Printf("WebSocket authentication error: %v", err)
			ctx = context.WithValue(ctx, AuthErrorCtxKey, err)
			return ctx, &initPayload, nil
		}

		// If authentication is successful, add the user and their role to the context.
		ctx = context.WithValue(ctx, UserCtxKey, user)
		ctx = context.WithValue(ctx, RoleCtxKey, user.Role)

		return ctx, &initPayload, nil
	}
}