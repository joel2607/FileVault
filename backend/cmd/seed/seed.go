// Package main is a standalone script for seeding the database with initial data.
package main

import (
	"log"

	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-joel2607/database"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-joel2607/models"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Could not read config file for seeding: %s. Relying on environment variables.", err)
	}

	viper.AutomaticEnv()
	viper.BindEnv("postgres.host", "POSTGRES_HOST")
	viper.BindEnv("postgres.port", "POSTGRES_PORT")
	viper.BindEnv("postgres.user", "POSTGRES_USER")
	viper.BindEnv("postgres.password", "POSTGRES_PASSWORD")
	viper.BindEnv("postgres.db", "POSTGRES_DB")
}

// Seeding Module to fill database with test data for local development.
func main() {
	log.Println("Starting database seeding...")

	// Override config values for local seeding against Docker container
	viper.Set("postgres.host", "localhost")
	viper.Set("postgres.port", "5433")

	// 1. Initialize the database connection (this also runs migrations)
	database.Init()
	db := database.DB // Use the global DB variable from the database package
	log.Println("Database connection and migration successful.")

	db.Exec("TRUNCATE TABLE users, files, folders, deduplicated_contents, file_sharings RESTART IDENTITY CASCADE")
	log.Println("Deleted existing records.")

	// 2. Create Users
	adminUser := models.User{
		Username:       "admin",
		Email:          "admin@example.com",
		PasswordHash:   "$2a$10$dqgg48GLMDj7AjLHTZ7n7uUO3Ksl9cQTiCWE9.KQWqGsMpUMNBNoG", // sha 256 for "admin" with server salt
		Role:           models.RoleAdmin,
		StorageQuotaKB: 1048576, // 1 GB
	}
	regularUser := models.User{
		Username:       "user",
		Email:          "user@example.com",
		PasswordHash:   "$2a$10$OuvU0Yy9wDLwFBDrdlkvMe0p6wtAWFS7IuuQ7c5q2b1NmFmvwykyW", // sha 256 for "user" with server salt
		Role:           models.RoleUser,
		StorageQuotaKB: 102400, // 100 MB
	}

	// Use a transaction to ensure all or nothing
	tx := db.Begin()

	if err := tx.Create(&adminUser).Error; err != nil {
		tx.Rollback()
		log.Fatalf("Could not create admin user: %v", err)
	}
	if err := tx.Create(&regularUser).Error; err != nil {
		tx.Rollback()
		log.Fatalf("Could not create regular user: %v", err)
	}
	log.Println("Successfully created users.")

	// 3. Create a Folder for the admin user
	rootFolder := models.Folder{
		UserID:     adminUser.ID,
		FolderName: "My Documents",
	}
	if err := tx.Create(&rootFolder).Error; err != nil {
		tx.Rollback()
		log.Fatalf("Could not create root folder: %v", err)
	}
	log.Println("Successfully created a folder.")

	// 4. Create a File record for the admin user in their new folder
	// First, create the deduplicated content record
	deduplicatedContent := models.DeduplicatedContent{
		SHA256Hash:     "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", // SHA256 of an empty string
		ReferenceCount: 1,
	}
	if err := tx.Create(&deduplicatedContent).Error; err != nil {
		tx.Rollback()
		log.Fatalf("Could not create deduplicated content: %v", err)
	}

	// Now, create the file metadata that points to the content
	sampleFile := models.File{
		UserID:          adminUser.ID,
		FileName:        "empty_file.txt",
		MIMEType:        "text/plain",
		Size:            0,
		DeduplicationID: deduplicatedContent.ID,
		FolderID:        &rootFolder.ID,
		Tags:            "[]",
	}
	if err := tx.Create(&sampleFile).Error; err != nil {
		tx.Rollback()
		log.Fatalf("Could not create sample file: %v", err)
	}
	log.Println("Successfully created a file record.")

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	log.Println("Database seeding completed successfully!")
}
