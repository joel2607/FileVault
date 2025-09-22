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
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-joel2607/handlers"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-joel2607/middleware"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-joel2607/services"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-joel2607/services/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
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
	viper.BindEnv("download_token_secret", "DOWNLOAD_TOKEN_SECRET")
	viper.BindEnv("app.base_url", "APP_BASE_URL")
	viper.BindEnv("ratelimit.limit", "RATELIMIT_LIMIT")
}

func main() {
	port := viper.GetString("server.port")
	if port == "" {
		port = defaultPort
	}

	db := database.Init()
	rdb := database.InitRedis()

	// Service Initialization
	authService := services.NewAuthService(db)
	storageProvider := storage.NewLocalStorageProvider(viper.GetString("app.base_url"))
	fileService := services.NewFileService(db, rdb, storageProvider)
	shareService := services.NewShareService(db)

	// Setup Chi Router
	router := chi.NewRouter()

	// Middleware
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
	}).Handler)
	router.Use(middleware.AuthMiddleware(authService))
	router.Use(middleware.RateLimitMiddleware(rdb, viper.GetInt("ratelimit.limit"), 1*time.Second))

	// Setup GraphQL Server
	resolver := &graphQL.Resolver{
		DB:           db,
		RDB:          rdb,
		AuthService:  authService,
		FileService:  fileService,
		ShareService: shareService,
	}
	srv := handler.NewDefaultServer(graphQL.NewExecutableSchema(graphQL.Config{Resolvers: resolver}))

	srv.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for WebSocket
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		InitFunc: middleware.WebSocketInitFunc(authService),
	})

	// Define Routes
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)
	router.Get("/downloads/*", handlers.DownloadHandler)

	// Start Server
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}