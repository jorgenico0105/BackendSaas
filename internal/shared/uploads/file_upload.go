package uploads

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	MaxFileSize     = 10 << 20 // 10 MB
	UploadDirectory = "storage/uploads"
)

var AllowedImageTypes = []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
var AllowedDocTypes = []string{".pdf", ".doc", ".docx", ".xls", ".xlsx"}

type UploadResult struct {
	FileName     string `json:"file_name"`
	OriginalName string `json:"original_name"`
	FilePath     string `json:"file_path"`
	FileSize     int64  `json:"file_size"`
	MimeType     string `json:"mime_type"`
}

func SaveFile(c *gin.Context, file *multipart.FileHeader, subDirectory string, allowedTypes []string) (*UploadResult, error) {
	ext := strings.ToLower(filepath.Ext(file.Filename))

	if !isAllowedType(ext, allowedTypes) {
		return nil, fmt.Errorf("file type %s is not allowed", ext)
	}

	if file.Size > MaxFileSize {
		return nil, fmt.Errorf("file size exceeds maximum allowed size of %d bytes", MaxFileSize)
	}

	uploadPath := filepath.Join(UploadDirectory, subDirectory)
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	newFileName := generateFileName(ext)
	fullPath := filepath.Join(uploadPath, newFileName)

	if err := c.SaveUploadedFile(file, fullPath); err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	return &UploadResult{
		FileName:     newFileName,
		OriginalName: file.Filename,
		FilePath:     fullPath,
		FileSize:     file.Size,
		MimeType:     file.Header.Get("Content-Type"),
	}, nil
}

func DeleteFile(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil
	}
	return os.Remove(filePath)
}

func generateFileName(ext string) string {
	return fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), generateRandomString(8), ext)
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[time.Now().UnixNano()%int64(len(charset))]
		time.Sleep(time.Nanosecond)
	}
	return string(result)
}

func isAllowedType(ext string, allowedTypes []string) bool {
	for _, allowed := range allowedTypes {
		if ext == allowed {
			return true
		}
	}
	return false
}
