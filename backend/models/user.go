package models

import "gorm.io/gorm"

// User corresponds to the "users" table in the database
type User struct {
	gorm.Model // Includes ID, CreatedAt, UpdatedAt, DeletedAt

	Username     string `gorm:"unique;not null"`
	IsAdmin      bool   `gorm:"default:false;not null"`
	PasswordHash string `gorm:"not null"`
	StorageUsed  int64  `gorm:"default:0;not null"`
	StorageQuota int64  `gorm:"default:10485760;not null"` // 10 MB default

	Files   []File   `gorm:"foreignKey:OwnerID"` // A User has many Files
	Folders []Folder `gorm:"foreignKey:OwnerID"` // A User has many Folders
}
