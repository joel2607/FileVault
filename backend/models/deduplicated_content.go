// Package models defines the data structures used in the application.
package models

// DeduplicatedContent stores a single instance of file content for deduplication.
// This table is central to the deduplication feature, tracking the hash of the content
// and the number of files that reference it.
type DeduplicatedContent struct {
	BaseModel
	SHA256Hash     string `gorm:"type:varchar(128);unique;not null"`
	ReferenceCount int    `gorm:"default:0"`
}