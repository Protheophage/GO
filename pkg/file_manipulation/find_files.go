package file_manipulation

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// FindFiles searches for files based on a pattern. It can search all drives or a specific directory.
func FindFiles(filesToFind string, searchAllDrives bool, checkThisDisk string) ([]string, error) {
	var files []string

	if searchAllDrives {
		// Get all drives on the system (Windows-specific logic)
		drives := getAllDrives()
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

// getAllDrives returns a list of all drives on the system.
func getAllDrives() []string {
	if os.PathSeparator == '/' {
		return []string{"/"} // Root directory for Linux
	}
	drives := []string{}
	for letter := 'A'; letter <= 'Z'; letter++ {
		drive := fmt.Sprintf("%c:\\", letter)
		if _, err := os.Stat(drive); err == nil {
			drives = append(drives, drive)
		}
	}
	return drives
}