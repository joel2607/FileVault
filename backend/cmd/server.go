package main

import (
	"log"
	"net/http"
	"os"

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

func main() {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	db := database.Init()
	
	authService := services.NewAuthService(db)

	router := chi.NewRouter()
	router.Use(middleware.AuthMiddleware(authService))

	srv := handler.NewDefaultServer(graphQL.NewExecutableSchema(graphQL.Config{Resolvers: &graphQL.Resolver{DB: db, AuthService: authService}}))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}