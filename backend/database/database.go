package database

import (
	"fmt"
	"log"

	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-joel2607/models"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() *gorm.DB {

	// Set up automatic environment variable reading
	viper.AutomaticEnv()

	// Get database credentials from Viper
	dbHost := viper.GetString("postgres.host")
	dbUser := viper.GetString("postgres.user")
	dbPassword := viper.GetString("postgres.password")
	dbName := viper.GetString("postgres.db")
	dbPort := viper.GetString("postgres.port")

	// Database connection
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName, dbPort)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// AutoMigrate the schema
	err = DB.AutoMigrate(&models.User{}, &models.DeduplicatedContent{}, &models.Folder{}, &models.File{}, &models.FileSharing{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	return DB
}
