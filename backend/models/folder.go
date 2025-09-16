// Package models defines the data structures used in the application.
package models

import "gorm.io/gorm"

// Folder represents a folder for organizing files.
// This table supports nested folders and public sharing of folders.
type Folder struct {
	gorm.Model
	UserID         uint   `gorm:"not null"`
	User           User   `gorm:"foreignkey:UserID"`
	FolderName     string `gorm:"type:varchar(255);not null"`
	ParentFolderID *uint  `gorm:"default:null"`
	IsPublic       bool   `gorm:"default:false"`
}