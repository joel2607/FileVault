// Package models defines the data structures used in the application.
package models

import "gorm.io/gorm"

// FileSharing manages file sharing with specific users.
// This table is used for the optional feature of sharing files with
// specific users and defining their permission levels.
type FileSharing struct {
	gorm.Model
	FileID           uint   `gorm:"not null"`
	File             File   `gorm:"foreignkey:FileID"`
	SharedWithUserID uint   `gorm:"not null"`
	SharedWithUser   User   `gorm:"foreignkey:SharedWithUserID"`
	PermissionLevel  string `gorm:"type:varchar(50)"`
}