package azure

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/google/uuid"
)

var blobClient *azblob.Client
var storageAccountName string
var containerName string
var cdnEndpoint string

func InitStorage() error {
	storageAccountName = os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
	containerName = os.Getenv("AZURE_STORAGE_CONTAINER_NAME")
	cdnEndpoint = os.Getenv("AZURE_CDN_ENDPOINT")

	if storageAccountName == "" || containerName == "" || cdnEndpoint == "" {
		return fmt.Errorf("AZURE_STORAGE_ACCOUNT_NAME, AZURE_STORAGE_CONTAINER_NAME, and AZURE_CDN_ENDPOINT must be set")
	}

	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", storageAccountName)

	var err error
	blobClient, err = azblob.NewClient(serviceURL, cred, nil)
	if err != nil {
		return fmt.Errorf("failed to create blob client: %v", err)
	}

	return nil
}

// UploadClubLogo uploads a club logo to Azure Blob Storage and returns the CDN URL
func UploadClubLogo(clubID string, file multipart.File, header *multipart.FileHeader) (string, error) {
	if blobClient == nil {
		return "", fmt.Errorf("blob client not initialized")
	}

	// Validate file type
	contentType := header.Header.Get("Content-Type")
	if !isValidImageType(contentType) {
		return "", fmt.Errorf("invalid file type: only PNG, JPEG, and WebP images are allowed")
	}

	// Validate file size (max 5MB)
	if header.Size > 5*1024*1024 {
		return "", fmt.Errorf("file too large: maximum size is 5MB")
	}

	// Generate unique blob name
	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = getExtensionFromContentType(contentType)
	}
	blobName := fmt.Sprintf("%s-%s%s", clubID, uuid.New().String(), ext)

	// Read file content
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	// Upload to blob storage
	ctx := context.Background()
	_, err = blobClient.UploadBuffer(ctx, containerName, blobName, fileContent, &azblob.UploadBufferOptions{
		HTTPHeaders: &blob.HTTPHeaders{
			BlobContentType: &contentType,
		},
		Metadata: map[string]*string{
			"club_id":     &clubID,
			"uploaded_at": func() *string { t := time.Now().UTC().Format(time.RFC3339); return &t }(),
		},
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload file: %v", err)
	}

	// Return CDN URL
	cdnURL := fmt.Sprintf("%s/club-assets/%s", strings.TrimSuffix(cdnEndpoint, "/"), blobName)
	return cdnURL, nil
}

// DeleteClubLogo deletes a club logo from Azure Blob Storage
func DeleteClubLogo(logoURL string) error {
	if blobClient == nil {
		return fmt.Errorf("blob client not initialized")
	}

	// Extract blob name from CDN URL
	blobName := extractBlobNameFromURL(logoURL)
	if blobName == "" {
		return fmt.Errorf("invalid logo URL")
	}

	ctx := context.Background()
	_, err := blobClient.DeleteBlob(ctx, containerName, blobName, nil)
	if err != nil {
		return fmt.Errorf("failed to delete blob: %v", err)
	}

	return nil
}

func isValidImageType(contentType string) bool {
	validTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/webp",
	}

	for _, validType := range validTypes {
		if contentType == validType {
			return true
		}
	}
	return false
}

func getExtensionFromContentType(contentType string) string {
	switch contentType {
	case "image/jpeg", "image/jpg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	default:
		return ".jpg"
	}
}

func extractBlobNameFromURL(url string) string {
	// Extract blob name from CDN URL
	// Assuming CDN URL format: https://cdn-endpoint.com/club-assets/blob-name
	parts := strings.Split(url, "/")
	if len(parts) < 2 {
		return ""
	}

	// Find the part that contains "club-assets"
	for i, part := range parts {
		if part == "club-assets" && i+1 < len(parts) {
			return parts[i+1]
		}
	}

	return ""
}
