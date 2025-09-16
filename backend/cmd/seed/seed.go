// Package main is a standalone script for seeding the database with initial data.
package main

import (
	"log"
	"os"

	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-joel2607/database"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-joel2607/models"
)

func main() {
	log.Println("Starting database seeding...")

	// Hardcode database connection details for local development seeding.
	os.Setenv("POSTGRES_HOST", "localhost")
	os.Setenv("POSTGRES_PORT", "5433")

	// 1. Initialize the database connection (this also runs migrations)
	database.Init()
	db := database.DB // Use the global DB variable from the database package
	log.Println("Database connection and migration successful.")

	// 2. Create Users
	adminUser := models.User{
		Email:          "admin@example.com",
		Role:           models.RoleAdmin,
		StorageQuotaMB: 1024, // 1 GB
	}
	regularUser := models.User{
		Email:          "user@example.com",
		Role:           models.RoleUser,
		StorageQuotaMB: 100, // 100 MB
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
		ParentFolderID:  &rootFolder.ID,
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
