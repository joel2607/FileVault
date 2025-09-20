package models

// FolderSharing manages folder sharing with specific users.
type FolderSharing struct {
	BaseModel
	FolderID         uint   `gorm:"not null"`
	Folder           Folder `gorm:"foreignkey:FolderID"`
	SharedWithUserID uint   `gorm:"not null"`
	SharedWithUser   User   `gorm:"foreignkey:SharedWithUserID"`
	PermissionLevel  string `gorm:"type:varchar(50)"`
}