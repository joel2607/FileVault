package services

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-joel2607/database"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-joel2607/models"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-joel2607/services/storage"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// FileService provides methods for file management, including uploads,
// deletions, and storage statistics.
type FileService struct {
	DB      *gorm.DB
	RDB     *redis.Client
	Storage storage.FileStorageProvider
}

// NewFileService creates a new instance of FileService.
func NewFileService(db *gorm.DB, rdb *redis.Client, storage storage.FileStorageProvider) *FileService {
	return &FileService{DB: db, RDB: rdb, Storage: storage}
}

func (s *FileService) GetStorageStatistics(userID uint) (*models.StorageStatistics, error) {
	var user models.User
	if err := s.DB.First(&user, userID).Error; err != nil {
		return nil, err
	}

	var percentageSaved float64
	if user.UsedStorageKB > 0 {
		percentageSaved = (user.SavedStorageKB / user.UsedStorageKB) * 100
	}

	return &models.StorageStatistics{
		UsedStorageKb:   user.UsedStorageKB,
		SavedStorageKb:  user.SavedStorageKB,
		PercentageSaved: percentageSaved,
	}, nil
}

func (s *FileService) publishStorageUpdate(user *models.User) {
	stats, err := s.GetStorageStatistics(user.ID)
	if err != nil {
		log.Printf("Error getting storage statistics for user %d: %v", user.ID, err)
		return
	}

	payload, err := json.Marshal(stats)
	if err != nil {
		log.Printf("Error marshalling storage statistics for user %d: %v", user.ID, err)
		return
	}

	channel := fmt.Sprintf("storage_updates_%d", user.ID)
	s.RDB.Publish(database.Ctx, channel, payload)
}

// GenerateDownloadURL handles the logic for creating a secure, temporary download link for a file.
// It checks for user permissions, increments the file's download count, and then delegates
// the actual URL creation to the configured storage provider.
func (s *FileService) GenerateDownloadURL(ctx context.Context, fileID string, user *models.User) (string, error) {
	id, err := strconv.ParseUint(fileID, 10, 64)
	if err != nil {
		return "", fmt.Errorf("invalid file ID")
	}
	uid := uint(id)

	var file models.File
	if err := s.DB.First(&file, uid).Error; err != nil {
		return "", fmt.Errorf("file not found")
	}

	// Authorization check: user must be the owner, the file must be public, or it must be shared with the user.
	isOwner := file.UserID == user.ID
	isPublic := file.IsPublic
	var isShared bool
	var fileShare models.FileSharing
	if err := s.DB.Where("file_id = ? AND shared_with_user_id = ?", file.ID, user.ID).First(&fileShare).Error; err == nil {
		isShared = true
	}

	if !isOwner && !isPublic && !isShared {
		return "", fmt.Errorf("access denied: you do not have permission to download this file")
	}

	// Atomically increment the download count
	if err := s.DB.Model(&models.File{}).Where("id = ?", file.ID).UpdateColumn("download_count", gorm.Expr("download_count + 1")).Error; err != nil {
		log.Printf("Failed to increment download count for file %d: %v", file.ID, err)
		// Do not fail the operation if this fails, just log it.
	} else {
		// Publish the update to Redis
		s.PublishDownloadCountUpdate(file.ID, int32(file.DownloadCount+1))
	}

	var content models.DeduplicatedContent
	if err := s.DB.First(&content, file.DeduplicationID).Error; err != nil {
		return "", fmt.Errorf("could not find file content")
	}

	// Delegate URL generation to the storage provider
	return s.Storage.GetDownloadURL(content.SHA256Hash, file.FileName)
}

// PublishDownloadCountUpdate publishes a message to Redis when a file's download count changes.
func (s *FileService) PublishDownloadCountUpdate(fileID uint, newCount int32) {
	payload, err := json.Marshal(models.DownloadCountUpdate{
		FileID:        strconv.FormatUint(uint64(fileID), 10),
		DownloadCount: int32(newCount),
	})
	if err != nil {
		log.Printf("Error marshalling download count update for file %d: %v", fileID, err)
		return
	}

	channel := fmt.Sprintf("download_updates_%d", fileID)
	if err := s.RDB.Publish(database.Ctx, channel, payload).Err(); err != nil {
		log.Printf("Error publishing download count update for file %d: %v", fileID, err)
	}
}

// SubscribeToFileDownloads handles the business logic for subscribing to file download updates.
// It includes authorization checks and Redis pub/sub management.
func (s *FileService) SubscribeToFileDownloads(ctx context.Context, fileID string, user *models.User) (<-chan *models.DownloadCountUpdate, error) {
	id, err := strconv.ParseUint(fileID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid file ID")
	}
	uid := uint(id)

	var file models.File
	if err := s.DB.First(&file, uid).Error; err != nil {
		return nil, fmt.Errorf("file not found")
	}

	// Authorization check
	isOwner := file.UserID == user.ID
	isPublic := file.IsPublic
	var isShared bool
	var fileShare models.FileSharing
	if err := s.DB.Where("file_id = ? AND shared_with_user_id = ?", file.ID, user.ID).First(&fileShare).Error; err == nil {
		isShared = true
	}
	isAdmin := user.Role == "ADMIN"

	if !isOwner && !isPublic && !isShared && !isAdmin {
		return nil, fmt.Errorf("access denied")
	}

	// Create channel and subscribe
	ch := make(chan *models.DownloadCountUpdate, 1)
	channel := fmt.Sprintf("download_updates_%d", file.ID)
	pubsub := s.RDB.Subscribe(ctx, channel)

	// Send initial data
	ch <- &models.DownloadCountUpdate{
		FileID:        fileID,
		DownloadCount: int32(file.DownloadCount),
	}

	// Goroutine to listen for updates
	go func() {
		defer pubsub.Close()
		defer close(ch)

		for {
			select {
			case msg, ok := <-pubsub.Channel():
				if !ok {
					return
				}
				var update models.DownloadCountUpdate
				if err := json.Unmarshal([]byte(msg.Payload), &update); err == nil {
					ch <- &update
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return ch, nil
}

// SubscribeToStorageStatistics handles the business logic for subscribing to storage statistics updates.
func (s *FileService) SubscribeToStorageStatistics(ctx context.Context, userID *string, currentUser *models.User) (<-chan *models.StorageStatistics, error) {
	var targetUserID uint
	if userID != nil {
		if currentUser.Role != "ADMIN" {
			return nil, fmt.Errorf("only admins can view other users' statistics")
		}
		id, err := strconv.ParseUint(*userID, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid user ID")
		}
		targetUserID = uint(id)
	} else {
		targetUserID = currentUser.ID
	}

	ch := make(chan *models.StorageStatistics, 1)
	channel := fmt.Sprintf("storage_updates_%d", targetUserID)
	pubsub := s.RDB.Subscribe(ctx, channel)

	// Send initial data
	initialStats, err := s.GetStorageStatistics(targetUserID)
	if err == nil {
		ch <- initialStats
	}

	go func() {
		defer pubsub.Close()
		defer close(ch)

		for {
			select {
			case msg, ok := <-pubsub.Channel():
				if !ok {
					return
				}
				var stats models.StorageStatistics
				if err := json.Unmarshal([]byte(msg.Payload), &stats); err == nil {
					ch <- &stats
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return ch, nil
}

// UploadFile handles the entire file upload process.
// It performs content hashing for deduplication, saves the file to storage if it's new,
// validates the MIME type, and creates the necessary metadata in the database.
// It also updates the user's storage usage statistics.
//
// Inputs:
// - ctx: The context for the request.
// - file: The graphql.Upload object containing the file data and metadata.
// - user: The user model of the person uploading the file.
// - parentFolderID: An optional string pointer to the ID of the folder where the file should be placed.
//
// Outputs:
// - A pointer to the created models.File object if successful.
// - An error if any part of the process fails.
func (s *FileService) UploadFile(ctx context.Context, file graphql.Upload, user *models.User, parentFolderID *string) (*models.File, error) {

	// Check for whether the user has enough storage quota
	if user.UsedStorageKB - user.SavedStorageKB + float64(file.Size)/1024 > user.StorageQuotaKB {
		return nil, fmt.Errorf("storage quota exceeded")
	}


	log.Printf("Uploading file: %s of size %d bytes\n", file.Filename, file.Size)
	// 1. Hashing
	hash := sha256.New()
	if _, err := io.Copy(hash, file.File); err != nil {
		return nil, err
	}
	file.File.Seek(0, 0) // Reset reader
	sha256Hash := fmt.Sprintf("%x", hash.Sum(nil))

	// 2. Deduplication Check
	var existingContent models.DeduplicatedContent
	if err := s.DB.Where("sha256_hash = ?", sha256Hash).First(&existingContent).Error; err == nil {
		// Content exists, create new file metadata and point to existing content
		newFile := &models.File{
			UserID:          user.ID,
			FileName:        file.Filename,
			MIMEType:        file.ContentType,
			Size:            file.Size,
			DeduplicationID: existingContent.ID,
		}
		if parentFolderID != nil {
			id, _ := strconv.ParseUint(*parentFolderID, 10, 64)
			uid := uint(id)
			newFile.FolderID = &uid
		}
		if err := s.DB.Create(newFile).Error; err != nil {
			return nil, err
		}
		existingContent.ReferenceCount++
		s.DB.Save(&existingContent)

		// Update user's storage usage. They save space by not uploading duplicate data.
		storageChangeKB := float64(file.Size) / 1024
		user.SavedStorageKB += storageChangeKB
		user.UsedStorageKB += storageChangeKB
		log.Printf("User used storage(saved): %f", user.UsedStorageKB)
		s.DB.Save(user)
		s.publishStorageUpdate(user)
		return newFile, nil
	}

	// 3. Save File and Create Metadata
	// For simplicity, saving to a local 'uploads' directory. In production, use a cloud storage service.
	uploadDir := "./uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, os.ModePerm)
	}
	filePath := filepath.Join(uploadDir, sha256Hash)
	out, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer out.Close()
	if _, err := io.Copy(out, file.File); err != nil {
		return nil, err
	}

	// 4. MIME Type Validation
	if !isValidMIME(filePath, file.ContentType) {
		os.Remove(filePath) // Clean up invalid file
		return nil, fmt.Errorf("invalid MIME type")
	}

	// Create new deduplicated content and file metadata
	newContent := &models.DeduplicatedContent{
		SHA256Hash:     sha256Hash,
		ReferenceCount: 1,
	}
	log.Println("Creating new deduplicated content with hash:", sha256Hash)
	if err := s.DB.Create(newContent).Error; err != nil {
		log.Printf("Could not create deduplicated content: %v", err)
		return nil, err
	}
	newFile := &models.File{
		UserID:          user.ID,
		FileName:        file.Filename,
		MIMEType:        file.ContentType,
		Size:            file.Size,
		DeduplicationID: newContent.ID,
	}
	if parentFolderID != nil {
		id, _ := strconv.ParseUint(*parentFolderID, 10, 64)
		uid := uint(id)
		newFile.FolderID = &uid
	}
	if err := s.DB.Create(newFile).Error; err != nil {
		return nil, err
	}

	// Update user's storage usage. They dont save space as this is new data.
	user.UsedStorageKB += float64(file.Size) / 1024
	log.Printf("User used storage: %f", user.UsedStorageKB)
	s.DB.Save(user)
	s.publishStorageUpdate(user)

	return newFile, nil
}

// CreateFolder creates a new folder for a given user.
//
// Inputs:
// - ctx: The context for the request.
// - input: The NewFolder input object containing the folder's name and optional parent ID.
// - user: The user who is creating the folder.
//
// Outputs:
// - A pointer to the created models.Folder object.
// - An error if the database operation fails.
func (s *FileService) CreateFolder(ctx context.Context, input models.NewFolder, user *models.User) (*models.Folder, error) {
	folder := &models.Folder{
		UserID:     user.ID,
		FolderName: input.FolderName,
	}
	if input.ParentFolderID != nil {
		id, _ := strconv.ParseUint(*input.ParentFolderID, 10, 64)
		uid := uint(id)
		folder.ParentFolderID = &uid
	}
	err := s.DB.Create(folder).Error
	return folder, err
}

// UpdateFile modifies an existing file's metadata, such as its name or parent folder.
// It ensures that the user attempting the update is the owner of the file.
//
// Inputs:
// - ctx: The context for the request.
// - input: The UpdateFile input object containing the file's ID and the new data.
// - user: The user requesting the update.
//
// Outputs:
// - A pointer to the updated models.File object.
// - An error if the file is not found or the database operation fails.
func (s *FileService) UpdateFile(ctx context.Context, input models.UpdateFile, user *models.User) (*models.File, error) {
	uid, err := strconv.ParseUint(input.ID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid file ID")
	}
	var file models.File
	if err := s.DB.First(&file, "id = ? AND user_id = ?", uid, user.ID).Error; err != nil {
		return nil, err
	}
	if input.FileName != nil {
		file.FileName = *input.FileName
	}
	if input.ParentFolderID != nil {
		id, _ := strconv.ParseUint(*input.ParentFolderID, 10, 64)
		parsedID := uint(id)
		file.FolderID = &parsedID
	}
	err = s.DB.Save(&file).Error
	return &file, err
}

// DeleteFile removes a file's metadata from the database and handles the deduplication logic.
// It decrements the reference count of the associated content. If the count reaches zero,
// it deletes the actual file from storage, the content record from the database and
// updates the user's storage usage statistics.
//
// Inputs:
// - ctx: The context for the request.
// - id: The string ID of the file to be deleted.
// - user: The user requesting the deletion.
//
// Outputs:
// - A boolean indicating whether the deletion was successful.
// - An error if the database operation fails.
func (s *FileService) DeleteFile(ctx context.Context, id string, user *models.User) (*models.File, error) {
	uid, _ := strconv.ParseUint(id, 10, 64)
	var file models.File
	if err := s.DB.First(&file, "id = ? AND user_id = ?", uid, user.ID).Error; err != nil {
		return nil, err
	}

	fileSizeKB := float64(file.Size) / 1024

	// First delete file record to prevent deduplicated content delete fail.
	err := s.DB.Delete(&file).Error
	if err != nil {
		return nil, err
	}
	// Decrement reference count
	var content models.DeduplicatedContent
	if err := s.DB.First(&content, file.DeduplicationID).Error; err == nil {
		content.ReferenceCount--
		s.DB.Save(&content)
		// Update user's storage usage
		user.UsedStorageKB -= fileSizeKB
		s.DB.Save(user)
		if content.ReferenceCount <= 0 {
			// Delete file from storage
			if err := s.DB.Delete(&content).Error; err != nil {
				log.Println(err)
				log.Println("Failed to delete deduplicated content with hash:", content.SHA256Hash)
			}
			filePath := filepath.Join("./uploads", content.SHA256Hash)
			os.Remove(filePath)
		} else {
			// Update user's storage usage. They save less space now as one less reference.
			user.SavedStorageKB -= fileSizeKB
			s.DB.Save(user)
		}
	}

	s.publishStorageUpdate(user)
	return &file, nil

}

// UpdateFolder modifies an existing folder's properties, such as its name or parent folder.
// It ensures that the user attempting the update is the owner of the folder.
//
// Inputs:
// - ctx: The context for the request.
// - input: The UpdateFolder input object containing the folder's ID and the new data.
// - user: The user requesting the update.
//
// Outputs:
// - A pointer to the updated models.Folder object.
// - An error if the folder is not found or the database operation fails.
func (s *FileService) UpdateFolder(ctx context.Context, input models.UpdateFolder, user *models.User) (*models.Folder, error) {
	uid, err := strconv.ParseUint(input.ID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid folder ID")
	}
	var folder models.Folder
	if err := s.DB.First(&folder, "id = ? AND user_id = ?", uid, user.ID).Error; err != nil {
		return nil, err
	}
	if input.FolderName != nil {
		folder.FolderName = *input.FolderName
	}
	if input.ParentFolderID != nil {
		id, _ := strconv.ParseUint(*input.ParentFolderID, 10, 64)
		parsedID := uint(id)
		folder.ParentFolderID = &parsedID
	}
	err = s.DB.Save(&folder).Error
	return &folder, err
}

// DeleteFolder removes a folder from the database.
// It ensures that the user attempting the deletion is the owner of the folder.
// Note: This is a simple implementation. A production system would need to handle orphaned files or subfolders.
//
// Inputs:
// - ctx: The context for the request.
// - id: The string ID of the folder to be deleted.
// - user: The user requesting the deletion.
//
// Outputs:
// - A boolean indicating whether the deletion was successful.
// - An error if the database operation fails.
func (s *FileService) DeleteFolder(ctx context.Context, id string, user *models.User) (*models.Folder, error) {
	uid, _ := strconv.ParseUint(id, 10, 64)
	var folder models.Folder
	if err := s.DB.First(&folder, "id = ? AND user_id = ?", uid, user.ID).Error; err != nil {
		return nil, err
	}

	// Find all files in the folder
	var files []*models.File
	if err := s.DB.Where("folder_id = ?", uid).Find(&files).Error; err != nil {
		return nil, err
	}

	// Delete all files in the folder
	for _, file := range files {
		fileID := strconv.FormatUint(uint64(file.ID), 10)
		if _, err := s.DeleteFile(ctx, fileID, user); err != nil {
			return nil, err
		}
	}

	// Find all subfolders in the folder
	var subfolders []*models.Folder
	if err := s.DB.Where("parent_folder_id = ?", uid).Find(&subfolders).Error; err != nil {
		return nil, err
	}

	// Recursively delete all subfolders
	for _, subfolder := range subfolders {
		subfolderID := strconv.FormatUint(uint64(subfolder.ID), 10)
		if _, err := s.DeleteFolder(ctx, subfolderID, user); err != nil {
			return nil, err
		}
	}

	err := s.DB.Delete(&folder).Error
	if err != nil {
		return nil, err
	}
	return &folder, nil
}

// isValidMIME validates the actual content type of a file against its declared MIME type.
// It reads the first 512 bytes of the file to determine the real MIME type and also checks
// the file extension as a fallback.
//
// Inputs:
// - filePath: The path to the saved file on the local filesystem.
// - declaredMIME: The MIME type that was declared by the client upon upload.
//
// Outputs:
// - A boolean that is true if the actual MIME type matches the declared one, and false otherwise.
func isValidMIME(filePath, declaredMIME string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	// Read the first 512 bytes to let http.DetectContentType determine the MIME type.
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil && err != io.EOF {
		return false
	}

	// Get the actual MIME type from the content.
	actualMIME := http.DetectContentType(buffer)

	// It's common for http.DetectContentType to return a generic MIME type.
	// We can also check the file extension.
	ext := filepath.Ext(filePath)
	extMIME := mime.TypeByExtension(ext)

	return declaredMIME == actualMIME || (extMIME != "" && declaredMIME == extMIME)
}
