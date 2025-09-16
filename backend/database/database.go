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

func Init() {
	// Set up Viper to read .env file
	viper.SetConfigFile("../.env")
	viper.AutomaticEnv()

	viper.SetDefault("POSTGRES_HOST", "localhost")
	viper.SetDefault("POSTGRES_PORT", "5432")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	// Get database credentials from Viper
	dbHost := viper.GetString("POSTGRES_HOST")
	dbUser := viper.GetString("POSTGRES_USER")
	dbPassword := viper.GetString("POSTGRES_PASSWORD")
	dbName := viper.GetString("POSTGRES_DB")
	dbPort := viper.GetString("POSTGRES_PORT")

	// Database connection
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName, dbPort)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	// AutoMigrate the schema
	err = DB.AutoMigrate(&models.User{}, &models.Folder{}, &models.File{})
	if err != nil {
		log.Fatal("failed to migrate database")
	}
}
