package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-joel2607/models"
	"gorm.io/gorm"
)

// ShareService handles business logic related to file and folder sharing.
type ShareService struct {
	DB *gorm.DB
}

// NewShareService creates and returns a new ShareService instance.
func NewShareService(db *gorm.DB) *ShareService {
	return &ShareService{DB: db}
}

// GetUserRoot retrieves the top-level files and folders for a regular user.
// It returns the user's own root-level items plus any public root-level items.
func (s *ShareService) GetUserRoot(ctx context.Context, user *models.User) (*models.Root, error) {
	var files []*models.File
	var folders []*models.Folder

	if err := s.DB.Where("(user_id = ? AND folder_id IS NULL)", user.ID, true).Find(&files).Error; err != nil {
		return nil, err
	}
	if err := s.DB.Where("(user_id = ? AND parent_folder_id IS NULL)", user.ID, true).Find(&folders).Error; err != nil {
		return nil, err
	}

	return &models.Root{Files: files, Folders: folders}, nil
}

// GetAdminRoot retrieves root-level files and folders for an admin user.
// If a userID is provided, it fetches the root for that specific user.
// Otherwise, it returns all root-level items in the system.
func (s *ShareService) GetAdminRoot(ctx context.Context, userID *string) (*models.Root, error) {
	var files []*models.File
	var folders []*models.Folder

	db := s.DB
	if userID != nil {
		uid, err := strconv.ParseUint(*userID, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid user ID")
		}
		db = s.DB.Where("user_id = ?", uid)
	}

	if err := db.Where("folder_id IS NULL").Find(&files).Error; err != nil {
		return nil, err
	}
	if err := db.Where("parent_folder_id IS NULL").Find(&folders).Error; err != nil {
		return nil, err
	}

	return &models.Root{Files: files, Folders: folders}, nil
}

// GetFolder retrieves a specific folder by its ID, enforcing access control.
// Admins can access any folder.
// Regular users can access folders they own or public folders.
func (s *ShareService) GetFolder(ctx context.Context, id string, user *models.User) (*models.Folder, error) {
	var folder models.Folder
	uid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid folder ID")
	}

	if err := s.DB.First(&folder, uid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("folder not found")
		}
		return nil, err
	}

	// Admins can access any folder
	if user.Role == models.RoleAdmin {
		return &folder, nil
	}

	// Public folders are accessible to anyone
	if folder.IsPublic {
		return &folder, nil
	}

	// Users can access their own folders
	if folder.UserID == user.ID {
		return &folder, nil
	}

	return nil, fmt.Errorf("access denied")
}

// GetFile retrieves a specific file by its ID, enforcing access control.
// Admins can access any file.
// Regular users can access files they own, public files, or files shared with them.
func (s *ShareService) GetFile(ctx context.Context, id string, user *models.User) (*models.File, error) {
	var file models.File
	uid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid file ID")
	}

	if err := s.DB.First(&file, uid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("file not found")
		}
		return nil, err
	}

	// Admins can access any file
	if user.Role == models.RoleAdmin {
		return &file, nil
	}

	// Public files are accessible to anyone
	if file.IsPublic {
		return &file, nil
	}

	// Users can access their own files
	if file.UserID == user.ID {
		return &file, nil
	}

	// Check if the file is explicitly shared with the user
	var shareRecord models.FileSharing
	if err := s.DB.Where("file_id = ? AND shared_with_user_id = ?", file.ID, user.ID).First(&shareRecord).Error; err == nil {
		return &file, nil
	}

	return nil, fmt.Errorf("access denied")
}

// SetFilePublic makes a file public. Only the file owner can perform this action.
func (s *ShareService) SetFilePublic(ctx context.Context, fileID string, user *models.User) (*models.File, error) {
	var file models.File
	uid, err := strconv.ParseUint(fileID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid file ID")
	}

	if err := s.DB.First(&file, "id = ? AND user_id = ?", uid, user.ID).Error; err != nil {
		return nil, fmt.Errorf("file not found or access denied")
	}

	file.IsPublic = true
	if err := s.DB.Save(&file).Error; err != nil {
		return nil, err
	}

	return &file, nil
}

// SetFilePrivate makes a file private. Only the file owner can perform this action.
func (s *ShareService) SetFilePrivate(ctx context.Context, fileID string, user *models.User) (*models.File, error) {
	var file models.File
	uid, err := strconv.ParseUint(fileID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid file ID")
	}

	if err := s.DB.First(&file, "id = ? AND user_id = ?", uid, user.ID).Error; err != nil {
		return nil, fmt.Errorf("file not found or access denied")
	}

	file.IsPublic = false
	if err := s.DB.Save(&file).Error; err != nil {
		return nil, err
	}

	return &file, nil
}

// ShareFileWithUser grants a user access to a private file.
// It returns an error if the file is public.
func (s *ShareService) ShareFileWithUser(ctx context.Context, fileID string, shareWithUserID string, user *models.User) (*models.FileSharing, error) {
	var file models.File
	uid, err := strconv.ParseUint(fileID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid file ID")
	}

	if err := s.DB.First(&file, "id = ? AND user_id = ?", uid, user.ID).Error; err != nil {
		return nil, fmt.Errorf("file not found or access denied")
	}

	if file.IsPublic {
		return nil, fmt.Errorf("cannot share a public file")
	}

	shareWithUID, err := strconv.ParseUint(shareWithUserID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID to share with")
	}

	share := &models.FileSharing{
		FileID:           uint(uid),
		SharedWithUserID: uint(shareWithUID),
		PermissionLevel:  "read", // Or any other permission level you want to implement
	}

	if err := s.DB.Create(share).Error; err != nil {
		return nil, err
	}

	return share, nil
}

// RemoveFileAccess removes a user's access to a shared file.
// Only the file owner can perform this action.
func (s *ShareService) RemoveFileAccess(ctx context.Context, fileID string, userIDToRemove string, user *models.User) (bool, error) {
	var file models.File
	uid, err := strconv.ParseUint(fileID, 10, 64)
	if err != nil {
		return false, fmt.Errorf("invalid file ID")
	}

	if err := s.DB.First(&file, "id = ? AND user_id = ?", uid, user.ID).Error; err != nil {
		return false, fmt.Errorf("file not found or access denied")
	}

	userIDToRemove_uid, err := strconv.ParseUint(userIDToRemove, 10, 64)
	if err != nil {
		return false, fmt.Errorf("invalid user ID to remove access from")
	}

	if err := s.DB.Where("file_id = ? AND shared_with_user_id = ?", uid, userIDToRemove_uid).Delete(&models.FileSharing{}).Error; err != nil {
		return false, err
	}

	return true, nil
}