// This module is cross-platform (Windows and Linux).

package file_manipulation

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Protheophage/GO/pkg/random_utilities"
)

// RemoveFiles deletes files matching specific criteria.
//
// Description:
// - Searches for files matching a pattern and deletes them.
// - Can search all drives or a specific directory.
//
// Parameters:
// - filesToDelete (string): The pattern of files to delete (e.g., "*.tmp").
// - searchAllDrives (bool): Whether to search all drives or a specific directory.
// - checkThisDisk (string): The specific directory or drive to search (ignored if searchAllDrives is true).
//
// Returns:
// - error: An error if the operation fails.
//
// Example Usage:
// ```go
// err := RemoveFiles("*.tmp", true, "")
//
//	if err != nil {
//	    fmt.Println("Error:", err)
//	} else {
//
//	    fmt.Println("Files removed successfully.")
//	}
//
// ```
func RemoveFiles(filesToDelete string, searchAllDrives bool, checkThisDisk string) error {
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
		fmt.Printf("Removing files in: %s\n", drive)
		err := filepath.Walk(drive, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil // Skip errors
			}
			if match, _ := filepath.Match(filesToDelete, info.Name()); match {
				if err := os.Remove(path); err != nil {
					fmt.Printf("Failed to remove file: %s. Error: %v\n", path, err)
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
