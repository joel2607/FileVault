// Package models defines the data structures used in the application.
package models

import (
	"gorm.io/gorm"
)

// UserRole defines the roles a user can have.
type UserRole string

const (
	// RoleUser is the default role for a regular user.
	RoleUser UserRole = "user"
	// RoleAdmin is the role for an administrator.
	RoleAdmin UserRole = "admin"
)

// User represents a user in the system.
// This table stores user information, including authentication details,
// storage quotas, and API rate limits.
type User struct {
	gorm.Model
	Username       string    `gorm:"type:varchar(255);unique;not null"`
	Email          string    `gorm:"type:varchar(255);unique;not null"`
	PasswordHash   string    `gorm:"type:varchar(255);not null"`
	StorageQuotaKB float64   `gorm:"default:10240"`
	UsedStorageKB  float64   `gorm:"default:0"`
	SavedStorageKB float64   `gorm:"default:0"`
	APIRateLimit   int       `gorm:"default:2"`
	Role           UserRole  `gorm:"type:varchar(50);default:'user'"`
}