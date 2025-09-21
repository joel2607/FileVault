package main

import (
	"log"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-joel2607/database"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-joel2607/graphQL"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-joel2607/middleware"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-joel2607/services"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

const defaultPort = "8080"

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Could not read config file: %s. Relying on environment variables.", err)
	}

	viper.AutomaticEnv()
	viper.BindEnv("postgres.host", "POSTGRES_HOST")
	viper.BindEnv("postgres.port", "POSTGRES_PORT")
	viper.BindEnv("postgres.user", "POSTGRES_USER")
	viper.BindEnv("postgres.password", "POSTGRES_PASSWORD")
	viper.BindEnv("postgres.db", "POSTGRES_DB")
	viper.BindEnv("redis.addr", "REDIS_ADDR")
	viper.BindEnv("auth.jwt_expiration_hours", "JWT_EXPIRATION_HOURS")
	viper.BindEnv("jwt_auth_secret", "JWT_AUTH_SECRET")
}

func main() {
	port := viper.GetString("server.port")
	if port == "" {
		port = defaultPort
	}

	db := database.Init()
	redisClient := database.InitRedis()

	authService := services.NewAuthService(db)
	fileService := services.NewFileService(db, redisClient)
	shareService := services.NewShareService(db)

	router := chi.NewRouter()
	router.Use(middleware.AuthMiddleware(authService))
	router.Use(middleware.RedisRateLimiter(redisClient, 2, 1*time.Second))

	srv := handler.NewDefaultServer(graphQL.NewExecutableSchema(graphQL.Config{Resolvers: &graphQL.Resolver{
		DB:           db,
		RDB:          redisClient,
		AuthService:  authService,
		FileService:  fileService,
		ShareService: shareService,
	}}))

	srv.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		InitFunc: middleware.WebSocketInitFunc(authService),
	})

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}