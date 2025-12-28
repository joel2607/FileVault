package graphQL

//go:generate go run github.com/99designs/gqlgen generate

import (
	"github.com/joel2607/FileVault/services"
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