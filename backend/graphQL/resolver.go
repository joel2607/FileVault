package graphQL

import (
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-joel2607/services"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB           *gorm.DB
	RDB          *redis.Client
	AuthService  *services.AuthService
	FileService  *services.FileService
	ShareService *services.ShareService
}