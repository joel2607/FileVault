package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-joel2607/database"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-joel2607/models"
	"github.com/go-redis/redis/v8"
)

func RedisRateLimiter(rdb *redis.Client, limit int, window time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value("user").(*models.User)
			if !ok {
				next.ServeHTTP(w, r)
				return
			}

			key := fmt.Sprintf("rate_limit:%d", user.ID)
			now := time.Now().UnixNano()
			windowStart := now - window.Nanoseconds()

			pipe := rdb.TxPipeline()
			pipe.ZRemRangeByScore(database.Ctx, key, "0", fmt.Sprintf("%d", windowStart))
			pipe.ZAdd(database.Ctx, key, &redis.Z{Score: float64(now), Member: now})
			pipe.ZCard(database.Ctx, key)
			pipe.Expire(database.Ctx, key, window)

			cmds, err := pipe.Exec(database.Ctx)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			count := cmds[2].(*redis.IntCmd).Val()

			if count > int64(limit) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				json.NewEncoder(w).Encode(map[string]string{"error": "You have exceeded the rate limit."})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}