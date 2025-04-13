package file_manipulation

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Protheophage/GO/pkg/random_utilities"
)

// GetFilesCount counts the number of files matching specific criteria.
func GetFilesCount(filesToFind string, searchAllDrives bool, checkThisDisk string) (int, error) {
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

	count := 0
	for _, drive := range drives {
		fmt.Printf("Counting files in: %s\n", drive)
		err := filepath.Walk(drive, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil // Skip errors
			}
			if match, _ := filepath.Match(filesToFind, info.Name()); match {
				count++
			}
			return nil
		})
		if err != nil {
			return 0, fmt.Errorf("error processing drive %s: %v", drive, err)
		}
	}

	return count, nil
}