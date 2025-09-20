package services

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-joel2607/models"
	"gorm.io/gorm"
)

// FileService handles business logic related to file and folder management.
// It interacts with the database to perform CRUD operations on files and folders.
type FileService struct {
	DB *gorm.DB
}

// NewFileService creates and returns a new FileService instance.
// It takes a GORM database connection as input.
// Returns a pointer to the newly created FileService.
func NewFileService(db *gorm.DB) *FileService {
	return &FileService{DB: db}
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
		FolderName: input.Name,
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
	var file models.File
	if err := s.DB.First(&file, "id = ? AND user_id = ?", input.ID, user.ID).Error; err != nil {
		return nil, err
	}
	if input.Name != nil {
		file.FileName = *input.Name
	}
	if input.ParentFolderID != nil {
		id, _ := strconv.ParseUint(*input.ParentFolderID, 10, 64)
		uid := uint(id)
		file.FolderID = &uid
	}
	err := s.DB.Save(&file).Error
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
func (s *FileService) DeleteFile(ctx context.Context, id string, user *models.User) (bool, error) {
	uid, _ := strconv.ParseUint(id, 10, 64)
	var file models.File
	if err := s.DB.First(&file, "id = ? AND user_id = ?", uid, user.ID).Error; err != nil {
		return false, err
	}

	storageChangeKB := float64(file.Size) / 1024

	// Decrement reference count
	var content models.DeduplicatedContent
	if err := s.DB.First(&content, file.DeduplicationID).Error; err == nil {
		content.ReferenceCount--
		s.DB.Save(&content)
		if content.ReferenceCount <= 0 {
			// Delete file from storage
			filePath := filepath.Join("./uploads", content.SHA256Hash)
			os.Remove(filePath)
			s.DB.Delete(&content)
			// Update user's storage usage
			user.UsedStorageKB -= storageChangeKB
			s.DB.Save(user)
		} else {
			// Update user's storage usage. They save less space now as one less reference.
			user.SavedStorageKB -= storageChangeKB
			user.UsedStorageKB -= storageChangeKB
			s.DB.Save(user)
		}
	}

	err := s.DB.Delete(&file).Error
	return err == nil, err
}

// GetFolder retrieves a specific folder by its ID for a given user.
// It ensures that only the owner of the folder can access it.
//
// Inputs:
// - ctx: The context for the request.
// - id: The string ID of the folder to retrieve.
// - user: The user requesting the folder.
//
// Outputs:
// - A pointer to the retrieved models.Folder object.
// - An error if the folder is not found or the user does not have permission.
func (s *FileService) GetFolder(ctx context.Context, id string, user *models.User) (*models.Folder, error) {
	var folder models.Folder
	uid, _ := strconv.ParseUint(id, 10, 64)
	err := s.DB.First(&folder, "id = ? AND user_id = ?", uid, user.ID).Error
	return &folder, err
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
	var folder models.Folder
	if err := s.DB.First(&folder, "id = ? AND user_id = ?", input.ID, user.ID).Error; err != nil {
		return nil, err
	}
	if input.Name != nil {
		folder.FolderName = *input.Name
	}
	if input.ParentFolderID != nil {
		id, _ := strconv.ParseUint(*input.ParentFolderID, 10, 64)
		uid := uint(id)
		folder.ParentFolderID = &uid
	}
	err := s.DB.Save(&folder).Error
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
func (s *FileService) DeleteFolder(ctx context.Context, id string, user *models.User) (bool, error) {
	uid, _ := strconv.ParseUint(id, 10, 64)
	err := s.DB.Delete(&models.Folder{}, "id = ? AND user_id = ?", uid, user.ID).Error
	return err == nil, err
}

// GetRoot retrieves the top-level files and folders for a given user.
// This represents the user's root directory.
//
// Inputs:
// - ctx: The context for the request.
// - user: The user whose root directory is being requested.
//
// Outputs:
// - A pointer to a models.Root object, containing slices of top-level files and folders.
// - An error if the database query fails.
func (s *FileService) GetRoot(ctx context.Context, user *models.User) (*models.Root, error) {
	var files []*models.File
	if err := s.DB.Where("user_id = ? AND folder_id IS NULL", user.ID).Find(&files).Error; err != nil {
		return nil, err
	}
	var folders []*models.Folder
	if err := s.DB.Where("user_id = ? AND parent_folder_id IS NULL", user.ID).Find(&folders).Error; err != nil {
		return nil, err
	}
	return &models.Root{Files: files, Folders: folders}, nil
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