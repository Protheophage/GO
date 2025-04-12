package file_manipulation

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SetFilesExtension changes the extension of files matching specific criteria.
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
		drives = getAllDrives()
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