// Package models defines the data structures used in the application.
package models

import (
	"gorm.io/gorm"
)

// File represents a file uploaded by a user.
// This table stores metadata for each file, but not the file content itself.
// It links to the user who uploaded it and the deduplicated content.
type File struct {
	gorm.Model
	UserID              uint      `gorm:"not null"`
	User                User      `gorm:"foreignkey:UserID"`
	FileName            string    `gorm:"type:varchar(255);not null"`
	MIMEType            string    `gorm:"type:varchar(100);not null"`
	Size                int64     `gorm:"not null"`
	DeduplicationID     uint      `gorm:"not null"`
	DeduplicatedContent DeduplicatedContent `gorm:"foreignkey:DeduplicationID"`
	IsPublic            bool      `gorm:"default:false"`
	DownloadCount       int       `gorm:"default:0"`
	Tags                string    `gorm:"type:jsonb"`
	FolderID            *uint     `gorm:"default:null"`
	Folder              *Folder   `gorm:"foreignkey:FolderID"`
}
