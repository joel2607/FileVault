package models

import "gorm.io/gorm"

// Folder corresponds to the "folders" table in the database
type Folder struct {
	gorm.Model // Includes ID, CreatedAt, UpdatedAt, DeletedAt

	Name     string `gorm:"not null"`
	IsPublic bool   `gorm:"default:false;not null"`

	OwnerID uint // Foreign key for the User who owns the folder
	Owner   User `gorm:"references:ID"`

	ParentID *uint    // Foreign key for the parent folder (nullable for root folders)
	Parent   *Folder  `gorm:"references:ID"`
	Subfolders []Folder `gorm:"foreignKey:ParentID"` // A Folder has many subfolders

	Files []File `gorm:"foreignKey:FolderID"` // A Folder has many Files
}
