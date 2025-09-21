package storage

import (
	"fmt"
	"net/url"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

// LocalStorageProvider implements the FileStorageProvider interface for serving
// files from the local filesystem.
type LocalStorageProvider struct {
	BaseURL string // The base URL of the server, e.g., "http://localhost:8080"
}

// NewLocalStorageProvider creates a new instance of LocalStorageProvider.
func NewLocalStorageProvider(baseURL string) *LocalStorageProvider {
	return &LocalStorageProvider{BaseURL: baseURL}
}

// GetDownloadURL generates a secure, temporary URL for a local file.
// It creates a short-lived JWT that encodes the file path, which is then
// validated by a dedicated download handler.
func (p *LocalStorageProvider) GetDownloadURL(filePath string, originalFilename string) (string, error) {
	// Create a short-lived JWT to authorize the download.
	// This token is specific to this file and expires quickly.
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &jwt.RegisteredClaims{
		Subject:   filePath,
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}

	secret := []byte(viper.GetString("DOWNLOAD_TOKEN_SECRET"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign download token: %w", err)
	}

	// Construct the final download URL pointing to our secure download endpoint.
	return fmt.Sprintf(
		"%s/downloads/%s?token=%s&filename=%s",
		p.BaseURL,
		url.PathEscape(filePath),
		tokenString,
		url.QueryEscape(originalFilename),
	), nil
}
