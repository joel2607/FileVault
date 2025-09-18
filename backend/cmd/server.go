package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-joel2607/database"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-joel2607/graphQL"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/viper"
)

func main() {
	// Initialize Viper
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	err := viper.ReadInConfig()   // Find and read the config file
	if err != nil {               // Handle errors reading the config file
		log.Fatalf("Error reading config file, %s", err)
	}
	viper.AutomaticEnv()
	port := viper.GetString("PORT")

	// Initialize the database connection
	database.Init()
	db := database.GetDB()

	// Create a Chi router
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Create a new GraphQL server
	srv := handler.NewDefaultServer(graphQL.NewExecutableSchema(graphQL.Config{Resolvers: &graphQL.Resolver{DB: db}}))

	// Add the GraphQL playground and query handlers to the Chi router
	r.Handle("/", playground.Handler("GraphQL playground", "/query"))
	r.Handle("/query", srv)

	// Simple hello world endpoint
	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World! The backend is running with chi."))
	} )

	// Start the server
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}