package utils

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func RemoveImage(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return err
	}
	return nil
}

func UploadImage(base664String string, dir string, uniqueID string) (*string, error) {
	// define parameters
	rootFolder := "cdn"
	targetFolder := fmt.Sprintf("%s/%s", rootFolder, dir)
	allowedExtensions := map[string]any{
		"image/png":  ".png",
		"image/jpeg": ".jpg",
		"image/jpg":  ".jpg",
		"image/webp": ".webp",
	}
	maxSize := 2 * 1024 * 1024 // 2MB

	// validating data
	mimeType, data, err := parseBase64(base664String)
	if err != nil {
		return nil, err
	}

	ext, ok := allowedExtensions[mimeType]
	if !ok {
		return nil, errors.New("unsupported image format")
	}

	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, errors.New("invalid base64 data")
	}

	if len(decoded) > maxSize {
		return nil, errors.New("image too large (max 2MB)")
	}

	// ensure folder exists
	if err := os.MkdirAll(targetFolder, 0755); err != nil {
		return nil, errors.New("failed to create folder")
	}

	// save image
	filename := fmt.Sprintf("%s%s", uniqueID, ext)
	fullPath := filepath.Join(targetFolder, filename)
	if err := os.WriteFile(fullPath, decoded, 0644); err != nil {
		return nil, errors.New("failed to save image")
	}

	return &fullPath, nil
}

func parseBase64(b64 string) (mimeType, data string, err error) {
	if !strings.HasPrefix(b64, "data:") {
		return "", "", errors.New("Base64 string must contain MIME type prefix")
	}

	parts := strings.SplitN(b64, ",", 2)
	if len(parts) != 2 {
		return "", "", errors.New("invalid base64 format")
	}

	header := parts[0]
	data = parts[1]

	if !strings.Contains(header, ";base64") {
		return "", "", errors.New("base64 string must contain ';base64'")
	}

	mimeType = strings.TrimPrefix(header, "data:")
	mimeType = strings.TrimSuffix(mimeType, ";base64")

	return mimeType, data, nil
}
