package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-joel2607/database"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-joel2607/graphQL"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-joel2607/middleware"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-joel2607/services"
	"github.com/go-chi/chi/v5"
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
	viper.BindEnv("auth.jwt_expiration_hours", "JWT_EXPIRATION_HOURS")
	viper.BindEnv("jwt_auth_secret", "JWT_AUTH_SECRET")
}

func main() {
	port := viper.GetString("server.port")
	if port == "" {
		port = defaultPort
	}

	db := database.Init()
	
	authService := services.NewAuthService(db)
	fileService := services.NewFileService(db)
	shareService := services.NewShareService(db)

	router := chi.NewRouter()
	router.Use(middleware.AuthMiddleware(authService))

	srv := handler.NewDefaultServer(graphQL.NewExecutableSchema(graphQL.Config{Resolvers: &graphQL.Resolver{
		DB: db, 
		AuthService: authService, 
		FileService: fileService, 
		ShareService: shareService,
	}}))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}