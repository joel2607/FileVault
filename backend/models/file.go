package models

import "gorm.io/gorm"

// File corresponds to the "files" table in the database
type File struct {
	gorm.Model // Includes ID, CreatedAt, UpdatedAt, DeletedAt

	Filename      string `gorm:"not null"`
	Mimetype      string `gorm:"not null"`
	Size          int64  `gorm:"not null"`
	Hash          string `gorm:"unique;not null"`
	DownloadCount int    `gorm:"default:0;not null"`
	IsPublic      bool   `gorm:"default:false;not null"`

	OwnerID uint // Foreign key for the User who owns the file
	Owner   User `gorm:"references:ID"`

	FolderID *uint  // Foreign key for the Folder (nullable)
	Folder   Folder `gorm:"references:ID"`
}
