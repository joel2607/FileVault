package storage

// FileStorageProvider defines the interface for a file storage backend.
// This abstraction allows for interchangeable storage solutions (e.g., local, S3)
// without changing the core business logic.
type FileStorageProvider interface {
	// GetDownloadURL generates a temporary, secure URL to access a file.
	// - filePath: The path of the file within the storage backend (e.g., "user_1/data.txt").
	// - originalFilename: The user-facing name for the file, used to set the
	//   Content-Disposition header for the download.
	GetDownloadURL(filePath string, originalFilename string) (string, error)
}
