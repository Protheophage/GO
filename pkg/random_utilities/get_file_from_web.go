// This module is cross-platform (Windows and Linux).

package random_utilities

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

// GetFileFromWeb downloads a file from the specified URL and saves it to the provided destination path.
// If overwrite is false, it will not overwrite an existing file.
func GetFileFromWeb(sourceURL, destinationPath string, overwrite bool) error {
	// Normalize the destination path for cross-platform compatibility
	destinationPath = filepath.Clean(destinationPath)

	// Ensure the destination directory exists
	destDir := filepath.Dir(destinationPath)
	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", destDir, err)
		}
	}

	// Check if the file already exists
	if !overwrite {
		if _, err := os.Stat(destinationPath); err == nil {
			return errors.New("file already exists and overwrite is disabled")
		}
	}

	// Download the file
	resp, err := http.Get(sourceURL)
	if err != nil {
		return fmt.Errorf("failed to download file from %s: %w", sourceURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: received status code %d", resp.StatusCode)
	}

	// Create the destination file with appropriate permissions
	outFile, err := os.Create(destinationPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", destinationPath, err)
	}
	defer outFile.Close()

	// Copy the content to the destination file
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save file to %s: %w", destinationPath, err)
	}

	// Print success message with platform-specific newline
	newline := "\n"
	if runtime.GOOS == "windows" {
		newline = "\r\n"
	}
	fmt.Printf("Download successful. The file is located at %s%s", destinationPath, newline)

	return nil
}
