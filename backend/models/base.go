package models

import (
	"time"
)

// BaseModel defines the common fields for all models, without the gorm.DeletedAt field.
// This ensures that all deletes are hard deletes.
type BaseModel struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
