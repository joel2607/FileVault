package graphQL

//go:generate go run github.com/99designs/gqlgen generate

import (
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-joel2607/services"
	"gorm.io/gorm"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB          *gorm.DB
	AuthService *services.AuthService
}