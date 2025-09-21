package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-joel2607/models"
	"github.com/go-redis/redis/v8"
)

// RateLimitMiddleware provides a Redis-backed rate limiter.
// It checks the number of requests from a user within a given time window
// and blocks requests that exceed the limit.
// It returns a structured GraphQL error if the request is for the GraphQL endpoint.
func RateLimitMiddleware(rdb *redis.Client, defaultLimit int, defaultWindow time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Try to get user from context
			user, ok := r.Context().Value(UserCtxKey).(*models.User)
			if !ok {
				// If no user, proceed without rate limiting
				next.ServeHTTP(w, r)
				return
			}

			// Use user's API rate limit, or the default
			limit := user.APIRateLimit
			if limit == 0 {
				limit = defaultLimit
			}
			window := defaultWindow

			key := fmt.Sprintf("rate_limit:%d", user.ID)
			now := time.Now().UnixNano()
			windowStart := now - window.Nanoseconds()

			// Use a Redis pipeline for efficiency
			pipe := rdb.TxPipeline()
			pipe.ZRemRangeByScore(r.Context(), key, "0", fmt.Sprintf("%d", windowStart))
			pipe.ZAdd(r.Context(), key, &redis.Z{Score: float64(now), Member: now})
			pipe.ZCard(r.Context(), key)
			pipe.Expire(r.Context(), key, window)
			cmds, err := pipe.Exec(r.Context())

			if err != nil {
				// If Redis fails, it's better to let the request through than to block everyone
				fmt.Printf("Redis error in rate limiter: %v", err)
				next.ServeHTTP(w, r)
				return
			}

			count := cmds[2].(*redis.IntCmd).Val()

			if count > int64(limit) {
				// Check if it's a GraphQL request
				isGraphQL := r.Header.Get("Content-Type") == "application/json"

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)

				if isGraphQL {
					// Send a proper GraphQL error response
					response := map[string]interface{}{
						"errors": []map[string]interface{}{
							{
								"message": "You have exceeded the rate limit.",
								"extensions": map[string]string{
									"code": "RATE_LIMIT_EXCEEDED",
								},
							},
						},
					}
					json.NewEncoder(w).Encode(response)
				} else {
					// Send a simple JSON error for non-GraphQL requests
					json.NewEncoder(w).Encode(map[string]string{"error": "You have exceeded the rate limit."})
				}
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}