package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func GetLastMonthLogFilePath() (string, error) {
	// Get Month
	lastMonth := time.Now().AddDate(0, -1, 0)
	fileName := fmt.Sprintf("wardrobe-%s-%d.log", lastMonth.Format("January"), lastMonth.Year())

	filePath := filepath.Join("logs", fileName)

	// Check Exist
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("log file not found: %s", filePath)
	}

	return filePath, nil
}

func DeleteFileByPath(path string) error {
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("failed to delete file %s: %w", path, err)
	}
	return nil
}
