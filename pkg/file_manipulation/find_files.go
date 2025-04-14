// This module is cross-platform (Windows and Linux).

package file_manipulation

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/Protheophage/GO/pkg/random_utilities"
)

// FindFiles searches for files based on a pattern.
//
// Description:
// - Searches for files matching a pattern in all drives or a specific directory.
// - Returns a list of matching file paths.
//
// Parameters:
// - filesToFind (string): The pattern of files to find (e.g., "*.txt").
// - searchAllDrives (bool): Whether to search all drives or a specific directory.
// - checkThisDisk (string): The specific directory or drive to search (ignored if searchAllDrives is true).
//
// Returns:
// - []string: A slice of matching file paths.
// - error: An error if the operation fails.
//
// Example Usage:
// ```go
// files, err := FindFiles("*.txt", false, "C:\\")
//
//	if err != nil {
//	    fmt.Println("Error:", err)
//	} else {
//
//	    fmt.Println("Found files:", files)
//	}
//
// ```
func FindFiles(filesToFind string, searchAllDrives bool, checkThisDisk string) ([]string, error) {
	var files []string

	if searchAllDrives {
		// Get all drives on the system (Windows-specific logic)
		drives := random_utilities.GetAllDrives()
		for _, drive := range drives {
			fmt.Printf("Searching: %s for %s\n", drive, filesToFind)
			err := filepath.Walk(drive, func(path string, info fs.FileInfo, err error) error {
				if err != nil {
					return nil // Skip errors
				}
				if match, _ := filepath.Match(filesToFind, info.Name()); match {
					files = append(files, path)
				}
				return nil
			})
			if err != nil {
				fmt.Printf("Error searching drive %s: %v\n", drive, err)
			}
		}
	} else {
		if checkThisDisk == "" {
			if os.PathSeparator == '/' {
				checkThisDisk = "/"
			} else {
				checkThisDisk = filepath.Join(os.Getenv("SystemDrive"), "")
			}
		}
		fmt.Printf("Searching: %s for %s\n", checkThisDisk, filesToFind)
		err := filepath.Walk(checkThisDisk, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return nil // Skip errors
			}
			if match, _ := filepath.Match(filesToFind, info.Name()); match {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("error searching directory %s: %v", checkThisDisk, err)
		}
	}

	return files, nil
}
