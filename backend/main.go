package main

import (
	"log"
	"net/http"

	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-joel2607/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Initialize the database connection
	database.Init()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	
	// Simple hello world endpoint
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World! The backend is running with chi."))
	})

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
