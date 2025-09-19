package services

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strconv"
	"net/http"
	"github.com/99designs/gqlgen/graphql"
	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-joel2607/models"
	"gorm.io/gorm"
)

type FileService struct {
	DB *gorm.DB
}

func NewFileService(db *gorm.DB) *FileService {
	return &FileService{DB: db}
}

func (s *FileService) UploadFile(ctx context.Context, file graphql.Upload, user *models.User, parentFolderID *string) (*models.File, error) {
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
	if err := s.DB.Create(newContent).Error; err != nil {
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

	return newFile, nil
}

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

func (s *FileService) DeleteFolder(ctx context.Context, id string, user *models.User) (bool, error) {
	uid, _ := strconv.ParseUint(id, 10, 64)
	err := s.DB.Delete(&models.Folder{}, "id = ? AND user_id = ?", uid, user.ID).Error
	return err == nil, err
}

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

func (s *FileService) DeleteFile(ctx context.Context, id string, user *models.User) (bool, error) {
	uid, _ := strconv.ParseUint(id, 10, 64)
	var file models.File
	if err := s.DB.First(&file, "id = ? AND user_id = ?", uid, user.ID).Error; err != nil {
		return false, err
	}

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
		}
	}

	err := s.DB.Delete(&file).Error
	return err == nil, err
}

func (s *FileService) GetFolder(ctx context.Context, id string, user *models.User) (*models.Folder, error) {
	var folder models.Folder
	uid, _ := strconv.ParseUint(id, 10, 64)
	err := s.DB.First(&folder, "id = ? AND user_id = ?", uid, user.ID).Error
	return &folder, err
}

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
