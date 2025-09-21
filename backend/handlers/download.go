package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

// DownloadHandler handles requests for secure file downloads.
// It validates a short-lived JWT from the query parameters to authorize the request,
// then serves the file from the local 'uploads' directory.
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.URL.Query().Get("token")
	filePath := chi.URLParam(r, "*")
	filename := r.URL.Query().Get("filename")

	if tokenStr == "" || filePath == "" || filename == "" {
		http.Error(w, "Forbidden: Missing required parameters", http.StatusForbidden)
		return
	}

	// Validate the temporary download token.
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(viper.GetString("DOWNLOAD_TOKEN_SECRET")), nil
	})

	if err != nil {
		http.Error(w, "Forbidden: Invalid token \n" + err.Error(), http.StatusForbidden)
		return
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		// Ensure the token was issued for the specific file being requested.
		if claims.Subject != filePath {
			http.Error(w, "Forbidden: Token does not match file path", http.StatusForbidden)
			return
		}

		// Set the Content-Disposition header to ensure the browser downloads the file with its original name.
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

		// Serve the file from the filesystem.
		http.ServeFile(w, r, filepath.Join("./uploads", filePath))
	} else {
		http.Error(w, "Forbidden: Invalid claims", http.StatusForbidden)
	}
}
