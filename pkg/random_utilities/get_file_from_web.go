// This module is cross-platform (Windows and Linux).

package random_utilities

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

// GetFileFromWeb downloads a file from the specified URL and saves it to the provided destination path.
//
// Description:
// - Ensures the destination directory exists before downloading.
// - Optionally overwrites the file if it already exists.
//
// Parameters:
// - sourceURL (string): The URL of the file to download.
// - destinationPath (string): The path where the file will be saved.
// - overwrite (bool): Whether to overwrite the file if it already exists.
//
// Returns:
// - error: An error if the download or file saving fails.
//
// Example Usage:
// ```go
// err := GetFileFromWeb("https://example.com/file.txt", "C:\\Downloads\\file.txt", true)
//
//	if err != nil {
//	    fmt.Println("Error:", err)
//	} else {
//
//	    fmt.Println("File downloaded successfully.")
//	}
//
// ```
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
			fmt.Printf("File %s already exists. Overwrite is disabled. Skipping download.\n", destinationPath)
			return nil
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
