// This module is cross-platform (Windows and Linux).

package file_manipulation

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Protheophage/GO/pkg/random_utilities"
)

// FindFilesByContent searches for files containing specific content.
//
// Description:
// - Searches for files of a specific type and size containing a given string.
// - Can search all drives or a specific directory.
//
// Parameters:
// - stringToFind (string): The string to search for within files.
// - fileTypeToSearch (string): The file extension to filter by (e.g., ".txt").
// - maxFileSizeKB (int): The maximum file size in kilobytes to search.
// - searchAllDrives (bool): Whether to search all drives or a specific directory.
// - checkThisDisk (string): The specific directory or drive to search (ignored if searchAllDrives is true).
//
// Returns:
// - []string: A slice of file paths containing the specified content.
// - error: An error if the operation fails.
//
// Example Usage:
// ```go
// files, err := FindFilesByContent("error", ".log", 1024, false, "C:\\")
//
//	if err != nil {
//	    fmt.Println("Error:", err)
//	} else {
//
//	    fmt.Println("Found files containing the string:", files)
//	}
//
// ```
func FindFilesByContent(stringToFind, fileTypeToSearch string, maxFileSizeKB int, searchAllDrives bool, checkThisDisk string) ([]string, error) {
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

	var foundFiles []string
	for _, drive := range drives {
		fmt.Printf("Searching for content in: %s\n", drive)
		err := filepath.Walk(drive, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil // Skip errors
			}
			if filepath.Ext(path) == fileTypeToSearch && info.Size() <= int64(maxFileSizeKB*1024) {
				file, err := os.Open(path)
				if err != nil {
					return nil
				}
				defer file.Close()

				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					if strings.Contains(scanner.Text(), stringToFind) {
						foundFiles = append(foundFiles, path)
						break
					}
				}
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("error processing drive %s: %v", drive, err)
		}
	}

	return foundFiles, nil
}
