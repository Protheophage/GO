// This module is cross-platform (Windows and Linux).

package file_manipulation

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Protheophage/GO/pkg/random_utilities"
)

// SetFilesExtension changes the extension of files matching specific criteria.
//
// Description:
// - Searches for files matching a pattern and changes their extensions.
// - Can search all drives or a specific directory.
//
// Parameters:
// - filesToFind (string): The pattern of files to find (e.g., "*.txt").
// - newExtension (string): The new extension to apply (e.g., ".log").
// - searchAllDrives (bool): Whether to search all drives or a specific directory.
// - checkThisDisk (string): The specific directory or drive to search (ignored if searchAllDrives is true).
//
// Returns:
// - error: An error if the operation fails.
//
// Example Usage:
// ```go
// err := SetFilesExtension("*.txt", ".log", false, "C:\\")
//
//	if err != nil {
//	    fmt.Println("Error:", err)
//	} else {
//
//	    fmt.Println("File extensions updated successfully.")
//	}
//
// ```
func SetFilesExtension(filesToFind, newExtension string, searchAllDrives bool, checkThisDisk string) error {
	if checkThisDisk == "" {
		if os.PathSeparator == '/' {
			checkThisDisk = "/"
		} else {
			checkThisDisk = filepath.Join(os.Getenv("SystemDrive"), "")
		}
	}

	var drives []string
	if searchAllDrives {
		drives = random_utilities.GetAllDrives()
	} else {
		drives = []string{checkThisDisk}
	}

	for _, drive := range drives {
		fmt.Printf("Changing file extensions on: %s\n", drive)
		err := filepath.Walk(drive, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil // Skip errors
			}
			if match, _ := filepath.Match(filesToFind, info.Name()); match {
				newPath := strings.TrimSuffix(path, filepath.Ext(path)) + newExtension
				if err := os.Rename(path, newPath); err != nil {
					fmt.Printf("Failed to rename file: %s. Error: %v\n", path, err)
				}
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("error processing drive %s: %v", drive, err)
		}
	}

	return nil
}
